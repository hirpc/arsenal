package hihttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

func (r *Request) formatPayload(ctx context.Context, data ...interface{}) (io.Reader, error) {
	var payload io.Reader
	if len(data) > 0 {
		switch r.header.Get(SerializationType) {
		case SerializationTypeJSON:
			var params string
			switch data[0].(type) {
			case string, []byte:
				params = fmt.Sprint(data[0])
			default:
				if b, err := json.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
			payload = strings.NewReader(params)
		case SerializationTypeWWWFrom:
			params := []string{}
			if dataMap, ok := data[0].(map[string]interface{}); ok {
				for k, v := range dataMap {
					params = append(params, fmt.Sprintf("%s=%v", k, v))
				}
			}else if dataMap, ok := data[0].(map[string]string); ok{
				for k, v := range dataMap {
					params = append(params, fmt.Sprintf("%s=%s", k, v))
				}
			}else if len(data) > 1 && len(data)%2 == 0 {
				for i := 1; i < len(data); i = i + 2 {
					params = append(params, fmt.Sprintf("%v=%v", data[i-1], data[i]))
				}
			}
			payload = strings.NewReader(strings.Join(params, "&"))
		default:
			payloadBuf := &bytes.Buffer{}
			writer := multipart.NewWriter(payloadBuf)
			// Set the default 'Content-Type'
			if r.header.Get(SerializationType) == "" {
				r.header.Set(SerializationType, writer.FormDataContentType())
			}

			if dataMap, ok := data[0].(map[string]interface{}); ok {
				for k, v := range dataMap {
					_ = writer.WriteField(k, fmt.Sprint(v))
				}
				if err := writer.Close(); err != nil {
					return nil, err
				}
			}
			payload = payloadBuf
		}
	}
	return payload, nil
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
func (r *Request) Get(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	if r.header.Get(SerializationType) == "" {
		// set default SerializationType.
		r.header.Set(SerializationType, SerializationTypeWWWFrom)
	}
	payload, err := r.formatPayload(ctx, data...)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(payload)
	if err != nil {
		return nil, err
	}
	r.baseUrl = urlStr + "?" + string(buf)
	r.method = GET
	req, err := r.execute(ctx, nil)
	if err != nil {
		r.retry(ctx, nil)
	}
	return req, err
}

// send post request
func (r *Request) Post(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	payload, err := r.formatPayload(ctx, data...)
	if err != nil {
		return nil, err
	}
	r.baseUrl = urlStr
	r.method = POST
	req, err := r.execute(ctx, payload)
	if err != nil {
		r.retry(ctx, payload)
	}
	return req, err
}

func (r *Request) Put(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	payload, err := r.formatPayload(ctx, data...)
	if err != nil {
		return nil, err
	}
	r.baseUrl = urlStr
	r.method = PUT
	req, err := r.execute(ctx, payload)
	if err != nil {
		r.retry(ctx, payload)
	}
	return req, err
}

func (r *Request) Delete(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	payload, err := r.formatPayload(ctx, data...)
	if err != nil {
		return nil, err
	}
	r.baseUrl = urlStr
	r.method = DELETE
	req, err := r.execute(ctx, payload)
	if err != nil {
		r.retry(ctx, payload)
	}
	return req, err
}

func (r *Request) Patch(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	payload, err := r.formatPayload(ctx, data...)
	if err != nil {
		return nil, err
	}
	r.baseUrl = urlStr
	r.method = PATCH
	req, err := r.execute(ctx, payload)
	if err != nil {
		r.retry(ctx, payload)
	}
	return req, err
}
