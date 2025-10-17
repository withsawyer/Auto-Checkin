package util

import (
	"auto-checkin/internal/logger"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestParams struct {
	Method             string
	URL                string
	QueryParams        map[string]string
	BodyData           interface{}
	Headers            map[string]string
	InsecureSkipVerify bool
	Timeout            int
	Proxy              string
}

// createHTTPClient 创建HTTP客户端
func createHTTPClient(insecureSkipVerify bool, timeout int, proxy string) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	if proxy != "" {
		parse, err := url.Parse(proxy)
		if err != nil {
			logger.Log().Errorf("parse proxy error: %v", err)
			return nil
		}
		transport.Proxy = http.ProxyURL(parse)
	}

	client := &http.Client{
		Transport: transport,
	}
	if timeout <= 0 {
		client.Timeout = time.Duration(30) * time.Second
	} else {
		client.Timeout = time.Duration(timeout) * time.Second
	}

	return client
}

// buildURL 构建带查询参数的URL
func buildURL(rawURL string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	if queryParams != nil {
		params := u.Query()
		for key, value := range queryParams {
			params.Add(key, value)
		}
		u.RawQuery = params.Encode()
	}
	return u.String(), nil
}

// createRequestBody 根据BodyData类型创建请求体
func createRequestBody(bodyData interface{}) (io.Reader, error) {
	if bodyData == nil {
		return nil, nil
	}

	switch v := bodyData.(type) {
	case []byte:
		return bytes.NewReader(v), nil
	case string:
		return strings.NewReader(v), nil
	case url.Values:
		return strings.NewReader(v.Encode()), nil
	default:
		marshal, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		return bytes.NewReader(marshal), nil
	}
}

// setContentType 根据BodyData类型设置Content-Type
func setContentType(request *http.Request, bodyData interface{}) {
	if bodyData == nil {
		return
	}

	switch bodyData.(type) {
	case map[string]interface{}, []byte:
		request.Header.Set("Content-Type", "application/json")
	case string:
		request.Header.Set("Content-Type", "text/plain")
	case url.Values:
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

func SendRequest(req *RequestParams) (map[string]interface{}, error) {
	client := createHTTPClient(req.InsecureSkipVerify, req.Timeout, req.Proxy)
	urlWithQuery, err := buildURL(req.URL, req.QueryParams)
	if err != nil {
		return nil, err
	}
	bodyData, err := createRequestBody(req.BodyData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(req.Method, urlWithQuery, bodyData)
	if err != nil {
		return nil, err
	}

	if req.Headers != nil {
		for key, value := range req.Headers {
			request.Header.Add(key, value)
		}
	}
	setContentType(request, req.BodyData)
	logger.Log().Debug("正在发送请求Request URL: ", urlWithQuery)
	resp, err := client.Do(request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "Client.Timeout exceeded") {
			return nil, fmt.Errorf("request timeout: %v", err)
		}
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP request failed with status code: %d  for URL: %s", resp.StatusCode, urlWithQuery)
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, nil
	}
	// 打印响应体内容
	bodyString := string(bodyBytes)
	logger.Log().Debugf("%s - Response Body: %s", urlWithQuery, bodyString)
	var result map[string]interface{}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}
	return result, nil
}
