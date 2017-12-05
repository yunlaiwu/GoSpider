package main

import (
	"fmt"
	"strconv"

	"github.com/opesun/goquery"
)

/*MovieCommentData 电影短评数据结构 */
type MovieCommentData struct {
	CommentID   string `json:"cid"`
	UserId      string `json:"userid"`
	UserName    string `json:"username"`
	UserPage    string `json:"userpage"`
	UserAvatar  string `json:"useravatar"`
	PublishDate string `json:"publish"`
	Rate        int    `json:"rate"`
	Useful      int    `json:"useful"`
	Content     string `json:"content"`
}

/*NewMovieCommentData 生成MovieCommentData */
func NewMovieCommentData() *MovieCommentData {
	return &MovieCommentData{}
}

/*String String()接口实现 */
func (self MovieCommentData) String() string {
	return fmt.Sprintf("电影短评%v: 用户:%v|%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v, 内容:%v", self.CommentID, self.UserName, self.UserId, self.UserPage, self.UserAvatar, self.PublishDate, self.Rate, self.Useful, self.Content)
	/*
		    ss := []string{self.commentid, self.userid, self.username, self.userpage, self.publish_date, Int2String(self.rate), Int2String(self.useful), self.content}
			for index, s := range ss {
				ss[index] = SanityString(s)
			}
		    return strings.Join(ss, "\t")
	*/
}

/*ParseMovieComment
 *从豆瓣电影的短评列表页中解析所有短评，https://movie.douban.com/subject/20495023/comments?start=0&limit=20&sort=time&status=P&percent_type=
 */
func ParseMovieComment(htm string) (comments []*MovieCommentData, err error) {
	nodes, err := goquery.ParseString(htm)
	if err != nil {
		fmt.Println("ParseMovieComment: failed parse html")
		return comments, err
	}

	commentsNodes := nodes.Find(".comment-list")

	//comment id
	ids := make([]string, 0)
	commentsNodes.Find(".comment-item").Each(func(index int, item *goquery.Node) {
		for _, attr := range item.Attr {
			if attr.Key == "data-cid" {
				ids = append(ids, attr.Val)
			}
		}
	})

	comments = make([]*MovieCommentData, len(ids))
	for i := range comments {
		comment := NewMovieCommentData()
		comment.CommentID = ids[i]
		comments[i] = comment
	}

	//用户名 用户小站
	commentsNodes.Find(".avatar").Each(func(i int, avatar *goquery.Node) {
		for _, child := range avatar.Child {
			//用户名，用户个人主页
			for _, attr := range child.Attr {
				if attr.Key == "title" && i < len(comments) {
					comments[i].UserName = attr.Val
				} else if attr.Key == "href" && i < len(comments) {
					comments[i].UserPage = attr.Val
				}
			}

			//用户头像
			for _, child2 := range child.Child {
				for _, attr2 := range child2.Attr {
					if attr2.Key == "src" && i < len(comments) {
						comments[i].UserAvatar = attr2.Val
					}
				}
			}
		}
	})

	//有用？
	commentsNodes.Find(".vote-count").Each(func(i int, voteCount *goquery.Node) {
		if len(voteCount.Child) > 0 {
			if num, err := strconv.Atoi(voteCount.Child[0].Data); err == nil {
				comments[i].Useful = num
			}
		}
	})

	//评论内容
	commentsNodes.Find(".comment-content").Each(func(i int, contentNode *goquery.Node) {
		if len(contentNode.Child) > 0 {
			comments[i].Content = contentNode.Child[0].Data
		}
	})

	//几个星？发表时间
	commentsNodes.Find(".comment-info").Each(func(i int, infoNode *goquery.Node) {
		for _, child := range infoNode.Child {
			for _, attr := range child.Attr {
				if attr.Key == "class" {
					comments[i].Rate = ParseRating(attr.Val)
				}
			}

			if child.Data == "span" {
				for _, child2 := range child.Child {
					comments[i].PublishDate = child2.Data
				}
			}
		}
	})

	//从用户主页和用户头像中解析出用户id
	for _, comment := range comments {
		userID, err := ParseUserID(comment.UserAvatar, comment.UserPage)
		if err == nil {
			comment.UserId = userID
		} else {
			logErrorf("failed to get userID, avatar:%v, page:%v", comment.UserAvatar, comment.UserPage)
		}
	}

	return comments, nil
}
