package security

import (
	"fmt"
	"testing"
)

func Testaes(t *testing.T) {
	msg := "12345"
	encrypted, err := EncryptAES(msg)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	fmt.Println("加密后的数据:", encrypted)
	decrypted, err := DecryptAES(encrypted)

	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	fmt.Println("解密后的数据:", decrypted)
}
