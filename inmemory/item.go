package inmemory

import "time"

type item struct {
	value     interface{}
	timestamp time.Time
}
