package hihttp

import (
	"context"
	"net/http"
	"time"
)

type hiclient struct {
	client     *http.Client
	opt        Options
	statusCode int
}
type RetryErrorFunc func(ctx context.Context, r *Request) error

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
var client = hiclient{
	client: &http.Client{
		Timeout: 5 * time.Second,
	},
	opt: Options{
		retryCount: 0,
		retryWait:  time.Duration(0),
		retryError: func(ctx context.Context, r *Request) error {
			return nil
		},
		timeout: 5 * time.Second,
	},
}

// Load 设置client的全局参数
func Load(opts ...Option) {
	for _, o := range opts {
		o(&client.opt)
	}
}

// SetHeader 以k-v格式设置header
func (r *Request) SetHeader(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

// SetHeaders 设置header参数 map[string]string{}
// 例如: c.Headers(map[string]string{"key":"value"})
func (r *Request) SetHeaders(args map[string]string) *Request {
	for k, v := range args {
		r.header.Add(k, v)
	}
	return r
}

// SetCookies 设置cookie
func (r *Request) SetCookies(hc ...*http.Cookie) *Request {
	r.cookies = append(r.cookies, hc...)
	return r
}
