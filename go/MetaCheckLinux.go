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
    "syscall"
    "time"
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

    ATimeFlag      bool
    CTimeFlag      int
    MTimeFlag      int
    CheckMd5Sum    int
    CheckDentryCnt int

    ModeFlag    int
    SizeFlag    int
    InoFlag     int
    BlksizeFlag int
    BlocksFlag  int
    UidFlag     int
    GidFlag     int
    NlinkFlag   int

    ThreadCnt int
)

type MetaInfo struct {
    Path string

    Atime time.Time
    Ctime time.Time
    Mtime time.Time

    Mode    os.FileMode
    Size    int64
    Ino     uint64
    Blksize int64
    Blocks  int64
    Uid     uint32
    Gid     uint32
    Nlink   uint64

    Check  bool
    IsFile bool
    Md5Sum string
}

var exit = make(chan int)
var Metabuf chan string

func init() {
    kingpin.Flag("thread", "How many threads verify a directory at the same time default (8)").
        Short('t').Default("16").IntVar(&ThreadCnt)
    kingpin.Flag("atime", "compare atime default false").Short('a').Default("false").BoolVar(&ATimeFlag)
    kingpin.Flag("ctime", "compare ctime default true").Short('c').Default("1").IntVar(&CTimeFlag)
    kingpin.Flag("mtime", "compare mtime default true").Short('m').Default("1").IntVar(&MTimeFlag)
    kingpin.Flag("mode", "compare mode default true").Short('M').Default("1").IntVar(&ModeFlag)
    kingpin.Flag("size", "compare size default true").Short('S').Default("1").IntVar(&SizeFlag)
    kingpin.Flag("inode", "compare inode number default true").
        Short('I').Default("1").IntVar(&InoFlag)
    kingpin.Flag("blksize", "compare block size default true").
        Short('k').Default("1").IntVar(&BlksizeFlag)
    kingpin.Flag("blocks", "compare blocks default true").Short('K').Default("1").IntVar(&BlocksFlag)
    kingpin.Flag("uid", "compare uid default true").Short('U').Default("1").IntVar(&UidFlag)
    kingpin.Flag("gid", "compare gid default true").Short('G').Default("1").IntVar(&GidFlag)
    kingpin.Flag("nlinks", "compare nlinks default true").Short('N').Default("1").IntVar(&NlinkFlag)
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
    var nameCnt = 0

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
        if err == io.EOF {
            break
        } else if err != nil {
            rlogger.Println(fmt.Sprintf("failed to readdir path %s "+
                "err is %s", m.Path, err))
            goto out
        }

        for _, val := range list {
            if val.Name() == BaseName {
                nameCnt++
            }
        }
    }

    if nameCnt != 1 {
        rlogger.Printf("name count = %d, name = %s", nameCnt, BaseName)
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

    if InoFlag != 0 && m.Ino != m2.Ino {
        clogger.Println(fmt.Sprintf("Ino compare failed src_path"+
            " %s Ino %v dsr_path %s Ino %v", m.Path, m.Ino, m2.Path, m2.Ino))
    }

    if BlksizeFlag != 0 && m.Blksize != m2.Blksize {
        clogger.Println(fmt.Sprintf("Blksize compare failed src_path"+
            " %s Blksize %v dsr_path %s Blksize %v", m.Path, m.Blksize, m2.Path, m2.Blksize))
    }

    if BlocksFlag != 0 && m.Blocks != m2.Blocks {
        clogger.Println(fmt.Sprintf("Blocks compare failed src_path"+
            " %s Blocks %v dsr_path %s Blocks %v", m.Path, m.Blocks, m2.Path, m2.Blocks))
    }

    if UidFlag != 0 && m.Uid != m2.Uid {
        clogger.Println(fmt.Sprintf("Uid compare failed src_path"+
            " %s Uid %v dsr_path %s Uid %v", m.Path, m.Uid, m2.Path, m2.Uid))
    }

    if GidFlag != 0 && m.Gid != m2.Gid {
        clogger.Println(fmt.Sprintf("Gid compare failed src_path"+
            " %s Gid %v dsr_path %s Gid %v", m.Path, m.Gid, m2.Path, m2.Gid))
    }

    if NlinkFlag != 0 && m.Nlink != m2.Nlink {
        clogger.Println(fmt.Sprintf("Nlink compare failed src_path"+
            " %s Nlink %v dsr_path %s Nlink %v", m.Path, m.Nlink, m2.Path, m2.Nlink))
    }

    if ATimeFlag && m.Atime != m2.Atime {
        clogger.Println(fmt.Sprintf("ATimeSince compare failed src_path"+
            " %s ATimeSince %v dsr_path %s ATimeSince %v", m.Path, m.Atime, m2.Path, m2.Atime))
    }

    if CTimeFlag != 0 && m.Ctime != m2.Ctime {
        clogger.Println(fmt.Sprintf("CTimeSince compare failed src_path"+
            " %s CTimeSince %v dsr_path %s CTimeSince %v", m.Path, m.Ctime, m2.Path, m2.Ctime))
    }

    if MTimeFlag != 0 && m.Mtime != m2.Mtime {
        clogger.Println(fmt.Sprintf("MTimeSince compare failed src_path"+
            " %s MTimeSince %v dsr_path %s MTimeSince %v", m.Path, m.Mtime, m2.Path, m2.Mtime))
    }

    if CheckMd5Sum != 0 && !m.Check && m.Md5Sum != m2.Md5Sum {
        clogger.Println(fmt.Sprintf("Md5Sum compare failed src_path"+
            " %s Md5Sum %v dsr_path %s Md5Sum %v", m.Path, m.Md5Sum, m2.Path, m2.Md5Sum))
    }

    if CheckDentryCnt != 0 && m.IsFile {
        ret := m2.CheckRepeatDentry()
        if ret == DentryCntFailed {
            clogger.Println(fmt.Sprintf("In path %s maybe have more dentry %s",
                filepath.Dir(m2.Path), filepath.Base(m2.Path)))
        }
    }
}

func (m *MetaInfo) GetMetaInfo(f os.FileInfo, path string) {
    m.Path = path

    m.Mode = f.Mode()
    m.Size = f.Size()
    m.Ino = f.Sys().(*syscall.Stat_t).Ino
    m.Blksize = f.Sys().(*syscall.Stat_t).Blksize
    m.Blocks = f.Sys().(*syscall.Stat_t).Blocks
    m.Uid = f.Sys().(*syscall.Stat_t).Uid
    m.Gid = f.Sys().(*syscall.Stat_t).Gid
    m.Nlink = f.Sys().(*syscall.Stat_t).Nlink

    AAtime := f.Sys().(*syscall.Stat_t).Atim
    ACtime := f.Sys().(*syscall.Stat_t).Ctim
    AMtime := f.Sys().(*syscall.Stat_t).Mtim
    m.Atime = time.Unix(AAtime.Sec, AAtime.Nsec)
    m.Ctime = time.Unix(ACtime.Sec, ACtime.Nsec)
    m.Mtime = time.Unix(AMtime.Sec, AMtime.Nsec)

    m.Check = !f.IsDir()
    m.IsFile = !f.IsDir()

    if m.Check {
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

    /* Why don't use Directory_walk FileInfo ?
     * Directory_walk get FileInfo by readdir,readdir return the stbuf
     * by hash_brick,not merged,so not explicit
     *
     * the dstpath get FileInfo bt stat,glusterfs stat will merge stbuf for
     * all brick,so if we use Directory_walk FileInfo,it will check err for
     * meta info
     */
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
