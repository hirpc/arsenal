// Package txcrypto should be used for tencent products.
// ref. https://www.kancloud.cn/idzqj/customer/1421041
package txcrypto

import (
	"bytes"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Param can store all decrypted values as a map
type Param map[string]interface{}

// String will try to get a value from decrypted map.
// Empty will be returned if not found.
func (p Param) String(k string) string {
	if val, ok := p[k]; ok {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return ""
}

// Value same as String, but returned an interface{}
func (p Param) Value(k string) interface{} {
	if val, ok := p[k]; ok {
		return val
	}
	return nil
}
