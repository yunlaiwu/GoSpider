package main

type BookInfo struct {
    bookId      string
    bookTitle   string
}

func NewBookInfo(id, title string) *BookInfo {
    return &BookInfo{
        bookId:id,
        bookTitle:title,
    }
}

func (self BookInfo) GetId() string {
    return self.bookId
}

func (self BookInfo) GetTitle() string {
    return self.bookTitle
}
