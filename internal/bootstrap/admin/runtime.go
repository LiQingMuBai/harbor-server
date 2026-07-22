package admin

import (
	adminmodels "cointrade/admin_models"
	adminmodules "cointrade/admin_modules"
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/internal/bootstrap/shared"
	"cointrade/models"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
)

type Options struct {
	Port       int
	RPCClients map[int]string
}

func OptionsFromEnv() Options {
	return Options{
		Port:       shared.GetenvInt("ADMIN_PORT", 8080),
		RPCClients: shared.ParseRPCClients(shared.Getenv("RPC_CLIENTS", "9010=127.0.0.1,9020=127.0.0.1")),
	}
}

func Run(options Options) {
	initializeAdmin()
	startAdminBackgroundJobs(options)
	httpServer := createAdminHTTPServer()
	httpServer.Run(options.Port)
}

func initializeAdmin() {
	models.InitData()
	adminmodels.SYSTEM_MODEL.LoadSiteConfig()
}

func startAdminBackgroundJobs(options Options) {
	if len(options.RPCClients) == 0 {
		return
	}
	go rpcClientTask(options.RPCClients)
}

func createAdminHTTPServer() *common.HttpModules {
	common.ModuleGlobal.EncodeFlag = false
	httpServer := common.CreateHttp()
	adminUser := new(adminmodules.AdminUserModule)
	httpServer.LoadModule(adminUser)
	return httpServer
}

func rpcClientTask(clients map[int]string) {
	for {
		popList := config.GlobalRedis.PopQueue(models.QUEUE_RPC_LIST)
		if popList == nil {
			continue
		}
		for _, item := range popList {
			request := new(models.RpcRequest)
			if err := json.Unmarshal([]byte(item), request); err != nil {
				continue
			}
			for port, ip := range clients {
				client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", ip, port))
				if err != nil {
					log.Printf("rpc dial error: err=%v ip=%s port=%d", err, ip, port)
					continue
				}
				var replay int
				client.Call("RpcStruct.RunSystemCmd", request.Cmd, &replay)
				client.Close()
			}
		}
	}
}
