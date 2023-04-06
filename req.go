package aime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type RequestData struct {
	Model       string     `json:"model"`
	Messages    []*Message `json:"messages"`
	Temperature float64    `json:"temperature"`
}

type ResponseData struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type authRT struct {
	key string
}

func (rt *authRT) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.key))
	return http.DefaultTransport.RoundTrip(req)
}

var _ http.RoundTripper = new(authRT)

func newAuthRTFromKey(key string) http.RoundTripper {
	return &authRT{key: key}
}

func newAuthRTFromKeyFile(keyfile string) (http.RoundTripper, error) {
	if b, err := ioutil.ReadFile(keyfile); err != nil {
		return nil, err
	} else {
		return newAuthRTFromKey(string(b)), nil
	}
}

func ToRequest(data *RequestData) (*http.Request, error) {
	const url = "https://api.openai.com/v1/chat/completions"
	if body, err := json.Marshal(data); err != nil {
		return nil, err
	} else if req, err := http.NewRequest("POST", url, bytes.NewBuffer(body)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

type ClientOption func(c *http.Client) error

func WithKey(key string) ClientOption {
	return func(c *http.Client) error {
		c.Transport = newAuthRTFromKey(key)
		return nil
	}
}

func WithKeyFile(keyfile string) ClientOption {
	rt, err := newAuthRTFromKeyFile(keyfile)
	return func(c *http.Client) error {
		c.Transport = rt
		return err
	}
}

func ToClient(opts ...ClientOption) (*http.Client, error) {
	var client = &http.Client{}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func toResponseData(rc io.ReadCloser) (*ResponseData, error) {
	defer func() { _ = rc.Close() }()
	var data = new(ResponseData)
	if b, err := ioutil.ReadAll(rc); err != nil {
		return nil, err
	} else if err := json.Unmarshal(b, data); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func MakeRequest(client *http.Client, req *http.Request) (*ResponseData, error) {
	if resp, err := client.Do(req); err != nil {
		return nil, err
	} else if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("missing response body")
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	} else {
		return toResponseData(resp.Body)
	}
}
