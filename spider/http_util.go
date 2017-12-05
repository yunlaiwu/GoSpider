package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func HttpGet(url string) (ret []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	/*
		req.SetBasicAuth("test", "123456")
		cookie := &http.Cookie{
			Name:  "test",
			Value: "12",
		}
		req.AddCookie(cookie)
	*/

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36") //设定ua

	//resp, err := http.Get(url)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logError("HTTP GET failed,", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logError("Read Response Body failed,", err)
		return nil, err
	}

	return body, nil
}

func HttpProxyGet(destUrl, proxy string) (ret []byte, err error) {
	if len(proxy) == 0 {
		return HttpGet(destUrl)
	}

	proxyFunc := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxy)
	}

	transport := &http.Transport{
		Proxy: proxyFunc,
		Dial: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(destUrl)
	if err != nil {
		logError("HTTP Proxy GET failed,", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logError("Read Response Body failed,", err)
		return nil, err
	}

	return body, nil
}

func HttpProxyGet2(destUrl, proxy string) (ret []byte, err error) {
	if len(proxy) == 0 {
		return HttpGet(destUrl)
	}

	proxyFunc := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxy)
	}

	transport := &http.Transport{
		Proxy: proxyFunc,
		Dial: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(destUrl)
	if err != nil {
		logError("HTTP Proxy GET failed,", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logError("Read Response Body failed,", err)
		return nil, err
	}

	return body, nil
}
