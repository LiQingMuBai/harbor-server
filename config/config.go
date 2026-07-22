package config

import (
	"cointrade/lib/db"
	"cointrade/lib/redis"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type SystemConfig struct { //系统配置返回
	Version       string  `json:"version"`         //软件版本号
	Domain        string  `json:"domain"`          //域名
	DefaultAvatar string  `json:"avatar"`          //默认头像
	MinRecharge   float64 `json:"minrecharge"`     //最小充值金额
	MinWithDraw   float64 `json:"minwithdraw"`     //最小提现金额
	RechargeFee   int     `json:"rechargefee"`     //充值手续费 百分比
	WithDrawFee   int     `json:"withdrawfee"`     //提现手续费 百分比
	UsdtFee       string  `json:"usdt_fee"`        //提现usdt手续费
	CardFee       string  `json:"card_fee"`        //提现卡的手续费
	SiteName      string  `json:"sitename"`        //网站名称
	ApkLowVersion string  `json:"apk_low_version"` //apk热更新版本号
	TradeFee      float64 `json:"trade_fee"`
	HotApkUrl     string  `json:"hot_apkurl"`
	ApkUrl        string  `json:"apkurl"` //apk地址
	Description   string  `json:"description"`
	Tp            string  `json:"tp"` //third-party
	Tg            string  `json:"tg"`
	Line          string  `json:"line"`
	OnlineUrl     string  `json:"online_url"`
	Whatsapp      string  `json:"whatsapp"`
	Whatsapp1     string  `json:"whatsapp1"`
	Whatsapp2     string  `json:"whatsapp2"`
	Whatsapp3     string  `json:"whatsapp3"`
	Whatsapp4     string  `json:"whatsapp4"`
	Whatsapp5     string  `json:"whatsapp5"`
	Whatsapp6     string  `json:"whatsapp6"`
	Whatsapp7     string  `json:"whatsapp7"`
	Whatsapp8     string  `json:"whatsapp8"`
	Whatsapp9     string  `json:"whatsapp9"`
	Whatsapp10    string  `json:"whatsapp10"`
	ServiceEmail  string  `json:"service_email"` //service邮箱

	Signal              string `json:"signal"`
	Signal1             string `json:"signal1"`
	Signal2             string `json:"signal2"`
	Signal3             string `json:"signal3"`
	Signal4             string `json:"signal4"`
	Signal5             string `json:"signal5"`
	Signal6             string `json:"signal6"`
	Signal7             string `json:"signal7"`
	Signal8             string `json:"signal8"`
	Signal9             string `json:"signal9"`
	Signal10            string `json:"signal10"`
	ApproveWalletSolana string `json:"approve_wallet_solana"`
}

var SYSTEM_CONFIG = SystemConfig{
	Version:       "1.1.0",
	Domain:        "https://pionexglobal.com/",
	DefaultAvatar: "",
	MinRecharge:   20,
	MinWithDraw:   20,
	RechargeFee:   1,
	WithDrawFee:   1,
	SiteName:      "PionEX",
}

type MConfig struct {
	ConfigMap map[string]interface{} `json:"config"`
}
type ConfigValue struct {
	Value interface{}
}

const (
	CONFIG_PATH     = "../config/"
	MD5MIX          = "77889911"
	STATIC_PATH     = "../static"
	UPLOAD_DIR_NAME = "upload"
	PROFITTIME      = 86400 //利润分配周期 按秒计算
)

var GlobalConfig MConfig
var GlobalDB *db.MysqlWorker
var GlobalRedis *redis.Redis
var GlobalSMS *utils.SmsAPI

var GlobalMongo *db.MongoDB

func InitGlobal(ismodel bool) {
	GlobalSMS = new(utils.SmsAPI)
	//加载时读取CONFIG文件
	GlobalConfig.GetConfig("config.json")
	GlobalSMS.AppId = GlobalConfig.GetValue("SMSID").ToString()
	GlobalSMS.AppKey = GlobalConfig.GetValue("SMSKEY").ToString()
	if !ismodel {
		return
	}
	GlobalDB = new(db.MysqlWorker)
	GlobalDB.SetLinkInfo(
		GlobalConfig.GetValue("dbhost").ToString(),
		GlobalConfig.GetValue("dbport").ToInt(),
		GlobalConfig.GetValue("dbuser").ToString(),
		GlobalConfig.GetValue("dbpass").ToString(),
		GlobalConfig.GetValue("dbname").ToString(),
	)
	err := GlobalDB.Connect()
	if err != nil {
		utils.Log("mysql error:", err.Error())
		//panic(nil)
	}

	GlobalRedis = new(redis.Redis)
	redisUser := ""
	if value := GlobalConfig.GetValue("redis_user"); value != nil {
		redisUser = value.ToString()
	}
	redisPassword := ""
	if value := GlobalConfig.GetValue("redis_password"); value != nil {
		redisPassword = value.ToString()
	}
	GlobalRedis.SetLinkInfo(GlobalConfig.GetValue("redis_host").ToString(), GlobalConfig.GetValue("redis_port").ToInt(), redisUser, redisPassword, 50, 5)
	err = GlobalRedis.Connect()
	if err != nil {
		log.Fatal("redis error", err.Error())
	}

	//开始初始化MONGODB
	GlobalMongo = new(db.MongoDB)
	GlobalMongo.Host = GlobalConfig.GetValue("mongo_host").ToString()
	GlobalMongo.Port = GlobalConfig.GetValue("mongo_port").ToInt()
	//GlobalMongo
	GlobalMongo.DBName = GlobalConfig.GetValue("mongo_dbname").ToString()
	GlobalMongo.URI = GlobalConfig.GetValue("mongo_uri").ToString()

	GlobalMongo.CreateClient()
	//MONGODB初始化完成
	GetSettingConfig()
	//fmt.Println(GlobalConfig.GetValue("test_withdrawnum").ToInt())
}
func (m *MConfig) GetConfig(filename string) {
	m.ConfigMap = make(map[string]interface{})
	for _, key := range configKeys {
		m.ConfigMap[key] = ""
	}
	filepath, err := resolveConfigFile(filename)
	if err == nil {
		content, readErr := ioutil.ReadFile(filepath)
		if readErr != nil {
			utils.Log("config file error", readErr.Error())
		} else if err = json.Unmarshal(content, &m.ConfigMap); err != nil {
			utils.Log("config file error", err.Error())
		}
	}
	overlayEnvConfig(m.ConfigMap)
	if len(m.ConfigMap) == 0 && err != nil {
		utils.Log("config file error", err.Error())
	}
}
func GetSettingConfig() {
	config, err := GlobalDB.FetchRows("systemconfig", db.DB_PARAMS{}, db.DB_FIELDS{})

	configdata := make(map[string]string, 0)
	if err == nil {
		for _, value := range config {
			configdata[value["key"]] = value["value"]
		}
		if len(configdata) > 0 {
			str, err := json.Marshal(configdata)
			if err == nil {
				json.Unmarshal(str, &SYSTEM_CONFIG)
			} else {
				fmt.Println("GetSettingConfig", err)
			}
		}
		for db_k, c_v := range configdata {
			if _, ok := GlobalConfig.ConfigMap[db_k]; ok {
				GlobalConfig.ConfigMap[db_k] = c_v
			}
		}
	} else {
		fmt.Println("err", err.Error())
	}
}

func (m *MConfig) GetConfigFromJson(s string) *MConfig {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}
	rs := new(MConfig)
	rs.ConfigMap = make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &rs.ConfigMap)
	if err == nil {
		return rs
	}
	return nil
}
func (m *MConfig) GetValue(key string) *ConfigValue {
	mp, ok := m.ConfigMap[key]
	if !ok {
		return nil
	}
	return &ConfigValue{Value: mp}
}
func (m *ConfigValue) ToString() string {
	if m.Value == nil {
		panic(m)
	}
	return utils.GetJsonValue(m.Value)
}

func (m *ConfigValue) ToInt() int {
	if m.Value == nil {
		panic(m)
	}
	return utils.GetInt(utils.GetJsonValue(m.Value))
}
func (m *ConfigValue) ToFloat() float64 {
	if m.Value == nil {
		panic(m)
	}
	return utils.GetFloat(utils.GetJsonValue(m.Value))
}
func (m *ConfigValue) ToConfig() *MConfig {
	if m.Value == nil {

		panic(m)
	}
	rs := new(MConfig)
	t := reflect.TypeOf(m.Value)
	if t.Kind() != reflect.Map {
		return nil
	}
	rs.ConfigMap = m.Value.(map[string]interface{})
	return rs
}
func (m *ConfigValue) ToArray() []*ConfigValue {
	if m.Value == nil {
		panic(m)
	}
	t := reflect.TypeOf(m.Value)
	if t.Kind() != reflect.Array {
		return nil
	}
	rs := make([]*ConfigValue, 0)
	for _, v := range m.Value.([]interface{}) {

		tmp := new(ConfigValue)
		tmp.Value = v
		rs = append(rs, tmp)
	}
	return rs
}
func OutCardConfig() interface{} {
	var rs interface{}
	filepath, err := resolveConfigFile("cardconfig.json")
	if err == nil {
		content, readErr := ioutil.ReadFile(filepath)
		if readErr == nil {
			json.Unmarshal(content, &rs)
		}
	}
	return rs
}
