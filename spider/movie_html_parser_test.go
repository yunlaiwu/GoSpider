package main_test

import (
	spider "GoSpider/spider"
	"fmt"
	"testing"
)

/*
func Test_ParseMovieCommentNumber(t *testing.T) {
	htm, err := spider.HttpGet("https://movie.douban.com/subject/20495023/comments?start=0&limit=20&sort=time&status=P&percent_type=")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	count, err := spider.ParseTotalMovieCommentsForWatched(string(htm))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("movie comment count: %v\n", count)
}
*/

func Test_ParseNextMovieCommentListPage(t *testing.T) {
	urls := make(map[string]bool)
	urls["https://movie.douban.com/subject/27034748/comments?start=0&limit=20&sort=time&status=P&percent_type="] = true
	urls["https://movie.douban.com/subject/27034748/comments?start=20&limit=20&sort=time&status=P&percent_type="] = false
	urls["https://movie.douban.com/subject/26737068/comments?sort=time&status=P"] = false

	for url, hasNext := range urls {
		htm, err := spider.HttpGet(url)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		next, err := spider.ParseNextMovieCommentListPage(string(htm))
		if err != nil && !hasNext {
			fmt.Printf("has no next page for url %v \n", url)
		} else if err == nil && hasNext {
			fmt.Printf("has next page of %v for url %v \n", next, url)
		} else {
			fmt.Println("failed!")
			t.FailNow()
		}
	}
}

func Test_ParseMovieComment(t *testing.T) {
	htm, err := spider.HttpGet("https://movie.douban.com/subject/27034748/comments?start=0&limit=20&sort=time&status=P&percent_type=")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	comments, err := spider.ParseMovieComment(string(htm))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("found %v comments\n", len(comments))
	for _, comment := range comments {
		fmt.Println(comment)
		fmt.Println(comment.ToJson())
	}
}
