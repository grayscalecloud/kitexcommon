package utils

import (
	"regexp"
	"strings"
)

// DesensitizeType 脱敏类型
type DesensitizeType int

const (
	// DesensitizeTypeNone 不脱敏
	DesensitizeTypeNone DesensitizeType = iota
	// DesensitizeTypePhone 手机号脱敏
	DesensitizeTypePhone
	// DesensitizeTypeEmail 邮箱脱敏
	DesensitizeTypeEmail
	// DesensitizeTypeIDCard 身份证脱敏
	DesensitizeTypeIDCard
	// DesensitizeTypeBankCard 银行卡脱敏
	DesensitizeTypeBankCard
	// DesensitizeTypeName 姓名脱敏
	DesensitizeTypeName
	// DesensitizeTypeAddress 地址脱敏
	DesensitizeTypeAddress
	// DesensitizeTypeCustom 自定义脱敏
	DesensitizeTypeCustom
)

// Desensitizer 脱敏器
type Desensitizer struct {
	// 脱敏字符，默认为*
	maskChar string
	// 是否保留首尾字符
	keepHeadTail bool
	// 保留的首字符数量
	keepHeadCount int
	// 保留的尾字符数量
	keepTailCount int
}

// NewDesensitizer 创建脱敏器
func NewDesensitizer() *Desensitizer {
	return &Desensitizer{
		maskChar:      "*",
		keepHeadTail:  true,
		keepHeadCount: 3,
		keepTailCount: 4,
	}
}

// SetMaskChar 设置脱敏字符
func (d *Desensitizer) SetMaskChar(char string) *Desensitizer {
	d.maskChar = char
	return d
}

// SetKeepHeadTail 设置是否保留首尾字符
func (d *Desensitizer) SetKeepHeadTail(keep bool) *Desensitizer {
	d.keepHeadTail = keep
	return d
}

// SetKeepCount 设置保留的首尾字符数量
func (d *Desensitizer) SetKeepCount(headCount, tailCount int) *Desensitizer {
	d.keepHeadCount = headCount
	d.keepTailCount = tailCount
	return d
}

// Desensitize 通用脱敏方法
func (d *Desensitizer) Desensitize(data string, desensitizeType DesensitizeType) string {
	if data == "" {
		return data
	}

	switch desensitizeType {
	case DesensitizeTypePhone:
		return d.DesensitizePhone(data)
	case DesensitizeTypeEmail:
		return d.DesensitizeEmail(data)
	case DesensitizeTypeIDCard:
		return d.DesensitizeIDCard(data)
	case DesensitizeTypeBankCard:
		return d.DesensitizeBankCard(data)
	case DesensitizeTypeName:
		return d.DesensitizeName(data)
	case DesensitizeTypeAddress:
		return d.DesensitizeAddress(data)
	case DesensitizeTypeCustom:
		return d.DesensitizeCustom(data)
	default:
		return data
	}
}

// DesensitizePhone 手机号脱敏
func (d *Desensitizer) DesensitizePhone(phone string) string {
	// 验证手机号格式
	if !d.isValidPhone(phone) {
		return phone
	}

	if len(phone) == 11 {
		// 中国手机号：保留前3位和后4位
		return phone[:3] + strings.Repeat(d.maskChar, 4) + phone[7:]
	}
	return d.DesensitizeCustom(phone)
}

// DesensitizeEmail 邮箱脱敏
func (d *Desensitizer) DesensitizeEmail(email string) string {
	// 验证邮箱格式
	if !d.isValidEmail(email) {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	// 用户名脱敏
	if len(username) <= 2 {
		username = strings.Repeat(d.maskChar, len(username))
	} else {
		username = username[:1] + strings.Repeat(d.maskChar, len(username)-2) + username[len(username)-1:]
	}

	return username + "@" + domain
}

// DesensitizeIDCard 身份证脱敏
func (d *Desensitizer) DesensitizeIDCard(idCard string) string {
	// 验证身份证格式
	if !d.isValidIDCard(idCard) {
		return idCard
	}

	if len(idCard) == 18 {
		// 18位身份证：保留前6位和后4位
		return idCard[:6] + strings.Repeat(d.maskChar, 8) + idCard[14:]
	} else if len(idCard) == 15 {
		// 15位身份证：保留前6位和后3位
		return idCard[:6] + strings.Repeat(d.maskChar, 6) + idCard[12:]
	}
	return d.DesensitizeCustom(idCard)
}

// DesensitizeBankCard 银行卡脱敏
func (d *Desensitizer) DesensitizeBankCard(bankCard string) string {
	// 移除空格和连字符
	cleanCard := strings.ReplaceAll(strings.ReplaceAll(bankCard, " ", ""), "-", "")

	// 验证银行卡格式（16-19位数字）
	if !d.isValidBankCard(cleanCard) {
		return bankCard
	}

	if len(cleanCard) >= 16 {
		// 保留前4位和后4位
		return cleanCard[:4] + strings.Repeat(d.maskChar, len(cleanCard)-8) + cleanCard[len(cleanCard)-4:]
	}
	return d.DesensitizeCustom(cleanCard)
}

// DesensitizeName 姓名脱敏
func (d *Desensitizer) DesensitizeName(name string) string {
	if name == "" || len(name) == 1 {
		return name
	}

	runes := []rune(name)
	length := len(runes)

	if length <= 1 {
		return strings.Repeat(d.maskChar, length)
	} else if length == 2 {
		return string(runes[0]) + d.maskChar
	} else {
		// 保留第一个字符，其余用*替代
		return string(runes[0]) + strings.Repeat(d.maskChar, length-1)
	}
}

// DesensitizeAddress 地址脱敏
func (d *Desensitizer) DesensitizeAddress(address string) string {
	if address == "" {
		return address
	}

	runes := []rune(address)
	length := len(runes)

	if length <= 6 {
		// 短地址：保留前2位，其余脱敏
		return string(runes[:2]) + strings.Repeat(d.maskChar, length-2)
	} else {
		// 长地址：保留前3位和后3位
		return string(runes[:3]) + strings.Repeat(d.maskChar, length-6) + string(runes[length-3:])
	}
}

// DesensitizeCustom 自定义脱敏
func (d *Desensitizer) DesensitizeCustom(data string) string {
	if data == "" {
		return data
	}

	runes := []rune(data)
	length := len(runes)

	if !d.keepHeadTail {
		// 不保留首尾，全部脱敏
		return strings.Repeat(d.maskChar, length)
	}

	if length <= d.keepHeadCount+d.keepTailCount {
		// 长度不足，只保留首字符
		if length <= 1 {
			return d.maskChar
		}
		return string(runes[0]) + strings.Repeat(d.maskChar, length-1)
	}

	// 保留首尾字符
	head := string(runes[:d.keepHeadCount])
	tail := string(runes[length-d.keepTailCount:])
	middle := strings.Repeat(d.maskChar, length-d.keepHeadCount-d.keepTailCount)

	return head + middle + tail
}

// 验证方法

// isValidPhone 验证手机号格式
func (d *Desensitizer) isValidPhone(phone string) bool {
	// 中国手机号正则：1开头的11位数字
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// isValidEmail 验证邮箱格式
func (d *Desensitizer) isValidEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

// isValidIDCard 验证身份证格式
func (d *Desensitizer) isValidIDCard(idCard string) bool {
	// 18位身份证正则
	matched18, _ := regexp.MatchString(`^\d{17}[\dXx]$`, idCard)
	// 15位身份证正则
	matched15, _ := regexp.MatchString(`^\d{15}$`, idCard)
	return matched18 || matched15
}

// isValidBankCard 验证银行卡格式
func (d *Desensitizer) isValidBankCard(bankCard string) bool {
	// 银行卡：16-19位数字
	matched, _ := regexp.MatchString(`^\d{16,19}$`, bankCard)
	return matched
}

// 便捷方法

// DesensitizePhone 手机号脱敏（便捷方法）
func DesensitizePhone(phone string) string {
	return NewDesensitizer().DesensitizePhone(phone)
}

// DesensitizeEmail 邮箱脱敏（便捷方法）
func DesensitizeEmail(email string) string {
	return NewDesensitizer().DesensitizeEmail(email)
}

// DesensitizeIDCard 身份证脱敏（便捷方法）
func DesensitizeIDCard(idCard string) string {
	return NewDesensitizer().DesensitizeIDCard(idCard)
}

// DesensitizeBankCard 银行卡脱敏（便捷方法）
func DesensitizeBankCard(bankCard string) string {
	return NewDesensitizer().DesensitizeBankCard(bankCard)
}

// DesensitizeName 姓名脱敏（便捷方法）
func DesensitizeName(name string) string {
	return NewDesensitizer().DesensitizeName(name)
}

// DesensitizeAddress 地址脱敏（便捷方法）
func DesensitizeAddress(address string) string {
	return NewDesensitizer().DesensitizeAddress(address)
}

// DesensitizeCustom 自定义脱敏（便捷方法）
func DesensitizeCustom(data string, headCount, tailCount int) string {
	return NewDesensitizer().
		SetKeepCount(headCount, tailCount).
		DesensitizeCustom(data)
}

// 检测脱敏数据的方法

// IsDesensitizedPhone 检测是否为脱敏的手机号
func IsDesensitizedPhone(phone string) bool {
	if phone == "" {
		return false
	}
	// 检查是否包含脱敏字符（*）且长度符合脱敏后的格式
	// 脱敏后的手机号格式：138****5678
	matched, _ := regexp.MatchString(`^1[3-9]\d\*\*\*\*\d{4}$`, phone)
	return matched
}

// IsDesensitizedEmail 检测是否为脱敏的邮箱
func IsDesensitizedEmail(email string) bool {
	if email == "" {
		return false
	}
	// 检查是否包含脱敏字符（*）且格式符合脱敏后的邮箱
	// 脱敏后的邮箱格式：t**t@example.com 或 *@b.com
	// 使用更严格的正则：必须包含至少2个连续的*号
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]*\*\*+[a-zA-Z0-9._%+-]*@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

// IsDesensitizedData 检测是否为脱敏数据（通用方法）
func IsDesensitizedData(data string) bool {
	if data == "" {
		return false
	}
	// 检查是否包含脱敏字符（*）且长度大于1
	return strings.Contains(data, "*") && len(data) > 1
}
