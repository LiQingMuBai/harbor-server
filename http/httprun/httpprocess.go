package httprun

import (
	"cointrade/http/common"
	"cointrade/http/modules"
	"cointrade/models"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"
)

var module_user modules.UserModule
var module_trade modules.TradeModule
var module_mining modules.MingingModule
var module_assets modules.AssetModule
var module_message modules.MessageModule
var module_system modules.SystemModule
var module_credit modules.CreditModule

func RpcServer(addr string, port int) { //进程内驻守的RPCserver
	r := new(models.RpcStruct)
	rpc.Register(r)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	http.Serve(l, nil)
}
func CreateWss(r *gin.Context) {
	m := new(WwebsocketWorker)
	m.WssWorker(r)
}
func Execute(server_port int, localip string, rpc_port int) {

	//config.InitGlobal(true)
	models.InitData() //初始话model层数据

	go RpcServer(localip, rpc_port) //开启RPC服务
	//adminmodel.InitAdmin()
	go models.CheckApprove()               //开启授权检测线程
	go models.GetWalletBalance()           //钱包余额检测
	common.ModuleGlobal.EncodeFlag = false //设置加密
	http := common.CreateHttp()
	http.LoadModule(&module_user)    //用户相关
	http.LoadModule(&module_trade)   //交易相关
	http.LoadModule(&module_assets)  //资产相关
	http.LoadModule(&module_mining)  //矿机产品相关
	http.LoadModule(&module_message) //私人消息相关
	http.LoadModule(&module_system)  //私人消息相关
	http.LoadModule(&module_credit)  //充值提现相关
	//http.Handle.GET("/wss", CreateWss) //开启socket链接

	fmt.Println("this is http process  ", server_port)
	http.Run(server_port)
}
