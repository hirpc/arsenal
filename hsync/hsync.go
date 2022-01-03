package hsync

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type HSync interface {
	// Lock tries to get a locker for a resource
	// the parameter of ctx should have a trigger that Done() can be happended
	Lock(ctx context.Context, id string) error
	// Unlock will release a locker
	Unlock(id string) error
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
)

type hsync struct {
	r   *redis.Client
	key string
	opt Options
}

func New(r *redis.Client, k string, opts ...Option) HSync {
	var opt = options
	for _, o := range opts {
		o(&opt)
	}
	return &hsync{
		r:   r,
		key: "{lock}:" + k,
		opt: opt,
	}
}

func (h hsync) Lock(ctx context.Context, id string) error {
	for {
		select {
		case <-ctx.Done():
			return ErrLockFailed
		default:
			res, err := redis.NewScript(lockLUA).Run(
				ctx, h.r, []string{h.key}, id, h.opt.timeout.Milliseconds(),
			).Int()
			if err == nil && res == Success {
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

func (h hsync) Unlock(id string) error {
	res, err := redis.NewScript(unlockLUA).Run(
		context.Background(), h.r, []string{h.key}, id,
	).Int()
	if err != nil {
		return err
	}
	if res == Success {
		return nil
	}
	return ErrUnlockFailed
}
