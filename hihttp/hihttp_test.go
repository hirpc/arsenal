package hihttp

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// 设置全局超时时间
	Load(
		WithTimeout(time.Second),
		WithRetryCount(1),
		WithRetryWait(time.Second),
		WithRetryError(func(ctx context.Context, r *Request) error {
			return nil
		}),
	)
	res, err := New().Get(context.Background(), "http://www.google.com", NewQueryParam(""))
	if err != nil {
		t.Error(1, err)
	}

	// 设置当前次请求的超时时间
	res, err = New(WithTimeout(3*time.Second)).Get(context.Background(), "http://www.google.com")
	if err != nil {
		t.Error(2, err)
	}

	// 不设置当前次请求的超时时间，则使用默认的全局超时时间，
	// 也就是SetClient(WithTimeout(time.Second))设置的
	res, err = New().
		SetHeader("token", "1234567890").
		SetCookies(&http.Cookie{
			Name:  "token",
			Value: "abcdefg",
		}).
		Get(context.Background(), "http://www.google.com")
	if err != nil {
		t.Error(3, err)
	}
	// end
	t.Log(4, string(res))
}

func TestPost(t *testing.T) {
	// 发送一个post请求
	res, err := New().Post(context.Background(), "http://www.yumontime.com/test/login", NewWWWFormPayload(map[string]interface{}{
		"username": "jankin",
	}))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res))
	//
	// New().Post(context.Background(), "http://www.yumontime.com/test/login", NewFormPayload(map[string]interface{}{
	// 	"user_name": "yumontime", "password": "123123",
	// }))

	// New().SetHeader(SerializationType, SerializationTypeJSON)
}

func TestPut(t *testing.T) {
	r := New().SetCookies()
	res, err := r.Put(context.Background(), "http://127.0.0.1:8080/test_http_method", NewFormPayload(map[string]interface{}{
		"name": "jankin",
	}))
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestDelete(t *testing.T) {
	res, err := New().Delete(context.Background(), "http://127.0.0.1:8080/test_http_method", NewMapParams(map[string]interface{}{
		"name": "jankin",
	}))
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestPatch(t *testing.T) {
	res, err := New().Patch(context.Background(), "http://127.0.0.1:8080/test_http_method", NewFormPayload(map[string]interface{}{
		"name": "jankin",
	}))
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestGET(t *testing.T) {
	res, err := New().Get(context.Background(), "http://127.0.0.1:8080/test_http_method", NewKVParam("name", "jankin"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res))
}

// TestTimeout 请求超时
func TestTimeout(t *testing.T) {
	urlStr := "https://www.google.com"

	// 1. 使用WithTimeout()方法来实现超时控制
	ctx := context.Background()
	res1, err := New(WithRetryCount(2), WithTimeout(3*time.Second)).Get(ctx, urlStr)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res1))

	// 2. 将传入的ctx本身设置上超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res2, err := New().Get(ctx, urlStr)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res2))
}

// TestRetry 重试
func TestRetry(t *testing.T) {
	urlStr := "https://www.google.com"

	// 直接使用
	// 请求失败后会再次进行重试请求
	ctx := context.Background()
	st := time.Now().Unix()
	res1, err := New(WithRetryCount(1), WithTimeout(3*time.Second)).Get(ctx, urlStr)
	if err != nil {
		t.Error(err)
	}
	en := time.Now().Unix()
	t.Log(string(res1), en-st)
	// 如果无需重试，3秒后超时，则en-st = 3
	// 如果重试1次，3秒后超时，重试又3秒，因此en-st=6
}
