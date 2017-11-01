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