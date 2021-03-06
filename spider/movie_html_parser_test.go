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
*/

/*
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
*/

/*
func Test_ParseMovieReviewListPage(t *testing.T) {
	//htm, err := spider.HttpGet("https://movie.douban.com/subject/27000061/reviews?sort=time&start=0")
	htm, err := spider.HttpGet("https://movie.douban.com/subject/27034748/reviews?sort=time&start=0")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	reviews, err := spider.ParseMovieReviewListPage(string(htm))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("found %v movie reviews\n", len(reviews))
	for _, review := range reviews {
		fmt.Println(review)
		//fmt.Println(review.ToJson())
	}
}
*/

/*
func Test_ParseMovieReviewDetailPage(t *testing.T) {
	details := make(map[string]string)
	//details["https://movie.douban.com/review/8832330/"] = "这是某个没有太多支线的3A大作录像"
	//details["https://movie.douban.com/review/8894086/"] = "或许你会想了解的电影里的墨西哥元素"
	//details["https://movie.douban.com/review/7553013/"] = "混蛋"
	//details["https://movie.douban.com/review/7659921/"] = "也没有那么差劲嘛"
	details["https://movie.douban.com/review/7657695/"] = "如果你有看过前面007的全部电影，你会会心一笑"

	for url, title := range details {
		htm, err := spider.HttpGet(url)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		fmt.Println(string(htm))
		detail, err := spider.ParseMovieReviewDetailPage(string(htm))
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		} else {
			t.Log(title)
			t.Log(detail)
		}
	}
}
*/

/*
func Test_ParseTotalReviews(t *testing.T) {
	htm, err := spider.HttpGet("https://movie.douban.com/subject/10345617/reviews?sort=time&start=0")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	//fmt.Println(string(htm))
	total, err := spider.ParseTotalReviews(string(htm))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	} else {
		fmt.Println("total:", total)
	}
}
*/

func Test_GetPageTitle(t *testing.T) {
	details := make(map[string]string)
	//details["https://movie.douban.com/review/8832330/"] = "这是某个没有太多支线的3A大作录像"
	//details["https://movie.douban.com/review/8894086/"] = "或许你会想了解的电影里的墨西哥元素"
	//details["https://movie.douban.com/review/7553013/"] = "混蛋"
	//details["https://movie.douban.com/review/7659921/"] = "也没有那么差劲嘛"
	details["https://movie.douban.com/review/7657695/"] = "如果你有看过前面007的全部电影，你会会心一笑"

	for url, _ := range details {
		htm, err := spider.HttpGet(url)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		title, err := spider.GetPageTitle(string(htm))
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		} else {
			t.Log(title)
		}
	}
}
