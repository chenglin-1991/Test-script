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

package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
	. "github.com/xtao/xtor/client"
	. "github.com/xtao/xtor/common"
)

var (

    /*
     * du command
     */
    du = kingpin.Command("du", "Environment configuration")
    du_path = du.Arg("path", "Specified path to show statistics").
	    Default(".").
	    String()

    du_obj = du.Flag("inodes", "List inode usage information instead of block usage").
        Short('i').
        Bool()

    /*
     * fsstat command
     */
    df = kingpin.Command("fsstat", "Statistics for file system")
    df_path = df.Arg("path", "Specified path to show statistics").
	    Default(".").
	    String()
    df_detail = df.Flag("detail", "Show detailed statistics for all bricks").
        Short('d').
        Bool()
)

func main() {
	cred := GetClientUserInfo()
	kingpin.CommandLine.HelpFlag.Short('h')
	args := kingpin.Parse()
	cmds := strings.Split(args, " ")

	switch cmds[0] {
	/*
         * du command
         */
	case "du":
		path, _ := filepath.Abs(*du_path)
		err:= os.Chdir(*du_path)
		if err != nil {
			/*
                         * Invalid path
                         */
			fmt.Printf("Path %s is invalid. %s\n", *du_path,
			    err.Error())
			return
		}
		err, svr, vol, _, rpath := ExtractPath(path)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		svr = fmt.Sprintf("%s:8765", svr)
		client := NewXtorClient(svr, cred)
		xtorCmd := newXtorCommand(client)
		xtorCmd.Du(vol, rpath, *du_obj)

	case "fsstat":
		path, _ := filepath.Abs(*df_path)
		err, svr, vol, mntpt, _ := ExtractPath(path)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		svr = fmt.Sprintf("%s:8765", svr)
		client := NewXtorClient(svr, cred)
		xtorCmd := newXtorCommand(client)

		xtorCmd.Fsstat(vol, mntpt, *df_detail)
	}
}
