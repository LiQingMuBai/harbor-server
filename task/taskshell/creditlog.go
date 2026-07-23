package taskshell

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"strings"
)

func CreditLog() {
	//用户账变扫描线程
	go UserLevels() //层级统计
	go Team()       //团队统计
	for {
		list := config.GlobalRedis.PopQueue(models.QUEUE_USER_COIN_LOG)
		if len(list) > 0 {
			for _, v := range list {
				mp := config.GlobalConfig.GetConfigFromJson(v)
				if mp == nil {
					continue
				}
				if mp.GetValue("type") == nil {
					continue
				}
				t := mp.GetValue("type").ToInt()
				uid := mp.GetValue("uid").ToInt()
				data := mp.GetValue("data").ToConfig()
				credit := data.GetValue("credit").ToFloat()
				lockcredit := data.GetValue("lockcredit").ToFloat()
				sn := data.GetValue("sn").ToString()
				createtime := data.GetValue("createtime").ToInt()
				cointype := data.GetValue("cointype").ToString()
				if cointype == "" {
					cointype = "usdt"
				}
				if credit == 0 && lockcredit == 0 {
					continue
				}
				models.MODEL_CREDIT_LOG.Add(&models.CreditLogInfo{
					Uid:        uid,
					Credit:     credit,
					LockCredit: lockcredit,
					Mode:       1,
					Sn:         sn,
					Type:       t,
					Createtime: createtime,
					CoinType:   cointype,
				})
			}
		}
	}
}
func UserLevels() {
	//用户层级
	for {
		list := config.GlobalRedis.PopQueue(models.QUEUE_USER_REGISTER)
		if len(list) > 0 {
			ntime := utils.GetNow()
			for _, v := range list {
				models.MODEL_CREDIT_LOG.AddSiteCount(ntime, map[string]float64{"register_num": 1}) //增加注册数
				mp := config.GlobalConfig.GetConfigFromJson(v)
				uid := mp.GetValue("uid").ToInt()
				models.MODEL_ASSETS.InitUserAssets(uid) //初始化用户资产
				invite_order := mp.GetValue("invite_order").ToString()
				if invite_order == "" {
					continue
				}
				channel_id := mp.GetValue("channel_id")
				if channel_id != nil {
					models.MODEL_CREDIT_LOG.AddAgentLog(channel_id.ToInt(), ntime, map[string]float64{"register_num": 1})
					models.MODEL_CREDIT_LOG.AddAgentLevelLog(channel_id.ToInt(), mp.GetValue("channel_level").ToInt(), ntime, map[string]float64{"register_num": 1})
				}
				tmp := strings.Split(invite_order, ",")
				n := len(tmp)
				for _, v := range tmp {
					ctmp := tmp[len(tmp)-n:]
					config.GlobalDB.InsertData(models.DB_TABLE_USER_LEVELS, db.DB_PARAMS{"puid": v, "uid": uid, "level": n, "levle_order": strings.Join(ctmp, ","), "createtime": ntime})
					values := map[string]float64{"register_num": 1}
					if n == 0 {
						values["direct_register_num"] = 1
					}
					models.MODEL_CREDIT_LOG.AddLevelCount(utils.GetInt(v), n, ntime, map[string]float64{"register_num": 1}) //增加层级统计的注册人数
					models.MODEL_CREDIT_LOG.AddUserCount(utils.GetInt(v), ntime, values)                                    //增加团队统计的注册人数

					n--
				}
			}
		}
	}
}
func Team() {
	for {
		list := config.GlobalRedis.PopQueue(models.QUEUE_TEAM_COIN_LOG)
		if len(list) > 0 {
			for _, v := range list {
				mp := config.GlobalConfig.GetConfigFromJson(v)
				uid := mp.GetValue("uid").ToInt()
				data := utils.GetJsonValue(mp.GetValue("data").Value)
				utils.ServiceInfo("team credit log payload:", data)
				logInfo := new(models.QueueTeamLog)
				err := json.Unmarshal([]byte(data), logInfo)
				if err == nil {
					models.MODEL_CREDIT_LOG.AddUserCountLog(uid, logInfo)
				} else {
					utils.ServiceError("team credit log unmarshal failed:", err)
				}
			}
		}
	}
}
