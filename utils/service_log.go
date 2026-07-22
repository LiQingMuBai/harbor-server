package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	serviceLogMu   sync.RWMutex
	serviceLogName = "default"
	runtimeLogFile *os.File
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
	if runtimeLogFile != nil {
		_ = runtimeLogFile.Close()
	}
	runtimeLogFile = file
	serviceLogMu.Unlock()

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	return nil
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
