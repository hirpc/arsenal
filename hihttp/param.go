package hihttp

import (
	"fmt"
	"strings"
	"sync"
)

type Param interface {
	Marshal() string
}

// 对参数进行合并
func mergeParams(ps ...Param) string {
	params := []string{}
	for _, v := range ps {
		if v.Marshal() != "" {
			params = append(params, v.Marshal())
		}
	}
	return strings.Join(params, "&")
}

type mapParams struct {
	params map[string]interface{}
	mu     sync.Mutex
}

func NewMapParams(m map[string]interface{}) *mapParams {
	return &mapParams{
		params: m,
	}
}

// 对map进行相应的编码（urlencode等等）
func (m mapParams) Marshal() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	params := []string{}
	for k, v := range m.params {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(params, "&")
}

type kvParams struct {
	Key   string
	Value interface{}
}

func NewKVParams(key string, value interface{}) *kvParams {
	return &kvParams{
		Key:   key,
		Value: value,
	}
}

// 对kv 格式进行相应的编码（urlencode等等）
func (m kvParams) Marshal() string {
	return fmt.Sprintf("%s&%v", m.Key, m.Value)
}

type queryParams struct {
	Query string
}

func NewQueryParams(query string) *queryParams {
	return &queryParams{
		Query: query,
	}
}

// 直接使用url的query参数 返回
func (m queryParams) Marshal() string {
	return m.Query
}
