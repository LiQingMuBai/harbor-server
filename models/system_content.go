package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
)

func (m *SystemModel) RuleText(ruleType string, lang string) db.DB_LIST_RESULT {
	cacheID := m.MakeCacheId(ruleType, lang)
	one := make(db.DB_LIST_RESULT, 0)
	err := config.GlobalRedis.GetObject(HASH_RULE_TEXT, cacheID, &one)
	if err == nil {
		return one
	}
	one, _ = config.GlobalDB.FetchRows(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"rule_type": ruleType, "lang": lang}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(HASH_RULE_TEXT, cacheID, one)
	}
	return one
}

func (m *SystemModel) CoinDesc(symbol string, lang string) map[string]string {
	cacheID := m.MakeCacheId(symbol, lang)
	cache := make(map[string]string, 0)
	config.GlobalRedis.GetObject(HASH_RULE_TEXT, cacheID, &cache)
	if len(cache) > 0 {
		return cache
	}
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_COIN_DESC, db.DB_PARAMS{"lang": lang, "symbol": symbol}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(DB_TABLE_COIN_DESC, cacheID, one)
	}
	return one
}

func (m *SystemModel) RuleOne(data map[string]interface{}) db.DB_ROW_RESULT {
	one := make(db.DB_ROW_RESULT, 0)
	where := ""
	if id, ok := data["id"]; ok && utils.GetInt(fmt.Sprintf("%v", id)) > 0 {
		where = fmt.Sprintf("id = '%v'", id)
	} else {
		lang := "en"
		if _, ok := data["lang"]; ok {
			lang = fmt.Sprintf("%v", data["lang"])
		}
		if _, ok := data["rule"]; !ok {
			return nil
		}
		where = fmt.Sprintf("lang = '%s' AND rule_type = '%s'", lang, data["rule"])
		if lang == "th" {
			one, _ = config.GlobalDB.FetchRow(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"_": where}, db.DB_FIELDS{})
			return one
		}
	}
	cacheID := m.MakeCacheId("rule", utils.Md5(where))
	err := config.GlobalRedis.GetObject(HASH_RULE_TEXT, cacheID, &one)
	if err == nil {
		return one
	}
	one, _ = config.GlobalDB.FetchRow(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"_": where}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(HASH_RULE_TEXT, cacheID, one)
	}
	return one
}
