package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// GenerateHash 生成密码哈希
func GenerateHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// VerifyPassword 验证密码
func VerifyPassword(password, hash string) bool {
	hashedPassword := GenerateHash(password)
	return hashedPassword == hash
}

// EncryptAES 使用AES-256加密敏感数据
func EncryptAES(plaintext []byte, key []byte) ([]byte, error) {
	// 确保密钥长度为32字节（AES-256）
	if len(key) != 32 {
		hashKey := sha256.Sum256(key)
		key = hashKey[:]
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES 使用AES-256解密敏感数据
func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	// 确保密钥长度为32字节（AES-256）
	if len(key) != 32 {
		hashKey := sha256.Sum256(key)
		key = hashKey[:]
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 检查密文长度
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("密文太短")
	}

	// 提取nonce
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptString 加密字符串
func EncryptString(plaintext string, key string) (string, error) {
	// 加密
	ciphertext, err := EncryptAES([]byte(plaintext), []byte(key))
	if err != nil {
		return "", err
	}

	// 编码为Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString 解密字符串
func DecryptString(ciphertext string, key string) (string, error) {
	// 解码Base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 解密
	plaintext, err := DecryptAES(data, []byte(key))
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}