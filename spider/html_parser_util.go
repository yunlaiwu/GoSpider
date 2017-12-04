package main

import (
	"errors"
	"fmt"
	"strings"
)

/*ParseUserIDFromAvatar does 从用户头像url中获取用户ID*/
func ParseUserIDFromAvatar(url string) (string, error) {
	//用户头像是类似这样：https://img1.doubanio.com/icon/u29691602-37.jpg or https://img1.doubanio.com/icon/user_normal.jpg
	if strings.Contains(url, "user_normal.jpg") {
		return "", errors.New("default user avatar")
	}
	url = strings.ToLower(url)
	url = strings.TrimSpace(url)
	lastSlash := strings.LastIndexByte(url, '/')
	if lastSlash == -1 {
		return "", errors.New("faild to find lastSlash")
	}
	lastShortLine := strings.LastIndexByte(url, '-')
	if lastShortLine == -1 {
		return "", errors.New("faild to find lastShortLine")
	}
	if lastSlash+2 > lastShortLine {
		return "", errors.New("lastSlash should be less than lastDot")
	}
	rs := []rune(url)
	return string(rs[lastSlash+2 : lastShortLine]), nil
}

/*ParseUserIDFromUserPage does 从用户个人主页的url中获取用户ID*/
func ParseUserIDFromUserPage(url string) (string, error) {
	//like "https://www.douban.com/people/48942518/"， so it get "48942518"
	url = strings.ToLower(url)
	url = strings.TrimSpace(url)
	url = strings.Replace(url, "https://www.douban.com/people/", "", -1)
	url = strings.Replace(url, "http://www.douban.com/people/", "", -1)
	url = strings.Trim(url, "/")
	return url, nil
}

/*ParseUserID 根据传入的用户的头像url和用户个人主页url解析出用户ID*/
func ParseUserID(userAvatar, userPage string) (userID string, err error) {
	userID, err = ParseUserIDFromAvatar(userAvatar)
	if err == nil && IsValidUserID(userID) {
		return userID, nil
	}
	userID, err = ParseUserIDFromUserPage(userPage)
	if err == nil && IsValidUserID(userID) {
		return userID, nil
	}
	return "", errors.New("failed to extract userID")
}

/*IsValidUserID 判断是否是正确的用户id，就是判断是否是纯数字字符串*/
func IsValidUserID(id string) bool {
	if len(id) == 0 {
		return false
	}
	digitalChars := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
	}
	for _, c := range id {
		if _, exist := digitalChars[c]; !exist {
			fmt.Printf("%c is not digital", c)
			return false
		}
	}
	return true
}
