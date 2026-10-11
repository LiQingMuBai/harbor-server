package main

import (
	"cointrade/config"
	"cointrade/models"
	"cointrade/utils"
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := utils.SetupServiceLogger("mongo-init"); err != nil {
		log.Fatal(err)
	}
	utils.ServiceStartupBanner("harbor-server mongo-init", "service", "mongo-init")
	if err := models.InitData(); err != nil {
		log.Fatal(err)
	}
	initMongo()
}

func initMongo() {
	klineOpts := new(options.CreateCollectionOptions)
	klineOpts.SetCapped(true)
	klineOpts.SetMaxDocuments(500)
	klineOpts.SetSizeInBytes(1024 * 1024 * 5)

	mbpOpts := new(options.CreateCollectionOptions)
	mbpOpts.SetCapped(true)
	mbpOpts.SetMaxDocuments(30)
	mbpOpts.SetSizeInBytes(1024 * 1024 * 5)

	nativeCoins := make([]map[string]string, 0, len(models.COIN_LIST))
	for _, v := range models.COIN_LIST {
		if v["isnative"] == "1" {
			nativeCoins = append(nativeCoins, v)
		}
	}

	utils.ServiceInfo("mongo init start:",
		"native_coins", len(nativeCoins),
		"periods", models.PERIOD_LIST,
		"kline_max_docs", 500,
		"mbp_max_docs", 30,
	)

	totalCreated := 0
	totalFailed := 0
	for _, v := range nativeCoins {
		pair := v["pair"]
		utils.ServiceInfo("init coin start:", pair)

		for _, period := range models.PERIOD_LIST {
			table := pair + "_kline_" + period
			if ok := dropAndCreate(table, klineOpts, "kline"); ok {
				totalCreated++
			} else {
				totalFailed++
			}
		}

		for _, suffix := range []string{"_mbp_buy", "_mbp_sell", "_tradedetail_buy", "_tradedetail_sell"} {
			table := pair + suffix
			kind := "mbp"
			if strings.HasPrefix(suffix, "_tradedetail") {
				kind = "tradedetail"
			}
			if ok := dropAndCreate(table, mbpOpts, kind); ok {
				totalCreated++
			} else {
				totalFailed++
			}
		}

		utils.ServiceInfo("init coin done:", pair)
	}

	utils.ServiceInfo("mongo init over:",
		"coins", len(nativeCoins),
		"created", totalCreated,
		"failed", totalFailed,
	)
}

func dropAndCreate(table string, opts *options.CreateCollectionOptions, kind string) bool {
	if err := config.GlobalMongo.DBHandle.Collection(table).Drop(context.TODO()); err != nil {
		utils.ServiceError("drop "+kind+" collection failed:", table, err)
	}
	if err := config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), table, opts); err != nil {
		utils.ServiceError("create "+kind+" collection failed:", table, err)
		return false
	}
	utils.ServiceInfo("create "+kind+" collection:", table)
	return true
}
