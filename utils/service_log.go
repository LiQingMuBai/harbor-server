package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	serviceLogMu   sync.RWMutex
	serviceLogName = "default"
	runtimeLogFile *os.File
	serviceLogger  *zap.Logger
	stdLogUndo     func()

	serviceInfoLogDate string
	serviceInfoLogTag  string
	serviceRotateOnce  sync.Once
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

func dailyInfoLogFilenameAt(serviceName string, now time.Time) string {
	name := strings.ReplaceAll(filepath.ToSlash(serviceName), "/", "-")
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		name = "default"
	}
	return name + ".info." + now.Format("2006-01-02") + ".log"
}

func SetupServiceLogger(name string) error {
	serviceName := normalizeServiceLogName(name)

	logDirName := serviceName
	logFileTag := serviceName
	if strings.HasPrefix(serviceName, "task/") {
		logDirName = "task"
		mode := strings.TrimPrefix(serviceName, "task/")
		mode = strings.TrimSpace(mode)
		if mode != "" {
			logFileTag = "task-" + mode
		} else {
			logFileTag = "task"
		}
	}

	SetServiceLogName(logDirName)
	now := time.Now()
	file, err := OpenServiceLog(dailyInfoLogFilenameAt(logFileTag, now))
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
	serviceInfoLogDate = now.Format("2006-01-02")
	serviceInfoLogTag = logFileTag
	serviceLogMu.Unlock()

	serviceRotateOnce.Do(func() {
		go rotateServiceInfoLogLoop()
	})

	log.SetFlags(0)
	return nil
}

func rotateServiceInfoLogLoop() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		sleepFor := time.Until(next)
		if sleepFor < 0 {
			sleepFor = time.Second
		}
		time.Sleep(sleepFor)
		rotateServiceInfoLogIfNeeded(time.Now())
	}
}

func rotateServiceInfoLogIfNeeded(now time.Time) {
	date := now.Format("2006-01-02")

	serviceLogMu.Lock()
	if serviceLogger == nil || serviceInfoLogDate == date {
		serviceLogMu.Unlock()
		return
	}
	serviceTag := serviceInfoLogTag
	if strings.TrimSpace(serviceTag) == "" {
		serviceTag = serviceLogName
	}
	serviceLogMu.Unlock()

	file, err := OpenServiceLog(dailyInfoLogFilenameAt(serviceTag, now))
	if err != nil {
		return
	}

	serviceLogMu.Lock()
	if serviceLogger == nil || serviceInfoLogDate == date {
		serviceLogMu.Unlock()
		_ = file.Close()
		return
	}

	oldFile := runtimeLogFile
	oldLogger := serviceLogger
	oldUndo := stdLogUndo

	runtimeLogFile = file
	serviceLogger = newServiceLogger(file)
	stdLogUndo = zap.RedirectStdLog(serviceLogger)
	serviceInfoLogDate = date
	serviceLogMu.Unlock()

	if oldUndo != nil {
		oldUndo()
	}
	if oldLogger != nil {
		_ = oldLogger.Sync()
	}
	if oldFile != nil {
		_ = oldFile.Close()
	}
}

type dailyFileWriter struct {
	mu         sync.Mutex
	date       string
	file       *os.File
	filenameAt func(now time.Time) string
}

func newDailyFileWriter(filenameAt func(now time.Time) string) (*dailyFileWriter, error) {
	w := &dailyFileWriter{filenameAt: filenameAt}
	if err := w.rotate(time.Now()); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *dailyFileWriter) rotate(now time.Time) error {
	date := now.Format("2006-01-02")
	if w.file != nil && w.date == date {
		return nil
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	file, err := OpenServiceLog(w.filenameAt(now))
	if err != nil {
		w.file = nil
		w.date = ""
		return err
	}
	w.file = file
	w.date = date
	return nil
}

func (w *dailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotate(time.Now()); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func dailyAccessLogFilenameAt(serviceName string, now time.Time) string {
	name := strings.ReplaceAll(filepath.ToSlash(serviceName), "/", "-")
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		name = "default"
	}
	return name + ".access." + now.Format("2006-01-02") + ".log"
}

func dailyPanicLogFilenameAt(serviceName string, now time.Time) string {
	name := strings.ReplaceAll(filepath.ToSlash(serviceName), "/", "-")
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		name = "default"
	}
	return name + ".panic." + now.Format("2006-01-02") + ".log"
}

func NewServiceAccessLogWriter() (io.Writer, error) {
	serviceLogMu.RLock()
	tag := strings.TrimSpace(serviceInfoLogTag)
	if tag == "" {
		tag = serviceLogName
	}
	serviceLogMu.RUnlock()
	return newDailyFileWriter(func(now time.Time) string {
		return dailyAccessLogFilenameAt(tag, now)
	})
}

func AppendServicePanicLog(text string) error {
	serviceLogMu.RLock()
	tag := strings.TrimSpace(serviceInfoLogTag)
	if tag == "" {
		tag = serviceLogName
	}
	serviceLogMu.RUnlock()

	file, err := OpenServiceLog(dailyPanicLogFilenameAt(tag, time.Now()))
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
