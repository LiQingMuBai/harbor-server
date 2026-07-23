package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

	ServiceStartupEvent(
		"service_dir", logDirName,
		"service_tag", logFileTag,
		"pid", os.Getpid(),
		"go_version", runtime.Version(),
		"info_log", ServiceLogPath(dailyInfoLogFilenameAt(logFileTag, now)),
		"access_log", ServiceLogPath(dailyAccessLogFilenameAt(logFileTag, now)),
		"panic_log", ServiceLogPath(dailyPanicLogFilenameAt(logFileTag, now)),
	)

	log.SetFlags(0)
	return nil
}

func ServiceStartupEvent(keysAndValues ...interface{}) {
	ServiceLogger().Sugar().Infow("startup", keysAndValues...)
}

func ServiceStartupBanner(title string, keysAndValues ...interface{}) {
	eventFields := make([]interface{}, 0, len(keysAndValues)+2)
	eventFields = append(eventFields, "title", title)
	eventFields = append(eventFields, keysAndValues...)
	ServiceStartupEvent(eventFields...)

	sugar := ServiceLogger().Sugar()
	boxWidth := 86
	innerWidth := boxWidth - 4
	avatarWidth := 18

	top := "+" + strings.Repeat("-", boxWidth-2) + "+"
	bottom := top
	serviceID := startupBannerServiceID(keysAndValues...)
	avatarLines := startupAvatarLines(serviceID)

	kv := formatKeyValues(keysAndValues...)
	textLines := make([]string, 0, len(avatarLines))
	textLines = append(textLines, "START "+title)
	textLines = append(textLines, wrapText(kv, innerWidth-avatarWidth-1, 4)...)

	rowCount := len(avatarLines)
	if len(textLines) > rowCount {
		rowCount = len(textLines)
	}

	sugar.Info(top)
	for i := 0; i < rowCount; i++ {
		left := ""
		if i < len(avatarLines) {
			left = avatarLines[i]
		}
		right := ""
		if i < len(textLines) {
			right = textLines[i]
		}
		sugar.Info("| " + padRight(joinBannerColumns(left, right, avatarWidth, innerWidth), innerWidth) + " |")
	}
	sugar.Info(bottom)
}

func startupBannerServiceID(keysAndValues ...interface{}) string {
	service := ""
	mode := ""
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		key := fmt.Sprint(keysAndValues[i])
		value := fmt.Sprint(keysAndValues[i+1])
		switch key {
		case "service":
			service = strings.TrimSpace(value)
		case "mode":
			mode = strings.TrimSpace(value)
		}
	}
	if service == "task" && mode != "" {
		return service + "/" + mode
	}
	if service != "" {
		return service
	}
	return strings.TrimSpace(GetServiceLogName())
}

func startupAvatarLines(serviceID string) []string {
	switch serviceID {
	case "api":
		return []string{
			"    _/\\\\/\\\\_       ",
			"   /  o  o \\\\      ",
			"  |    --   |      ",
			"  |  \\\\__/  |      ",
			"   \\\\__==__/       ",
		}
	case "admin":
		return []string{
			"    .-^^^^-.      ",
			"   /  -  - \\\\     ",
			"  |   [__]  |     ",
			"  |  \\\\____/ |     ",
			"   '-.==.-'      ",
		}
	case "wss":
		return []string{
			"    _/~~~~\\\\_     ",
			"   /  0  0  \\\\    ",
			"  |    <>    |    ",
			"  |  \\\\____/  |    ",
			"   \\\\__==__//     ",
		}
	case "task/data":
		return []string{
			"    _/::::\\\\_     ",
			"   /  o  .. \\\\    ",
			"  |   <==>   |    ",
			"  |  \\\\____/  |    ",
			"   \\\\__==__//     ",
		}
	case "task/jobs":
		return []string{
			"    _/####\\\\_     ",
			"   /  -  oo \\\\    ",
			"  |   [==]   |    ",
			"  |  \\\\____/  |    ",
			"   \\\\__==__//     ",
		}
	case "cdn":
		return []string{
			"    _/@@@@\\\\_     ",
			"   /  ^  ^  \\\\    ",
			"  |   \\__/   |    ",
			"  |  /____\\  |    ",
			"   \\\\__==__//     ",
		}
	default:
		return []string{
			"    _/----\\\\_     ",
			"   /  .  .  \\\\    ",
			"  |    --    |    ",
			"  |  \\\\____/  |    ",
			"   \\\\__==__//     ",
		}
	}
}

func joinBannerColumns(left string, right string, leftWidth int, totalWidth int) string {
	left = padRight(left, leftWidth)
	rightWidth := totalWidth - leftWidth - 1
	if rightWidth < 0 {
		rightWidth = 0
	}
	return left + " " + padRight(right, rightWidth)
}

func formatKeyValues(keysAndValues ...interface{}) string {
	if len(keysAndValues) == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < len(keysAndValues); i += 2 {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i+1 >= len(keysAndValues) {
			b.WriteString(fmt.Sprint(keysAndValues[i]))
			break
		}
		b.WriteString(fmt.Sprint(keysAndValues[i]))
		b.WriteByte('=')
		b.WriteString(fmt.Sprint(keysAndValues[i+1]))
	}
	return b.String()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func wrapText(s string, width int, maxLines int) []string {
	s = strings.TrimSpace(s)
	if s == "" || width <= 0 || maxLines <= 0 {
		return []string{""}
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, maxLines)
	var current strings.Builder

	flush := func() {
		if current.Len() == 0 {
			return
		}
		lines = append(lines, current.String())
		current.Reset()
	}

	for _, w := range words {
		if len(lines) >= maxLines {
			break
		}
		if len(w) > width {
			w = w[:width]
		}

		if current.Len() == 0 {
			current.WriteString(w)
			continue
		}
		if current.Len()+1+len(w) <= width {
			current.WriteByte(' ')
			current.WriteString(w)
			continue
		}

		flush()
		if len(lines) >= maxLines {
			break
		}
		current.WriteString(w)
	}
	flush()

	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	if len(lines) == maxLines && len(words) > 0 {
		last := lines[maxLines-1]
		if len(last) >= 3 && last[len(last)-3:] != "..." && len(words) > 0 {
			if len(last) > width-3 {
				last = last[:width-3]
			}
			lines[maxLines-1] = strings.TrimRight(last, " ") + "..."
		}
	}

	return lines
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
