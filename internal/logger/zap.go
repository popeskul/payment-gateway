package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zap: zapLogger}, nil
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Debugw(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Warnw(msg, keysAndValues...)
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Fatalw(msg, keysAndValues...)
}
