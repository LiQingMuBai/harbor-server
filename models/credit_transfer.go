package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func (m *CreditModel) AddWallet(uid int, rq *WalletAddressRequest) *BaseResponse {
	rs := new(BaseResponse)
	if rq.Address == "" || rq.CoinType == "" || rq.Contract == "" {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	count := config.GlobalDB.GetCount(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"cointype": rq.CoinType, "contract": rq.Contract})
	if count >= 10 {
		rs.State = STATE_FAILD
		rs.Msg = "too more"
		return rs
	}
	insertData := db.DB_PARAMS{
		"uid":        uid,
		"createtime": utils.GetNow(),
		"cointype":   rq.CoinType,
		"contract":   rq.Contract,
		"address":    rq.Address,
		"title":      rq.Title,
	}
	config.GlobalDB.InsertData(DB_TABLE_USER_WITHDRAW_WALLET, insertData)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *CreditModel) DeleteWallet(uid, id int) *BaseResponse {
	config.GlobalDB.Delete(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"uid": uid, "id": id})
	return &BaseResponse{State: STATE_SUCCESS, Msg: "success"}
}

func (m *CreditModel) GetWalletList(uid int) db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return list
}

func (m *CreditModel) TransFer(uid int, trans *TransferRequest) *BaseResponse {
	ntime := utils.GetNow()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	trans.Coin = strings.ToLower(trans.Coin)
	toAddress := ""

	if uinfo == nil {
		return &BaseResponse{State: STATE_SYSTEM_ERROR, Msg: "SYSTEM ERROR"}
	}
	if trans.Amount <= 0 {
		return &BaseResponse{State: RECHARGE_STATE_MIN, Msg: "too min"}
	}

	if strings.Index(trans.Coin, "usdt") >= 0 {
		trans.Coin = "usdt"
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_TRANSFER)
	insertData := db.DB_PARAMS{"uid": uid, "sn": sn, "createtime": ntime}
	asset := MODEL_ASSETS.GetOneAsset(uid, trans.Coin)

	if trans.Coin != "usdt" && trans.Coin != "usdc" && asset == nil {
		return &BaseResponse{State: STATE_FAILD, Msg: "error assets"}
	}
	if trans.Coin != "usdt" && trans.Coin != "usdc" && asset != nil {
		if asset.IsTrans != 1 {
			return &BaseResponse{State: STATE_FAILD, Msg: "error assets  12"}
		}
		toAddress = asset.Address
	} else if trans.ToAddress != "" {
		toAddress = trans.ToAddress
	} else {
		toAddress = uinfo.WalletAddress
	}

	if trans.Direction == TRANSFER_DIRECTION_OUT {
		toAddress = trans.ToAddress
		if toAddress == "" {
			return &BaseResponse{State: RECHARGE_STATE_ERROR_ADDRESS, Msg: "error address"}
		}
		day := time.Now()
		today := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local).Unix()
		count := config.GlobalDB.GetCount(DB_TABLE_TRANSFER, db.DB_PARAMS{"direction": TRANSFER_DIRECTION_OUT, "_": fmt.Sprintf("createtime >= %d and state != 2 and uid = %d", today, uinfo.Id)})
		if count > config.GlobalConfig.GetValue("max_withdrawnum").ToInt() {
			return &BaseResponse{State: RECHARGE_STATE_ERROR_MAX_WITHDRAW, Msg: "max withdraw"}
		}

		if trans.Coin != "usdt" {
			if asset.Count < trans.Amount {
				return &BaseResponse{State: RECHARGE_STATE_ERROR_MONEY, Msg: "not enough assets"}
			}
			MODEL_ASSETS.AddAssets(uid, &Assets{
				Coin:    trans.Coin,
				Pair:    trans.Coin + "usdt",
				Num:     -1 * trans.Amount,
				LockNum: trans.Amount,
				Mode:    USER_MODE_REAL,
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: COIN_LOG_USER_WITHDRAW,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   trans.Coin,
				},
			})
		} else {
			if trans.Amount > uinfo.Credit {
				return &BaseResponse{State: RECHARGE_STATE_ERROR_MONEY, Msg: "not enough assets"}
			}
			if !MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: COIN_LOG_EXCHANGE_ACCOUNT_OUT,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   "usdt",
				},
			}) {
				return &BaseResponse{State: STATE_SYSTEM_ERROR, Msg: "ERROR addcredit"}
			}
		}
	}

	insertData["to_address"] = toAddress
	insertData["direction"] = trans.Direction
	insertData["amount"] = trans.Amount
	insertData["coin_symbol"] = trans.Coin
	_, err := config.GlobalDB.InsertData(DB_TABLE_TRANSFER, insertData)
	if err == nil {
		return &BaseResponse{State: STATE_SUCCESS, Msg: "ok"}
	}
	return &BaseResponse{State: STATE_SYSTEM_ERROR, Msg: "error"}
}

func (m *CreditModel) ExchangeAccount(uid int, rq *ExchangeAccountRequest) *BaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if rq.Drection == EXCHANGE_DIRECTION_CONTRACT {
		return m.RechargeByApprove(uid, rq.Amount)
	}
	if rq.Drection == EXCHANGE_DIRECTION_ACCOUNT {
		if uinfo.IsWithDraw == 0 {
			return &BaseResponse{State: WIDTHDRAW_STATE_ERROR_LOCKED, Msg: uinfo.WithDrawMsg}
		}
		return m.TransFer(uid, &TransferRequest{
			Coin:      "usdt",
			Amount:    rq.Amount,
			Direction: TRANSFER_DIRECTION_OUT,
			ToAddress: uinfo.WalletAddress,
		})
	}
	return nil
}

func (m *CreditModel) ExchangeAccount2(uid int, rq *ExchangeAccountRequest2) *BaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo.IsWithDraw == 0 {
		return &BaseResponse{State: WIDTHDRAW_STATE_ERROR_LOCKED, Msg: uinfo.WithDrawMsg}
	}

	amount, err := strconv.ParseFloat(rq.Amount, 64)
	if err != nil {
		return &BaseResponse{State: STATE_SYSTEM_ERROR, Msg: "SYSTEM ERROR"}
	}

	if rq.Symbol == "" || strings.ToLower(rq.Symbol) != "usdc" {
		rq.Symbol = "usdt"
	}
	return m.TransFer(uid, &TransferRequest{
		Coin:      strings.ToLower(rq.Symbol),
		Amount:    amount,
		Direction: TRANSFER_DIRECTION_OUT,
		ToAddress: rq.Address,
	})
}

func (m *CreditModel) TransFerLogs(uid int, rq *TransFerLogsRequest) *PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	if rq.Direction > -1 {
		condition["direction"] = rq.Direction
	}
	count := config.GlobalDB.GetCount(DB_TABLE_TRANSFER, condition)
	limit := 15
	pagesize := int(math.Ceil(float64(count) / float64(limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_TRANSFER, condition, db.DB_FIELDS{}, "order by id desc", fmt.Sprintf("limit %d,%d", (rq.Page-1)*limit, limit))
	rs := new(PageBaseResponse)
	rs.Limit = limit
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Total = count
	rs.PageTotal = pagesize
	rs.Page = rq.Page
	rs.List = list
	return rs
}

func (m *CreditModel) TransFerDetail(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_TRANSFER, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}
