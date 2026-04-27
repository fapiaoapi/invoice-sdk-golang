// Use of this source code is governed by a Apache-2.0
// license that can be found in the LICENSE file.

package invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

const (
	DefaultBaseURL = "https://api.fa-piao.com"
)

// Client 数电发票API客户端
type Client struct {
	BaseURL    string
	AppKey     string
	AppSecret  string
	Token      string
	Debug      bool
	HTTPClient *http.Client
}

// NewClient 创建新的客户端实例
func NewClient(appKey, appSecret string, debug bool) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		AppKey:     appKey,
		AppSecret:  appSecret,
		Debug:      debug,
		HTTPClient: &http.Client{Timeout: 150 * time.Second},
	}
}

// SetBaseURL 设置API基础URL
func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

// SetToken 设置授权令牌
func (c *Client) SetToken(token string) {
	c.Token = token
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(method, path string, params map[string]string) (*Response, error) {
	// 生成签名参数
	randomString := generateRandomString(20)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// 计算签名
	signature := calculateSignature(method, path, randomString, timestamp, c.AppKey, c.AppSecret)

	// 创建请求体
	payload, contentType := createRequestBody(params)

	// 创建HTTP请求
	req, err := http.NewRequest(method, c.BaseURL+path, payload)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("AppKey", c.AppKey)
	req.Header.Set("Sign", signature)
	req.Header.Set("TimeStamp", timestamp)
	req.Header.Set("RandomString", randomString)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Sdk", "Go1016")

	// 如果有授权令牌，添加到请求头
	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}

	if c.Debug {
		c.printDebugRequest(method, c.BaseURL+path, req.Header, params)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	response, responseBody, err := handleResponse(resp)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		c.printDebugResponse(resp.StatusCode, responseBody)
	}

	return response, nil
}

func (c *Client) doRequestWithFields(method, path string, fields []formField) (*Response, error) {
	randomString := generateRandomString(20)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	signature := calculateSignature(method, path, randomString, timestamp, c.AppKey, c.AppSecret)

	payload, contentType := createMultipartBodyWithFields(fields)

	req, err := http.NewRequest(method, c.BaseURL+path, payload)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("AppKey", c.AppKey)
	req.Header.Set("Sign", signature)
	req.Header.Set("TimeStamp", timestamp)
	req.Header.Set("RandomString", randomString)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Sdk", "Go1016")

	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}

	if c.Debug {
		c.printDebugRequest(method, c.BaseURL+path, req.Header, fields)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	response, responseBody, err := handleResponse(resp)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		c.printDebugResponse(resp.StatusCode, responseBody)
	}

	return response, nil
}

func (c *Client) printDebugRequest(method, url string, header http.Header, params interface{}) {
	fmt.Printf("[invoice-sdk debug] request method=%s url=%s\n", method, url)
	fmt.Printf("[invoice-sdk debug] request header=%s\n", formatHeader(header))
	fmt.Printf("[invoice-sdk debug] request params=%v\n", params)
}

func (c *Client) printDebugResponse(statusCode int, body []byte) {
	fmt.Printf("[invoice-sdk debug] response status=%d\n", statusCode)
	fmt.Printf("[invoice-sdk debug] response body=%s\n", bytes.TrimSpace(body))
}

func formatHeader(header http.Header) string {
	keys := make([]string, 0, len(header))
	for key := range header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make(map[string][]string, len(header))
	for _, key := range keys {
		result[key] = header[key]
	}

	headerJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("%v", header)
	}

	return string(headerJSON)
}

// 解析响应数据到指定结构
func parseResponseData(response *Response, target interface{}) error {
	if !response.IsSuccess() {
		return fmt.Errorf(response.Error())
	}

	if err := json.Unmarshal(response.Data, target); err != nil {
		return fmt.Errorf("解析数据失败: %v", err)
	}

	return nil
}

// ParseResponseData 解析响应数据到指定结构
func ParseResponseData(response *Response, target interface{}) error {
	if !response.IsSuccess() {
		return fmt.Errorf(response.Error())
	}

	if err := json.Unmarshal(response.Data, target); err != nil {
		return fmt.Errorf("解析数据失败: %v", err)
	}

	return nil
}
