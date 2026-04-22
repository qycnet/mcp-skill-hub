package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client API 客户端
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient 创建 API 客户端
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchSkills 搜索技能
func (c *Client) SearchSkills(query string, page, pageSize int) (*SearchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/search?q=%s&page=%d&page_size=%d",
		c.baseURL, query, page, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSkill 获取技能详情
func (c *Client) GetSkill(id string) (*SkillDetail, error) {
	url := fmt.Sprintf("%s/api/v1/skills/%s", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result SkillDetail
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// InstallSkill 安装技能
func (c *Client) InstallSkill(id string) (*InstallResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills/%s/download", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result InstallResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PublishSkill 发布技能
func (c *Client) PublishSkill(skill *PublishRequest) (*PublishResponse, error) {
	url := fmt.Sprintf("%s/api/v1/skills", c.baseURL)

	data, err := json.Marshal(skill)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result PublishResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Login 用户登录
func (c *Client) Login(username, password string) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/login", c.baseURL)

	data, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result LoginResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API 错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// 响应类型

type SearchResponse struct {
	Skills   []SkillSummary `json:"skills"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type SkillSummary struct {
	ID           uint    `json:"id"`
	UUID         string  `json:"uuid"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Downloads    int64   `json:"downloads"`
	Rating       float64 `json:"rating"`
	QualityScore float64 `json:"quality_score"`
}

type SkillDetail struct {
	SkillSummary
	Author      AuthorInfo     `json:"author"`
	Tags        []TagInfo      `json:"tags"`
	Versions    []VersionInfo  `json:"versions"`
	DownloadURL string         `json:"download_url"`
}

type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type TagInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type VersionInfo struct {
	ID        uint   `json:"id"`
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
}

type InstallResponse struct {
	DownloadURL string `json:"download_url"`
	Message     string `json:"message"`
}

type PublishRequest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	License     string   `json:"license"`
	Repository  string   `json:"repository"`
	Homepage    string   `json:"homepage"`
}

type PublishResponse struct {
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Message  string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
