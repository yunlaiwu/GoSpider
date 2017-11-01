package main


type Callback interface {
    OnFinished(handledUrl string, success bool) (newUrl string)
}

type SpiderEngine struct{
    ProxyMgr    *ProxyMgr
    urlMgr      *UrlMgr
}

