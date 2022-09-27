# 介绍

该库对Go原生http库做了封装，添加了错误处理，超时处理及重试功能。

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
	hihttp.WithRetryError(func(ctx context.Context, r hihttp.Request) error {
		return nil
	}),
)
```



# GET示例

```go
// 正常发送一个Get请求
res, err := hihttp.New().Get(context.Background(), "http://www.google.com")
if err != nil {
	fmt.Println(err)
}

// 添加header和cookie
hihttp.New().SetHeader("token", "1234567890").SetCookies(&http.Cookie{
			Name:  "token",
			Value: "abcdefg",
	}).Get(context.Background(), "http://www.google.com")


// 正常发送一个Get请求,追加get参数，以key-value格式
res, err := hihttp.New().Get(context.Background(), "http://www.google.com",hihttp.NewKVParam("name", "jankin"))
if err != nil {
	fmt.Println(err)
}

// 正常发送一个Get请求,追加get参数，以map格式
res, err := hihttp.New().Get(context.Background(), "http://www.google.com",hihttp.NewMapParams(map[string]interface{}{
		"name": "jankin",
	}))
if err != nil {
	fmt.Println(err)
}

// 正常发送一个Get请求,追加get参数，以字符串格式
res, err := hihttp.New().Get(context.Background(), "http://www.google.com",hihttp.NewQueryParam("name=jankin"))
if err != nil {
	fmt.Println(err)
}
```



# POST 示例

```go
// -- application/x-www-form-urlencoded -- // 
// 以map的形式添加post参数
hihttp.New().Post(context.Background(), "http://www.yumontime.com/test/login",hihttp.NewWWWFormPayload(map[string]interface{}{
		"username": "jankin",
	}))


// -- application/json -- //
hihttp.New().SetHeader(hihttp.SerializationType,hihttp.SerializationTypeJSON).Post(context.Background(), "http://www.yumontime.com/test/login", hihttp.NewJSONPayload("username=jankin"))

```

# 超时处理
hihttp共有两种方式可以实现超时处理，具体说明和使用方式如下所示
1. hihttp.WithTimeout(time.Second)
```Go 
res, err := New(WithTimeout(time.Second)).Get(ctx, urlStr)
```
2. 在Get、Post等方法里传入一个带有timeout的context
```Go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
res, err := New(WithTimeout(6*time.Second)).Get(ctx, urlStr)
```
	注意：第二种方式的优先级将优于第一种，也就是说如果两种方式同时在使用，则以传入ctx的时间为实际超时时间。

# 请求重试
hihttp集成了请求重试，如果在请求失败(含请求超时)后，可以进行请求重试，即在延时一段时间后(可以是0秒)，重新发起请求。
```Go
urlStr := "https://www.google.com"
// 请求失败后会再次进行重试请求
ctx := context.Background()
res1, err := New(WithRetryCount(2)).Get(ctx, urlStr)
if err != nil {
	t.Error(err)
}
t.Log(string(res1))
```