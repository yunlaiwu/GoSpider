package main

import (
    "fmt"
    "flag"
    "os"
    "os/signal"
)

var (
    spe      *SpiderEngine
    storeMgr IResStorer
    doneChan chan error
)

func main() {
    //解析命令行
    restype := flag.String("restype", "", "resource to download")
    resfile := flag.String("resfile", "", "resource file which has resource info")
    savedir := flag.String("savedir", "", "directory to save")

    flag.Parse()

    //检查命令行
    const example = "usage: spider -restype=book-comment or book-review -resfile=./book-comment.txt -savedir=./book"

    switch *restype {
    case "book-comment", "book-review":
        fmt.Println("restype:", *restype)
    default:
        fmt.Println("invalid restype")
        fmt.Println(example)
        os.Exit(-1)
    }

    if FileExist(*resfile) == false {
        fmt.Println("invalid resfile or file not exist")
        fmt.Println(example)
        os.Exit(-1)
    }else {
        fmt.Println("resfile:", *resfile)
    }

    if len(*savedir) == 0 || CreateDirIfNotExist(*savedir) != nil {
        fmt.Println("invalid savedir")
        fmt.Println(example)
        os.Exit(-1)
    }else {
        fmt.Println("savedir:", *savedir)
    }

    //init logger
    fmt.Println("init spider logger")
    initLogger()

    //初始化engine
    spe = NewSpiderEngine()
    spe.Start()

    go func() {
        if *restype == "book-comment" {
            storeMgr = NewBookCommentStore()
            storeMgr.Start(*resfile, *savedir)
        }
    }()

    //进程收到的退出信号
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    doneChan = make(chan error, 1)

    select {
    case <-signals:
        logInfo("Received OS signal, stop spider")
        spe.Stop()
        break
    case <- doneChan:
        logInfo("Received done signal, stop spider")
        spe.Stop()
        break
    }

    logInfo("exit spider process")
}

