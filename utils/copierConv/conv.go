package copierConv

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
)

// GetTimeToUnixConverter 返回一个类型转换器，用于将 time.Time 转换为 Unix 时间戳（秒）
// 该转换器将 time.Time 对象转换为表示 Unix 时间戳（从1970年起的秒数）的 int64 类型
func GetTimeToUnixConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(time.Time)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return s.Unix(), nil
		},
	}
}

// GetUnixToTimeConverter 返回一个类型转换器，用于将 Unix 时间戳（秒）转换为 time.Time
// 该转换器将 int64 类型的 Unix 时间戳转换为 time.Time 对象
func GetUnixToTimeConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int64(0),
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(int64)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return time.Unix(s, 0), nil
		},
	}
}

// GetTimeToUnixMilliConverter 返回一个类型转换器，用于将 time.Time 转换为 Unix 毫秒时间戳
// 该转换器将 time.Time 对象转换为表示 Unix 毫秒时间戳的 int64 类型
func GetTimeToUnixMilliConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(time.Time)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return s.UnixMilli(), nil
		},
	}
}

// GetUnixMilliToTimeConverter 返回一个类型转换器，用于将 Unix 毫秒时间戳转换为 time.Time
// 该转换器将 int64 类型的 Unix 毫秒时间戳转换为 time.Time 对象
func GetUnixMilliToTimeConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int64(0),
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(int64)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return time.UnixMilli(s), nil
		},
	}
}

// GetTimeToStringConverter 返回一个类型转换器，用于将 time.Time 转换为格式化的日期字符串
// 该转换器将 time.Time 对象转换为格式为 "2006-01-02 15:04:05" 的字符串
func GetTimeToStringConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: "",
		Fn: func(src interface{}) (interface{}, error) {
			t, ok := src.(time.Time)
			if !ok {
				return nil, nil
			}
			return t.Format("2006-01-02 15:04:05"), nil
		},
	}
}

// GetStringToTimeConverter 返回一个类型转换器，用于将格式化的日期字符串转换为 time.Time
// 该转换器将格式为 "2006-01-02 15:04:05" 的字符串转换为 time.Time 对象
// 如果字符串格式不正确，将返回错误
func GetStringToTimeConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: "",
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, nil
			}
			return time.Parse("2006-01-02 15:04:05", s)
		},
	}
}

// GetStringToInt64Converter 返回一个类型转换器，用于将字符串转换为int64类型
// 该转换器将字符串解析为int64数值，如果解析失败则返回错误
func GetStringToInt64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: "",
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			if s == "" {
				return int64(0), nil
			}
			return strconv.ParseInt(s, 10, 64)
		},
	}
}

// GetInt64ToStringConverter 返回一个类型转换器，用于将int64类型转换为字符串
// 该转换器将int64数值转换为其字符串表示
func GetInt64ToStringConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int64(0),
		DstType: "",
		Fn: func(src interface{}) (interface{}, error) {
			i, ok := src.(int64)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return strconv.FormatInt(i, 10), nil
		},
	}
}

// GetStringToFloat64Converter 返回一个类型转换器，用于将字符串转换为float64类型
// 该转换器将字符串解析为float64数值，如果解析失败则返回错误
func GetStringToFloat64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: "",
		DstType: float64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			if s == "" {
				return float64(0), nil
			}
			return strconv.ParseFloat(s, 64)
		},
	}
}

// GetFloat64ToStringConverter 返回一个类型转换器，用于将float64类型转换为字符串
// 该转换器将float64数值转换为其字符串表示，保留小数点后的精度
func GetFloat64ToStringConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: float64(0),
		DstType: "",
		Fn: func(src interface{}) (interface{}, error) {
			f, ok := src.(float64)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return strconv.FormatFloat(f, 'f', -1, 64), nil
		},
	}
}

// GetInt64ToFloat64Converter 返回一个类型转换器，用于将int64类型转换为float64类型
func GetInt64ToFloat64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int64(0),
		DstType: float64(0),
		Fn: func(src interface{}) (interface{}, error) {
			i, ok := src.(int64)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return float64(i), nil
		},
	}
}

// GetFloat64ToInt64Converter 返回一个类型转换器，用于将float64类型转换为int64类型
// 注意：该转换会截断小数部分
func GetFloat64ToInt64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: float64(0),
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			f, ok := src.(float64)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return int64(f), nil
		},
	}
}

// GetInt32ToBoolConverter 返回一个类型转换器，用于将int32类型转换为bool类型
// 转换规则：非零值转换为true，零值转换为false
func GetInt32ToBoolConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int32(0),
		DstType: false,
		Fn: func(src interface{}) (interface{}, error) {
			i, ok := src.(int32)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return i != 0, nil
		},
	}
}

// GetBoolToInt32Converter 返回一个类型转换器，用于将bool类型转换为int32类型
// 转换规则：true转换为1，false转换为0
func GetBoolToInt32Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: false,
		DstType: int32(0),
		Fn: func(src interface{}) (interface{}, error) {
			b, ok := src.(bool)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			if b {
				return int32(1), nil
			}
			return int32(0), nil
		},
	}
}

// GetInt8ToBoolConverter 返回一个类型转换器，用于将int8类型转换为bool类型
// 转换规则：非零值转换为true，零值转换为false
func GetInt8ToBoolConverter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: int8(0),
		DstType: false,
		Fn: func(src interface{}) (interface{}, error) {
			i, ok := src.(int8)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			return i != 0, nil
		},
	}
}

// GetBoolToInt32Converter 返回一个类型转换器，用于将bool类型转换为int8类型
// 转换规则：true转换为1，false转换为0
func GetBoolToInt8Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: false,
		DstType: int8(0),
		Fn: func(src interface{}) (interface{}, error) {
			b, ok := src.(bool)
			if !ok {
				return nil, errors.New("源类型不匹配")
			}
			if b {
				return int8(1), nil
			}
			return int8(0), nil
		},
	}
}
