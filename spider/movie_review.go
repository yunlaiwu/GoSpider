package main

/*
 * 这里对应于一部电影的所有影评
 * 总体上，和书评差不多，都是先从影评列表页的首页获取总的影评个数以便计算出总的分页，然后遍历每一页列表页。
 * 不一样的就是，获取影评详情时，图书是从一个json接口，电影我不知道有没有这个接口，所以是解析影评详情页获取评论内容
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
	//MOVIE_REVIEW_LISTPAGE_URL_FORMAT 第一个参数是书电影的mid，第二个参数注意，这里的start通常是0 20 40 60....
	MOVIE_REVIEW_LISTPAGE_URL_FORMAT = "https://movie.douban.com/subject/%v/reviews?sort=time&start=%v"
	//MOVIE_REVIEW_DETAIL_URL_FORMAT 参数的reviewID
	MOVIE_REVIEW_DETAIL_URL_FORMAT = "https://movie.douban.com/review/%v/"
)

/*MovieReview ... */
type MovieReview struct {
	movieId    string
	movieTitle string
	baseFolder string

	totalPage           int
	totalReview         int
	totalFinishedReview int
	pageMapLock         sync.Mutex
	pageMap             map[string]([]string)
	reviewMapLock       sync.Mutex
	reviewMap           map[string]*MovieReviewData
}

/*NewMovieReview ... */
func NewMovieReview(movieId, movieTitle string, baseFolder string) *MovieReview {
	return &MovieReview{
		movieId:             movieId,
		movieTitle:          movieTitle,
		baseFolder:          baseFolder,
		totalPage:           -1,
		totalReview:         -1,
		totalFinishedReview: -1,
		pageMap:             make(map[string]([]string)),
		reviewMap:           make(map[string]*MovieReviewData),
	}
}

/*Start ... */
func (self *MovieReview) Start() {
	logInfof("%v|%v, start!", self.movieId, self.movieTitle)
	spe.Register(self.getResId(), self)
	spe.Do(self.getResId(), self.getListPageUrl(1), map[string]string{"mid": self.movieId, "title": self.movieTitle, "res": "movie-review", "page": strconv.Itoa(1)})
}

/*OnFinished ... */
func (self MovieReview) OnFinished() {
	self.saveToFile()
	storeMgr.OnFinished(self.movieId)
}

/*OnResponse ... */
func (self *MovieReview) OnResponse(url string, resp []byte, params map[string]string) {
	logInfof("MovieReview:OnResponse, url:%v, params:%v", url, params)
	if page, exist := params["page"]; exist {
		if page == "1" {
			//第一页解析总的评论数，并计算总的页数
			count, err := ParseTotalReviews(string(resp))
			if err != nil {
				logErrorf("%v|%v, failed to get page count, %v", self.movieId, self.movieTitle, err)
				self.OnFinished()
				return
			}
			self.totalPage = (count + 19) / 20
			logInfof("%v|%v, total page %v", self.movieId, self.movieTitle, self.totalPage)

			for i := 2; i <= self.totalPage; i++ {
				spe.Do(self.getResId(), self.getListPageUrl(i), map[string]string{"mid": self.movieId, "title": self.movieTitle, "res": "movie-review", "page": strconv.Itoa(i)})
			}
		}

		reviews, err := ParseMovieReviewListPage(string(resp))
		if len(reviews) == 0 || err != nil {
			logErrorf("%v|%v, parse html for reviews failed, %v", self.movieId, self.movieTitle, err)
			self.OnFinished()
		} else {
			self.addPageReviews(page, reviews)
		}

	} else if reviewId, exist := params["rid"]; exist {
		self.reviewMapLock.Lock()
		defer self.reviewMapLock.Unlock()
		if review, exist := self.reviewMap[reviewId]; exist && review != nil {
			if content, err := ParseMovieReviewDetailPage(string(resp)); err == nil {
				review.Content = content
				self.totalFinishedReview++

				//检查是否完成
				self.checkFinish()
			} else {
				logErrorf("%v|%v, parse html for review %v failed, %v, %v", self.movieId, self.movieTitle, reviewId, err, string(resp))
				// 豆瓣有可能返回错误信息，由于UA或者访问过多什么原因，这里重试
				respString := string(resp)
				if strings.Contains(respString, "<html") && strings.Contains(respString, "<title>") && strings.Contains(respString, "没有访问权限") {
					//这种情况是这个影评的详情无权访问，书评的例子：https://book.douban.com/j/review/5440030/full
					review.Content = ""
					review.Useful = 0
					review.Useless = 0
					self.totalFinishedReview++

					//检查是否完成
					self.checkFinish()
				} else {
					spe.Do(self.getResId(), self.getDetailUrl(review.ReviewID), map[string]string{"mid": self.movieId, "title": self.movieTitle, "res": "movie-review", "rid": review.ReviewID})
				}
			}
		} else {
			if !exist {
				logErrorf("%v|%v, reviewId %v not exist", reviewId)
			} else {
				logErrorf("%v|%v, reviewId %v exist but it is nil", reviewId)
			}

			self.OnFinished()
		}
	} else {
		logErrorf("%v|%v, param error, no page no rid, %v", self.movieId, self.movieTitle, params)
	}
}

func (self MovieReview) getMovieId() string {
	return self.movieId
}

func (self MovieReview) getMovieTitle() string {
	return self.movieTitle
}

func (self MovieReview) getResId() string {
	return RES_MOVIE_REVIEW + "-" + self.movieId
}

func (self MovieReview) getListPageUrl(page int) string {
	return fmt.Sprintf(MOVIE_REVIEW_LISTPAGE_URL_FORMAT, self.movieId, (page-1)*20)
}

func (self MovieReview) getDetailUrl(rid string) string {
	return fmt.Sprintf(MOVIE_REVIEW_DETAIL_URL_FORMAT, rid)
}

func (self *MovieReview) checkFinish() {
	self.pageMapLock.Lock()
	defer self.pageMapLock.Unlock()
	if self.totalPage == len(self.pageMap) && self.totalFinishedReview == self.totalReview {
		logInfof("%v|%v, download finished, %v pages with %v reviews", self.movieId, self.movieTitle, self.totalPage, self.totalReview)
		go func() {
			self.OnFinished()
		}()
	}
}

func (self *MovieReview) addPageReviews(page string, reviews []*MovieReviewData) {
	logInfof("MovieReview:addPageReviews, add %d reviews for page %v", len(reviews), page)

	self.pageMapLock.Lock()
	self.reviewMapLock.Lock()
	defer self.pageMapLock.Unlock()
	defer self.reviewMapLock.Unlock()

	if _, exist := self.pageMap[page]; exist {
		logErrorf("%v|%v, page %v maybe downloaed more than once", self.movieId, self.movieTitle, page)
	} else {
		reviewIds := make([]string, 0)
		for _, review := range reviews {
			reviewIds = append(reviewIds, review.ReviewID)
			self.reviewMap[review.ReviewID] = review
		}
		self.pageMap[page] = reviewIds
		self.totalReview += len(reviews)
	}

	for _, review := range reviews {
		spe.Do(self.getResId(), self.getDetailUrl(review.ReviewID), map[string]string{"mid": self.movieId, "title": self.movieTitle, "res": "movie-review", "rid": review.ReviewID})
	}
}

func (self MovieReview) saveToFile() error {
	fullpath := GetFullPath(self.baseFolder)
	err := CreateDirIfNotExist(fullpath)
	if err != nil {
		logErrorf("MovieReview:saveToFile, failed to create folder %v, err:", fullpath, err)
		return err
	}

	fullfile := filepath.Join(fullpath, SanityStringForFileName(self.movieId+"_"+self.movieTitle)+".txt")
	f, err := os.OpenFile(fullfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logErrorf("MovieReview:saveToFile, failed to create file %v, err:", fullpath, err)
		return err
	}

	defer f.Close()

	for i := 1; i <= self.totalPage; i++ {
		if reviewIDs, exist := self.pageMap[fmt.Sprintf("%v", i)]; exist {
			for _, reviewID := range reviewIDs {
				if review, exist := self.reviewMap[reviewID]; exist {
					jstr, err := review.ToJson()
					if err == nil {
						f.WriteString(SanityString(jstr) + "\n")
					} else {
						logErrorf("MovieReview:saveToFile, failed to marshal to json, reviewId %v", reviewID)
					}
				}
			}
		}
	}

	logInfof("MovieReview:saveToFile, %v|%v, save to file %v successfully, totally %v reviews in %v pages", self.movieId, self.movieTitle, fullfile, self.totalFinishedReview, self.totalPage)
	return nil
}
