package hihttp

import (
	"context"
	"net/http"
	"time"
)

type hiclient struct {
	client  *http.Client
	header  *http.Header
	cookies []*http.Cookie
	// API接口前缀
	baseUrl string
	method  string
	// 重试次数
	retryCount int
	// 重试等待时间
	retryWaitTime time.Duration
	retryError    RetryErrorFunc
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

func (c hiclient) SetHeader(key, value string) hiclient {
	c.header.Add(key, value)
	return c
}

// Setting the header parameter should be 'map[string]string{}'
// Usually you need to set up 'Content-Type'
// Example:
// c.Headers(map[string]string{"key":"value"})
func (c hiclient) SetHeaders(args map[string]string) hiclient {
	for k, v := range args {
		c.header.Add(k, v)
	}
	return c
}

// SetCookies 设置cookie
func (c hiclient) SetCookies(hc ...*http.Cookie) hiclient {
	c.cookies = append(c.cookies, hc...)
	return c
}
