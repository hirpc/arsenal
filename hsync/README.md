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
    "github.com/hirpc/arsenal/hsync"
    "github.com/go-redis/redis/v8"
)

func Foo(ctx context.Context) error {
    locker := hsync.New(redis.Client, "test")
    if err := locker.Lock(ctx, "id"); err != nil {
        return err
    }
    defer locker.Unlock("id")
    // ....
}
```

## TODO
1. 支持自我续命