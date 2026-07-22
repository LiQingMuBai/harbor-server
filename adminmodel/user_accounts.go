package adminmodel

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
	"unsafe"
)

func (m *UserModel) SaveParentMemo(rq P) *AdminResponse {
	t := rq.Ts()
	user := models.MODEL_USER.GetBaseInfo(t.Get("topid").ToInt())
	if user == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前用户不存在!",
		}
	}
	models.MODEL_USER.Update(t.Get("topid").ToInt(), db.DB_PARAMS{"memo": t.Get("memo").ToString()})
	models.MODEL_USER.ClearCache(t.Get("uid").ToInt())
	return &AdminResponse{
		State: SUCCESS,
		Data:  "修改成功!",
	}
}

func (m *UserModel) UserList(rq P) *AdminResponse {
	where := make([]string, 0)
	pdata := rq.Ts()
	if v := pdata.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(username ='%s' OR id = '%s' OR wallet_address ='%s' OR memo like '%s')", v, v, v, v))
	}
	if v := pdata.Get("invite_code").ToString(); v != "" {
		where = append(where, fmt.Sprintf("invite_code  = %s", v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("createtime BETWEEN %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("isagent").ToString(); v != "" {
		where = append(where, fmt.Sprintf("is_agent  = %s", v))
	}
	if v := pdata.Get("user_type").ToString(); v != "" {
		where = append(where, fmt.Sprintf("user_type  = %s", v))
	}
	if v := pdata.Get("approve_state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("approve_state  = %s", v))
	}
	if v := pdata.Get("iswithdraw").ToString(); v != "" {
		where = append(where, fmt.Sprintf("iswithdraw  = %s", v))
	}
	if v := pdata.Get("status").ToString(); v != "" {
		where = append(where, fmt.Sprintf("status  = %s", v))
	}
	if v := pdata.Get("online").ToString(); v != "" {
		where = append(where, fmt.Sprintf("online = %d", pdata.Get("online").ToInt()))
	}
	if v := pdata.Get("approve_state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("approve_state  = %s", v))
	}

	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_USER, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Order(rq.Ts().Get("sort", " online desc,id desc").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	for k, v := range list {
		uid := db.DBValue{Value: v["id"]}
		if v["createip"] != "" {
			createip := db.DBValue{Value: v["createip"]}
			v["createip"] = utils.Long2Ip(createip.ToInt())
		}
		if v["loginip"] != "" {
			loginip := db.DBValue{Value: v["loginip"]}
			v["loginip"] = utils.Long2Ip(loginip.ToInt())
		}
		uinfo := models.MODEL_USER.GetBaseInfo(uid.ToInt())
		bank, _ := config.GlobalDB.FetchRow(models.DB_TABLE_BANKINFO, db.DB_PARAMS{"uid": uid.ToInt()}, db.DB_FIELDS{})
		if len(bank) > 0 {
			for bk, bv := range bank {
				v[bk] = bv
			}
		}
		v["parent_name"] = uinfo.ParentName
		list[k] = v
	}

	assetList := SYSTEM_MODEL.CoinKeyValPair()
	assetList[100] = "usdt"
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":             list,
			"count":            count,
			"auth_level":       SYSTEM_MODEL.UserAuthLevel(),
			"user_mode":        SYSTEM_MODEL.UserTypePair(),
			"online_state":     SYSTEM_MODEL.OnlinePair(),
			"withdrawState":    SYSTEM_MODEL.WithdrawStatus(),
			"userState":        SYSTEM_MODEL.UserStatus(),
			"assetlist":        assetList,
			"controller_state": SYSTEM_MODEL.ControllerState(),
		},
	}
}

func (m *UserModel) UserCoinLog(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个用户!",
		}
	}
	where = append(where, fmt.Sprintf("uid = '%d'", t.Get("uid").ToInt()))
	if v := t.Get("type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("type = '%d'", t.Get("type").ToInt()))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("cointype = '%s'", strings.ToLower(t.Get("cointype").ToString())))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_CREDIT_LOG, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Order(t.Get("order", " id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_CREDIT_LOG, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	coinlist := SYSTEM_MODEL.CoinKeyValPair()
	coinlist[100] = "usdt"
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":        list,
			"count":       count,
			"coinlist":    coinlist,
			"change_type": SYSTEM_MODEL.COINLOG_TYPELIST(),
		},
	}
}

func (m *UserModel) AdjustUserCredit(rq P) *AdminResponse {
	t := rq.Ts()
	assetname := strings.ToLower(t.Get("assetname").ToString())
	if assetname == "" {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "添加失败",
		}
	}

	if v := t.Get("id").ToInt(); v > 0 {
		coin := t.Get("coin").ToFloat()
		ntime := utils.GetNow()
		if assetname != "usdt" {
			pair := fmt.Sprintf("%susdt", assetname)
			coinPriceInfo := models.MODEL_SYSTEM.GetLastCoinInfo(pair)
			price := 0.0
			if coinPriceInfo == nil {
				coininfo := models.MODEL_SYSTEM.GetCoinInfo(assetname, pair)
				price = utils.GetFloat(coininfo["f_price"])
			} else {
				price = coinPriceInfo["close"].(float64)
			}
			if models.MODEL_ASSETS.AddAssets(t.Get("id").ToInt(), &models.Assets{
				Coin:    assetname,
				Pair:    pair,
				Num:     coin,
				LockNum: 0,
				Price:   price,
				Mode:    1,
			}) {
				models.MODEL_USER.AddCredit(t.Get("id").ToInt(), &models.CreditValue{
					UserCoinLogType: models.COIN_LOG_BACKEND,
					UserCoinLogInfo: models.QueueCreditLog{
						Credit:     coin,
						CoinType:   assetname,
						CreateTime: ntime,
					},
				})
			}
		} else {
			models.MODEL_USER.AddCredit(v, &models.CreditValue{
				Credit:          coin,
				UserCoinLogType: models.COIN_LOG_BACKEND,
				UserCoinLogInfo: models.QueueCreditLog{
					Credit:     t.Get("coin").ToFloat(),
					CoinType:   assetname,
					CreateTime: ntime,
				},
				TeamCoinLogType: models.TEAM_LOG_RECHARGE,
				TeamCoinLogInfo: models.QueueTeamLog{
					Recharge:   coin,
					CreateTime: ntime,
				},
			})
		}
	}

	return &AdminResponse{
		State: SUCCESS,
		Data:  "修改成功",
	}
}

func (m *UserModel) UserControllerExp(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "需要指定一个控制的用户",
		}
	}
	models.MODEL_USER.Update(t.Get("id").ToInt(), db.DB_PARAMS{"explode_state": t.Get("explode_state").ToInt()})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "修改设置成功!",
	}
}

func (m *UserModel) SaveUser(rq P) *AdminResponse {
	t := rq.Ts()
	up := make(P)
	if v := t.Get("auth_lv").ToInt(); v > 0 {
		up["auth_lv"] = v
	}
	if v := t.Get("nickname").ToString(); v != "" {
		up["nickname"] = v
	}
	if v := t.Get("avatar").ToString(); v != "" {
		up["avatar"] = v
	}
	if v := t.Get("password").ToString(); v != "" {
		up["password"] = utils.Md5(fmt.Sprintf("%s%s", models.PASSMIX, v))
	}
	if v := t.Get("cash_password").ToString(); v != "" {
		up["cash_password"] = utils.Md5(v)
	}
	if v := t.Get("google_serect").ToString(); v != "" {
		up["google_serect"] = ""
	}
	if v := t.Get("mode").ToInt(); v > 0 {
		up["mode"] = v
	}
	if v := t.Get("credit_coin").ToFloat(); v > 0 {
		up["credit_coin"] = v
	}
	if v := t.Get("user_type").ToInt(); v > 0 {
		up["user_type"] = v
	}
	if v := t.Get("is_agent").ToString(); v != "" {
		up["is_agent"] = v
	}
	if v := t.Get("iswithdraw").ToString(); v != "" {
		up["iswithdraw"] = v
	}
	if v := t.Get("status").ToString(); v != "" {
		up["status"] = v
	}
	if v := t.Get("email").ToString(); v != "" {
		up["email"] = v
	}
	if v := t.Get("memo").ToString(); v != "" {
		up["memo"] = v
	}
	if v := t.Get("memo").ToString(); v == "" {
		up["memo"] = ""
	}
	if v := t.Get("avatar").ToString(); v != "" {
		up["avatar"] = v
	}
	if v := t.Get("withdraw_msg").ToString(); v != "" {
		up["withdraw_msg"] = v
	}

	bank := make(map[string]interface{}, 0)
	if v := t.Get("bankname").ToString(); v != "" {
		bank["bankname"] = v
	}
	if v := t.Get("realname").ToString(); v != "" {
		bank["realname"] = v
	}
	if v := t.Get("account").ToString(); v != "" {
		bank["account"] = v
	}
	if v := t.Get("router_num").ToString(); v != "" {
		bank["router_num"] = v
	}
	if v := t.Get("swift_code").ToString(); v != "" {
		bank["swift_code"] = v
	}
	if v := t.Get("bank_address").ToString(); v != "" {
		bank["bank_address"] = v
	}

	rs := new(AdminResponse)
	if t.Get("avatar").ToString() == "" {
		up["avatar"] = SYSTEM_MODEL.SettingGet("avatar").ToString()
	}
	inviteCode := t.Get("invite_code").ToString()
	if inviteCode != "" {
		if exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"invite_code": inviteCode}); exists > 0 {
			rs.State = ERROR
			rs.Data = "该邀请码已经被使用"
			return rs
		}
	}

	username := t.Get("username").ToString()
	if username == "" {
		rs.State = ERROR
		rs.Data = "用户名不能为空!"
		return rs
	}
	if len(up) == 0 {
		rs.State = ERROR
		rs.Data = "请输入一个要操作的用户信息!"
	}

	var err error
	var uid int
	up["username"] = username
	if v := t.Get("id").ToInt(); v == 0 {
		if inviteCode == "" {
			inviteCode = models.MODEL_USER.GetInvateCode()
		}
		up["createtime"] = utils.GetNow()
		if exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"username": username}); exists > 0 {
			rs.State = ERROR
			rs.Data = "该用户已经存在!"
			return rs
		}
		up["invite_code"] = inviteCode
		lastID, _ := config.GlobalDB.InsertData(models.DB_TABLE_USER, up)
		if lastID == 0 {
			rs.State = ERROR
			rs.Data = "添加用户失败"
			return rs
		}
		uid = int(lastID)
	} else {
		user := models.MODEL_USER.GetBaseInfo(t.Get("id").ToInt())
		if user == nil {
			rs.State = ERROR
			rs.Data = "该用户不存在!"
			return rs
		}
		if user.UserName != username {
			if exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"username": username}); exists > 0 {
				rs.State = ERROR
				rs.Data = "该用户已经存在!!!!"
				return rs
			}
		}
		aa := (*db.DB_PARAMS)(unsafe.Pointer(&up))
		models.MODEL_USER.Update(t.Get("id").ToInt(), *aa)
		uid = t.Get("id").ToInt()
	}

	if len(bank) > 0 {
		if !m.UpdateUserBankInfo(uid, bank) {
			rs.State = ERROR
			rs.Data = "操作用户银行卡信息失败"
		}
	}

	toptype := t.Get("type").ToInt()
	if toptype == 3 {
		models.MODEL_USER.Update(uid, db.DB_PARAMS{"parent_order": "", "parent_uid": 0})
		config.GlobalDB.Delete(models.DB_TABLE_USER_LEVEL, db.DB_PARAMS{"uid": uid})
		rs.State = SUCCESS
		rs.Data = "格式化用户上级成功!"
		return rs
	}
	if v := t.Get("parent_code").ToString(); v != "" {
		var topid int
		if toptype == 1 {
			topid = models.MODEL_USER.GetInviteUser(v)
		}
		if toptype == 2 {
			topid = t.Get("parent_code").ToInt()
		}
		if topid == uid {
			rs.State = ERROR
			rs.Data = "自己不能成為自己的上級\n推薦人邀請碼/ID錯誤!"
			return rs
		}
		if topid > 0 {
			m.UpTopInfo(uid, topid)
		} else {
			rs.State = ERROR
			rs.Data = "推荐人邀请码不存在\n 错误!"
			return rs
		}
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作用户信息错误!"
		return rs
	}
	rs.State = SUCCESS
	rs.Data = "操作用户信息成功!"
	return rs
}

func (m *UserModel) UpdateUserBankInfo(uid int, bank map[string]interface{}) bool {
	if uid == 0 || len(bank) == 0 {
		return false
	}
	cacheID := models.MODEL_USER.MakeCacheId("bank", uid)
	var err error
	if models.MODEL_CREDIT.GetBankInfo(uid) == nil {
		bank["uid"] = uid
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_BANKINFO, bank)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_BANKINFO, bank, db.DB_PARAMS{"uid": uid})
	}
	if err != nil {
		return false
	}
	config.GlobalRedis.Del(models.HASH_USER_BANK, cacheID)
	return true
}

func (m *UserModel) UpTopInfo(uid int, topid int) {
	top := models.MODEL_USER.GetBaseInfo(topid)
	if top == nil {
		return
	}

	var order string
	if top.ParentOrder != "" && top.ParentOrder != "0" {
		order = fmt.Sprintf("%d,%d", top.Id, topid)
	} else {
		order = fmt.Sprintf("%d", topid)
	}

	tmp := strings.Split(order, ",")
	if len(tmp) > 3 {
		tmp = tmp[len(tmp)-3:]
	}

	r := db.DB_PARAMS{"parent_order": strings.Join(tmp, ","), "parent_uid": topid}
	if top.ChaneelId != "0" && top.ChaneelId != "" {
		topChannel := models.MODEL_USER.GetBaseInfo(utils.GetInt(top.ChaneelId))
		if topChannel != nil {
			r["channel_id"] = top.ChaneelId
			r["channel_username"] = topChannel.Email
			r["channel_level"] = top.Level + 1
		}
	} else if top.IsAgent == 1 {
		r["channel_id"] = top.Id
		r["channel_username"] = top.Email
		r["channel_level"] = 1
	}

	models.MODEL_USER.Update(uid, r)
	config.GlobalDB.Delete(models.DB_TABLE_USER_LEVEL, db.DB_PARAMS{"uid": uid})
	j := 1
	levelOrder := make([]string, 0)
	for l := len(tmp) - 1; l >= 0; l-- {
		levelOrder = append(levelOrder, tmp[l])
		config.GlobalDB.InsertData(models.DB_TABLE_USER_LEVEL, db.DB_PARAMS{
			"puid":        tmp[l],
			"uid":         uid,
			"level":       j,
			"levle_order": strings.Join(levelOrder, ","),
		})
		j++
	}
}

func (m *UserModel) UserAssetByDirect(uid int) P {
	assetlist, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	p := make(P)
	for _, asset := range assetlist {
		p[asset.Get("coin_symbol").ToString()] = asset.Get("coin_symbol").ToFloat()
	}
	return p
}

func (m *UserModel) UpdateUserAssetWallet(rq P) *AdminResponse {
	t := rq.Ts()
	if v := t.Get("id").ToInt(); v == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个要修改的信息",
		}
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_USERASSETS, db.DB_PARAMS{"wallet_address": t.Get("wallet_address").ToString()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "修改用户钱包信息失败!",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "修改用户钱包信息成功",
	}
}

func (m *UserModel) UserAssetList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(u.username = '%s' or u.id = '%s' or c.wallet_address like '%%%s%%') ", v, v, v))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("c.coin_symbol = '%s'", v))
	}
	if v := t.Get("is_trans").ToString(); v != "" {
		where = append(where, fmt.Sprintf("c.is_trans = '%s'", v))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_USERASSETS+" as c", models.DB_TABLE_USER+" as u ", " c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"u.username", "u.memo", "c.*"}, utils.Order(t.Get("sort", " c.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_USERASSETS+" as c", models.DB_TABLE_USER+" as u ", " c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		rs = append(rs, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  rs,
			"count": count,
			"coin":  SYSTEM_MODEL.CoinKeyValPair(),
		},
	}
}

func (m *UserModel) UserProfitLog(uid int) *AdminResponse {
	if uid == 0 {
		return &AdminResponse{State: ERROR}
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_PROFIT_LOG, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	count := config.GlobalDB.GetCount(models.DB_TABLE_PROFIT_LOG, db.DB_PARAMS{"uid": uid})
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  list,
			"count": count,
		},
	}
}

func (m *UserModel) UserTeamLevelCount(rq P) *AdminResponse {
	t := rq.Ts()
	uid := t.Get("uid").ToInt()
	if uid == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个查看的用户",
		}
	}

	where := make([]string, 0)
	subWhere := ""
	if v := t.Get("level").ToInt(); v > 0 {
		subWhere = fmt.Sprintf("  and ul.level = %d", v)
	}
	if v := t.Get("email").ToString(); v != "" {
		subWhere += fmt.Sprintf("  and u.username like '%%%s%%'", v)
	}
	count := config.GlobalDB.GetCount("`users` AS u , `user_levels` AS ul", db.DB_PARAMS{"_": "u.id =ul.uid and ul.puid=" + t.Get("uid").ToString() + " " + subWhere})
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("uc.daytime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.JoinTable("(SELECT u.id,u.username,u.memo, u.createtime,ul.level FROM `users` AS u , `user_levels` AS ul WHERE u.id =ul.uid and ul.puid="+t.Get("uid").ToString()+" "+subWhere+" ) AS ut", "`user_count` AS uc", "uc.uid = ut.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"ANY_VALUE(ut.id) as id",
		"ANY_VALUE(ut.username) username",
		"ANY_VALUE(ut.createtime) createtime",
		"ANY_VALUE(ut.level) level", "SUM(uc.recharge) as recharge, SUM(uc.withdraw) as withdraw, SUM(uc.trade) as trade, SUM(uc.trade_profit) trade_profit, SUM(uc.mining_count) mining_count, SUM(uc.mining_profit) mining_profit",
	}, " group by ut.id")

	l := make([]db.DB_PARAMS, 0)
	for _, item := range list {
		it := make(map[string]interface{}, 0)
		item.SetInterface(&it)
		l = append(l, it)
	}

	total := SYSTEM_MODEL.UserLevelListCount(rq)
	tss := make(map[string]interface{})
	total.SetInterface(&tss)
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  l,
			"count": count,
			"total": tss,
		},
	}
}

func (m *UserModel) GetUserAssetOverview(uid int) *AdminResponse {
	uinfo := models.MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return &AdminResponse{State: ERROR, Data: "用户不能存在!"}
	}

	asset, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"(num+lock_num) as amount, ANY_VALUE(coin_pair) as coin_pair"})
	minnermoney, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "uid": uid}, db.DB_FIELDS{"SUM(amount) as amount"})
	contract, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "2", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	explode, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "1", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})

	var assetTotal float64
	for _, assetItem := range asset {
		itemPrice := models.MODEL_SYSTEM.GetLastCoinInfo(assetItem.Get("coin_pair").ToString())
		if itemPrice != nil && assetItem.Get("amount").ToFloat() > 0 {
			assetTotal += assetItem.Get("amount").ToFloat() * itemPrice["close"].(float64)
		}
	}

	total := P{
		"asset":    assetTotal,
		"minner":   minnermoney.Get("amount").ToFloat(),
		"contract": contract.Get("amount").ToFloat(),
		"explode":  explode.Get("amount").ToFloat(),
		"lock":     uinfo.LockCredit,
		"valid":    uinfo.Credit,
	}
	rs := make([]P, 0)
	totalMoney := 0.00
	for key, val := range total {
		rs = append(rs, P{"project": key, "asset": val})
		totalMoney += val.(float64)
	}
	rs = append(rs, P{"project": "total", "asset": totalMoney})
	return &AdminResponse{State: SUCCESS, Data: P{
		"asset": rs,
		"project": P{
			"asset":    "钱包折合资产",
			"minner":   "矿机",
			"contract": "永续",
			"explode":  "交割",
			"lock":     "冻结",
			"valid":    "可用",
			"total":    "总计",
		},
	}}
}
