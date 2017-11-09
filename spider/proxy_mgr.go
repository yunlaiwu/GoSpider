package main

import (
    "encoding/json"
    "time"
)

//这个是代理那边的请求地址
const (
    GET_PROXY_ADDRESS = "http://api.xdaili.cn/xdaili-api//privateProxy/applyStaticProxy?spiderId=43c5671e16f3401aa11911d91dfb0c0c&returnType=2&count=4"
)

//代理那边返回的json结构
type ProxyItem struct {
    IP      string  `json:"ip"`
    Port    string  `json:"port"`
}

type JsonItem struct {
    Code    string       `json:"ERRORCODE"`
    Result  []ProxyItem  `json:"RESULT"`
}

type JsonFailedItem struct {
    Code    string       `json:"ERRORCODE"`
    Result  string       `json:"RESULT"`
}

/*
type JsonItem struct {
    Code    string       `json:"ERRORCODE"`
    Result  []struct{
        IP      string  `json:"ip"`
        Port    string  `json:"port"`
    }  `json:"RESULT"`
}
*/

type ProxyMgr struct {
    proxyList   []string
    proxyIndex  int
    getChan     chan struct{}
    retChan     chan string
    finishChan  chan struct{}
}

func NewProxyMgr() *ProxyMgr  {
    return &ProxyMgr {
        proxyList:  make([]string, 0),
        proxyIndex: 0,
        finishChan: make(chan struct{}),
        getChan:  make(chan struct{}),
        retChan:  make(chan string),
    }
}

func (self *ProxyMgr) Start()  {
    logInfo("ProxyMgr:Start, start")
    //启动就更新proxy，防止启动顺序不一致，导致启动时拿不到
    self.updateProxy()
    go self.run()
}

func (self *ProxyMgr) Stop()  {
    logInfo("ProxyMgr:Stop, to stop...")
    self.finishChan <- struct{}{}
    close(self.finishChan)
    close(self.getChan)
    close(self.retChan)
}

func (self * ProxyMgr) Get() string {
    self.getChan <- struct{}{}
    return <-self.retChan
}

func (self *ProxyMgr) run()  {
    for {
        select {
        case <-self.finishChan:
            return
        case  <-self.getChan:
            self.retChan <- self.getProxy()
        case <-time.After(15 * time.Second):
            self.updateProxy()
        }
    }
}

func (self *ProxyMgr) updateProxy() {
    if IsMAC() {
        logInfof("DO NOT update proxy on MAC!")
        return
    }

    resp, err := HttpGet(GET_PROXY_ADDRESS)
    if err != nil {
        logErrorf("failed to update proxy, http get failed %v", err)
        return
    }

    var proxy JsonItem
    if err := json.Unmarshal(resp, &proxy); err != nil {
        logErrorf("failed to update proxy, decode json failed, %v", err)
        var msg JsonFailedItem
        if err = json.Unmarshal(resp, &msg); err != nil {
            logErrorf("also failed to decode error msg, %v", err)
        }else {
            logErrorf("request failed, err:%v, msg:%v", msg.Code, msg.Result)
        }
        return
    }

    if proxy.Code != "0" {
        logErrorf("failed to update proxy, return code %v", proxy.Code)
        return
    }

    self.proxyList = make([]string, 0)
    for _, item := range proxy.Result {
        if len(item.IP) > 0 && len(item.Port) > 0 {
            self.proxyList = append(self.proxyList, item.IP + ":" + item.Port)
        }
    }

    if len(self.proxyList) == 0 {
        logError("failed to update proxy, no proxy updated")
        return
    }

    self.proxyIndex = 0
    logInfof("update proxy success, total %v", len(self.proxyList))
}

func (self *ProxyMgr) getIndex() int {
    index := self.proxyIndex
    self.proxyIndex++
    if self.proxyIndex >= len(self.proxyList) {
        self.proxyIndex = 0
    }
    return index
}

func (self *ProxyMgr) getProxy() string {
    if len(self.proxyList) == 0 {
        return ""
    }
    return self.proxyList[self.getIndex()]
}