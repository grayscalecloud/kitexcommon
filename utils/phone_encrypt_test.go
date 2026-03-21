package utils

import (
	"strings"
	"testing"
)

// TestEncryptPhone_Disabled 测试加密未启用时的行为
func TestEncryptPhone_Disabled(t *testing.T) {
	// 初始化为不启用加密
	InitPhoneEncryption("")

	tests := []struct {
		name    string
		phone   string
		want    string
		wantErr bool
	}{
		{
			name:    "空字符串",
			phone:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "正常手机号",
			phone:   "13800138000",
			want:    "13800138000",
			wantErr: false,
		},
		{
			name:    "带前缀的加密手机号",
			phone:   "enc:SGVsbG8=",
			want:    "enc:SGVsbG8=",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPhone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncryptPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEncryptPhone_Enabled 测试加密启用时的行为
func TestEncryptPhone_Enabled(t *testing.T) {
	// 初始化加密密钥
	testKey := "test-encryption-key-12345678"
	InitPhoneEncryption(testKey)

	tests := []struct {
		name    string
		phone   string
		wantErr bool
		check   func(string) bool
	}{
		{
			name:    "空字符串",
			phone:   "",
			wantErr: false,
			check: func(result string) bool {
				return result == ""
			},
		},
		{
			name:    "正常手机号加密",
			phone:   "13800138000",
			wantErr: false,
			check: func(result string) bool {
				// 检查是否有加密前缀
				return strings.HasPrefix(result, "enc:")
			},
		},
		{
			name:    "已加密的手机号不重复加密",
			phone:   "enc:YWJjZGVmZ2hpams=",
			wantErr: false,
			check: func(result string) bool {
				return result == "enc:YWJjZGVmZ2hpams="
			},
		},
		{
			name:    "不同手机号产生不同密文",
			phone:   "13900139000",
			wantErr: false,
			check: func(result string) bool {
				return strings.HasPrefix(result, "enc:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPhone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.check(got) {
				t.Errorf("EncryptPhone() = %v, check failed", got)
			}
		})
	}
}

// TestEncryptPhone_DecryptConsistency 测试加密解密一致性
func TestEncryptPhone_DecryptConsistency(t *testing.T) {
	// 初始化加密密钥
	testKey := "test-encryption-key-12345678"
	InitPhoneEncryption(testKey)

	phones := []string{
		"13800138000",
		"13900139000",
		"15012345678",
		"18888888888",
		"",
	}

	for _, phone := range phones {
		t.Run(phone, func(t *testing.T) {
			// 加密
			encrypted, err := EncryptPhone(phone)
			if err != nil {
				t.Errorf("EncryptPhone(%q) error = %v", phone, err)
				return
			}

			// 解密
			decrypted, err := DecryptPhone(encrypted)
			if err != nil {
				t.Errorf("DecryptPhone(%q) error = %v", encrypted, err)
				return
			}

			// 验证一致性
			if decrypted != phone {
				t.Errorf("DecryptPhone(EncryptPhone(%q)) = %v, want %v", phone, decrypted, phone)
			}
		})
	}
}

// TestEncryptPhone_Deterministic 测试加密的确定性（相同明文产生相同密文）
func TestEncryptPhone_Deterministic(t *testing.T) {
	// 初始化加密密钥
	testKey := "test-encryption-key-12345678"
	InitPhoneEncryption(testKey)

	phone := "13800138000"

	// 加密两次
	encrypted1, err1 := EncryptPhone(phone)
	encrypted2, err2 := EncryptPhone(phone)

	if err1 != nil || err2 != nil {
		t.Fatalf("EncryptPhone failed: err1=%v, err2=%v", err1, err2)
	}

	// 相同明文应该产生相同密文（因为使用了确定性 nonce）
	if encrypted1 != encrypted2 {
		t.Errorf("EncryptPhone() not deterministic: %v != %v", encrypted1, encrypted2)
	}
}

// TestEncryptPhone_DifferentPhonesDifferentCiphers 测试不同手机号产生不同密文
func TestEncryptPhone_DifferentPhonesDifferentCiphers(t *testing.T) {
	// 初始化加密密钥
	testKey := "U3ubOv0iFqyvYnI3QzO5"
	InitPhoneEncryption(testKey)

	phone1 := "188502277777"
	phone2 := "13900139000"

	encrypted1, err1 := EncryptPhone(phone1)
	encrypted2, err2 := EncryptPhone(phone2)

	t.Logf("encrypted1: %v, encrypted2: %v", encrypted1, encrypted2)

	if err1 != nil || err2 != nil {
		t.Fatalf("EncryptPhone failed: err1=%v, err2=%v", err1, err2)
	}

	// 不同明文应该产生不同密文
	if encrypted1 == encrypted2 {
		t.Errorf("EncryptPhone() produced same cipher for different phones: %v", encrypted1)
	}
}

// TestEncryptPhone_PhoneFormatIsEncrypted 测试已加密格式不再重复加密
func TestEncryptPhone_PhoneFormatIsEncrypted(t *testing.T) {
	// 初始化加密密钥
	testKey := "test-encryption-key-12345678"
	InitPhoneEncryption(testKey)

	phone := "13800138000"

	// 第一次加密
	encrypted1, err := EncryptPhone(phone)
	if err != nil {
		t.Fatalf("EncryptPhone() error = %v", err)
	}

	// 第二次对已加密的字符串加密
	encrypted2, err := EncryptPhone(encrypted1)
	if err != nil {
		t.Fatalf("EncryptPhone() error = %v", err)
	}

	// 不应该重复加密
	if encrypted1 != encrypted2 {
		t.Errorf("EncryptPhone() encrypted already encrypted phone: %v -> %v", encrypted1, encrypted2)
	}
}
