package main

import (
    "os"
    "container/list"
    "sync"
    "errors"
    "time"
    "fmt"
)

type DownTask struct {
    res    string
    url    string
    cb     IResHunter
    params map[string]string
    resp   []byte
}

func NewDownTask(res string, url string, cb IResHunter, params map[string]string) *DownTask {
    return &DownTask{
        res: res,
        url: url,
        cb: cb,
        params:params,
    }
}

func NewDownTaskWihtouParam(res string, url string, cb IResHunter) *DownTask {
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

func (self DownTask) String() string {
    return fmt.Sprintf("%v, %v, %v", self.res, self.url, self.params)
}

func (self DownTask) Process() {
    if self.cb != nil {
        self.cb.OnResponse(self.url, self.resp, self.params)
    }
}

type UrlMgr struct {
    taskList *list.List
    taskListLock sync.Mutex

    doneRecord bool
    doneFile *os.File
    doneMap  map[string]int64
    finishChan chan struct{}
    urlChan    chan string
}

func NewUrlMgr() *UrlMgr {
    return &UrlMgr{
        doneRecord:false,
    }
}

func (self *UrlMgr) Start(path string) (err error) {
    logInfo("UrlMgr:Start, start")

    //初始化，启动记录线程
    self.taskList = list.New()
    self.urlChan = make(chan string, 1024)
    self.finishChan = make(chan struct{})

    //加载已经完成的url记录文件，并初始化完成url的map
    //liujia: os.Open just open file for reading, for writing it need os.OpenFile
    self.doneFile, err = os.Open(path)
    if err != nil {
        logInfof("UrlMgr:Start, open file failed, %v", err)
        return err
    }

    l, err := ReadFileLines(self.doneFile)
    if err != nil {
        logInfof("UrlMgr:Start, read file failed, %v", err)
        return err
    }

    defer l.Init()

    self.doneMap = make(map[string]int64)
    for elem := l.Front(); elem != nil; elem = elem.Next() {
        self.doneMap[elem.Value.(string)] = 1
    }

    self.doneRecord = true
    go self.run()

    return nil
}

func (self *UrlMgr) Stop()  {
    logInfo("UrlMgr:Stop, to stop...")

    if self.doneRecord {
        //self.finishChan <- struct{}{}
        close(self.finishChan)
        close(self.urlChan)
    }
}

func (self *UrlMgr) Push(task *DownTask) {
    if task == nil || !task.Valid() {
        return
    }

    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()
    self.taskList.PushBack(task)

    if self.doneRecord {
        if _, exist := self.doneMap[task.url]; exist {
            return
        }
    }
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

func (self *UrlMgr) PopAll() (tasks []*DownTask, err error) {
    self.taskListLock.Lock()
    defer self.taskListLock.Unlock()

    if self.taskList.Len() > 0 {
        tasks = make([]*DownTask, 0)
        for elem := self.taskList.Front(); elem != nil; elem = elem.Next() {
            tasks = append(tasks, elem.Value.(*DownTask))
        }
        self.taskList.Init()
        return tasks, nil
    }else {
        return nil, errors.New("no url at this time")
    }
}

func (self *UrlMgr) UrlDone(handledUrl string) {
    if self.doneRecord {
        self.doneMap[handledUrl] = TimeMillSecond()
    }
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
