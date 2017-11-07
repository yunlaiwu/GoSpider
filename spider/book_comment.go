package main

import (
    "fmt"
    "strconv"
    "sync"
)

const (
    BOOK_COMMENT_URL_FORMAT = "https://book.douban.com/subject/%v/comments/new?p=%v"
)

type BookComment struct {
    bookId string
    bookTitle string
    baseFolder string

    totalPage int
    pageMap   sync.Map
}

func NewBookComment(bookId, bookTitle string, baseFolder string) *BookComment {
    return &BookComment{
        bookId:bookId,
        bookTitle:bookTitle,
        baseFolder:baseFolder,
        totalPage:-1,
    }
}

func (self *BookComment) Start() {
    logInfof("%v|%v, start!", self.bookId, self.bookTitle)
    spe.Register(self.getResId(), self)
    spe.Do(self.getResId(), self.getUrl(1), map[string]string{"id":self.bookId, "title":self.bookTitle, "res":"book-comment", "page":strconv.Itoa(1)})
}

func (self *BookComment) OnResponse(url string, resp []byte, params map[string]string) {
    logInfof("BookComment:OnResponse, url:%v, params:%v", url, params)
    if page,exist := params["page"]; exist {
        if page == "1" {
            //第一页解析总的评论数，并计算总的页数
            count, err := ParseTotalComments(string(resp))
            if err != nil {
                logErrorf("%v|%v, failed to get page count, %v", self.bookId, self.bookTitle, err)
                storeMgr.OnFinished(self.bookId)
                return
            }
            self.totalPage = (count+19)/20
            logInfof("%v|%v, total page %v", self.bookId, self.bookTitle, self.totalPage)

            for i := 2; i <= self.totalPage; i++ {
                spe.Do(self.getResId(), self.getUrl(i), map[string]string{"id":self.bookId, "title":self.bookTitle, "res":"book-comment", "page":strconv.Itoa(i)})
            }
        }

        comments, err := ParseBookComment(string(resp))
        if len(comments) == 0 || err != nil {
            logErrorf("%v|%v, parse html for comments failed, %v", self.bookId, self.bookTitle, err)
            storeMgr.OnFinished(self.bookId)
        }else {
            self.addComments(page, comments)
        }

    }else {
        logErrorf("%v|%v, param error, no page, %v", self.bookId, self.bookTitle, params)
    }
}

func (self BookComment) getResId() string {
    return RES_BOOK_COMMENT+"-"+self.bookId
}

func (self BookComment) getUrl(page int) (string) {
    return fmt.Sprintf(BOOK_COMMENT_URL_FORMAT, self.bookId, page)
}

func (self *BookComment) addComments(page string, comments []*BOOK_COMMENT) {
    logInfof("BookComment:addComments, add %d comments for page %v", len(comments), page)
    _, loaded := self.pageMap.LoadOrStore(page, comments)
    if loaded == true {
        logErrorf("%v|%v, page %v maybe downloaed more than once", self.bookId, self.bookTitle, page)
    }

    n := 0
    total := 0
    self.pageMap.Range(func(key, value interface{}) bool {
        n +=1
        total += len(value.([]*BOOK_COMMENT))
        return true
    })

    if n == self.totalPage {
        logInfof("%v|%v, download finished, total %v", self.bookId, self.bookTitle, total)
        storeMgr.OnFinished(self.bookId)
    }
}