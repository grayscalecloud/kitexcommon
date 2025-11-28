package utils

// 预定义的位置换表（0-63 的随机排列）
var permutation = [64]uint{
	7, 22, 13, 8, 30, 24, 17, 2,
	28, 19, 11, 29, 5, 20, 15, 31,
	0, 12, 25, 21, 4, 10, 16, 1,
	27, 23, 6, 14, 9, 3, 26, 18,
	45, 58, 41, 50, 62, 56, 49, 34,
	60, 51, 43, 61, 37, 52, 47, 63,
	32, 44, 57, 53, 36, 42, 48, 33,
	59, 55, 38, 46, 40, 35, 54, 39,
}

// InversePermutation 逆置换表（导出供测试使用）
var InversePermutation [64]uint

func init() {
	// 生成逆置换表
	for i, p := range permutation {
		InversePermutation[p] = uint(i)
	}
}

type IntObfuscator struct {
	xorKey int64
}

// NewObfuscator 创建混淆器，建议传入随机种子作为 xorKey
// 如果不提供 xorKey，使用默认值 0x5a5a5a5a5a5a5a5a
func NewObfuscator(xorKey int64) *IntObfuscator {
	if xorKey == 0 {
		xorKey = 0x5a5a5a5a5a5a5a5a // 默认异或密钥
	}
	return &IntObfuscator{xorKey: xorKey}
}

// permuteBitsFast 快速位置换实现（内联优化）
// 使用位操作技巧，减少分支和循环开销
// 编译器会自动展开循环并优化
func permuteBitsFast(value uint64, perm *[64]uint) uint64 {
	var result uint64
	// 使用位操作技巧，避免条件判断和分支预测失败
	// 直接计算每一位的映射，编译器可以更好地优化
	// 使用指针访问数组，减少数组拷贝
	p := perm
	for i := uint(0); i < 64; i++ {
		result |= ((value >> i) & 1) << p[i]
	}
	return result
}

// PermuteBits 根据置换表重排位（导出供测试使用）
func PermuteBits(value uint64, perm [64]uint) uint64 {
	return permuteBitsFast(value, &perm)
}

// Obfuscate 混淆整型
func (o *IntObfuscator) Obfuscate(id int64) int64 {
	// 先异或，再置换
	t := uint64(id ^ o.xorKey)
	return int64(permuteBitsFast(t, &permutation))
}

// Deobfuscate 还原整型
func (o *IntObfuscator) Deobfuscate(code int64) int64 {
	// 先逆置换，再异或
	t := permuteBitsFast(uint64(code), &InversePermutation)
	return int64(t) ^ o.xorKey
}
