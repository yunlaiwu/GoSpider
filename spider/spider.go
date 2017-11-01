package main

import (
    "fmt"
    "os"
    "os/signal"
)

func main() {
    //init logger
    fmt.Println("init spider logger")
    initLogger()

    //进程收到的退出信号
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    select {
    case <-signals:
        logInfo("Received OS signal, stop spider")
    }

    logInfo("exit spider process")
}

