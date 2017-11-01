package main

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "fmt"
)

func HttpGet(url string) (ret []byte, err error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
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

    transport := &http.Transport{Proxy: proxyFunc}
    client := &http.Client{Transport: transport}

    resp, err := client.Get(destUrl)
    if err != nil {
        fmt.Println("Get failed", err)
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Read Response Body failed", err)
        return nil, err
    }

    return body, nil
}