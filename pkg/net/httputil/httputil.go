package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxConnsPerHost:       100,
	MaxIdleConnsPerHost:   50,
}

// SendRequest 发送http请求
func SendRequest(method string, reqUrl string, header map[string]string, reqBody []byte) (statusCode int, respBody []byte, err error) {
	client := &http.Client{
		Timeout:   120 * time.Second,
		Transport: DefaultTransport,
	}
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return 400, nil, fmt.Errorf("生成请求失败 url:%v err:%v", reqUrl, err)
	}
	if header != nil {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return 400, nil, fmt.Errorf("发送请求失败 url:%v err:%v", reqUrl, err)
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("发送请求成功，读取响应结果失败 url:%v err:%v", reqUrl, err)
	}
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("发送请求成功，响应失败 url:%v statusCode:%v resp:%v", reqUrl, resp.StatusCode, utils.BytesToString(respBody))
	// }
	return http.StatusOK, respBody, nil
}
