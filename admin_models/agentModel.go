package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"strings"
)

var MODEL_AGENT AgentModel

type AgentModel struct{}

/**
 *	代理列表
 */
func (a *AgentModel) AgentList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" u.username = '%s'", v))
	}
	where = append(where, "u.is_agent = 1")
	l, _ := config.GlobalDB.FetchAll(models.DB_TABLE_USER+" as u", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		fmt.Sprintf("*, ( SELECT COUNT(*) from %s where user_type = 2 and channel_id = u.id) AS employer ", models.DB_TABLE_USER),
		fmt.Sprintf("( SELECT COUNT(*) from %s where user_type = 1 and channel_id = u.id) AS custom ", models.DB_TABLE_USER),
	})
	count := config.GlobalDB.GetCount(models.DB_TABLE_USER+" as u", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	agentlist := make([]map[string]interface{}, 0)
	for _, agent := range l {
		u := make(map[string]interface{}, 0)
		agent.SetInterface(&u)
		agentlist = append(agentlist, u)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":      agentlist,
			"userState": SYSTEM_MODEL.UserStatus(),
			"count":     count,
		},
	}
}

/***
 *	代理统计
 */
func (a *AgentModel) AgentCountList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("c.daytime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	if v := t.Get("uid").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("c.uid = '%d'", v))
	}
	groupby := " GROUP BY c.uid"
	countBy := "COUNT(DISTINCT(c.uid)) as count"
	if v := t.Get("uid").ToInt(); v > 0 {
		groupby = "GROUP BY c.daytime"
		countBy = "COUNT(*) AS count"
	}
	if v := t.Get("sum").ToInt(); v > 0 {
		groupby = ""
	}
	rl, _ := config.GlobalDB.JoinTable(models.DB_TABLE_AGENT_COUNT+" as c", models.DB_TABLE_USER+" as u", "u.id = c.uid", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"ANY_VALUE(u.username) as username",
		"ANY_VALUE(c.daytime) AS daytime",
		"SUM(c.recharge) AS recharge",
		"SUM(c.withdraw) as withdraw",
		"SUM(c.trade) AS trade",
		"SUM(c.trade_profit) as trade_profit",
		"SUM(c.minning_count) as mining",
		"SUM(c.minning_profit) as mining_profit",
		"SUM(c.register_num) as register_num",
		"SUM(c.pro_register_num) as pro_register_num"}, groupby, utils.Order(t.Get("order", " daytime desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	where = append(where, "1")
	count, _ := config.GlobalDB.FetchOne(models.DB_TABLE_AGENT_COUNT+" as c", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{countBy})

	rs := make([]map[string]interface{}, 0)

	for _, item := range rl {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		rs = append(rs, r)
	}
	re := db.DB_PARAMS{"list": rs, "count": count.Get("count").ToInt()}

	if v := t.Get("sum").ToInt(); v == 0 {
		rq["sum"] = 1
		total := a.AgentCountList(rq)
		countre := total.Data.(db.DB_PARAMS)

		re["total"] = countre["list"]
	}

	return &AdminResponse{
		State: SUCCESS,
		Data:  re,
	}
}

func (a *AgentModel) EmpoyerList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("agent_id").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf("u.channel_id = %d", v))
	}
	if v := t.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf("u.username like '%%%s%%'", v))
	}
	selecttype := 1
	custom_field := " 1"
	if v := t.Get("user_type").ToInt(); v > 0 {
		selecttype = v
	}
	if selecttype == 1 {
		custom_field = fmt.Sprintf("(select count(*)  from %s where parent_uid = u.id) as custom", models.DB_TABLE_USER)
	}
	where = append(where, fmt.Sprintf(" user_type = %d", selecttype))
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_USER+" as u", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"*", custom_field})
	count := config.GlobalDB.GetCount(models.DB_TABLE_USER+" as u", db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":          list,
			"count":         count,
			"auth_level":    SYSTEM_MODEL.UserAuthLevel(),
			"user_mode":     SYSTEM_MODEL.UserTypePair(),
			"online_state":  SYSTEM_MODEL.OnlinePair(),
			"withdrawState": SYSTEM_MODEL.WithdrawStatus(),
			"userState":     SYSTEM_MODEL.UserStatus(),
		},
	}
}

func (a *AgentModel) DelAgent(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个删除的代理信息",
		}
	}
	count := config.GlobalDB.GetCount(models.DB_TABLE_USER, db.DB_PARAMS{"channel_id": id})
	if count > 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前代理下有员工或者客户信息 无法删除!",
		}
	}
	config.GlobalDB.Delete(models.DB_TABLE_USER, db.DB_PARAMS{"id": id, "is_agent": 1})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除成功!",
	}
}
