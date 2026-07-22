package httprun

import "cointrade/models"

//初始化文件

func init() {
	USER_SOCKET_LIST = make(map[int]*WwebsocketWorker) //用户SOCKET映射
	ADMIN_SOCKET_LIST = make(map[int]*WwebsocketWorker)
	DB_MESSAGE_CHAN = make(chan *DBMessage, 512)                           //消息入库通道
	SERVICE_MESSAGE_RECIVE_CHAN = make(chan *models.DBServiceMessage, 512) //客服消息通道
}
