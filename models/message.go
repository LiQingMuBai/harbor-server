package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"strings"
)

type MessageModel struct {
	ModelBase
}

const (
	MESSAGE_TYPE_TEXT       = 1 //文本消息
	MESSAGE_TYPE_CREDIT     = 2 //金额变动
	MESSAGE_TYPE_MSG_CANCEL = 3
)

type MessageText struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}
type MessageCredit struct {
	Credit      float64      `json:"credit"`      //用户余额
	LockCredit  float64      `json:"lockcredit"`  //用户冻结余额
	VCredit     float64      `json:"vcredit"`     //用户虚拟余额
	LockVCredit float64      `json:"lockvcredit"` //用户虚拟冻余额
	Text        *MessageText `json:"text"`        //携带的文本消息
}
type SendMessage struct {
	MsgType int         `json:"msgtype"`
	Msg     interface{} `json:"msg"`
}
type DBServiceMessage struct {
	Uid        int    `json:"uid"`     //用户ID
	Content    string `json:"content"` //内容
	SnId       string `json:"sn_id"`
	CreateTime int    `json:"createtime"` //创建时间
	Flag       int    `json:"flag"`       //发送还是收取 1 服务端收到的 2 服务端发出去的
}

func (m *MessageModel) PushMessage(uid int, content interface{}, cmd int) { //消息推入队列
	//queue_name := fmt.Sprintf("%s_%d", HASH_USER_MESSAGE, uid)
	config.GlobalRedis.PushQueue(HASH_USER_MESSAGE, map[string]interface{}{"uid": uid, "content": content, "cmd": cmd})
}
func (m *MessageModel) GetList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	//私人消息列表
	condition := db.DB_PARAMS{"uid": uid, "msg_type": MESSAGE_TYPE_TEXT}
	count := config.GlobalDB.GetCount(DB_TABLE_MESSAGE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	offset := (rq.Page - 1) * rq.Limit
	limitstr := fmt.Sprintf("limit %d,%d", offset, rq.Limit)
	list, err := config.GlobalDB.FetchRows(DB_TABLE_MESSAGE, condition, db.DB_FIELDS{}, "order by createtime desc", limitstr)
	if err == nil {
		return &PageBaseResponse{
			Total:     count,
			Page:      rq.Page,
			PageTotal: pagesize,
			Limit:     rq.Limit,
			List:      list,
		}
	}
	return nil
}
func (m *MessageModel) GetUnreadNum(uid int) int { //这里缓存后面加
	//获取私人消息未读数量
	condition := db.DB_PARAMS{"uid": uid, "state": 0}
	return config.GlobalDB.GetCount(DB_TABLE_MESSAGE, condition)
}
func (m *MessageModel) ChangeState(uid int, messageid int) bool { //这里缓存后面加
	config.GlobalDB.UpdateData(DB_TABLE_MESSAGE, db.DB_PARAMS{"state": 1, "readtime": utils.GetNow()}, db.DB_PARAMS{"uid": uid, "id": messageid})
	return true
}
func (m *MessageModel) GetServiceList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	//获得用户客服消息列表
	rq.Limit = 50
	condition := db.DB_PARAMS{"uid": uid}
	count := config.GlobalDB.GetCount(DB_TABLE_SERVICE_MESSAGE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_SERVICE_MESSAGE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := make([]map[string]string, 0)
	for _, v := range list {
		v["content"] = strings.Replace(v["content"], "\n", " ", -1)
		v["content"] = strings.Replace(v["content"], "\t", " ", -1)
		rs = append(rs, v)
	}
	return &PageBaseResponse{
		Page:      rq.Page,
		Limit:     rq.Limit,
		Total:     count,
		PageTotal: pagesize,
		List:      rs,
		BaseResponse: BaseResponse{
			State: STATE_SUCCESS,
			Msg:   "success",
		},
	}
}
func (m *MessageModel) ClearServiceUnread(uid int) *BaseResponse {
	config.GlobalDB.UpdateData(DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"read_state": 1}, db.DB_PARAMS{"uid": uid, "flag": 2, "read_state": 0})
	return &BaseResponse{
		State: STATE_SUCCESS,
		Msg:   "success!",
	}
}
func (m *MessageModel) GetServiceUnreadCount(uid int) int { //得到客服消息未读数量
	return config.GlobalDB.GetCount(DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"uid": uid, "read_state": 0, "flag": 2})
}
