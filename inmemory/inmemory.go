package inmemory

import (
	"errors"
	"sync"
)

type inmemory struct {
	namedItems map[string]interface{}
	mu         sync.RWMutex
}

func (i *inmemory) nset(name string, v interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.namedItems[name] = v
}

func (i *inmemory) nget(name string) interface{} {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if v, ok := i.namedItems[name]; ok {
		return v
	}
	return nil
}

var storage inmemory

func New() *inmemory {
	storage = inmemory{
		namedItems: make(map[string]interface{}),
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

var ErrNotExist = errors.New("Does NOT found")

func NGet(name string) (interface{}, error) {
	if v := storage.nget(name); v != nil {
		return v, nil
	}
	return nil, ErrNotExist
}
