package main

/*
 * 电影短评获取方式
 * 拿到一个movieID后，生成一个MovieComment对象，保存所有的这些MovieComment对象
 * 首先去完成一个MovieComment对象，调用其Start()接口
 * 这个MovieComment对象完成后，会调用OnFinished()接口，我们记录一下(记录个数，并且记录完成的那些ID)，如果都完成了，就退出
 * 否则取一个新的MovieComment，并Start()
 */

import (
	"container/list"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*MovieCommentStore ...*/
type MovieCommentStore struct {
	movieFile string
	saveDir   string

	movieList     *list.List
	movieListLock sync.Mutex
	doneMap       sync.Map
	totalCount    int
	doneCount     int
}

/*NewMovieCommentStore ...*/
func NewMovieCommentStore() *MovieCommentStore {
	return &MovieCommentStore{
		movieList:  list.New(),
		totalCount: 0,
		doneCount:  0,
	}
}

/*Start ...*/
func (self *MovieCommentStore) Start(movieFile, saveDir string) (err error) {
	logInfo("MovieCommentStore:Start, start")
	self.movieFile = movieFile
	self.saveDir = GetFullPath(filepath.Join(saveDir, "./movie-comment/"))
	err = CreateDirIfNotExist(self.saveDir)
	if err != nil {
		logErrorf("MovieCommentStore:saveToFile, failed to create folder %v", self.saveDir)
		return err
	}

	taskFile, err := os.Open(self.movieFile)
	if err != nil {
		logErrorf("MovieCommentStore:Start, failed to open movieFile %v", err)
		return err
	}

	defer taskFile.Close()

	lines, err := ReadFileLines(taskFile)
	if err != nil {
		logErrorf("MovieCommentStore:Start, failed to read movieFile %v", err)
		return err
	}

	midDone := loadDoneTask(self.saveDir)

	for elem := lines.Front(); elem != nil; elem = elem.Next() {
		//每行是用\t分割的 movieID和movieTitle
		parts := strings.Split(elem.Value.(string), "\t")
		if len(parts) != 2 {
			//report error here
			continue
		}

		if _, exist := midDone[parts[0]]; !exist {
			self.movieList.PushBack(NewMovieComment(parts[0], parts[1], self.saveDir))
		} else {
			logInfof("movie comment for %v|%v is already downloaded", parts[0], parts[1])
		}
	}

	self.totalCount = self.movieList.Len()
	if self.totalCount == 0 {
		//啥都没有
		logInfof("NO TASK! total download %v resources", self.doneCount)
		doneChan <- nil
	}

	comments := self.getCommentTask(3)
	for _, comment := range comments {
		comment.Start()
	}

	return nil
}

/*OnFinished ...*/
func (self *MovieCommentStore) OnFinished(id string) {
	self.doneCount++
	self.doneMap.Store(id, TimeMillSecond())

	logInfof("One Task is Done! downloaded %v resources now", self.doneCount)

	comments := self.getCommentTask(1)
	for _, comment := range comments {
		comment.Start()
	}

	if self.totalCount == self.doneCount {
		//都完成了
		logInfof("All Task is Done! total download %v resources", self.doneCount)
		doneChan <- nil
	}
}

func (self *MovieCommentStore) getCommentTask(n int) (comments []*MovieComment) {
	comments = make([]*MovieComment, 0)
	if n < 1 {
		return comments
	}

	for {
		elem := self.movieList.Front()
		if elem == nil {
			break
		}
		comments = append(comments, elem.Value.(*MovieComment))
		self.movieList.Remove(elem)
		if len(comments) == n {
			return comments
		}
	}

	return comments
}
