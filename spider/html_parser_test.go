package main_test

import (
    "testing"
    spider "GoSpider/spider"
    "fmt"
)

/*
func Test_ParseBookComment(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/comments/new")
    if err != nil {
        t.FailNow()
    }

    spider.ParseBookComment(string(htm))
}
*/

func Test_ParseBookReview(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/reviews?sort=time")
    if err != nil {
        t.FailNow()
    }

    reviewIds, err := spider.ParseBookReviewListPage(string(htm))
    if err != nil {
        t.FailNow()
    }

    for _,reviewId := range reviewIds {
        fmt.Println(reviewId)
    }

    reviewId := reviewIds[0]
    reviewUrl := fmt.Sprintf("https://book.douban.com/review/%v/#comments", reviewId)
    fmt.Println("review full url", reviewUrl)
    htm, err = spider.HttpGet(reviewUrl)
    if err != nil {
        t.FailNow()
    }
    spider.ParseBookReviewPage(string(htm))
}