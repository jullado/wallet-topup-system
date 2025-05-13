package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type appLogs struct {
	log *zap.Logger
}

func NewAppZapLogs() AppLog {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// ข้ามตำแหน่งการเรียกใช้ logger
	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return appLogs{log: log}
}

func (l appLogs) Info(msg string) {
	l.log.Info(msg)
}

func (l appLogs) Debug(msg string) {
	l.log.Debug(msg)
}

func (l appLogs) Warning(msg string) {
	l.log.Warn(msg)
}

func (l appLogs) Error(msg interface{}) {
	switch v := msg.(type) {
	case error:
		l.log.Error(v.Error())
	case string:
		l.log.Error(v)
	}
}
