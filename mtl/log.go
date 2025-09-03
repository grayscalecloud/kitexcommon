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

package mtl

import (
	"io"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/grayscalecloud/kitexcommon/logger"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(ioWriter io.Writer, rootPath string, useTrace bool) {
	var opts []kitexzap.Option
	var output zapcore.WriteSyncer

	if os.Getenv("GO_ENV") == "dev" {
		// 开发环境使用 ConsoleEncoder 并输出到控制台不输出到文件
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		consoleOutput := zapcore.AddSync(os.Stdout)
		output = zapcore.NewMultiWriteSyncer(consoleOutput)
	} else {
		// 生产环境使用 JSONEncoder，并输出到控制台和带缓冲的文件
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		consoleOutput := zapcore.AddSync(os.Stdout)
		fileOutput := &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ioWriter),
			FlushInterval: time.Minute,
		}
		output = zapcore.NewMultiWriteSyncer(consoleOutput, fileOutput)
	}

	server.RegisterShutdownHook(func() {
		output.Sync() //nolint:errcheck
	})

	if useTrace {
		opts = append(opts, kitexzap.WithRecordStackTraceInSpan(true))
		log := logger.NewTraceLogger(rootPath)
		klog.SetLogger(log)
	} else {
		log := kitexzap.NewLogger(opts...)
		klog.SetLogger(log)
	}
	klog.SetLevel(klog.LevelTrace)
	klog.SetOutput(output)
}
