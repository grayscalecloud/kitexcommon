package tools

import (
	"strconv"
)

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func StringSliceToInt32Slice(strings []string) []int32 {
	var ints []int32
	for _, str := range strings {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil
		}
		ints = append(ints, int32(num))
	}
	return ints
}
