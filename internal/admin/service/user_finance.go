package service

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
	"unsafe"
)

func (m *UserModel) ListCoinApplications(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("c.coin_symbol ='%s", v))
	}
	if v := pdata.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("c.state = %s", v))
	}
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("u.uid = '%s' OR u.username like '%%%s%%'", v, v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("c.createtime between %d and %d ", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	count := config.GlobalDB.JoinCount(models.DB_TABLE_BUY_COIN_ORDER+" as c", models.DB_TABLE_USER+" as u ", "c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_BUY_COIN_ORDER+" as c", models.DB_TABLE_USER+" as u ", "c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{"c.*", "u.username", "u.memo", "u.memo"}, utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))

	r := make([]map[string]interface{}, 0)
	cointype := make([]string, 0)
	coin, _ := config.GlobalDB.FetchAll(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{}, db.DB_FIELDS{"distinct(coin_symbol)"})
	for _, co := range coin {
		cointype = append(cointype, co.Get("coin_symbol").ToString())
	}
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		i["parent_name"] = uinfo.ParentName
		r = append(r, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"count":     count,
			"list":      r,
			"coin_type": cointype,
		},
	}
}

func (m *UserModel) DeleteCoinApplication(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的信息"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "删除信息失败!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "删除信息成功!"}
}

func (m *UserModel) ReviewCoinApplication(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请确定一个要操作的申请信息!"}
	}
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{State: ERROR, Data: "操作密码错误"}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "不存在该申请信息或者该信息已经处理完成"}
	}

	direct := pdata.Get("state").ToInt()
	if direct == 1 {
		models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{Credit: 0, LockCredit: one.Get("all_price").ToFloat() * -1, UserCoinLogType: models.COIN_LOG_BUY_COIN})
		models.MODEL_ASSETS.AddAssets(one.Get("uid").ToInt(), &models.Assets{
			Coin:    one.Get("coin_symbol").ToString(),
			Pair:    one.Get("coin_pair").ToString(),
			Num:     one.Get("amount").ToFloat(),
			LockNum: 0,
			Price:   one.Get("price").ToFloat(),
			Mode:    1,
			IsTrans: 0,
		})
		config.GlobalDB.AddValue(models.DB_TABLE_COINS, map[string]float64{"selled_amount": one.Get("amount").ToFloat()}, map[string]interface{}{"symbol": one.Get("coin_symbol").ToString()})
	} else {
		models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{Credit: one.Get("all_price").ToFloat(), LockCredit: one.Get("all_price").ToFloat() * -1, UserCoinLogType: models.COIN_LOG_BUY_COIN})
	}

	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"state": direct, "reson": pdata.Get("reason").ToString()}, db.DB_PARAMS{"id": pdata.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "操作信息失败"}
	}
	return &AdminResponse{State: SUCCESS, Data: "操作成功"}
}

func (m *UserModel) ListLoanOrders(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("sn").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(o.sn like  '%%%s%%' OR u.username like  '%%%s%%')", v, v))
	}
	if v := pdata.Get("circle").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.circle  = %d", v))
	}
	if v := pdata.Get("state").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("o.state = %d", v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("o.createtime between %d and  %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	count := config.GlobalDB.JoinCount(models.DB_TABLE_LOAN_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_LOAN_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{"u.username", "u.memo", "o.*"}, utils.Order(pdata.Get("orderby").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))

	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		r["parent_name"] = uinfo.ParentName
		rs = append(rs, r)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: db.DB_PARAMS{
			"count": count,
			"list":  rs,
			"state": SYSTEM_MODEL.LoanState(),
		},
	}
}

func (m *UserModel) DeleteLoanOrder(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要操作的信息!"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt()}); err == nil {
		return &AdminResponse{State: SUCCESS, Data: "删除订单成功!"}
	}
	return &AdminResponse{State: ERROR, Data: "删除订单失败!"}
}

func (m *UserModel) ReviewLoanOrder(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要操作的信息!"}
	}
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{State: ERROR, Data: "操作密码错误!"}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "当前订单不存在或者已经处理完毕!"}
	}

	up := P{
		"state":      pdata.Get("state").ToInt(),
		"reason":     pdata.Get("reason").ToString(),
		"finishtime": utils.GetNow(),
	}
	if pdata.Get("state").ToInt() == 1 {
		up["interest_time"] = utils.GetNow() + 86400
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_LOAN_ORDER, *(*map[string]interface{})(unsafe.Pointer(&up)), db.DB_PARAMS{"id": pdata.Get("id").ToInt()})
	if err == nil {
		if pdata.Get("state").ToInt() == 1 {
			models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
				Credit: one.Get("amount").ToFloat(),
				UserCoinLogInfo: models.QueueCreditLog{
					Credit:     one.Get("amount").ToFloat(),
					CoinType:   "usdt",
					CreateTime: utils.GetNow(),
				},
				UserCoinLogType: models.COIN_LOG_LORA_IN,
			})
		}
		return &AdminResponse{State: SUCCESS, Data: "审核订单信息成功!"}
	}
	return &AdminResponse{State: ERROR, Data: "审核订单信息失败!"}
}

func (m *UserModel) ListRechargeApprovals(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("r.finishtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(r.sn like '%%%s%%' OR u.username like '%%%s%%')", v, v))
	}
	if v := pdata.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("r.state = %s", v))
	}
	count := config.GlobalDB.JoinCount(models.DB_TABLE_RECHARGE_APPROVE+" as r", models.DB_TABLE_USER+" as u ", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_RECHARGE_APPROVE+" as r", models.DB_TABLE_USER+" as u ", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND")}, db.DB_FIELDS{"r.*", "u.username"}, utils.Order(pdata.Get("orderby").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))

	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		r["parent_name"] = uinfo.ParentName
		rs = append(rs, r)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  db.DB_PARAMS{"count": count, "list": rs, "state": SYSTEM_MODEL.ApproveState()},
	}
}

func (m *UserModel) GetUserAddressBalanceUsdt(address string) float64 {
	return 0
}

func (m *UserModel) ListMiningOrders(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(mo.sn like '%%%s%%' OR u.username like '%%%s%%')", v, v))
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("mo.state = %d", t.Get("state").ToInt()))
	}
	if v := t.Get("pid").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" mo.pid = %d", t.Get("pid").ToInt()))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf(" mo.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_MINING_ORDER+" as mo", models.DB_TABLE_USER+" as u", "mo.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{"mo.*", "u.user_type", "u.username"}, utils.Order(t.Get("sort", "mo.id desc, mo.state asc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_MINING_ORDER+" as mo", models.DB_TABLE_USER+" as u", "mo.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")})
	l := make([]map[string]interface{}, 0)
	for _, item := range list {
		o := make(map[string]interface{}, 0)
		item.SetInterface(&o)
		l = append(l, o)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":        l,
			"count":       count,
			"state_list":  []string{"进行中", "已完成"},
			"minner_type": map[int]string{1: "自定义", 2: "固定金额"},
			"minner_pair": SYSTEM_MODEL.MinnerPair(),
		},
	}
}

func (m *UserModel) StopMiningOrder(id int, pass string) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请指定一个要停止的订单号"}
	}
	if pass != OPERATION_PASSWORD {
		return &AdminResponse{State: ERROR, Data: "操作密码错误"}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "当前信息不存在!"}
	}
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 1, "unlocktime": utils.GetNow()}, db.DB_PARAMS{"id": one.Get("id").ToInt()}); err != nil {
		return &AdminResponse{State: ERROR, Data: "停止订单失败!"}
	}
	models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
		Credit:          one.Get("amount").ToFloat(),
		LockCredit:      0,
		VCrdit:          0,
		LockVCredit:     0,
		UserCoinLogType: models.COIN_LOG_USER_MINING_PROFIT,
		UserCoinLogInfo: models.QueueCreditLog{
			Credit:     one.Get("amount").ToFloat(),
			LockCredit: 0,
			Sn:         one.Get("sn").ToString(),
			CreateTime: utils.GetNow(),
		},
	})
	return &AdminResponse{State: SUCCESS, Data: "停止矿机成功!"}
}

func (m *UserModel) ReviewTransfer(rq P) *AdminResponse {
	pdata := rq.Ts()
	state := pdata.Get("state").ToInt()
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{State: ERROR, Data: "操作密码错误"}
	}
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "不存在该订单信息"}
	}
	if pdata.Get("money").ToFloat() == 0 && pdata.Get("direction").ToInt() == 1 && state == 1 {
		return &AdminResponse{State: ERROR, Data: "请输入用户实际到账金额"}
	}

	symbol := strings.ToLower(one.Get("coin_symbol").ToString())
	toPriceInfo := 1.0
	if symbol != "usdt" && symbol != "usdc" {
		p := models.MODEL_SYSTEM.GetLastCoinInfo(symbol + "usdt")
		toPriceInfo = p["close"].(float64)
	}
	asset := &models.Assets{Coin: symbol, Pair: symbol + "usdt", Price: toPriceInfo, Mode: 1}
	creditlog := &models.QueueCreditLog{CoinType: symbol, CreateTime: utils.GetNow()}

	if pdata.Get("direction").ToInt() == 1 && state == 1 {
		one["amount"] = pdata.Get("money")
		asset.Num = pdata.Get("money").ToFloat()
		models.MODEL_ASSETS.AddAssets(one.Get("uid").ToInt(), asset)
		creditlog.Credit = one.Get("amount").ToFloat()
	}
	if pdata.Get("direction").ToInt() == 2 && symbol != "usdt" {
		if state == 1 {
			asset.LockNum = one.Get("amount").ToFloat() * -1
		} else {
			asset.Num = one.Get("amount").ToFloat()
			asset.LockNum = one.Get("amount").ToFloat() * -1
			creditlog.Credit = one.Get("amount").ToFloat()
			creditlog.LockCredit = one.Get("amount").ToFloat() * -1
		}
		models.MODEL_ASSETS.AddAssets(one.Get("uid").ToInt(), asset)
	}
	if pdata.Get("direction").ToInt() == 2 && symbol == "usdt" {
		creditlog.Credit = one.Get("amount").ToFloat()
		if state == 1 {
			models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
				Credit:     0,
				LockCredit: one.Get("amount").ToFloat() * -1,
			})
		} else {
			creditlog.Credit = one.Get("amount").ToFloat()
			creditlog.LockCredit = one.Get("amount").ToFloat() * -1
			models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
				Credit:     one.Get("amount").ToFloat(),
				LockCredit: one.Get("amount").ToFloat() * -1,
			})
		}
	}
	models.MODEL_QUEUE.InputUserQueue(one.Get("uid").ToInt(), models.COIN_LOG_USER_RECHARGE, creditlog)
	if state == 1 && symbol == "usdt" {
		models.MODEL_QUEUE.InputTeamQueue(one.Get("uid").ToInt(), models.TEAM_LOG_WITHDRAW, models.QueueTeamLog{
			Recharge:   one.Get("amount").ToFloat(),
			CreateTime: utils.GetNow(),
		})
	}

	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_TRANSFER, db.DB_PARAMS{
		"state":      pdata.Get("state").ToInt(),
		"reson":      pdata.Get("reason").ToString(),
		"amount":     one.Get("amount").ToFloat(),
		"info":       pdata.Get("reason").ToString(),
		"finishtime": utils.GetNow(),
	}, db.DB_PARAMS{"id": one.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "操作该订单失败!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "操作订单成功!"}
}

func (m *UserModel) TransferList(rq P) *AdminResponse {
	pdata := rq.Ts()
	direct := pdata.Get("direction").ToInt()
	where := []string{fmt.Sprintf("t.direction = %d", direct)}
	if v := pdata.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("t.state = %s", v))
	}
	if v := pdata.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("t.coin_symbol = '%s'", strings.ToLower(v)))
	}
	if v := pdata.Get("user_type").ToString(); v != "" {
		where = append(where, fmt.Sprintf("u.user_type = %s", v))
	}
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(u.username like '%%%s%%' OR t.sn = '%s')", v, v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("t.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("coin_symbol").ToString(); v != "" {
		where = append(where, fmt.Sprintf("t.coin_symbol= '%s'", v))
	}

	count := config.GlobalDB.JoinCount(models.DB_TABLE_TRANSFER+" as t", models.DB_TABLE_USER+" as u  ", "t.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_TRANSFER+" as t", models.DB_TABLE_USER+" as u ", "t.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"t.*", "u.username, u.memo"}, " order by  id desc", utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))

	r := make([]map[string]interface{}, 0)
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		i["parent_name"] = uinfo.ParentName
		i["user_type"] = uinfo.UserType
		r = append(r, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"count":     count,
			"list":      r,
			"coin_list": SYSTEM_MODEL.TranferCoin(direct),
			"state":     SYSTEM_MODEL.UserStatePair(),
			"user_type": SYSTEM_MODEL.UserTypePair(),
		},
	}
}

func (m *UserModel) RechargeList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (u.username like '%%%s%%' OR r.sn like' %%%s%%' OR r.uid = '%s') ", v, v, v))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" r.cointype = '%s'", v))
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("r.state = '%s'", v))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("  r.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	count := config.GlobalDB.JoinCount("recharge as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	list, _ := config.GlobalDB.JoinTable("recharge as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"u.username", "u.credit", "u.user_type", "r.*"}, utils.Order(t.Get("sort", "r.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		a := make(map[string]interface{}, 0)
		item.SetInterface(&a)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		a["parent_name"] = uinfo.ParentName
		l = append(l, a)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":      l,
			"count":     count,
			"coin_list": SYSTEM_MODEL.CoinTypePair(),
			"state":     SYSTEM_MODEL.UserStatePair(),
		},
	}
}

func (m *UserModel) ReviewRecharge(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("id").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "确认状态失败!"
		return rs
	}
	if v := t.Get("password").ToString(); v != OPERATION_PASSWORD {
		rs.State = ERROR
		rs.Data = "操作密码错误!"
		return rs
	}
	rg := m.GetRechargeByID(t.Get("id").ToInt())
	if rg == nil {
		rs.State = ERROR
		rs.Data = "该充值信息没有找到!"
		return rs
	}
	state := t.Get("state").ToInt()
	if state == 0 {
		rs.State = ERROR
		rs.Data = "进行中的状态无法手动修改!"
		return rs
	}
	var err error
	if state == 1 {
		rechargeCredit := rg.Credit
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE, db.DB_PARAMS{"state": state, "finishtime": utils.GetNow()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
		if rg.CoinType != "USDT" {
			pair := fmt.Sprintf("%susdt", strings.ToLower(rg.CoinType))
			coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(pair)
			models.MODEL_ASSETS.AddAssets(rg.Uid, &models.Assets{
				Coin:    strings.ToLower(rg.CoinType),
				Pair:    pair,
				Num:     rg.Credit,
				LockNum: 0,
				Price:   coinPriceInfo["close"].(float64),
				Mode:    1,
			})
			rechargeCredit = 0
		}
		models.MODEL_USER.AddCredit(rg.Uid, &models.CreditValue{
			Credit:          rechargeCredit,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: models.COIN_LOG_USER_RECHARGE,
			UserCoinLogInfo: models.QueueCreditLog{
				Credit:     rg.Credit,
				LockCredit: 0,
				Sn:         rg.Sn,
				CoinType:   strings.ToLower(rg.CoinType),
				CreateTime: utils.GetNow(),
			},
			TeamCoinLogType: models.TEAM_LOG_RECHARGE,
			TeamCoinLogInfo: models.QueueTeamLog{
				Recharge:   rg.FactCredit,
				CreateTime: utils.GetNow(),
			},
		})
		models.MODEL_USER.ClearCache(rg.Uid)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE, db.DB_PARAMS{"state": 2, "reason": t.Get("reason").ToString()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作充值信息失败!"
		return rs
	}
	rs.State = SUCCESS
	rs.Data = "操作成功!"
	return rs
}

func (m *UserModel) GetRechargeByID(id int) *models.Recharge {
	if id == 0 {
		return nil
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_RECHARGE, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	rg := new(models.Recharge)
	one.SetObj(rg)
	return rg
}

func (m *UserModel) ListUserWallets(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("address").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (ua.address like '%%%s%%' OR u.username like '%%%s%%' OR u.id = '%s')", v, v, v))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" ua.cointype ='%s'", v))
	}
	if v := t.Get("contract").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" ua.contract ='%s'", v))
	}
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_USER_WITHDRAW_WALLET+" as ua ", models.DB_TABLE_USER+" AS u", " u.id = ua.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"ua.*", "u.username", "u.memo"}, utils.Order(t.Get("sort", "ua.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_USER_WITHDRAW_WALLET+" as ua ", models.DB_TABLE_USER+" AS u", " u.id = ua.uid", db.DB_PARAMS{})
	l := make([]interface{}, 0)
	for _, item := range list {
		n := make(map[string]interface{}, 0)
		item.SetInterface(&n)
		l = append(l, n)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":      l,
			"count":     count,
			"chan_type": SYSTEM_MODEL.ContractFlag(),
			"coin_type": SYSTEM_MODEL.CoinKeyValPair(),
		},
	}
}

func (m *UserModel) DeleteUserWallet(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确定有一个要删除的信息!"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"id": id}); err == nil {
		return &AdminResponse{State: ERROR, Data: "删除用户钱包成功!"}
	}
	return &AdminResponse{State: ERROR, Data: "删除用户钱包失败!"}
}

func (m *UserModel) WithdrawList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (u.username like '%%%s%%' OR u.id = '%s')", v, v))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("w.cointype=%s", v))
	}
	if v := t.Get("contract").ToString(); v != "" {
		if v == "1" {
			where = append(where, fmt.Sprintf("w.type='%s'", v))
		} else {
			where = append(where, fmt.Sprintf("w.contract=%s", v))
		}
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("w.state = %s", t.Get("state").ToString()))
	}
	if v := t.Get("sn").ToString(); v != "" {
		where = append(where, fmt.Sprintf("w.sn = %s", v))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("w.createtime between  %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_WITHDRAW+" as w", models.DB_TABLE_USER+" as u", "u.id = w.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"u.username", "u.memo", "u.user_type", "w.*"}, utils.Order(t.Get("sort", "w.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_WITHDRAW+" as w", models.DB_TABLE_USER+" as u", "u.id = w.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	lr := make([]*models.Withdraw, 0)
	for _, v := range list {
		i := new(models.Withdraw)
		v.SetObj(i)
		uinfo := models.MODEL_USER.GetBaseInfo(v.Get("uid").ToInt())
		i.ParentName = uinfo.ParentName
		lr = append(lr, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":          lr,
			"count":         count,
			"withdraw_pair": SYSTEM_MODEL.WithdrawPair(),
			"coin_list":     SYSTEM_MODEL.CoinKeyValPair(),
			"chan_type":     SYSTEM_MODEL.ContractFlag(),
		},
	}
}

func (m *UserModel) SaveWithdraw(rq P) *AdminResponse {
	t := rq.Ts()
	up := make(db.DB_PARAMS, 0)
	id := t.Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认要修改的信息"}
	}
	if v := t.Get("contract").ToString(); v != "" {
		up["contract"] = v
	}
	if v := t.Get("cointype").ToString(); v != "" {
		up["cointype"] = v
	}
	if v := t.Get("address").ToString(); v != "" {
		up["address"] = v
	}
	if v := t.Get("fact_credit").ToFloat(); v > 0 {
		up["fact_credit"] = v
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_WITHDRAW, up, db.DB_PARAMS{"id": id})
	if err == nil {
		return &AdminResponse{State: SUCCESS, Data: "修改 成功!"}
	}
	return &AdminResponse{State: ERROR, Data: "修改失败!"}
}

func (m *UserModel) ReviewWithdraw(id int, state int, info string, password string) *AdminResponse {
	rs := new(AdminResponse)
	if id == 0 {
		rs.State = ERROR
		rs.Data = "请确认一个要处理的请求!"
		return rs
	}
	if password != OPERATION_PASSWORD {
		rs.State = ERROR
		rs.Data = "操作密码不正确"
		return rs
	}
	one := m.GetWithdrawByID(id)
	if one.State != 0 {
		rs.State = ERROR
		rs.Data = "该提现请求已经处理"
		return rs
	}

	up := db.DB_PARAMS{"state": state, "finishtime": utils.GetNow(), "info": info}
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_WITHDRAW, up, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "处理提现信息失败!"}
	}
	if state == 1 {
		models.MODEL_USER.AddCredit(one.Uid, &models.CreditValue{
			Credit:          0,
			LockCredit:      -1 * one.Credit,
			VCrdit:          0,
			LockVCredit:     0,
			TeamCoinLogType: models.TEAM_LOG_WITHDRAW,
			TeamCoinLogInfo: models.QueueTeamLog{
				Withdraw:   one.Credit,
				CreateTime: utils.GetNow(),
			},
		})
	} else {
		models.MODEL_USER.AddCredit(one.Uid, &models.CreditValue{
			Credit:          one.Credit,
			LockCredit:      -1 * one.Credit,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: models.COIN_LOG_USER_WITHDRAW_FAILD,
			UserCoinLogInfo: models.QueueCreditLog{
				Credit:     one.Credit,
				LockCredit: -1 * one.Credit,
				Sn:         one.Sn,
				CreateTime: utils.GetNow(),
			},
		})
	}
	return &AdminResponse{State: SUCCESS, Data: "处理提现信息成功!"}
}

func (m *UserModel) GetWithdrawByID(id int) *models.Withdraw {
	if id == 0 {
		return nil
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_WITHDRAW, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	r := new(models.Withdraw)
	one.SetObj(r)
	return r
}
