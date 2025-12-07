package utils

import (
	"regexp"
	"strings"
)

// 操作系统常量
const (
	OSWindows = "Windows"
	OSmacOS   = "macOS"
	OSLinux   = "Linux"
	OSiOS     = "iOS"
	OSAndroid = "Android"
	OSUnknown = "Unknown"
)

// 浏览器常量
const (
	BrowserChrome  = "Chrome"
	BrowserFirefox = "Firefox"
	BrowserSafari  = "Safari"
	BrowserEdge    = "Edge"
	BrowserOpera   = "Opera"
	BrowserWeChat  = "WeChat"
	BrowserQQ      = "QQ Browser"
	BrowserUC      = "UC Browser"
	Browser360     = "360 Browser"
	BrowserSogou   = "Sogou Browser"
	BrowserBaidu   = "Baidu Browser"
	BrowserSamsung = "Samsung Internet"
	BrowserYandex  = "Yandex"
	BrowserBrave   = "Brave"
	BrowserUnknown = "Unknown"
)

// 设备类型常量
const (
	DeviceMobile  = "Mobile"
	DeviceTablet  = "Tablet"
	DeviceDesktop = "Desktop"
)

// 渲染引擎常量
const (
	EngineWebKit   = "WebKit"
	EngineGecko    = "Gecko"
	EngineBlink    = "Blink"
	EngineTrident  = "Trident"
	EngineEdgeHTML = "EdgeHTML"
	EngineUnknown  = "Unknown"
)

// 小程序平台常量
const (
	MiniProgramWeChat    = "WeChat"
	MiniProgramAlipay    = "Alipay"
	MiniProgramByteDance = "ByteDance"
	MiniProgramBaidu     = "Baidu"
	MiniProgramQQ        = "QQ"
	MiniProgramKuaishou  = "Kuaishou"
	MiniProgramJD        = "JD"
	MiniProgramMeituan   = "Meituan"
	MiniProgramDingTalk  = "DingTalk"
)

// 设备制造商常量
const (
	ManufacturerApple     = "Apple"
	ManufacturerSamsung   = "Samsung"
	ManufacturerHuawei    = "Huawei"
	ManufacturerXiaomi    = "Xiaomi"
	ManufacturerOPPO      = "OPPO"
	Manufacturervivo      = "vivo"
	ManufacturerOnePlus   = "OnePlus"
	ManufacturerMeizu     = "Meizu"
	ManufacturerLenovo    = "Lenovo"
	ManufacturerMicrosoft = "Microsoft"
	ManufacturerUnknown   = "Unknown"
)

// CPU 架构常量
const (
	ArchX64   = "x64"
	ArchArm64 = "arm64"
	ArchX86   = "x86"
	ArchArmv7 = "armv7"
)

// 平台标识常量
const (
	PlatformiOS        = "iOS"
	PlatformWindows    = "Windows"
	PlatformWin64      = "Win64"
	PlatformAndroid    = "Android"
	PlatformMacIntel   = "MacIntel"
	PlatformMacARM     = "MacARM"
	PlatformLinux      = "Linux"
	PlatformLinuxX8664 = "Linux x86_64"
	PlatformLinuxX86   = "Linux x86"
)

// 应用名称常量
const (
	AppWeChat = "WeChat"
	AppAlipay = "Alipay"
	AppTikTok = "TikTok/Douyin"
	AppWeibo  = "Weibo"
	AppQQ     = "QQ"
	AppBaidu  = "Baidu App"
)

// UserAgentInfo 用户环境信息结构体
type UserAgentInfo struct {
	OS             string // 操作系统 (Windows, macOS, Linux, iOS, Android 等)
	OSVersion      string // 操作系统版本
	Browser        string // 浏览器 (Chrome, Firefox, Safari, Edge 等)
	BrowserVersion string // 浏览器版本
	Device         string // 设备类型 (Mobile, Tablet, Desktop)
	IsMobile       bool   // 是否为移动设备
	IsBot          bool   // 是否为爬虫/机器人
	Raw            string // 原始 User-Agent 字符串

	// 新增字段
	Engine        string // 渲染引擎 (WebKit, Gecko, Blink, Trident, EdgeHTML 等)
	EngineVersion string // 渲染引擎版本
	DeviceModel   string // 设备型号 (iPhone14,2, SM-G973F 等)
	Manufacturer  string // 设备制造商 (Apple, Samsung, Huawei, Xiaomi 等)
	Architecture  string // CPU 架构 (x64, arm64, x86, armv7 等)
	Language      string // 语言代码 (zh-CN, en-US 等)
	AppName       string // 应用名称 (WeChat, Alipay, TikTok 等，如果是应用内浏览器)
	AppVersion    string // 应用版本
	Platform      string // 平台标识 (Win64, MacIntel, Linux x86_64 等)
	MiniProgram   string // 小程序平台 (WeChat, Alipay, ByteDance, Baidu, QQ 等)
	IsMiniProgram bool   // 是否为小程序环境
}

// ParseUserAgent 解析 User-Agent 字符串，返回用户环境信息
func ParseUserAgent(ua string) *UserAgentInfo {
	info := &UserAgentInfo{
		Raw: ua,
	}

	uaLower := strings.ToLower(ua)

	// 检测是否为机器人/爬虫
	info.IsBot = isBot(uaLower)

	// 检测设备类型 - 先检查平板，再检查移动设备
	if isTablet(uaLower) {
		info.Device = DeviceTablet
		info.IsMobile = false
	} else {
		info.IsMobile = isMobile(uaLower)
		if info.IsMobile {
			info.Device = DeviceMobile
		} else {
			info.Device = DeviceDesktop
		}
	}

	// 解析操作系统
	info.OS, info.OSVersion = parseOS(ua, uaLower)

	// 解析浏览器
	info.Browser, info.BrowserVersion = parseBrowser(ua, uaLower)

	// 解析渲染引擎
	info.Engine, info.EngineVersion = parseEngine(ua, uaLower)

	// 解析设备信息
	info.DeviceModel, info.Manufacturer = parseDeviceInfo(ua, uaLower)

	// 解析 CPU 架构
	info.Architecture = parseArchitecture(uaLower)

	// 解析语言
	info.Language = parseLanguage(ua)

	// 解析平台标识
	info.Platform = parsePlatform(ua, uaLower)

	// 解析应用信息（应用内浏览器）
	info.AppName, info.AppVersion = parseAppInfo(ua, uaLower)

	// 解析小程序信息
	info.MiniProgram, info.IsMiniProgram = parseMiniProgram(ua, uaLower)

	return info
}

// isBot 检测是否为机器人/爬虫
func isBot(uaLower string) bool {
	botKeywords := []string{
		"bot", "crawler", "spider", "scraper", "curl", "wget",
		"python", "java", "go-http", "httpclient", "apache",
		"postman", "insomnia", "rest client",
	}
	// 排除一些常见应用，它们包含 "client" 但不是机器人
	excludeKeywords := []string{
		"alipayclient", "micromessenger", "qq", "wechat",
	}

	for _, keyword := range botKeywords {
		if strings.Contains(uaLower, keyword) {
			// 检查是否在排除列表中
			isExcluded := false
			for _, exclude := range excludeKeywords {
				if strings.Contains(uaLower, exclude) {
					isExcluded = true
					break
				}
			}
			if !isExcluded {
				return true
			}
		}
	}
	return false
}

// isMobile 检测是否为移动设备（不包括平板）
func isMobile(uaLower string) bool {
	// 排除平板设备
	if isTablet(uaLower) {
		return false
	}

	mobileKeywords := []string{
		"mobile", "android", "iphone", "ipod", "blackberry",
		"windows phone", "opera mini", "iemobile",
	}
	for _, keyword := range mobileKeywords {
		if strings.Contains(uaLower, keyword) {
			return true
		}
	}
	return false
}

// isTablet 检测是否为平板设备
func isTablet(uaLower string) bool {
	// iPad
	if strings.Contains(uaLower, "ipad") {
		return true
	}

	// Android 平板：包含 android 但不包含 mobile，且通常有特定标识
	if strings.Contains(uaLower, "android") && !strings.Contains(uaLower, "mobile") {
		// 检查是否有平板标识
		if strings.Contains(uaLower, "tablet") {
			return true
		}
		// 某些 Android 平板型号（如 SM-T 系列）
		if match := regexp.MustCompile(`sm-t\d+`).FindStringSubmatch(uaLower); len(match) > 0 {
			return true
		}
	}

	// 其他平板设备
	if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "playbook") || strings.Contains(uaLower, "kindle") {
		return true
	}

	return false
}

// parseOS 解析操作系统信息
func parseOS(ua, uaLower string) (os, version string) {
	// iOS - 必须在 macOS 之前检查，因为 iOS UA 包含 "like Mac OS X"
	if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ipod") {
		os = OSiOS
		// 尝试多种格式匹配版本号
		// 格式1: CPU iPhone OS 16_0 或 CPU OS 16_0 (iPad)
		if match := regexp.MustCompile(`cpu\s+(?:iphone|ipad|ipod)?\s*os\s+([\d_]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = strings.ReplaceAll(match[1], "_", ".")
		} else if match := regexp.MustCompile(`cpu\s+(?:iphone|ipad|ipod)\s+os\s+([\d_]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = strings.ReplaceAll(match[1], "_", ".")
		} else if match := regexp.MustCompile(`os ([\d_]+)`).FindStringSubmatch(ua); len(match) > 1 {
			version = strings.ReplaceAll(match[1], "_", ".")
		} else if match := regexp.MustCompile(`iphone os ([\d_]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = strings.ReplaceAll(match[1], "_", ".")
		}
		return
	}

	// Windows
	if strings.Contains(uaLower, "windows") {
		os = OSWindows
		if match := regexp.MustCompile(`windows nt ([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		} else if strings.Contains(uaLower, "windows 10") {
			version = "10"
		} else if strings.Contains(uaLower, "windows 11") {
			version = "11"
		}
		return
	}

	// Android
	if strings.Contains(uaLower, "android") {
		os = OSAndroid
		// 尝试多种格式匹配版本号
		if match := regexp.MustCompile(`android ([\d.]+)`).FindStringSubmatch(ua); len(match) > 1 {
			version = match[1]
		} else if match := regexp.MustCompile(`android\s+([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// macOS - 必须在 iOS 之后检查
	if strings.Contains(uaLower, "mac os x") || strings.Contains(uaLower, "macintosh") {
		os = OSmacOS
		// 尝试多种格式匹配版本号，按优先级从高到低
		// 格式1: Intel Mac OS X 10_15_7 (下划线格式，如 Safari) - 必须包含下划线
		if match := regexp.MustCompile(`intel mac os x ([\d_]+)`).FindStringSubmatch(uaLower); len(match) > 1 && strings.Contains(match[1], "_") {
			version = strings.ReplaceAll(match[1], "_", ".")
		} else if match := regexp.MustCompile(`intel mac os x ([\d]+(?:\.[\d]+)+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			// 格式2: Intel Mac OS X 10.15 (点号格式，如 Firefox) - 必须包含点号
			version = match[1]
		} else if match := regexp.MustCompile(`mac os x ([\d_]+)`).FindStringSubmatch(uaLower); len(match) > 1 && strings.Contains(match[1], "_") {
			// 格式3: Mac OS X 10_15_7 (下划线格式) - 必须包含下划线
			version = strings.ReplaceAll(match[1], "_", ".")
		} else if match := regexp.MustCompile(`mac os x ([\d]+(?:\.[\d]+)+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			// 格式4: Mac OS X 10.15 (点号格式) - 必须包含点号
			version = match[1]
		} else if match := regexp.MustCompile(`intel mac os x ([\d]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			// 格式5: Intel Mac OS X 10 (只有主版本号)
			version = match[1]
		} else if match := regexp.MustCompile(`mac os x ([\d]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			// 格式6: Mac OS X 10 (只有主版本号)
			version = match[1]
		}
		return
	}

	// Linux
	if strings.Contains(uaLower, "linux") && !strings.Contains(uaLower, "android") {
		os = OSLinux
		// Linux 版本信息通常在 User-Agent 中不明确
		return
	}

	// 其他
	return OSUnknown, ""
}

// parseBrowser 解析浏览器信息
func parseBrowser(ua, uaLower string) (browser, version string) {
	// Edge Chromium - 必须在 Chrome 之前检查
	if strings.Contains(uaLower, "edg/") || strings.Contains(uaLower, "edgios/") || strings.Contains(uaLower, "edga/") {
		browser = BrowserEdge
		if match := regexp.MustCompile(`edg[eiosa]?/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// QQ 浏览器 - 必须在 Chrome 之前检查（因为包含 chrome/）
	if strings.Contains(uaLower, "mqqbrowser") {
		browser = BrowserQQ
		if match := regexp.MustCompile(`mqqbrowser/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// UC 浏览器 - 必须在 Chrome 之前检查
	if strings.Contains(uaLower, "ucbrowser") || strings.Contains(uaLower, "ucweb") {
		browser = BrowserUC
		if match := regexp.MustCompile(`(?:ucbrowser|ucweb)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// 360 浏览器 - 必须在 Chrome 之前检查
	if strings.Contains(uaLower, "360se") || strings.Contains(uaLower, "360ee") || strings.Contains(uaLower, "qihu") {
		browser = Browser360
		if match := regexp.MustCompile(`(?:360se|360ee|qihu)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Opera - 必须在 Chrome 之前检查（因为包含 chrome/）
	if strings.Contains(uaLower, "opera/") || strings.Contains(uaLower, "opr/") {
		browser = BrowserOpera
		if match := regexp.MustCompile(`(?:opera|opr)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Chrome
	if strings.Contains(uaLower, "chrome/") && !strings.Contains(uaLower, "chromium") {
		browser = BrowserChrome
		if match := regexp.MustCompile(`chrome/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Safari
	if strings.Contains(uaLower, "safari/") && !strings.Contains(uaLower, "chrome") {
		browser = BrowserSafari
		if match := regexp.MustCompile(`version/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Firefox
	if strings.Contains(uaLower, "firefox/") {
		browser = BrowserFirefox
		if match := regexp.MustCompile(`firefox/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// 微信内置浏览器
	if strings.Contains(uaLower, "micromessenger") {
		browser = BrowserWeChat
		if match := regexp.MustCompile(`micromessenger/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// 搜狗浏览器
	if strings.Contains(uaLower, "metasr") || strings.Contains(uaLower, "sogou") {
		browser = BrowserSogou
		if match := regexp.MustCompile(`(?:metasr|sogou)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// 百度浏览器
	if strings.Contains(uaLower, "baiduboxapp") || strings.Contains(uaLower, "baidubrowser") {
		browser = BrowserBaidu
		if match := regexp.MustCompile(`(?:baiduboxapp|baidubrowser)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Samsung Internet
	if strings.Contains(uaLower, "samsungbrowser") {
		browser = BrowserSamsung
		if match := regexp.MustCompile(`samsungbrowser/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Yandex Browser
	if strings.Contains(uaLower, "yabrowser") {
		browser = BrowserYandex
		if match := regexp.MustCompile(`yabrowser/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Brave Browser
	if strings.Contains(uaLower, "brave") {
		browser = BrowserBrave
		if match := regexp.MustCompile(`brave/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// 其他
	return BrowserUnknown, ""
}

// parseMiniProgram 解析小程序平台信息
func parseMiniProgram(ua, uaLower string) (platform string, isMiniProgram bool) {
	// 微信小程序
	// 微信小程序的 User-Agent 通常包含 "miniprogram" 或 "miniProgram"
	// 例如: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.x NetType/WIFI Language/zh_CN miniProgram"
	if strings.Contains(uaLower, "miniprogram") {
		// 确认是微信环境
		if strings.Contains(uaLower, "micromessenger") {
			return MiniProgramWeChat, true
		}
	}

	// 支付宝小程序
	// 支付宝小程序的 User-Agent 通常包含 "alipay" 和 "miniprogram"
	// 例如: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 AlipayClient/10.x.x AlipayApp/10.x.x Language/zh-Hans miniProgram"
	if strings.Contains(uaLower, "alipay") && (strings.Contains(uaLower, "miniprogram") || strings.Contains(uaLower, "alipayclient")) {
		// 进一步确认是否是小程序（通常会有特定的标识）
		if strings.Contains(uaLower, "miniprogram") || strings.Contains(uaLower, "alipayapp") {
			return MiniProgramAlipay, true
		}
	}

	// 抖音小程序 / 字节跳动小程序
	// 抖音小程序的 User-Agent 通常包含 "ttwebview" 或 "bytedancewebview"
	// 例如: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 ByteDanceWebview/8.x.x"
	if strings.Contains(uaLower, "ttwebview") || strings.Contains(uaLower, "bytedancewebview") {
		return MiniProgramByteDance, true
	}

	// 百度小程序
	// 百度小程序的 User-Agent 通常包含 "swan" 或 "baiduboxapp" 和 "miniprogram"
	// 例如: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 swan/2.x.x"
	if strings.Contains(uaLower, "swan") || (strings.Contains(uaLower, "baiduboxapp") && strings.Contains(uaLower, "miniprogram")) {
		return MiniProgramBaidu, true
	}

	// QQ 小程序
	// QQ 小程序的 User-Agent 通常包含 "qq" 和 "miniprogram"
	// 例如: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 QQ/8.x.x NetType/WIFI miniProgram"
	if strings.Contains(uaLower, "qq/") && strings.Contains(uaLower, "miniprogram") {
		return MiniProgramQQ, true
	}

	// 快手小程序
	// 快手小程序的 User-Agent 通常包含 "kswebview" 或 "kwai"
	if strings.Contains(uaLower, "kswebview") || (strings.Contains(uaLower, "kwai") && strings.Contains(uaLower, "miniprogram")) {
		return MiniProgramKuaishou, true
	}

	// 京东小程序
	// 京东小程序的 User-Agent 通常包含 "jdapp" 和 "miniprogram"
	if strings.Contains(uaLower, "jdapp") && strings.Contains(uaLower, "miniprogram") {
		return MiniProgramJD, true
	}

	// 美团小程序
	// 美团小程序的 User-Agent 通常包含 "meituan" 和 "miniprogram"
	if strings.Contains(uaLower, "meituan") && strings.Contains(uaLower, "miniprogram") {
		return MiniProgramMeituan, true
	}

	// 钉钉小程序
	// 钉钉小程序的 User-Agent 通常包含 "dingtalk" 和 "miniprogram"
	if strings.Contains(uaLower, "dingtalk") && strings.Contains(uaLower, "miniprogram") {
		return MiniProgramDingTalk, true
	}

	return "", false
}

// parseEngine 解析渲染引擎信息
func parseEngine(ua, uaLower string) (engine, version string) {
	// Blink (Chrome 28+) - 必须在 WebKit 之前检查，因为 Blink 基于 WebKit
	if strings.Contains(uaLower, "chrome") && strings.Contains(uaLower, "webkit") {
		engine = EngineBlink
		if match := regexp.MustCompile(`chrome/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// WebKit (Safari 等，但不包括 Chrome)
	if strings.Contains(uaLower, "webkit") && !strings.Contains(uaLower, "chrome") {
		engine = EngineWebKit
		if match := regexp.MustCompile(`webkit/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Gecko (Firefox) - 必须在 WebKit 之后检查，避免 "like Gecko" 误判
	if strings.Contains(uaLower, "gecko") && !strings.Contains(uaLower, "webkit") {
		engine = EngineGecko
		if match := regexp.MustCompile(`gecko/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		} else if match := regexp.MustCompile(`rv:([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// Trident (IE)
	if strings.Contains(uaLower, "trident") {
		engine = EngineTrident
		if match := regexp.MustCompile(`trident/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	// EdgeHTML (旧版 Edge)
	if strings.Contains(uaLower, "edge/") && strings.Contains(uaLower, "edgehtml") {
		engine = EngineEdgeHTML
		if match := regexp.MustCompile(`edgehtml/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			version = match[1]
		}
		return
	}

	return EngineUnknown, ""
}

// parseDeviceInfo 解析设备型号和制造商
func parseDeviceInfo(ua, uaLower string) (model, manufacturer string) {
	// iPhone 型号
	if strings.Contains(uaLower, "iphone") {
		manufacturer = ManufacturerApple
		if match := regexp.MustCompile(`iphone(\d+,\d+)`).FindStringSubmatch(ua); len(match) > 1 {
			model = "iPhone" + match[1]
		} else {
			model = "iPhone"
		}
		return
	}

	// iPad 型号
	if strings.Contains(uaLower, "ipad") {
		manufacturer = ManufacturerApple
		if match := regexp.MustCompile(`ipad(\d+,\d+)`).FindStringSubmatch(ua); len(match) > 1 {
			model = "iPad" + match[1]
		} else {
			model = "iPad"
		}
		return
	}

	// iPod
	if strings.Contains(uaLower, "ipod") {
		manufacturer = ManufacturerApple
		model = "iPod"
		return
	}

	// Android 设备
	if strings.Contains(uaLower, "android") {
		// 尝试提取设备型号 - 格式: (Linux; Android 13; SM-G973F) 或 (Linux; U; Android 13; zh-cn; SM-G973F)
		// 方法1: 匹配括号内最后一个分号后的内容（通常是设备型号）
		if match := regexp.MustCompile(`\([^)]*android[^)]*;\s*([a-z0-9\s\-_]+)\)`).FindStringSubmatch(uaLower); len(match) > 1 {
			candidate := strings.TrimSpace(match[1])
			// 排除语言代码等非设备型号的内容
			if candidate != "" && candidate != "u" && len(candidate) > 1 {
				// 检查是否是语言代码（通常是2-5个字符，可能包含连字符）
				if !regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`).MatchString(candidate) {
					model = candidate
				}
			}
		}
		// 方法2: 如果方法1失败，尝试匹配 Android 版本号后的设备型号
		if model == "" {
			if match := regexp.MustCompile(`android\s+[\d.]+\s*;\s*([a-z0-9\s\-_]+?)(?:\s+applewebkit|\s+build|$|;|\s+mqqbrowser|\s+ucbrowser)`).FindStringSubmatch(uaLower); len(match) > 1 {
				candidate := strings.TrimSpace(match[1])
				if candidate != "" && candidate != "u" && len(candidate) > 1 {
					if !regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`).MatchString(candidate) {
						model = candidate
					}
				}
			}
		}
		// 方法3: 如果方法1和2都失败，尝试其他格式
		if model == "" {
			if match := regexp.MustCompile(`;\s*([a-z0-9\s\-_]+)\s+build`).FindStringSubmatch(uaLower); len(match) > 1 {
				model = strings.TrimSpace(match[1])
			} else if match := regexp.MustCompile(`;\s*([a-z0-9\s\-_]+)\s+miui`).FindStringSubmatch(uaLower); len(match) > 1 {
				model = strings.TrimSpace(match[1])
				manufacturer = ManufacturerXiaomi
			} else if match := regexp.MustCompile(`;\s*([a-z0-9\s\-_]+)\s+android`).FindStringSubmatch(uaLower); len(match) > 1 {
				model = strings.TrimSpace(match[1])
			}
		}

		// 识别常见制造商（必须在设备型号提取之后）
		if strings.Contains(uaLower, "samsung") || strings.Contains(uaLower, "sm-") {
			manufacturer = ManufacturerSamsung
		} else if strings.Contains(uaLower, "huawei") || strings.Contains(uaLower, "honor") {
			manufacturer = ManufacturerHuawei
		} else if strings.Contains(uaLower, "xiaomi") || strings.Contains(uaLower, "redmi") || strings.Contains(uaLower, "mi ") || strings.Contains(uaLower, "miui") {
			manufacturer = ManufacturerXiaomi
		} else if strings.Contains(uaLower, "oppo") {
			manufacturer = ManufacturerOPPO
		} else if strings.Contains(uaLower, "vivo") {
			manufacturer = Manufacturervivo
		} else if strings.Contains(uaLower, "oneplus") {
			manufacturer = ManufacturerOnePlus
		} else if strings.Contains(uaLower, "meizu") {
			manufacturer = ManufacturerMeizu
		} else if strings.Contains(uaLower, "lenovo") {
			manufacturer = ManufacturerLenovo
		} else if manufacturer == "" {
			manufacturer = ManufacturerUnknown
		}

		return
	}

	// macOS
	if strings.Contains(uaLower, "mac os x") || strings.Contains(uaLower, "macintosh") {
		manufacturer = ManufacturerApple
		model = "Mac"
		return
	}

	// Windows
	if strings.Contains(uaLower, "windows") {
		manufacturer = ManufacturerMicrosoft
		model = "PC"
		return
	}

	return "", ""
}

// parseArchitecture 解析 CPU 架构
func parseArchitecture(uaLower string) string {
	// x64
	if strings.Contains(uaLower, "x64") || strings.Contains(uaLower, "x86_64") || strings.Contains(uaLower, "win64") || strings.Contains(uaLower, "wow64") {
		return ArchX64
	}

	// macOS Intel 通常是 x64
	if strings.Contains(uaLower, "intel") && (strings.Contains(uaLower, "mac") || strings.Contains(uaLower, "macintosh")) {
		return ArchX64
	}

	// arm64
	if strings.Contains(uaLower, "arm64") || strings.Contains(uaLower, "aarch64") {
		return ArchArm64
	}

	// x86
	if strings.Contains(uaLower, "x86") && !strings.Contains(uaLower, "x64") {
		return ArchX86
	}

	// armv7
	if strings.Contains(uaLower, "armv7") {
		return ArchArmv7
	}

	// 移动设备通常使用 ARM
	if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "android") {
		return ArchArm64
	}

	return ""
}

// parseLanguage 解析语言代码
func parseLanguage(ua string) string {
	// 语言代码通常在 UA 末尾，格式如 "zh-CN" 或 "en-US"
	if match := regexp.MustCompile(`[;\s]([a-z]{2}(?:-[a-z]{2})?)(?:[,;]|$)`).FindStringSubmatch(strings.ToLower(ua)); len(match) > 1 {
		// 检查是否是有效的语言代码位置
		lang := match[1]
		if len(lang) >= 2 && len(lang) <= 5 {
			return lang
		}
	}
	return ""
}

// parsePlatform 解析平台标识
func parsePlatform(ua, uaLower string) string {
	// iOS - 必须在 macOS 之前检查
	if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ipod") {
		return PlatformiOS
	}

	// Windows
	if strings.Contains(uaLower, "windows") {
		if strings.Contains(uaLower, "win64") || strings.Contains(uaLower, "wow64") {
			return PlatformWin64
		}
		return PlatformWindows
	}

	// Android
	if strings.Contains(uaLower, "android") {
		return PlatformAndroid
	}

	// macOS - 必须在 iOS 之后检查
	if strings.Contains(uaLower, "mac os x") || strings.Contains(uaLower, "macintosh") {
		if strings.Contains(uaLower, "intel") {
			return PlatformMacIntel
		} else if strings.Contains(uaLower, "arm") {
			return PlatformMacARM
		}
		return PlatformMacIntel
	}

	// Linux
	if strings.Contains(uaLower, "linux") && !strings.Contains(uaLower, "android") {
		if strings.Contains(uaLower, "x86_64") {
			return PlatformLinuxX8664
		} else if strings.Contains(uaLower, "x86") {
			return PlatformLinuxX86
		}
		return PlatformLinux
	}

	return ""
}

// parseAppInfo 解析应用内浏览器信息
func parseAppInfo(ua, uaLower string) (appName, appVersion string) {
	// 微信
	if strings.Contains(uaLower, "micromessenger") {
		appName = AppWeChat
		if match := regexp.MustCompile(`micromessenger/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	// 支付宝
	if strings.Contains(uaLower, "alipayclient") {
		appName = AppAlipay
		if match := regexp.MustCompile(`alipayclient/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	// 抖音
	if strings.Contains(uaLower, "aweme") || strings.Contains(uaLower, "douyin") {
		appName = AppTikTok
		if match := regexp.MustCompile(`(?:aweme|douyin)/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	// 微博
	if strings.Contains(uaLower, "weibo") {
		appName = AppWeibo
		if match := regexp.MustCompile(`weibo/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	// QQ
	if strings.Contains(uaLower, "qq/") {
		appName = AppQQ
		if match := regexp.MustCompile(`qq/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	// 百度 App
	if strings.Contains(uaLower, "baiduboxapp") {
		appName = AppBaidu
		if match := regexp.MustCompile(`baiduboxapp/([\d.]+)`).FindStringSubmatch(uaLower); len(match) > 1 {
			appVersion = match[1]
		}
		return
	}

	return "", ""
}
