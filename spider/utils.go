package main

import (
    "container/list"
    "bufio"
    "io"
    "strings"
    "os"
    "errors"
    "time"
    "fmt"
    "path/filepath"
    "runtime"
    "strconv"
)

func ReadFileLines(fi *os.File) (ret *list.List, err error) {
    if fi == nil {
        return ret, errors.New("invalid file handle")
    }

    ret = list.New()
    buf := bufio.NewReader(fi)
    for {
        line, err := buf.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            //report error
            return ret, err
        }

        ret.PushBack(strings.TrimSpace(line))
    }

    return ret, nil
}

func TimeMillSecond() int64 {
    return time.Now().UnixNano() / 1000000
}

func TimeSecond() int64 {
    return time.Now().Unix()
}

func Time2String(ts int64) string  {
    var secs int64 = ts/1000
    var left = ts%1000
    var s string = time.Unix(secs, 0).Format("2006-01-02 15:04:05")
    return s+fmt.Sprintf(":%v", left)
}

//得到当前工作目录，其实就是当前进程的work directory
func GetCurrentDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        fmt.Println(err)
    }
    return strings.Replace(dir, "\\", "/", -1)
}

//得到全路径，如果relPath已经是Absolute的就直接用，否则跟当前工作目录拼成Abs的
func GetFullPath(relPath string) (fullPath string) {
    if filepath.IsAbs(relPath) == false {
        return filepath.Join(GetCurrentDirectory(), relPath)
    } else {
        return relPath
    }
}

//文件是否存在
func FileExist(path string) bool {
    if fi, err := os.Stat(path); err == nil {
        if fi.IsDir() == false {
            return true
        }
    }
    return false
}

//文件夹是否存在
func DirExist(path string) bool {
    if fi, err := os.Stat(path); err == nil {
        if fi.IsDir() == true {
            return true
        }
    }
    return false
}

//不存在则创建
func CreateDirIfNotExist(path string) (err error) {
    if DirExist(path) == false {
        err = os.MkdirAll(path, 0777)
        if err != nil {
            return err
        }
    }
    return nil
}

/*
 * 获取当前的goroutine的id，这个俺也是从网上抄的
 */
func GoID() int {
    var buf [64]byte
    n := runtime.Stack(buf[:], false)
    idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
    id, err := strconv.Atoi(idField)
    if err != nil {
        panic(fmt.Sprintf("cannot get goroutine id: %v", err))
    }
    return id
}

//去掉\t \n \r
func SanityString(s string) string {
    s = strings.Replace(s, "\t", "", -1)
    s = strings.Replace(s, "\n", "", -1)
    s = strings.Replace(s, "\r", "", -1)
    return s
}

func Int2String(i int) string {
    return fmt.Sprintf("%v", i)
}

func String2Int(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

func IsMAC() bool {
    return runtime.GOOS == "darwin"
}