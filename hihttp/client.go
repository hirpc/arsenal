package hihttp

import (
	"context"
	"io"
	"net/http"
	"time"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"
	// MethodPost HTTP method
	MethodPost = "POST"
	// MethodPut HTTP method
	MethodPut = "PUT"
	// MethodDelete HTTP method
	MethodDelete = "DELETE"
)
const (
	SerializationType          string = "Content-Type"
	SerializationTypeFormData  string = "multipart/form-data"
	SerializationTypeJSON      string = "application/json"
	SerializationTypeWWWFrom   string = "application/x-www-form-urlencoded"
	SerializationTypePlainText string = "text/plain; charset=utf-8"
)

type hiclient struct {
	client  *http.Client
	header  *http.Header
	cookies []*http.Cookie
	// API接口前缀
	prefix string
	route  string
	// baseUrl = prefix + route
	baseUrl string
	method  string
	payload io.Reader
	// 重试次数
	RetryCount int
	// 重试等待时间
	RetryWaitTime time.Duration
	RetryError    RetryErrorFunc
	// 超时时间
	Timeout time.Duration
}
type RetryErrorFunc func(ctx context.Context, c hiclient) error

var defaultTrimChars = string([]byte{
	'\t', // Tab.
	'\v', // Vertical tab.
	'\n', // New line (line feed).
	'\r', // Carriage return.
	'\f', // New page.
	' ',  // Ordinary space.
	0x00, // NUL-byte.
	0x85, // Delete.
	0xA0, // Non-breaking space.
})
var client = &hiclient{}

func Client() *hiclient {
	return client
}

// Setting the header parameter should be 'map[string]string{}' or '[]string{}'
// Usually you need to set up 'Content-Type'
// Example:
// c.Headers([]string{"key","value"})
// c.Headers(map[string]string{"key":"value"})
func (c hiclient) SetHeader(args interface{}) hiclient {
	if argsMap, ok := args.(map[string]string); ok {
		for k, v := range argsMap {
			c.header.Add(k, v)
		}
		return c
	}

	if argsSlice, ok := args.([]string); ok {
		if len(argsSlice)%2 != 0 {
			return c
		}
		for i := 1; i < len(argsSlice); i += 2 {
			c.header.Add(argsSlice[i-1], argsSlice[i])
		}
	}

	return c
}

// SetCookies 设置cookie
func (c hiclient) SetCookies(hc ...*http.Cookie) hiclient {
	c.cookies = append(c.cookies, hc...)
	return c
}
