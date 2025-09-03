package tools

import (
	"strings"
)

// ConvImagePath 将相对图片路径转换为绝对URL
//
// 如果路径已经是以 http 或 https 开头的绝对URL，则直接返回
// 否则，将 baseUrl 和 path 拼接成完整URL
//
// 参数:
//   - path: 图片路径，可以是相对路径或绝对URL
//   - baseUrl: 基础URL，用于拼接相对路径
//
// 返回:
//   - 完整的图片URL
func ConvImagePath(path string, baseUrl string) string {
	// 如果路径已经是绝对URL，直接返回
	if strings.HasPrefix(path, "http") {
		return path
	}

	// 确保baseUrl以/结尾
	if baseUrl != "" && !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}

	// 拼接路径
	return baseUrl + path
}
