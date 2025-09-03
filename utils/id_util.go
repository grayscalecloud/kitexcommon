package utils

import (
	"sync"
	"time"
)

const (
	// 定义字符集
	charSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Snowflake 雪花算法结构体
type Snowflake struct {
	mu         sync.Mutex
	nodeID     int64
	sequence   int64
	lastMillis int64
}

// NewSnowflake 创建新的雪花算法实例
func NewSnowflake(nodeID int64) *Snowflake {
	return &Snowflake{
		nodeID: nodeID,
	}
}

// NextID 生成下一个 ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1_000_000 // 毫秒级时间

	if now == s.lastMillis {
		s.sequence = (s.sequence + 1) & 0xFFF // 12 位序列号
	} else {
		s.sequence = 0 // 不同毫秒重置序列号
	}

	s.lastMillis = now

	// 组合 ID
	id := ((now & 0x1FFFFFFFFFFF) << 22) | (s.nodeID << 12) | s.sequence
	return id
}

// EncodeToShortID 将 ID 转换为短标识符
func (s *Snowflake) EncodeToShortID() string {
	id := s.NextID()
	base := int64(len(charSet))
	var shortID string

	for id > 0 {
		remainder := id % base
		shortID = string(charSet[remainder]) + shortID
		id /= base
	}

	return shortID
}
