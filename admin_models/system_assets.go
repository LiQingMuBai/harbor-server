package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

func (s *SystemModel) WalletAddressList(rq P) *AdminResponse {
	where := ""
	if v := rq.Ts().Get("search").ToString(); v != "" {
		where = fmt.Sprintf(" cointype like '%%%s%%' OR contract like '%%%s%%'", v, v)
	}
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"_": where}, db.DB_FIELDS{}, utils.Limit(rq.Ts().Get("paghe").ToInt(), rq.Ts().Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{})
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"list":      list,
			"count":     count,
			"chan_type": s.ContractFlag(),
		},
	}
}

func (s *SystemModel) OpenAddr(rq P) *AdminResponse {
	t := rq.Ts()
	if v := t.Get("id").ToInt(); v == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要开启的通道"}
	}
	if c := config.GlobalDB.GetCount(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"id": t.Get("id").ToInt()}); c == 9 {
		return &AdminResponse{State: ERROR, Data: "系统无法找到该收款信息"}
	}
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": t.Get("state").ToInt()}, db.DB_PARAMS{"id": t.Get("id").ToInt()}); err == nil {
		return &AdminResponse{State: SUCCESS, Data: "修改收款状态成功!"}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	return &AdminResponse{State: ERROR, Data: "修改失败了！"}
}

func (s *SystemModel) DeleteWalletAddress(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的收款信息"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "删除收款信息失败!"}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	return &AdminResponse{State: SUCCESS, Data: "删除收款信息成功!"}
}

func (s *SystemModel) SaveWalletAddress(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("cointype").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种类型不能为空!"
		return rs
	}
	if v := t.Get("contract").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "合约标识不能为空!"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "logo不能为空!"
	}
	if v := t.Get("address").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "收款通道地址不能为空!"
		return rs
	}
	insert := P{
		"cointype":     t.Get("cointype").ToString(),
		"logo":         t.Get("logo").ToString(),
		"contract":     t.Get("contract").ToString(),
		"address":      t.Get("address").ToString(),
		"state":        t.Get("state").ToInt(),
		"min":          t.Get("min").ToFloat(),
		"withdraw_min": t.Get("withdraw_min").ToFloat(),
	}
	var err error
	if t.Get("id").ToInt() > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_ADDRESS, insert, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_RECHARGE_ADDRESS, insert)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "插入数据库失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	rs.State = SUCCESS
	rs.Data = "操作收款钱包信息成功!"
	return rs
}

func (s *SystemModel) CoinKeyValPair() map[int]string {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make(map[int]string, 0)
	for _, v := range list {
		rs[v.Get("id").ToInt()] = v.Get("symbol").ToString()
	}
	return rs
}

func (s *SystemModel) CoinTypePair() map[string]*RechargeAddress {
	r := make(map[string]*RechargeAddress, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, item := range list {
		rg := new(RechargeAddress)
		item.SetObj(rg)
		r[item.Get("cointype").ToString()] = rg
	}
	return r
}

func (s *SystemModel) TranferCoin(direct int) map[string]map[string]interface{} {
	p := db.DB_PARAMS{}
	if direct == 1 {
		p["is_in"] = 1
	} else {
		p["is_out"] = 1
	}
	r := make(map[string]map[string]interface{}, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, p, db.DB_FIELDS{})
	for _, item := range list {
		rg := map[string]interface{}{}
		item.SetObj(rg)
		r[item.Get("symbol").ToString()] = rg
	}
	return r
}

func (s *SystemModel) ContractFlag() []string {
	return []string{"ETH", "TRON", "BITCOIN", "SOLANA"}
}

func (s *SystemModel) ContractList() []string {
	l, _ := config.GlobalDB.FetchAll(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]string, 0)
	for _, i := range l {
		rs = append(rs, i.Get("contract").ToString())
	}
	return rs
}

func (s *SystemModel) CoinList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("open_coin2coin").ToString(); v != "" {
		where = append(where, fmt.Sprintf("open_coin2coin = %s", v))
	}
	if v := t.Get("open_trade").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" open_trade = %s ", v))
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{}, utils.Order(t.Get("sort").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_COINS, db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	l := make([]*Coin, 0)
	for _, item := range list {
		coin := new(Coin)
		item.SetObj(coin)
		l = append(l, coin)
	}
	return &AdminResponse{State: SUCCESS, Data: P{"list": l, "count": count}}
}

func (s *SystemModel) CoinDescList(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("lang").ToString(); v != "" {
		where = append(where, fmt.Sprintf("lang = '%s'", v))
	}
	if v := pdata.Get("symbol").ToString(); v != "" {
		where = append(where, fmt.Sprintf("symbol  = '%s'", v))
	}
	count := config.GlobalDB.GetCount(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Limit(pdata.Get("page").ToInt()))
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"count":    count,
			"list":     list,
			"coinList": s.CoinKeyValPair(),
			"langList": s.LangeList(),
		},
	}
}

func (s *SystemModel) SaveCoinDesc(rq P) *AdminResponse {
	info := rq.Ts()
	coininfo := P{
		"symbol":     info.Get("symbol").ToString(),
		"lang":       info.Get("lang").ToString(),
		"desc":       info.Get("desc").ToString(),
		"pubtime":    utils.TimeToint64(info.Get("pubtime").ToString()),
		"totalnum":   info.Get("totalnum").ToInt(),
		"whitepaper": info.Get("whitepaper").ToString(),
		"website":    info.Get("website").ToString(),
	}
	var descErr error
	if v := info.Get("id").ToInt(); v > 0 {
		_, descErr = config.GlobalDB.UpdateData(models.DB_TABLE_COIN_DESC, coininfo, db.DB_PARAMS{"id": v})
	} else {
		exists := config.GlobalDB.GetCount(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"lang": coininfo["lang"], "symbol": coininfo["symbol"]})
		if exists > 0 {
			return &AdminResponse{State: ERROR, Data: "该币种介绍已经存在 ，请勿重复添加"}
		}
		_, descErr = config.GlobalDB.InsertData(models.DB_TABLE_COIN_DESC, coininfo)
	}
	if descErr != nil {
		return &AdminResponse{State: ERROR, Data: descErr.Error()}
	}
	cacheID := models.MODEL_SYSTEM.MakeCacheId(coininfo["symbol"], coininfo["lang"])
	config.GlobalRedis.Del(models.HASH_COIN_DESC, cacheID)
	return &AdminResponse{State: SUCCESS, Data: "操作币种信息成功!"}
}

func (s *SystemModel) DeleteCoinDesc(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的币种信息"}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "该介绍信息不存在!"}
	}
	config.GlobalDB.Delete(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"id": id})
	cacheID := models.MODEL_SYSTEM.MakeCacheId(one.Get("symbol").ToString(), one.Get("lang").ToString())
	config.GlobalRedis.Del(models.HASH_COIN_DESC, cacheID)
	return &AdminResponse{State: SUCCESS, Data: "删除币种信息成!"}
}

func (s *SystemModel) DeleteCoin(rq P) *AdminResponse {
	id := rq.Ts().Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的币种信息"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_COINS, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "删除信息失败!"}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_COIN_LIST})
	return &AdminResponse{State: SUCCESS, Data: "删除币种信息成功!"}
}

func (s *SystemModel) SaveCoin(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("name").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种名称请填写"
		return rs
	}
	if v := t.Get("symbol").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种简称请填写"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种LOGO请上传"
	}
	if v := t.Get("address").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种对应的钱包地址不能为空"
	}
	if isnative := t.Get("isnative").ToInt(); isnative == 0 {
		if v := t.Get("vpair").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发币币种必须填写相对币种行情"
			return rs
		}
	}

	insert := P{
		"name":              t.Get("name").ToString(),
		"symbol":            t.Get("symbol").ToString(),
		"pair":              fmt.Sprintf("%susdt", t.Get("symbol").ToString()),
		"logo":              t.Get("logo").ToString(),
		"desc":              t.Get("desc").ToString(),
		"open_coin2coin":    t.Get("open_coin2coin").ToInt(),
		"open_trade":        t.Get("open_trade").ToInt(),
		"isnative":          t.Get("isnative").ToInt(),
		"dnum":              t.Get("dnum", "2").ToInt(),
		"sort":              t.Get("sort").ToInt(),
		"cnum":              t.Get("cnum", 6).ToInt(),
		"baseprice":         t.Get("baseprice").ToFloat(),
		"min_price_float":   t.Get("min_price_float").ToFloat(),
		"max_price_float":   t.Get("max_price_float").ToFloat(),
		"max_float":         t.Get("max_float").ToFloat(),
		"vpair":             t.Get("vpair").ToString(),
		"address":           t.Get("address").ToString(),
		"is_market":         t.Get("is_market").ToInt(),
		"is_f":              t.Get("is_f").ToInt(),
		"f_price":           t.Get("f_price").ToFloat(),
		"is_in":             t.Get("is_in").ToInt(),
		"is_out":            t.Get("is_out").ToInt(),
		"is_new":            t.Get("is_new").ToInt(),
		"pubtime":           utils.TimeToint64(t.Get("pubtime").ToString()),
		"all_amount":        t.Get("all_amount").ToFloat(),
		"contorl_price_min": t.Get("contorl_price_min").ToFloat(),
		"contorl_price_max": t.Get("contorl_price_max").ToFloat(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_COINS, insert, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_COINS, insert)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作币种信息失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_COIN_LIST})
	rs.State = SUCCESS
	rs.Data = "操作成功!"
	return rs
}

func (s *SystemModel) CurrencyList() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_CURRENCY, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]*models.Currency, 0)
	for _, v := range list {
		n := new(models.Currency)
		v.SetObj(n)
		rs = append(rs, n)
	}
	return &AdminResponse{State: SUCCESS, Data: rs}
}

func (s *SystemModel) DeleteCurrency(id string) *AdminResponse {
	if id == "" {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的货币信息"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_CURRENCY, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "删除货币信息失败!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "删除货币信息成功!"}
}

func (s *SystemModel) SaveCurrency(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("symbol").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "货币标识不能为空!"
		return rs
	}
	if v := t.Get("rate").ToFloat(); v == 0 {
		rs.State = ERROR
		rs.Data = "换算汇率不能为空！"
		return rs
	}
	if v := t.Get("country").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "货币所属国家不能为空!"
		return rs
	}
	in := P{
		"symbol":  t.Get("symbol").ToString(),
		"rate":    t.Get("rate").ToFloat(),
		"country": t.Get("country").ToString(),
		"memo":    t.Get("memo").ToString(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_CURRENCY, in, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_CURRENCY, in)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作货币信息失败!"
		return rs
	}
	rs.State = SUCCESS
	rs.Data = "操作货币信息成功!"
	return rs
}

func (s *SystemModel) CurrentCyPair() map[int]string {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make(map[int]string, 0)
	for _, v := range list {
		rs[v.Get("id").ToInt()] = v.Get("symbol").ToString()
	}
	return rs
}

func (s *SystemModel) CurrentcyList(rq P) *AdminResponse {
	t := rq.Ts()
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_CURRENCY, db.DB_PARAMS{}, db.DB_FIELDS{}, utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_CURRENCY, db.DB_PARAMS{})
	return &AdminResponse{State: SUCCESS, Data: P{"list": list, "count": count}}
}

func (s *SystemModel) LangeList() map[string]string {
	l := make(map[string]string)
	l["zh"] = "中文"
	l["fr"] = "法语"
	l["zh-tw"] = "台湾"
	l["es"] = "西班牙语"
	l["en"] = "英文"
	l["th"] = "泰文"
	l["ja"] = "日文"
	l["ko"] = "韩文"
	l["ar"] = "阿拉伯语"
	l["vi"] = "越南语"
	return l
}
