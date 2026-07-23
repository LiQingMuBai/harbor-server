package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	serviceLogMu   sync.RWMutex
	serviceLogName = "default"
	runtimeLogFile *os.File
	serviceLogger  *zap.Logger
	stdLogUndo     func()
)

func normalizeServiceLogName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\\", "/")
	name = filepath.Clean(name)
	name = strings.TrimPrefix(name, "/")
	name = strings.TrimPrefix(name, "./")
	if name == "" || name == "." {
		return "default"
	}
	if strings.HasPrefix(name, "..") {
		return "default"
	}
	return name
}

func SetServiceLogName(name string) {
	serviceLogMu.Lock()
	defer serviceLogMu.Unlock()
	serviceLogName = normalizeServiceLogName(name)
}

func GetServiceLogName() string {
	serviceLogMu.RLock()
	defer serviceLogMu.RUnlock()
	return serviceLogName
}

func ServiceLogDir() string {
	return filepath.Join("logs", filepath.FromSlash(GetServiceLogName()))
}

func EnsureServiceLogDir() (string, error) {
	dir := ServiceLogDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func ServiceLogPath(filename string) string {
	dir, err := EnsureServiceLogDir()
	if err != nil {
		return filepath.Join("logs", "default", filename)
	}
	return filepath.Join(dir, filename)
}

func OpenServiceLog(filename string) (*os.File, error) {
	return os.OpenFile(ServiceLogPath(filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func SetupServiceLogger(name string) error {
	SetServiceLogName(name)
	file, err := OpenServiceLog("runtime.log")
	if err != nil {
		return err
	}

	serviceLogMu.Lock()
	if serviceLogger != nil {
		_ = serviceLogger.Sync()
	}
	if stdLogUndo != nil {
		stdLogUndo()
		stdLogUndo = nil
	}
	if runtimeLogFile != nil {
		_ = runtimeLogFile.Close()
	}
	runtimeLogFile = file
	serviceLogger = newServiceLogger(file)
	stdLogUndo = zap.RedirectStdLog(serviceLogger)
	serviceLogMu.Unlock()

	log.SetFlags(0)
	return nil
}

func newServiceLogger(file *os.File) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(file),
	)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zap.InfoLevel,
	)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func ServiceLogger() *zap.Logger {
	serviceLogMu.RLock()
	logger := serviceLogger
	serviceLogMu.RUnlock()
	if logger != nil {
		return logger
	}
	return zap.NewNop()
}

func ServiceInfo(args ...interface{}) {
	ServiceLogger().Sugar().Info(args...)
}

func ServiceWarn(args ...interface{}) {
	ServiceLogger().Sugar().Warn(args...)
}

func ServiceError(args ...interface{}) {
	ServiceLogger().Sugar().Error(args...)
}

func ServiceInfof(template string, args ...interface{}) {
	ServiceLogger().Sugar().Infof(template, args...)
}

func ServiceWarnf(template string, args ...interface{}) {
	ServiceLogger().Sugar().Warnf(template, args...)
}

func ServiceErrorf(template string, args ...interface{}) {
	ServiceLogger().Sugar().Errorf(template, args...)
}

func AppendServiceLog(filename string, text string) error {
	file, err := OpenServiceLog(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	_, err = file.WriteString(text)
	return err
}
