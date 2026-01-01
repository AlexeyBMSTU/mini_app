package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/middleware"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

type Client struct {
	client      *http.Client
	baseURL     string
	timeout     time.Duration
	logger      logger.Logger
	authToken   string
	authHeaders map[string]string
}

type Option func(*Client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithAuthToken(token string) Option {
	return func(c *Client) {
		c.authToken = token
	}
}

func WithAuthHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.authHeaders = headers
	}
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		client:    &http.Client{},
		timeout:   30 * time.Second,
		logger:    logger.GetLogger(),
		authToken: "",
		authHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	client.client.Timeout = client.timeout

	return client
}

func (c *Client) DoRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*Response, error) {
	reqURL := url
	if c.baseURL != "" && !isAbsoluteURL(url) {
		reqURL = c.baseURL + url
	}

	requestLogger := c.logger
	if requestID := ctx.Value(middleware.RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			requestLogger = c.logger.WithRequestID(id)
		}
	}

	requestLogger.Debugf("Making %s request to %s", method, reqURL)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		requestLogger.Errorf("Error creating request: %v", err)
		return nil, err
	}

	for key, value := range c.authHeaders {
		req.Header.Set(key, value)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", c.authToken)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		requestLogger.Errorf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		requestLogger.Errorf("Error reading response body: %v", err)
		return nil, err
	}

	requestLogger.Debugf("Response status: %d, body: %s", resp.StatusCode, string(respBody))

	response := &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.DoRequest(ctx, "GET", url, nil, headers)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*Response, error) {
	return c.DoRequest(ctx, "POST", url, body, headers)
}

func (c *Client) Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*Response, error) {
	return c.DoRequest(ctx, "PUT", url, body, headers)
}

func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.DoRequest(ctx, "DELETE", url, nil, headers)
}

func (c *Client) GetJSON(ctx context.Context, url string, target interface{}, headers map[string]string) error {
	resp, err := c.Get(ctx, url, headers)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp)
	}

	return json.Unmarshal(resp.Body, target)
}

func (c *Client) PostJSON(ctx context.Context, url string, body interface{}, target interface{}, headers map[string]string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.Post(ctx, url, bytes.NewReader(jsonBody), headers)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp)
	}

	if target != nil {
		return json.Unmarshal(resp.Body, target)
	}

	return nil
}

func (c *Client) handleErrorResponse(resp *Response) error {
	var errorResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(resp.Body, &errorResp); err == nil {
		if errorResp.Message != "" {
			return &HTTPError{
				StatusCode: resp.StatusCode,
				Message:    errorResp.Message,
			}
		}
		if errorResp.Error != "" {
			return &HTTPError{
				StatusCode: resp.StatusCode,
				Message:    errorResp.Error,
			}
		}
	}

	return &HTTPError{
		StatusCode: resp.StatusCode,
		Message:    string(resp.Body),
	}
}

func isAbsoluteURL(url string) bool {
	return len(url) > 7 && url[:7] == "http://" || 
		   len(url) > 8 && url[:8] == "https://"
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}