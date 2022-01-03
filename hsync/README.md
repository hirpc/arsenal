# 介绍
库封装了redis的分布式锁

## 特点
1. 锁安全，支持传入线程、进程、请求ID等等作为id，对某个资源加锁，避免被他人释放
2. 锁的操作采用LUA脚本保证其原子性
3. 支持redis的cluster环境
4. 支持自旋等待尝试加锁

## 常规使用方式
``` foo.go
import (
    "time"
    "github.com/hirpc/arsenal/hsync"
    "github.com/go-redis/redis/v8"
)

func Foo(ctx context.Context) error {
    locker := hsync.New(
        redis.Client, "test",
        // 以下均为可选选项
        // WithTimeout 设置key的默认超时时间
        // 一般设置为最大超时时间
        // 默认10秒
        WithTimeout(time.Second * 5),
        // WithWaitingPeriod 设置自旋等待间隔，一般为毫秒
        // 当加锁失败时候，会尝试自旋等待
        // 默认50毫秒
        WithWaitingPeriod(time.Second),
        // WithDisableRetry 禁止加锁重试
        // 是否关闭自旋重试
        // 如果携带此选项，则加锁失败后，会直接返回err
        // 默认开启
        WithDisableRetry(),
    )
    if err := locker.Lock(ctx, "id"); err != nil {
        return err
    }
    defer locker.Unlock("id")
    // ....
}
```

## TODO
1. 支持自我续命
2. 支持重入
3. 读写锁

## 参考
1. https://blog.csdn.net/grandachn/article/details/89032815