package task

import (
	"cointrade/config"
	"cointrade/internal/bootstrap/shared"
	"cointrade/lib/db"
	"cointrade/lib/payment"
	"cointrade/models"
	"cointrade/task/taskshell"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Options struct {
	Mode string
}

func OptionsFromArgs(args []string) (Options, error) {
	if len(args) == 0 {
		mode := strings.TrimSpace(shared.Getenv("TASK_MODE", ""))
		if mode == "" {
			return Options{}, errors.New("missing task mode")
		}
		return Options{Mode: mode}, nil
	}
	return Options{Mode: args[0]}, nil
}

func Run(options Options) error {
	models.InitData()
	switch options.Mode {
	case "data":
		runDataMode()
	case "approve":
		runApproveMode()
	case "task":
		runTaskMode()
	default:
		return fmt.Errorf("unknown task mode: %s", options.Mode)
	}
	return nil
}

func runDataMode() {
	taskshell.ControlPriceStruct = make(map[string]map[int]float64)
	go traceGoroutines()
	refreshCoinInfo()
	go refreshCoinInfoLoop()
	taskshell.ChanTradeData = make(chan *taskshell.TradeData, 1024)
	taskshell.DataChan = make(chan *taskshell.KlineData, 512)
	go taskshell.GetKine()
	taskshell.GetTradeData()
}

func runApproveMode() {
	taskshell.ApproveRecharge()
}

func runTaskMode() {
	go taskshell.UpdateUserAsset()
	go updateCurrency()
	go taskshell.ClearExplodeTrade()
	go taskshell.ClearDelegateTrade()
	go taskshell.Minging()
	go taskshell.LoanCount()
	go taskshell.ClearKeepCross()
	taskshell.CreditLog()
}

func refreshCoinInfoLoop() {
	for {
		refreshCoinInfo()
		time.Sleep(5 * time.Second)
	}
}

func refreshCoinInfo() {
	taskshell.CoinInfoMap = make(map[string]db.DB_ROW_RESULT)
	for _, value := range models.MODEL_SYSTEM.GetAllCoins() {
		taskshell.CoinInfoMap[value["pair"]] = value
	}
	taskshell.PAIR_MAP = models.MODEL_SYSTEM.GetPairMap()
}

func traceGoroutines() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func updateCurrency() {
	for {
		rates := payment.GetCurrencyRate()
		if rates != nil {
			for key, value := range rates {
				one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_CURRENCY, db.DB_PARAMS{"symbol": strings.ToUpper(key)}, db.DB_FIELDS{})
				if one != nil {
					config.GlobalDB.UpdateData(models.DB_TABLE_CURRENCY, db.DB_PARAMS{"rate": value}, db.DB_PARAMS{"id": one["id"].Value})
				}
			}
		}
		time.Sleep(60 * time.Minute)
	}
}
