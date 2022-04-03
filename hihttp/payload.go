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
	ContentType() string
}

type jsonPayload struct {
	Payload string
}

func (p *jsonPayload) Serialize() io.Reader {
	return strings.NewReader(p.Payload)
}
func (p *jsonPayload) ContentType() string {
	return SerializationTypeJSON
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
	Payload map[string]string
}

func (p *formPayload) Serialize() io.Reader {
	paylod := &bytes.Buffer{}
	writer := multipart.NewWriter(paylod)
	for k, v := range p.Payload {
		_ = writer.WriteField(k, v)
	}
	if err := writer.Close(); err != nil {
		return nil
	}
	return paylod
}
func (p *formPayload) ContentType() string {
	return SerializationTypeFormData
}

// NewFormPayload 会根据序列化类型，生成一个payload
func NewFormPayload(data map[string]interface{}) *formPayload {
	p := formPayload{
		Payload: map[string]string{},
	}
	for k, v := range data {
		p.Payload[k] = fmt.Sprint(v)
	}

	return &p
}