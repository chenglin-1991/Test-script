package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

const urls = "http://192.168.60.228:8888/seed/var/log/glusterfs/cli.log"

func main() {
    resp, err := http.Get(urls)
    ErrPrint(err)
    defer resp.Body.Close()

    //[<a href="/seed/var/log/samba/old/">old/</a>]
    //[<a href="/seed/var/log/samba/log.10.61.52.214">log.10.61.52.214</a>]
    //<a href="/seed/var/log/samba/glusterfs/">glusterfs/</a>
    //dirreg := regexp.MustCompile(`^<a href=".+/">.+</a>`)
    //filereg := regexp.MustCompile(`^<a href=".+">.+</a>`)
    //
    //buff := bufio.NewReader(resp.Body)
    //for {
    //    data, _, eof := buff.ReadLine()
    //    if eof == io.EOF {
    //        break
    //    }
    //
    //    str1 := string(data)
    //
    //    result1 := reg.FindAllString(str1, -1)
    //    fmt.Printf("%v\n", result1)
    //
    //    fmt.Println("=========",result1,len(result1))
    //}

    fmt.Println(resp.Header.Get("Content-Type"))
    fmt.Println(resp.Header.Get("Content-Disposition"))
    fmt.Println(resp.Header.Get("Content-Length"))
}

func ErrPrint(err error) {
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
}
