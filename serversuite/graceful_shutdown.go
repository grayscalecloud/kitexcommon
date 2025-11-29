// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serversuite

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

// GracefulShutdownOptions 优雅退出配置选项
type GracefulShutdownOptions struct {
	// ShutdownTimeout 优雅退出超时时间，默认 30 秒
	// 注意：Kitex 的 Stop() 方法本身没有超时控制，这里用于包装超时逻辑
	ShutdownTimeout time.Duration
	// CleanupFunc 自定义清理函数，在服务器停止后执行
	// 该函数会通过 Kitex 的 RegisterShutdownHook 注册，在 Stop() 时自动执行
	CleanupFunc func() error
	// BeforeShutdownFunc 在开始关闭前执行的函数
	// 该函数在收到退出信号后、调用 Stop() 之前执行
	// 可以用于：停止接收新请求、设置健康检查状态为不健康、通知其他服务等
	// 注意：Kitex 的 Stop() 会自动注销注册中心，所以不需要在这里手动注销
	BeforeShutdownFunc func()
	// BeforeShutdownTimeout BeforeShutdownFunc 执行的超时时间，默认 5 秒
	BeforeShutdownTimeout time.Duration
}

// defaultGracefulShutdownOptions 返回默认的优雅退出配置
func defaultGracefulShutdownOptions() *GracefulShutdownOptions {
	return &GracefulShutdownOptions{
		ShutdownTimeout:       30 * time.Second,
		BeforeShutdownTimeout: 5 * time.Second,
	}
}

// RunWithGracefulShutdown 以优雅退出的方式运行 Kitex 服务器
// 该方法基于 Kitex 内置的优雅退出机制，提供增强功能：
// 1. 自动监听系统信号（SIGTERM, SIGINT）
// 2. 在停止前执行自定义钩子函数（BeforeShutdownFunc）
// 3. 调用 Kitex 的 Stop() 方法（会自动注销注册中心并执行 ShutdownHook）
// 4. 提供超时控制，防止无限等待
// 5. 支持自定义清理函数（通过 RegisterShutdownHook 注册）
//
// 注意：
//   - Kitex 的 Stop() 方法会自动注销注册中心（如果通过 WithRegistry 配置了）
//   - CleanupFunc 会通过 RegisterShutdownHook 注册，在 Stop() 时自动执行
//   - BeforeShutdownFunc 在 Stop() 之前执行，可以用于设置健康检查状态等
//
// 参数：
//   - svr: Kitex 服务器实例
//   - opts: 优雅退出配置选项，如果为 nil 则使用默认配置
//
// 示例：
//
//	svr := server.NewServer(
//		server.WithRegistry(nacosRegistry), // Kitex 会自动注销
//		// ... 其他选项
//	)
//	serversuite.RunWithGracefulShutdown(svr, &serversuite.GracefulShutdownOptions{
//		ShutdownTimeout: 30 * time.Second,
//		BeforeShutdownFunc: func() {
//			// 停止接收新请求、设置健康检查状态等
//			healthCheck.SetUnhealthy()
//		},
//		BeforeShutdownTimeout: 5 * time.Second,
//		CleanupFunc: func() error {
//			// 关闭数据库连接、清理资源等
//			// 注意：这个函数会通过 RegisterShutdownHook 注册
//			return db.Close()
//		},
//	})
func RunWithGracefulShutdown(svr server.Server, opts *GracefulShutdownOptions) {
	// 规范化配置选项
	opts = normalizeOptions(opts)

	// 注册清理函数到 Kitex 的 ShutdownHook（如果提供）
	if opts.CleanupFunc != nil {
		server.RegisterShutdownHook(func() {
			if err := opts.CleanupFunc(); err != nil {
				klog.Errorf("清理函数执行失败: %v", err)
			} else {
				klog.Infof("清理函数执行完成")
			}
		})
	}

	// 创建信号通道并配置 Kitex 的退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// 在 goroutine 中启动服务器
	errChan := make(chan error, 1)
	go func() {
		klog.Infof("Kitex 服务器正在启动...")
		if err := svr.Run(); err != nil {
			errChan <- err
		}
	}()

	// 等待服务器启动完成或出错
	select {
	case err := <-errChan:
		klog.Errorf("服务器启动失败: %v", err)
		os.Exit(1)
		return
	case <-time.After(100 * time.Millisecond):
		// 给服务器一点时间启动
		klog.Infof("Kitex 服务器已启动，等待退出信号...")
	}

	// 等待退出信号
	sig := <-sigChan
	klog.Infof("收到退出信号: %v，开始优雅退出...", sig)

	// 执行关闭前的钩子函数（在 Stop() 之前）
	executeBeforeShutdownFunc(opts)

	// 停止服务器（Kitex 会自动注销注册中心并执行 ShutdownHook）
	stopServerWithTimeout(svr, opts.ShutdownTimeout)

	klog.Infof("优雅退出完成")
}

// normalizeOptions 规范化配置选项，设置默认值
func normalizeOptions(opts *GracefulShutdownOptions) *GracefulShutdownOptions {
	if opts == nil {
		opts = defaultGracefulShutdownOptions()
	}
	if opts.ShutdownTimeout <= 0 {
		opts.ShutdownTimeout = 30 * time.Second
	}
	if opts.BeforeShutdownTimeout <= 0 {
		opts.BeforeShutdownTimeout = 5 * time.Second
	}
	return opts
}

// executeBeforeShutdownFunc 执行关闭前的钩子函数
func executeBeforeShutdownFunc(opts *GracefulShutdownOptions) {
	if opts.BeforeShutdownFunc == nil {
		return
	}

	klog.Infof("执行关闭前钩子函数...")
	beforeShutdownCtx, beforeShutdownCancel := context.WithTimeout(context.Background(), opts.BeforeShutdownTimeout)
	defer beforeShutdownCancel()

	beforeShutdownDone := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("关闭前钩子函数发生 panic: %v", r)
			}
			close(beforeShutdownDone)
		}()
		opts.BeforeShutdownFunc()
	}()

	select {
	case <-beforeShutdownDone:
		klog.Infof("关闭前钩子函数执行完成")
	case <-beforeShutdownCtx.Done():
		klog.Warnf("关闭前钩子函数执行超时（%v），继续执行关闭流程", opts.BeforeShutdownTimeout)
	}
}

// stopServerWithTimeout 停止服务器，带超时控制
// Kitex 的 Stop() 方法会自动：
// 1. 执行所有通过 RegisterShutdownHook 注册的钩子函数
// 2. 注销注册中心（如果通过 WithRegistry 配置了）
// 3. 停止服务器
func stopServerWithTimeout(svr server.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("停止服务器时发生 panic: %v", r)
			}
			close(done)
		}()
		klog.Infof("正在停止 Kitex 服务器（将自动注销注册中心并执行清理函数）...")
		if err := svr.Stop(); err != nil {
			klog.Errorf("停止服务器时出错: %v", err)
		} else {
			klog.Infof("Kitex 服务器已停止")
		}
	}()

	select {
	case <-done:
		// 服务器已停止，Kitex 已自动处理注册中心注销和 ShutdownHook
		klog.Infof("优雅退出流程完成")
	case <-ctx.Done():
		klog.Warnf("优雅退出超时（%v），但服务器可能仍在关闭中", timeout)
	}
}

// RunWithGracefulShutdownAndCleanup 以优雅退出的方式运行 Kitex 服务器，并执行清理函数
// 这是 RunWithGracefulShutdown 的简化版本，只需要提供清理函数
//
// 参数：
//   - svr: Kitex 服务器实例
//   - cleanupFunc: 清理函数，在服务器停止后执行
//
// 示例：
//
//	svr := server.NewServer(...)
//	serversuite.RunWithGracefulShutdownAndCleanup(svr, func() error {
//		return db.Close()
//	})
func RunWithGracefulShutdownAndCleanup(svr server.Server, cleanupFunc func() error) {
	RunWithGracefulShutdown(svr, &GracefulShutdownOptions{
		CleanupFunc: cleanupFunc,
	})
}

// RunWithGracefulShutdownSimple 以优雅退出的方式运行 Kitex 服务器（使用默认配置）
// 这是最简单的使用方式，不需要任何额外配置
//
// 参数：
//   - svr: Kitex 服务器实例
//
// 示例：
//
//	svr := server.NewServer(...)
//	serversuite.RunWithGracefulShutdownSimple(svr)
func RunWithGracefulShutdownSimple(svr server.Server) {
	RunWithGracefulShutdown(svr, nil)
}

// RunWithNacosGracefulShutdown 使用 NacosServerSuite 以优雅退出的方式运行 Kitex 服务器
// 这是专门为 NacosServerSuite 提供的便捷方法
//
// 注意：如果服务器是通过 NacosServerSuite.Options() 创建的，注册中心已经通过 WithRegistry 配置，
// Kitex 的 Stop() 方法会自动注销，无需手动处理。
//
// 参数：
//   - svr: Kitex 服务器实例（应该已经通过 NacosServerSuite.Options() 配置了注册中心）
//   - nacosSuite: NacosServerSuite 实例
//   - cleanupFunc: 可选的清理函数
//
// 示例：
//
//	nacosSuite := serversuite.NacosServerSuite{
//		CurrentServiceName: "your-service",
//		RegistryAddr:      "127.0.0.1:8848",
//		// ... 其他配置
//	}
//	opts := nacosSuite.Options() // 这里已经配置了注册中心
//	svr := server.NewServer(opts...)
//	serversuite.RunWithNacosGracefulShutdown(svr, nacosSuite, func() error {
//		return db.Close()
//	})
func RunWithNacosGracefulShutdown(svr server.Server, nacosSuite NacosServerSuite, cleanupFunc func() error) {
	// 直接使用优雅退出，Kitex 会自动处理注册中心注销
	RunWithGracefulShutdown(svr, &GracefulShutdownOptions{
		CleanupFunc: cleanupFunc,
	})
}
