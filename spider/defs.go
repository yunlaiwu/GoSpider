package main

const (
    RES_BOOK_COMMENT = "book-comment"
    RES_BOOK_REVIEW = "book-review"
)

/*
 * 管理某个资源，通常是某本书的所有短评，某本书的所有书评等。其实是对应一系列带分页的url
 */
type IResHunter interface {
    OnResponse(url string, resp []byte, params map[string]string)
}

/*
 * 资源管理器，通常管理某一类资源，管理url怎么生成，怎么分页，理解成一类任务也可以
 */
type IResStorer interface {
    Start(resfile string, saveDir string) (err error)
    OnFinished(id string)
}

/*
 * 下载器，通常是engine实现这个接口，用来下载某个url并回调对应res的callbacker
 * Register注册某类资源的回调
 * Do用来下载某类资源的某个url，param1 param2用来带某些参数
 */
type IDownloader interface {
    Register(res string, callbacker IResHunter)
    Do(res string, url string, params map[string]string)
}