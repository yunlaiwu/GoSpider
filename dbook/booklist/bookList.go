package main

import (
    "fmt"
    "flag"
    "os"
    "bufio"
    "io"
    "strings"
    "encoding/json"
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

//从图书的url中解析图书的ID，https://book.douban.com/subject/2194049/
func DecodeId(durl string) string {
    durl = strings.Replace(durl, "https://book.douban.com/subject/", "", -1)
    durl = strings.Replace(durl, "/", "", -1)
    return durl
}

//去掉\t \n \r
func SanityString(s string) string {
    s = strings.Replace(s, "\t", "", -1)
    s = strings.Replace(s, "\n", "", -1)
    s = strings.Replace(s, "\r", "", -1)
    return s
}

type DBookJson struct {
    Url    string       `json:"durl"`
    Title  string       `json:"title"`
}

func main() {
    //解析命令行
    ifile := flag.String("f", "", "input book file")
    ofile := flag.String("o", "", "output booklist file")

    flag.Parse()

    bookfile := *ifile
    outfile := *ofile

    //检查命令行
    if FileExist(bookfile) == false {
        fmt.Println("bookfile not exit", bookfile)
        return
    }

    fi, err := os.OpenFile(bookfile, os.O_RDONLY, 0666)
    if err != nil {
        fmt.Println("failed to open bookfile", bookfile)
        return
    }

    defer fi.Close()

    fo, err := os.OpenFile(outfile, os.O_WRONLY | os.O_CREATE, 0666)
    if err != nil {
        fmt.Println("failed to create output file", outfile)
        return
    }

    defer fo.Close()

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

        var book DBookJson
        if err := json.Unmarshal([]byte(line), &book); err != nil {
            fmt.Println("failed to decode json", err)
        }else {
            bookId := DecodeId(book.Url)
            s := fmt.Sprintf("%v\t%v\n", SanityString(bookId), SanityString(book.Title))
            fmt.Println(book.Title, book.Url)
            fo.WriteString(s)
        }
    }

    fo.Sync()
    fmt.Println("execute successfully")
}


