package utils

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"runtime"
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

// GetMachineNodeID 根据机器信息自动生成nodeID
// 使用机器的主机名、MAC地址、IP地址等信息生成唯一的nodeID
func GetMachineNodeID() (int64, error) {
	// 收集机器标识信息
	var machineInfo string

	// 1. 获取主机名
	if hostname, err := os.Hostname(); err == nil {
		machineInfo += hostname
	}

	// 2. 获取MAC地址
	if macAddr, err := getMACAddress(); err == nil {
		machineInfo += macAddr
	}

	// 3. 获取本机IP地址
	if localIP, err := getLocalIP(); err == nil {
		machineInfo += localIP
	}

	// 4. 获取操作系统信息
	machineInfo += runtime.GOOS
	machineInfo += runtime.GOARCH

	// 5. 获取进程ID
	machineInfo += fmt.Sprintf("%d", os.Getpid())

	// 如果无法获取足够的信息，使用时间戳作为备选
	if machineInfo == "" {
		machineInfo = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// 使用MD5哈希生成固定长度的标识
	hash := md5.Sum([]byte(machineInfo))

	// 将哈希的前8字节转换为int64，并确保在有效范围内
	nodeID := int64(hash[0])<<24 | int64(hash[1])<<16 | int64(hash[2])<<8 | int64(hash[3])

	// 确保nodeID在雪花算法的有效范围内 (0-1023，10位)
	nodeID = nodeID & 0x3FF // 取低10位

	return nodeID, nil
}

// getMACAddress 获取第一个非回环网络接口的MAC地址
func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过回环接口和无效接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取MAC地址
		if len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("未找到有效的MAC地址")
}

// getLocalIP 获取本机IP地址
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// NewSnowflakeWithAutoNodeID 使用自动生成的nodeID创建雪花算法实例
func NewSnowflakeWithAutoNodeID() (*Snowflake, error) {
	nodeID, err := GetMachineNodeID()
	if err != nil {
		return nil, fmt.Errorf("获取机器nodeID失败: %w", err)
	}

	return NewSnowflake(nodeID), nil
}

// GetNodeIDFromEnv 从环境变量获取nodeID，如果不存在则自动生成
func GetNodeIDFromEnv() (int64, error) {
	// 首先尝试从环境变量获取
	if nodeIDStr := os.Getenv("NODE_ID"); nodeIDStr != "" {
		var nodeID int64
		if _, err := fmt.Sscanf(nodeIDStr, "%d", &nodeID); err == nil {
			// 确保nodeID在有效范围内
			if nodeID >= 0 && nodeID <= 1023 {
				return nodeID, nil
			}
		}
	}

	// 如果环境变量无效或不存在，则自动生成
	return GetMachineNodeID()
}

// NewSnowflakeWithEnvNodeID 优先使用环境变量中的nodeID，否则自动生成
func NewSnowflakeWithEnvNodeID() (*Snowflake, error) {
	nodeID, err := GetNodeIDFromEnv()
	if err != nil {
		return nil, fmt.Errorf("获取nodeID失败: %w", err)
	}

	return NewSnowflake(nodeID), nil
}
