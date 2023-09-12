package esapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
	Client  *http.Client
	Auth    Authorization
	Headers map[string]string
}

func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		Client:  http.DefaultClient,
	}
}

func (c *Client) SkipTLS() {
	c.Client.Transport = &http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func (c *Client) SetAuth(auth Authorization) {
	c.Auth = auth
}

func buildRequest(method, url string, body []byte) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, url, nil)
	}
	buf := bytes.NewBuffer(body)
	return http.NewRequest(method, url, buf)
}

func (c *Client) Do(method, endpoint string, body []byte) ([]byte, error) {
	url := renderEndpoint(c.BaseURL, endpoint)
	method = strings.ToUpper(method)
	req, err := buildRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("buildRequest() = %s", err)
	}
	if c.Auth != nil {
		req.Header.Set("Authorization", c.Auth.AuthorizationHeader())
	}
	if c.Headers != nil {
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func renderEndpoint(baseURL, endpoint string) string {
	base := strings.TrimRight(baseURL, "/")
	path := strings.TrimLeft(endpoint, "/")
	return base + "/" + path
}
