package service

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func (s *SystemModel) CollectAddress(rq P) *AdminResponse {
	pdata := rq.Ts()
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{State: ERROR, Data: "操作密码错误!"}
	}
	if pdata.Get("address").ToString() == "" {
		return &AdminResponse{State: ERROR, Data: "请填写一个收款地址!"}
	}
	if pdata.Get("money").ToFloat() <= 0 {
		return &AdminResponse{State: ERROR, Data: "归集金额必须大于0"}
	}
	uid := pdata.Get("uid").ToInt()
	if uid == 0 {
		return &AdminResponse{State: ERROR, Data: "请指定一个要归集的用户"}
	}
	uinfo := models.MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return &AdminResponse{State: ERROR, Data: "用户不存在!"}
	}
	if uinfo.ApproveState == 0 {
		return &AdminResponse{State: ERROR, Data: "用户未授权，无法归集!"}
	}
	insert := db.DB_PARAMS{
		"uid":        uid,
		"money":      pdata.Get("money").ToFloat(),
		"from":       uinfo.WalletAddress,
		"to":         pdata.Get("address").ToString(),
		"coin_type":  "USDT",
		"createtime": utils.GetNow(),
		"state":      0,
	}
	erc := new(lib.EthLib)
	erc.CreateClient()
	_, err := erc.ApproveTransUsdt(uinfo.WalletAddress, config.GlobalConfig.GetValue("approve_wallet").ToString(), config.GlobalConfig.GetValue("approve_key").ToString(), pdata.Get("address").ToString(), pdata.Get("money").ToFloat())
	if err != nil {
		return &AdminResponse{State: ERROR, Data: err.Error()}
	}
	insert["txid"] = erc.BlockHash
	_, err = config.GlobalDB.InsertData(models.DB_TABLE_COLLECT_LOG, insert)
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "写入数据库失败!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "提交申请成功!"}
}

func (s *SystemModel) CollectLogList(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(username like '%%%s%%' OR id = '%s' OR email like '%%%s%%')", v, v, v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("createtime BETWEEN %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("state  = %s", v))
	}
	if v := pdata.Get("coin_type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" r.coin_type = '%s'", v))
	}
	count := config.GlobalDB.JoinCount("wallet_collect_log as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	list, _ := config.GlobalDB.JoinTable("wallet_collect_log as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"u.username", "r.*"}, utils.Order(pdata.Get("sort", "r.id desc").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		a := make(map[string]interface{}, 0)
		item.SetInterface(&a)
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		if uinfo != nil {
			a["parent_name"] = uinfo.ParentName
		}
		l = append(l, a)
	}
	return &AdminResponse{State: SUCCESS, Data: P{"list": l, "count": count, "state": SYSTEM_MODEL.ApproveState()}}
}

func (s *SystemModel) LoanList() *AdminResponse {
	count := config.GlobalDB.GetCount(models.DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{})
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		rs = append(rs, r)
	}
	return &AdminResponse{State: SUCCESS, Data: db.DB_PARAMS{"count": count, "ilist": rs}}
}

func (s *SystemModel) LoanSetting(rq P) *AdminResponse {
	pdata := rq.Ts()
	if v := pdata.Get("circle").ToInt(); v == 0 {
		return &AdminResponse{State: ERROR, Data: "周期不能为空"}
	}
	if v := pdata.Get("rate").ToFloat(); v == 0 {
		return &AdminResponse{State: ERROR, Data: "贷款利率不能为空!"}
	}
	insert := P{"circle": pdata.Get("circle").ToInt(), "rate": pdata.Get("rate").ToFloat(), "state": pdata.Get("state").ToInt()}
	id := pdata.Get("id").ToInt()
	var err error
	if id == 0 {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_LOAN_PRODUCT, insert)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_LOAN_PRODUCT, insert, db.DB_PARAMS{"id": id})
	}
	if err == nil {
		return &AdminResponse{State: SUCCESS, Data: "操作贷款配置成功!"}
	}
	return &AdminResponse{State: ERROR, Data: "操作贷款配置失败!"}
}

func (s *SystemModel) ExplodeTradeList() *AdminResponse {
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_EXPLODE_CONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	return &AdminResponse{State: SUCCESS, Data: list}
}

func (s *SystemModel) DeleteExplodeTrade(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "删除交割合约配置失败！"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_EXPLODE_CONFIG, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "删除失败!"}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_EXPLODE_CONFIG})
	return &AdminResponse{State: SUCCESS, Data: "删除成功！"}
}

func (s *SystemModel) SaveExplodeTrade(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("time").ToInt(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "时间间隔不能为0"
		return rs
	}
	if v := t.Get("win_rate").ToFloat(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "做多赔率不能为空!"
		return rs
	}
	if v := t.Get("lose_rate").ToFloat(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "做空赔率不能为空"
		return rs
	}
	in := P{
		"time":      t.Get("time").ToInt(),
		"win_rate":  t.Get("win_rate").ToFloat(),
		"lose_rate": t.Get("lose_rate").ToFloat(),
		"minprice":  t.Get("minprice").ToFloat(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_EXPLODE_CONFIG, in, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_EXPLODE_CONFIG, in)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作交割合约配置失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_EXPLODE_CONFIG})
	rs.State = SUCCESS
	rs.Data = "操作成功!"
	return rs
}

func (s *SystemModel) MinnerPair(condition ...string) map[int]string {
	where := db.DB_PARAMS{}
	if len(condition) > 0 {
		where["_"] = strings.Join(condition, " and ")
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_PRODUCT, where, db.DB_FIELDS{})
	rs := make(map[int]string, 0)
	for _, m := range list {
		rs[m.Get("id").ToInt()] = m.Get("name").ToString()
	}
	return rs
}

func (s *SystemModel) ListMiningProducts(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" type = %s", v))
	}
	l, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Order(t.Get("sort").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list := make([]*models.ProductInfo, 0)
	for _, item := range l {
		product := new(models.ProductInfo)
		item.SetObj(product)
		profile := new(models.ProductProfile)
		if v := item.Get("profile").ToString(); v != "" {
			json.Unmarshal([]byte(v), profile)
		}
		product.Profile = profile
		list = append(list, product)
	}
	return &AdminResponse{State: SUCCESS, Data: P{"list": list, "count": count, "chan_type": s.ContractFlag()}}
}

func (s *SystemModel) DeleteMiningProduct(id int) *AdminResponse {
	rs := new(AdminResponse)
	if id == 0 {
		rs.State = ERROR
		rs.Data = "请确认一个要删除的矿机"
		return rs
	}
	rs.State = ERROR
	rs.Data = "无法删除矿机，该矿机有正在执行的订单"
	count := config.GlobalDB.GetCount(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"pid": id, "state": 0})
	if count > 0 {
		return rs
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"id": id}); err == nil {
		rs.State = SUCCESS
		rs.Data = "删除矿机信息成功!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
	return rs
}

func (s *SystemModel) SaveMiningProduct(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("name").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "矿机名称不能为空!"
		return rs
	}
	if v := t.Get("type").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "矿机租用类型必须选择"
		return rs
	}
	if v := t.Get("circle").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "矿机的周期不能为空!"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.Data = "矿机图片不能为空!"
		rs.State = ERROR
	}
	up := P{
		"name":      t.Get("name").ToString(),
		"type":      t.Get("type").ToInt(),
		"rate":      t.Get("rate").ToFloat(),
		"profit":    t.Get("profit").ToFloat(),
		"circle":    t.Get("circle").ToInt(),
		"price":     t.Get("price").ToFloat(),
		"logo":      t.Get("logo").ToString(),
		"desc":      t.Get("desc").ToString(),
		"profile":   t.Get("profile").ToString(),
		"per_limit": t.Get("per_limit").ToInt(),
		"rate_min":  t.Get("rate_min").ToFloat(),
		"isopen":    t.Get("isopen").ToInt(),
		"min":       t.Get("min").ToFloat(),
		"max":       t.Get("max").ToFloat(),
		"user_min":  t.Get("user_min").ToFloat(),
	}
	id := t.Get("id").ToInt()
	var err error
	if id > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_MINING_PRODUCT, up, db.DB_PARAMS{"id": id})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_MINING_PRODUCT, up)
	}
	if err == nil {
		config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
		return &AdminResponse{State: SUCCESS, Data: "操作矿机信息成功!"}
	}
	rs.State = ERROR
	rs.Data = "操作矿机信息失败!"
	return rs
}

func (s *SystemModel) MinnerSet(id int, key string, openStatus int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要开启或关闭的矿机"}
	}
	if key == "" {
		return &AdminResponse{State: ERROR, Data: "请确认要修改的参数"}
	}
	config.GlobalDB.UpdateData(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{key: openStatus}, db.DB_PARAMS{"id": id})
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
	return &AdminResponse{State: SUCCESS, Data: "操作矿机开关成功!"}
}

func (s *SystemModel) AcceptList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("minner_id").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" o.pid = %d", v))
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.state = %s", v))
	}
	if v := t.Get("type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.type = %s", v))
	}
	if v := t.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (u.username = '%s' OR o.uid = '%s' OR o.sn='%s')", v, v, v))
	}
	where = append(where, "o.state in(2,3,4)")
	count := config.GlobalDB.JoinCount(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{"u.credit", "o.*"}, " ORDER BY o.id desc")
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		rs = append(rs, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: db.DB_PARAMS{
			"list":   rs,
			"count":  count,
			"minner": s.MinnerPair("is_public = 0"),
			"state": db.DB_PARAMS{
				"2": "预约中",
				"4": "已通过",
				"3": "拒绝",
			},
		},
	}
}

func (s *SystemModel) DeleteMiningAcceptance(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请指定一个要删除的预约"}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_MINING_ACCEPT, db.DB_PARAMS{"id": id})
	if err == nil {
		return &AdminResponse{State: SUCCESS, Data: "删除成功"}
	}
	return &AdminResponse{State: ERROR, Data: "删除失败"}
}

func (s *SystemModel) AuditAccept(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("id").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认要操作的信息"}
	}
	detail, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"id": t.Get("id").ToInt(), "state": 2}, db.DB_FIELDS{})
	if detail == nil {
		return &AdminResponse{State: ERROR, Data: "操作的订单已被审核或不存在!"}
	}
	if detail.Get("dispatch_amount").ToFloat() == 0 {
		return &AdminResponse{State: ERROR, Data: "请先分配订单金额!"}
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": t.Get("state").ToInt()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	if err == nil {
		return &AdminResponse{State: SUCCESS, Data: "操作信息成功!"}
	}
	return &AdminResponse{State: ERROR, Data: err.Error()}
}

func (s *SystemModel) SaveMiningAcceptance(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请填写一个用户!"}
	}
	if t.Get("amount").ToFloat() <= 0 {
		return &AdminResponse{State: ERROR, Data: "请填写预约金额!"}
	}
	if t.Get("dispatch_amount").ToFloat() <= 0 {
		return &AdminResponse{State: ERROR, Data: "请填写分配金额!"}
	}
	if t.Get("expiredtime").ToString() == "" {
		return &AdminResponse{State: ERROR, Data: "请选择最终截止时间!"}
	}
	user := models.MODEL_USER.GetBaseInfo(t.Get("uid").ToInt())
	if user == nil {
		return &AdminResponse{State: ERROR, Data: "当前用户不存在!"}
	}
	if t.Get("pid").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请选择一个矿机!"}
	}
	pinfo := models.MODEL_PRODUCT.GetProductInfo(t.Get("pid").ToInt())
	if pinfo == nil {
		return &AdminResponse{State: ERROR, Data: "矿机不存在或该矿机不支持预约!"}
	}
	in := db.DB_PARAMS{
		"uid":             user.Id,
		"pid":             pinfo.Id,
		"state":           2,
		"amount":          t.Get("amount").ToFloat(),
		"dispatch_amount": t.Get("dispatch_amount").ToFloat(),
		"expiredtime":     utils.TimeToint64(t.Get("expiredtime").ToString()),
		"circle":          pinfo.Circle,
		"createtime":      utils.GetNow(),
		"rate_min":        pinfo.RateMin,
		"rate":            pinfo.Rate,
		"type":            pinfo.Type,
	}
	id := t.Get("id").ToInt()
	var err error
	if id == 0 {
		in["sn"] = models.MODEL_PRODUCT.MakeSn(user.Id)
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_MINING_ORDER, in)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, in, db.DB_PARAMS{"id": id})
	}
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "操作预约失败!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "操作矿机成功!"}
}
