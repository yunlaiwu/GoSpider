package main

import (
	"errors"
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
	//下面这个是哥自己的douban的cookie，直接整体设置进去，上面那种设置方式太麻烦了
	req.Header.Set("cookie", `bid=IxTlxUVtHMw; gr_user_id=92c3b064-ccee-4842-809d-5ed41b30c3f4; ll="108288"; viewed="1083428_26794880_25723658_1856494_27104286_27104764_27139942_20413077_3288908_27071483"; CNZZDATA1258025570=2002166084-1511423754-https%253A%252F%252Fmovie.douban.com%252F%7C1511423754; cn_d6168da03fa1ahcc4e86_dplus=%7B%22distinct_id%22%3A%20%2215fe7f629d22-000ff054c0eb0c-31627c01-13c680-15fe7f629d3840%22%2C%22sp%22%3A%20%7B%22%24id%22%3A%20%2236749039%22%2C%22%24_sessionid%22%3A%200%2C%22%24_sessionTime%22%3A%201511425192%2C%22%24dp%22%3A%200%2C%22%24_sessionPVTime%22%3A%201511425192%7D%2C%22initial_view_time%22%3A%20%221511423754%22%2C%22initial_referrer%22%3A%20%22https%3A%2F%2Fmovie.douban.com%2Fsubject_search%3Fsearch_text%3D%25E7%25BC%2598%25E5%2588%2586%26cat%3D1002%22%2C%22initial_referrer_domain%22%3A%20%22movie.douban.com%22%7D; UM_distinctid=15fe7f629d22-000ff054c0eb0c-31627c01-13c680-15fe7f629d3840; gsScrollPos-2969=; gsScrollPos-3069=0; gsScrollPos-5492=0; ps=y; ct=y; gsScrollPos-6398=; _ga=GA1.2.1311487014.1496308178; _gid=GA1.2.1317096402.1512630351; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1512642587%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; __utmt=1; ue="liujia.gl@gmail.com"; dbcl2="36749039:Hnxw5oUC3x0"; ck=htMd; _vwo_uuid_v2=65AC86A9AECDE5E405C81E6DF843AE62|147624e7df109db1d8d120243538ddf0; _pk_id.100001.4cf6=9f25967c1f8f7d09.1507316253.36.1512642765.1512640174.; _pk_ses.100001.4cf6=*; __utma=30149280.1311487014.1496308178.1512638020.1512642587.102; __utmb=30149280.0.10.1512642587; __utmc=30149280; __utmz=30149280.1512630215.100.30.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; __utmv=30149280.3674; __utma=223695111.534550196.1503552685.1512638038.1512642587.36; __utmb=223695111.7.10.1512642587; __utmc=223695111; __utmz=223695111.1512638038.35.14.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; push_noty_num=0; push_doumail_num=0; ap=1`)

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
		//拿不到ip就报错重试
		//return HttpGet(destUrl)
		return nil, errors.New("no proxy")
	}

	proxyFunc := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxy)
	}

	transport := &http.Transport{
		Proxy: proxyFunc,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
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
