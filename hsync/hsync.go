package hsync

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type HSync interface {
	// Lock tries to get a locker for a resource
	// the parameter of ctx should have a trigger that Done() can be happended
	Lock(ctx context.Context, id any) error
	// Unlock will release a locker
	Unlock(id any) error

	// RLock will try to add a read lock to a resource
	RLock(ctx context.Context, id any) error
	// RUnlock will try to release the read lock
	RUnlock(id any) error
}

const (
	Success = 1
	Fail    = 0
)

var (
	ErrLockFailed   = errors.New("lock failed")
	ErrUnlockFailed = errors.New("unlock failed")
)

var (
	//go:embed scripts/lock.lua
	lockLUA string
	//go:embed scripts/unlock.lua
	unlockLUA string
	//go:embed scripts/rlock.lua
	rlockLUA string
	//go:embed scripts/runlock.lua
	runlockLUA string
)

type hsync struct {
	r                        *redis.Client
	wkey, reentrantKey, rkey string
	opt                      Options
}

func New(r *redis.Client, key any, opts ...Option) HSync {
	var opt = options
	for _, o := range opts {
		o(&opt)
	}
	k := fmt.Sprint(key)
	return &hsync{
		r:            r,
		wkey:         "write_{lock}:" + k,
		rkey:         "read_{lock}:",
		reentrantKey: "reentrant_{lock}:" + k,
		opt:          opt,
	}
}

func (h hsync) Lock(ctx context.Context, id any) error {
	for {
		select {
		case <-ctx.Done():
			return ErrLockFailed
		default:
			res, err := redis.NewScript(lockLUA).Run(
				ctx, h.r, []string{h.wkey, h.rkey, h.reentrantKey}, id, h.opt.timeout.Milliseconds(),
			).Int()
			if err == nil && res > 0 {
				return nil
			}
			if !h.opt.retry {
				// Maybe err or nil depends on the result got from Run()
				return err
			}
			time.Sleep(h.opt.waitingPeriod)
		}
	}
}

func (h hsync) Unlock(id any) error {
	res, err := redis.NewScript(unlockLUA).Run(
		context.Background(), h.r, []string{h.wkey, h.reentrantKey}, id,
	).Int()
	if err != nil {
		return err
	}
	if res != -1 {
		return nil
	}
	return ErrUnlockFailed
}

func (h hsync) RLock(ctx context.Context, id any) error {
	for {
		select {
		case <-ctx.Done():
			return ErrLockFailed
		default:
			res, err := redis.NewScript(rlockLUA).Run(
				ctx, h.r, []string{h.rkey, h.wkey}, id, h.opt.timeout.Milliseconds(),
			).Int()
			if err == nil && res > 0 {
				return nil
			}
			if !h.opt.retry {
				// Maybe err or nil depends on the result got from Run()
				return err
			}
			time.Sleep(h.opt.waitingPeriod)
		}
	}
}

func (h hsync) RUnlock(id any) error {
	res, err := redis.NewScript(runlockLUA).Run(
		context.Background(), h.r, []string{h.rkey, h.wkey}, id,
	).Int()
	if err != nil {
		return err
	}
	if res != -1 {
		return nil
	}
	return ErrUnlockFailed
}
