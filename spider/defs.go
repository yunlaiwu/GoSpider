package main

/*
 * 回调接口
 */
type ICallbacker interface {
    OnResponse(url string, resp []byte, params map[string]string)
}

/*
 * 下载器，通常是engine实现这个接口，用来下载某个url并回调对应res的callbacker
 * Register注册某类资源的回调
 * Do用来下载某类资源的某个url，param1 param2用来带某些参数
 */
type IDownloader interface {
    Register(res string, callbacker ICallbacker)
    Do(res string, url string, params map[string]string)
}
