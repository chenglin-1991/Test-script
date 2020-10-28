package server

import (
    "net/http"
    . "github.com/xtao/xtor/common"
)

type XtRestServer struct {
    addr      string
}

var GlobalRestServer *XtRestServer = nil

func NewRESTServer(addr string) *XtRestServer {
    GlobalRestServer = &XtRestServer{
        addr: addr,
    }

    return GlobalRestServer
}

func (server *XtRestServer) StartRESTServer() {
    router := NewRouterFn()
    Logger.Fatal(http.ListenAndServe(server.addr, router))
}
