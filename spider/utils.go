package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ReadFileLines(fi *os.File) (ret *list.List, err error) {
	if fi == nil {
		return ret, errors.New("invalid file handle")
	}

	ret = list.New()
	buf := bufio.NewReader(fi)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			//report error
			return ret, err
		}

		ret.PushBack(strings.TrimSpace(line))
	}

	return ret, nil
}

func TimeMillSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func TimeSecond() int64 {
	return time.Now().Unix()
}

func Time2String(ts int64) string {
	var secs int64 = ts / 1000
	var left = ts % 1000
	var s string = time.Unix(secs, 0).Format("2006-01-02 15:04:05")
	return s + fmt.Sprintf(":%v", left)
}

//得到当前工作目录，其实就是当前进程的work directory
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//得到全路径，如果relPath已经是Absolute的就直接用，否则跟当前工作目录拼成Abs的
func GetFullPath(relPath string) (fullPath string) {
	if filepath.IsAbs(relPath) == false {
		return filepath.Join(GetCurrentDirectory(), relPath)
	}
	return relPath
}

/*FileExist does 检查文件是否存在*/
func FileExist(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() == false {
			return true
		}
	}
	return false
}

/*DirExist does 文件夹是否存在*/
func DirExist(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() == true {
			return true
		}
	}
	return false
}

//不存在则创建
func CreateDirIfNotExist(path string) (err error) {
	if DirExist(path) == false {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * 获取当前的goroutine的id，这个俺也是从网上抄的
 */
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

//去掉\t \n \r
func SanityString(s string) string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

//去掉\t \n \r
//去掉其它文件系统不允许的作为文件名的字符：. .. 和空格 和 / \  ？ * & $ 替换为_，
func SanityStringForFileName(s string) string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "..", "-", -1)
	s = strings.Replace(s, ".", "-", -1)
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "\\", "-", -1)
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "?", "-", -1)
	s = strings.Replace(s, "*", "-", -1)
	s = strings.Replace(s, "&", "-", -1)
	s = strings.Replace(s, "$", "-", -1)
	return s
}

func Int2String(i int) string {
	return fmt.Sprintf("%v", i)
}

func String2Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IsMAC() bool {
	return runtime.GOOS == "darwin"
}

func loadDoneTask(dirpath string) (ids map[string]bool) {
	ids = make(map[string]bool)
	totalFile := 0
	filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			logErrorf("loadDoneTask, invalid file or directory [%v], %v", f.Name(), err)
			return err
		}
		totalFile += 1
		if strings.HasSuffix(path, ".txt") == false {
			logErrorf("loadDoneTask, not txt file, %v", f.Name())
			return err
		}
		filename := filepath.Base(path)
		//just like 1059419_海边的卡夫卡.txt
		first := strings.Index(filename, "_")
		if first == -1 {
			logErrorf("loadDoneTask, failed to extract id from file %v", f.Name())
			return err
		}

		id := filename[0:first]
		if _, exist := ids[id]; exist {
			logErrorf("loadDoneTask, id %v already exist, %v", id, f.Name())
		} else {
			ids[id] = true
		}

		return nil
	})
	logInfof("loadDoneTask, finally find %v files with correct id", totalFile)
	return ids
}
