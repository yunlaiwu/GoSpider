package main

/*
 * 这里和短评的工作方式是完全一样的，唯一就是操作的是MovieReview而不是MovieComment对象
 */

import (
	"container/list"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*MovieReviewStore ... */
type MovieReviewStore struct {
	moviesFile string
	saveDir    string

	movieList     *list.List
	movieListLock sync.Mutex
	doneMap       sync.Map
	totalCount    int
	doneCount     int
}

/*NewMovieReviewStore ... */
func NewMovieReviewStore() *MovieReviewStore {
	return &MovieReviewStore{
		movieList:  list.New(),
		totalCount: 0,
		doneCount:  0,
	}
}

/*Start ... */
func (self *MovieReviewStore) Start(moviesFile, saveDir string) (err error) {
	logInfo("MovieReviewStore:Start, start")
	self.moviesFile = moviesFile
	self.saveDir = GetFullPath(filepath.Join(saveDir, "./movie-review/"))

	movieFile, err := os.Open(self.moviesFile)
	if err != nil {
		logErrorf("MovieReviewStore:Start, failed to open moviesFile %v", err)
		return err
	}

	defer movieFile.Close()

	lines, err := ReadFileLines(movieFile)
	if err != nil {
		logErrorf("MovieReviewStore:Start, failed to read moviesFile %v", err)
		return err
	} else {
		logInfof("we got %v movies listed in config file", lines.Len())
	}
	logInfof("we got %v movies in task file", lines.Len())

	midDone := loadDoneTask(self.saveDir)
	logInfof("we got %v movies already downloaded", len(midDone))

	for elem := lines.Front(); elem != nil; elem = elem.Next() {
		//每行是用\t分割的 movieID和movieTitle
		parts := strings.Split(elem.Value.(string), "\t")
		if len(parts) != 2 {
			//report error here
			logErrorf("invalid line in task file, %v", elem.Value.(string))
			continue
		}

		if _, exist := midDone[parts[0]]; !exist {
			self.movieList.PushBack(NewMovieReview(parts[0], parts[1], self.saveDir))
		} else {
			logInfof("movie review for %v|%v is already downloaded", parts[0], parts[1])
		}
	}

	self.totalCount = self.movieList.Len()
	logInfof("we got %v movies to download this time", self.totalCount)
	for elem := self.movieList.Front(); elem != nil; elem = elem.Next() {
		review := elem.Value.(*MovieReview)
		logInfof("movie review for %v|%v need to download", review.getMovieId(), review.getMovieTitle())
	}

	reviews := self.getReviewTask(3)
	if len(reviews) == 0 {
		logInfof("no task, exit!")
		doneChan <- nil
	} else {
		for _, review := range reviews {
			review.Start()
		}
	}

	return nil
}

/*OnFinished ... */
func (self *MovieReviewStore) OnFinished(id string) {
	self.doneCount += 1
	self.doneMap.Store(id, TimeMillSecond())

	logInfof("One Task is Done! Already downloaded %v movies now", self.doneCount)

	if self.totalCount == self.doneCount {
		//都完成了
		logInfof("All Task is Done! total download %v resources", self.doneCount)
		doneChan <- nil

		//only for debug
		logInfof("Check movie list....")
		for elem := self.movieList.Front(); elem != nil; elem = elem.Next() {
			review := elem.Value.(*MovieReview)
			logInfof("movie review for %v|%v still in movie list", review.getMovieId(), review.getMovieTitle())
		}

	} else {
		reviews := self.getReviewTask(1)
		for _, review := range reviews {
			review.Start()
		}
	}
}

func (self *MovieReviewStore) getReviewTask(n int) (reviews []*MovieReview) {
	reviews = make([]*MovieReview, 0)
	if n < 1 {
		return reviews
	}

	for {
		elem := self.movieList.Front()
		if elem == nil {
			break
		}
		reviews = append(reviews, elem.Value.(*MovieReview))
		self.movieList.Remove(elem)
		if len(reviews) == n {
			return reviews
		}
	}

	return reviews
}
