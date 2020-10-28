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

package common

import (
	"io"
	"os"
	"log"
	"fmt"
	"os/user"
	"errors"
	"reflect"
	"net/http"
	"encoding/json"
)

const (
	XTOR_DEFAULT_SERVER = "127.0.0.1:8765"
)

type LoggerConfig struct {
    Logfile string
}

var Logger *log.Logger = nil

func LoggerInit(config *LoggerConfig) error {
    logFile, err := os.OpenFile(config.Logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
    if err != nil {
        fmt.Printf("open file error=%s\r\n", err.Error())
        return err
    }

    Logger = log.New(logFile, "\n", log.Ldate | log.Ltime | log.Lshortfile)
    
    return nil
}


func Contain(obj interface{}, target interface{}) (bool, error) {
    targetValue := reflect.ValueOf(target)
    switch reflect.TypeOf(target).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if targetValue.Index(i).Interface() == obj {
                return true, nil
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            return true, nil
        }
    }

    return false, errors.New("not in array")
}

func IndexOf(array interface{}, obj interface{}) (int, error) {
	value := reflect.ValueOf(array)

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0 ; i < value.Len(); i++ {
			if value.Index(i).Interface() == obj {
				return i, nil
			}
		}
		return -1, errors.New("not in array")
	default:
		return 0, errors.New("not supported type")
	}
}

func WriteJSON(status int, data interface{}, w http.ResponseWriter) error {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

func CopyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}

type UserAccountInfo struct {
    Username string       `json:"user"`
    Groupname string      `json:"group"`
    Uid string            `json:"uid"`
    Gid string            `json:"gid"`
}

func GetClientUserInfo() *UserAccountInfo {
    curUser, err := user.Current()
    if err != nil {
        fmt.Printf("Fail to get current user: %s\n",
            err.Error())
        return &UserAccountInfo{}
    }

    info := &UserAccountInfo {
        Username: curUser.Username,
        Uid: curUser.Uid,
        Gid: curUser.Gid,
    }

    return info
}

func ShowSize(input float64) string {
	unit := ""
	size := input

	for {
		if size / 1024 < 1 {
			break
		}

		size = size / 1024
		switch unit {
		case "":
			unit = "K"
		case "K":
			unit = "M"
		case "M":
			unit = "G"
		case "G":
			unit = "T"
		case "T":
			unit = "P"
		}
	}

	return fmt.Sprintf("%.2f%s", size, unit)
}
