package hihttp

import (
	"context"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	// API接口前缀
	baseUrl string
	method  string
	header  *http.Header
	cookies []*http.Cookie
	client  hiclient
	opt     Options
}

// 公共参
func New(opts ...Option) *Request {
	request := Request{
		baseUrl: "",
		method:  "",
		header:  &http.Header{},
		cookies: []*http.Cookie{},
		client:  client,
	}
	opt := client.opt

	for _, o := range opts {
		o(&opt)
	}

	request.opt = opt
	return &request
}

type HiHTTP interface {
	Get(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
	Post(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
	Put(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
	Delete(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
	Patch(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
}

func (r *Request) execute(ctx context.Context, payload io.Reader) ([]byte, error) {
	if len(r.baseUrl) > 0 {
		r.baseUrl = strings.Trim(r.baseUrl, defaultTrimChars)
	}

	// var payload io.Reader
	httpCtx, cancel := context.WithTimeout(ctx, r.opt.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(httpCtx, r.method, r.baseUrl, payload)
	if err != nil {
		return nil, err
	}
	if r.header != nil {
		req.Header = *r.header
	}
	if len(r.cookies) > 0 {
		for _, cookie := range r.cookies {
			req.AddCookie(cookie)
		}
	}

	res, err := r.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r.client.statusCode = res.StatusCode
	return body, nil
}

// send get request
// 也可以把参数直接放到URL后面，则data传nil即可
func (r *Request) Get(ctx context.Context, urlStr string, data ...Param) ([]byte, error) {
	r.baseUrl = urlStr + MergeParams(data...)
	r.method = GET
	req, err := r.execute(ctx, nil)
	if err != nil {
		r.retry(ctx, nil)
	}
	return req, err
}

// send post request
func (r *Request) Post(ctx context.Context, urlStr string, p Payload) ([]byte, error) {
	r.baseUrl = urlStr
	r.method = POST
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}

func (r *Request) Put(ctx context.Context, urlStr string, p Payload) ([]byte, error) {
	r.baseUrl = urlStr
	r.method = PUT
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}

func (r *Request) Delete(ctx context.Context, urlStr string, data ...Param) ([]byte, error) {
	r.baseUrl = urlStr + MergeParams(data...)
	r.method = DELETE
	req, err := r.execute(ctx, nil)
	if err != nil {
		r.retry(ctx, nil)
	}
	return req, err
}

func (r *Request) Patch(ctx context.Context, urlStr string, p Payload) ([]byte, error) {
	r.baseUrl = urlStr
	r.method = PATCH
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}
