package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// MD5Hash MD5哈希值
func MD5Hash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5HashString MD5哈希值
func MD5HashString(s string) string {
	return MD5Hash([]byte(s))
}

// SHA1Hash SHA1哈希值
func SHA1Hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1HashString SHA1哈希值
func SHA1HashString(s string) string {
	return SHA1Hash([]byte(s))
}

// SaltedHmac salt SHA1哈希值
func SaltedHmac(keySalt string, value string) string {
	h := sha1.New()
	h.Write([]byte(keySalt))
	key := h.Sum(nil)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	t := mac.Sum(nil)
	s := hex.EncodeToString(t)
	return s
}
