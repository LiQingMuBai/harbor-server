package service

import (
	"cointrade/config"
	userdomain "cointrade/internal/domain/user"
	useridentityrepo "cointrade/internal/useridentity/repo"
	useridentityservice "cointrade/internal/useridentity/service"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

var adminUserIdentitySvc = useridentityservice.NewService(
	useridentityrepo.NewDBRepository(),
	adminUserIdentityUserGateway{},
	adminUserIdentityNotifier{},
)

type adminUserIdentityUserGateway struct{}

func (adminUserIdentityUserGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return models.MODEL_USER.GetBaseInfo(uid)
}

func (adminUserIdentityUserGateway) Update(uid int, data db.DB_PARAMS) {
	models.MODEL_USER.Update(uid, data)
}

type adminUserIdentityNotifier struct{}

func (adminUserIdentityNotifier) IncrementNotify(typ int, num int) {
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: typ, Num: num})
}

func (m *UserModel) SaveAuth(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{State: ERROR, Data: "请填写一个用户id"}
	}
	if len(t.Get("realname").ToString()) < 4 || len(t.Get("realname").ToString()) > 40 {
		return &AdminResponse{State: ERROR, Data: "真实姓名过长"}
	}
	if t.Get("card_front").ToString() == "" || t.Get("card_back").ToString() == "" || t.Get("card_hand").ToString() == "" {
		return &AdminResponse{State: ERROR, Data: "证件照不全"}
	}
	if err := adminUserIdentitySvc.AdminSaveLv1(
		t.Get("uid").ToInt(),
		t.Get("realname").ToString(),
		t.Get("inid").ToString(),
		t.Get("card_front").ToString(),
		t.Get("card_back").ToString(),
		t.Get("card_hand").ToString(),
	); err != nil {
		switch err {
		case useridentityservice.ErrUserNotFound:
			return &AdminResponse{State: ERROR, Data: "该用户不存在 "}
		case useridentityservice.ErrAuthAlreadyExists:
			return &AdminResponse{State: ERROR, Data: "当前用户的认证信息已经存在"}
		default:
			return &AdminResponse{State: ERROR, Data: "添加用户认证失败"}
		}
	}
	return &AdminResponse{State: SUCCESS, Data: "添加用户认证成功"}
}

func (m *UserModel) DeleteUserAuth(id int, tp int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的认证信息!"}
	}
	if err := adminUserIdentitySvc.AdminDeleteAuth(id, tp); err != nil {
		switch err {
		case useridentityservice.ErrAuthNotFound:
			return &AdminResponse{State: ERROR, Data: "当前认证信息不存在!"}
		default:
			return &AdminResponse{State: SUCCESS, Data: "删除认证信息失败!"}
		}
	}
	return &AdminResponse{State: SUCCESS, Data: "删除认证信息成功!"}
}

func (m *UserModel) ReviewUserAuth(rq P, tp int) *AdminResponse {
	t := rq.Ts()
	id := t.Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "审核的id不存在!"}
	}
	if err := adminUserIdentitySvc.AdminReviewAuth(id, tp, t.Get("process_state").ToInt(), t.Get("reason").ToString()); err != nil {
		switch err {
		case useridentityservice.ErrAuthNotFound:
			return &AdminResponse{State: ERROR, Data: "当前审核信息不存在!"}
		default:
			return &AdminResponse{State: ERROR, Data: "审核当前信息失败！"}
		}
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

func (m *UserModel) ListAdvancedUserAuth(rq P) *AdminResponse {
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
