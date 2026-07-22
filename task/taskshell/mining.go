package taskshell

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
	"time"
)

func Minging() {
	//用户挖矿利润分配
	go MiningR() //预约订单扫描
	for {
		ntime := utils.GetNow()
		list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "_": fmt.Sprintf("profittime<=%d", ntime), "unlocktime": 0}, db.DB_FIELDS{}, "limit 0,500")
		for _, v := range list {
			if v["profittimes"].ToInt() >= v["circle"].ToInt() {
				//config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 1, "unlocktime": ntime}, db.DB_PARAMS{"id": v["id"].Value})
				models.MODEL_PRODUCT.Unlock(v["uid"].ToInt(), v["sn"].ToString())
				/*models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{ //体验机
					Credit:          v["allprofit"].ToFloat(),
					LockCredit:      0,
					VCrdit:          0,
					LockVCredit:     0,
					UserCoinLogType: models.COIN_LOG_USER_MINING_PROFIT,
					UserCoinLogInfo: models.QueueCreditLog{
						Credit:     v["allprofit"].ToFloat(),
						LockCredit: 0,
						Sn:         v["sn"].ToString(),
						CreateTime: ntime,
					},
					TeamCoinLogType: models.TEAM_LOG_MINING_PROFIT,
					TeamCoinLogInfo: models.QueueTeamLog{
						MiningProfit: v["allprofit"].ToFloat(),
						CreateTime:   ntime,
					},
				})*/
				continue
			}
			//累计利润
			uinfo := models.MODEL_USER.GetBaseInfo(v["uid"].ToInt())
			if uinfo == nil {
				continue
			}
			profit := 0.0
			amount := v["amount"].ToFloat()
			if v["rate"].ToFloat() > 0 {
				if v["dispatch_amount"].ToFloat() > 0 {
					amount = v["dispatch_amount"].ToFloat()
				}
				diffrand := int((v["rate"].ToFloat() - v["rate_min"].ToFloat()) * 100000)
				profit = amount * (v["rate_min"].ToFloat() + float64(diffrand)/100000)
			} else {
				profit = v["profit"].ToFloat()
			}

			if uinfo.ParentOrder != "" {
				//开始团队返利
				list := strings.Split(uinfo.ParentOrder, ",")
				n := len(list)
				for _, vv := range list {
					if rates, ok := models.MINING_INCOME_RATES[n]; ok {
						m_profit_income := profit * rates[0]
						models.MODEL_CREDIT_LOG.AddLevelCount(utils.GetInt(vv), n, ntime, map[string]float64{"mining_income": m_profit_income})
						config.GlobalDB.AddValue(models.DB_TABLE_USER, map[string]float64{"mining_income": m_profit_income}, db.DB_PARAMS{"id": utils.GetInt(vv)})
						config.GlobalDB.InsertData(models.DB_TABLE_INCOME_LOG, db.DB_PARAMS{
							"uid":         utils.GetInt(vv),
							"income":      m_profit_income,
							"type":        models.INCOME_TYPE_MINING_PROFIT,
							"child_uid":   uinfo.Id,
							"child_email": uinfo.Email,
							"createtime":  ntime,
							"level":       n,
						})
					}
					n--
				}
			}
			config.GlobalDB.AddValue(models.DB_TABLE_MINING_ORDER, map[string]float64{"allprofit": profit, "profittime": models.CIRCLE_TIME, "profittimes": 1}, db.DB_PARAMS{"id": v["id"].Value})
			models.MODEL_USER.AddCredit(uinfo.Id, &models.CreditValue{
				Credit:          profit,
				LockCredit:      0,
				UserCoinLogType: models.COIN_LOG_USER_MINING_PROFIT,
				UserCoinLogInfo: models.QueueCreditLog{
					Credit:     profit,
					CreateTime: ntime,
				},
				TeamCoinLogType: models.TEAM_LOG_MINING_PROFIT,
				TeamCoinLogInfo: models.QueueTeamLog{
					MiningProfit: profit,
					CreateTime:   ntime,
				},
			})
			if v["profittimes"].ToInt()+1 >= v["circle"].ToInt() {
				models.MODEL_PRODUCT.Unlock(v["uid"].ToInt(), v["sn"].ToString())
			}
		}
		time.Sleep(1 * time.Second)
	}

}
func MiningR() {
	//挖矿预约扫描线程
	for {
		ntime := utils.GetNow()
		list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 4}, db.DB_FIELDS{}, "limit 0,1000")
		for _, v := range list {
			if v["expiredtime"].ToInt() == 0 || v["expiredtime"].ToInt() > ntime {
				//这里开始判断余额是否满足分配的金额
				uinfo := models.MODEL_USER.GetBaseInfo(v["uid"].ToInt())
				/*if uinfo.Credit >= v["dispatch_amount"].ToFloat() && v["dispatch_amount"].ToFloat() > 0 {
					//这时候预约成功 修改订单状态 扣除用户余额
					if models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
						Credit:          -1 * v["dispatch_amount"].ToFloat(),
						UserCoinLogType: models.COIN_LOG_USER_BUY_MINING,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     -1 * v["dispatch_amount"].ToFloat(),
							CreateTime: ntime,
							CoinType:   "usdt",
							Sn:         v["sn"].ToString(),
						},
					}) {
						config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "amount": v["dispatch_amount"].ToFloat(), "profittime": ntime + models.CIRCLE_TIME}, db.DB_PARAMS{"id": v["id"].Value}) //修改订单状态为可执行
					}
				}*/
				diff_amount := v["dispatch_amount"].ToFloat() - v["recive_amount"].ToFloat()
				if diff_amount <= 0 {
					if v["state"].ToInt() == 4 {
						config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "amount": v["dispatch_amount"].ToFloat(), "profittime": ntime + models.CIRCLE_TIME}, db.DB_PARAMS{"id": v["id"].Value})
					}
					continue
				}
				if uinfo.Credit > diff_amount {
					if models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
						Credit:          -1 * diff_amount,
						UserCoinLogType: models.COIN_LOG_USER_REVERVATION,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     -1 * diff_amount,
							CreateTime: ntime,
							CoinType:   "usdt",
							Sn:         v["sn"].ToString(),
						},
					}) {
						config.GlobalDB.AddValue(models.DB_TABLE_MINING_ORDER, map[string]float64{"recive_amount": diff_amount}, db.DB_PARAMS{"id": v["id"].Value})
					}

				} else {
					if models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
						Credit:          -1 * uinfo.Credit,
						UserCoinLogType: models.COIN_LOG_USER_REVERVATION,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     -1 * uinfo.Credit,
							CreateTime: ntime,
							CoinType:   "usdt",
							Sn:         v["sn"].ToString(),
						},
					}) {
						config.GlobalDB.AddValue(models.DB_TABLE_MINING_ORDER, map[string]float64{"recive_amount": uinfo.Credit}, db.DB_PARAMS{"id": v["id"].Value})
					}
				}
			} else {
				//修改订单状态为已过期
				config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 3}, db.DB_PARAMS{"id": v["id"].Value})
			}
		}
		time.Sleep(10 * time.Second)
	}
}
