package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type Crypto struct {
	gcm cipher.AEAD
}

func NewEncrypt(secretKey string) (*Crypto, error) {
	block, err := aes.NewCipher([]byte(secretKey))

	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Crypto{
		gcm: gcm,
	}, nil
}

func (c *Crypto) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("encrypt a one-time key err %v", err)
	}
	ciphertext := c.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, err
}
func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	fmt.Println(" ht ", string(nonce), string(encryptedData), len(nonce), len(encryptedData), nonceSize)
	result, err := c.gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}
	return result, err
}
