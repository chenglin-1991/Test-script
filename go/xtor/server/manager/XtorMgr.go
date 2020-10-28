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

package manager

import (
    "encoding/json"
    "errors"
    "fmt"
    . "github.com/xtao/xtor/common"
    "gopkg.in/ini.v1"
    "io"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "sort"
    "strings"
)

type XtorMgr struct{}

const (
    XTOR_XD_ALAMO_EXEC string = "xd-alamo"
    XtorConfigPath            = "/etc/xtorsvr/xtorsvr.conf"
    XtorLogPath               = "/var/log/xtorsvr/xtorsvr.log"
)

var GlobalXtorMgr *XtorMgr = nil

func NewXtorMgr() *XtorMgr {
    XtorMgr := &XtorMgr{}
    GlobalXtorMgr = XtorMgr

    return XtorMgr
}

func GetXtorMgr() *XtorMgr {
    return GlobalXtorMgr
}

var GlobleConfig *XtorConfig = nil

func NewXtorConfig() *XtorConfig {
    XtorConfig := &XtorConfig{}
    GlobleConfig = XtorConfig

    return GlobleConfig
}

func GetXtorConfig() *XtorConfig {
    return GlobleConfig
}

type QuotaOutput struct {
    Available  string `json:"Available"`
    Path       string `json:"Path"`
    Used       string `json:"Used"`
    Dirs       string `json:"Dirs"`
    Files      string `json:"Files"`
    Hard_Limit string `json: "Hard_Limit"`
    Soft_Limit string `json: "Soft_Limit"`
}

type XTXtorQuota struct {
    Output []QuotaOutput `json:"output"`
    Errmsg string        `json:"errmsg"`
    Error  int           `json:"error"`
}

type XTXtorXdOutput struct {
    Output string `json:"output"`
    Errmsg string `json:"errmsg"`
    Error  int    `json:"error"`
}

func (a *XtorMgr) Du(req XTXtorDuAPIRequest) (error, *XTXtorDuReply) {

    var cmd *exec.Cmd

    if req.Obj == false {
        cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
            req.VolName, "-p", req.Path)
    } else {
        cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
            req.VolName, "-o", "-p", req.Path)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Start(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    var result XTXtorQuota
    decoder := json.NewDecoder(stdout)
    err = decoder.Decode(&result)
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Wait(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if result.Error != 0 {
        return errors.New(result.Errmsg), nil
    }

    reply := &XTXtorDuReply{
        Size:  result.Output[0].Used,
        Dirs:  result.Output[0].Dirs,
        Files: result.Output[0].Files,
    }

    return nil, reply
}

type XTXtorFsstat struct {
    Output XTXtorFsstatReply `json:"output"`
    Errmsg string            `json:"errmsg"`
    Error  int               `json:"error"`
}

func (a *XtorMgr) Fsstat(req XTXtorAPIRequest) (error, *XTXtorFsstatReply) {

    var cmd *exec.Cmd

    cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "volume", "stat",
        req.VolName)

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Start(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    var result XTXtorFsstat
    decoder := json.NewDecoder(stdout)
    err = decoder.Decode(&result)
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Wait(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    Logger.Print(result)

    if result.Error != 0 {
        return errors.New(result.Errmsg), nil
    }

    return nil, &result.Output
}

func (a *XtorMgr) QuotaList(req XTXtorQuotaListAPIRequest) (error, *XTXtorAPIReply) {

    var cmd *exec.Cmd

    if req.Path == "" {
        if req.Obj == false {
            cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
                req.VolName)
            Logger.Println("quota list ", req.VolName)
        } else {
            cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
                "-o", req.VolName)
            Logger.Printf("quota list -o %s\n", req.VolName)
        }
    } else {
        if req.Obj == false {
            cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
                req.VolName, "-p", req.Path)
            Logger.Printf("quota list %s -p %s", req.VolName, req.Path)
        } else {
            cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "list",
                req.VolName, "-o", "-p", req.Path)
        }
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Start(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    var result XTXtorQuota
    decoder := json.NewDecoder(stdout)
    err = decoder.Decode(&result)
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Wait(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if result.Error != 0 {
        return errors.New(result.Errmsg), nil
    }

    output, err := json.Marshal(result.Output)
    reply := &XTXtorAPIReply{
        Status: result.Error,
        Errmsg: result.Errmsg,
        Result: string(output),
    }

    return nil, reply
}

func (a *XtorMgr) QuotaSet(req XTXtorQuotaSetAPIRequest) (error, *XTXtorAPIReply) {

    var cmd *exec.Cmd

    if req.Obj == false {
        cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "set",
            req.VolName, req.Path, req.Limit)
    } else {
        cmd = exec.Command(XTOR_XD_ALAMO_EXEC, "quota", "set",
            "-o", req.VolName, req.Path, req.Limit)
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Start(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    var result XTXtorQuota
    decoder := json.NewDecoder(stdout)
    err = decoder.Decode(&result)
    if err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if err := cmd.Wait(); err != nil {
        Logger.Println(err.Error())
        return err, nil
    }

    if result.Error != 0 {
        return errors.New(result.Errmsg), nil
    }

    output, err := json.Marshal(result.Output)

    reply := &XTXtorAPIReply{
        Status: result.Error,
        Errmsg: result.Errmsg,
        Result: string(output),
    }

    return nil, reply
}

type XtorConfig struct {
    Node    []string `ini:"node" comment:"which inode you want to monitor"`
    LogPath []string `ini:"logpath" comment:"which log path you want to download"`
    Port    int      `ini:"port" comment:"the http server listen port"`
}

func (p *XtorConfig) UnPraiseConfigFile() {
    cfg, err := ini.Load(XtorConfigPath)
    if err != nil {
        Logger.Println(fmt.Sprintf("failed to Load config file %s", err))
        os.Exit(-1)
    }

    err = cfg.MapTo(p)
    if err != nil {
        fmt.Println(err)
        Logger.Println(fmt.Sprintf("failed to Map config file %s", err))
        os.Exit(-1)
    }
}

func (p *XtorConfig) ClientDownload(w http.ResponseWriter, HostName string,
    UrlPath string) {
    var IsDir bool

    uu := fmt.Sprintf("http://%s:%d%s", HostName, p.Port, UrlPath)

    b := strings.HasSuffix(UrlPath, "/")
    if b == true {
        IsDir = true
    }

    resp, err := http.Get(uu)
    if err != nil {
        Logger.Println(fmt.Sprintf("failed to get %s %s", uu, err))
        p.Error(w, http.StatusInternalServerError)
        return
    }

    defer func() { _ = resp.Body.Close() }()

    if IsDir == true {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
    } else {
        w.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))

        w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

        w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
    }
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        Logger.Println(fmt.Sprintf("failed to write "+
            "%s respond %s", uu, err))
        p.Error(w, http.StatusInternalServerError)
    }
}

func (p *XtorConfig) Error(w http.ResponseWriter, code int) {
    w.WriteHeader(code)
}

func (p *XtorConfig) DirList(w http.ResponseWriter, f http.File, UrlPath string) {
    var vurl url.URL

    dirs, err := f.Readdir(-1)
    if err != nil {
        p.Error(w, http.StatusInternalServerError)
        return
    }
    sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    _, _ = fmt.Fprintf(w, "<pre>\n")
    for _, d := range dirs {
        name := d.Name()
        if d.IsDir() {
            name += "/"
        }

        vurl = url.URL{Path: fmt.Sprintf("%s/%s", UrlPath, name)}
        _, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), name)
    }
    _, _ = fmt.Fprintf(w, "</pre>\n")
}

func (p *XtorConfig) ToHttpError(err error) int {
    if os.IsNotExist(err) {
        return http.StatusNotFound
    }

    return http.StatusInternalServerError
}

func (p *XtorConfig) GetContentType(fileName string) (extension, contentType string) {
    arr := strings.Split(fileName, ".")

    // to see: https://tool.oschina.net/commons/
    if len(arr) >= 2 {
        extension = arr[len(arr)-1]
        switch extension {
        case "jpeg", "jpe", "jpg":
            contentType = "image/jpeg"
        case "png":
            contentType = "image/png"
        case "gif":
            contentType = "image/gif"
        case "mp4":
            contentType = "video/mpeg4"
        case "mp3":
            contentType = "audio/mp3"
        case "wav":
            contentType = "audio/wav"
        case "pdf":
            contentType = "application/pdf"
        case "js":
            contentType = "application/javascript"
        case "xml":
            contentType = "text/xml"
        case "doc", "":
            contentType = "application/msword"
        case "html":
            contentType = "text/html"
        default:
            contentType = "application/octet-stream"
        }
    }

    return
}
