/* 
 Copyright (c) 2016-2017 XTAO technology <www.xtaotech.com>
 All rights reserved.

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions
 are met:
  1. Redistributions of source code must retain the above copyright
     notice, this list of conditions and the following disclaimer.
  2. Redistributions in binary form must reproduce the above copyright
     notice, this list of conditions and the following disclaimer in the
     documentation and/or other materials provided with the distribution.
 
  THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
  ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
  OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
  HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
  LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
  OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
  SUCH DAMAGE.
*/

package handler

import "C"
import (
	"encoding/json"
	"fmt"
	. "github.com/xtao/xtor/common"
	. "github.com/xtao/xtor/server/manager"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Du(w http.ResponseWriter, req *http.Request) {

	var result XTXtorAPIDuResult

	err, _ := AuthHeaderFromClient(req)
	if err != nil {
		result.Status = XT_API_RET_FAIL
		result.Errmsg = "Authorization failure"
		WriteJSON(500, result, w)
		return
	}

	decoder := json.NewDecoder(req.Body)

	defer req.Body.Close()

	var duReq XTXtorDuAPIRequest

	err = decoder.Decode(&duReq)
	if err != nil {
		return
	}

	mgr := GetXtorMgr()

	err, reply := mgr.Du(duReq)
	if err != nil {
		Logger.Printf("Can't handle du: %s\n", err.Error())
		result.Status = XT_API_RET_FAIL
		result.Errmsg = err.Error()
		WriteJSON(500, result, w)
	} else {
		Logger.Printf("Successfully get du\n")
		result.Status = XT_API_RET_OK
		result.Result = *reply
		WriteJSON(http.StatusAccepted, result, w)
	}
}

func Fsstat(w http.ResponseWriter, req *http.Request) {

	var result XTXtorAPIFsstatResult

	err, _ := AuthHeaderFromClient(req)
	if err != nil {
		result.Status = XT_API_RET_FAIL
		result.Errmsg = "Authorization failure"
		WriteJSON(500, result, w)
		return
	}

	decoder := json.NewDecoder(req.Body)

	defer req.Body.Close()

	var dfReq XTXtorAPIRequest

	err = decoder.Decode(&dfReq)
	if err != nil {
		return
	}

	mgr := GetXtorMgr()

	err, reply := mgr.Fsstat(dfReq)
	if err != nil {
		Logger.Printf("Can't handle fsstat: %s\n", err.Error())
		result.Status = XT_API_RET_FAIL
		result.Errmsg = err.Error()
		WriteJSON(500, result, w)
	} else {
		Logger.Printf("Successfully get fsstat\n")
		result.Status = XT_API_RET_OK
		result.Result = *reply
		WriteJSON(http.StatusAccepted, result, w)
	}
}

func QuotaList(w http.ResponseWriter, req *http.Request) {

	var result XTXtorAPIReply

	err, _ := AuthHeaderFromClient(req)
	if err != nil {
		result.Status = -1
		result.Errmsg = "Authorization failure"
		WriteJSON(500, result, w)
		return
	}

	decoder := json.NewDecoder(req.Body)

	defer req.Body.Close()

	var apiReq XTXtorQuotaListAPIRequest

	err = decoder.Decode(&apiReq)
	if err != nil {
		return
	}

	mgr := GetXtorMgr()

	err, reply := mgr.QuotaList(apiReq)
	if err != nil {
		Logger.Printf("Can't handle quota list: %s\n", err.Error())
		result.Status = -1
		result.Errmsg = err.Error()
		WriteJSON(500, result, w)
	} else {
		Logger.Printf("Successfully quota list\n")
		WriteJSON(http.StatusAccepted, reply, w)
	}
}

func QuotaSet(w http.ResponseWriter, req *http.Request) {

	var result XTXtorAPIReply

	err, _ := AuthHeaderFromClient(req)
	if err != nil {
		result.Status = -1
		result.Errmsg = "Authorization failure"
		WriteJSON(500, result, w)
		return
	}

	decoder := json.NewDecoder(req.Body)

	defer req.Body.Close()

	var apiReq XTXtorQuotaSetAPIRequest

	err = decoder.Decode(&apiReq)
	if err != nil {
		return
	}

	mgr := GetXtorMgr()

	err, reply := mgr.QuotaSet(apiReq)
	if err != nil {
		Logger.Printf("Can't handle quota set: %s\n", err.Error())
		result.Status = -1
		result.Errmsg = err.Error()
		WriteJSON(500, result, w)
	} else {
		Logger.Printf("Successfully quota set\n")
		WriteJSON(http.StatusAccepted, reply, w)
	}
}

func ShowSource(w http.ResponseWriter, r *http.Request) {
	var UrlPath = ""
	UrlPath = r.URL.Path

	s := strings.Split(r.URL.Path, "/")
	host := s[3]

	Logger.Println(fmt.Sprintf("get req for %s", UrlPath))

	Cf := GetXtorConfig()

	//show host
	if len(s) == 4 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := fmt.Fprintf(w, "<pre>\n")
		if err != nil {
			Logger.Println(fmt.Sprintf("failed to write respond for %s", UrlPath))
			goto out
		}
		for _, v := range Cf.LogPath {
			if v[len(v)-1] == '/' {
				v = v[:len(v)-1]
			}

			vurl := url.URL{Path: fmt.Sprintf("%s%s/", UrlPath, v)}
			_, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), v)
		}
		_, _ = fmt.Fprintf(w, "</pre>\n")
	} else {
		var path string
		s = s[4:]
		for _, v := range s {
			path = path + "/" + v
		}

		hostname, _ := os.Hostname()
		if hostname == host {
			f, err := os.Open(path)
			if err != nil {
				Cf.Error(w, Cf.ToHttpError(err))
				goto out
			}
			defer f.Close()

			stat, err := f.Stat()
			if err != nil {
				Cf.Error(w, Cf.ToHttpError(err))
				goto out
			}

			if stat.IsDir() {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = fmt.Fprintf(w, "<pre>\n")
				Cf.DirList(w, f, UrlPath)
				_, _ = fmt.Fprintf(w, "</pre>\n")
			} else {
				filename := filepath.Base(path)
				_, contentType := Cf.GetContentType(filename)

				w.Header().Set("Content-Disposition", "attachment; filename="+filename)

				w.Header().Set("Content-Type", contentType)

				w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

				_, _ = f.Seek(0, 0)

				_, _ = io.Copy(w, f)
			}
		} else {
			Cf.ClientDownload(w, host, UrlPath)
		}
	}

out:
	return
}

func ShowRootPath(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	Cf := GetXtorConfig()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<pre>\n")
	for _, n := range Cf.Node {
		vurl := url.URL{Path: fmt.Sprintf("%s/%s/", path, n)}
		_, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), n)
	}
	_, _ = fmt.Fprintf(w, "</pre>\n")
}
