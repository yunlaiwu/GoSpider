package main

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
    "net"
)

func HttpGet(url string) (ret []byte, err error) {
    resp, err := http.Get(url)
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
        Dial:(&net.Dialer{
                Timeout : 15 * time.Second,
                KeepAlive : 15 * time.Second,
            }).Dial,
        TLSHandshakeTimeout: 5 * time.Second,
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