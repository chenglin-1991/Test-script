package main

import (
    "fmt"
    "gopkg.in/ini.v1"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

/*
   https://www.cnblogs.com/wind-zhou/p/12821176.html
   https://cloud.tencent.com/developer/article/1066126
   https://juejin.im/post/6844904048764649479
   https://github.com/go-ini/ini
   https://www.jianshu.com/p/f95558a49e98
   https://blog.csdn.net/xmcy001122/article/details/104654096
*/

type Config struct {
    Node    []string `ini:"node" comment:"which inode you want to monitor"`
    LogPath []string `ini:"logpath" comment:"which log path you want to download"`
    Port    int      `ini:"port" comment:"the http server listen port"`
}

var extensionToContentType = map[string]string{
    ".html": "text/html; charset=utf-8",
    ".css":  "text/css; charset=utf-8",
    ".js":   "application/javascript",
    ".xml":  "text/xml; charset=utf-8",
    ".jpg":  "image/jpeg",
}

func (p *Config) GenerateTmpConfigFile() {
    cfg := ini.Empty()

    c := Config{
        Node:    []string{"seed", "AlamoD1N8serverc802", "xt3"},
        LogPath: []string{"/var/log/glusterfs", "/var/log/samba"},
        Port:    8888,
    }

    err := ini.ReflectFrom(cfg, &c)
    if err != nil {
        fmt.Println("ReflectFrom failed: ", err)
        return
    }

    err = cfg.SaveTo("my-copy.ini")
    if err != nil {
        fmt.Println("SaveTo failed: ", err)
        return
    }
}

func (p *Config) UnPraiseConfigFile() {
    cfg, err := ini.Load("my-copy.ini")
    if err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }

    err = cfg.MapTo(p)
    if err != nil {
        fmt.Println(err)
        return
    }
}

func (p *Config) toHTTPError(err error) int {
    if os.IsNotExist(err) {
        return http.StatusNotFound
    }
    if os.IsPermission(err) {
        return http.StatusForbidden
    }
    return http.StatusInternalServerError
}

func (p *Config) Error(w http.ResponseWriter, code int) {
    w.WriteHeader(code)
}

func (p *Config) DirList(w http.ResponseWriter, f http.File, UrlPath string) {
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

        vurl = url.URL{Path: fmt.Sprintf("%s%s", UrlPath, name)}
        _, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), name)
    }
    _, _ = fmt.Fprintf(w, "</pre>\n")
}

func (p *Config) ShowRootPath(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    fmt.Println("=====2", path)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    _, _ = fmt.Fprintf(w, "<pre>\n")
    for _, n := range p.Node {
        n += "/"
        vurl := url.URL{Path: n}
        _, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), n)
    }
    _, _ = fmt.Fprintf(w, "</pre>\n")
}

func getContentType(fileName string) (extension, contentType string) {
    arr := strings.Split(fileName, ".")
    fmt.Println(arr, len(arr))
    fmt.Println(arr[len(arr)-1])

    // see: https://tool.oschina.net/commons/
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
        case "doc", "":
            contentType = "application/msword"
        default:
            contentType = "application/octet-stream"
        }
    }

    return
}

func (p *Config) ShowSource(w http.ResponseWriter, r *http.Request) {
    var UrlPath = ""
    UrlPath = r.URL.Path

    s := strings.Split(r.URL.Path, "/")
    host := s[1]

    fmt.Println(r.URL.Path, s, len(s))

    if len(s) == 3 {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        _, _ = fmt.Fprintf(w, "<pre>\n")
        for _, v := range p.LogPath {
            vurl := url.URL{Path: fmt.Sprintf("/%s/%s/", host, v)}
            _, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), v)
        }
        _, _ = fmt.Fprintf(w, "</pre>\n")
    } else {
        var path string
        s = s[2:]
        for _, v := range s {
            path = path + "/" + v
        }

        fmt.Println(path)
        hostname, _ := os.Hostname()
        if hostname == host {
            f, err := os.Open(path)
            if err != nil {
                p.Error(w, p.toHTTPError(err))
                return
            }
            defer f.Close()

            stat, err := f.Stat()
            if err != nil {
                fmt.Println(err)
            }

            if stat.IsDir() {
                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                _, _ = fmt.Fprintf(w, "<pre>\n")
                p.DirList(w, f, UrlPath)
                _, _ = fmt.Fprintf(w, "</pre>\n")
            } else {
                filename := filepath.Base(path)
                fmt.Println(filename)

                _, contentType := getContentType(filename)
                fmt.Println(contentType)

                w.Header().Set("Content-Disposition", "attachment; filename="+filename)

                w.Header().Set("Content-Type", contentType)

                w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

                _, _ = f.Seek(0, 0)

                _, _ = io.Copy(w, f)
            }
        } else {
            p.ClientDownload(w, host, UrlPath)
        }
    }
}

func (p *Config) ClientDownload(w http.ResponseWriter, HostName string, UrlPath string) {
    var IsDir bool
    //http://seed:8888/xt2/var/log/glusterfs/
    uu := fmt.Sprintf("http://%s:%d%s", HostName, p.Port, UrlPath)
    fmt.Println(uu)

    b := strings.HasSuffix(UrlPath, "/")
    if b == true {
        IsDir = true
    }

    resp, err := http.Get(uu)
    if err != nil {
        fmt.Println()
        os.Exit(-1)
    }

    defer resp.Body.Close()

    fmt.Println(IsDir)
    if IsDir == true {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
    } else {
        w.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))

        w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

        w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
    }
    _, _ = io.Copy(w, resp.Body)
}

func main() {

    c := Config{}

    c.GenerateTmpConfigFile()

    c.UnPraiseConfigFile()

    http.HandleFunc("/", c.ShowRootPath)

    for _, v := range c.Node {
        http.HandleFunc(fmt.Sprintf("/%s/", v), c.ShowSource)
    }

    port := fmt.Sprintf(":%d", c.Port)
    _ = http.ListenAndServe(port, nil)
}
