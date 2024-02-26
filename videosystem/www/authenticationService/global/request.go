package global

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

/*
封装远程请求结构体以及方法
*/

type RequestClinet struct {
	Client *http.Client
}

// 创建一个新的 RequestClient 实例， 配置超时
func NewRequestClient(timeout time.Duration) *RequestClinet {
	return &RequestClinet{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *RequestClinet) DoRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		go SendLogs("error", "创建请求错误", err)
		return nil, err
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发出请求
	resp, err := r.Client.Do(req)
	if err != nil {
		go SendLogs("error", "发送请求错误", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		go SendLogs("error", "解析response错误", err)
		return nil, err
	}
	return responseBody, nil
}
