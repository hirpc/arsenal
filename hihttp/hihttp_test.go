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
		WithRetryError(func(ctx context.Context, c hiclient) error {
			return nil
		}),
	)
	res, err := New().Get(context.Background(), "http://www.google.com")
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
	New().Post(context.Background(), "http://www.yumontime.com/test/login", "user_name", "yumontime", "password", "123123")

	//
	New().Post(context.Background(), "http://www.yumontime.com/test/login", map[string]interface{}{
		"user_name": "yumontime", "password": "123123",
	})

	New().SetHeader(SerializationType, SerializationTypeJSON)
}

func TestPut(t *testing.T) {
	r := New().SetCookies()
	res, err := r.Put(context.Background(), "http://127.0.0.1:8080/test_http_method", map[string]interface{}{
		"name": "jankin",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestDelete(t *testing.T) {
	res, err := New().Delete(context.Background(), "http://127.0.0.1:8080/test_http_method", map[string]interface{}{
		"name": "jankin",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestPatch(t *testing.T) {
	res, err := New().Patch(context.Background(), "http://127.0.0.1:8080/test_http_method", map[string]interface{}{
		"name": "jankin",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestGET(t *testing.T) {
	res, err := New().SetHeader(SerializationType, SerializationTypeWWWFrom).Get(context.Background(), "http://127.0.0.1:8080/test_http_method", "name", "jankin")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}
