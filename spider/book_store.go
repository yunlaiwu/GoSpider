package main

import (
    "container/list"
    "os"
    "sync"
    "strings"
)

type BookStore struct {
    bookList *list.List
    doneFile *os.File
    doneMap  sync.Map

    fullReviewMap sync.Map
}

func NewBookStore() *BookStore {
    return &BookStore{}
}

func (self *BookStore) Start(bookListFile, doneFilePath string) (err error) {
    bookFile, err := os.Open(bookListFile)
    if err != nil {
        return err
    }

    defer bookFile.Close()

    lines, err := ReadFileLines(bookFile)
    if err != nil {
        return err
    }

    for elem := lines.Front(); elem != nil; elem = elem.Next() {
        parts := strings.Split(elem.Value.(string) , "\t")
        if len(parts) != 2 {
            //report error here
            continue
        }
        self.doneMap.Store(parts[0], parts[1])
    }

    return nil
}
