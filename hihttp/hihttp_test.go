package hihttp

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	New(Config{
		Prefix:        "",
		RetryCount:    3,
		RetryWaitTime: time.Second,
		RetryError: func(ctx context.Context, c hiclient) error {
			fmt.Println(ctx, c.baseUrl, c.retryCount)
			return nil
		},
	})

	res, err := Client().Get(context.Background(), "https://www.google.com/")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestPost(t *testing.T) {
	New(Config{
		Prefix:        "",
		RetryCount:    3,
		RetryWaitTime: time.Second,
		RetryError: func(ctx context.Context, c hiclient) error {
			fmt.Println(ctx, c.baseUrl, c.retryCount)
			return nil
		},
	})

	res, err := Client().SetHeader(SerializationType, SerializationTypeFormData).Post(context.Background(), "https://www.google.com")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}
