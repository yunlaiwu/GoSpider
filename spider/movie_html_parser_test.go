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

/*
func Test_ParseNextMovieCommentListPage(t *testing.T) {
	htm, err := spider.HttpGet("https://movie.douban.com/subject/27034748/comments?start=0&limit=20&sort=time&status=P&percent_type=")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	url, err := spider.ParseNextMovieCommentListPage(string(htm))
	if err != nil {
		fmt.Println("failed, err:", err)
		t.FailNow()
	}else {
		fmt.Println("next url is", url)
	}

	htm, err = spider.HttpGet("https://movie.douban.com/subject/27034748/comments?start=20&limit=20&sort=time&status=P&percent_type=")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	url, err = spider.ParseNextMovieCommentListPage(string(htm))
	if err != nil {
		fmt.Println("failed, err:", err)
		t.FailNow()
	}else {
		fmt.Println("next url is", url)
	}
}
*/

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
	}
}

