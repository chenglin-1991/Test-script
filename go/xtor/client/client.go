package client

import (
	"crypto/tls"
	"net/http"
	. "github.com/xtao/xtor/common"
 )

type XtorClient struct {
	httpClient *http.Client
	account *UserAccountInfo
	server string
}

func NewXtorClient(server string, account *UserAccountInfo) *XtorClient {
	httpClient := &http.Client{}
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient.Transport = tr

	return &XtorClient {
		httpClient: httpClient,
		server: server,
		account: account,
	}
}

