package tools

import (
	"sort"
	"strconv"
	"strings"
)

func FormatIds(ids []int32) []string {
	if len(ids) == 0 {
		return []string{}
	}

	// 对ids 排序从小到大
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	var result []string
	start := ids[0]
	prev := ids[0]

	for i := 1; i < len(ids); i++ {
		if ids[i] != prev+1 {
			if start == prev {
				result = append(result, strconv.Itoa(int(start)))
			} else {
				result = append(result, strconv.Itoa(int(start))+"-"+strconv.Itoa(int(prev)))
			}
			start = ids[i]
		}
		prev = ids[i]
	}

	// 处理最后一组
	if start == prev {
		result = append(result, strconv.Itoa(int(start)))
	} else {
		result = append(result, strconv.Itoa(int(start))+"-"+strconv.Itoa(int(prev)))
	}
	return result
}

func ParseIds(ids []string) []int32 {
	if len(ids) == 0 {
		return []int32{}
	}

	var result []int32
	for _, id := range ids {
		// 检查是否包含 "-" 分隔符
		parts := strings.Split(id, "-")
		if len(parts) == 1 {
			// 单个数字的情况
			num, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			result = append(result, int32(num))
		} else if len(parts) == 2 {
			// 范围数字的情况
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				continue
			}
			// 添加范围内的所有数字
			for i := start; i <= end; i++ {
				result = append(result, int32(i))
			}
		}
	}

	return result
}

func MergeInt32Slices(slice1, slice2 []int32) []int32 {
	merged := make([]int32, 0, len(slice1)+len(slice2))
	merged = append(merged, slice1...)
	merged = append(merged, slice2...)
	return RemoveDuplicateInt32(merged)
}

func RemoveDuplicateInt32(slice []int32) []int32 {
	if len(slice) == 0 {
		return slice
	}
	// 先排序
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})

	// 去重
	j := 0
	for i := 1; i < len(slice); i++ {
		if slice[j] != slice[i] {
			j++
			slice[j] = slice[i]
		}
	}

	return slice[:j+1]
}
