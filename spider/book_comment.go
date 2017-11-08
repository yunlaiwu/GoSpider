package main

/*
 * 这里对应于一本书的所有短评
 * Start()是入口，首先生成第一页的下载任务，这个任务执行完毕后，在完成回调中，从页面中获取共有多少短评(20个一页)，并计算出总的页数
 * 然后根据总的页数，生成对应的从第二页到最后一页的任务去执行
 * 对每个任务，都记录去完成的结果，当完成后，告知BookCommentStore这本书的短评都下完了(OnFinished()接口)，并且对所有短评记录到文件中
 */

import (
    "sync"
    "strconv"
)

import (
    "fmt"
    "path/filepath"
    "os"
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
    spe.Do(self.getResId(), self.getUrl(1), map[string]string{"bid":self.bookId, "title":self.bookTitle, "res":"book-comment", "page":strconv.Itoa(1)})
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
                spe.Do(self.getResId(), self.getUrl(i), map[string]string{"bid":self.bookId, "title":self.bookTitle, "res":"book-comment", "page":strconv.Itoa(i)})
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
        go func() {
            self.saveToFile()
            storeMgr.OnFinished(self.bookId)
            logInfof("%v|%v, write to file finished", self.bookId, self.bookTitle)
        }()
    }
}

func (self BookComment) saveToFile() error {
    fullpath := GetFullPath(filepath.Join(self.baseFolder, "./book-comment/"))
    err := CreateDirIfNotExist(fullpath)
    if err != nil {
        logErrorf("BookComment:saveToFile, failed to create folder %v", fullpath)
        return err
    }

    fullfile := filepath.Join(fullpath, SanityString(self.bookId + "_" + self.bookTitle + ".txt"))
    f, err := os.OpenFile(fullfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }

    defer f.Close()

    for i:=1; i<=self.totalPage; i++ {
        v, ok := self.pageMap.Load(fmt.Sprintf("%v", i))
        if ok {
            comments := v.([]*BOOK_COMMENT)
            for _, comment := range comments {
                f.WriteString(comment.String() + "\n")
            }
        }
    }

    logErrorf("BookComment:saveToFile, save to file %v successfully", fullfile)
    return nil
}