package main

import (
    //"strings"
    //"sync"
    "time"
)

type SpiderEngine struct{
    proxyMgr    *ProxyMgr
    urlMgr      *UrlMgr
    workPool    *WorkerPool

    callbacks map[string]IResHunter

    //cond        *sync.Cond
    finishChan    chan struct{}
    processChan   chan *DownTask
}

func NewSpiderEngine() *SpiderEngine {
    proxyMgr := NewProxyMgr()
    urlMgr := NewUrlMgr()
    workPool := NewWorkerPool(5)

    return &SpiderEngine{
        proxyMgr: proxyMgr,
        urlMgr: urlMgr,
        workPool: workPool,
        callbacks: make(map[string]IResHunter),
        processChan: make(chan *DownTask, 10000),
        finishChan: make(chan struct{}, 1),
        //cond: sync.NewCond(new(sync.Mutex)),
    }
}

func (self *SpiderEngine) Start() {
    self.urlMgr.Start("")
    self.proxyMgr.Start()

    go self.dispatch()
    go self.process()
}

func (self *SpiderEngine) Stop() {
    logInfo("SpiderEngine:Stop, to stop...")
    close(self.finishChan)
    close(self.processChan)
    //self.cond.Signal()
    self.urlMgr.Stop()
    self.proxyMgr.Stop()
    logInfo("SpiderEngine:Stop, stopped")
}

func (self *SpiderEngine) Register(res string, callbacker IResHunter) {
    if len(res) == 0 && callbacker == nil {
        return
    }
    self.callbacks[res] = callbacker
}

func (self *SpiderEngine) Do(res string, url string, params map[string]string) {
    if cb, ok := self.callbacks[res]; ok {
        task := NewDownTask(res, url, cb, params)
        logInfof("SpiderEngine:Do, new task, %v", task)
        self.urlMgr.Push(task)
        //self.cond.Signal()
    }
}

func (self *SpiderEngine) RetryTask(task *DownTask) {
    if task == nil {
        return
    }
    logInfof("SpiderEngine:RetryTask, retry task %v", task)
    self.urlMgr.Push(task)
}

func (self *SpiderEngine) dispatch() {
    for{
        /*
        self.cond.L.Lock()
        self.cond.Wait()
        if self.exit {
            logInfo("SpiderEngine:run, recv exit signal!")
            self.cond.L.Unlock()
            return
        }
        self.cond.L.Unlock()
        */

        select {
        case <-self.finishChan:
            logInfo("SpiderEngine:dispatch, recv exit signal!")
            return
        case <-time.After(1 * time.Second):
            if tasks, err := self.urlMgr.PopAll(); err == nil {
                for _, task_  := range tasks {
                    task := task_
                    logInfof("SpiderEngine:dispatch, got task, %v", task)
                    self.workPool.Put(func(){
                        resp, err := HttpProxyGet(task.url, self.proxyMgr.getProxy())
                        if err != nil {
                            logWarningf("SpiderEngine:dispatch, failed to do HTTP GET for %v, %v", task.url, err)
                            self.RetryTask(task)
                        }else {
                            logInfof("SpiderEngine:dispatch, recv HTTP resp for %v, len %v", task.url, len(resp))
                            //task.cb.OnResponse(task.url, resp, task.params)
                            task.resp = resp
                            self.processChan <- task
                        }
                    })
                }
            }
        }
    }
}

func (self *SpiderEngine) process() {
    for{
        select {
        case <-self.finishChan:
            logInfo("SpiderEngine:process, recv exit signal!")
            return
        case task := <-self.processChan:
            task.Process()
        }
    }
}