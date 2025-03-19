package fuel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
	logger     Logger
}

type Options func(*Client)

func WithHTTPClient(cli *http.Client) Options {
	return func(c *Client) {
		c.httpClient = cli
	}
}

func WithLogger(logger Logger) Options {
	return func(c *Client) {
		c.logger = logger
	}
}

func NewClient(endpoint string, opts ...Options) *Client {
	c := &Client{endpoint: endpoint}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	return c
}

func NewClientWithLogger(endpoint string, logger Logger) *Client {
	return NewClient(endpoint, WithLogger(logger))
}

type QueryErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (loc QueryErrorLocation) String() string {
	return fmt.Sprintf("(line:%d,column:%d)", loc.Line, loc.Column)
}

type QueryError struct {
	Message   string               `json:"message"`
	Locations []QueryErrorLocation `json:"locations"`
}

func (e QueryError) String() string {
	var buf bytes.Buffer
	for i, loc := range e.Locations {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(loc.String())
	}
	buf.WriteString(": ")
	buf.WriteString(e.Message)
	return buf.String()
}

type QueryErrors []QueryError

func (e QueryErrors) Error() string {
	var buf bytes.Buffer
	for i, ei := range e {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.WriteString(ei.String())
	}
	return fmt.Sprintf("execute query failed: %s", buf.String())
}

func ExecuteQuery[DATA any](ctx context.Context, cli *Client, query string) (data DATA, err error) {
	if cli.logger != nil {
		cli.logger.Infof("execute query: %s", query)
	}
	start := time.Now()
	var reqBody bytes.Buffer
	if err = json.NewEncoder(&reqBody).Encode(map[string]any{"query": query}); err != nil {
		return data, fmt.Errorf("build request failed: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", cli.endpoint, &reqBody)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return data, fmt.Errorf("build request failed: %w", err)
	}

	var resp *http.Response
	resp, err = cli.httpClient.Do(req)
	if err != nil {
		return data, fmt.Errorf("send request failed: %w", err)
	}

	defer resp.Body.Close()
	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("read response body failed: %w", err)
	}
	if cli.logger != nil {
		cli.logger.Infof("query result(len: %d, used: %s): %s", len(respBody), time.Since(start), string(respBody))
	}

	var result struct {
		Data   DATA        `json:"data"`
		Errors QueryErrors `json:"errors"`
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return data, fmt.Errorf("parse response body failed: %w", err)
	}
	if len(result.Errors) > 0 {
		return data, result.Errors
	}
	return result.Data, nil
}
