package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	envOnce   sync.Once
	envValues map[string]string
)

var configKeys = []string{
	"credit_coin",
	"mysql_dsn",
	"dbname",
	"dbhost",
	"dbuser",
	"dbpass",
	"dbport",
	"mongo_host",
	"mongo_port",
	"mongo_dbname",
	"mongo_uri",
	"avatar",
	"default_check_times",
	"description",
	"min_recharge",
	"min_withdraw",
	"recharge_fee",
	"trade_fee",
	"recharge_income_rates",
	"mining_income_rates",
	"redis_host",
	"redis_port",
	"redis_user",
	"redis_password",
	"sitename",
	"smtp_host",
	"smtp_pass",
	"smtp_user",
	"team_card_rate",
	"team_sme_rate",
	"withdraw_fee",
	"SMSID",
	"SMSKEY",
	"collection_wallet",
	"approve_wallet",
	"approve_key",
	"max_withdrawnum",
	"tgbot",
	"ethkey",
}

func loadDotEnv() {
	envOnce.Do(func() {
		envValues = make(map[string]string)
		for _, path := range findEnvCandidates() {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			parseDotEnv(content, envValues)
			for key, value := range envValues {
				if _, ok := os.LookupEnv(key); ok {
					continue
				}
				_ = os.Setenv(key, value)
			}
			return
		}
	})
}

func EnsureDotEnvLoaded() {
	loadDotEnv()
}

func findEnvCandidates() []string {
	candidates := []string{
		".env",
		"../.env",
		"../../.env",
	}
	if envPath := strings.TrimSpace(os.Getenv("APP_ENV_FILE")); envPath != "" {
		candidates = append([]string{envPath}, candidates...)
	}
	return candidates
}

func parseDotEnv(content []byte, target map[string]string) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		index := strings.Index(line, "=")
		if index <= 0 {
			continue
		}
		key := normalizeEnvKey(strings.TrimSpace(line[:index]))
		value := strings.TrimSpace(line[index+1:])
		value = strings.Trim(value, `"'`)
		value = stripInlineComment(value)
		if key != "" {
			target[key] = value
		}
	}
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

func lookupEnvValue(key string) (string, bool) {
	loadDotEnv()
	envKey := normalizeEnvKey(key)
	if value, ok := os.LookupEnv(envKey); ok {
		return stripInlineComment(value), true
	}
	value, ok := envValues[envKey]
	if !ok {
		return "", false
	}
	return stripInlineComment(value), true
}

func normalizeEnvKey(key string) string {
	replacer := strings.NewReplacer(".", "_", "-", "_")
	return strings.ToUpper(replacer.Replace(key))
}

func overlayEnvConfig(target map[string]interface{}) {
	loadDotEnv()
	for _, key := range configKeys {
		if value, ok := lookupEnvValue(key); ok {
			target[key] = value
		}
	}
}

func resolveConfigFile(filename string) (string, error) {
	dirCandidates := []string{
		"config",
		"../config",
		"../../config",
	}
	if configDir := strings.TrimSpace(os.Getenv("CONFIG_DIR")); configDir != "" {
		dirCandidates = append([]string{configDir}, dirCandidates...)
	}
	for _, dir := range dirCandidates {
		fullPath := filepath.Join(dir, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("config file not found: %s", filename)
}
