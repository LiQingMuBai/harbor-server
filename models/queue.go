package models

import "cointrade/config"

//队列模块 ntime 是队列产生时间 并非事实发生时间

type QueueModel struct {
	ModelBase
}

func (m *QueueModel) InputTeamQueue(uid, t int, data interface{}) { //推入团队账变队列
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo.UserType != 1 {
		//return
	}
	config.GlobalRedis.PushQueue(QUEUE_TEAM_COIN_LOG, map[string]interface{}{"type": t, "data": data, "uid": uid})
}
func (m *QueueModel) InputUserQueue(uid, t int, data interface{}) { //推入用户账变队列

	config.GlobalRedis.PushQueue(QUEUE_USER_COIN_LOG, map[string]interface{}{"type": t, "data": data, "uid": uid})
}
