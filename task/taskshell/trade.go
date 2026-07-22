package taskshell

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func ClearDelegateTrade() { //开始清理委托单
	//go ClearExplodeTrade() //清理交割合约
	for {
		start := 1
		limit := 500 //一次处理500单
		offset := (start - 1) * limit
		limitstr := fmt.Sprintf("limit %d,%d", offset, limit)
		list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 0, "is_f": 0}, db.DB_FIELDS{}, limitstr)
		if len(list) > 0 {
			for _, v := range list {
				OpDelegateTrade(v)
			}
		}

		//start++

		time.Sleep(1000 * time.Millisecond)
	}
}

func OpDelegateTrade(one db.DBValues) bool {
	//操作委托单
	ntime := utils.GetNow()

	uid := one["uid"].ToInt()
	if ntime-one["createtime"].ToInt() >= 7*24*60*60 {
		models.MODEL_TRADE.CancleDelegate(uid, one["sn"].ToString())
		return false
	}
	dtype := one["trade_type"].ToInt()
	mode := one["mode"].ToInt()
	coin := one["coin_symbol"].ToString()
	//coinPriceInfo := config.GlobalMongo.GetOne("lastdata", bson.M{"pair": one["coinpair"].ToString()}, nil)
	coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(one["coinpair"].ToString())
	if coinPriceInfo["close"] == nil {
		return false
	}
	coinprice := coinPriceInfo["close"].(float64)
	v_credit := 0.0      //虚拟额度
	credit := 0.0        //真实额度
	lock_v_credit := 0.0 //冻结虚拟额度
	lock_credit := 0.0   //冻结真实额度

	if mode == models.USER_MODE_REAL {
		credit = one["credit"].ToFloat()
		lock_credit = one["credit"].ToFloat() + one["fee"].ToFloat()
	} else {
		v_credit = one["credit"].ToFloat()
		lock_v_credit = one["credit"].ToFloat() + one["fee"].ToFloat()
	}
	switch dtype {
	case models.OPEN_TYPE_BB: //币币交易
		if one["delegate_type"].ToInt() == models.DELEGATE_TYPE_BUY {
			//买单
			if coinprice <= one["price"].ToFloat() {
				//当前价格到达委托价格时 委托成功 减去用户冻结余额 增加用户资产
				if (models.MODEL_ASSETS.AddAssets(uid, &models.Assets{
					Coin:    coin,
					Pair:    one["coinpair"].ToString(),
					Num:     one["num"].ToFloat(),
					LockNum: 0,
					Price:   one["price"].ToFloat(),
					Mode:    mode,
				})) {
					models.MODEL_USER.AddCredit(uid, &models.CreditValue{ //用户余额处理
						Credit:          0,
						LockCredit:      -1 * lock_credit,
						VCrdit:          0,
						LockVCredit:     -1 * lock_v_credit,
						UserCoinLogType: models.COIN_LOG_BB_TRADE,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     one["num"].ToFloat(),
							CoinType:   one["coin_symbol"].ToString(),
							CreateTime: ntime,
						},
						TeamCoinLogType: models.TEAM_LOG_TRADE,
						TeamCoinLogInfo: models.QueueTeamLog{
							TradeBB:    lock_credit,
							CreateTime: ntime,
						},
					})

					config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": one["id"].ToInt()}) //更改委托状态
					return true
				}

			}
		} else {
			//卖单 减去用户冻结资产 加入用户余额
			if coinprice >= one["price"].ToFloat() {
				assetInfo := models.MODEL_ASSETS.GetAllAssets(uid, one["mode"].ToInt())

				bb_profit := credit - one["num"].ToFloat()*assetInfo[one["coin_symbol"].ToString()].O_Price
				if (models.MODEL_ASSETS.AddAssets(uid, &models.Assets{
					Coin:    coin,
					Pair:    one["coinpair"].ToString(),
					Num:     0,
					LockNum: -1 * one["num"].ToFloat(),
					Price:   -1 * one["price"].ToFloat(),
					Mode:    mode,
				})) {
					models.MODEL_USER.AddCredit(uid, &models.CreditValue{ //用户余额处理 加入用户余额
						Credit:          credit,
						LockCredit:      0,
						VCrdit:          v_credit,
						LockVCredit:     0,
						UserCoinLogType: models.COIN_LOG_USER_CLOSE,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     credit,
							LockCredit: 0,
							Sn:         one["sn"].ToString(),
							CreateTime: ntime,
						},
						TeamCoinLogType: models.TEAM_LOG_TRADE_PROFIT,
						TeamCoinLogInfo: models.QueueTeamLog{
							CreateTime:     ntime,
							TradeBB_Profit: bb_profit,
						},
					})
					models.MODEL_USER.AddCredit(uid, &models.CreditValue{ //用户余额处理 加入用户余额
						Credit:          0,
						LockCredit:      0,
						VCrdit:          0,
						LockVCredit:     0,
						UserCoinLogType: models.COIN_LOG_USER_CLOSE,
						UserCoinLogInfo: models.QueueCreditLog{
							Credit:     -1 * one["num"].ToFloat(),
							Sn:         one["sn"].ToString(),
							CoinType:   one["coin_symbol"].ToString(),
							CreateTime: ntime,
						},
					})
					config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": one["id"].ToInt()}) //更改委托状态到已完成
					return true
				}
			}
		} //币币交易处理完成
	case models.OPEN_TYPE_EXPLODE: //交割合约处理 这里只有买单处理 系统自动交割处理持仓
		if (one["flag"].ToInt() == models.DIRECT_TYPE_BIG && coinprice <= one["price"].ToFloat()) || (one["flag"].ToInt() == models.DIRECT_TYPE_SMALL && coinprice >= one["price"].ToFloat()) {
			//价格达到临界点
			explodeConfig, ok := models.EXPLODE_CONFIG[one["close_time"].ToInt()]
			if !ok {
				return false
			}
			insertData := db.DB_PARAMS{}
			insertData["uid"] = one["uid"].ToInt()
			insertData["sn"] = one["sn"].ToString()
			insertData["trade_type"] = models.OPEN_TYPE_EXPLODE
			insertData["flag"] = one["flag"].ToInt()
			insertData["openprice"] = one["price"].ToFloat()
			insertData["closeprice"] = 0
			insertData["coinid"] = one["coinid"].ToInt()
			insertData["coinpair"] = one["coinpair"].ToString()
			insertData["coin_symbol"] = one["coin_symbol"].ToString()
			insertData["close_time"] = one["close_time"].ToInt()
			insertData["close_real_time"] = ntime + one["close_time"].ToInt()
			insertData["clear_time"] = 0
			insertData["createtime"] = ntime
			insertData["ganggan"] = one["ganggan"].ToInt()
			insertData["credit"] = one["credit"].ToFloat()
			insertData["profit"] = 0
			insertData["win_rate"] = explodeConfig.Winrate
			insertData["lose_rate"] = explodeConfig.Loserate
			insertData["num"] = one["num"].ToFloat()
			insertData["mode"] = mode
			_, err := config.GlobalDB.InsertData(models.DB_TABLE_OPENED_TRADE, insertData) //增加交割合约的持仓
			if err == nil {
				if models.MODEL_USER.AddCredit(uid, &models.CreditValue{
					Credit:          0,
					LockCredit:      -1 * lock_credit,
					VCrdit:          0,
					LockVCredit:     -1 * lock_v_credit,
					UserCoinLogType: 0,
					UserCoinLogInfo: nil,
					TeamCoinLogType: models.TEAM_LOG_TRADE,
					TeamCoinLogInfo: models.QueueTeamLog{
						TradeExplode: lock_credit,
						CreateTime:   ntime,
					},
				}) {
					config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": one["id"].ToInt()})
					return true

				}

			}
		}
	case models.OPEN_TYPE_KEEP: //永续合约处理
		if one["delegate_type"].ToInt() == models.DELEGATE_TYPE_BUY {
			//买单处理
			fmt.Println("coinprice:", coinprice)
			if (one["flag"].ToInt() == models.DIRECT_TYPE_BIG && coinprice <= one["price"].ToFloat()) || (one["flag"].ToInt() == models.DIRECT_TYPE_SMALL && coinprice >= one["price"].ToFloat()) {
				//价格到达了临界点
				if models.MODEL_USER.AddCredit(uid, &models.CreditValue{
					Credit:          0,
					LockCredit:      -1 * lock_credit,
					VCrdit:          0,
					LockVCredit:     -1 * lock_v_credit,
					UserCoinLogType: 0,
					UserCoinLogInfo: nil,
					TeamCoinLogType: models.TEAM_LOG_TRADE,
					TeamCoinLogInfo: models.QueueTeamLog{
						TradeKeep:  lock_credit,
						CreateTime: ntime,
					},
				}) {
					models.MODEL_TRADE.AddKeepOpend(one)
					return true
				}

			}
		} else {
			//处理卖单

			if (one["flag"].ToInt() == models.DIRECT_TYPE_BIG && coinprice >= one["price"].ToFloat()) || (one["flag"].ToInt() == models.DIRECT_TYPE_SMALL && coinprice <= one["price"].ToFloat()) {
				//价格到达临界点 开始处理卖单
				//开始减仓
				var appendinfo *models.OpenedInfo
				if one["ganggan_sn"].ToString() != "0" && one["ganggan"].ToInt() > 1 {
					appendinfo = models.MODEL_TRADE.GetOpendBySn(one["uid"].ToInt(), one["ganggan_sn"].ToString())
				} else {
					appendinfo = models.MODEL_TRADE.GetOpendOne(one["uid"].ToInt(), coin, one["trade_type"].ToInt(), one["flag"].ToInt(), mode, one["ganggan"].ToInt()) //取得持仓信息
				}

				if appendinfo == nil {
					return false
				}
				profit := ((one["price"].ToFloat() - appendinfo.OpenPrice) / appendinfo.OpenPrice) * one["num"].ToFloat() * 1000 //计算出利益卖多的利益
				if one["flag"].ToInt() != models.DIRECT_TYPE_BIG {
					//卖空
					profit = ((one["num"].ToFloat()*1000)/one["price"].ToFloat())*appendinfo.OpenPrice - one["num"].ToFloat()*1000 //计算卖空的利益
				}
				profit = profit * float64(appendinfo.Ganggan)
				all_real_credit := profit + one["num"].ToFloat()*1000
				insertData := db.DB_PARAMS{}
				insertData["uid"] = uid
				insertData["sn"] = one["sn"].ToString()
				insertData["coin_symbol"] = one["coin_symbol"].ToString()
				insertData["trade_type"] = one["trade_type"].ToInt()
				insertData["flag"] = one["flag"].ToInt()
				insertData["amount"] = one["num"].ToFloat()
				insertData["close_price"] = one["price"].ToFloat()
				insertData["createtime"] = ntime
				insertData["num"] = one["num"].ToFloat()
				insertData["mode"] = mode
				insertData["allprice"] = all_real_credit
				insertData["profit"] = profit
				insertData["o_price"] = appendinfo.OpenPrice
				_, err := config.GlobalDB.InsertData(models.DB_TABLE_CLOSE_TRADE, insertData) //插入平仓表
				if err == nil {
					//这里要剪掉持仓数量
					config.GlobalDB.AddValue(models.DB_TABLE_OPENED_TRADE, map[string]float64{"lock_num": -1 * one["num"].ToFloat(), "credit": -1 * one["num"].ToFloat() * 1000}, db.DB_PARAMS{"id": appendinfo.Id})
					if mode == models.USER_MODE_REAL {
						//真实盘下
						models.MODEL_USER.AddCredit(uid, &models.CreditValue{
							Credit:          all_real_credit,
							LockCredit:      0,
							VCrdit:          0,
							LockVCredit:     0,
							UserCoinLogType: models.COIN_LOG_USER_CLOSE,
							UserCoinLogInfo: models.QueueCreditLog{
								Credit:     all_real_credit,
								Sn:         one["sn"].ToString(),
								CreateTime: ntime,
							},
							TeamCoinLogType: models.TEAM_LOG_TRADE_PROFIT,
							TeamCoinLogInfo: models.QueueTeamLog{
								TradeKeep_Profit: profit,
								CreateTime:       ntime,
							},
						})
					} else {
						//虚拟盘下
						models.MODEL_USER.AddCredit(uid, &models.CreditValue{
							Credit:          0,
							LockCredit:      0,
							VCrdit:          all_real_credit,
							LockVCredit:     0,
							UserCoinLogType: 0,
							UserCoinLogInfo: nil,
							TeamCoinLogType: 0,
							TeamCoinLogInfo: nil,
						})
					}
					config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": one["id"].ToInt()}) //修改委托状态
					return true
				}

			}
		}
	}
	return false
}
func EqualsPrice(myprice, coinprice float64, isbig bool) bool {
	if isbig {
		return myprice <= coinprice
	} else {
		return myprice >= coinprice
	}

}
func CheckStop(v db.DBValues, coinprice float64, isbig bool) bool {
	if v["stop_up_price"].ToFloat() > 0 && EqualsPrice(v["stop_up_price"].ToFloat(), coinprice, isbig) { //触发止盈
		delegate_price := v["stop_up_delegate"].ToFloat()
		if delegate_price == 0 {
			delegate_price = coinprice
		}
		if v["stop_up_delegate"].ToFloat() > 0 {
			rs := models.MODEL_TRADE.DelegateTrade(v["uid"].ToInt(), &models.TradeDelegateRequest{
				OpenType:     models.OPEN_TYPE_KEEP,
				DelegateType: models.DELEGATE_TYPE_SELL,
				Pair:         v["coinpair"].ToString(),
				Coin:         v["coin_symbol"].ToString(),
				PriceType:    models.PRICE_TYPE_LIMIT,
				GangGan:      v["ganggan"].ToInt(),
				Amount:       v["num"].ToFloat(),
				Price:        delegate_price,
				Sn:           v["sn"].ToString(),
			})
			if rs.State == models.STATE_SUCCESS {
				config.GlobalDB.UpdateData(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"auto_delegate": 1}, db.DB_PARAMS{"id": v["id"].Value})
			}
			return true
		}
	}

	if v["stop_down_price"].ToFloat() > 0 && EqualsPrice(v["stop_up_price"].ToFloat(), coinprice, !isbig) { //触发止亏
		delegate_price := v["stop_down_delegate"].ToFloat()
		if delegate_price == 0 {
			delegate_price = coinprice
		}
		if v["stop_down_delegate"].ToFloat() > 0 {
			rs := models.MODEL_TRADE.DelegateTrade(v["uid"].ToInt(), &models.TradeDelegateRequest{
				OpenType:     models.OPEN_TYPE_KEEP,
				DelegateType: models.DELEGATE_TYPE_SELL,
				Pair:         v["coinpair"].ToString(),
				Coin:         v["coin_symbol"].ToString(),
				PriceType:    models.PRICE_TYPE_LIMIT,
				GangGan:      v["ganggan"].ToInt(),
				Amount:       v["num"].ToFloat(),
				Price:        delegate_price,
				Sn:           v["sn"].ToString(),
			})
			if rs.State == models.STATE_SUCCESS {
				config.GlobalDB.UpdateData(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"auto_delegate": 1}, db.DB_PARAMS{"id": v["id"].Value})
			}
			return true
		}
	}
	return false
}
func ClearKeepCross() {
	//扫描永续杠杆穿仓的订单 然后强制平仓 亏损大于投资金额的80%强行平仓
	for {
		limit := 1000 //一次处理1000张订单
		condition := db.DB_PARAMS{"trade_type": models.OPEN_TYPE_KEEP, "_": "num>0 and ganggan>1 and auto_delegate=0"}
		count := config.GlobalDB.GetCount(models.DB_TABLE_OPENED_TRADE, condition)
		pagesize := int(math.Ceil(float64(count) / float64(limit)))
		for i := 1; i <= pagesize; i++ {
			list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_OPENED_TRADE, condition, db.DB_FIELDS{}, fmt.Sprintf("limit %d,%d", (i-1)*limit, limit))
			for _, v := range list {
				ntime := utils.GetNow()
				coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(v["coinpair"].ToString())
				coinprice := coinPriceInfo["close"].(float64)
				profit := 0.0
				if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
					//做多
					if CheckStop(v, coinprice, true) {
						continue
					}
					profit = ((coinprice - v["openprice"].ToFloat()) / v["openprice"].ToFloat()) * v["num"].ToFloat() * 1000 //计算出利益卖多的利益
				} else {
					//做空
					if CheckStop(v, coinprice, false) {
						continue
					}
					profit = ((v["num"].ToFloat()*1000)/coinprice)*v["openprice"].ToFloat() - v["num"].ToFloat()*1000 //计算出利益卖空的利益
				}
				profit = profit * v["ganggan"].ToFloat()

				if profit < 0 && (math.Abs(profit)/(v["num"].ToFloat()*1000)) >= 1 {
					//当亏损大于临界值 开始强制平仓 并且计算穿仓金额 从余额中扣除
					all_real_credit := profit + v["num"].ToFloat()*1000
					insertData := db.DB_PARAMS{}
					insertData["uid"] = v["uid"].ToInt()
					insertData["sn"] = v["sn"].ToString()
					insertData["coin_symbol"] = v["coin_symbol"].ToString()
					insertData["trade_type"] = v["trade_type"].ToInt()
					insertData["flag"] = v["flag"].ToInt()
					insertData["amount"] = v["num"].ToFloat()
					insertData["close_price"] = coinprice
					insertData["createtime"] = ntime
					insertData["num"] = v["num"].ToFloat()
					insertData["mode"] = models.USER_MODE_REAL
					insertData["allprice"] = all_real_credit
					insertData["profit"] = profit
					insertData["o_price"] = v["openprice"].ToFloat()
					_, err := config.GlobalDB.InsertData(models.DB_TABLE_CLOSE_TRADE, insertData) //插入平仓表
					if err == nil {
						//开始更改持仓表和用户金额
						config.GlobalDB.AddValue(models.DB_TABLE_OPENED_TRADE, map[string]float64{"num": -1 * v["num"].ToFloat(), "credit": -1 * v["num"].ToFloat() * 1000}, db.DB_PARAMS{"id": v["id"].Value}) //清空持仓
						models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
							Credit:          all_real_credit,
							LockCredit:      0,
							VCrdit:          0,
							LockVCredit:     0,
							UserCoinLogType: models.COIN_LOG_KEEP_BREAK, //穿仓
							UserCoinLogInfo: models.QueueCreditLog{
								Credit:     all_real_credit,
								Sn:         v["sn"].ToString(),
								CreateTime: ntime,
							},
							TeamCoinLogType: models.TEAM_LOG_TRADE_PROFIT,
							TeamCoinLogInfo: models.QueueTeamLog{
								TradeKeep_Profit: profit,
								CreateTime:       ntime,
							},
						})
					}
				}
			}
		}
		time.Sleep(3 * time.Second) //三秒一次穿仓扫描
	}
}
func ClearExplodeTrade() {
	//处理交割合约的持仓
	for {
		ntime := utils.GetNow()
		list, err := config.GlobalDB.FetchAll(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"trade_type": models.OPEN_TYPE_EXPLODE, "clear_time": 0, "_": fmt.Sprintf("close_real_time<=%d", ntime)}, db.DB_FIELDS{}, "limit 0,500") //每次处理500单
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for _, v := range list {
			//coinPriceInfo := config.GlobalMongo.GetOne("lastdata", bson.M{"pair": v["coinpair"].ToString()}, nil)
			coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(v["coinpair"].ToString())
			coinprice := coinPriceInfo["close"].(float64)
			coinInfo := models.MODEL_SYSTEM.GetCoinInfo(v["coin_symbol"].ToString(), v["coinpair"].ToString())

			iswin := false
			backCredit := 0.0
			profit := 0.0
			if (v["flag"].ToInt() == models.DIRECT_TYPE_BIG && coinprice > v["openprice"].ToFloat()) || (v["flag"].ToInt() == models.DIRECT_TYPE_SMALL && coinprice < v["openprice"].ToFloat()) {
				iswin = true
			} else {
				iswin = false
			}
			control := int32(0)
			//这里需要控制逻辑 后台操控输赢的逻辑
			user_explode_state_info, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER, db.DB_PARAMS{"id": v["uid"].Value}, db.DB_FIELDS{"explode_state"})
			if user_explode_state_info != nil && user_explode_state_info["explode_state"].ToInt() > 0 {
				switch user_explode_state_info["explode_state"].ToInt() {
				case 1:
					control = 1
				case 2:
					control = 2
				case 3:
					if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
						control = 1
					}
				case 4:
					if v["flag"].ToInt() == models.DIRECT_TYPE_SMALL {
						control = 1
					}
				case 5:
					if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
						control = 1
					} else {
						control = 2
					}
				case 6:
					if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
						control = 2
					} else {
						control = 1
					}
				}
			}
			sn_control := models.MODEL_SYSTEM.GetControlExplode(v["sn"].ToString())
			if sn_control > 0 {
				control = sn_control
			}
			diffPrice := float64((1 + rand.Intn(100))) / float64(math.Pow10(utils.GetInt(coinInfo["dnum"])))
			if control_price_min, bb := coinInfo["contorl_price_min"]; bb {
				if control_price_max, bb_max := coinInfo["contorl_price_max"]; bb_max {

					d_min := utils.GetFloat(control_price_min)
					d_max := utils.GetFloat(control_price_max)
					if d_min > 0 && d_max > 0 {
						d_num_pow := math.Pow10(utils.GetInt(coinInfo["dnum"]))
						statr_min := d_min * d_num_pow
						end_max := d_max * d_num_pow
						diffPrice = (statr_min + float64(rand.Intn(int(end_max-statr_min)))) / d_num_pow
					}

				}
			}
			if control == 1 {
				iswin = true
				if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
					coinprice = v["openprice"].ToFloat() + diffPrice
				} else {
					coinprice = v["openprice"].ToFloat() - diffPrice
				}
			} else if control == 2 {
				iswin = false
				if v["flag"].ToInt() == models.DIRECT_TYPE_BIG {
					coinprice = v["openprice"].ToFloat() - diffPrice
				} else {
					coinprice = v["openprice"].ToFloat() + diffPrice
				}
			}
			//config.GlobalMongo.Update("explode_control", map[string]interface{}{"result_time": utils.GetNow()}, bson.M{"sn": v["sn"].ToString()}) //更新结束控制
			//设置为任务完成状态
			config.GlobalMongo.DBHandle.Collection(models.COIN_CONTROLLER).DeleteOne(context.TODO(), bson.M{"sn": v["sn"].ToString()})
			if iswin {
				profit = v["credit"].ToFloat() * (v["win_rate"].ToFloat() / float64(100))
			} else {
				profit = -1 * v["credit"].ToFloat() * (v["lose_rate"].ToFloat() / float64(100))
			}
			backCredit = v["credit"].ToFloat() + profit
			if v["mode"].ToInt() == models.USER_MODE_REAL {
				fmt.Println(" 交割真实订单.....")
				models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
					Credit:          backCredit,
					LockCredit:      0,
					VCrdit:          0,
					LockVCredit:     0,
					UserCoinLogType: models.COIN_LOG_USER_CLOSE, //平仓账变
					UserCoinLogInfo: models.QueueCreditLog{
						Credit:     backCredit,
						LockCredit: 0,
						Sn:         v["sn"].ToString(),
						CreateTime: ntime,
					},
					TeamCoinLogType: models.TEAM_LOG_TRADE_PROFIT,
					TeamCoinLogInfo: models.QueueTeamLog{
						TradeExplode_Profit: profit,
						CreateTime:          ntime,
					},
				})
			} else {
				models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
					Credit:          0,
					LockCredit:      0,
					VCrdit:          backCredit,
					LockVCredit:     0,
					UserCoinLogType: 0,
					UserCoinLogInfo: nil,
					TeamCoinLogType: 0,
					TeamCoinLogInfo: nil,
				})
			}
			config.GlobalDB.UpdateData(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"clear_time": v["close_real_time"].ToInt(), "profit": profit, "closeprice": coinprice}, db.DB_PARAMS{"id": v["id"].ToInt()})
			//开始进入平仓表
			insertData := db.DB_PARAMS{}
			insertData["uid"] = v["uid"].ToInt()
			insertData["sn"] = v["sn"].ToString()
			insertData["coin_symbol"] = v["coin_symbol"].ToString()
			insertData["trade_type"] = v["trade_type"].ToInt()
			insertData["flag"] = v["flag"].ToInt()
			insertData["amount"] = v["num"].ToFloat()
			insertData["close_price"] = coinprice
			insertData["createtime"] = v["close_real_time"].ToInt()
			insertData["num"] = v["num"].ToFloat()
			insertData["mode"] = v["mode"].ToInt()
			insertData["allprice"] = backCredit
			insertData["o_price"] = v["openprice"].ToFloat()
			insertData["profit"] = profit
			config.GlobalDB.InsertData(models.DB_TABLE_CLOSE_TRADE, insertData)
			//models.MODEL_MESSAGE.PushMessage(v["uid"].ToInt(), models.MessageText{Content: v["sn"].ToString(), Title: fmt.Sprintf("explode:%s|%f", v["sn"].ToString(), profit)}, models.MESSAGE_TYPE_TEXT)
			//交割合约平仓的消息 推送给客户端

		}
		time.Sleep(500 * time.Millisecond) //100毫秒处理一次
	}
}
