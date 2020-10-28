package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"regexp"
	"strconv"
	"strings"
	"path"
	"path/filepath"
)

type Entry struct {
	Name          string
	Directory     string
	Types         []string
	Options       []string
	DumpFrequency int
	PassNumber    int
}

func parse(filename string) ([]*Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	entries := make([]*Entry, 0, 4)

	for {
		line, rserr := buf.ReadString('\n')
		if rserr != nil && rserr != io.EOF {
			return nil, rserr
		}
		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			entries = append(entries, entry)
		}

		if rserr == io.EOF {
			break
		}
	}

	return entries, nil
}

func unescape(s string) string {
	s = strings.Replace(s, "\\011", "\t", -1)
	s = strings.Replace(s, "\\012", "\n", -1)
	s = strings.Replace(s, "\\040", " ", -1)
	s = strings.Replace(s, "\\134", "\\", -1)
	return s
}

var splitRegExp = regexp.MustCompile("\\s+")

func parseLine(untrimmedLine string) (*Entry, error) {
	line := strings.TrimSpace(untrimmedLine)
	if len(line) == 0 {
		return nil, nil
	}
	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	fields := splitRegExp.Split(line, -1)
	if len(fields) != 6 {
		return nil, errors.New(fmt.Sprintf("Each line must consist 6 fields but got %d", len(fields)))
	}

	entry := &Entry{}
	entry.Name = unescape(fields[0])
	entry.Directory = unescape(fields[1])
	entry.Types = strings.Split(unescape(fields[2]), ",")
	entry.Options = strings.Split(unescape(fields[3]), ",")

	num, err := strconv.ParseUint(fields[4], 10, 31)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't parse dump frequency field: %s", err))
	}
	entry.DumpFrequency = int(num)

	num, err = strconv.ParseUint(fields[5], 10, 31)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't parse pass number field: %s", err))
	}
	entry.PassNumber = int(num)

	return entry, nil
}

func mountPoint(p string) string {
    pi, err := os.Stat(p)
    if err != nil {
        return ""
    }

    odev := pi.Sys().(*syscall.Stat_t).Dev

    for p != "/" {
        _path := filepath.Dir(p)

        in, err := os.Stat(_path)
        if err != nil {
            return ""
        }

        if odev != in.Sys().(*syscall.Stat_t).Dev {
            break
        }

        p = _path
    }

    return p
}

func ExtractPath(p string) (error, string, string, string, string) {
	var server, volstr, volume, subdir string

	mntpt := mountPoint(p)
	entries, err := parse("/etc/mtab")
	if err != nil {
		errMsg := "Failed to parse /etc/mtab"
		fmt.Println(errMsg)
		return errors.New(errMsg), "", "", "", ""
	}

	volstr = ""
	server = ""
	subdir = ""
	for _, ent := range entries {
		if ent.Directory == mntpt {
			findNFS := false
			for _, a := range ent.Types {
				if a == "nfs" {
					findNFS = true
					break
				}
			}
			if findNFS == false {
				errmsg := fmt.Sprintf("Invalid: %s is not a NFS Filesystem", p)
				return errors.New(errmsg), "", "", "", ""
			}
			segs := strings.Split(ent.Name, ":")
			server = segs[0]
			volstr = strings.TrimPrefix(segs[1], "/")
			volume = strings.Split(volstr, "/")[0]
			subdir = strings.TrimPrefix(volstr, volume)
			break
		}
	}

	if server == "" {
		errMsg := "Failed to match the path entry in mntab"
		return errors.New(errMsg), "", "", "", ""
	}

	relativePath := "/"
	if p != mntpt {
		relativePath = strings.TrimPrefix(p, mntpt)
	}

	if subdir != "" {
		relativePath = path.Join(subdir, relativePath)
	}

	return nil, server, volume, mntpt, relativePath
}

/*
func main() {

	path, _ := filepath.Abs("./")
	
	err, svr, vol, rpath := ExtractPath(path)

	fmt.Printf("Vol Name:%s\n", vol)
	fmt.Printf("\tSvr:%s\n", svr)
	fmt.Printf("\tpath:%s\n",rpath)
}
*/
