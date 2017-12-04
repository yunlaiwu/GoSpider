package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opesun/goquery"
)

/*
 * 电影短评
 */
type MOVIE_COMMENT struct {
	commentid    string
	userid       string
	username     string
	userpage     string
	useravatar   string
	publish_date string
	rate         int
	useful       int
	content      string
}

func NewMOVIE_COMMENT() *MOVIE_COMMENT {
	return &MOVIE_COMMENT{}
}

func MOVIE_COMMENT_FROM_STRING(s string) (movieComment *MOVIE_COMMENT) {
	s = strings.Trim(s, "\n")
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "\t")
	if len(parts) < 8 {
		return nil
	}

	movieComment = NewMOVIE_COMMENT()
	movieComment.commentid = parts[0]
	movieComment.userid = parts[1]
	movieComment.username = parts[2]
	movieComment.userpage = parts[3]
	movieComment.publish_date = parts[4]
	movieComment.rate = String2Int(parts[5])
	movieComment.useful = String2Int(parts[6])
	movieComment.content = parts[7]

	return movieComment
}

func (self MOVIE_COMMENT) String() string {
	//return fmt.Sprintf("电影短评%v: 用户:%v|%v|%v, 发表日期:%v, 评分:%v, 有用:%v, 内容:%v", self.commentid, self.username, self.userid, self.userpage, self.publish_date, self.rate, self.useful, self.content)
	ss := []string{self.commentid, self.userid, self.username, self.userpage, self.publish_date, Int2String(self.rate), Int2String(self.useful), self.content}
	for index, s := range ss {
		ss[index] = SanityString(s)
	}
	return strings.Join(ss, "\t")
}

/*ParseMovieComment
 * 从豆瓣电影的短评列表页中解析所有短评，https://movie.douban.com/subject/20495023/comments?start=0&limit=20&sort=time&status=P&percent_type=
 */
func ParseMovieComment(htm string) (comments []*MOVIE_COMMENT, err error) {
	nodes, err := goquery.ParseString(htm)
	if err != nil {
		fmt.Println("ParseMovieComment: failed parse html")
		return comments, err
	}

	commentsNodes := nodes.Find(".comment-list")

	//用户id
	ids := make([]string, 0)
	commentsNodes.Find(".comment-item").Each(func(index int, item *goquery.Node) {
		for _, attr := range item.Attr {
			if attr.Key == "data-cid" {
				ids = append(ids, attr.Val)
			}
		}
	})

	comments = make([]*MOVIE_COMMENT, len(ids))
	for i, _ := range comments {
		comment := NewMOVIE_COMMENT()
		comment.commentid = ids[i]
		comments[i] = comment
	}

	//用户名 用户小站
	commentsNodes.Find(".avatar").Each(func(i int, avatar *goquery.Node) {
		for _, child := range avatar.Child {
			for _, attr := range child.Attr {
				if attr.Key == "title" && i < len(comments) {
					comments[i].username = attr.Val
				} else if attr.Key == "href" && i < len(comments) {
					comments[i].userpage = attr.Val
				}
			}

			//用户头像暂时没用
			for _, child2 := range child.Child {
				for _, attr2 := range child2.Attr {
					if attr2.Key == "src" && i < len(comments) {
						comments[i].useravatar = attr2.Val
						comments[i].userid, _ = ParseUserIDFromAvatar(attr2.Val)
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

	/*
	   for _, comment := range comments {
	       fmt.Println(comment)
	   }*/

	return comments, nil
}
