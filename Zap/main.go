// https://github.com/uber-go/zap
package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 支持7中日志级别：Debug Info Warn Error DPanic Panic Fatal
// DPanic 是指在开发环境下记录日志后会进行 panic

func main() {
	Example4()
}

func Example1() {
	// 将日志输出到控制台
	log, _ := zap.NewProduction()

	log.Debug("This is a DEBUG message") // 在Production下，该日志不会输出
	log.Info("This is a INFO message", zap.String("url", "http://127.0.0.1"), zap.Int("attempt", 3), zap.Duration("backoff", time.Second))
	log.Warn("This is a WARN message")
	log.Error("This is a ERROR message")
	log.Fatal("This is a FATAL message")
	fmt.Println("")
}

// ====================================================================================================
var log *zap.Logger
var logSugared *zap.SugaredLogger

func Example2() {
	encoder := getEncoder()
	writerSyncer := getWriterSyncer()
	core := zapcore.NewCore(encoder, writerSyncer, zap.DebugLevel)
	log = zap.New(core, zap.AddCallerSkip(0), zap.Fields(zap.String("xxx", "yyy")))
	logSugared = log.Sugar()

	log.Info("this is a info message", zap.String("url", "http://127.0.0.1"))
	logSugared.Infof("this is a info message, url: %s", "http://127.0.0.1")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncer() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}

// ====================================================================================================
func Example3() {
	logger, _ := newCustomLogger()
	defer logger.Sync()

	// 增加一个 skip 选项，触发 zap 内部 error，将错误输出到 error.log
	// logger = logger.WithOptions(zap.AddCallerSkip(100))

	logger.Info("Info msg")
	logger.Error("Error msg")
}

func newCustomLogger() (*zap.Logger, error) {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "./log/test.log"},
		ErrorOutputPaths: []string{"./log/error.log"},
	}
	return cfg.Build(zap.Fields(zap.String("xxx", "yyy")), zap.AddCallerSkip(1), zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
}

// ====================================================================================================
var logX *zap.Logger

func Example4() {
	InitLoggerX()

	for i := 0; i < 100000; i++ {
		logX.Info("this is a info message")
	}
}

func InitLoggerX() {
	encoder := getEncoderX()

	// 对日志级别进行判断
	highPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zap.ErrorLevel && l >= zap.DebugLevel
	})

	infoFileWriteSyncer := getInfoWriterSyncerX()
	errorFileWriteSyncer := getErrorWriterSyncerX()

	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)
	// core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writerSyncer), zap.DebugLevel)
	cores := []zapcore.Core{infoFileCore, errorFileCore}

	logX = zap.New(zapcore.NewTee(cores...), zap.Fields(zap.String("serverName", "awesome web")), zap.AddCallerSkip(1), zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
}

func getEncoderX() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncerX() zapcore.WriteSyncer {
	// file, _ := os.OpenFile("./log/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// return zapcore.AddSync(file)

	// 引入第三方库 Lumberjack 加入日志切割功能
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./log/log.log", // 日志文件名
		MaxSize:    10,              // 日志文件大小
		MaxBackups: 100,             // 最大备份文件数
		MaxAge:     30,              // 日志保存最大天数
		Compress:   false,           // 不压缩
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberWriteSyncer)
}

// 记录 error 以下日志级别的文件
func getInfoWriterSyncerX() zapcore.WriteSyncer {
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./log/info.log", // 日志文件名
		MaxSize:    1,                // 日志文件大小
		MaxBackups: 100,              // 最大备份文件数
		MaxAge:     30,               // 日志保存最大天数
		Compress:   false,            // 不压缩
	}

	return zapcore.AddSync(lumberWriteSyncer)
}

// 记录 error 及以上日志级别的文件
func getErrorWriterSyncerX() zapcore.WriteSyncer {
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./log/error.log", // 日志文件名
		MaxSize:    10,                // 日志文件大小
		MaxBackups: 100,               // 最大备份文件数
		MaxAge:     30,                // 日志保存最大天数
		Compress:   false,             // 不压缩
	}

	return zapcore.AddSync(lumberWriteSyncer)
}
