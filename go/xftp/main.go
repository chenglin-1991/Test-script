package main

import (
    "fmt"
    "net/http"
    "xftp/utils"
)

/*
https://github.com/gccostabr/go-bin-rpm/blob/master/demo/rpm.json
https://zhuanlan.zhihu.com/p/127353769
https://github.com/gccostabr/go-bin-rpm
*/

func main() {
    c := &utils.Config{}

    l := utils.LogInit(utils.XftpLogPath)

    utils.Rlog = l

    //c.GenerateTmpConfigFile(l)

    c.UnPraiseConfigFile(l)

    http.HandleFunc("/", c.ShowRootPath)

    for _, v := range c.Node {
        http.HandleFunc(fmt.Sprintf("/%s/", v), c.ShowSource)
    }

    port := fmt.Sprintf(":%d", c.Port)

    _ = http.ListenAndServe(port, nil)
}
