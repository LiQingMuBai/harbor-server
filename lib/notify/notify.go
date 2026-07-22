package notify

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"encoding/json"
	"fmt"
)

var NOTIFY Notify

type ServiceMsg struct {
	Uid        int    `json:"uid"`        //发送者
	Msg        string `json:"msg"`        //消息内容
	SendName   string `json:"sendname"`   //发送者名字
	Type       string `json:"type"`       //消息类型
	CreateTime int    `json:"createtime"` //消息创建时间
}

type NotifyItem struct {
	Type  int         `json:"type"`
	Title string      `json:"title"`
	Msg   string      `json:"msg"`
	Music int         `json:"music"`
	Md5   string      `json:"md5"`
	Num   int         `json:"num"`  //消息数量
	Datad *ServiceMsg `json:"data"` //详细信息
}

type Notify struct{}

/**
 *	格式化消息
 */
func (s *Notify) FormatNotify(notify *NotifyItem) *NotifyItem {
	var message map[int]string = map[int]string{
		1: "有%d条新的提现待处理",
		2: "有%d条新的充值待确认",
		3: "有%d条新的初级认证待审核",
		4: "有%d条新的高级认证待审核",
		5: "%s 发来条新的消息",
	}

	if str, ok := message[notify.Type]; ok {
		var msg string
		if notify.Type == 5 {
			if notify.Datad.Type == "pic" {
				msg = fmt.Sprintf("<img src='%s' width='50'>", notify.Datad.Msg)
			} else {
				msg = fmt.Sprintf("<i>消息内容</i>: %s", notify.Datad.Msg)
			}
			notify.Msg = msg
			notify.Music = 2
			notify.Title = fmt.Sprintf(str, notify.Datad.SendName)
		} else {
			notify.Title = fmt.Sprintf(str, notify.Num)
			notify.Music = 1
		}
	}

	return notify
}

/**
 *	添加通知
 */
func (s *Notify) AddNotify(one *NotifyItem) {
	key := utils.Md5(fmt.Sprintf("%d", one.Type))
	old := new(NotifyItem)
	config.GlobalRedis.GetObject("notify_list", key, old)

	if old.Num > 0 {
		one.Num = old.Num + one.Num
	} else {
		one.Num = old.Num + 1
	}
	one.Md5 = key
	config.GlobalRedis.SetValue("notify_list", key, one)
}

func (s *Notify) NewMsg(msg db.DB_PARAMS) {
	ServiceMsg := new(ServiceMsg)
	if len(msg) == 0 {
		return
	}
	if uid, ok := msg["uid"]; !ok {
		fmt.Println("uid 错误 无法添加")
		return
	} else {
		ServiceMsg.Uid = uid.(int)
		uinfo, _ := config.GlobalDB.FetchOne("users", db.DB_PARAMS{"id": ServiceMsg.Uid}, db.DB_FIELDS{})
		if uinfo != nil {
			ServiceMsg.SendName = uinfo.Get("username").ToString()
		}

	}
	if content, ok := msg["content"]; ok {
		c := fmt.Sprintf("%v", content)
		if err := json.Unmarshal([]byte(c), ServiceMsg); err == nil {
			if err != nil {
				fmt.Println("error", err)
			}
		} else {
			fmt.Println("content 错误 ", err)
		}
	}
	key := utils.Md5(fmt.Sprintf("%d", ServiceMsg.Uid))
	num := 1
	config.GlobalRedis.SetValue("notify_list", key, &NotifyItem{
		Type:  5,
		Num:   num,
		Md5:   key,
		Datad: ServiceMsg,
	})
}
func (s *Notify) NotifyList() []*NotifyItem {
	rs := make([]*NotifyItem, 0)
	list := config.GlobalRedis.GetAll("notify_list")
	for _, item := range list {

		notify := new(NotifyItem)
		if err := json.Unmarshal([]byte(item), notify); err == nil {

			if notify.Num > 0 {
				rs = append(rs, s.FormatNotify(notify))
			}

		}
		config.GlobalRedis.Delete("notify_list")

	}
	return rs

}

func (s *Notify) ClearNotify(tp string) {
	config.GlobalRedis.Del("notify_list", tp)
}
