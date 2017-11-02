package main

import (
    "github.com/opesun/goquery"
    "fmt"
    "strconv"
    "strings"
)

type BOOK_COMMENT struct {
    username        string
    userid          string
    userpage        string
    useravatar      string
    publish_date    string
    publish_time    int64
    rate            int
    content         string
    useful          int
}

type BOOK_REVIEW struct {
    review_id       string
    username        string
    userid          string
    userpage        string
    useravatar      string
    publish_date    string
    publish_time    int64
    rate            int
    title           string
    content         string
    useful          int
    useless         int
}

func NewBookComment() *BOOK_COMMENT {
    return &BOOK_COMMENT{}
}

func NewBookReview() *BOOK_REVIEW {
    return &BOOK_REVIEW{}
}

func (self BOOK_COMMENT) String() string {
    return fmt.Sprintf("短评: 用户:%v|%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v, 内容:%v", self.username, self.userid, self.userpage, self.useravatar, self.publish_date, self.rate, self.useful, self.content)
}

func (self BOOK_REVIEW) String() string {
    return fmt.Sprintf("书评 %v: 用户:%v|%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v, 标题:%v, 内容:%v", self.review_id, self.username, self.userid, self.userpage, self.useravatar, self.publish_date, self.rate, self.useful, self.title, self.content)
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

func GetUserIdFromUserPage(r string) string {
    //like "https://www.douban.com/people/48942518/"， so it get "48942518"
    r = strings.ToLower(r)
    r = strings.Replace(r, "https://www.douban.com/people/", "", -1)
    r = strings.Replace(r, "http://www.douban.com/people/", "", -1)
    r = strings.TrimSpace(r)
    r = strings.Trim(r, "/")
    return r
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
        comment := NewBookComment()
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

            for _, child2 := range child.Child {
                for _, attr2 := range child2.Attr {
                    if attr2.Key == "src" && i < len(comments) {
                        comments[i].useravatar = attr2.Val
                    }
                }
            }
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

    for _, comment := range comments {
        fmt.Println(comment)
    }

    return comments, nil
}

/*
 * 豆瓣图书的书评，
 * 这个url获取列表，但评论是折叠的，https://book.douban.com/subject/1083428/reviews?sort=time
 * 从这个获取详情，https://book.douban.com/j/review/8883176/full
 */
func ParseBookReview(htm string) (reviews []*BOOK_REVIEW, err error) {
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
        review := NewBookReview()
        review.review_id = reviewIds[i]
        reviews[i] = review
    }

    //评论title
    reviewNodes.Find(".title-link").Each(func(i int, item *goquery.Node) {
        if len(item.Child) > 0 {
            reviews[i].title = item.Child[0].Data
        }
    })

    /*
    reviewNodes.Find(".author-avatar").Each(func(i int, item *goquery.Node) {
        for _, child := range item.Child {
            if child.Data == "img" {
                for _, attr := range child.Attr {
                    if attr.Key == "src" {
                        reviews[i].useravatar = attr.Val
                    }
                }
            }
        }
    })

    //用户id,用户名，评分，发布时间
    reviewNodes.Find(".author").Each(func(i int, item *goquery.Node) {
        for _, attr := range item.Attr {
            if attr.Key == "href" {
                reviews[i].userpage = attr.Val
                reviews[i].userid = GetUserIdFromUserPage(attr.Val)
            }
        }

        for _, child := range item.Child {
            if child.Data == "span" {
                if len(child.Child) > 0 {
                    reviews[i].username = child.Child[0].Data
                }
            }
        }
    })
    */

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
                }else if class == "author" {
                    reviews[i].userpage = href
                    reviews[i].userid = GetUserIdFromUserPage(href)
                    for _, child2 := range child.Child {
                        if child2.Data == "span" {
                            if len(child2.Child) > 0 {
                                reviews[i].username = child2.Child[0].Data
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
                    reviews[i].rate = ParseRating(class)
                }else if property == "v:dtreviewed" && class == "main-meta" {
                    if len(child.Child) > 0 {
                        reviews[i].publish_date = child.Child[0].Data
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