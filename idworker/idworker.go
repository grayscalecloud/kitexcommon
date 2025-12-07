package idworker

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// IdWorker 雪花算法ID生成器(完整版) 支持两个参数，workerId 和 datacenterId
// workerId 是机器ID，datacenterId 是数据中心ID
// workerId 和 datacenterId 的取值范围是 0-31
// workerId 和 datacenterId 的取值范围是 0-31
type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func (iw *IdWorker) InitIdWorker(workerId, datacenterId int64) error {

	var baseValue int64 = -1
	iw.startTime = 1463834116272
	iw.workerIdBits = 5
	iw.datacenterIdBits = 5
	iw.maxWorkerId = baseValue ^ (baseValue << iw.workerIdBits)
	iw.maxDatacenterId = baseValue ^ (baseValue << iw.datacenterIdBits)
	iw.sequenceBits = 12
	iw.workerIdLeftShift = iw.sequenceBits
	iw.datacenterIdLeftShift = iw.workerIdBits + iw.workerIdLeftShift
	iw.timestampLeftShift = iw.datacenterIdBits + iw.datacenterIdLeftShift
	iw.sequenceMask = baseValue ^ (baseValue << iw.sequenceBits)
	iw.sequence = 0
	iw.lastTimestamp = -1
	iw.signMask = ^baseValue + 1

	iw.idLock = &sync.Mutex{}

	// 先赋值，再检查
	iw.workerId = workerId
	iw.datacenterId = datacenterId

	if iw.workerId < 0 || iw.workerId > iw.maxWorkerId {
		return fmt.Errorf("workerId[%v] is less than 0 or greater than maxWorkerId[%v]", workerId, iw.maxWorkerId)
	}
	if iw.datacenterId < 0 || iw.datacenterId > iw.maxDatacenterId {
		return fmt.Errorf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d]", datacenterId, iw.maxDatacenterId)
	}
	return nil
}

func (iw *IdWorker) NextId() (int64, error) {
	iw.idLock.Lock()
	defer iw.idLock.Unlock()

	// 使用毫秒级时间戳，与 startTime 保持一致
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if timestamp < iw.lastTimestamp {
		return -1, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", iw.lastTimestamp-timestamp)
	}

	if timestamp == iw.lastTimestamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask
		if iw.sequence == 0 {
			timestamp = iw.tilNextMillis()
			iw.sequence = 0
		}
	} else {
		iw.sequence = 0
	}

	iw.lastTimestamp = timestamp

	id := ((timestamp - iw.startTime) << iw.timestampLeftShift) |
		(iw.datacenterId << iw.datacenterIdLeftShift) |
		(iw.workerId << iw.workerIdLeftShift) |
		iw.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (iw *IdWorker) tilNextMillis() int64 {
	// 获取当前毫秒级时间戳
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	// 如果当前时间戳小于等于上次时间戳，则等待到下一个毫秒
	for timestamp <= iw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}

// GetWorkerIdFromEnv 从环境变量获取 workerId，如果不存在则返回默认值和错误
func GetWorkerIdFromEnv() (int64, error) {
	if workerIdStr := os.Getenv("IDWORKER_WORKER_ID"); workerIdStr != "" {
		workerId, err := strconv.ParseInt(workerIdStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("环境变量 IDWORKER_WORKER_ID 格式错误: %w", err)
		}
		if workerId < 0 || workerId > 31 {
			return 0, fmt.Errorf("环境变量 IDWORKER_WORKER_ID 超出范围 [0-31]: %d", workerId)
		}
		return workerId, nil
	}
	return 0, fmt.Errorf("环境变量 IDWORKER_WORKER_ID 未设置")
}

// GetDatacenterIdFromEnv 从环境变量获取 datacenterId，如果不存在则返回默认值和错误
func GetDatacenterIdFromEnv() (int64, error) {
	if datacenterIdStr := os.Getenv("IDWORKER_DATACENTER_ID"); datacenterIdStr != "" {
		datacenterId, err := strconv.ParseInt(datacenterIdStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("环境变量 IDWORKER_DATACENTER_ID 格式错误: %w", err)
		}
		if datacenterId < 0 || datacenterId > 31 {
			return 0, fmt.Errorf("环境变量 IDWORKER_DATACENTER_ID 超出范围 [0-31]: %d", datacenterId)
		}
		return datacenterId, nil
	}
	return 0, fmt.Errorf("环境变量 IDWORKER_DATACENTER_ID 未设置")
}

// GetMachineWorkerId 根据机器信息自动生成 workerId (0-31)
// 使用机器的主机名、MAC地址、IP地址等信息生成唯一的 workerId
func GetMachineWorkerId() (int64, error) {
	machineInfo, err := collectMachineInfo()
	if err != nil {
		return 0, fmt.Errorf("收集机器信息失败: %w", err)
	}

	// 使用MD5哈希生成固定长度的标识
	hash := md5.Sum([]byte(machineInfo + "_worker"))

	// 将哈希转换为 int64，并确保在有效范围内 (0-31)
	workerId := int64(hash[0]) & 0x1F // 取低5位 (0-31)

	return workerId, nil
}

// GetMachineDatacenterId 根据机器信息自动生成 datacenterId (0-31)
// 使用机器的主机名、MAC地址、IP地址等信息生成唯一的 datacenterId
func GetMachineDatacenterId() (int64, error) {
	machineInfo, err := collectMachineInfo()
	if err != nil {
		return 0, fmt.Errorf("收集机器信息失败: %w", err)
	}

	// 使用MD5哈希生成固定长度的标识
	hash := md5.Sum([]byte(machineInfo + "_datacenter"))

	// 将哈希转换为 int64，并确保在有效范围内 (0-31)
	datacenterId := int64(hash[0]) & 0x1F // 取低5位 (0-31)

	return datacenterId, nil
}

// collectMachineInfo 收集机器标识信息
func collectMachineInfo() (string, error) {
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

	return machineInfo, nil
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

// NewIdWorkerFromEnv 从环境变量创建 IdWorker 实例
// 优先使用环境变量 IDWORKER_WORKER_ID 和 IDWORKER_DATACENTER_ID
// 如果环境变量未设置，则自动生成
func NewIdWorkerFromEnv() (*IdWorker, error) {
	iw := &IdWorker{}

	// 尝试从环境变量获取 workerId
	workerId, err := GetWorkerIdFromEnv()
	if err != nil {
		// 如果环境变量未设置，自动生成
		workerId, err = GetMachineWorkerId()
		if err != nil {
			return nil, fmt.Errorf("获取 workerId 失败: %w", err)
		}
	}

	// 尝试从环境变量获取 datacenterId
	datacenterId, err := GetDatacenterIdFromEnv()
	if err != nil {
		// 如果环境变量未设置，自动生成
		datacenterId, err = GetMachineDatacenterId()
		if err != nil {
			return nil, fmt.Errorf("获取 datacenterId 失败: %w", err)
		}
	}

	err = iw.InitIdWorker(workerId, datacenterId)
	if err != nil {
		return nil, fmt.Errorf("初始化 IdWorker 失败: %w", err)
	}

	return iw, nil
}

// NewIdWorkerWithAutoId 使用自动生成的 workerId 和 datacenterId 创建 IdWorker 实例
func NewIdWorkerWithAutoId() (*IdWorker, error) {
	iw := &IdWorker{}

	workerId, err := GetMachineWorkerId()
	if err != nil {
		return nil, fmt.Errorf("自动生成 workerId 失败: %w", err)
	}

	datacenterId, err := GetMachineDatacenterId()
	if err != nil {
		return nil, fmt.Errorf("自动生成 datacenterId 失败: %w", err)
	}

	err = iw.InitIdWorker(workerId, datacenterId)
	if err != nil {
		return nil, fmt.Errorf("初始化 IdWorker 失败: %w", err)
	}

	return iw, nil
}
