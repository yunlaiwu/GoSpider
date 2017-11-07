package main

import (
    "fmt"
    "strings"

    "GoSpider/libs/log"
)

const (
    LOG_CONFIG_FILE = "spider_log.conf"
    PRINT_TO_CONSOLE = true
)

var (
    gLog *logger.Logger
)

func initLogger() {
    conf := logger.NewLogConfig(LOG_CONFIG_FILE)
    err := conf.LoadConfig()
    if err != nil {
        fmt.Errorf("load log conf file failed: %v, use default config", err)
    }
    gLog = logger.NewLogger(conf)
    gLog.Infof("init logger %s done, %v", "ddd", conf)
}

func logDebug(args ...interface{}) {
    if gLog != nil {
        gLog.Debug(args...)
    }
    myPrintln(args...)
}

func logDebugf(format string, args ...interface{}) {
    if gLog != nil {
        gLog.Debugf(format, args...)
    }
    myPrintf(format, args...)
}

func logInfo(args ...interface{}) {
    if gLog != nil {
        gLog.Info(args...)
    }
    myPrintln(args...)
}

func logInfof(format string, args ...interface{}) {
    if gLog != nil {
        gLog.Infof(format, args...)
    }
    myPrintf(format, args...)
}

func logWarning(args ...interface{}) {
    gLog.Warning(args...)
    myPrintln(args...)
}

func logWarningf(format string, args ...interface{}) {
    if gLog != nil {
        gLog.Warningf(format, args...)
    }
    myPrintf(format, args...)
}

func logError(args ...interface{}) {
    if gLog != nil {
        gLog.Error(args...)
    }
    myPrintln(args...)
}

func logErrorf(format string, args ...interface{}) {
    if gLog != nil {
        gLog.Errorf(format, args...)
    }
    myPrintf(format, args...)
}

func logFatal(args ...interface{}) {
    if gLog != nil {
        gLog.Fatal(args...)
    }
    myPrintln(args...)
}

func logFatalf(format string, args ...interface{}) {
    if gLog != nil {
        gLog.Fatalf(format, args...)
    }
    myPrintf(format, args...)
}

func myPrintln(args ...interface{})  {
    if PRINT_TO_CONSOLE == false {
        return
    }
    fmt.Println(args...)
}

func myPrintf(format string, args ...interface{})  {
    if PRINT_TO_CONSOLE == false {
        return
    }
    s := fmt.Sprintf(format, args...)
    fmt.Printf("[thread %v] %v", GoID(), s)
    //哥就是为了输出看着舒服些。。。
    if strings.Contains(format, "\n") == false {
        fmt.Println("")
    }
}
