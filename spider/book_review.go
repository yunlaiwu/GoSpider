package main

/*
 * 这里对应于一本书的所有书评
 * 书评和短评不一样的地方在于，短评从各个分页列表里可以拿到所有的结果，但书评从分页列表里拿不到最详细的结果(主要是我正则太烂，不想去解析那坨恶心的html了)
 * 所以需要调用一个接口，去获取这个书评的细节(就是书评内容，含有少量的标签，https://book.douban.com/j/review/8883176/full)
 * 所以我们需要一个map，根据reviewId，去缓存一个书评的内容(从分页列表中获取的)，然后细节获取后，补充到这个缓存的对象中(嗯，对象应该是*BOOK_REVIEW)
 *
 * 所以流程是这样的：
 * Start()是入口，首先生成第一页的下载任务，这个任务执行完毕后，在完成回调中，从页面中获取共有多少书评(20个一页)，并计算出总的页数
 * 然后根据总的页数，生成对应的从第二页到最后一页的任务去执行
 * 对于每一页解析出来的BooReview对象，获取其reviewId，然后生成对应的url去拉取评论详情(任务的params参数里记录这两种任务类型，一个是book-review，一个是book-review-detail)
 * 对应分页的任务，跟短评一样，记录完成结果，对应单独的reviewId，也记录
 * 当两个完成了，才是真正完成
 * 当完成后，告知BookReviewStore这本书的短评都下完了(OnFinished()接口)，并且对所有短评记录到文件中
 */

import (
    "sync"
    "strconv"
)

import (
    "fmt"
    "path/filepath"
    "os"
    "strings"
)

const (
    //第一个参数是书的ID，第二个参数注意，这里的start通常是0 20 40 60....
    BOOK_REVIEW_LISTPAGE_URL_FORMAT = "https://book.douban.com/subject/%v/reviews?sort=time&start=%v"
    //参数的reviewID
    BOOK_REVIEW_DETAIL_URL_FORMAT = "https://book.douban.com/j/review/%v/full"
)

type BookReview struct {
    bookId string
    bookTitle string
    baseFolder string

    totalPage   int
    totalReview int
    totalFinishedReview int
    pageMapLock sync.Mutex
    pageMap   map[string]([]string)
    reviewMapLock sync.Mutex
    reviewMap map[string]*BOOK_REVIEW
}

func NewBookReview(bookId, bookTitle string, baseFolder string) *BookReview {
    return &BookReview{
        bookId:bookId,
        bookTitle:bookTitle,
        baseFolder:baseFolder,
        totalPage:-1,
        totalReview:-1,
        totalFinishedReview:-1,
        pageMap:make(map[string]([]string)),
        reviewMap:make(map[string]*BOOK_REVIEW),
    }
}

func (self *BookReview) Start() {
    logInfof("%v|%v, start!", self.bookId, self.bookTitle)
    spe.Register(self.getResId(), self)
    spe.Do(self.getResId(), self.getListPageUrl(1), map[string]string{"id":self.bookId, "title":self.bookTitle, "res":"book-review", "page":strconv.Itoa(1)})
}

func (self *BookReview) checkFinish() {
    self.pageMapLock.Lock()
    defer self.pageMapLock.Unlock()
    if self.totalPage == len(self.pageMap) && self.totalFinishedReview == self.totalReview {
        logInfof("%v|%v, download finished, %v pages with %v reviews", self.bookId, self.bookTitle, self.totalPage, self.totalReview)
        go func() {
            self.onFinished()
        }()
    }
}

func (self *BookReview) OnResponse(url string, resp []byte, params map[string]string) {
    logInfof("BookReview:OnResponse, url:%v, params:%v", url, params)
    if page,exist := params["page"]; exist {
        if page == "1" {
            //第一页解析总的评论数，并计算总的页数
            count, err := ParseTotalReviews(string(resp))
            if err != nil {
                logErrorf("%v|%v, failed to get page count, %v", self.bookId, self.bookTitle, err)
                self.onFinished()
                return
            }
            self.totalPage = (count+19)/20
            logInfof("%v|%v, total page %v", self.bookId, self.bookTitle, self.totalPage)

            for i := 2; i <= self.totalPage; i++ {
                spe.Do(self.getResId(), self.getListPageUrl(i), map[string]string{"id":self.bookId, "title":self.bookTitle, "res":"book-review", "page":strconv.Itoa(i)})
            }
        }

        reviews, err := ParseBookReviewListPage(string(resp))
        if len(reviews) == 0 || err != nil {
            logErrorf("%v|%v, parse html for reviews failed, %v", self.bookId, self.bookTitle, err)
            self.onFinished()
        }else {
            self.addPageReviews(page, reviews)
        }

    }else if reviewId, exist := params["rid"]; exist {
        self.reviewMapLock.Lock()
        defer self.reviewMapLock.Unlock()
        if review, exist := self.reviewMap[reviewId]; exist && review != nil {
            if content, useful, useless, err := ParseReviewJson(resp); err == nil {
                review.Content = content
                review.Useful = useful
                review.Useless = useless
                self.totalFinishedReview +=1

                //检查是否完成
                self.checkFinish()
            }else {
                logErrorf("%v|%v, parse html for review %v failed, %v, %v", self.bookId, self.bookTitle, reviewId, err, string(resp))
                // 豆瓣有可能返回错误信息，由于UA或者访问过多什么原因，这里重试
                respString := string(resp)
                if strings.Contains(respString, "<html") && strings.Contains(respString, "<title>") && strings.Contains(respString, "没有访问权限") {
                    //这种情况是这个书评的详情无权访问，例子：view-source:https://book.douban.com/j/review/5440030/full
                    review.Content = ""
                    review.Useful = 0
                    review.Useless = 0
                    self.totalFinishedReview += 1

                    //检查是否完成
                    self.checkFinish()
                } else {
                    spe.Do(self.getResId(), self.getDetailUrl(review.ReviewId), map[string]string{"bid": self.bookId, "title": self.bookTitle, "res": "book-review", "rid": review.ReviewId})
                }
            }
        }else {
            if !exist {
                logErrorf("%v|%v, reviewId %v not exist", reviewId)
            }else {
                logErrorf("%v|%v, reviewId %v exist but it is nil", reviewId)
            }

            self.onFinished()
        }
    } else {
        logErrorf("%v|%v, param error, no page no rid, %v", self.bookId, self.bookTitle, params)
    }
}

func (self BookReview) getResId() string {
    return RES_BOOK_REVIEW+"-"+self.bookId
}

func (self BookReview) getListPageUrl(page int) (string) {
    return fmt.Sprintf(BOOK_REVIEW_LISTPAGE_URL_FORMAT, self.bookId, (page-1)*20)
}

func (self BookReview) getDetailUrl(rid string) (string) {
    return fmt.Sprintf(BOOK_REVIEW_DETAIL_URL_FORMAT, rid)
}

func (self *BookReview) addPageReviews(page string, reviews []*BOOK_REVIEW) {
    logInfof("BookReview:addPageReviews, add %d reviews for page %v", len(reviews), page)

    self.pageMapLock.Lock()
    self.reviewMapLock.Lock()
    defer self.pageMapLock.Unlock()
    defer self.reviewMapLock.Unlock()

    if _, exist := self.pageMap[page];exist {
        logErrorf("%v|%v, page %v maybe downloaed more than once", self.bookId, self.bookTitle, page)
    }else {
        reviewIds := make([]string, 0)
        for _, review := range reviews {
            reviewIds = append(reviewIds, review.ReviewId)
            self.reviewMap[review.ReviewId] = review
        }
        self.pageMap[page] = reviewIds
        self.totalReview += len(reviews)
    }

    for _, review := range reviews {
        spe.Do(self.getResId(), self.getDetailUrl(review.ReviewId), map[string]string{"bid":self.bookId, "title":self.bookTitle, "res":"book-review", "rid":review.ReviewId})
    }
}

func (self BookReview) onFinished() {
    self.saveToFile()
    storeMgr.OnFinished(self.bookId)
}

func (self BookReview) saveToFile() error {
    fullpath := GetFullPath(filepath.Join(self.baseFolder, "./book-review/"))
    err := CreateDirIfNotExist(fullpath)
    if err != nil {
        logErrorf("BookReview:saveToFile, failed to create folder %v", fullpath)
        return err
    }

    fullfile := filepath.Join(fullpath, SanityString(self.bookId + "_" + self.bookTitle + ".txt"))
    f, err := os.OpenFile(fullfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }

    defer f.Close()

    for i:=1; i<=self.totalPage; i++ {
        if reviewIds, exist := self.pageMap[fmt.Sprintf("%v", i)]; exist {
            for _, reviewId := range reviewIds {
                if review, exist := self.reviewMap[reviewId]; exist {
                    jstr, err := review.Json()
                    if err == nil {
                        f.WriteString(SanityString(jstr) + "\n")
                    }else {
                        logErrorf("BookReview:saveToFile, failed to marshal to json, reviewId %v", reviewId)
                    }
                }
            }
        }
    }

    logInfof("BookReview:saveToFile, %v|%v, save to file %v successfully, totally %v reviews in %v pages", self.bookId, self.bookTitle, fullfile, self.totalFinishedReview, self.totalPage)
    return nil
}
