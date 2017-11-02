package main_test

import (
    "testing"
    spider "GoSpider/spider"
    "fmt"
)

func Test_HttpProxyGet(t *testing.T) {
    proxyMgr := spider.NewProxyMgr()
    proxyMgr.Start()
    proxy := proxyMgr.Get()
    if len(proxy) == 0 {
        t.Errorf("Failed to get proxy")
        t.FailNow()
    }else {
        fmt.Println("Got proxy", proxy)
    }

    ret, err := spider.HttpProxyGet("http://www.baidu.com", proxy)
    if err != nil {
        t.Errorf("Get failed, %v", err)
        fmt.Println("Get failed, ", err)
    }else {
        fmt.Println("Proxy Get success")
        fmt.Println(string(ret))
    }
}