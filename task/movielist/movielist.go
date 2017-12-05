package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

//文件是否存在
func FileExist(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() == false {
			return true
		}
	}
	return false
}

//去掉\t \n \r
func SanityString(s string) string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

type DMovieJson struct {
	Mid   string `json:"mid"`
	Url   string `json:"durl"`
	Title string `json:"title"`
}

func main() {
	//解析命令行
	ifile := flag.String("f", "", "input movie file")
	ofile := flag.String("o", "", "output movielist file")
	num := flag.Int("n", 0, "number of item")
	dfile := flag.String("d", "duplicated.txt", "output duplicated item file")

	flag.Parse()

	moviefile := *ifile
	outfile := *ofile
	n := *num
	dupfile := *dfile

	//检查命令行
	if FileExist(moviefile) == false {
		fmt.Println("moviefile not exit", moviefile)
		return
	}

	fi, err := os.OpenFile(moviefile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("failed to open moviefile", moviefile)
		return
	}

	defer fi.Close()

	fo, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("failed to create output file", outfile)
		return
	}

	defer fo.Close()

	fd, err := os.OpenFile(dupfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("failed to create duplicated item file", dupfile)
		return
	}

	defer fd.Close()

	idm := make(map[string]int)

	buf := bufio.NewReader(fi)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("failed to read file", err)
			break
		}

		var movie DMovieJson
		if err := json.Unmarshal([]byte(line), &movie); err != nil {
			fmt.Println("failed to decode json", err)
		} else {
			movieID := SanityString(movie.Mid)
			if _, exist := idm[movieID]; exist {
				s := fmt.Sprintf("found duplicated item, id %v, title %v\n", movieID, movie.Title)
				fd.WriteString(s)
			} else {
				s := fmt.Sprintf("%v\t%v\n", movieID, SanityString(movie.Title))
				fo.WriteString(s)
				idm[movieID] = 1
				if len(idm) == n {
					fmt.Printf("meet number %v, quit\n", n)
					break
				}
			}
		}
	}

	fo.Sync()
	fmt.Println("execute successfully")
}
