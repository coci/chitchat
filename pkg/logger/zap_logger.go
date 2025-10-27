package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	logger, _ := zap.NewProduction()
	return &ZapLogger{
		logger: logger,
	}
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Info(msg, zapFields...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Warn(msg, zapFields...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Error(msg, zapFields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	z.logger.Fatal(msg, zapFields...)
}
