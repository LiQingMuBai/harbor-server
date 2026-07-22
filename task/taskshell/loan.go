package taskshell

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"time"
)

func LoanCount() {
	//贷款计算
	for {
		ntime := utils.GetNow()
		list, err := config.GlobalDB.FetchAll(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"state": 1, "_": fmt.Sprintf("interest_time<%d", ntime)}, db.DB_FIELDS{})
		if err == nil {
			for _, v := range list {
				if v["circle"].ToInt() == v["process"].ToInt() {
					//借款到期了
					credit := v["amount"].ToFloat() + v["all_interest"].ToFloat()
					models.MODEL_USER.AddCredit(
						v["uid"].ToInt(),
						&models.CreditValue{
							Credit:          -1 * credit,
							LockCredit:      0,
							UserCoinLogType: models.COIN_LOG_LOAN_BACK,
							UserCoinLogInfo: models.QueueCreditLog{
								Credit:     -1 * credit,
								CoinType:   "usdt",
								CreateTime: ntime,
								Sn:         v["sn"].ToString(),
							},
						},
					)
					config.GlobalDB.UpdateData(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"state": 3}, db.DB_PARAMS{"id": v["id"].Value})
					continue
				}
				f := v["amount"].ToFloat() * (v["rate"].ToFloat() / float64(100))

				config.GlobalDB.AddValue(models.DB_TABLE_LOAN_ORDER, map[string]float64{"all_interest": f, "interest_time": v["interest_time"].ToFloat() + 24*60*60, "process": 1}, db.DB_PARAMS{"id": v["id"].Value})
			}
		}
		time.Sleep(10 * time.Second)
	}
}
