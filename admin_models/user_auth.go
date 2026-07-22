package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

func (m *UserModel) SaveAuth(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请填写一个用户id"}
	}
	uinfo := models.MODEL_USER.GetBaseInfo(t.Get("uid").ToInt())
	if uinfo == nil {
		return &AdminResponse{State: ERROR, Data: "该用户不存在 "}
	}
	if exists := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH, db.DB_PARAMS{"uid": t.Get("uid").ToInt()}); exists > 0 {
		return &AdminResponse{State: ERROR, Data: "当前用户的认证信息已经存在"}
	}
	if len(t.Get("realname").ToString()) < 4 || len(t.Get("realname").ToString()) > 40 {
		return &AdminResponse{State: ERROR, Data: "真实姓名过长"}
	}
	if t.Get("card_front").ToString() == "" || t.Get("card_back").ToString() == "" || t.Get("card_hand").ToString() == "" {
		return &AdminResponse{State: ERROR, Data: "证件照不全"}
	}

	insertData := db.DB_PARAMS{
		"uid":           t.Get("uid").ToInt(),
		"realname":      t.Get("realname").ToString(),
		"inid":          t.Get("inid").ToString(),
		"card_front":    t.Get("card_front").ToString(),
		"card_back":     t.Get("card_back").ToString(),
		"card_hand":     t.Get("card_hand").ToString(),
		"process_state": 1,
		"createtime":    utils.GetNow(),
		"passtime":      utils.GetNow(),
	}
	config.GlobalDB.InsertData(models.DB_TABLE_USERAUTH, insertData)
	return &AdminResponse{State: SUCCESS, Data: "添加用户认证成功"}
}

func (m *UserModel) UauthDel(id int, tp int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的认证信息!"}
	}
	table := map[int]string{
		1: models.DB_TABLE_USERAUTH,
		2: models.DB_TABLE_USERAUTH_LV2,
	}
	one, _ := config.GlobalDB.FetchOne(table[tp], db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "当前认证信息不存在!"}
	}
	if _, err := config.GlobalDB.Delete(table[tp], db.DB_PARAMS{"id": id}); err == nil {
		models.MODEL_USER.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp - 1})
		return &AdminResponse{State: SUCCESS, Data: "删除认证信息成功!"}
	}
	return &AdminResponse{State: SUCCESS, Data: "删除认证信息失败!"}
}

func (m *UserModel) UauthOp(rq P, tp int) *AdminResponse {
	t := rq.Ts()
	id := t.Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "审核的id不存在!"}
	}
	table := map[int]string{
		1: models.DB_TABLE_USERAUTH,
		2: models.DB_TABLE_USERAUTH_LV2,
	}
	one, _ := config.GlobalDB.FetchOne(table[tp], db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{State: ERROR, Data: "当前审核信息不存在!"}
	}

	up := db.DB_PARAMS{"process_state": t.Get("process_state").ToInt(), "reason": t.Get("reason").ToString()}
	if tp == 2 {
		up = db.DB_PARAMS{"state": t.Get("process_state").ToInt(), "reason": t.Get("reason").ToString()}
	}
	up["passtime"] = utils.GetNow()
	if _, err := config.GlobalDB.UpdateData(table[tp], up, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{State: ERROR, Data: "审核当前信息失败！"}
	}
	if t.Get("process_state").ToInt() == 1 {
		usup := db.DB_PARAMS{"auth_lv": tp, "credit_coin": 80}
		if tp == 2 {
			usup["credit_coin"] = 100
		}
		models.MODEL_USER.Update(one.Get("uid").ToInt(), usup)
	}
	return &AdminResponse{State: SUCCESS, Data: "审核当前信息成功!"}
}

func (m *UserModel) UserAuthList(rq P) *AdminResponse {
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
		parent := m.GetParentUser(item.Get("uid").ToInt())
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

func (m *UserModel) GetParentUser(uid int) *models.UserBaseInfo {
	user := models.MODEL_USER.GetBaseInfo(uid)
	return user
}

func (m *UserModel) UserAuthDouble(rq P) *AdminResponse {
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
