package tools

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"
)

// GenerateFlowNo 生成流水号
func GenerateFlowNo(prefix string) string {
	// 生成以SF开头的服务费订单号，格式为 SF + YYYYMMDDHHMMSS + 6位随机数
	// 使用Unix毫秒时间戳替代纳秒时间，提高可读性和语义性
	now := time.Now()
	dateStr := now.Format("20060102150405")
	prefix = strings.ToUpper(prefix)
	// 使用crypto/rand生成更安全的随机数
	randomBytes := make([]byte, 4)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		// fallback到内置随机数生成器
		randomStr := fmt.Sprintf("%08d", now.Nanosecond()%100000000)
		return prefix + dateStr + randomStr
	}

	randomStr := fmt.Sprintf("%08d", (int(randomBytes[0])<<24|int(randomBytes[1])<<16|int(randomBytes[2])<<8|int(randomBytes[3]))%100000000)
	return prefix + dateStr + randomStr
}
