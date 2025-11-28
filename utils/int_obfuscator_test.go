package utils

import (
	"fmt"
	"testing"
	"time"
)

// TestIntObfuscator_ShowExamples 展示混淆值和还原值的示例
func TestIntObfuscator_ShowExamples(t *testing.T) {
	// 创建混淆器
	obfuscator := NewObfuscator(0x12345678)

	// 测试用例
	testValues := []int64{
		1,
		100,
		1000,
		123456789,
		-1,
		-100,
		-123456789,
		1,
		2,
		3,
		10,
		20,
		30,
		40,
		50,
		0,
		9223372036854775807,  // int64 最大值
		-9223372036854775808, // int64 最小值
	}

	fmt.Println("\n=== 混淆和还原示例 ===")
	fmt.Printf("%-25s | %-25s | %-25s\n", "原值", "混淆值", "还原值")
	fmt.Println("-------------------------------------------------------------------------------------------")

	for _, original := range testValues {
		obfuscated := obfuscator.Obfuscate(original)
		deobfuscated := obfuscator.Deobfuscate(obfuscated)

		fmt.Printf("%-25d | %-25d | %-25d\n", original, obfuscated, deobfuscated)

		// 验证还原正确性
		if deobfuscated != original {
			t.Errorf("还原失败: 原值=%d, 混淆值=%d, 还原值=%d", original, obfuscated, deobfuscated)
		}
	}

	fmt.Println("\n=== 不同密钥的混淆效果 ===")
	original := int64(123456789)
	obfuscator1 := NewObfuscator(0x1111111111111111)
	obfuscator2 := NewObfuscator(0x2222222222222222)
	obfuscator3 := NewObfuscator(0xABCDEF01)

	obf1 := obfuscator1.Obfuscate(original)
	obf2 := obfuscator2.Obfuscate(original)
	obf3 := obfuscator3.Obfuscate(original)

	fmt.Printf("原值: %d\n", original)
	fmt.Printf("密钥 0x1111111111111111 -> 混淆值: %d -> 还原值: %d\n", obf1, obfuscator1.Deobfuscate(obf1))
	fmt.Printf("密钥 0x2222222222222222 -> 混淆值: %d -> 还原值: %d\n", obf2, obfuscator2.Deobfuscate(obf2))
	fmt.Printf("密钥 0xABCDEF01 -> 混淆值: %d -> 还原值: %d\n", obf3, obfuscator3.Deobfuscate(obf3))
}

func TestIntObfuscator_ObfuscateAndDeobfuscate(t *testing.T) {
	// 测试用例：不同的 xorKey 值
	testCases := []struct {
		name   string
		xorKey int64
		values []int64
	}{
		{
			name:   "默认配置",
			xorKey: 0x1234567812345678,
			values: []int64{1, 100, 1000, 10000, 123456789, -1, -100, -123456789},
		},
		{
			name:   "大密钥值",
			xorKey: 0xABCDEF01,
			values: []int64{0, 1, -1, 9223372036854775807, -9223372036854775808},
		},
		{
			name:   "全1密钥（低32位）",
			xorKey: 0xFFFFFFFF,
			values: []int64{42, 999, -42, -999},
		},
		{
			name:   "零值测试",
			xorKey: 0x1234567812345678,
			values: []int64{0},
		},
		{
			name:   "默认密钥",
			xorKey: 0, // 会使用默认密钥
			values: []int64{12345, -12345},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obfuscator := NewObfuscator(tc.xorKey)

			for _, originalValue := range tc.values {
				// 混淆
				obfuscated := obfuscator.Obfuscate(originalValue)

				// 验证混淆后的值与原值不同（除非是特殊情况）
				if obfuscated == originalValue && tc.xorKey != 0 {
					t.Errorf("混淆后的值应该与原值不同: 原值=%d, 混淆值=%d", originalValue, obfuscated)
				}

				// 还原
				deobfuscated := obfuscator.Deobfuscate(obfuscated)

				// 验证还原后的值与原值相同
				if deobfuscated != originalValue {
					t.Errorf("还原失败: 原值=%d, 混淆值=%d, 还原值=%d", originalValue, obfuscated, deobfuscated)
				}
			}
		})
	}
}

func TestIntObfuscator_DefaultKey(t *testing.T) {
	// 测试默认密钥（当 xorKey 为 0 时）
	t.Run("默认密钥", func(t *testing.T) {
		obfuscator := NewObfuscator(0)
		if obfuscator.xorKey != 0x5a5a5a5a5a5a5a5a {
			t.Errorf("期望默认密钥=0x5a5a5a5a5a5a5a5a, 实际密钥=0x%x", obfuscator.xorKey)
		}

		// 验证功能正常
		original := int64(12345)
		obfuscated := obfuscator.Obfuscate(original)
		deobfuscated := obfuscator.Deobfuscate(obfuscated)
		if deobfuscated != original {
			t.Errorf("默认密钥功能异常: 原值=%d, 还原值=%d", original, deobfuscated)
		}
	})

	t.Run("自定义密钥", func(t *testing.T) {
		customKey := int64(0x1234567812345678)
		obfuscator := NewObfuscator(customKey)
		if obfuscator.xorKey != customKey {
			t.Errorf("期望密钥=0x%x, 实际密钥=0x%x", customKey, obfuscator.xorKey)
		}

		// 验证功能正常
		original := int64(12345)
		obfuscated := obfuscator.Obfuscate(original)
		deobfuscated := obfuscator.Deobfuscate(obfuscated)
		if deobfuscated != original {
			t.Errorf("自定义密钥功能异常: 原值=%d, 还原值=%d", original, deobfuscated)
		}
	})
}

func TestIntObfuscator_DifferentKeys(t *testing.T) {
	// 测试不同的 xorKey 产生不同的混淆结果
	original := int64(123456789)

	obfuscator1 := NewObfuscator(0x1111111111111111)
	obfuscator2 := NewObfuscator(0x2222222222222222)

	obfuscated1 := obfuscator1.Obfuscate(original)
	obfuscated2 := obfuscator2.Obfuscate(original)

	if obfuscated1 == obfuscated2 {
		t.Errorf("不同的 xorKey 应该产生不同的混淆结果")
	}

	// 验证各自能正确还原
	if obfuscator1.Deobfuscate(obfuscated1) != original {
		t.Errorf("obfuscator1 还原失败")
	}
	if obfuscator2.Deobfuscate(obfuscated2) != original {
		t.Errorf("obfuscator2 还原失败")
	}
}

func TestIntObfuscator_EdgeCases(t *testing.T) {
	obfuscator := NewObfuscator(0x1234567812345678)

	edgeCases := []int64{
		0,
		1,
		-1,
		9223372036854775807,  // int64 最大值
		-9223372036854775808, // int64 最小值
	}

	for _, value := range edgeCases {
		t.Run("边界值测试", func(t *testing.T) {
			obfuscated := obfuscator.Obfuscate(value)
			deobfuscated := obfuscator.Deobfuscate(obfuscated)

			if deobfuscated != value {
				t.Errorf("边界值还原失败: 原值=%d, 还原值=%d", value, deobfuscated)
			}
		})
	}
}

func BenchmarkIntObfuscator_Obfuscate(b *testing.B) {
	obfuscator := NewObfuscator(0x1234567812345678)
	value := int64(123456789)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = obfuscator.Obfuscate(value)
	}
}

func BenchmarkIntObfuscator_Deobfuscate(b *testing.B) {
	obfuscator := NewObfuscator(0x1234567812345678)
	value := int64(123456789)
	obfuscated := obfuscator.Obfuscate(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = obfuscator.Deobfuscate(obfuscated)
	}
}

// TestIntObfuscator_SingleOperationTime 测试单个操作的执行时间
func TestIntObfuscator_SingleOperationTime(t *testing.T) {
	obfuscator := NewObfuscator(0x1234567812345678)
	value := int64(123456789)

	// 预热
	for i := 0; i < 1000; i++ {
		_ = obfuscator.Obfuscate(value)
	}

	// 测试混淆操作
	iterations := 10000000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = obfuscator.Obfuscate(value)
	}
	obfuscateDuration := time.Since(start)
	obfuscateAvg := float64(obfuscateDuration.Nanoseconds()) / float64(iterations)

	// 测试还原操作
	obfuscated := obfuscator.Obfuscate(value)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = obfuscator.Deobfuscate(obfuscated)
	}
	deobfuscateDuration := time.Since(start)
	deobfuscateAvg := float64(deobfuscateDuration.Nanoseconds()) / float64(iterations)

	fmt.Printf("\n=== 单个操作执行时间 ===\n")
	fmt.Printf("测试迭代次数: %d\n", iterations)
	fmt.Printf("\n混淆操作 (Obfuscate):\n")
	fmt.Printf("  总耗时: %v\n", obfuscateDuration)
	fmt.Printf("  平均每次: %.2f ns (%.4f μs)\n", obfuscateAvg, obfuscateAvg/1000)
	fmt.Printf("  每秒可执行: %.0f 次\n", float64(iterations)/obfuscateDuration.Seconds())
	fmt.Printf("\n还原操作 (Deobfuscate):\n")
	fmt.Printf("  总耗时: %v\n", deobfuscateDuration)
	fmt.Printf("  平均每次: %.2f ns (%.4f μs)\n", deobfuscateAvg, deobfuscateAvg/1000)
	fmt.Printf("  每秒可执行: %.0f 次\n", float64(iterations)/deobfuscateDuration.Seconds())
	fmt.Printf("\n性能总结:\n")
	fmt.Printf("  - 混淆操作: %.2f ns/op\n", obfuscateAvg)
	fmt.Printf("  - 还原操作: %.2f ns/op\n", deobfuscateAvg)
	fmt.Printf("  - 两者都非常快，适合高频调用\n")
}

// TestIntObfuscator_SecurityAnalysis 分析算法的安全性
// 演示在什么情况下可以被破解
func TestIntObfuscator_SecurityAnalysis(t *testing.T) {
	fmt.Println("\n=== 算法安全性分析 ===")

	// 目标：破解这个混淆器
	unknownXorKey := int64(0x1234567812345678)
	targetObfuscator := NewObfuscator(unknownXorKey)

	// 场景1: 如果攻击者有一个已知的明文-密文对
	fmt.Println("\n【场景1】已知明文-密文对破解")
	knownPlaintext := int64(123456789)
	knownCiphertext := targetObfuscator.Obfuscate(knownPlaintext)
	fmt.Printf("已知: 明文=%d, 密文=%d\n", knownPlaintext, knownCiphertext)
	fmt.Println("注意: 置换表是固定的（在代码中可见）")

	// 由于置换表是固定的，攻击者可以反推 xorKey
	fmt.Println("\n尝试破解 xorKey...")
	// 从密文反推：先逆置换，然后与明文异或得到 xorKey
	reversed := PermuteBits(uint64(knownCiphertext), InversePermutation)
	crackedXorKey := int64(reversed) ^ knownPlaintext

	// 验证这个候选密钥是否正确
	testObfuscator := NewObfuscator(crackedXorKey)
	if testObfuscator.Obfuscate(knownPlaintext) == knownCiphertext {
		// 用另一个值验证
		testValue := int64(999999)
		if testObfuscator.Obfuscate(testValue) == targetObfuscator.Obfuscate(testValue) {
			fmt.Printf("✓ 破解成功! xorKey=0x%X\n", crackedXorKey)
			fmt.Printf("警告: 算法已被破解! 仅需 1 个明文-密文对即可破解\n")
			fmt.Println("原因: 置换表是固定的，攻击者可以从代码中获取")
		}
	}

	// 场景2: 如果攻击者只有密文，没有明文
	fmt.Println("\n【场景2】仅密文攻击（暴力破解）")
	unknownCiphertext := targetObfuscator.Obfuscate(int64(999888777))
	fmt.Printf("只有密文: %d\n", unknownCiphertext)
	fmt.Println("需要尝试:")
	fmt.Printf("  - xorKey 值: 2^64 种可能 (int64 范围)\n")
	fmt.Printf("  - 总搜索空间: 2^64 ≈ 1.84 * 10^19 种组合\n")
	fmt.Println("结论: 仅密文攻击在计算上不可行（需要暴力搜索整个密钥空间）")

	// 场景3: 如果有多个明文-密文对
	fmt.Println("\n【场景3】多个明文-密文对")
	plaintexts := []int64{1, 100, 1000, 123456789}
	ciphertexts := make([]int64, len(plaintexts))
	for i, p := range plaintexts {
		ciphertexts[i] = targetObfuscator.Obfuscate(p)
	}
	fmt.Println("已知多个明文-密文对:")
	for i := range plaintexts {
		fmt.Printf("  明文=%d, 密文=%d\n", plaintexts[i], ciphertexts[i])
	}
	fmt.Println("多个明文-密文对可以:")
	fmt.Println("  1. 更快地验证破解结果")
	fmt.Println("  2. 提高破解的准确性")
	fmt.Println("  3. 但破解方法相同（反推 xorKey）")

	// 场景4: 置换表的安全性
	fmt.Println("\n【场景4】置换表的安全性")
	fmt.Println("置换表是固定的（硬编码在代码中）")
	fmt.Println("这意味着:")
	fmt.Println("  - 攻击者可以从代码中直接看到置换表")
	fmt.Println("  - 只需要破解 xorKey 即可")
	fmt.Println("  - 如果有 1 个明文-密文对，可以立即破解 xorKey")
	fmt.Println("改进建议:")
	fmt.Println("  - 可以将置换表作为密钥的一部分")
	fmt.Println("  - 或者使用动态生成的置换表")

	// 总结
	fmt.Println("\n=== 安全性总结 ===")
	fmt.Println("✓ 优点:")
	fmt.Println("  - 快速混淆，性能好")
	fmt.Println("  - 可逆（知道密钥可以还原）")
	fmt.Println("  - 位置换提供了额外的混淆层")
	fmt.Println("\n✗ 缺点:")
	fmt.Println("  - 不是加密算法，只是混淆算法")
	fmt.Println("  - 置换表是固定的（在代码中可见）")
	fmt.Println("  - 如果有 1 个已知明文-密文对，可以立即破解 xorKey")
	fmt.Println("  - 仅密文攻击虽然困难，但理论上可以通过暴力破解")
	fmt.Println("  - 不适合用于安全敏感的场景")
	fmt.Println("\n建议:")
	fmt.Println("  - 适用于: ID 混淆、URL 参数混淆、防止用户直接看到真实 ID")
	fmt.Println("  - 不适用于: 密码、敏感数据、需要真正加密的场景")
	fmt.Println("  - 如果用于生产环境，建议:")
	fmt.Println("    1. 定期更换 xorKey")
	fmt.Println("    2. 不要使用默认密钥")
	fmt.Println("    3. 使用随机生成的强密钥")
	fmt.Println("    4. 不要在客户端暴露密钥")
	fmt.Println("    5. 考虑将置换表作为密钥的一部分（动态生成）")
}
