package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/stulzq/hexo-deploy-agent/config"
	"github.com/stulzq/hexo-deploy-agent/util"
)

// Logger is the interface for Logger types
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

var (
	innerLogger Logger
	loggerLock  = &sync.Mutex{}
)

func init() {
	logLevel := config.GetString("log:level")

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}
	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	zapLoggerEncoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "message",
		StacktraceKey:    "stacktrace",
		EncodeCaller:     customCallerEncoder,
		EncodeTime:       customTimeEncoder,
		EncodeLevel:      customLevelEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		LineEnding:       "\n",
		ConsoleSeparator: " ",
	}

	var level zap.AtomicLevel
	var syncWriter zapcore.WriteSyncer
	if util.IsDebug() {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapLoggerEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		syncWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
		syncWriter = &zapcore.BufferedWriteSyncer{
			WS: zapcore.AddSync(&lumberjack.Logger{
				Filename:  filepath.Join("logs", "app.log"), // ⽇志⽂件路径
				MaxSize:   50,                               // 单位为MB,默认为512MB
				MaxAge:    5,                                // 文件最多保存多少天
				LocalTime: true,                             // 采用本地时间
				Compress:  false,                            // 是否压缩日志
			}),
			Size: 4096,
		}
	}

	if logLevel != "" {
		l := zap.AtomicLevel{}
		if err := l.UnmarshalText([]byte(logLevel)); err == nil {
			level = l
		}
	}

	zapCore := zapcore.NewCore(zapcore.NewConsoleEncoder(zapLoggerEncoderConfig), syncWriter, level)
	zapLogger := zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(1))

	SetLogger(zapLogger.Sugar())

	fmt.Printf("hexo[logger]: Use log level \x1b[0;32m%s\x1b[0m\n", strings.ToUpper(level.String()))
}

func SetLogger(logger Logger) {
	innerLogger = logger
}

func GetLogger() Logger {
	return innerLogger
}
