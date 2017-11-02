package main

import (
    "container/list"
    "os"
    "sync"
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

func (self *BookStore) Start(doneFilePath string) {

}