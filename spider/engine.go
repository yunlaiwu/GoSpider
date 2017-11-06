package main

import (
    //"strings"
    "sync"
)

type SpiderEngine struct{
    proxyMgr    *ProxyMgr
    urlMgr      *UrlMgr
    workPool    *WorkerPool

    callbacks map[string]ICallbacker

    cond        *sync.Cond
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
        callbacks: make(map[string]ICallbacker),
        cond: sync.NewCond(new(sync.Mutex)),
    }
}

func (self *SpiderEngine) Start() {
    self.urlMgr.Start("")
    self.proxyMgr.Start()

    go self.run()
}

func (self *SpiderEngine) Register(res string, callbacker ICallbacker) {
    if len(res) == 0 && callbacker == nil {
        return
    }
    self.callbacks[res] = callbacker
}

func (self *SpiderEngine) Do(res string, url string, params map[string]string) {
    if cb, ok := self.callbacks[res]; ok {
        self.urlMgr.Push(NewDownTask(res, url, cb, params))
        self.cond.Signal()
    }
}

func (self *SpiderEngine) run() {
    self.cond.L.Lock()
    if self.exit {
        self.cond.L.Unlock()
        return
    }else {
        self.cond.Wait()
        if task, err := self.urlMgr.Pop(); err != nil {
            self.workPool.Put(func(){
                resp, err := HttpProxyGet(task.url, self.proxyMgr.getProxy())
                if err != nil {
                    //report error
                }else {
                    task.cb.OnResponse(task.url, resp, task.params)
                }
            })
        }
    }
}