package idworker

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestInitIdWorker_Success(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 应该成功，但返回错误: %v", err)
	}

	if iw.workerId != 1 {
		t.Errorf("workerId 应该是 1，实际是 %d", iw.workerId)
	}
	if iw.datacenterId != 1 {
		t.Errorf("datacenterId 应该是 1，实际是 %d", iw.datacenterId)
	}
	if iw.maxWorkerId != 31 {
		t.Errorf("maxWorkerId 应该是 31，实际是 %d", iw.maxWorkerId)
	}
	if iw.maxDatacenterId != 31 {
		t.Errorf("maxDatacenterId 应该是 31，实际是 %d", iw.maxDatacenterId)
	}
}

func TestInitIdWorker_WorkerIdOutOfRange(t *testing.T) {
	tests := []struct {
		name     string
		workerId int64
	}{
		{"负数", -1},
		{"超出最大值", 32},
		{"超出最大值2", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IdWorker{}
			err := iw.InitIdWorker(tt.workerId, 1)
			if err == nil {
				t.Errorf("InitIdWorker 应该返回错误，但没有返回")
			}
		})
	}
}

func TestInitIdWorker_DatacenterIdOutOfRange(t *testing.T) {
	tests := []struct {
		name         string
		datacenterId int64
	}{
		{"负数", -1},
		{"超出最大值", 32},
		{"超出最大值2", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iw := &IdWorker{}
			err := iw.InitIdWorker(1, tt.datacenterId)
			if err == nil {
				t.Errorf("InitIdWorker 应该返回错误，但没有返回")
			}
		})
	}
}

func TestInitIdWorker_BoundaryValues(t *testing.T) {
	// 测试边界值：0 和最大值
	iw1 := &IdWorker{}
	err := iw1.InitIdWorker(0, 0)
	if err != nil {
		t.Errorf("InitIdWorker(0, 0) 应该成功，但返回错误: %v", err)
	}

	iw2 := &IdWorker{}
	err = iw2.InitIdWorker(31, 31)
	if err != nil {
		t.Errorf("InitIdWorker(31, 31) 应该成功，但返回错误: %v", err)
	}
}

func TestNextId_Basic(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	id, err := iw.NextId()
	if err != nil {
		t.Fatalf("NextId 应该成功，但返回错误: %v", err)
	}

	if id <= 0 {
		t.Errorf("生成的 ID 应该大于 0，实际是 %d", id)
	}
}

func TestNextId_Uniqueness(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	count := 1000
	ids := make(map[int64]bool)

	for i := 0; i < count; i++ {
		id, err := iw.NextId()
		if err != nil {
			t.Fatalf("NextId 失败: %v", err)
		}

		if ids[id] {
			t.Errorf("生成的 ID 重复: %d", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("应该生成 %d 个唯一 ID，实际生成了 %d 个", count, len(ids))
	}
}

func TestNextId_Monotonicity(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	count := 100
	prevId := int64(0)

	for i := 0; i < count; i++ {
		id, err := iw.NextId()
		if err != nil {
			t.Fatalf("NextId 失败: %v", err)
		}

		if id <= prevId {
			t.Errorf("ID 应该单调递增，但 %d <= %d", id, prevId)
		}
		prevId = id
	}
}

func TestNextId_Concurrent(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	goroutines := 10
	idsPerGoroutine := 100
	totalIds := goroutines * idsPerGoroutine

	var wg sync.WaitGroup
	idsChan := make(chan int64, totalIds)
	errorsChan := make(chan error, totalIds)

	// 启动多个 goroutine 并发生成 ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := iw.NextId()
				if err != nil {
					errorsChan <- err
					return
				}
				idsChan <- id
			}
		}()
	}

	wg.Wait()
	close(idsChan)
	close(errorsChan)

	// 检查错误
	for err := range errorsChan {
		t.Fatalf("并发生成 ID 时出错: %v", err)
	}

	// 检查唯一性
	ids := make(map[int64]bool)
	for id := range idsChan {
		if ids[id] {
			t.Errorf("并发生成的 ID 重复: %d", id)
		}
		ids[id] = true
	}

	if len(ids) != totalIds {
		t.Errorf("应该生成 %d 个唯一 ID，实际生成了 %d 个", totalIds, len(ids))
	}
}

func TestNextId_DifferentWorkers(t *testing.T) {
	// 测试不同 workerId 和 datacenterId 生成的 ID 不同
	iw1 := &IdWorker{}
	err := iw1.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	iw2 := &IdWorker{}
	err = iw2.InitIdWorker(2, 2)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	// 在同一毫秒内生成 ID，应该因为 workerId 和 datacenterId 不同而不同
	id1, err := iw1.NextId()
	if err != nil {
		t.Fatalf("NextId 失败: %v", err)
	}

	// 稍微延迟以确保时间戳可能相同
	time.Sleep(time.Millisecond)

	id2, err := iw2.NextId()
	if err != nil {
		t.Fatalf("NextId 失败: %v", err)
	}

	if id1 == id2 {
		t.Errorf("不同 worker 生成的 ID 应该不同，但都是 %d", id1)
	}
}

func TestNextId_RapidGeneration(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	// 快速生成大量 ID，测试序列号机制
	count := 5000
	ids := make([]int64, count)

	start := time.Now()
	for i := 0; i < count; i++ {
		id, err := iw.NextId()
		if err != nil {
			t.Fatalf("NextId 失败: %v", err)
		}
		ids[i] = id
	}
	duration := time.Since(start)

	// 检查唯一性
	idMap := make(map[int64]bool)
	for _, id := range ids {
		if idMap[id] {
			t.Errorf("快速生成时出现重复 ID: %d", id)
		}
		idMap[id] = true
	}

	t.Logf("生成 %d 个 ID 耗时: %v", count, duration)
}

func TestTilNextMillis(t *testing.T) {
	iw := &IdWorker{}
	err := iw.InitIdWorker(1, 1)
	if err != nil {
		t.Fatalf("InitIdWorker 失败: %v", err)
	}

	// 设置一个未来的时间戳
	iw.lastTimestamp = time.Now().UnixNano()/int64(time.Millisecond) + 100

	start := time.Now()
	nextTimestamp := iw.tilNextMillis()
	duration := time.Since(start)

	if nextTimestamp <= iw.lastTimestamp {
		t.Errorf("tilNextMillis 应该返回大于 lastTimestamp 的值，但 %d <= %d", nextTimestamp, iw.lastTimestamp)
	}

	// 应该至少等待了接近 100 毫秒
	if duration < 90*time.Millisecond {
		t.Logf("警告: tilNextMillis 等待时间可能不足，实际等待 %v", duration)
	}
}

func TestGetWorkerIdFromEnv(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("IDWORKER_WORKER_ID")
	defer os.Setenv("IDWORKER_WORKER_ID", originalValue)

	// 测试有效值
	os.Setenv("IDWORKER_WORKER_ID", "15")
	workerId, err := GetWorkerIdFromEnv()
	if err != nil {
		t.Errorf("GetWorkerIdFromEnv 应该成功，但返回错误: %v", err)
	}
	if workerId != 15 {
		t.Errorf("workerId 应该是 15，实际是 %d", workerId)
	}

	// 测试无效值（超出范围）
	os.Setenv("IDWORKER_WORKER_ID", "100")
	_, err = GetWorkerIdFromEnv()
	if err == nil {
		t.Errorf("GetWorkerIdFromEnv 应该返回错误，但没有返回")
	}

	// 测试未设置环境变量
	os.Unsetenv("IDWORKER_WORKER_ID")
	_, err = GetWorkerIdFromEnv()
	if err == nil {
		t.Errorf("GetWorkerIdFromEnv 应该返回错误，但没有返回")
	}
}

func TestGetDatacenterIdFromEnv(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("IDWORKER_DATACENTER_ID")
	defer os.Setenv("IDWORKER_DATACENTER_ID", originalValue)

	// 测试有效值
	os.Setenv("IDWORKER_DATACENTER_ID", "10")
	datacenterId, err := GetDatacenterIdFromEnv()
	if err != nil {
		t.Errorf("GetDatacenterIdFromEnv 应该成功，但返回错误: %v", err)
	}
	if datacenterId != 10 {
		t.Errorf("datacenterId 应该是 10，实际是 %d", datacenterId)
	}

	// 测试未设置环境变量
	os.Unsetenv("IDWORKER_DATACENTER_ID")
	_, err = GetDatacenterIdFromEnv()
	if err == nil {
		t.Errorf("GetDatacenterIdFromEnv 应该返回错误，但没有返回")
	}
}

func TestGetMachineWorkerId(t *testing.T) {
	workerId, err := GetMachineWorkerId()
	if err != nil {
		t.Fatalf("GetMachineWorkerId 应该成功，但返回错误: %v", err)
	}

	if workerId < 0 || workerId > 31 {
		t.Errorf("workerId 应该在 [0-31] 范围内，实际是 %d", workerId)
	}

	// 测试多次调用应该返回相同的值（基于机器信息）
	workerId2, err := GetMachineWorkerId()
	if err != nil {
		t.Fatalf("GetMachineWorkerId 第二次调用失败: %v", err)
	}
	if workerId != workerId2 {
		t.Errorf("同一机器应该生成相同的 workerId，但 %d != %d", workerId, workerId2)
	}
}

func TestGetMachineDatacenterId(t *testing.T) {
	datacenterId, err := GetMachineDatacenterId()
	if err != nil {
		t.Fatalf("GetMachineDatacenterId 应该成功，但返回错误: %v", err)
	}

	if datacenterId < 0 || datacenterId > 31 {
		t.Errorf("datacenterId 应该在 [0-31] 范围内，实际是 %d", datacenterId)
	}

	// 测试多次调用应该返回相同的值
	datacenterId2, err := GetMachineDatacenterId()
	if err != nil {
		t.Fatalf("GetMachineDatacenterId 第二次调用失败: %v", err)
	}
	if datacenterId != datacenterId2 {
		t.Errorf("同一机器应该生成相同的 datacenterId，但 %d != %d", datacenterId, datacenterId2)
	}
}

func TestNewIdWorkerFromEnv(t *testing.T) {
	// 保存原始环境变量
	originalWorkerId := os.Getenv("IDWORKER_WORKER_ID")
	originalDatacenterId := os.Getenv("IDWORKER_DATACENTER_ID")
	defer func() {
		if originalWorkerId != "" {
			os.Setenv("IDWORKER_WORKER_ID", originalWorkerId)
		} else {
			os.Unsetenv("IDWORKER_WORKER_ID")
		}
		if originalDatacenterId != "" {
			os.Setenv("IDWORKER_DATACENTER_ID", originalDatacenterId)
		} else {
			os.Unsetenv("IDWORKER_DATACENTER_ID")
		}
	}()

	// 测试从环境变量创建
	os.Setenv("IDWORKER_WORKER_ID", "5")
	os.Setenv("IDWORKER_DATACENTER_ID", "10")
	iw, err := NewIdWorkerFromEnv()
	if err != nil {
		t.Fatalf("NewIdWorkerFromEnv 应该成功，但返回错误: %v", err)
	}
	if iw.workerId != 5 {
		t.Errorf("workerId 应该是 5，实际是 %d", iw.workerId)
	}
	if iw.datacenterId != 10 {
		t.Errorf("datacenterId 应该是 10，实际是 %d", iw.datacenterId)
	}

	// 测试环境变量未设置时自动生成
	os.Unsetenv("IDWORKER_WORKER_ID")
	os.Unsetenv("IDWORKER_DATACENTER_ID")
	iw2, err := NewIdWorkerFromEnv()
	if err != nil {
		t.Fatalf("NewIdWorkerFromEnv 应该自动生成 ID，但返回错误: %v", err)
	}
	if iw2.workerId < 0 || iw2.workerId > 31 {
		t.Errorf("自动生成的 workerId 应该在 [0-31] 范围内，实际是 %d", iw2.workerId)
	}
	if iw2.datacenterId < 0 || iw2.datacenterId > 31 {
		t.Errorf("自动生成的 datacenterId 应该在 [0-31] 范围内，实际是 %d", iw2.datacenterId)
	}

	// 测试生成的 ID 可以正常工作
	id, err := iw2.NextId()
	if err != nil {
		t.Fatalf("NextId 应该成功，但返回错误: %v", err)
	}
	if id <= 0 {
		t.Errorf("生成的 ID 应该大于 0，实际是 %d", id)
	}
}

func TestNewIdWorkerWithAutoId(t *testing.T) {
	iw, err := NewIdWorkerWithAutoId()
	if err != nil {
		t.Fatalf("NewIdWorkerWithAutoId 应该成功，但返回错误: %v", err)
	}

	if iw.workerId < 0 || iw.workerId > 31 {
		t.Errorf("自动生成的 workerId 应该在 [0-31] 范围内，实际是 %d", iw.workerId)
	}
	if iw.datacenterId < 0 || iw.datacenterId > 31 {
		t.Errorf("自动生成的 datacenterId 应该在 [0-31] 范围内，实际是 %d", iw.datacenterId)
	}

	// 测试生成的 ID 可以正常工作
	id, err := iw.NextId()
	if err != nil {
		t.Fatalf("NextId 应该成功，但返回错误: %v", err)
	}
	if id <= 0 {
		t.Errorf("生成的 ID 应该大于 0，实际是 %d", id)
	}

	// 测试同一机器多次调用应该生成相同的 workerId 和 datacenterId
	iw2, err := NewIdWorkerWithAutoId()
	if err != nil {
		t.Fatalf("NewIdWorkerWithAutoId 第二次调用失败: %v", err)
	}
	if iw.workerId != iw2.workerId {
		t.Errorf("同一机器应该生成相同的 workerId，但 %d != %d", iw.workerId, iw2.workerId)
	}
	if iw.datacenterId != iw2.datacenterId {
		t.Errorf("同一机器应该生成相同的 datacenterId，但 %d != %d", iw.datacenterId, iw2.datacenterId)
	}
}
