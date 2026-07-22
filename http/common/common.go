package common

import (
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HTTPARG              = "p"
	HTTP_CODE_NOLOGIN    = 10001 //未登陆
	HTTP_CODE_ERRORPARAM = 10002 //错误的参数
	HTTP_CODE_SUCCESS    = 0     //成功提交
)

type HandleArray []gin.HandlerFunc
type MODULEHANDLELIST []*ModuleHandles

var ModuleGlobal ModuleBase

type TError struct {
	Msg string
}

func (t *TError) Error() string {
	return t.Msg
}

type HttpResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
type MContxt struct { //扩展gin.context 实现
	gin.Context
	Sid  string
	Uid  int
	Data interface{}
}
type Params struct {
	Sid  string      `json:"sid"`
	Uid  interface{} `json:"uid"`
	Data interface{} `json:"data"`
}
type ModuleHandles struct {
	Path    string
	Handles HandleArray
	Method  string
}
type ModuleBase struct {
	EncodeFlag bool
}
type HttpModuleInterface interface {
	ModuleList() MODULEHANDLELIST //路径 HANDLE
}
type HttpModules struct {
	Handle *gin.Engine
}

func (m *ModuleBase) GetP(r *gin.Context) { //获取参数
	str := r.DefaultQuery(HTTPARG, "")
	if str == "" {
		str = r.DefaultPostForm(HTTPARG, "")
	}

	if m.EncodeFlag {

		s, err := utils.DecryptByAes(str)
		if err != nil {
			m.SendResponse(r, HTTP_CODE_ERRORPARAM, nil)
			r.Abort()
		}
		str = string(s)

	}

	var rs Params
	err := json.Unmarshal([]byte(str), &rs)
	if err != nil {
		r.Next()
		return
	}
	//utils.Log("getstr", str)
	r.Set("sid", rs.Sid)
	r.Set("uid", utils.GetInt(utils.GetJsonValue(rs.Uid)))
	r.Set("data", rs.Data)
	r.Next()
}

func (m *ModuleBase) NeedLogin(r *gin.Context) {
	uid := models.MODEL_USER.CheckSessionId(r.GetString("sid"))
	if uid <= 0 || uid != r.GetInt("uid") {
		m.SendResponse(r, HTTP_CODE_NOLOGIN, nil)
		r.Abort()
		return
	}
	r.Next()
}
func (m *ModuleBase) CrossDomain(r *gin.Context) {
	//跨域
	r.Header("Access-Control-Allow-Origin", "*")
	r.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	r.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	r.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	r.Header("Access-Control-Allow-Credentials", "true")
	r.Next()
}
func (m *ModuleBase) ConvertObject(r *gin.Context, obj interface{}) error {
	data, b := r.Get("data")
	if b {
		if data == nil {
			return &TError{Msg: "no this param"}
		}
		dataStr, err := json.Marshal(data)

		if err == nil {
			err := json.Unmarshal(dataStr, obj)
			return err
		} else {
			utils.Log("UnmarshalErr", err)
			return err
		}
	}
	return &TError{Msg: "no this param"}
}
func (m *ModuleBase) SendResponse(r *gin.Context, code int, data interface{}) {

	r.JSON(200, HttpResponse{Code: code, Data: data})
}

func (m *ModuleBase) GetValue(r *gin.Context, key string) string {
	data, b := r.Get("data")
	if b {
		if reflect.TypeOf(data).Kind() == reflect.Map {
			s, ok := data.(map[string]interface{})[key]
			if ok {
				return utils.GetJsonValue(s)
			}
		}
	}
	return ""
}
func (m *ModuleBase) GetInt(r *gin.Context, key string) int {
	data, b := r.Get("data")
	if b {
		if reflect.TypeOf(data).Kind() == reflect.Map {
			s, ok := data.(map[string]interface{})[key]
			if ok {
				return utils.GetInt(utils.GetJsonValue(s))
			}
		}
	}
	return 0
}
func (m *ModuleBase) GetFloat(r *gin.Context, key string) float64 {
	data, b := r.Get("data")
	if b {
		if reflect.TypeOf(data).Kind() == reflect.Map {
			s, ok := data.(map[string]interface{})[key]
			if ok {
				return utils.GetFloat(utils.GetJsonValue(s))
			}
		}
	}
	return 0
}

func (m *ModuleBase) LogWrite(r *gin.Engine) {
	f, _ := os.OpenFile("./www.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	r.Use(gin.LoggerWithWriter(io.MultiWriter(f, os.Stdout)))
	formatter := func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,http协议版本:%s,请求状态码:%d,响应时间:%s,客户端:%s，错误信息:%s\n",
			param.ClientIP,
			param.TimeStamp.Format("2006年01月02日 15:03:04"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}
	r.Use(gin.LoggerWithFormatter(formatter))
}

func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//response := route_response.Response{}
				//response.Data.List = []interface{}{} // 初始化为空切片，而不是空引用
				traceId := c.Writer.Header().Get("X-Request-Trace-Id")
				stackMsg := string(debug.Stack())
				logField := map[string]interface{}{
					"trace_id":    traceId, //  鉴权之后可以得到唯一跟踪ID和用户名
					"user":        c.Writer.Header().Get("X-Request-User"),
					"uri":         c.Request.URL.Path,
					"remote_addr": c.ClientIP(),
					"stack":       stackMsg, // 打印堆栈信息
				}
				c.Abort()
				//response.Code, response.Message = configure.ApiInnerResponseError, fmt.Sprintf("Api内部报错，请联系管理员(id=%s", traceId)
				//log.WithFields(logField).Error(err) // 输出panic 信息
				utils.WriteFile("/home/gin_panic.log", utils.GetJsonValue(logField))
				//dao.ModelClient.RedisClient.HMSet(traceId, redisField) // 上报redis
				//c.JSON(http.StatusUnauthorized, response)
				return
			}
		}()

		c.Next()
	}
}
func (h *HttpModules) LoadModule(module HttpModuleInterface) {
	handleList := module.ModuleList()
	for _, v := range handleList {
		switch strings.ToLower(v.Method) {
		case "get":
			h.Handle.GET(v.Path, v.Handles...)
		case "post":
			h.Handle.POST(v.Path, v.Handles...)
		}
	}
}
func (h *HttpModules) Run(port int) {

	h.Handle.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

func CreateHttp() *HttpModules {
	gin.SetMode(gin.ReleaseMode)
	rs := new(HttpModules)
	rs.Handle = gin.Default()
	ModuleGlobal.LogWrite(rs.Handle)

	rs.Handle.Use(ModuleGlobal.CrossDomain, ModuleGlobal.GetP)
	return rs
}
