package copierConv

import (
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

func TestStringToTimeConverter(t *testing.T) {
	// 创建转换器
	converter := GetStringToTimeConverter()

	// 测试用例
	testCases := []struct {
		name     string
		input    string
		expected time.Time
		hasError bool
	}{
		// 标准格式
		{
			name:     "标准格式",
			input:    "2023-12-25 15:30:45",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// RFC3339 格式
		{
			name:     "RFC3339格式",
			input:    "2023-12-25T15:30:45Z",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// RFC3339 带毫秒
		{
			name:     "RFC3339带毫秒",
			input:    "2023-12-25T15:30:45.123Z",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC),
			hasError: false,
		},
		// RFC3339 带时区
		{
			name:     "RFC3339带时区",
			input:    "2023-12-25T15:30:45+08:00",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.FixedZone("+08:00", 8*3600)),
			hasError: false,
		},
		// 日期格式
		{
			name:     "日期格式",
			input:    "2023-12-25",
			expected: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		// 时间格式
		{
			name:     "时间格式",
			input:    "15:30:45",
			expected: time.Date(0, 1, 1, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// 斜杠分隔格式
		{
			name:     "斜杠分隔格式",
			input:    "2023/12/25 15:30:45",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// 美式格式
		{
			name:     "美式格式",
			input:    "12/25/2023 15:30:45",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// 欧式格式
		{
			name:     "欧式格式",
			input:    "25/12/2023 15:30:45",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// 带毫秒的标准格式
		{
			name:     "带毫秒的标准格式",
			input:    "2023-12-25 15:30:45.123",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC),
			hasError: false,
		},
		// ISO 格式（无时区）
		{
			name:     "ISO格式无时区",
			input:    "2023-12-25T15:30:45",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			hasError: false,
		},
		// ISO 格式带毫秒（无时区）
		{
			name:     "ISO格式带毫秒无时区",
			input:    "2023-12-25T15:30:45.123",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC),
			hasError: false,
		},
		// Unix 时间戳（秒）
		{
			name:     "Unix时间戳秒",
			input:    "1703511045",
			expected: time.Unix(1703511045, 0),
			hasError: false,
		},
		// Unix 时间戳（毫秒）
		{
			name:     "Unix时间戳毫秒",
			input:    "1703511045123",
			expected: time.Unix(1703511045, 123000000),
			hasError: false,
		},
		// Unix 时间戳（浮点数秒）
		{
			name:     "Unix时间戳浮点数秒",
			input:    "1703511045.123",
			expected: time.Unix(1703511045, 123000000),
			hasError: false,
		},
		// 无效格式
		{
			name:     "无效格式",
			input:    "invalid-time",
			expected: time.Time{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.Fn(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("期望错误，但得到了结果: %v", result)
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误，但得到了: %v", err)
				return
			}

			actual, ok := result.(time.Time)
			if !ok {
				t.Errorf("期望 time.Time 类型，但得到了: %T", result)
				return
			}

			// 比较时间（对于浮点数时间戳，允许微小的纳秒差异）
			if tc.name == "Unix时间戳浮点数秒" {
				// 对于浮点数时间戳，允许1毫秒的差异
				diff := actual.Sub(tc.expected)
				if diff < -time.Millisecond || diff > time.Millisecond {
					t.Errorf("期望: %v, 实际: %v, 差异: %v", tc.expected, actual, diff)
				}
			} else {
				// 其他情况使用精确比较
				if !actual.Equal(tc.expected) {
					t.Errorf("期望: %v, 实际: %v", tc.expected, actual)
				}
			}
		})
	}
}

func TestStringToTimeConverterWithCopier(t *testing.T) {
	// 测试与 copier 库的集成
	type Source struct {
		TimeStr string
	}

	type Destination struct {
		TimeStr time.Time // 字段名相同，但类型不同
	}

	// 使用转换器选项
	option := copier.Option{
		Converters: []copier.TypeConverter{GetStringToTimeConverter()},
	}

	source := Source{
		TimeStr: "2023-12-25 15:30:45",
	}

	var dest Destination
	err := copier.CopyWithOption(&dest, &source, option)
	if err != nil {
		t.Errorf("复制失败: %v", err)
		return
	}

	expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	if !dest.TimeStr.Equal(expected) {
		t.Errorf("期望: %v, 实际: %v", expected, dest.TimeStr)
	}
}

func TestStringToTimeConverterEdgeCases(t *testing.T) {
	converter := GetStringToTimeConverter()

	// 测试边界情况
	edgeCases := []struct {
		name  string
		input string
		valid bool
	}{
		{"空字符串", "", false},
		{"只有数字", "123", true}, // Unix 时间戳
		{"只有日期", "2023-12-25", true},
		{"只有时间", "15:30:45", true},
		{"带时区的RFC3339", "2023-12-25T15:30:45+08:00", true},
		{"带毫秒和时区的RFC3339", "2023-12-25T15:30:45.123+08:00", true},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.Fn(tc.input)

			if tc.valid {
				if err != nil {
					t.Errorf("期望成功解析，但得到错误: %v", err)
				} else if result == nil {
					t.Errorf("期望得到结果，但得到 nil")
				}
			} else {
				if err == nil && result != nil {
					t.Errorf("期望解析失败，但得到结果: %v", result)
				}
			}
		})
	}
}
