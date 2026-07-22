package adminmodel

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func (m *UserModel) MsgList(rq P) *AdminResponse {
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

func (m *UserModel) CustomServiceList(rq P) *AdminResponse {
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

func (m *UserModel) UserByMessage(uid int) *AdminResponse {
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
			}
		}
	}
	return &AdminResponse{State: SUCCESS, Data: rs}
}

func (m *UserModel) DelUserNotice(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的信息"}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_USER_NOTICE, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "删除信息失败"}
	}
	return &AdminResponse{State: SUCCESS, Data: "删除信息成功"}
}

func (m *UserModel) SendUserNotice(msg *UserNoticeMsg) *AdminResponse {
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
		return &AdminResponse{State: ERROR, Data: "发送通知失败!" + err.Error()}
	}
	return &AdminResponse{State: SUCCESS, Data: "发送个人通知成功"}
}

func (m *UserModel) SendMsg(msg *CustomMsg) *AdminResponse {
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

func (m *UserModel) DelCustomMsg(uid int) *AdminResponse {
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
