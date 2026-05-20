package loginfra

import (
	"context"
	"os"

	"github.com/datpham2001/mb-user-service/internal/infrastructure/configinfra"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	logger *logrus.Logger
}

func New(cfg *configinfra.Config) *Logger {
	l := &Logger{
		logger: logrus.New(),
	}

	l.init(cfg.Server.Env)
	return l
}

func (l *Logger) init(env string) {
	l.logger.SetOutput(os.Stdout)

	if env == "production" {
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		l.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	switch env {
	case "production":
		l.logger.SetLevel(logrus.InfoLevel)
	default:
		l.logger.SetLevel(logrus.DebugLevel)
	}
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, args ...any) {
	l.withContext(ctx).Infof(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(args...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, args ...any) {
	l.withContext(ctx).Errorf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) withContext(ctx context.Context) *logrus.Entry {
	entry := l.logger.WithContext(ctx)

	if requestID, ok := ctx.Value("request_id").(string); ok {
		entry = entry.WithField("request_id", requestID)
	}

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		entry = entry.WithField("trace_id", span.SpanContext().TraceID().String())
		entry = entry.WithField("span_id", span.SpanContext().SpanID().String())
	}

	return entry
}

func (l *Logger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...any) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.logger.Panicf(format, args...)
}
