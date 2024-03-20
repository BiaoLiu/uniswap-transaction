package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AesEncrypt 加密
func AesEncrypt(sKey string, data []byte) ([]byte, error) {
	key := []byte(sKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	encrypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return encrypted, nil
}

// AesDecrypt 解密
func AesDecrypt(sKey string, data []byte) ([]byte, error) {
	key := []byte(sKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(data))
	blockMode.CryptBlocks(decrypted, data)
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data error")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}
