package httprun

import (
	"bytes"
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket
const (
	SERVER_NAME          = "cb1" //服务名
	WS_CMD_LOGIN         = 10001 //登陆消息
	WS_CMD_EXIT          = 10002 //退出消息
	WS_CMD_PAIRDATA      = 10003 //交易对的当前价格
	WS_CMD_KILNEDATA     = 10004 //交易对的K线数据
	WS_CMD_DEPDATA       = 10005 //交易对的深度数据 挂单盘状况
	WS_CMD_TRADEDETAIL   = 10006 //交易对的成交数据
	WS_CMD_PING          = 10007 //ping
	WS_CMD_DATA          = 10008
	WS_CMD_MESSAGE       = 10009 //发送给用户的通知消息
	WS_CMD_SERVICE       = 10010 //客服消息
	WS_CMD_ENTER_SERVICE = 10011 //用户进入客服系统
	WS_CMD_EXIT_SERVICE  = 10012 //用户退出客服系统

	DATA_GET_STATE_PAIR  = 1 //用户获取全部交易对现在行情状态码 0001 & 0011=0
	DATA_GET_STATE_TRADE = 2 //用户获取历史成交记录状态码 0010
	DATA_GET_STATE_DEP   = 4 //用户获取深度数据状态码 0100
	DATA_GET_STATE_KLINE = 8 //用户获取K线图状态码 1000

)

var ADMIN_SOCKET_LIST map[int]*WwebsocketWorker
var USER_SOCKET_LIST map[int]*WwebsocketWorker //用户与链接的分布表
var WS_GLOBAL_LOCK sync.RWMutex
var WS_ADMIN_LOCK sync.RWMutex

type BaseRequestMessage struct {
	CMD  int         `json:"cmd"`  //命令码
	Data interface{} `json:"data"` //请求参数
}

type PingResponse struct {
	Ts int `json:"ts"` //当前服务器的ping时间
}
type DBMessage struct {
	Uid     int
	Message *models.SendMessage
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var DB_MESSAGE_CHAN chan *DBMessage
var SERVICE_MESSAGE_RECIVE_CHAN chan *models.DBServiceMessage //客服消息接收通道

var COIN_LAST_MAP map[string]interface{}                //币种最新行情
var COIN_HISTORY_MAP map[string]interface{}             //币种历史行情 1000条
var MBP_MAP map[string]interface{}                      //深度数据 50条
var TRADE_DETAIL_MAP map[string]interface{}             //成交数据 50条
var KLINE_LAST_MAP map[string]map[string]interface{}    //K线最新数据
var KLINE_HISTORY_MAP map[string]map[string]interface{} //K线历史数据 1000条

type WwebsocketWorker struct {
	WriteLock    sync.RWMutex
	Ws           *websocket.Conn
	WebState     *gin.Context
	LastSendTime int
}

func StartWSSBackgroundJobs() {
	DataUpdateFunc()
	go DBMessageService()
	go MessageService()
	go DBServiceMessageReciveFunc()
	go DataService()
}

func CreateWSSHTTPServer() *common.HttpModules {
	httpServer := common.CreateHttp()
	httpServer.Handle.GET("/wss", CreateWss)
	return httpServer
}

func DataUpdateFunc() {
	//数据单线程更新 加快处理速度
	COIN_LAST_MAP = models.MODEL_SYSTEM.GetLastCoinData()

	//COIN_HISTORY_MAP = models.MODEL_SYSTEM.GetCoinHistoryData()

	MBP_MAP = models.MODEL_SYSTEM.GetCoinMbp()

	TRADE_DETAIL_MAP = models.MODEL_SYSTEM.GetCoinTradeDetail()

	KLINE_HISTORY_MAP = models.MODEL_SYSTEM.GetCoinHitoryKline()

	KLINE_LAST_MAP = models.MODEL_SYSTEM.GetCoinLastKline()
	go func() {
		for {
			COIN_LAST_MAP = models.MODEL_SYSTEM.GetLastCoinData()
			KLINE_LAST_MAP = models.MODEL_SYSTEM.GetCoinLastKline()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {

			MBP_MAP = models.MODEL_SYSTEM.GetCoinMbp()
			TRADE_DETAIL_MAP = models.MODEL_SYSTEM.GetCoinTradeDetail()
			KLINE_HISTORY_MAP = models.MODEL_SYSTEM.GetCoinHitoryKline()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	/*go func() {
		for {
			TRADE_DETAIL_MAP = models.MODEL_SYSTEM.GetCoinTradeDetail()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for {
			fmt.Println("kline start...")
			KLINE_HISTORY_MAP = models.MODEL_SYSTEM.GetCoinHitoryKline()
			fmt.Println("kline end...")
			time.Sleep(500 * time.Millisecond)
		}
	}()*/

}
func snapshotUserSockets() []*WwebsocketWorker {
	WS_GLOBAL_LOCK.RLock()
	list := make([]*WwebsocketWorker, 0, len(USER_SOCKET_LIST))
	for _, worker := range USER_SOCKET_LIST {
		list = append(list, worker)
	}
	WS_GLOBAL_LOCK.RUnlock()
	return list
}

func DataService() { //单协程发送，避免高频派生 goroutine 压垮调度器
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	pingCycle := 0
	for range ticker.C {
		ntime := utils.GetNow()
		workers := snapshotUserSockets()
		shouldCheckPing := pingCycle >= 50
		for _, v := range workers {
			if shouldCheckPing {
				v.CheckPing(v.WebState, v.Ws)
			}
			if ntime-v.LastSendTime >= 1 {
				v.LastSendTime = ntime
				v.SendData(v.WebState, v.Ws)
			}
		}
		if shouldCheckPing {
			pingCycle = 0
			continue
		}
		pingCycle++
	}
}
func DBMessageService() {
	//消息入库
	insertData := db.DB_PARAMS{}
	for {
		msg := <-DB_MESSAGE_CHAN
		insertData["uid"] = msg.Uid
		insertData["msg_type"] = msg.Message.MsgType
		insertData["state"] = 0
		insertData["createtime"] = utils.GetNow()
		insertData["readtime"] = 0
		insertData["content"] = msg.Message.Msg

		config.GlobalDB.InsertData(models.DB_TABLE_MESSAGE, insertData)
	}
}
func GZIPEn(str string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(str)); err != nil {
		utils.ServiceError("gzip write failed:", err)
		return []byte(str)
	}
	if err := gz.Close(); err != nil {
		utils.ServiceError("gzip close failed:", err)
		return []byte(str)
	}
	return b.Bytes()
}

func parseRequestConfig(data interface{}) *config.MConfig {
	if data == nil {
		return nil
	}
	return (&config.ConfigValue{Value: data}).ToConfig()
}
func DBServiceMessageReciveFunc() {
	utils.ServiceInfo("db service msg server start")
	//客服消息处理
	for {
		message := <-SERVICE_MESSAGE_RECIVE_CHAN
		insertData := db.DB_PARAMS{
			"uid":        message.Uid,
			"content":    message.Content,
			"createtime": message.CreateTime,
			"flag":       message.Flag,
			"id":         0,
		}
		utils.ServiceInfo("service msg:", insertData)
		id, err := config.GlobalDB.InsertData(models.DB_TABLE_SERVICE_MESSAGE, insertData)
		if err == nil {
			insertData["id"] = id
		} else {
			utils.ServiceError("service db error:", err)
		}

		//notify.NOTIFY.NewMsg(insertData)
		go func() {
			for _, v := range ADMIN_SOCKET_LIST {
				utils.ServiceInfo("push admin service msg:", insertData)
				v.SendMessage(v.Ws, BaseRequestMessage{CMD: WS_CMD_SERVICE, Data: insertData})
			}
		}()

	}
}
func DBServiceMessageSendFunc() {
	//客服消息发送
	for {
		list := config.GlobalRedis.PopQueue(models.HASH_USER_SERVICE_MESSAGE)
		if len(list) > 0 {
			for _, v := range list {
				var message models.DBServiceMessage
				err := json.Unmarshal([]byte(v), &message)
				if err == nil {
					so, ok := USER_SOCKET_LIST[message.Uid]
					service_state := 0
					message.Content = strings.Replace(message.Content, "\n", " ", -1)
					if ok {
						so.SendMessage(so.Ws, BaseRequestMessage{CMD: WS_CMD_SERVICE, Data: message.Content}) //发送给客户端
						service_state = so.WebState.GetInt("service_state")
					}
					insertData := db.DB_PARAMS{
						"uid":        message.Uid,
						"content":    message.Content,
						"sn_id":      message.SnId,
						"createtime": message.CreateTime,
						"flag":       message.Flag,
						"read_state": service_state,
					}

					config.GlobalDB.InsertData(models.DB_TABLE_SERVICE_MESSAGE, insertData)
				}

			}
		}
	}
}
func MessageService() { //消息服务
	for {
		list := config.GlobalRedis.PopQueue(models.HASH_USER_MESSAGE)
		if len(list) > 0 {
			for _, v := range list {
				mp := config.GlobalConfig.GetConfigFromJson(v)
				if mp.GetValue("uid") != nil {

					uid := mp.GetValue("uid").ToInt()

					so, ok := USER_SOCKET_LIST[uid]
					sendmessage := models.SendMessage{MsgType: mp.GetValue("cmd").ToInt(), Msg: mp.GetValue("content").Value}
					if ok {
						so.SendMessage(so.Ws, BaseRequestMessage{CMD: WS_CMD_MESSAGE, Data: sendmessage}) //发送给客户端
					}
					DB_MESSAGE_CHAN <- &DBMessage{Uid: uid, Message: &sendmessage}
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
func (m *WwebsocketWorker) Close(r *gin.Context) {

	uid := r.GetInt("uid")
	if uid != 0 {
		if r.GetInt("admin") == 1 { //管理员登录的SOCKET
			WS_ADMIN_LOCK.Lock()
			delete(ADMIN_SOCKET_LIST, uid)
			WS_ADMIN_LOCK.Unlock()
			return
		}
		mm, ok := USER_SOCKET_LIST[uid]
		if ok && mm == m {
			WS_GLOBAL_LOCK.Lock()
			delete(USER_SOCKET_LIST, uid)
			WS_GLOBAL_LOCK.Unlock()
			if uid > 0 {
				models.MODEL_USER.Update(uid, db.DB_PARAMS{"online": 0})
			}

		}

	}
	m.Ws.Close()
}
func (m *WwebsocketWorker) WssWorker(r *gin.Context) {
	//WebSockets 处理模块
	ws, err := upGrader.Upgrade(r.Writer, r.Request, nil)
	if err != nil {

		r.Writer.Write([]byte(err.Error()))
		return
	}
	defer m.Close(r)
	m.Ws = ws
	m.WebState = r
	for {
		//开始死循环处理客户端发来的消息
		var msg BaseRequestMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			utils.ServiceError("websocket read json failed:", err)
			//m.Close(r)
			break
		}
		switch msg.CMD {
		case WS_CMD_LOGIN:
			m.Login(&msg, r, ws) //链接SOCKET
		case WS_CMD_PING:
			m.Ping(&msg, r, ws) //心跳
		case WS_CMD_DATA: //客户端要求获取行情数据
			m.ChangeDataGet(&msg, r, ws) //客户端改变数据接收的模式
		case WS_CMD_SERVICE:
			m.ServiceMessage(r, &msg)
		case WS_CMD_ENTER_SERVICE:
			r.Set("service_state", 1)
		case WS_CMD_EXIT_SERVICE:
			r.Set("service_state", 0)
		}
		time.Sleep(50 * time.Millisecond) //还是停50毫秒
	}
}
func (m *WwebsocketWorker) SendMessage(ws *websocket.Conn, data interface{}) error { //发送消息
	bs, err := json.Marshal(data)
	if err != nil {
		utils.ServiceError("send message marshal failed:", err)
		return err
	}
	send_data := GZIPEn(string(bs))
	m.WriteLock.Lock()
	defer m.WriteLock.Unlock()
	return ws.WriteMessage(websocket.BinaryMessage, send_data)
}
func (m *WwebsocketWorker) ServiceMessage(r *gin.Context, rq *BaseRequestMessage) {
	if !m.CheckLoginState(r) { //没有登录的话要立即断开连接
		m.Close(r)
		return
	}
	if rq == nil || rq.Data == nil {
		return
	}
	text := utils.GetJsonValue(rq.Data)
	uid := r.GetInt("uid")
	if uid <= 0 { //游客不给发消息
		return
	}
	message := new(models.DBServiceMessage)
	message.Content = text
	message.CreateTime = utils.GetNow()
	message.Flag = 1
	message.Uid = uid
	SERVICE_MESSAGE_RECIVE_CHAN <- message

}
func (m *WwebsocketWorker) CheckLoginState(r *gin.Context) bool { //检查登录状态
	return r.GetInt("loginstate") == 1
}
func (m *WwebsocketWorker) CheckPing(r *gin.Context, ws *websocket.Conn) {
	o_uid := r.GetInt("uid")
	isadmin := r.GetInt("admin")
	if o_uid > 0 && isadmin == 0 {
		_uid := models.MODEL_USER.CheckSessionId(r.GetString("sid"))
		if _uid <= 0 { //被退出后要主动断开
			m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_EXIT, Data: nil})
			m.Ws.Close()
			return
		}

	}

	ntime := utils.GetNow()
	err := m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_PING, Data: PingResponse{Ts: ntime}})
	if err != nil {
		m.Close(r)
		return
	}
	r.Set("lastping", ntime)

}
func (m *WwebsocketWorker) ChangeDataGet(rq *BaseRequestMessage, r *gin.Context, ws *websocket.Conn) {
	/*if !m.CheckLoginState(r) { //没有登录的话要立即断开连接
		m.Close(r)
		return
	}*/
	cfg := parseRequestConfig(rq.Data)
	if cfg == nil {
		return
	}
	dataCode := cfg.GetValue("datacode")
	only := cfg.GetValue("only")
	period := cfg.GetValue("period")
	history := cfg.GetValue("history")
	if dataCode != nil { //设置行情码
		r.Set("datacode", dataCode.ToInt())
	}
	if only != nil { //设置交易对
		r.Set("only", only.ToString())
	}
	if period != nil {
		r.Set("period", period.ToString())
	}
	if history != nil {
		r.Set("history", history.ToInt())
	}
}
func (m *WwebsocketWorker) Ping(rq *BaseRequestMessage, r *gin.Context, ws *websocket.Conn) {
	//处理ping的消息
	/*if !m.CheckLoginState(r) { //没有登录的话要立即断开连接
		m.Close(r)
		return
	}*/
	if r.GetInt("admin") == 1 {
		return
	}
	if rq.Data != nil {
		pingmap := parseRequestConfig(rq.Data)
		if pingmap == nil {
			return
		}
		if ts := pingmap.GetValue("ts"); ts != nil {
			if ts.ToInt() != r.GetInt("lastping") {
				m.Close(r)
				return
			}
		}
	}
}

func (m *WwebsocketWorker) Login(rq *BaseRequestMessage, r *gin.Context, ws *websocket.Conn) {
	//处理登录SOCKET
	cfg := parseRequestConfig(rq.Data)
	if cfg == nil {
		return
	}
	sid := cfg.GetValue("sid")
	uid := 0
	if uidValue := cfg.GetValue("uid"); uidValue != nil {
		uid = uidValue.ToInt()
	}
	admin_pass := cfg.GetValue("admin_pass")
	r.Set("admin", 0)
	if sid != nil {
		_uid := models.MODEL_USER.CheckSessionId(sid.ToString())
		if _uid == uid && _uid > 0 {
			ks := r.Keys
			if _, ok := ks["uid"]; ok { //解除游客状态
				o_uid := r.GetInt("uid")
				if o_uid < 0 {
					WS_GLOBAL_LOCK.Lock()
					delete(USER_SOCKET_LIST, o_uid)
					WS_GLOBAL_LOCK.Unlock()
				}

			}
			m_s, ok := USER_SOCKET_LIST[_uid]
			if ok {
				m_s.SendMessage(m_s.Ws, BaseRequestMessage{CMD: WS_CMD_EXIT, Data: nil})
				//m_s.Close(m_s.WebState)
			}

			r.Set("loginstate", 1)
			r.Set("sid", sid.ToString())
			r.Set("uid", _uid)
			r.Set("datacode", 1) //设置数据获取码
			r.Set("only", "")    //设置单独获取的交易对 比如K线 深度数据这些只获取一个的就行了
			r.Set("period", "")
			r.Set("history", 0)
			WS_GLOBAL_LOCK.Lock()
			USER_SOCKET_LIST[_uid] = m
			WS_GLOBAL_LOCK.Unlock()                                                                                          //增加映射
			m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_LOGIN, Data: models.BaseResponse{State: models.STATE_SUCCESS}}) //登录成功

			//go m.CheckPing(r, ws)    //开始5秒钟一次的ping 并检测登陆状态
			//go m.SendData(rq, r, ws) //每个用户一个单独的数据播放携程
			models.MODEL_USER.Update(uid, db.DB_PARAMS{"online": 1})
			return
		} else {
			expectedAdminPass := strings.TrimSpace(os.Getenv("WS_ADMIN_PASS"))
			if expectedAdminPass != "" && admin_pass != nil && admin_pass.ToString() == expectedAdminPass {
				admin_id := utils.GetNow()
				r.Set("loginstate", 1)
				r.Set("admin", 1)
				r.Set("uid", admin_id)
				WS_ADMIN_LOCK.Lock()
				ADMIN_SOCKET_LIST[admin_id] = m
				WS_ADMIN_LOCK.Unlock()
				utils.ServiceInfo("admin socket registered")
				return
			}
			_uid = utils.GetInt(fmt.Sprintf("%d%d", -1*utils.GetNow(), rand.Intn(1000)))
			r.Set("loginstate", 1)
			r.Set("sid", "")
			r.Set("uid", _uid)
			r.Set("datacode", 1) //设置数据获取码
			r.Set("only", "")    //设置单独获取的交易对 比如K线 深度数据这些只获取一个的就行了
			r.Set("period", "")
			r.Set("history", 0)
			WS_GLOBAL_LOCK.Lock()
			USER_SOCKET_LIST[_uid] = m
			WS_GLOBAL_LOCK.Unlock()
			m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_LOGIN, Data: models.BaseResponse{State: models.STATE_SUCCESS}}) //登录成功
			return
		}
	}
	m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_LOGIN, Data: models.BaseResponse{State: models.STATE_FAILD}}) //登陆失败
	m.Close(r)

}
func (m *WwebsocketWorker) SendData(r *gin.Context, ws *websocket.Conn) { //一次性输出所有交易对最新的价格信息
	//处理单一的交易价格数据

	dataCode := r.GetInt("datacode")
	only := r.GetString("only")
	period := r.GetString("period")

	history_state := r.GetInt("history")

	if dataCode&DATA_GET_STATE_KLINE > 0 && only != "" {
		//获取K线图数据
		if period == "" { //取1000条实时数据
			//data=config.GlobalMongo.

			if history_state == 0 {

				//fmt.Println("kline data:", COIN_LAST_MAP[only].(primitive.M)["close"])
				m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_KILNEDATA, Data: COIN_LAST_MAP[only]}) //向用户端推送最新行情信息
			} else {

				m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_KILNEDATA, Data: COIN_HISTORY_MAP[only]}) //向用户端推送最新行情信息
				r.Set("history", 0)
			}

		} else {
			if history_state == 0 {
				//data := config.GlobalMongo.GetOne("kline", bson.M{"pair": only, "period": period}, bson.M{"id": -1})

				if d, ok := KLINE_LAST_MAP[only][period]; ok {
					//fmt.Println("kline data:", KLINE_LAST_MAP[only][period].(primitive.M)["close"])
					m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_KILNEDATA, Data: d})
				}

			} else {
				//data := config.GlobalMongo.GetList("kline", bson.M{"pair": only, "period": period}, bson.M{"id": -1}, limit)
				if d, ok := KLINE_HISTORY_MAP[only][period]; ok {
					m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_KILNEDATA, Data: d})
				}
				r.Set("history", 0)
			}

		}
	}
	if dataCode&DATA_GET_STATE_DEP > 0 && only != "" {
		//获取深度数据
		if d, ok := MBP_MAP[only]; ok {
			m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_DEPDATA, Data: d})
		}

	}
	if dataCode&DATA_GET_STATE_TRADE > 0 && only != "" {
		if d, ok := TRADE_DETAIL_MAP[only]; ok {
			m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_TRADEDETAIL, Data: d})
		}

	}
	if dataCode&DATA_GET_STATE_PAIR > 0 { //当用户获取行情时
		m.SendMessage(ws, BaseRequestMessage{CMD: WS_CMD_PAIRDATA, Data: COIN_LAST_MAP}) //向用户端推送最新行情信息

	}
	//m.LastSendTime = utils.GetNow()
}
