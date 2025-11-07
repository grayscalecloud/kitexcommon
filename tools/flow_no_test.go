package tools

import (
	"regexp"
	"strings"
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

	// 验证格式：前缀 + YYYYMMDDHHMMSS + 8位数字
	matched, _ := regexp.MatchString(`^[A-Z]+[0-9]{14}[0-9]{8}$`, flowNo)
	if !matched {
		t.Errorf("Flow number format is incorrect: %s", flowNo)
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
	}
}

func BenchmarkGenerateFlowNo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		no := GenerateFlowNo("test")
		b.Logf("Generated flow number: %s", no)
	}
}
