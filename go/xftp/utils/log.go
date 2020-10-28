package utils

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
)

type Log struct {
    Path    string
    Rlogger *log.Logger
}

var Rlog *Log

func LogInit(path string) *Log {
    err := os.Mkdir(filepath.Dir(XftpLogPath), 0666)
    if err != nil && !os.IsExist(err) {
        os.Exit(-1)
    }

    rlogFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
        log.Panic(fmt.Sprintf("failed to open log file %s", path))
    }
    rlogger := log.New(rlogFile, "", log.Ldate|log.Ltime|log.Lshortfile)
    l := &Log{Rlogger: rlogger}

    return l
}
