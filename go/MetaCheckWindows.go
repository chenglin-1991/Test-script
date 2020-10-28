package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	DentryCntErr = iota
	DentryCntSuccess
	DentryCntFailed
	ReaddirSize = 4096
)

var (
	srcDirectoryPath string
	dstDirectoryPath string
	CheckLogPath     string
	RunLogPath       string
	rlogger          *log.Logger
	clogger          *log.Logger

	ATimeFlag   bool
	CTimeFlag   int
	MTimeFlag   int
	CheckMd5Sum int
	CheckDentryCnt int

	ModeFlag int
	SizeFlag int

	ThreadCnt int
)

type MetaInfo struct {
	Path string

	FileAttributes uint32
	CreationTime   int64
	LastAccessTime int64
	LastWriteTime  int64

	Mode   os.FileMode
	Size   int64
	Check  bool
	IsFile bool
	Md5Sum string
}

var exit = make(chan int)
var Metabuf chan string

func init() {
	kingpin.Flag("thread", "How many threads verify a directory at the same default (8)").
		Short('t').Default("8").IntVar(&ThreadCnt)
	kingpin.Flag("atime", "compare atime default false").Short('a').Default("false").BoolVar(&ATimeFlag)
	kingpin.Flag("ctime", "compare ctime default true").Short('c').Default("1").IntVar(&CTimeFlag)
	kingpin.Flag("mtime", "compare mtime default true").Short('m').Default("1").IntVar(&MTimeFlag)
	kingpin.Flag("mode", "compare mode default true").Short('M').Default("1").IntVar(&ModeFlag)
	kingpin.Flag("size", "compare size default true").Short('S').Default("1").IntVar(&SizeFlag)
	kingpin.Flag("md5", "compare md5 default true").Short('n').Default("1").IntVar(&CheckMd5Sum)
	kingpin.Flag("dentry_cnt", "compare the dentry cnt").Short('R').Default("0").IntVar(&CheckDentryCnt)

	kingpin.Flag("SrcPath", "the source path must give").Short('s').Required().StringVar(&srcDirectoryPath)
	kingpin.Flag("DstPath", "the dest path must give").Short('d').Required().StringVar(&dstDirectoryPath)

	kingpin.Flag("RunLogPath", "the runing log path").Short('r').Default("./run.log").StringVar(&RunLogPath)
	kingpin.Flag("CheckPath", "the check log path").Short('C').Default("./check.log").StringVar(&CheckLogPath)

	kingpin.Version("0.0.1")
	kingpin.Parse()

	rlogFile, err := os.OpenFile(RunLogPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to open log file %s", RunLogPath))
	}
	rlogger = log.New(rlogFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	clogFile, err := os.OpenFile(CheckLogPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to open log file %s", CheckLogPath))
	}
	clogger = log.New(clogFile, "", log.Ldate|log.Ltime)

	Metabuf = make(chan string, ThreadCnt)
}


func (m *MetaInfo) CheckRepeatDentry() int {
	var ret = DentryCntErr
	var NameCnt = 0

	ParentPath := filepath.Dir(m.Path)
	BaseName := filepath.Base(m.Path)

	f, err := os.Open(ParentPath)
	if err != nil {
		rlogger.Println(fmt.Sprintf("failed to open path %s "+
			"err is %s", m.Path, err))
		goto out
	}
	defer f.Close()

	for {
		list, err := f.Readdir(ReaddirSize)
		if err == io.EOF || len(list) == 0 {
			break
		} else if err != nil {
			rlogger.Println(fmt.Sprintf("failed to readdir path %s "+
				"err is %s", m.Path, err))
			goto out
		}

		for _, val := range list {
			cret := strings.Compare(val.Name(), BaseName)
			if cret == 0 {
				NameCnt++
			}
		}
	}

	if NameCnt != 1 {
		ret = DentryCntFailed
	} else {
		ret = DentryCntSuccess
	}

out:
	return ret
}

func (m *MetaInfo) CompareMetaInfo(m2 *MetaInfo) {
	if m2 == nil {
		return
	}

	if ModeFlag != 0 && m.Mode != m2.Mode {
		clogger.Println(fmt.Sprintf("mode compare failed src_path"+
			" %s Mode %v dsr_path %s Mode %v", m.Path, m.Mode, m2.Path, m2.Mode))
	}

	if SizeFlag != 0 && m.Size != m2.Size {
		clogger.Println(fmt.Sprintf("Size compare failed src_path"+
			" %s Size %v dsr_path %s Size %v", m.Path, m.Size, m2.Path, m2.Size))
	}

	if ATimeFlag && m.LastAccessTime != m2.LastAccessTime {
		clogger.Println(fmt.Sprintf("LastAccessTime compare failed src_path"+
			" %s LastAccessTime %v dsr_path %s LastAccessTime %v", m.Path, m.LastAccessTime, m2.Path, m2.LastAccessTime))
	}

	if MTimeFlag != 0 && m.LastWriteTime != m2.LastWriteTime {
		clogger.Println(fmt.Sprintf("LastWriteTime compare failed src_path"+
			" %s LastWriteTime %v dsr_path %s LastWriteTime %v", m.Path, m.LastWriteTime, m2.Path, m2.LastWriteTime))
	}

	if CTimeFlag != 0 && m.CreationTime != m2.CreationTime {
		clogger.Println(fmt.Sprintf("CreationTime compare failed src_path"+
			" %s CreationTime %v dsr_path %s CreationTime %v", m.Path, m.CreationTime, m2.Path, m2.CreationTime))
	}

	if CheckMd5Sum != 0 && !m.Check && m.Md5Sum != m2.Md5Sum {
		clogger.Println(fmt.Sprintf("Md5Sum compare failed src_path"+
			" %s Md5Sum %v dsr_path %s Md5Sum %v", m.Path, m.Md5Sum, m2.Path, m2.Md5Sum))
	}

	if CheckDentryCnt != 0 && m.IsFile {
		ret := m2.CheckRepeatDentry()
		if ret == DentryCntFailed {
			clogger.Println(fmt.Sprintf("In path %s maybe have more dentry %s"+
				filepath.Dir(m2.Path), filepath.Base(m2.Path)))
		}
	}
}

func (m *MetaInfo) GetMetaInfo(f os.FileInfo, path string) {
	m.Path = path

	m.Mode = f.Mode()
	m.Size = f.Size()

	fileSys := f.Sys().(*syscall.Win32FileAttributeData)
	m.FileAttributes = fileSys.FileAttributes

	nanoseconds := fileSys.LastAccessTime.Nanoseconds()
	m.LastAccessTime = nanoseconds / 1e6

	nanoseconds = fileSys.LastWriteTime.Nanoseconds()
	m.LastWriteTime = nanoseconds / 1e6

	nanoseconds = fileSys.CreationTime.Nanoseconds()
	m.CreationTime = nanoseconds / 1e6

	m.Check = !f.IsDir()

	if !m.Check {
		pf, err := os.Open(path)
		if err != nil {
			m.Check = false
			goto out
		}

		defer pf.Close()

		md5h := md5.New()
		_, err = io.Copy(md5h, pf)
		if err != nil {
			m.Check = false
			goto out
		}

		m.Md5Sum = hex.EncodeToString(md5h.Sum(nil))
	}

out:
	return
}

func CheckMetaInfo(path string) {
	srcPath := path
	dstPath1 := path[len(srcDirectoryPath):]
	if len(dstPath1) == 0 {
		return
	}
	dstPath := fmt.Sprintf("%s%s", dstDirectoryPath, dstPath1)

	AFileinfo := new(MetaInfo)
	BFileinfo := new(MetaInfo)

	sfileinfo, err := os.Lstat(srcPath)
	if err != nil {
		rlogger.Println(err.Error())
		return
	}

	dfileinfo, err := os.Lstat(dstPath)
	if err != nil {
		rlogger.Println(err.Error())
		return
	}

	AFileinfo.GetMetaInfo(sfileinfo, srcPath)
	BFileinfo.GetMetaInfo(dfileinfo, dstPath)
	AFileinfo.CompareMetaInfo(BFileinfo)

	return
}

func HandleDentry() {
	for {
		path, ok := <-Metabuf
		if ok {
			CheckMetaInfo(path)
		} else {
			exit <- 1
			break
		}
	}
}

func Directory_walk(path string, sfileinfo os.FileInfo, err error) error {
	if err != nil {
		rlogger.Println(fmt.Sprintf(" readdir failed %s", err.Error()))
		return nil
	}

	Metabuf <- path

	return nil
}

func main() {
	srcDirectoryPath, _ = filepath.Abs(srcDirectoryPath)
	stbuf, err := os.Lstat(srcDirectoryPath)
	if err != nil || !stbuf.IsDir() {
		if err != nil {
			fmt.Printf("failed to opendir %s", err.Error())
			rlogger.Println(fmt.Sprintf("failed to opendir %s", err.Error()))
		} else {
			fmt.Printf("Is it a directory %s", srcDirectoryPath)
			rlogger.Println(fmt.Sprintf("Is it a directory %s", srcDirectoryPath))
		}
		os.Exit(1)
	}

	dstDirectoryPath, _ = filepath.Abs(dstDirectoryPath)
	stbuf, err = os.Lstat(dstDirectoryPath)
	if err != nil || !stbuf.IsDir() {
		if err != nil {
			fmt.Printf("failed to opendir %s", err.Error())
			rlogger.Println(fmt.Sprintf("failed to opendir %s", err.Error()))
		} else {
			fmt.Printf("Is it a directory %s", dstDirectoryPath)
			rlogger.Println(fmt.Sprintf("Is it a directory %s", dstDirectoryPath))
		}
		os.Exit(1)
	}

	for i := 0; i < ThreadCnt; i++ {
		go HandleDentry()
	}

	err = filepath.Walk(srcDirectoryPath, Directory_walk)
	if err != nil {
		fmt.Printf("failed to walk %s", err.Error())
		os.Exit(1)
	}

	close(Metabuf)

	for i := 0; i < ThreadCnt; i++ {
		<-exit
	}

	fmt.Printf("Check %s with %s complete! to see %s and %s for result!\n",
		srcDirectoryPath, dstDirectoryPath, RunLogPath, CheckLogPath)
}
