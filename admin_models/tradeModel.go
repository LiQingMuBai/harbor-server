package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

type TradeModel struct{}

/**
 *	开仓信息
 */
func (m *TradeModel) TradeList(rq P, isAdmin bool, agentCode string) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if !isAdmin {
		//where = append(where, fmt.Sprintf("parent_uid  = %s", agentCode))
		//fmt.Printf("不是管理员，归属团队%s\n", agentCode)
		where = append(where, fmt.Sprintf("u.team like '%%%s%%'", agentCode))

	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("o.createtime between %d AND %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("trade_type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.trade_type = %d", v))
	}
	if v := t.Get("coinid").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.coinid = %d", v))
	}
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (o.sn = '%s' OR u.username like '%s' u.id like '%s')", v, v, v))
	}
	if v := t.Get("flag").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.flag = %d", v))
	}
	if v := t.Get("mode").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.mode > %d", v))
	}
	if v := t.Get("state").ToInt(); v > 0 {
		if v == 1 { //当前订单状态
			where = append(where, " o.clear_time = 0")
		} else { //平仓
			where = append(where, " o.clear_time >0 ")
		}
	}
	rs := make([]*models.OpenedInfo, 0)
	count := config.GlobalDB.JoinCount(models.DB_TABLE_OPENED_TRADE+" as o ", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})

	if count == 0 {
		return &AdminResponse{
			State: SUCCESS,
			Data: P{
				"list":       rs,
				"count":      count,
				"flag_list":  SYSTEM_MODEL.DirectPair(),
				"trade_type": SYSTEM_MODEL.TradePair(),
				"mode_list":  SYSTEM_MODEL.ModePair(),
				"win_state":  map[int]string{1: "盈利", 2: "亏损"},
				"coin_list":  SYSTEM_MODEL.CoinKeyValPair(),
			},
		}
	}
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_OPENED_TRADE+" as o ", models.DB_TABLE_USER+" as   u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"o.*", "u.username", "u.user_type", "u.memo"}, utils.Order(t.Get("sort", "o.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))

	for _, r := range list {
		trade := new(models.OpenedInfo)
		r.SetObj(trade)
		trade.TradeType = r.Get("trade_type").ToInt()
		trade.CoinSymbol = r.Get("coin_symbol").ToString()
		trade.CoinPair = r.Get("coinpair").ToString()
		rs = append(rs, trade)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":       rs,
			"count":      count,
			"flag_list":  SYSTEM_MODEL.DirectPair(),
			"trade_type": SYSTEM_MODEL.TradePair(),
			"mode_list":  SYSTEM_MODEL.ModePair(),
			"win_state":  map[int]string{1: "盈利", 2: "亏损"},
			"coin_list":  SYSTEM_MODEL.CoinKeyValPair(),
		},
	}
}

/**
 *	历史平仓
 */
func (m *TradeModel) HistoryCloseRradeList(rq P, isAdmin bool, agentCode string) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)

	if !isAdmin {
		//where = append(where, fmt.Sprintf("parent_uid  = %s", agentCode))
		//fmt.Printf("不是管理员，归属团队%s\n", agentCode)
		where = append(where, fmt.Sprintf("u.team like '%%%s%%'", agentCode))

	}
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("  (o.sn like '%%%s%%' OR  u.username like '%%%s%%' OR  u.id='%s')", v, v, v))
	}
	if v := t.Get("win_state").ToInt(); v > 0 {
		if v == 1 {
			where = append(where, "o.profit > 0")
		} else {
			where = append(where, "o.profit < 0")
		}
	}
	if v := t.Get("flag").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" o.flag = %d", v))
	}
	if v := t.Get("type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" o.trade_type = %d", v))
	}
	if v := t.Get("coin_pair").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.coin_symbol = '%s'", v))
	}
	if v := t.Get("mode").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.mode = %d", v))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("o.createtime between  %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	rs := make([]*models.CloseTrade, 0)
	count := config.GlobalDB.JoinCount(models.DB_TABLE_CLOSE_TRADE+" as o", models.DB_TABLE_USER+" as u", " o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})

	if count == 0 {
		return &AdminResponse{
			State: SUCCESS,
			Data: P{
				"list":       rs,
				"count":      count,
				"flag_list":  SYSTEM_MODEL.DirectPair(),
				"trade_type": SYSTEM_MODEL.TradePair(),
				"mode_list":  SYSTEM_MODEL.ModePair(),
				"win_state":  map[int]string{1: "盈利", 2: "亏损"},
				"coin_list":  SYSTEM_MODEL.CoinKeyValPair(),
			},
		}
	}
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_CLOSE_TRADE+" as o", models.DB_TABLE_USER+" AS   u", " o.uid =u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"o.*", "u.username", "u.user_type", "u.memo"},
		utils.Order(t.Get("sort", "o.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))

	for _, r := range list {
		trade := new(models.CloseTrade)
		r.SetObj(trade)
		rs = append(rs, trade)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":       rs,
			"count":      count,
			"flag_list":  SYSTEM_MODEL.DirectPair(),
			"trade_type": SYSTEM_MODEL.TradePair(),
			"mode_list":  SYSTEM_MODEL.ModePair(),
			"win_state":  map[int]string{1: "盈利", 2: "亏损"},
			"coin_list":  SYSTEM_MODEL.CoinKeyValPair(),
		},
	}
}

func (m *TradeModel) DelegateHistoryDel(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确定要删除的委托信息",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "委托信息不存在!",
		}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"id": id}); err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除委托信息成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "委托信息删除失败!",
	}
}

func (m *TradeModel) OpSpot(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前订单信息不存在!",
		}
	}
	deleteinfo, _ := config.GlobalDB.FetchOne(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "is_f": 1}, db.DB_FIELDS{})

	update := map[string]interface{}{}
	if pdata.Get("state").ToInt() == 1 {
		if pdata.Get("days").ToString() == "" {
			return &AdminResponse{
				State: ERROR,
				Data:  "解锁时间不能为空",
			}
		}
		if deleteinfo == nil {
			return &AdminResponse{
				State: ERROR,
				Data:  "当前订单不存在或不为自发币订单",
			}
		}
		if pdata.Get("num").ToInt() > deleteinfo.Get("num").ToInt() {
			return &AdminResponse{
				State: ERROR,
				Data:  "审核数量不能超过用户申请数量",
			}
		}
		price := pdata.Get("price").ToFloat() * pdata.Get("num").ToFloat()
		if price > deleteinfo.Get("credit").ToFloat() {
			return &AdminResponse{
				State: ERROR,
				Data:  "通过总额不能超过用户原始冻结的金额",
			}
		}
		plus := deleteinfo.Get("credit").ToFloat() - price
		if plus >= 0 {
			models.MODEL_USER.AddCredit(deleteinfo.Get("uid").ToInt(), &models.CreditValue{
				Credit:          plus,
				LockCredit:      -1 * deleteinfo.Get("credit").ToFloat(),
				UserCoinLogType: models.COIN_LOG_SPOT_BACK,
			})
		}
		//加入资产
		b := models.MODEL_ASSETS.AddAssets(deleteinfo.Get("uid").ToInt(), &models.Assets{
			Coin:          deleteinfo.Get("coin_symbol").ToString(),
			Pair:          deleteinfo.Get("coinpair").ToString(), //交易对
			Num:           pdata.Get("num").ToFloat(),            //数量
			LockNum:       0,
			Price:         pdata.Get("price").ToFloat(), //开仓价格
			Mode:          1,                            //模式
			IsTrans:       0,                            //是否可以交易划转
			OpenTransTime: int(utils.TimeToint64(pdata.Get("days").ToString())),
		})
		fmt.Println("加入资产", b, int(utils.TimeToint64(pdata.Get("days").ToString())))
		update["state"] = 1
		//update["price"] = pdata.Get("price").ToFloat()
		//update["num"] = pdata.Get("num").ToFloat()
		update["changetime"] = utils.GetNow()
	} else {
		models.MODEL_USER.AddCredit(deleteinfo.Get("uid").ToInt(), &models.CreditValue{
			Credit:          deleteinfo.Get("credit").ToFloat(),
			LockCredit:      -1 * deleteinfo.Get("credit").ToFloat(),
			UserCoinLogType: models.COIN_LOG_SPOT_BACK,
		})
		update["state"] = 2
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, update, db.DB_PARAMS{"id": deleteinfo.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  err.Error(),
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "审核成功",
	}
}

/**
 *	历史委托表
 */
func (m *TradeModel) HistoryDelegateList(rq P, isAdmin bool, agentCode string) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)

	if !isAdmin {
		//where = append(where, fmt.Sprintf("parent_uid  = %s", agentCode))
		//fmt.Printf("不是管理员，归属团队%s\n", agentCode)
		where = append(where, fmt.Sprintf("u.team like '%%%s%%'", agentCode))

	}

	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(o.sn like '%%%s%%' OR u.username like '%%%s%%')", v, v))
	}
	if v := t.Get("delegate_type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.delegate_type = %d", v))
	}
	if v := t.Get("type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.trade_type = %d", v))
	}
	if v := t.Get("flag").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.flag = %d", v))
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("o.state = %d", t.Get("state").ToInt()))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf(" o.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("coinid").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.coinid= %d", v))
	}
	if v := t.Get("mode").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.mode = %d", v))
	}
	if v := t.Get("is_f").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.is_f = %d", v))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_DELEGATE_TRADE+" as o", models.DB_TABLE_USER+" as  u ", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"o.*", "u.username", "u.user_type", "u.memo"}, utils.Order(t.Get("sort", "o.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	rs := make([]*models.DelegateInfo, 0)
	for _, v := range list {
		delegate := new(models.DelegateInfo)
		v.SetObj(delegate)
		rs = append(rs, delegate)
	}
	count := config.GlobalDB.JoinCount(models.DB_TABLE_DELEGATE_TRADE+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":              rs,
			"count":             count,
			"coin_list":         SYSTEM_MODEL.CoinKeyValPair(),
			"delegatetype_list": SYSTEM_MODEL.DelegateType(),
			"flag_list":         SYSTEM_MODEL.DirectPair(),
			"mode_list":         SYSTEM_MODEL.ModePair(),
			"trade_type":        SYSTEM_MODEL.TradePair(),
		},
	}
}

func (m *TradeModel) ManualOperationTrade(uid int, sn string) *AdminResponse {
	if sn == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确定要手工成交的委托信息",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "手工成交委托信息不存在!",
		}
	}
	update := map[string]interface{}{}
	update["state"] = 1
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_DELEGATE_TRADE, update, db.DB_PARAMS{"sn": sn}); err == nil {
		insertData := db.DB_PARAMS{}
		insertData["uid"] = one.Get("uid").ToInt()
		insertData["sn"] = one.Get("sn").ToString()
		insertData["trade_type"] = one.Get("trade_type").ToInt()
		insertData["flag"] = one.Get("flag").ToInt()
		insertData["openprice"] = one.Get("price").ToFloat()
		insertData["closeprice"] = 0
		insertData["coinid"] = one.Get("coinid").ToInt()
		insertData["coinpair"] = one.Get("coinpair").ToString()
		insertData["coin_symbol"] = one.Get("coin_symbol").ToString()
		insertData["close_time"] = one.Get("close_time").ToInt()
		insertData["close_real_time"] = utils.GetNow() + one.Get("close_time").ToInt()
		insertData["clear_time"] = 0
		insertData["createtime"] = utils.GetNow()
		insertData["ganggan"] = 1
		insertData["credit"] = one.Get("credit").ToFloat()
		insertData["profit"] = 0
		if explodeConfig, ok := models.EXPLODE_CONFIG[one.Get("close_time").ToInt()]; ok {
			insertData["win_rate"] = explodeConfig.Winrate
			insertData["lose_rate"] = explodeConfig.Loserate
		} else {
			insertData["win_rate"] = 100
			insertData["lose_rate"] = 100
		}

		insertData["num"] = one.Get("num").ToFloat()
		insertData["mode"] = one.Get("mode").ToInt()

		if _, err := config.GlobalDB.InsertData(models.DB_TABLE_OPENED_TRADE, insertData); err == nil {
			return &AdminResponse{
				State: SUCCESS,
				Data:  "手工成交委托信息成功!",
			}
		}
	}

	return &AdminResponse{
		State: ERROR,
		Data:  "委托信息手工成交失败!",
	}
}
