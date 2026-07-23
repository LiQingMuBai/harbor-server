package models

import (
	"cointrade/config"
	shareddomain "cointrade/internal/domain/shared"
	"cointrade/lib/db"
	"cointrade/utils"
)

func (m *SystemModel) BuyCoin(uid int, coinID int, amount float64) *BaseResponse {
	ntime := utils.GetNow()
	coininfo, _ := config.GlobalDB.FetchOne(DB_TABLE_COINS, db.DB_PARAMS{"id": coinID}, db.DB_FIELDS{})
	if coininfo == nil {
		return &BaseResponse{State: STATE_FAILD, Msg: shareddomain.MsgCoinNotFound}
	}

	leaveAmount := coininfo["all_amount"].ToInt() - coininfo["selled_amount"].ToInt()
	if amount > float64(leaveAmount) {
		return &BaseResponse{State: COIN_BUY_STATE_NOTENGOUGH, Msg: shareddomain.MsgInsufficient}
	}

	allprice := amount * coininfo["f_price"].ToFloat()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if allprice <= 0 {
		return &BaseResponse{State: STATE_FAILD, Msg: shareddomain.MsgCoinNotFound}
	}
	if uinfo.Credit < allprice {
		return &BaseResponse{State: COIN_BUY_STATE_NOMONEY, Msg: shareddomain.MsgInsufficient}
	}

	if MODEL_USER.AddCredit(uid, &CreditValue{
		Credit:          -1 * allprice,
		LockCredit:      allprice,
		UserCoinLogType: COIN_LOG_BUY_COIN,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     -1 * allprice,
			LockCredit: allprice,
			CreateTime: ntime,
			CoinType:   "usdt",
		},
	}) {
		insertData := db.DB_PARAMS{
			"uid":         uid,
			"coin_id":     coininfo["id"].ToInt(),
			"coin_symbol": coininfo["symbol"].ToString(),
			"coin_pair":   coininfo["pair"].ToString(),
			"amount":      amount,
			"price":       coininfo["f_price"].ToFloat(),
			"all_price":   allprice,
			"createtime":  ntime,
		}
		config.GlobalDB.InsertData(DB_TABLE_BUY_COIN_ORDER, insertData)
		return &BaseResponse{State: STATE_SUCCESS, Msg: shareddomain.MsgOK}
	}
	return nil
}

func (m *SystemModel) GetBuyCoinOrders(uid int) db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return list
}
