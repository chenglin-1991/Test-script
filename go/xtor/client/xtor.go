package client

import (
	"errors"
	"net/http"
	"fmt"
	"strings"
	. "github.com/xtao/xtor/common"
	"encoding/json"
)

func (client *XtorClient)Du(vol string, rpath string, obj bool) (string, error) {

	server := client.server

	duReq := XTXtorDuAPIRequest {
		VolName: vol,
		Path: rpath,
		Obj: obj,
	}

	url := fmt.Sprintf("http://%s/v1/du", server)

	js, err := json.Marshal(duReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", url, strings.NewReader(string(js)))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	err = SignXtorReqHeader(req, client.account)
	if err != nil {
		return "", err
	}

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return "", err
	}


	decoder := json.NewDecoder(resp.Body)
	var result XTXtorAPIDuResult
	err = decoder.Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Status == XT_API_RET_OK {
		var output string
		if obj != true {
			output = result.Result.Size
		} else {
			output = fmt.Sprintf("%s files\n%s dirs",
				result.Result.Files, result.Result.Dirs)
		}
		return output, nil
	} else {
		return "", errors.New(result.Errmsg)
	}
}

func (client *XtorClient)Fsstat(vol string) (*XTXtorFsstatReply, error) {

	server := client.server

	dfReq := XTXtorAPIRequest {
		VolName: vol,
	}

	url := fmt.Sprintf("http://%s/v1/fsstat", server)

	js, err := json.Marshal(dfReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, strings.NewReader(string(js)))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	err = SignXtorReqHeader(req, client.account)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var result XTXtorAPIFsstatResult
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Status != XT_API_RET_OK {
		return nil, errors.New(result.Errmsg)
	}

	return &result.Result, nil
}
