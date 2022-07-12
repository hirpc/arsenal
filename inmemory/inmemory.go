package inmemory

import (
	"errors"
	"sync"
	"time"
)

type inmemory struct {
	namedItems map[string]*item
	opt        Options
	mu         sync.RWMutex
}

func (i *inmemory) ndel(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.namedItems, name)
}

func (i *inmemory) nset(name string, v interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.namedItems[name] = &item{
		value:     v,
		timestamp: time.Now(),
	}
}

func (i *inmemory) nget(name string) interface{} {
	i.mu.RLock()
	if v, ok := i.namedItems[name]; ok {
		i.mu.RUnlock()
		if i.opt.maxLife == time.Duration(0) {
			return v.value
		}
		if v.timestamp.Add(i.opt.maxLife).After(time.Now()) {
			return v.value
		}
		// expired
		i.ndel(name)
		return nil
	}
	i.mu.RUnlock()
	return nil
}

var storage inmemory

func New(opts ...Option) *inmemory {
	opt := Options{
		maxLife: time.Duration(0),
	}
	for _, o := range opts {
		o(&opt)
	}
	storage = inmemory{
		namedItems: make(map[string]*item),
		opt:        opt,
	}
	return &storage
}

func (i *inmemory) Load() error {
	return nil
}

func (i *inmemory) Name() string {
	return "hrpc-inmemory"
}

func (i *inmemory) DependsOn() []string {
	return nil
}

func NSet(name string, v interface{}) {
	storage.nset(name, v)
}

var ErrNotExist = errors.New("does NOT found")

func NGet(name string) (interface{}, error) {
	if v := storage.nget(name); v != nil {
		return v, nil
	}
	return nil, ErrNotExist
}
