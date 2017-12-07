package main_test

/*
func Test_ParseBookComment(t *testing.T) {
	htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/comments/new")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	comments, err := spider.ParseBookComment(string(htm))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, comment := range comments {
		fmt.Println(comment)
		fmt.Println("")
	}
}
*/

/*
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
*/

/*
func Test_ParseBookReview2(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1083428/reviews?sort=time")
    if err != nil {
        t.FailNow()
    }

    reviews, err := spider.ParseBookReviewListPage(string(htm))
    if err != nil {
        t.FailNow()
    }

    for _, review := range reviews {
        fmt.Println(review.GetId())
    }

    reviewId := reviews[0].GetId()
    reviewUrl := fmt.Sprintf("https://book.douban.com/j/review/%v/full", reviewId)
    fmt.Println("review full url", reviewUrl)
    htm, err = spider.HttpGet(reviewUrl)
    if err != nil {
        t.FailNow()
    }
    spider.ParseBookReviewJson(htm)
}
*/

/*
func Test_ParseBookTotalComments(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/27104286/comments/new?p=1")
    if err != nil {
        t.FailNow()
    }

    count, err := spider.ParseTotalComments(string(htm))
    if err != nil {
        t.FailNow()
    }else {
        fmt.Println(count)
    }
}
*/

/*
func Test_ParseBookTotalReviews(t *testing.T) {
    htm, err := spider.HttpGet("https://book.douban.com/subject/1000323/reviews?sort=time&start=0")
    if err != nil {
        t.FailNow()
    }

    count, err := spider.ParseTotalReviews(string(htm))
    if err != nil {
        t.FailNow()
    }else {
        fmt.Println(count)
    }
}
*/
/*
func Test_ParseBookReviewFullJson(t *testing.T) {
    //htm, err := spider.HttpGet("https://book.douban.com/j/review/6653108/full")
    //htm, err := spider.HttpGet("https://book.douban.com/j/review/7636302/full")
    htm, err := spider.HttpGet("https://book.douban.com/j/review/5440030/full")
    if err != nil {
        t.FailNow()
    }

    fmt.Println()
    fmt.Println(string(htm))
    fmt.Println()

    content, useful, useless, err := spider.ParseBookReviewJson(htm)
    if err != nil {
        fmt.Println("parse json failed", err)
        return
    }
    fmt.Println("content:", content)
    fmt.Println("usefull:", useful)
    fmt.Println("useless:", useless)
}
*/
