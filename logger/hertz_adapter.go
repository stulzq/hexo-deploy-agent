package logger

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
)

type HertzFullLogger struct {
}

func NewHertzFullLogger() hlog.FullLogger {
	return &HertzFullLogger{}
}

func (l *HertzFullLogger) Trace(v ...interface{}) {
	GetLogger().Debug(v...)
}

func (l *HertzFullLogger) Debug(v ...interface{}) {
	GetLogger().Debug(v...)
}

func (l *HertzFullLogger) Info(v ...interface{}) {
	GetLogger().Info(v...)
}

func (l *HertzFullLogger) Notice(v ...interface{}) {
	GetLogger().Info(v...)
}

func (l *HertzFullLogger) Warn(v ...interface{}) {
	GetLogger().Warn(v...)
}

func (l *HertzFullLogger) Error(v ...interface{}) {
	GetLogger().Error(v...)
}

func (l *HertzFullLogger) Fatal(v ...interface{}) {
	GetLogger().Fatal(v...)
}

func (l *HertzFullLogger) Tracef(format string, v ...interface{}) {
	GetLogger().Debugf(format, v...)
}

func (l *HertzFullLogger) Debugf(format string, v ...interface{}) {
	GetLogger().Debugf(format, v...)
}

func (l *HertzFullLogger) Infof(format string, v ...interface{}) {
	GetLogger().Infof(format, v...)
}

func (l *HertzFullLogger) Noticef(format string, v ...interface{}) {
	GetLogger().Infof(format, v...)
}

func (l *HertzFullLogger) Warnf(format string, v ...interface{}) {
	GetLogger().Warnf(format, v...)
}

func (l *HertzFullLogger) Errorf(format string, v ...interface{}) {
	GetLogger().Errorf(format, v...)
}

func (l *HertzFullLogger) Fatalf(format string, v ...interface{}) {
	GetLogger().Fatalf(format, v...)
}

func (l *HertzFullLogger) CtxTracef(_ context.Context, format string, v ...interface{}) {
	GetLogger().Debugf(format, v...)
}

func (l *HertzFullLogger) CtxDebugf(_ context.Context, format string, v ...interface{}) {
	GetLogger().Debugf(format, v...)
}

func (l *HertzFullLogger) CtxInfof(_ context.Context, format string, v ...interface{}) {
	GetLogger().Infof(format, v...)
}

func (l *HertzFullLogger) CtxNoticef(_ context.Context, format string, v ...interface{}) {
	GetLogger().Infof(format, v...)
}

func (l *HertzFullLogger) CtxWarnf(_ context.Context, format string, v ...interface{}) {
	GetLogger().Warnf(format, v...)
}

func (l *HertzFullLogger) CtxErrorf(_ context.Context, format string, v ...interface{}) {
	GetLogger().Errorf(format, v...)
}

func (l *HertzFullLogger) CtxFatalf(_ context.Context, format string, v ...interface{}) {
	GetLogger().Fatalf(format, v...)
}

func (l *HertzFullLogger) SetLevel(_ hlog.Level) {

}
func (l *HertzFullLogger) SetOutput(_ io.Writer) {

}
