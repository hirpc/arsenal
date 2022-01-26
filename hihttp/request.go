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

type HiHTTP interface {
	Get(ctx context.Context, urlStr string) ([]byte, error)
	Post(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error)
	execute(ctx context.Context, payload io.Reader) ([]byte, error)
}

func (c hiclient) execute(ctx context.Context, payload io.Reader) ([]byte, error) {
	if len(c.baseUrl) > 0 {
		c.baseUrl = strings.Trim(c.baseUrl, defaultTrimChars)
	}

	// var payload io.Reader
	httpCtx, cancel := context.WithTimeout(ctx, c.Timeout)
	req, err := http.NewRequestWithContext(httpCtx, c.method, c.baseUrl, payload)
	defer cancel()
	if err != nil {
		return nil, err
	}

	if c.header != nil {
		req.Header = *c.header
	}
	if len(c.cookies) > 0 {
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// send get request
func (c hiclient) Get(ctx context.Context, urlStr string) ([]byte, error) {
	c.baseUrl = urlStr
	c.method = GET
	req, err := c.execute(ctx, nil)
	if err != nil {
		c.retry(ctx, c, nil)
	}
	return req, err
}

// send post request
func (c hiclient) Post(ctx context.Context, urlStr string, data ...interface{}) ([]byte, error) {
	var payload io.Reader
	if len(data) > 0 {
		switch c.header.Get(SerializationType) {
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
			if c.header.Get(SerializationType) == "" {
				c.header.Set(SerializationType, SerializationTypeWWWFrom)
			}
			params := []string{}
			if len(data) > 1 && len(data)%2 == 0 {
				for i := 1; i < len(data); i = i + 2 {
					params = append(params, fmt.Sprintf("%v=%v", data[i-1], data[i]))
				}
			}
			payload = strings.NewReader(strings.Join(params, "&"))
		default:
			payloadBuf := &bytes.Buffer{}
			writer := multipart.NewWriter(payloadBuf)
			// Set the default 'Content-Type'
			if c.header.Get(SerializationType) == "" {
				c.header.Set(SerializationType, writer.FormDataContentType())
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
	c.baseUrl = urlStr
	c.method = POST
	req, err := c.execute(ctx, payload)
	if err != nil {
		c.retry(ctx, c, payload)
	}

	return req, err
}
