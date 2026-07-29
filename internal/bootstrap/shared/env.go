package shared

import (
	"cointrade/config"
	"os"
	"strconv"
	"strings"
	"sync"
)

var envOnce sync.Once

func ensureEnvLoaded() {
	envOnce.Do(func() {
		config.EnsureDotEnvLoaded()
	})
}

func Getenv(key string, defaultValue string) string {
	ensureEnvLoaded()
	value := strings.TrimSpace(stripInlineComment(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}
	return value
}

func GetenvInt(key string, defaultValue int) int {
	ensureEnvLoaded()
	value := strings.TrimSpace(stripInlineComment(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return number
}

func stripInlineComment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	for i := 0; i < len(value); i++ {
		if value[i] != '#' {
			continue
		}
		if i == 0 {
			return ""
		}
		prev := value[i-1]
		if prev == ' ' || prev == '\t' {
			return strings.TrimSpace(value[:i])
		}
	}
	return value
}

func ParseRPCClients(value string) map[int]string {
	result := map[int]string{}
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		result[port] = strings.TrimSpace(parts[1])
	}
	return result
}
