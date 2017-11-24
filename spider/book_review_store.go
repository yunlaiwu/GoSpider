package main

/*
 * 这里和短评的工作方式是完全一样的，唯一就是操作的是BookReview而不是BookComment对象
 */

import (
    "container/list"
    "os"
    "sync"
    "strings"
    "path/filepath"
)

type BookReviewStore struct {
    booksFile string
    saveDir  string

    bookList *list.List
    bookListLock sync.Mutex
    doneMap  sync.Map
    totalCount int
    doneCount int
}

func NewBookReviewStore() *BookReviewStore {
    return &BookReviewStore{
        bookList:list.New(),
        totalCount:0,
        doneCount: 0,
    }
}

func (self *BookReviewStore) Start(booksFile, saveDir string) (err error) {
    logInfo("BookReviewStore:Start, start")
    self.booksFile = booksFile
    self.saveDir = saveDir

    bookFile, err := os.Open(self.booksFile)
    if err != nil {
        logErrorf("BookReviewStore:Start, failed to open booksFile %v", err)
        return err
    }

    defer bookFile.Close()

    lines, err := ReadFileLines(bookFile)
    if err != nil {
        logErrorf("BookReviewStore:Start, failed to read booksFile %v", err)
        return err
    }

    bidDone := loadDoneTask(GetFullPath(filepath.Join(self.saveDir, "./book-review/")))

    for elem := lines.Front(); elem != nil; elem = elem.Next() {
        //每行是用\t分割的 bookID和bookTitle
        parts := strings.Split(elem.Value.(string) , "\t")
        if len(parts) != 2 {
            //report error here
            continue
        }

        if _, exist := bidDone[parts[0]]; !exist {
            self.bookList.PushBack(NewBookReview(parts[0], parts[1], self.saveDir))
        }else {
            logInfof("book review for %v|%v is already downloaded", parts[0], parts[1])
        }
    }

    self.totalCount = self.bookList.Len()

    reviews := self.getReviewTask(3)
    for _, review := range reviews {
        review.Start()
    }

    return nil
}

func (self *BookReviewStore) OnFinished(id string) {
    self.doneCount +=1
    self.doneMap.Store(id, TimeMillSecond())

    logInfof("One Task is Done! downloaded %v resources now", self.doneCount)

    reviews := self.getReviewTask(1)
    for _, review := range reviews {
        review.Start()
    }

    if self.totalCount == self.doneCount {
        //都完成了
        logInfof("All Task is Done! total download %v resources", self.doneCount)
        doneChan <- nil
    }
}

func (self *BookReviewStore) getReviewTask(n int) (reviews []*BookReview) {
    reviews = make([]*BookReview, 0)
    if n < 1 {
        return reviews
    }

    for{
        elem := self.bookList.Front()
        if elem == nil {
            break
        }
        reviews = append(reviews, elem.Value.(*BookReview))
        self.bookList.Remove(elem)
        if len(reviews) == n {
            return reviews
        }
    }

    return reviews
}

