# 介绍

该库对go愿生http库的封装，添加了错误处理，超时处理及重试功能。

# Todo

- [x] 错误处理；
- [x] 超时处理；
- [x] 重试功能；
- [x] 错误回调；
- [x] GET、POST支持；
- [x] DELETE、PUT、PATCH支持；



# 全局参数设置

```go
hihttp.Load(
		hihttp.WithTimeout(time.Second), // 设置全局超时时间
		hihttp.WithRetryCount(1), // 设置全局重试次数
		hihttp.WithRetryWait(time.Second),// 设置全局重试等待时间
  	// 设置全局重试错误时的回调方法
		hihttp.WithRetryError(func(ctx context.Context, c hiclient) error {
			return nil
		}),
	)
```



# GET示例

```go
// 正常发送一个Get请求
res, err := hihttp.New().Get(context.Background(), "http://www.google.com")
if err != nil {
	t.Error(err)
}

// 添加header和cookie
hihttp.New().
	SetHeader("token", "1234567890").
	SetCookies(&http.Cookie{
			Name:  "token",
			Value: "abcdefg",
	}).
	Get(context.Background(), "http://www.google.com")

// 添加cookie
```





# POST 示例

```go
// -- multipart/form-data -- // 
// 以map的形式添加post参数
hihttp.New().Post(context.Background(), "http://www.yumontime.com/test/login", map[string]interface{}{
		"user_name": "yumontime", "password": "123123",
})
// 以[]string的形式添加post参数，第一个参数对应的是key，第二个参数对应的是value，以此类推
hihttp.New().Post(context.Background(), "http://www.yumontime.com/test/login", "user_name", "yumontime", "password", "123123")

// -- application/json -- //
hihttp.New().SetHeader(hihttp.SerializationType,hihttp.SerializationTypeJSON).Post(context.Background(), "http://www.yumontime.com/test/login", "user_name=yumontime&password=123123")

```



