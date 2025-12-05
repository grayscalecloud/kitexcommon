package tools

import (
	"regexp"
	"strings"
	"sync"
	"testing"
)

func TestGenerateFlowNo(t *testing.T) {
	// 测试正常情况下生成的流水号
	prefix := "sf"
	flowNo := GenerateFlowNo(prefix)

	// 验证前缀是否转为大写
	expectedPrefix := "SF"
	if !strings.HasPrefix(flowNo, expectedPrefix) {
		t.Errorf("Expected prefix %s, but got %s", expectedPrefix, flowNo[:len(expectedPrefix)])
	}

	// 验证格式：前缀 + YYYYMMDDHHMMSS(14位) + 百毫秒(2位) + 序列号(3位) + 随机数(3位)
	// 总格式：前缀 + 22位数字（与原来保持一致）
	matched, _ := regexp.MatchString(`^[A-Z]+[0-9]{14}[0-9]{2}[0-9]{3}[0-9]{3}$`, flowNo)
	if !matched {
		t.Errorf("Flow number format is incorrect: %s (expected format: PREFIX + 14位日期 + 2位百毫秒 + 3位序列号 + 3位随机数)", flowNo)
	}

	// 验证长度（与原来保持一致：前缀 + 22位）
	expectedLength := len(prefix) + 22
	if len(flowNo) != expectedLength {
		t.Errorf("Flow number length is incorrect: expected %d, got %d", expectedLength, len(flowNo))
	}

	// 测试不同的前缀
	testCases := []struct {
		input    string
		expected string
	}{
		{"Ab", "AB"},
		{"test", "TEST"},
		{"XYZ", "XYZ"},
	}

	for _, tc := range testCases {
		result := GenerateFlowNo(tc.input)
		if !strings.HasPrefix(result, tc.expected) {
			t.Errorf("Expected prefix %s for input %s, but got %s", tc.expected, tc.input, result[:len(tc.expected)])
		}
		// 验证长度（与原来保持一致：前缀 + 22位）
		expectedLen := len(tc.expected) + 22
		if len(result) != expectedLen {
			t.Errorf("Flow number length is incorrect for prefix %s: expected %d, got %d", tc.input, expectedLen, len(result))
		}
	}
}

func TestGenerateFlowNo_Uniqueness(t *testing.T) {
	// 测试生成多个流水号的唯一性
	count := 1000
	flowNos := make(map[string]bool)

	for i := 0; i < count; i++ {
		flowNo := GenerateFlowNo("TEST")
		t.Logf("\nflowNo: %s", flowNo)
		if flowNos[flowNo] {
			t.Errorf("Generated duplicate flow number: %s", flowNo)
		}
		flowNos[flowNo] = true
	}

	if len(flowNos) != count {
		t.Errorf("Expected %d unique flow numbers, but got %d", count, len(flowNos))
	}
}

func TestGenerateFlowNo_Concurrent(t *testing.T) {
	// 测试并发安全性
	goroutines := 10
	idsPerGoroutine := 100
	totalIds := goroutines * idsPerGoroutine

	var wg sync.WaitGroup
	flowNosChan := make(chan string, totalIds)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				flowNo := GenerateFlowNo("CONC")
				flowNosChan <- flowNo
			}
		}()
	}

	wg.Wait()
	close(flowNosChan)

	// 检查唯一性
	flowNos := make(map[string]bool)
	for flowNo := range flowNosChan {
		if flowNos[flowNo] {
			t.Errorf("Generated duplicate flow number in concurrent test: %s", flowNo)
		}
		flowNos[flowNo] = true
	}

	if len(flowNos) != totalIds {
		t.Errorf("Expected %d unique flow numbers, but got %d", totalIds, len(flowNos))
	}
}

func BenchmarkGenerateFlowNo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateFlowNo("test")
	}
}
