package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	// 密钥环境变量名
	encryptKeyEnvName = "PHONE_ENCRYPT_KEY"
)

var (
	// phoneEncryptKey 手机号加密密钥（32字节，用于AES-256）
	phoneEncryptKey []byte
	// isEncryptionEnabled 是否启用加密（如果未配置密钥，则不加密）
	isEncryptionEnabled bool
	// isKeyInitialized 密钥是否已初始化（防止重复初始化）
	isKeyInitialized bool
)

// init 初始化加密密钥（从环境变量读取）
func init() {
	// 优先从环境变量读取
	key := os.Getenv(encryptKeyEnvName)
	if key == "" {
		// 如果没有配置密钥，不启用加密，直接存储明文
		isEncryptionEnabled = false
		phoneEncryptKey = nil
		isKeyInitialized = true
		return
	}

	// 配置了密钥，启用加密
	initEncryptionKey(key)
}

// InitPhoneEncryption 初始化手机号加密密钥
// 如果 key 为空，则不启用加密（存储明文）
// 如果 key 不为空，则启用加密
// 注意：此函数应该在应用启动时调用
func InitPhoneEncryption(key string) {
	initEncryptionKey(key)
}

// initEncryptionKey 内部函数，用于初始化加密密钥
func initEncryptionKey(key string) {
	if key == "" {
		// 如果没有配置密钥，不启用加密，直接存储明文
		isEncryptionEnabled = false
		phoneEncryptKey = nil
		isKeyInitialized = true
		return
	}

	// 配置了密钥，启用加密
	isEncryptionEnabled = true
	// 将密钥转换为32字节（AES-256需要32字节密钥）
	hash := sha256.Sum256([]byte(key))
	phoneEncryptKey = hash[:]
	isKeyInitialized = true
}

// EncryptPhone 加密手机号码
// 使用 AES-256-GCM 加密算法，返回 base64 编码的加密字符串
// 如果未配置加密密钥，则直接返回明文（不加密）
func EncryptPhone(phone string) (string, error) {
	if phone == "" {
		return "", nil
	}

	// 如果未启用加密，直接返回明文
	if !isEncryptionEnabled {
		return phone, nil
	}

	// 检查是否为已加密的手机号（避免重复加密）
	if IsEncryptedPhone(phone) {
		return phone, nil
	}

	// 创建 AES cipher
	block, err := aes.NewCipher(phoneEncryptKey)
	if err != nil {
		return "", fmt.Errorf("创建加密器失败: %w", err)
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %w", err)
	}

	// 生成确定性 nonce（基于明文和密钥，确保相同明文产生相同密文）
	// 使用 HMAC-SHA256 生成固定长度的 nonce
	nonceSize := gcm.NonceSize()
	h := hmac.New(sha256.New, phoneEncryptKey)
	h.Write([]byte(phone))
	nonceHash := h.Sum(nil)
	nonce := nonceHash[:nonceSize] // 取前12字节作为 nonce

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(phone), nil)

	// 返回 base64 编码的加密字符串，添加前缀标识
	return "enc:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPhone 解密手机号码
func DecryptPhone(encryptedPhone string) (string, error) {
	if encryptedPhone == "" {
		return "", nil
	}

	// 检查是否为加密的手机号
	if !IsEncryptedPhone(encryptedPhone) {
		// 如果不是加密格式，直接返回（兼容旧数据）
		return encryptedPhone, nil
	}

	// 移除前缀
	encryptedPhone = strings.TrimPrefix(encryptedPhone, "enc:")

	// 解码 base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPhone)
	if err != nil {
		return "", fmt.Errorf("解码失败: %w", err)
	}

	// 创建 AES cipher
	block, err := aes.NewCipher(phoneEncryptKey)
	if err != nil {
		return "", fmt.Errorf("创建解密器失败: %w", err)
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %w", err)
	}

	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	// 提取 nonce 和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}

// IsEncryptedPhone 检查手机号是否已加密
func IsEncryptedPhone(phone string) bool {
	return strings.HasPrefix(phone, "enc:")
}

// GetEncryptKey 获取当前使用的加密密钥（用于调试，不返回实际密钥值）
func GetEncryptKey() string {
	if !isEncryptionEnabled {
		return "未配置加密密钥，手机号以明文存储"
	}
	key := os.Getenv(encryptKeyEnvName)
	if key == "" {
		return "使用默认密钥（开发环境）"
	}
	return "使用环境变量配置的密钥"
}

// IsEncryptionEnabled 检查是否启用了加密
func IsEncryptionEnabled() bool {
	return isEncryptionEnabled
}

// IsPhoneNumber 检查字符串是否是手机号格式（11位数字，1开头）
func IsPhoneNumber(s string) bool {
	if len(s) != 11 {
		return false
	}
	// 检查是否全为数字且以1开头
	if s[0] != '1' {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// NormalizePhoneForQuery 规范化手机号用于查询
// 如果传入的是明文手机号，则加密后返回
// 如果传入的是密文（带 enc: 前缀），则直接返回
// 这样可以让查询方法同时支持明文和密文查询
func NormalizePhoneForQuery(phone string) (string, error) {
	if phone == "" {
		return "", nil
	}

	// 如果已经是加密格式，直接返回
	if IsEncryptedPhone(phone) {
		return phone, nil
	}

	// 否则加密后返回
	return EncryptPhone(phone)
}
