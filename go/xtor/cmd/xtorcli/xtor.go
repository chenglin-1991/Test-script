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
	"fmt"
	"github.com/crackcell/gotabulate"
	. "github.com/xtao/xtor/client"
	. "github.com/xtao/xtor/common"
)


type xtorCommand struct {
	client *XtorClient
}

func newXtorCommand(c *XtorClient) *xtorCommand {
	return &xtorCommand{
		client: c,
	}
}

func (cmd *xtorCommand) Du(vol string, rpath string, obj bool) {
	reply, err := cmd.client.Du(vol, rpath, obj)
	if err != nil {
		fmt.Printf("Error: Failed to retrieve du of %s:%s\n", vol, rpath)
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Println(reply)
	}
}

func (cmd *xtorCommand) Fsstat(vol string, mntpt string, detail bool) {
	reply, err := cmd.client.Fsstat(vol)
	if err != nil {
		fmt.Printf("Error: Failed to retrieve fsstat of %s\n", vol)
		fmt.Printf("%s\n", err.Error())
		return
	}

	tabulator := gotabulate.NewTabulator()
	tabulator.SetFirstRowHeader(true)
	tabulator.SetFormat("grid")

	var stable [][]string

	stable = append(stable,
		[]string{
			"Volume",
			"Total",
			"Used",
			"Available",
			"Use%%",
			"Mounted on"})

	free := reply.Free
	total := reply.Total
	used := total - free
	percent := fmt.Sprintf("%.2f%%",(total - free) * 100 / total)

	stable = append(stable,
		[]string{
			vol,
			ShowSize(total),
			ShowSize(used),
			ShowSize(free),
			percent,
			mntpt})


	fmt.Print(tabulator.Tabulate(stable))

	if detail != true {
		return
	}

	fmt.Println("")
	fmt.Println("Detailed Brick Statistics:")

	var table [][]string
	table = append(table,
		[]string{
			"Device",
			"Dir",
			"Total",
			"Avail"})
	for _, brick := range reply.Bricks {
		table = append(table,
			[]string{
				brick.Device,
				brick.Dir,
				brick.Total,
				brick.Free})
	}
	fmt.Print(tabulator.Tabulate(table))
}
