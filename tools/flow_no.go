package tools

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	// 序列号生成器，用于同一毫秒内的唯一性保证
	sequenceCounter int64
	sequenceMutex   sync.Mutex
	lastTimestamp   int64
)

// GenerateFlowNo 生成流水号
// 格式：前缀(大写) + YYYYMMDDHHMMSS(14位) + 百毫秒(2位) + 序列号(3位) + 随机数(3位)
// 总长度：前缀长度 + 22位数字（与原来保持一致）
// 例如：SF2024010112000012345678
// 特点：
//   - 时间戳精确到秒，百毫秒部分单独处理（提高精度）
//   - 使用序列号 + 随机数确保唯一性
//   - 支持高并发场景
func GenerateFlowNo(prefix string) string {
	now := time.Now()
	prefix = strings.ToUpper(prefix)

	// 格式化为 YYYYMMDDHHMMSS
	dateStr := now.Format("20060102150405")

	// 获取百毫秒部分（0-99），即毫秒的前两位
	hundredMillis := (now.Nanosecond() / 1000000) / 10
	hundredMillisStr := fmt.Sprintf("%02d", hundredMillis)

	// 获取当前百毫秒级时间戳（用于序列号管理）
	currentTimestamp := now.UnixNano() / int64(100*time.Millisecond)

	// 生成序列号（同一百毫秒内递增）
	sequenceMutex.Lock()
	if currentTimestamp != lastTimestamp {
		sequenceCounter = 0
		lastTimestamp = currentTimestamp
	} else {
		sequenceCounter++
		// 序列号限制在 0-999，如果超过则等待下一百毫秒
		if sequenceCounter >= 1000 {
			// 等待到下一百毫秒
			for now.UnixNano()/(int64(100*time.Millisecond)) <= currentTimestamp {
				now = time.Now()
			}
			currentTimestamp = now.UnixNano() / int64(100*time.Millisecond)
			dateStr = now.Format("20060102150405")
			hundredMillis = (now.Nanosecond() / 1000000) / 10
			hundredMillisStr = fmt.Sprintf("%02d", hundredMillis)
			sequenceCounter = 0
			lastTimestamp = currentTimestamp
		}
	}
	seq := sequenceCounter
	sequenceMutex.Unlock()

	// 生成3位随机数（0-999，使用拒绝采样确保均匀分布）
	randomNum := generateUniformRandom(1000)

	// 组合：前缀 + 日期时间(秒，14位) + 百毫秒(2位) + 序列号(3位) + 随机数(3位)
	// 总长度：前缀长度 + 14 + 2 + 3 + 3 = 前缀长度 + 22（与原来保持一致）
	flowNo := fmt.Sprintf("%s%s%s%03d%03d", prefix, dateStr, hundredMillisStr, seq, randomNum)

	return flowNo
}

// generateUniformRandom 生成 [0, max) 范围内的均匀分布随机数
// 使用拒绝采样方法确保分布均匀
func generateUniformRandom(max int) int {
	if max <= 0 {
		return 0
	}

	// 计算需要的字节数
	// 为了确保均匀分布，我们需要生成足够大的随机数
	// 然后使用拒绝采样
	bytesNeeded := 4 // 4字节可以表示 0-4294967295
	randomBytes := make([]byte, bytesNeeded)

	// 计算拒绝采样的上限（最大的 max 的倍数，不超过 uint32 最大值）
	limit := (4294967295 / max) * max

	for {
		_, err := io.ReadFull(rand.Reader, randomBytes)
		if err != nil {
			// fallback: 使用时间戳 + 纳秒作为随机源
			now := time.Now()
			fallback := int(now.UnixNano() % int64(max))
			return fallback
		}

		// 将字节转换为 uint32
		randomValue := uint32(randomBytes[0])<<24 |
			uint32(randomBytes[1])<<16 |
			uint32(randomBytes[2])<<8 |
			uint32(randomBytes[3])

		// 拒绝采样：如果值在有效范围内，返回；否则重新生成
		if randomValue < uint32(limit) {
			return int(randomValue % uint32(max))
		}
	}
}
