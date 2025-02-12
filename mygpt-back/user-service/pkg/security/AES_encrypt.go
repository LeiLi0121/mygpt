package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// AES 密钥（必须是 16、24 或 32 字节）
var aesKey = []byte("your-32-byte-secret-key.") // 长度必须是 16/24/32 字节

// PKCS7Padding 对明文进行 PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去除 PKCS7 填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// encryptAES 使用 AES-CBC 模式加密
func EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize) // IV 这里固定为 0，你可以替换为随机 IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// PKCS7 填充
	paddedText := pkcs7Padding([]byte(plaintext), aes.BlockSize)

	// 加密
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptAES 使用 AES-CBC 模式解密
func DecryptAES(cipherText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize) // IV 必须和加密时相同
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	plainText := make([]byte, len(decoded))
	mode.CryptBlocks(plainText, decoded)

	// 去除填充
	plainText, err = pkcs7UnPadding(plainText)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
