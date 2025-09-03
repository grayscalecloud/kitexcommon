package tools

import (
	"strconv"
	"strings"
)

// BoolFormat 三目运算
func BoolFormat(expr bool, a, b interface{}) interface{} {
	if expr {
		return a
	}
	return b
}

func BoolFormatInt32(expr bool, a, b int32) int32 {
	if expr {
		return a
	}
	return b
}

func BoolFormatStr(expr bool, a, b string) string {
	if expr {
		return a
	}
	return b
}

func Int64SliceToString(goldIds []int64) string {
	strSlice := make([]string, len(goldIds))
	for i, id := range goldIds {
		strSlice[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(strSlice, ",")
}
