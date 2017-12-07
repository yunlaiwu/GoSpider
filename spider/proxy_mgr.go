package main

import (
	"encoding/json"
	"time"
)

//这个是代理那边的请求地址
const (
	//原始的代理提供商那里的返回的ip
	GET_PROXY_ADDRESS = "http://api.xdaili.cn/xdaili-api//privateProxy/applyStaticProxy?spiderId=43c5671e16f3401aa11911d91dfb0c0c&returnType=2&count=4"
	//自己内网使用的，封装了上面代理商的
	INTERNAL_GET_PROXY_ADDRESS = "http://10.10.10.149:8099/catcher/proxy?type=douban&rn=5"
)

//代理那边返回的json结构
type ProxyItem struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type JsonItem struct {
	Code   string      `json:"ERRORCODE"`
	Result []ProxyItem `json:"RESULT"`
}

type JsonFailedItem struct {
	Code   string `json:"ERRORCODE"`
	Result string `json:"RESULT"`
}

//内部返回的json结构,{"errno":0,"errmsg":"","data":["58.208.114.47:20596","115.217.255.242:33151","121.206.68.65:20430","115.202.71.160:42858","117.86.8.169:25688"],"s":"c5eb96a418a454796f2b5566b1d1025f1494"} d
type JsonItem2 struct {
	Error int      `json:"error"`
	Msg   string   `json:"errmsg"`
	Data  []string `json:"data"`
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
	proxyList  []string
	proxyIndex int
	getChan    chan struct{}
	retChan    chan string
	finishChan chan struct{}
}

func NewProxyMgr() *ProxyMgr {
	return &ProxyMgr{
		proxyList:  make([]string, 0),
		proxyIndex: 0,
		finishChan: make(chan struct{}),
		getChan:    make(chan struct{}),
		retChan:    make(chan string),
	}
}

func (self *ProxyMgr) Start() {
	logInfo("ProxyMgr:Start, start")
	//启动就更新proxy，防止启动顺序不一致，导致启动时拿不到
	self.updateProxy2()
	go self.run()
}

func (self *ProxyMgr) Stop() {
	logInfo("ProxyMgr:Stop, to stop...")
	self.finishChan <- struct{}{}
	close(self.finishChan)
	close(self.getChan)
	close(self.retChan)
}

func (self *ProxyMgr) Get() string {
	self.getChan <- struct{}{}
	return <-self.retChan
}

func (self *ProxyMgr) run() {
	for {
		select {
		case <-self.finishChan:
			return
		case <-self.getChan:
			self.retChan <- self.getProxy()
		case <-time.After(5 * time.Second):
			self.updateProxy2()
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
		} else {
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
			self.proxyList = append(self.proxyList, item.IP+":"+item.Port)
		}
	}

	if len(self.proxyList) == 0 {
		logError("failed to update proxy, no proxy updated")
		return
	}

	self.proxyIndex = 0
	logInfof("update proxy success, total %v", len(self.proxyList))
}

func (self *ProxyMgr) updateProxy2() {
	if IsMAC() {
		logInfof("DO NOT update proxy on MAC!")
		return
	}

	resp, err := HttpGet(INTERNAL_GET_PROXY_ADDRESS)
	if err != nil {
		logErrorf("failed to update proxy, http get failed %v", err)
		return
	}

	var proxy JsonItem2
	if err := json.Unmarshal(resp, &proxy); err != nil {
		logErrorf("failed to update proxy, decode json failed, %v", err)
		return
	}

	if proxy.Error != 0 {
		logErrorf("failed to update proxy, return %v, %v", proxy.Error, proxy.Msg)
		return
	}

	self.proxyList = make([]string, 0)
	for _, item := range proxy.Data {
		if len(item) > 0 {
			self.proxyList = append(self.proxyList, item)
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
