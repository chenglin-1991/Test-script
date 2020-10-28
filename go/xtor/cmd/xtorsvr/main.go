package main

import (
	"fmt"
    "os"
    . "github.com/xtao/xtor"
)


func main() {
    fmt.Println("Xtor server starts now")

    err := ServerStart()
    if err != nil {
        fmt.Printf("Fail to start Xtor server: %s\n",
            err.Error())
        os.Exit(1)
    }
}
