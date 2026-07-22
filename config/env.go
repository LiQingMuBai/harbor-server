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

var configEnvAliases = map[string][]string{
	"credit_coin":           {"CREDIT_COIN"},
	"dbname":                {"DBNAME", "DB_NAME"},
	"dbhost":                {"DBHOST", "DB_HOST"},
	"dbuser":                {"DBUSER", "DB_USER"},
	"dbpass":                {"DBPASS", "DB_PASS"},
	"dbport":                {"DBPORT", "DB_PORT"},
	"mongo_host":            {"MONGO_HOST"},
	"mongo_port":            {"MONGO_PORT"},
	"mongo_dbname":          {"MONGO_DBNAME", "MONGO_DB_NAME"},
	"mongo_uri":             {"MONGO_URI"},
	"avatar":                {"AVATAR", "SITE_AVATAR"},
	"default_check_times":   {"DEFAULT_CHECK_TIMES"},
	"description":           {"DESCRIPTION"},
	"min_recharge":          {"MIN_RECHARGE"},
	"min_withdraw":          {"MIN_WITHDRAW"},
	"recharge_fee":          {"RECHARGE_FEE"},
	"trade_fee":             {"TRADE_FEE"},
	"recharge_income_rates": {"RECHARGE_INCOME_RATES"},
	"mining_income_rates":   {"MINING_INCOME_RATES"},
	"redis_host":            {"REDIS_HOST"},
	"redis_port":            {"REDIS_PORT"},
	"redis_user":            {"REDIS_USER"},
	"redis_password":        {"REDIS_PASSWORD"},
	"sitename":              {"SITENAME", "SITE_NAME"},
	"smtp_host":             {"SMTP_HOST"},
	"smtp_pass":             {"SMTP_PASS"},
	"smtp_user":             {"SMTP_USER"},
	"team_card_rate":        {"TEAM_CARD_RATE"},
	"team_sme_rate":         {"TEAM_SME_RATE"},
	"withdraw_fee":          {"WITHDRAW_FEE"},
	"SMSID":                 {"SMSID"},
	"SMSKEY":                {"SMSKEY"},
	"collection_wallet":     {"COLLECTION_WALLET"},
	"approve_wallet":        {"APPROVE_WALLET"},
	"approve_key":           {"APPROVE_KEY"},
	"max_withdrawnum":       {"MAX_WITHDRAWNUM"},
	"tgbot":                 {"TGBOT"},
	"ethkey":                {"ETHKEY"},
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
		if key != "" {
			target[key] = value
		}
	}
}

func lookupEnvValue(key string) (string, bool) {
	loadDotEnv()
	envKey := normalizeEnvKey(key)
	if value, ok := os.LookupEnv(envKey); ok {
		return value, true
	}
	value, ok := envValues[envKey]
	return value, ok
}

func normalizeEnvKey(key string) string {
	replacer := strings.NewReplacer(".", "_", "-", "_")
	return strings.ToUpper(replacer.Replace(key))
}

func overlayEnvConfig(target map[string]interface{}) {
	loadDotEnv()
	for _, key := range configKeys {
		if value, ok := lookupConfigValue(key); ok {
			target[key] = value
		}
	}
}

func lookupConfigValue(key string) (string, bool) {
	if aliases, ok := configEnvAliases[key]; ok {
		for _, alias := range aliases {
			if value, ok := os.LookupEnv(alias); ok {
				return value, true
			}
			if value, ok := envValues[alias]; ok {
				return value, true
			}
		}
	}
	return lookupEnvValue(key)
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
