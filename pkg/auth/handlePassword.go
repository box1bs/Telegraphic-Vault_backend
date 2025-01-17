package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func GenerateServerKey() (string, error) {
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil { // AES-256 key
		return "", err
	}

	return base64.StdEncoding.EncodeToString(randBytes), nil
}

func Decode(encrypted, tempKey string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(tempKey)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	plainText, err := unpad(data, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[length-1])
	if padding > blockSize || padding == 0 {
		return nil, errors.New("invalid padding")
	}

	for _, p := range data[length-padding:] {
		if int(p) != padding {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:length-padding], nil
}