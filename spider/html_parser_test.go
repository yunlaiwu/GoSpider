package main_test

import (
    "testing"
    spider "GoSpider/spider"
)

func Test_ParseBookComment(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/comments/new")
    if err != nil {
        t.FailNow()
    }

    spider.ParseBookComment(string(htm))
}

func Test_ParseBookReview(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/reviews?sort=time")
    if err != nil {
        t.FailNow()
    }

    spider.ParseBookReview(string(htm))
}