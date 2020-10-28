package utils

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

const (
    TmpXftpConfigPath = "/etc/xftp/xftp.back"
    XftpConfigPath    = "/etc/xftp/xftp.conf"
    XftpLogPath       = "/var/log/xftp/xftp.log"
)

type Config struct {
    Node    []string `ini:"node" comment:"which inode you want to monitor"`
    LogPath []string `ini:"logpath" comment:"which log path you want to download"`
    Port    int      `ini:"port" comment:"the http server listen port"`
}

func (p *Config) GenerateTmpConfigFile(log *Log) {
    err := os.Mkdir(filepath.Dir(TmpXftpConfigPath), 0666)
    if err != nil && !os.IsExist(err) {
        log.Rlogger.Println(fmt.Sprintf("failed to mkdir path %s "+
            "err is %s", TmpXftpConfigPath, err))
        fmt.Println(err)
        os.Exit(-1)
    }

    cfg := ini.Empty()

    c := Config{
        Node:    []string{"seed", "AlamoD1N8serverc802", "xt3"},
        LogPath: []string{"/var/log/glusterfs", "/var/log/samba"},
        Port:    8888,
    }

    err = ini.ReflectFrom(cfg, &c)
    if err != nil {
        log.Rlogger.Println(fmt.Sprintf("failed to ReflectFrom %s", err))
        os.Exit(-1)
    }

    err = cfg.SaveTo(TmpXftpConfigPath)
    if err != nil {
        log.Rlogger.Println(fmt.Sprintf("failed to SaveTo %s", err))
        os.Exit(-1)
        return
    }
}

func (p *Config) UnPraiseConfigFile(log *Log) {
    err := os.Mkdir(filepath.Dir(XftpConfigPath), 0666)
    if err != nil && !os.IsExist(err) {
        log.Rlogger.Println(fmt.Sprintf("failed to mkdir path %s "+
            "err is %s", XftpConfigPath, err))
        os.Exit(-1)
    }

    cfg, err := ini.Load(XftpConfigPath)
    if err != nil {
        log.Rlogger.Println(fmt.Sprintf("failed to Load %s", err))
        os.Exit(-1)
    }

    err = cfg.MapTo(p)
    if err != nil {
        fmt.Println(err)
        log.Rlogger.Println(fmt.Sprintf("failed to Map %s", err))
        os.Exit(-1)
    }
}

func (p *Config) ClientDownload(w http.ResponseWriter, HostName string,
    UrlPath string) {
    var IsDir bool

    uu := fmt.Sprintf("http://%s:%d%s", HostName, p.Port, UrlPath)

    b := strings.HasSuffix(UrlPath, "/")
    if b == true {
        IsDir = true
    }

    resp, err := http.Get(uu)
    if err != nil {
        Rlog.Rlogger.Println(fmt.Sprintf("failed to get %s", uu))
        return
    }

    defer resp.Body.Close()

    if IsDir == true {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
    } else {
        w.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))

        w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

        w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
    }
    n, err := io.Copy(w, resp.Body)
    if err != nil || n == 0 {
        Rlog.Rlogger.Println(fmt.Sprintf("failed to write "+
            "%s respond %s", uu, err))
    }
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

func (p *Config) toHTTPError(err error) int {
    if os.IsNotExist(err) {
        return http.StatusNotFound
    }
    if os.IsPermission(err) {
        return http.StatusForbidden
    }
    return http.StatusInternalServerError
}

func (p *Config) ShowSource(w http.ResponseWriter, r *http.Request) {
    var UrlPath = ""
    UrlPath = r.URL.Path

    fmt.Println("showSource", UrlPath)

    s := strings.Split(r.URL.Path, "/")
    host := s[1]

    Rlog.Rlogger.Println(fmt.Sprintf("get req for %s", UrlPath))

    //show host
    if len(s) == 3 {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        _, err := fmt.Fprintf(w, "<pre>\n")
        if err != nil {
            Rlog.Rlogger.Println(fmt.Sprintf("failed to write respond for %s", UrlPath))
            goto out
        }

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
                _, contentType := getContentType(filename)

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

out:
    return
}

func getContentType(fileName string) (extension, contentType string) {
    arr := strings.Split(fileName, ".")

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

func (p *Config) ShowRootPath(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    fmt.Println("==ShowRootPath===2", path)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    _, _ = fmt.Fprintf(w, "<pre>\n")
    for _, n := range p.Node {
        n += "/"
        vurl := url.URL{Path: n}
        _, _ = fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", vurl.String(), n)
    }
    _, _ = fmt.Fprintf(w, "</pre>\n")
}
