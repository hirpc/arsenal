package hihttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

type Payload interface {
	Serialize() io.Reader
}

type jsonPayload struct {
	Payload string
}
func (p *jsonPayload) Serialize() io.Reader {
	return strings.NewReader(p.Payload)
}
// NewPayload 会根据序列化类型，生成一个payload
func NewJSONPayload(data interface{}) *jsonPayload {
	p := jsonPayload{}
	switch data.(type) {
	case string, []byte:
		p.Payload = fmt.Sprint(data)
	default:
		if b, err := json.Marshal(data); err != nil {
			return nil
		} else {
			p.Payload = string(b)
		}
	}
	return &p
}

type formPayload struct {
	Payload bytes.Buffer
}
func (p *formPayload) Serialize() io.Reader {
	return &p.Payload
}
// NewFormPayload 会根据序列化类型，生成一个payload
func NewFormPayload(data map[string]interface{}) *formPayload {
	p := formPayload{
		Payload: bytes.Buffer{},
	}
	writer := multipart.NewWriter(&p.Payload)
	for k, v := range data {
		_ = writer.WriteField(k, fmt.Sprint(v))
	}
	if err := writer.Close(); err != nil {
		return nil
	}
	return &p
}

// NewXMLPayload todo
func NewXMLPayload(v interface{})
