package copierConv

import (
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

func TestTimeToStringConverter(t *testing.T) {
	// 创建测试时间
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)

	// 测试默认转换器
	converter := GetTimeToStringConverter()
	result, err := converter.Fn(testTime)
	if err != nil {
		t.Errorf("转换失败: %v", err)
		return
	}

	expected := "2023-12-25 15:30:45"
	if result != expected {
		t.Errorf("期望: %s, 实际: %s", expected, result)
	}
}

func TestTimeToStringConverterWithFormat(t *testing.T) {
	// 创建测试时间
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)

	// 测试各种格式
	testCases := []struct {
		name     string
		format   string
		expected string
	}{
		{"标准格式", "2006-01-02 15:04:05", "2023-12-25 15:30:45"},
		{"RFC3339格式", "2006-01-02T15:04:05Z", "2023-12-25T15:30:45Z"},
		{"RFC3339带毫秒", "2006-01-02T15:04:05.000Z", "2023-12-25T15:30:45.123Z"},
		{"RFC3339带时区", "2006-01-02T15:04:05Z07:00", "2023-12-25T15:30:45Z"},
		{"日期格式", "2006-01-02", "2023-12-25"},
		{"时间格式", "15:04:05", "15:30:45"},
		{"斜杠分隔格式", "2006/01/02 15:04:05", "2023/12/25 15:30:45"},
		{"美式格式", "01/02/2006 15:04:05", "12/25/2023 15:30:45"},
		{"欧式格式", "02/01/2006 15:04:05", "25/12/2023 15:30:45"},
		{"Unix时间戳秒", "unix", "1703518245"},
		{"Unix时间戳毫秒", "unixms", "1703518245123"},
		{"Unix时间戳纳秒", "unixns", "1703518245123000000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var converter copier.TypeConverter

			// 根据格式类型选择不同的转换器
			if tc.format == "unix" || tc.format == "unixms" || tc.format == "unixns" {
				converter = GetTimeToUnixStringConverter(tc.format)
			} else {
				converter = GetTimeToStringConverterWithFormat(tc.format)
			}

			result, err := converter.Fn(testTime)
			if err != nil {
				t.Errorf("转换失败: %v", err)
				return
			}

			if result != tc.expected {
				t.Errorf("期望: %s, 实际: %s", tc.expected, result)
			}
		})
	}
}

func TestTimeToStringConverterWithCustomFormat(t *testing.T) {
	// 创建测试时间
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)

	// 测试自定义格式
	converter := GetTimeToStringConverterWithFormat("2006年01月02日 15:04:05")
	result, err := converter.Fn(testTime)
	if err != nil {
		t.Errorf("转换失败: %v", err)
		return
	}

	expected := "2023年12月25日 15:30:45"
	if result != expected {
		t.Errorf("期望: %s, 实际: %s", expected, result)
	}
}

func TestTimeToStringConverterWithCopier(t *testing.T) {
	// 定义源结构体
	type Source struct {
		CreateTime time.Time `json:"create_time"`
		UpdateTime time.Time `json:"update_time"`
		LoginTime  time.Time `json:"login_time"`
		Birthday   time.Time `json:"birthday"`
		LastActive time.Time `json:"last_active"`
	}

	// 定义目标结构体
	type Destination struct {
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		LoginTime  string `json:"login_time"`
		Birthday   string `json:"birthday"`
		LastActive string `json:"last_active"`
	}

	// 创建源数据
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)
	source := Source{
		CreateTime: testTime,
		UpdateTime: testTime,
		LoginTime:  testTime,
		Birthday:   testTime,
		LastActive: testTime,
	}

	// 使用标准格式转换器
	option := copier.Option{
		Converters: []copier.TypeConverter{GetTimeToStringConverter()},
	}

	var dest Destination
	err := copier.CopyWithOption(&dest, &source, option)
	if err != nil {
		t.Errorf("复制失败: %v", err)
		return
	}

	expected := "2023-12-25 15:30:45"
	if dest.CreateTime != expected {
		t.Errorf("期望: %s, 实际: %s", expected, dest.CreateTime)
	}
}

func TestTimeToStringConverterWithDifferentFormats(t *testing.T) {
	// 定义源结构体
	type Source struct {
		CreateTime time.Time `json:"create_time"`
		UpdateTime time.Time `json:"update_time"`
		LoginTime  time.Time `json:"login_time"`
		Birthday   time.Time `json:"birthday"`
		LastActive time.Time `json:"last_active"`
	}

	// 定义目标结构体
	type Destination struct {
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		LoginTime  string `json:"login_time"`
		Birthday   string `json:"birthday"`
		LastActive string `json:"last_active"`
	}

	// 创建源数据
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)
	source := Source{
		CreateTime: testTime,
		UpdateTime: testTime,
		LoginTime:  testTime,
		Birthday:   testTime,
		LastActive: testTime,
	}

	// 使用标准格式转换器（copier 库会使用第一个匹配的转换器）
	option := copier.Option{
		Converters: []copier.TypeConverter{GetTimeToStringConverter()},
	}

	var dest Destination
	err := copier.CopyWithOption(&dest, &source, option)
	if err != nil {
		t.Errorf("复制失败: %v", err)
		return
	}

	// 验证标准格式
	expected := "2023-12-25 15:30:45"
	if dest.CreateTime != expected {
		t.Errorf("CreateTime 期望: %s, 实际: %s", expected, dest.CreateTime)
	}
	if dest.UpdateTime != expected {
		t.Errorf("UpdateTime 期望: %s, 实际: %s", expected, dest.UpdateTime)
	}
	if dest.LoginTime != expected {
		t.Errorf("LoginTime 期望: %s, 实际: %s", expected, dest.LoginTime)
	}
	if dest.Birthday != expected {
		t.Errorf("Birthday 期望: %s, 实际: %s", expected, dest.Birthday)
	}
	if dest.LastActive != expected {
		t.Errorf("LastActive 期望: %s, 实际: %s", expected, dest.LastActive)
	}
}

func TestTimeToStringConverterEdgeCases(t *testing.T) {
	// 测试边界情况
	edgeCases := []struct {
		name     string
		time     time.Time
		format   string
		expected string
	}{
		{
			name:     "零时间",
			time:     time.Time{},
			format:   "2006-01-02 15:04:05",
			expected: "0001-01-01 00:00:00",
		},
		{
			name:     "Unix时间戳零时间",
			time:     time.Time{},
			format:   "unix",
			expected: "-62135596800",
		},
		{
			name:     "最大时间",
			time:     time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			format:   "2006-01-02 15:04:05",
			expected: "9999-12-31 23:59:59",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			var converter copier.TypeConverter

			// 根据格式类型选择不同的转换器
			if tc.format == "unix" || tc.format == "unixms" || tc.format == "unixns" {
				converter = GetTimeToUnixStringConverter(tc.format)
			} else {
				converter = GetTimeToStringConverterWithFormat(tc.format)
			}

			result, err := converter.Fn(tc.time)
			if err != nil {
				t.Errorf("转换失败: %v", err)
				return
			}

			if result != tc.expected {
				t.Errorf("期望: %s, 实际: %s", tc.expected, result)
			}
		})
	}
}
