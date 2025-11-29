package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		expected *UserAgentInfo
	}{
		{
			name: "Chrome on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expected: &UserAgentInfo{
				OS:             OSWindows,
				OSVersion:      "10.0",
				Browser:        BrowserChrome,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineBlink,
				Architecture:   ArchX64,
				Platform:       PlatformWin64,
				Manufacturer:   ManufacturerMicrosoft,
				DeviceModel:    "PC",
			},
		},
		{
			name: "Firefox on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
			expected: &UserAgentInfo{
				OS:             OSmacOS,
				OSVersion:      "10.15",
				Browser:        BrowserFirefox,
				BrowserVersion: "121.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineGecko,
				Architecture:   ArchX64,
				Platform:       PlatformMacIntel,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "Mac",
			},
		},
		{
			name: "Safari on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			expected: &UserAgentInfo{
				OS:             OSmacOS,
				OSVersion:      "10.15.7",
				Browser:        BrowserSafari,
				BrowserVersion: "17.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineWebKit,
				Platform:       PlatformMacIntel,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "Mac",
			},
		},
		{
			name: "Edge on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			expected: &UserAgentInfo{
				OS:             OSWindows,
				OSVersion:      "10.0",
				Browser:        BrowserEdge,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineBlink,
				Architecture:   ArchX64,
				Platform:       PlatformWin64,
				Manufacturer:   ManufacturerMicrosoft,
				DeviceModel:    "PC",
			},
		},
		{
			name: "iPhone Safari",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
			expected: &UserAgentInfo{
				OS:             OSiOS,
				OSVersion:      "16.0",
				Browser:        BrowserSafari,
				BrowserVersion: "16.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Engine:         EngineWebKit,
				Architecture:   ArchArm64,
				Platform:       PlatformiOS,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "iPhone",
			},
		},
		{
			name: "Android Chrome",
			ua:   "Mozilla/5.0 (Linux; Android 13; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			expected: &UserAgentInfo{
				OS:             OSAndroid,
				OSVersion:      "13",
				Browser:        BrowserChrome,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Engine:         EngineBlink,
				Architecture:   ArchArm64,
				Platform:       PlatformAndroid,
				Manufacturer:   ManufacturerSamsung,
				DeviceModel:    "sm-g973f",
			},
		},
		{
			name: "iPad Safari",
			ua:   "Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
			expected: &UserAgentInfo{
				OS:             OSiOS,
				OSVersion:      "16.0",
				Browser:        BrowserSafari,
				BrowserVersion: "16.0",
				Device:         DeviceTablet,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineWebKit,
				Architecture:   ArchArm64,
				Platform:       PlatformiOS,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "iPad",
			},
		},
		{
			name: "微信小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.0 NetType/WIFI Language/zh_CN miniProgram",
			expected: &UserAgentInfo{
				OS:             OSiOS,
				OSVersion:      "16.0",
				Browser:        BrowserWeChat,
				BrowserVersion: "8.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				AppName:        AppWeChat,
				AppVersion:     "8.0.0",
				MiniProgram:    MiniProgramWeChat,
				IsMiniProgram:  true,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "iPhone",
				Platform:       PlatformiOS,
			},
		},
		{
			name: "支付宝小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 AlipayClient/10.2.0 AlipayApp/10.2.0 Language/zh-Hans miniProgram",
			expected: &UserAgentInfo{
				OS:            OSiOS,
				OSVersion:     "16.0",
				Device:        DeviceMobile,
				IsMobile:      true,
				IsBot:         false,
				AppName:       AppAlipay,
				AppVersion:    "10.2.0",
				MiniProgram:   MiniProgramAlipay,
				IsMiniProgram: true,
				Manufacturer:  ManufacturerApple,
				DeviceModel:   "iPhone",
				Platform:      PlatformiOS,
			},
		},
		{
			name: "抖音小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 ByteDanceWebview/8.0.0",
			expected: &UserAgentInfo{
				OS:            OSiOS,
				OSVersion:     "16.0",
				Device:        DeviceMobile,
				IsMobile:      true,
				IsBot:         false,
				MiniProgram:   MiniProgramByteDance,
				IsMiniProgram: true,
				Manufacturer:  ManufacturerApple,
				DeviceModel:   "iPhone",
				Platform:      PlatformiOS,
			},
		},
		{
			name: "百度小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 swan/2.0.0",
			expected: &UserAgentInfo{
				OS:            OSiOS,
				OSVersion:     "16.0",
				Device:        DeviceMobile,
				IsMobile:      true,
				IsBot:         false,
				MiniProgram:   MiniProgramBaidu,
				IsMiniProgram: true,
				Manufacturer:  ManufacturerApple,
				DeviceModel:   "iPhone",
				Platform:      PlatformiOS,
			},
		},
		{
			name: "QQ小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 QQ/8.0.0 NetType/WIFI miniProgram",
			expected: &UserAgentInfo{
				OS:            OSiOS,
				OSVersion:     "16.0",
				Device:        DeviceMobile,
				IsMobile:      true,
				IsBot:         false,
				AppName:       AppQQ,
				AppVersion:    "8.0.0",
				MiniProgram:   MiniProgramQQ,
				IsMiniProgram: true,
				Manufacturer:  ManufacturerApple,
				DeviceModel:   "iPhone",
				Platform:      PlatformiOS,
			},
		},
		{
			name: "微信浏览器（非小程序）",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.0 NetType/WIFI Language/zh_CN",
			expected: &UserAgentInfo{
				OS:             OSiOS,
				OSVersion:      "16.0",
				Browser:        BrowserWeChat,
				BrowserVersion: "8.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				AppName:        AppWeChat,
				AppVersion:     "8.0.0",
				MiniProgram:    "",
				IsMiniProgram:  false,
				Manufacturer:   ManufacturerApple,
				DeviceModel:    "iPhone",
				Platform:       PlatformiOS,
			},
		},
		{
			name: "QQ浏览器",
			ua:   "Mozilla/5.0 (Linux; Android 13; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36 MQQBrowser/10.0",
			expected: &UserAgentInfo{
				OS:             OSAndroid,
				OSVersion:      "13",
				Browser:        BrowserQQ,
				BrowserVersion: "10.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Manufacturer:   ManufacturerSamsung,
				DeviceModel:    "sm-g973f",
				Platform:       PlatformAndroid,
			},
		},
		{
			name: "UC浏览器",
			ua:   "Mozilla/5.0 (Linux; U; Android 13; zh-cn; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36 UCBrowser/13.0.0",
			expected: &UserAgentInfo{
				OS:             OSAndroid,
				OSVersion:      "13",
				Browser:        BrowserUC,
				BrowserVersion: "13.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Manufacturer:   ManufacturerSamsung,
				DeviceModel:    "sm-g973f",
				Platform:       PlatformAndroid,
			},
		},
		{
			name: "360浏览器",
			ua:   "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 360SE/13.0",
			expected: &UserAgentInfo{
				OS:             OSWindows,
				OSVersion:      "10.0",
				Browser:        Browser360,
				BrowserVersion: "13.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Architecture:   ArchX64,
				Platform:       PlatformWin64,
				Manufacturer:   ManufacturerMicrosoft,
				DeviceModel:    "PC",
			},
		},
		{
			name: "机器人 - Googlebot",
			ua:   "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expected: &UserAgentInfo{
				Device:       DeviceDesktop,
				IsMobile:     false,
				IsBot:        true,
				Manufacturer: "",
				DeviceModel:  "",
			},
		},
		{
			name: "机器人 - curl",
			ua:   "curl/7.68.0",
			expected: &UserAgentInfo{
				Device:       DeviceDesktop,
				IsMobile:     false,
				IsBot:        true,
				Manufacturer: "",
				DeviceModel:  "",
			},
		},
		{
			name: "Linux Chrome",
			ua:   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expected: &UserAgentInfo{
				OS:             OSLinux,
				Browser:        BrowserChrome,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Engine:         EngineBlink,
				Architecture:   ArchX64,
				Platform:       PlatformLinuxX8664,
			},
		},
		{
			name: "华为手机",
			ua:   "Mozilla/5.0 (Linux; Android 13; HUAWEI P50 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			expected: &UserAgentInfo{
				OS:             OSAndroid,
				OSVersion:      "13",
				Browser:        BrowserChrome,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Manufacturer:   ManufacturerHuawei,
				Platform:       PlatformAndroid,
			},
		},
		{
			name: "小米手机",
			ua:   "Mozilla/5.0 (Linux; Android 13; Mi 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			expected: &UserAgentInfo{
				OS:             OSAndroid,
				OSVersion:      "13",
				Browser:        BrowserChrome,
				BrowserVersion: "120.0.0.0",
				Device:         DeviceMobile,
				IsMobile:       true,
				IsBot:          false,
				Manufacturer:   ManufacturerXiaomi,
				Platform:       PlatformAndroid,
			},
		},
		{
			name: "Opera浏览器",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
			expected: &UserAgentInfo{
				OS:             OSWindows,
				OSVersion:      "10.0",
				Browser:        BrowserOpera,
				BrowserVersion: "106.0.0.0",
				Device:         DeviceDesktop,
				IsMobile:       false,
				IsBot:          false,
				Architecture:   ArchX64,
				Platform:       PlatformWin64,
				Manufacturer:   ManufacturerMicrosoft,
				DeviceModel:    "PC",
			},
		},
		{
			name: "空字符串",
			ua:   "",
			expected: &UserAgentInfo{
				Raw:      "",
				Device:   DeviceDesktop,
				IsMobile: false,
				IsBot:    false,
			},
		},
		{
			name: "未知浏览器",
			ua:   "SomeUnknownBrowser/1.0",
			expected: &UserAgentInfo{
				Browser:  BrowserUnknown,
				Device:   DeviceDesktop,
				IsMobile: false,
				IsBot:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUserAgent(tt.ua)

			// 验证基本字段
			if got.Raw != tt.ua {
				t.Errorf("Raw = %v, want %v", got.Raw, tt.ua)
			}

			// 验证期望的字段
			if tt.expected.OS != "" && got.OS != tt.expected.OS {
				t.Errorf("OS = %v, want %v", got.OS, tt.expected.OS)
			}
			if tt.expected.OSVersion != "" && got.OSVersion != tt.expected.OSVersion {
				t.Errorf("OSVersion = %v, want %v", got.OSVersion, tt.expected.OSVersion)
			}
			if tt.expected.Browser != "" && got.Browser != tt.expected.Browser {
				t.Errorf("Browser = %v, want %v", got.Browser, tt.expected.Browser)
			}
			if tt.expected.BrowserVersion != "" && got.BrowserVersion != tt.expected.BrowserVersion {
				t.Errorf("BrowserVersion = %v, want %v", got.BrowserVersion, tt.expected.BrowserVersion)
			}
			if got.Device != tt.expected.Device {
				t.Errorf("Device = %v, want %v", got.Device, tt.expected.Device)
			}
			if got.IsMobile != tt.expected.IsMobile {
				t.Errorf("IsMobile = %v, want %v", got.IsMobile, tt.expected.IsMobile)
			}
			if got.IsBot != tt.expected.IsBot {
				t.Errorf("IsBot = %v, want %v", got.IsBot, tt.expected.IsBot)
			}
			if tt.expected.Engine != "" && got.Engine != tt.expected.Engine {
				t.Errorf("Engine = %v, want %v", got.Engine, tt.expected.Engine)
			}
			if tt.expected.Architecture != "" && got.Architecture != tt.expected.Architecture {
				t.Errorf("Architecture = %v, want %v", got.Architecture, tt.expected.Architecture)
			}
			if tt.expected.Platform != "" && got.Platform != tt.expected.Platform {
				t.Errorf("Platform = %v, want %v", got.Platform, tt.expected.Platform)
			}
			if tt.expected.Manufacturer != "" && got.Manufacturer != tt.expected.Manufacturer {
				t.Errorf("Manufacturer = %v, want %v", got.Manufacturer, tt.expected.Manufacturer)
			}
			if tt.expected.DeviceModel != "" && got.DeviceModel != tt.expected.DeviceModel {
				t.Errorf("DeviceModel = %v, want %v", got.DeviceModel, tt.expected.DeviceModel)
			}
			if tt.expected.AppName != "" && got.AppName != tt.expected.AppName {
				t.Errorf("AppName = %v, want %v", got.AppName, tt.expected.AppName)
			}
			if tt.expected.AppVersion != "" && got.AppVersion != tt.expected.AppVersion {
				t.Errorf("AppVersion = %v, want %v", got.AppVersion, tt.expected.AppVersion)
			}
			if tt.expected.MiniProgram != "" && got.MiniProgram != tt.expected.MiniProgram {
				t.Errorf("MiniProgram = %v, want %v", got.MiniProgram, tt.expected.MiniProgram)
			}
			if got.IsMiniProgram != tt.expected.IsMiniProgram {
				t.Errorf("IsMiniProgram = %v, want %v", got.IsMiniProgram, tt.expected.IsMiniProgram)
			}
		})
	}
}

func TestParseUserAgent_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		ua    string
		check func(t *testing.T, info *UserAgentInfo)
	}{
		{
			name: "Windows 11",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			check: func(t *testing.T, info *UserAgentInfo) {
				if info.OS != OSWindows {
					t.Errorf("Expected Windows, got %s", info.OS)
				}
			},
		},
		{
			name: "包含语言信息",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36; zh-CN",
			check: func(t *testing.T, info *UserAgentInfo) {
				// 语言解析可能不总是准确，这里只检查不报错
				if info.OS != OSWindows {
					t.Errorf("Expected Windows, got %s", info.OS)
				}
			},
		},
		{
			name: "Android平板",
			ua:   "Mozilla/5.0 (Linux; Android 13; SM-T970) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			check: func(t *testing.T, info *UserAgentInfo) {
				if info.Device != DeviceTablet {
					t.Errorf("Expected Tablet, got %s", info.Device)
				}
				if info.IsMobile {
					t.Errorf("Expected IsMobile=false for tablet, got true")
				}
			},
		},
		{
			name: "快手小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 KuaishouApp/10.0.0 kswebview miniProgram",
			check: func(t *testing.T, info *UserAgentInfo) {
				if !info.IsMiniProgram {
					t.Errorf("Expected IsMiniProgram=true, got false")
				}
				if info.MiniProgram != MiniProgramKuaishou {
					t.Errorf("Expected MiniProgram=Kuaishou, got %s", info.MiniProgram)
				}
			},
		},
		{
			name: "京东小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 jdapp/10.0.0 miniProgram",
			check: func(t *testing.T, info *UserAgentInfo) {
				if !info.IsMiniProgram {
					t.Errorf("Expected IsMiniProgram=true, got false")
				}
				if info.MiniProgram != MiniProgramJD {
					t.Errorf("Expected MiniProgram=JD, got %s", info.MiniProgram)
				}
			},
		},
		{
			name: "美团小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 meituan/10.0.0 miniProgram",
			check: func(t *testing.T, info *UserAgentInfo) {
				if !info.IsMiniProgram {
					t.Errorf("Expected IsMiniProgram=true, got false")
				}
				if info.MiniProgram != MiniProgramMeituan {
					t.Errorf("Expected MiniProgram=Meituan, got %s", info.MiniProgram)
				}
			},
		},
		{
			name: "钉钉小程序",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 dingtalk/10.0.0 miniProgram",
			check: func(t *testing.T, info *UserAgentInfo) {
				if !info.IsMiniProgram {
					t.Errorf("Expected IsMiniProgram=true, got false")
				}
				if info.MiniProgram != MiniProgramDingTalk {
					t.Errorf("Expected MiniProgram=DingTalk, got %s", info.MiniProgram)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ParseUserAgent(tt.ua)
			if tt.check != nil {
				tt.check(t, info)
			}
		})
	}
}

func TestParseUserAgent_Consistency(t *testing.T) {
	// 测试相同 UA 多次解析结果一致
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	first := ParseUserAgent(ua)
	second := ParseUserAgent(ua)

	if !reflect.DeepEqual(first, second) {
		t.Errorf("ParseUserAgent results are not consistent:\nFirst:  %+v\nSecond: %+v", first, second)
	}
}

func TestParseUserAgent_RealWorldExamples(t *testing.T) {
	// 真实世界的 User-Agent 示例
	realWorldUAs := []string{
		// Chrome Windows
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		// Firefox macOS
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
		// Safari iOS
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
		// Android Chrome
		"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		// 微信小程序
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.0 NetType/WIFI Language/zh_CN miniProgram",
	}

	for i, ua := range realWorldUAs {
		t.Run(fmt.Sprintf("RealWorld_%d", i), func(t *testing.T) {
			info := ParseUserAgent(ua)

			// 基本验证：不应该 panic，应该有原始 UA
			if info == nil {
				t.Fatal("ParseUserAgent returned nil")
			}
			if info.Raw != ua {
				t.Errorf("Raw UA mismatch: got %q, want %q", info.Raw, ua)
			}

			// 验证设备类型应该被正确识别
			if info.Device == "" {
				t.Error("Device should not be empty")
			}

			// 验证设备类型值应该是有效的
			validDevices := map[string]bool{DeviceMobile: true, DeviceTablet: true, DeviceDesktop: true}
			if !validDevices[info.Device] {
				t.Errorf("Invalid device type: %s", info.Device)
			}
		})
	}
}
