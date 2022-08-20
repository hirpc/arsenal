package txcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

type miniprogram struct {
	key, iv string
}

// NewMiniProgram is a constructor that should be used to create miniprogram struct.
func NewMiniProgram(key, iv string) *miniprogram {
	return &miniprogram{
		key: key,
		iv:  iv,
	}
}

func (m *miniprogram) Decrypt(data string) ([]byte, error) {
	k, err := base64.StdEncoding.DecodeString(m.key)
	if err != nil {
		return nil, err
	}
	i, err := base64.StdEncoding.DecodeString(m.iv)
	if err != nil {
		return nil, err
	}
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	b, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(b, i)
	var decrypted = make([]byte, len(d))
	mode.CryptBlocks(decrypted, d)
	return PKCS7UnPadding(decrypted), nil
}

func (m *miniprogram) DecryptAsParam(data string) (Param, error) {
	decryptedData, err := m.Decrypt(data)
	if err != nil {
		return nil, err
	}
	var p = make(Param)
	if err := json.Unmarshal(decryptedData, &p); err != nil {
		return nil, err
	}
	return p, nil
}
