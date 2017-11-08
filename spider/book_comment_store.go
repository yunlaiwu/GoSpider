package main

/*
 * 短评获取方式
 * 拿到一个bookId后，生成一个BookComment对象，保存所有的这些BookComment对象
 * 首先去完成一个BookComment对象，调用其Start()接口
 * 这个BookComment对象完成后，会调用OnFinished()接口，我们记录一下(记录个数，并且记录完成的那些ID)
 * 然后取一个新的BookComment，并Start()
 */

import (
    "container/list"
    "os"
    "sync"
    "strings"
)

type BookCommentStore struct {
    booksFile string
    saveDir  string

    bookList *list.List
    bookListLock sync.Mutex
    doneMap  sync.Map
    totalCount int
    doneCount int
}

func NewBookCommentStore() *BookCommentStore {
    return &BookCommentStore{
        bookList:list.New(),
        totalCount:0,
        doneCount: 0,
    }
}

func (self *BookCommentStore) Start(booksFile, saveDir string) (err error) {
    logInfo("BookCommentStore:Start, start")
    self.booksFile = booksFile
    self.saveDir = saveDir

    bookFile, err := os.Open(self.booksFile)
    if err != nil {
        logErrorf("BookCommentStore:Start, failed to open booksFile %v", err)
        return err
    }

    defer bookFile.Close()

    lines, err := ReadFileLines(bookFile)
    if err != nil {
        logErrorf("BookCommentStore:Start, failed to read booksFile %v", err)
        return err
    }

    for elem := lines.Front(); elem != nil; elem = elem.Next() {
        //每行是用\t分割的 bookID和bookTitle
        parts := strings.Split(elem.Value.(string) , "\t")
        if len(parts) != 2 {
            //report error here
            continue
        }

        self.bookList.PushBack(NewBookComment(parts[0], parts[1], self.saveDir))
    }

    self.totalCount = self.bookList.Len()

    comments := self.getCommentTask(3)
    for _, comment := range comments {
        comment.Start()
    }

    return nil
}

func (self *BookCommentStore) OnFinished(id string) {
    self.doneCount +=1
    self.doneMap.Store(id, TimeMillSecond())

    logInfof("One Task is Done! downloaded %v resources now", self.doneCount)

    comments := self.getCommentTask(1)
    for _, comment := range comments {
        comment.Start()
    }

    if self.totalCount == self.doneCount {
        //都完成了
        logInfof("All Task is Done! total download %v resources", self.doneCount)
        doneChan <- nil
    }
}

func (self *BookCommentStore) getCommentTask(n int) (comments []*BookComment) {
    comments = make([]*BookComment, 0)
    if n < 1 {
        return comments
    }

    for{
        elem := self.bookList.Front()
        if elem == nil {
            break
        }
        comments = append(comments, elem.Value.(*BookComment))
        self.bookList.Remove(elem)
        if len(comments) == n {
            return comments
        }
    }

    return comments
}
