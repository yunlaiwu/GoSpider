package main

import (
    "os"
    "container/list"
    "sync"
    "errors"
    "time"
)

type UrlMgr struct {
    doneFile *os.File
    doneMap  map[string]int64

    taskList *list.List
    taskListLock sync.Mutex

    finishChan chan struct{}
    urlChan    chan string
}

func NewUrlMgr() *UrlMgr {
    return &UrlMgr{

    }
}

func (self *UrlMgr) Start(path string) (err error) {
    //加载已经完成的url记录文件，并初始化完成url的map
    //liujia: os.Open just open file for reading, for writing it need os.OpenFile
    self.doneFile, err = os.Open(path)
    if err != nil {
        return err
    }

    l, err := ReadFileLines(self.doneFile)
    if err != nil {
        return err
    }

    defer l.Init()

    self.doneMap = make(map[string]int64)
    for elem := l.Front(); elem != nil; elem = elem.Next() {
        self.doneMap[elem.Value.(string)] = 1
    }

    //初始化，启动记录线程
    self.taskList = list.New()
    self.urlChan = make(chan string, 1024)
    self.finishChan = make(chan struct{})
    go self.run()

    return nil
}

func (self *UrlMgr) Stop()  {
    self.finishChan <- struct{}{}

    close(self.finishChan)
    close(self.urlChan)
}

func (self *UrlMgr) Push(url string) {
    if len(url) == 0 {
        return
    }

    if _, exist := self.doneMap[url]; exist {
        return
    }

    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()
    self.taskList.PushBack(url)
}

func (self *UrlMgr) Pop() (url string, err error) {
    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()

    if self.taskList.Len() > 0 {
        first := self.taskList.Front()
        url = first.Value.(string)
        self.taskList.Remove(first)
        return url, nil
    }else {
        return "", errors.New("no url at this time")
    }
}

func (self *UrlMgr) OnFinished(handledUrl string, success bool) (newUrl string) {
    if !success {
        logErrorf("Failed to download HTML for %v", handledUrl)
        self.Push(handledUrl)
        return
    }

    self.doneMap[handledUrl] = TimeMillSecond()
    self.urlChan <- handledUrl
    return ""
}

func (self *UrlMgr) run() {
    if self.doneFile == nil {
        panic("done file is not open yet")
    }

    defer func() {
        //do something to clean up
        self.doneFile.Sync()
        self.doneFile.Close()
    }()

    for{
        select {
        case <-self.finishChan:
            return
        case url := <-self.urlChan:
            self.doneFile.WriteString(url+"\n")
        case <-time.After(10 * time.Second):
            self.doneFile.Sync()
        }
    }
}
