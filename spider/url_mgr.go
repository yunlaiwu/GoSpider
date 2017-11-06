package main

import (
    "os"
    "container/list"
    "sync"
    "errors"
    "time"
)

type DownTask struct {
    res string
    url string
    cb  ICallbacker
    params map[string]string
}

func NewDownTask(res string, url string, cb  ICallbacker, params map[string]string) *DownTask {
    return &DownTask{
        res: res,
        url: url,
        cb: cb,
        params:params,
    }
}

func NewDownTaskWihtouParam(res string, url string, cb  ICallbacker) *DownTask {
    return &DownTask{
        res: res,
        url: url,
        cb: cb,
        params:nil,
    }
}

func (self DownTask) Valid() bool {
    if len(self.res) == 0 || len(self.url) == 0 || self.cb == nil {
        return false
    }
    return true
}

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

func (self *UrlMgr) Push(task *DownTask) {
    if task == nil || !task.Valid() {
        return
    }

    if _, exist := self.doneMap[task.url]; exist {
        return
    }

    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()
    self.taskList.PushBack(task)
}

func (self *UrlMgr) Pop() (task *DownTask, err error) {
    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()

    if self.taskList.Len() > 0 {
        first := self.taskList.Front()
        task = first.Value.(*DownTask)
        self.taskList.Remove(first)
        return task, nil
    }else {
        return nil, errors.New("no url at this time")
    }
}

func (self *UrlMgr) UrlDone(handledUrl string) {
    self.doneMap[handledUrl] = TimeMillSecond()
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
