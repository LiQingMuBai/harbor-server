package taskshell

import (
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"time"
)

func ClearDelegateTrade() { //开始清理委托单
	//go ClearExplodeTrade() //清理交割合约
	for {
		limit := 500 //一次处理500单
		list, _ := models.MODEL_TRADE.ListPendingDelegates(limit)
		if len(list) > 0 {
			for _, v := range list {
				OpDelegateTrade(v)
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func OpDelegateTrade(one db.DBValues) bool {
	return models.MODEL_TRADE.OperateDelegateTrade(one)
}
func EqualsPrice(myprice, coinprice float64, isbig bool) bool {
	return models.MODEL_TRADE.EqualsPrice(myprice, coinprice, isbig)
}
func CheckStop(v db.DBValues, coinprice float64, isbig bool) bool {
	return models.MODEL_TRADE.CheckStop(v, coinprice, isbig)
}
func ClearKeepCross() {
	//扫描永续杠杆穿仓的订单 然后强制平仓 亏损大于投资金额的80%强行平仓
	for {
		limit := 1000 //一次处理1000张订单
		_, pagesize := models.MODEL_TRADE.ListKeepCrossOpened(1, limit)
		for i := 1; i <= pagesize; i++ {
			list, _ := models.MODEL_TRADE.ListKeepCrossOpened(i, limit)
			for _, v := range list {
				ntime := utils.GetNow()
				coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(v["coinpair"].ToString())
				coinprice := coinPriceInfo["close"].(float64)
				_ = models.MODEL_TRADE.SettleKeepCross(v, coinprice, ntime)
			}
		}
		time.Sleep(3 * time.Second) //三秒一次穿仓扫描
	}
}
func ClearExplodeTrade() {
	//处理交割合约的持仓
	for {
		ntime := utils.GetNow()
		list, err := models.MODEL_TRADE.ListDueExplodeTrades(500, ntime)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for _, v := range list {
			_ = models.MODEL_TRADE.SettleExplodeTrade(v, ntime)
		}
		time.Sleep(500 * time.Millisecond) //100毫秒处理一次
	}
}
