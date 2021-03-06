package main

/*
 * 这里对应于一部电影的所有短评
 * Start()是入口，首先生成第一页的下载任务，这个任务执行完毕后，在完成回调中，从页面中获取共有多少短评(20个一页)，并计算出总的页数
 * 然后根据总的页数，生成对应的从第二页到最后一页的任务去执行
 * 对每个任务，都记录去完成的结果，当完成后，告知MovieCommentStore这本书的短评都下完了(OnFinished()接口)，并且对所有短评记录到文件中
 *
 * NOTE：豆瓣电影的短评不同于图书短评，有看过和想看两种短评，分别对应于下面URL的status=P和status=F，下面只考虑了看过
 */

import (
	"strconv"
	"sync"
)

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	MOVIE_COMMENT_FIRST_PAGE_URL_FORMAT = "https://movie.douban.com/subject/%v/comments?start=0&limit=20&sort=time&status=P&percent_type="
	MOVIE_COMMENT_BASE_URL_FORMAT       = "https://movie.douban.com/subject/%v/comments%v"
)

/*MovieComment ...*/
type MovieComment struct {
	movieID    string
	movieTitle string
	baseFolder string

	totalPage int
	pageMap   sync.Map
}

/*NewMovieComment ...*/
func NewMovieComment(movieID, movieTitle string, baseFolder string) *MovieComment {
	return &MovieComment{
		movieID:    movieID,
		movieTitle: movieTitle,
		baseFolder: baseFolder,
		totalPage:  0,
	}
}

/*Start 启动*/
func (self *MovieComment) Start() {
	logInfof("%v|%v, start!", self.movieID, self.movieTitle)
	spe.Register(self.getResId(), self)
	spe.Do(self.getResId(), self.getFirstPageUrl(), map[string]string{"mid": self.movieID, "title": self.movieTitle, "res": "movie-comment", "page": strconv.Itoa(1)})
}

/*OnFinished ...*/
func (self MovieComment) OnFinished() {
	go func() {
		n := 0
		total := 0
		self.pageMap.Range(func(key, value interface{}) bool {
			n++
			total += len(value.([]*MovieCommentData))
			return true
		})

		logInfof("%v|%v, download finished, total %v comments in %v pages ", self.movieID, self.movieTitle, total, n)

		self.saveToFile()
		storeMgr.OnFinished(self.movieID)
	}()
}

/*OnResponse http回调处理，解析出评论*/
func (self *MovieComment) OnResponse(url string, resp []byte, params map[string]string) {
	logInfof("MovieComment:OnResponse, url:%v, params:%v", url, params)
	if p, exist := params["page"]; exist {
		comments, err := ParseMovieComment(string(resp))
		if err != nil {
			logErrorf("%v|%v, parse html for movie comments failed, %v", self.movieID, self.movieTitle, err)
			//重试，不能直接完成，否则后面有些http请求还在队列里，但这里已经标记结束了
			//更好的办法应该是对每页的重试次数进行统计，到了一定的重试后，就算错误
			//self.OnFinished()
			spe.Do(self.getResId(), url, params)
			return
		}

		page, err := strconv.Atoi(p)
		if err != nil {
			logErrorf("%v|%v, parse page %v fail, %v", self.movieID, self.movieTitle, p, err)
			self.OnFinished()
			return
		}

		if len(comments) > 0 {
			self.addComments(page, comments)
		} else {
			logErrorf("%v|%v, no comments in page %v, finish, %v", self.movieID, self.movieTitle, p, err)
			self.OnFinished()
			return
		}

		href, err := ParseNextMovieCommentListPage(string(resp))
		if err != nil || len(href) == 0 || !strings.HasPrefix(href, "?") {
			logErrorf("%v|%v, parse html to get next page fail, %v", self.movieID, self.movieTitle, err)
			self.OnFinished()
			return
		}

		nextUrl := self.getNextPageBaseUrl(href)
		spe.Do(self.getResId(), nextUrl, map[string]string{"mid": self.movieID, "title": self.movieTitle, "res": "movie-comment", "page": strconv.Itoa(page + 1)})
	} else {
		logErrorf("%v|%v, param error %v, no page, retry", self.movieID, self.movieTitle, params)
		spe.Do(self.getResId(), url, params)
	}
}

func (self MovieComment) getResId() string {
	return RES_MOVIE_COMMENT + "-" + self.movieID
}

func (self MovieComment) getFirstPageUrl() string {
	return fmt.Sprintf(MOVIE_COMMENT_FIRST_PAGE_URL_FORMAT, self.movieID)
}

func (self MovieComment) getNextPageBaseUrl(href string) string {
	return fmt.Sprintf(MOVIE_COMMENT_BASE_URL_FORMAT, self.movieID, href)
}

func (self *MovieComment) addComments(page int, comments []*MovieCommentData) {
	logInfof("MovieComment:addComments, add %d comments for page %v", len(comments), page)
	_, loaded := self.pageMap.LoadOrStore(strconv.Itoa(page), comments)
	if loaded == true {
		logErrorf("%v|%v, page %v maybe downloaed more than once", self.movieID, self.movieTitle, page)
	}

	//totalPage保存最大的页码
	if page > self.totalPage {
		self.totalPage = page
	}
}

func (self MovieComment) saveToFile() error {
	fullpath := filepath.Join(self.baseFolder, SanityStringForFileName(self.movieID+"_"+self.movieTitle)+".txt")
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logErrorf("MovieComment:saveToFile, failed to create file %v, err:", fullpath, err)
		return err
	}

	defer f.Close()

	totalComments := 0
	for i := 1; i <= self.totalPage; i++ {
		v, ok := self.pageMap.Load(fmt.Sprintf("%v", i))
		if ok {
			comments := v.([]*MovieCommentData)
			for _, comment := range comments {
				jstr, err := comment.ToJson()
				if err == nil {
					f.WriteString(SanityString(jstr) + "\n")
				} else {
					logErrorf("MovieComment:saveToFile, failed to marshal to json, commentID %v", comment.CommentID)
				}
			}
			totalComments += len(comments)
		}
	}

	logInfof("MovieComment:saveToFile, save to file %v successfully, totally %v comments in %v pages", fullpath, totalComments, self.totalPage)
	return nil
}
