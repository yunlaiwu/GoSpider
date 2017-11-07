package main

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
    doneCount int
}

func NewBookCommentStore() *BookCommentStore {
    return &BookCommentStore{
        bookList:list.New(),
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

    comments := self.getComments(3)
    for _, comment := range comments {
        comment.Start()
    }

    return nil
}

func (self *BookCommentStore) OnFinished(id string) {
    self.doneCount +=1
    self.doneMap.Store(id, TimeMillSecond())
    comments := self.getComments(1)
    if len(comments) == 0 {
        //都完成了
        logInfof("All Task is Done! total download %v resources", self.doneCount)
    }else {
        for _, comment := range comments {
            comment.Start()
        }
    }
}

func (self *BookCommentStore) getComments(n int) (comments []*BookComment) {
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
