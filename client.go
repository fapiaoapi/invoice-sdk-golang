package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	HTTPClient *http.Client
}

// NewClient 创建新的客户端实例
func NewClient(appKey, appSecret string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		AppKey:     appKey,
		AppSecret:  appSecret,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
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

	// 如果有授权令牌，添加到请求头
	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	return handleResponse(resp)
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