// Package random for creating random strings
// Ref. https://www.oschina.net/code/snippet_170216_50650
package random

import (
	"math/rand"
	"time"
)

// Category represents the type of random output
type Category int

// Int to int
func (c Category) Int() int {
	return int(c)
}

const (
	// Number represents all outputs will be numbers
	Number Category = iota
	// LowerCase represents all outputs will contain characters in lower case
	LowerCase
	// UpperCase represents all outputs will contain characters in upper case
	UpperCase
	// AllCase represents all outputs will contain characters in any case
	AllCase
)

// Characters will generate a random string
func Characters(size int, kind Category) []byte {
	ikind, kinds, result := kind.Int(), [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < size; i++ {
		if kind == AllCase {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
