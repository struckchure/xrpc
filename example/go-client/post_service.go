package main

import (
	"encoding/json"
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"net/url"
)

type PostServiceClient struct {
	client *resty.Client
}

func structToQueryParams(input any) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}
	var mapData map[string]any
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	query := url.Values{}
	for key, value := range mapData {
		if value == nil {
			continue
		}
		strVal := fmt.Sprintf("%v", value)
		query.Add(key, strVal)
	}
	return query.Encode(), nil
}

type MapError map[string]any

func (m MapError) Error() string {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("failed to marshal error map: %v", err)
	}
	return string(data)
}

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListPostInput struct {
	Skip  *int `json:"skip"`
	Limit *int `json:"limit"`
}

func (c *PostServiceClient) PostList(input ListPostInput) (*[]Post, error) {
	queryParams, err := structToQueryParams(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.R().SetQueryString(queryParams).SetError(&MapError{}).SetResult(&[]Post{}).Get("/post/list/")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, resp.Error().(*MapError)
	}
	return resp.Result().(*[]Post), nil
}

type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *PostServiceClient) PostCreate(input CreatePostInput) (*Post, error) {
	resp, err := c.client.R().SetBody(input).SetError(&MapError{}).SetResult(&Post{}).Post("/post/create/")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, resp.Error().(*MapError)
	}
	return resp.Result().(*Post), nil
}

type GetPostInput struct {
	Id       int    `json:"id"`
	AuthorId string `json:"author_id"`
}

func (c *PostServiceClient) PostGet(input GetPostInput) (*Post, error) {
	queryParams, err := structToQueryParams(input)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.R().SetQueryString(queryParams).SetError(&MapError{}).SetResult(&Post{}).Get("/post/get/")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, resp.Error().(*MapError)
	}
	return resp.Result().(*Post), nil
}

func NewPostServiceClient() *PostServiceClient {
	client := resty.New()
	client.SetBaseURL("http://localhost:9090")

	return &PostServiceClient{client: client}
}
