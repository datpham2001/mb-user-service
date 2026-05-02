package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	base       string
	httpClient *http.Client
	headers    map[string]string
}

type Option func(*Client)

func WithBaseURL(base string) Option {
	return func(c *Client) { c.base = strings.TrimRight(base, "/") }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

func WithHeader(key, value string) Option {
	return func(c *Client) { c.headers[key] = value }
}

func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		headers:    make(map[string]string),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type RequestOption func(*http.Request)

func WithBearerToken(token string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+token)
	}
}

func WithRequestHeader(key, value string) RequestOption {
	return func(r *http.Request) { r.Header.Set(key, value) }
}

func (c *Client) Get(ctx context.Context, path string, dst any, opts ...RequestOption) error {
	return c.do(ctx, http.MethodGet, path, nil, "", dst, opts...)
}

func (c *Client) PostForm(ctx context.Context, path string, form url.Values, dst any, opts ...RequestOption) error {
	return c.do(ctx, http.MethodPost, path,
		strings.NewReader(form.Encode()),
		"application/x-www-form-urlencoded",
		dst, opts...)
}

func (c *Client) PostJSON(ctx context.Context, path string, body any, dst any, opts ...RequestOption) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("httpclient: marshal request body: %w", err)
	}
	return c.do(
		ctx, http.MethodPost, path,
		bytes.NewReader(data),
		"application/json",
		dst, opts...)
}

func (c *Client) do(
	ctx context.Context,
	method, path string,
	body io.Reader,
	contentType string,
	dst any,
	opts ...RequestOption,
) error {
	target := c.base + path

	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return fmt.Errorf("httpclient: build request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for _, o := range opts {
		o(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpclient: %s %s: %w", method, target, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("httpclient: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{
			Method:     method,
			URL:        target,
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	if dst != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("httpclient: decode response: %w", err)
		}
	}

	return nil
}

type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("httpclient: %s %s → HTTP %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}

func IsStatus(err error, code int) bool {
	var e *HTTPError
	if err == nil {
		return false
	}
	e, ok := err.(*HTTPError)
	return ok && e.StatusCode == code
}
