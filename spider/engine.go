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
    exit        bool
}

func NewSpiderEngine() *SpiderEngine {
    proxyMgr := NewProxyMgr()
    urlMgr := NewUrlMgr()
    workPool := NewWorkerPool(5)

    return &SpiderEngine{
        exit: false,
        proxyMgr: proxyMgr,
        urlMgr: urlMgr,
        workPool: workPool,
        callbacks: make(map[string]IResHunter),
        //cond: sync.NewCond(new(sync.Mutex)),
    }
}

func (self *SpiderEngine) Start() {
    self.urlMgr.Start("")
    self.proxyMgr.Start()

    go self.run()
}

func (self *SpiderEngine) Stop() {
    logInfo("SpiderEngine:Stop, to stop...")
    self.exit = true
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

func (self *SpiderEngine) run() {
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
        case <-time.After(1 * time.Second):
            if self.exit {
                logInfo("SpiderEngine:run, recv exit signal!")
                return
            }
            if tasks, err := self.urlMgr.PopAll(); err == nil {
                for _, task_  := range tasks {
                    task := task_
                    logInfof("SpiderEngine:run, got task, %v", task)
                    self.workPool.Put(func(){
                        resp, err := HttpProxyGet(task.url, self.proxyMgr.getProxy())
                        if err != nil {
                            logInfof("SpiderEngine:run, failed to do HTTP GET for %v, %v", task.url, err)
                        }else {
                            logInfof("SpiderEngine:run, recv HTTP resp for %v, len %v", task.url, len(resp))
                            task.cb.OnResponse(task.url, resp, task.params)
                        }
                    })
                }
            }
        }
    }
}