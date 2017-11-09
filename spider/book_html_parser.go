package main

import (
    "github.com/opesun/goquery"
    "fmt"
    "strconv"
    "strings"
    "encoding/json"
    "errors"
)

type BOOK_COMMENT struct {
    userid          string
    username        string
    userpage        string
    publish_date    string
    publish_time    int64
    rate            int
    useful          int
    content         string
}

type BOOK_REVIEW struct {
    ReviewId     string     `json:"rid"`
    UserName     string     `json:"username"`
    UserId       string     `json:"userid"`
    UserPage     string     `json:"userpage"`
    ReviewTitle  string     `json:"title"`
    PublishDate  string     `json:"publish"`
    Rate         int        `json:"rate"`
    Content      string     `json:"content"`
    Useful       int        `json:"useful"`
    Useless      int        `json:"useless"`
}

func NewBOOK_COMMENT() *BOOK_COMMENT {
    return &BOOK_COMMENT{}
}

func COMMENT_FROM_STRING(s string) (bookComment *BOOK_COMMENT) {
    s = strings.Trim(s, "\n")
    s = strings.TrimSpace(s)
    parts := strings.Split(s, "\t")
    if len(parts) < 7 {
        return nil
    }

    bookComment = NewBOOK_COMMENT()
    bookComment.userid = parts[0]
    bookComment.username = parts[1]
    bookComment.userpage = parts[2]
    bookComment.publish_date = parts[3]
    bookComment.rate = String2Int(parts[4])
    bookComment.useful = String2Int(parts[5])
    bookComment.content = parts[6]

    return bookComment
}

func (self BOOK_COMMENT) String() string {
    //return fmt.Sprintf("短评: 用户:%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v, 内容:%v", self.username, self.userid, self.userpage, self.publish_date, self.rate, self.useful, self.content)
    ss := []string{self.userid, self.username, self.userpage, self.publish_date, Int2String(self.rate), Int2String(self.useful), self.content}
    for index, s := range ss {
        ss[index] = SanityString(s)
    }
    return strings.Join(ss, "\t")
}

func NewBOOK_REVIEW() *BOOK_REVIEW {
    return &BOOK_REVIEW{}
}

func (self BOOK_REVIEW) String() string {
    return fmt.Sprintf("书评 %v: 用户:%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v|%v, 标题:%v, 内容:%v", self.ReviewId, self.UserName, self.UserId, self.UserPage, self.PublishDate, self.Rate, self.Useful, self.Useless,self.ReviewTitle, self.Content)
}

func (self BOOK_REVIEW) Json() (string, error) {
    j, err := json.Marshal(self)
    if err != nil {
        return "", err
    }
    return string(j), nil
}

func (self BOOK_REVIEW) GetId() string {
    return self.ReviewId
}

func ParseRating(r string) int {
    //短评是 "user-stars allstar40 rating"， so it get 40 and return as integer
    //书评是 "allstar50 main-title-rating"
    r = strings.ToLower(r)
    r = strings.Replace(r, "main-title-rating", "", -1)
    r = strings.Replace(r, "user-stars", "", -1)
    r = strings.Replace(r, "allstar", "", -1)
    r = strings.Replace(r, "rating", "", -1)
    r = strings.TrimSpace(r)
    rate, err := strconv.Atoi(r)
    if err == nil {
        return rate
    }else {
        fmt.Println(r)
        return 0
    }
}

func ParseUserIdFromUserPage(r string) string {
    //like "https://www.douban.com/people/48942518/"， so it get "48942518"
    r = strings.ToLower(r)
    r = strings.Replace(r, "https://www.douban.com/people/", "", -1)
    r = strings.Replace(r, "http://www.douban.com/people/", "", -1)
    r = strings.TrimSpace(r)
    r = strings.Trim(r, "/")
    return r
}

func ParseUseful(r string) int {
    //like "有用 0" "没用 0"
    r = strings.Replace(r, "有用", "", -1)
    r = strings.Replace(r, "没用", "", -1)
    r = strings.TrimSpace(r)
    if count, err := strconv.Atoi(r); err == nil {
        return count
    }

    return 0
}

func ParseReviewCount(s string) (count int, err error) {
    //like "电视人的书评 (42)"
    parts1 := strings.Split(s, "(")
    if len(parts1) < 2 {
        return 0, errors.New("failed to find \"(\" ....")
    }

    parts2 := strings.Split(parts1[len(parts1)-1], ")")
    if len(parts2) < 2 {
        return 0, errors.New("failed to find \")\" ....")
    }

    return strconv.Atoi(parts2[0])
}

/*
 * 豆瓣图书的短评，https://book.douban.com/subject/1083428/comments/new
 */
func ParseBookComment(htm string) (comments []*BOOK_COMMENT, err error) {
    nodes, err := goquery.ParseString(htm)
    if err != nil {
        fmt.Println("ParseBookComment: failed parse html")
        return comments, err
    }

    commentsNodes := nodes.Find(".comment-list")

    //用户id
    userids := make([]string, 0)
    commentsNodes.Find(".comment-item").Each(func(index int, item *goquery.Node) {
        for _, attr := range item.Attr {
            if attr.Key == "data-cid" {
                userids = append(userids, attr.Val)
            }
        }
    })

    comments = make([]*BOOK_COMMENT, len(userids))
    for i, _ := range comments {
        comment := NewBOOK_COMMENT()
        comment.userid = userids[i]
        comments[i] = comment
    }

    //用户名 用户小站
    commentsNodes.Find(".avatar").Each(func(i int, avatar *goquery.Node) {
        for _, child := range avatar.Child {
            for _, attr := range child.Attr {
                if attr.Key == "title" && i < len(comments) {
                    comments[i].username = attr.Val
                }else if attr.Key == "href" && i < len(comments) {
                    comments[i].userpage = attr.Val
                }
            }

            /*
            //用户头像暂时没用
            for _, child2 := range child.Child {
                for _, attr2 := range child2.Attr {
                    if attr2.Key == "src" && i < len(comments) {
                        comments[i].useravatar = attr2.Val
                    }
                }
            }
            */
        }
    })

    //有用？
    commentsNodes.Find(".vote-count").Each(func(i int, voteCount *goquery.Node) {
        if len(voteCount.Child) > 0 {
            if num, err := strconv.Atoi(voteCount.Child[0].Data); err == nil {
                comments[i].useful = num
            }
        }
    })

    //评论内容
    commentsNodes.Find(".comment-content").Each(func(i int, contentNode *goquery.Node) {
        if len(contentNode.Child) > 0 {
            comments[i].content = contentNode.Child[0].Data
        }
    })

    //几个星？发表时间
    commentsNodes.Find(".comment-info").Each(func(i int, infoNode *goquery.Node) {
        for _, child := range infoNode.Child {
            for _, attr := range child.Attr {
                if attr.Key == "class" {
                    comments[i].rate = ParseRating(attr.Val)
                }
            }

            if child.Data == "span" {
                for _, child2 := range child.Child {
                    comments[i].publish_date = child2.Data
                }
            }
        }
    })

    /*
    for _, comment := range comments {
        fmt.Println(comment)
    }*/

    return comments, nil
}

/*
 * 豆瓣图书的书评的列表页，https://book.douban.com/subject/1083428/reviews?sort=time
 * 这个url获取列表，但评论是折叠的而且是不完全的
 * 从这个获取详情，https://book.douban.com/j/review/8883176/full
 * liujia: 为了少缓存东西，从书评的分页page中获取每个书评的id，然后从书评详情中抓取
 */
func ParseBookReviewListPage(htm string) (reviews []*BOOK_REVIEW, err error) {
    nodes, err := goquery.ParseString(htm)
    if err != nil {
        fmt.Println("ParseBookReview: failed parse html")
        return reviews, err
    }

    reviewNodes := nodes.Find(".review-list")

    //用户id
    reviewIds := make([]string, 0)
    reviewNodes.Find(".main").Each(func(index int, item *goquery.Node) {
        for _, attr := range item.Attr {
            if attr.Key == "id" {
                reviewIds = append(reviewIds, attr.Val)
            }
        }
    })

    reviews = make([]*BOOK_REVIEW, len(reviewIds))
    for i := range reviews {
        review := NewBOOK_REVIEW()
        review.ReviewId = reviewIds[i]
        reviews[i] = review
    }

    //评论title
    reviewNodes.Find(".title-link").Each(func(i int, item *goquery.Node) {
        if len(item.Child) > 0 {
            reviews[i].ReviewTitle = item.Child[0].Data
        }
    })

    //用户id，头像，个人站点，名字，发布时间，评分
    reviewNodes.Find(".header-more").Each(func(i int, item *goquery.Node) {
        for _, child := range item.Child {
            if child.Data == "a" {
                class := ""
                href := ""
                for _, attr := range child.Attr {
                    if attr.Key == "class" {
                        class = attr.Val
                    }else if attr.Key == "href" {
                        href = attr.Val
                    }
                }

                /*
                if class == "author-avatar" {
                    for _, child2 := range child.Child {
                        if child2.Data == "img" {
                            for _, attr := range child2.Attr {
                                if attr.Key == "src" {
                                    reviews[i].useravatar = attr.Val
                                }
                            }
                        }
                    }
                }else*/ if class == "author" {
                    reviews[i].UserPage = href
                    reviews[i].UserId = ParseUserIdFromUserPage(href)
                    for _, child2 := range child.Child {
                        if child2.Data == "span" {
                            if len(child2.Child) > 0 {
                                reviews[i].UserName = child2.Child[0].Data
                            }
                        }
                    }
                }
            }else if child.Data == "span" {
                property := ""
                class := ""
                for _, attr := range child.Attr {
                    if attr.Key == "property" {
                        property = attr.Val
                    }else if attr.Key == "class" {
                        class = attr.Val
                    }
                }

                if property == "v:rating" {
                    reviews[i].Rate = ParseRating(class)
                }else if property == "v:dtreviewed" && class == "main-meta" {
                    if len(child.Child) > 0 {
                        reviews[i].PublishDate = child.Child[0].Data
                    }
                }
            }
        }
    })

    for _, comment := range reviews {
        fmt.Println(comment)
    }

    return reviews, nil
}


/*
 * 豆瓣图书的书评详情，https://book.douban.com/review/8857096/#comments(https://book.douban.com/review/8857096/也行)
 * 这个url:https://book.douban.com/j/review/8883176/full，拿到json样子的，里面body字段如同上面页面的信息
 * 这个有问题，因为评论内容拿不全
 */
func ParseBookReviewPage(htm string) (bookReview *BOOK_REVIEW, err error) {
    nodes, err := goquery.ParseString(htm)
    if err != nil {
        fmt.Println("ParseBookReview: failed parse html")
        return bookReview, err
    }

    bookReview = NewBOOK_REVIEW()

    //获取reviewId，其实也可以外部填写，不过这里获取一下可以和外部比较，作为是否正确抓取的一个对照
    reviewItemNodes := nodes.Find(".main")
    for _, item := range reviewItemNodes {
        for _, attr := range item.Attr {
            if attr.Key == "id" {
                bookReview.ReviewId = attr.Val
            }
        }
    }

    //书评名
    titleNodes := nodes.Find(".book-content")
    for _, item := range titleNodes {
        for _, child:= range item.Child {
            isTitle := false
            for _, attr := range child.Attr {
                if attr.Key == "id" && attr.Val == "content" {
                    isTitle = true
                }
            }

            if isTitle {
                isTitle = false
                for _, child2 := range child.Child {
                    if child2.Data == "h1" && len(child2.Child) > 0 {
                        for _, child3 := range child2.Child {
                            for _, attr := range child3.Attr {
                                if attr.Key == "property" && attr.Val == "v:summary" {
                                    isTitle = true
                                }
                            }
                            if isTitle && len(child3.Child) > 0 {
                                bookReview.ReviewTitle = child3.Child[0].Data
                            }
                        }
                    }
                }
            }
        }
    }

    //发表书皮用户的id 主页，发布时间
    mainHDNodes := nodes.Find(".main-hd")
	for _, item := range mainHDNodes {
		for _, child := range item.Child {
            for _, attr := range child.Attr {
                if attr.Key == "class" && strings.HasPrefix(attr.Val, "allstar") {
                    bookReview.Rate = ParseRating(attr.Val)
                }else if attr.Key == "class" && attr.Val == "main-meta" {
                    if len(child.Child) > 0 {
                        bookReview.PublishDate = child.Child[0].Data
                    }
                }
            }

			for _, child2 := range child.Child {
				for _, attr := range child2.Attr {
					if attr.Key == "property" && attr.Val == "v:reviewer" {
						if len(child2.Child) > 0 {
							bookReview.UserName = child2.Child[0].Data
							for _, attr := range child.Attr {
								if attr.Key == "href" {
									bookReview.UserPage = attr.Val
									bookReview.UserId = ParseUserIdFromUserPage(attr.Val)
								}
							}
						}
					}
				}
			}
		}
	}

	//内容 (这里还可以拿用户名，不过上面已经拿过了)
    contentNodes := nodes.Find(".review-content")
    for _, item := range contentNodes {
        for _, child := range item.Child {
            bookReview.Content += child.Data
        }
    }

    //有用？无用？
    usefulNodes := nodes.Find(".main-panel-useful")
    for _, item := range usefulNodes {
        for _, child := range item.Child {
            if child.Data == "button" {
                for _, attr := range child.Attr {
                    if attr.Key == "class" && strings.HasPrefix(attr.Val, "btn") && len(child.Child) > 0 {
                        if strings.Contains(attr.Val, "useful_count") {
                            bookReview.Useful = ParseUseful(child.Child[0].Data)
                        }else if strings.Contains(attr.Val, "useless_count") {
                            bookReview.Useless = ParseUseful(child.Child[0].Data)
                        }
                    }
                }
            }
        }
    }

    fmt.Println(bookReview)

	return bookReview, nil
}

func ParseReviewJson(resp []byte) (content string, useful, useless int, err error) {
    type JsonReview struct {
        //Body    string       `json:"body"`
        Votes   struct{
            Useful     int  `json:"useful_count"`
            Useless    int  `json:"useless_count"`
        }  `json:"votes"`
        Html    string      `json:"html"`
    }

    var review JsonReview
    if err := json.Unmarshal(resp, &review); err == nil {
        return review.Html, review.Votes.Useful, review.Votes.Useless, nil
    }else {
        return "", 0, 0, err
    }
}

/*
 * 从图书短评分页列表页，获取总的短评个数
 */
func ParseTotalComments(resp string) (totalComments int, err error) {
    nodes, err := goquery.ParseString(resp)
    if err != nil {
        fmt.Println("ParseTotalComments: failed parse html")
        return 0, err
    }

    commentsNodes := nodes.Find(".comments-wrapper")
    for _, item := range commentsNodes {
        for _, child := range item.Child {
            for _, child2  := range child.Child {
                for _, attr := range child2.Attr {
                    if attr.Key == "id" && attr.Val == "total-comments" {
                        if len(child2.Child) > 0 {
                            s := child2.Child[0].Data
                            s = strings.Replace(s, "全部共", "", -1)
                            s = strings.Replace(s, "条", "", -1)
                            s = strings.TrimSpace(s)
                            return strconv.Atoi(s)
                        }
                    }
                }
            }
        }
    }
    return 0, errors.New("ParseTotalComments: parse html failed, cannnot found")
}

/*
 * 从图书书评分页列表页，获取总的书评个数
 */
func ParseTotalReviews(resp string) (totalComments int, err error) {
    nodes, err := goquery.ParseString(resp)
    if err != nil {
        fmt.Println("ParseTotalReviews: failed parse html")
        return 0, err
    }

    titleNodes := nodes.Find("title")
    for _, item := range titleNodes {
        for _, child := range item.Child {
            return ParseReviewCount(child.Data)
        }
    }

    return 0, nil
}

/*
fmt.Println("info data:", item.Data)
        for _, attr := range item.Attr {
            fmt.Println("info attr:", attr.Key, attr.Val)
        }

        for _, child := range item.Child {
            fmt.Println("info child data:", child.Data)
            for _, attr := range child.Attr {
                fmt.Println("info child attr:", attr.Key, attr.Val)
            }

            for _, child2 := range child.Child {
                fmt.Println("info child2 attr:", child2.Data)
                for _, attr := range child2.Attr {
                    fmt.Println("info child2 attr:", attr.Key, attr.Val)
                }
            }
        }
 */