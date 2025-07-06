package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Global logger
var Log *zap.Logger

// func InitLogger() {
// 	// 设置日志级别
// 	level := zapcore.InfoLevel
//
// 	// 配置输出格式：console（开发）或 json（生产）
// 	encoderCfg := zap.NewDevelopmentEncoderConfig()
// 	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
// 	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色输出
// 	encoderCfg.CallerKey = "caller"
// 	encoderCfg.FunctionKey = "func"
//
// 	// 日志核心配置
// 	core := zapcore.NewCore(
// 		zapcore.NewConsoleEncoder(encoderCfg), // 编码器
// 		zapcore.AddSync(os.Stdout),            // 输出目标
// 		level,
// 	)
//
// 	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
// 	zap.ReplaceGlobals(Log) // 替换全局 zap.L()
// }

func InitLogger() *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()

	// 彩色时间（蓝色）
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("\x1b[34m%s\x1b[0m", t.Format("2006-01-02 15:04:05")))
	}

	// 彩色 caller（青色）
	encoderCfg.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("\x1b[36m%s\x1b[0m", caller.TrimmedPath()))
	}

	// 彩色 level（黄色/红色等）
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoderCfg.FunctionKey = "func"
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
	return logger
}

