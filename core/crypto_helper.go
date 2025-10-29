package core

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"errors"
)

// pkcs7Pad pads the data to a multiple of blockSize.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad removes padding from the data.
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:(length - unpadding)], nil
}

// ECBEncrypt encrypts data using AES ECB mode.
func ECBEncrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = pkcs7Pad(data, aes.BlockSize)
	encrypted := make([]byte, len(data))
	for bs, be := 0, block.BlockSize(); bs < len(data); bs, be = bs+be, be+be {
		block.Encrypt(encrypted[bs:be], data[bs:be])
	}
	return encrypted, nil
}

// ECBDecrypt decrypts data using AES ECB mode.
func ECBDecrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	for bs, be := 0, block.BlockSize(); bs < len(data); bs, be = bs+be, be+be {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return pkcs7Unpad(decrypted)
}

// AESCipher provides AES encryption and decryption.
type AESCipher struct {
	key []byte
}

// NewAESCipher creates a new AESCipher.
func NewAESCipher(key []byte) *AESCipher {
	return &AESCipher{key: key}
}

// Encrypt encrypts data.
func (c *AESCipher) Encrypt(raw []byte, useBase64 bool) ([]byte, error) {
	encrypted, err := ECBEncrypt(c.key, raw)
	if err != nil {
		return nil, err
	}
	if useBase64 {
		return []byte(base64.StdEncoding.EncodeToString(encrypted)), nil
	}
	return encrypted, nil
}

// Decrypt decrypts data.
func (c *AESCipher) Decrypt(enc []byte, useBase64 bool) ([]byte, error) {
	if useBase64 {
		var err error
		enc, err = base64.StdEncoding.DecodeString(string(enc))
		if err != nil {
			return nil, err
		}
	}
	return ECBDecrypt(c.key, enc)
}

// MD5 calculates the MD5 hash of data.
func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}
