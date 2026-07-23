package admin

import (
	"cointrade/config"
	"cointrade/http/common"
	adminhandler "cointrade/internal/admin/handler"
	adminservice "cointrade/internal/admin/service"
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

func Run(options Options) error {
	if err := initializeAdmin(); err != nil {
		return err
	}
	startAdminBackgroundJobs(options)
	httpServer := createAdminHTTPServer()
	httpServer.Run(options.Port)
	return nil
}

func initializeAdmin() error {
	if err := models.InitData(); err != nil {
		return err
	}
	adminservice.SYSTEM_MODEL.LoadSiteConfig()
	return nil
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
	adminUser := new(adminhandler.AdminUserModule)
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
