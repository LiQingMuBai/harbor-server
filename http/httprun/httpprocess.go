package httprun

import (
	"cointrade/http/common"
	"cointrade/models"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"
)

func RpcServer(addr string, port int) { //进程内驻守的RPCserver
	r := new(models.RpcStruct)
	rpc.Register(r)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Printf("rpc listen error: %v", err)
		panic(err)
	}
	http.Serve(l, nil)
}
func CreateWss(r *gin.Context) {
	m := new(WwebsocketWorker)
	m.WssWorker(r)
}

func StartAPIBackgroundJobs(localIP string, rpcPort int) {
	go RpcServer(localIP, rpcPort)
	go models.CheckApprove()
	go models.GetWalletBalance()
}

func CreateAPIHTTPServer() *common.HttpModules {
	common.ModuleGlobal.EncodeFlag = false //设置加密
	httpServer := common.CreateHttp()
	registerAPIModules(httpServer)
	return httpServer
}
