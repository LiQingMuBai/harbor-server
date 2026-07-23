package mongoinit

import (
	"cointrade/config"
	"cointrade/models"
	"cointrade/utils"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var periods = []string{"1min", "5min", "15min", "30min", "60min", "4hour", "1day", "1mon", "1week", "1year"}

func Run() {
	models.InitData()

	klineHistoryOpts := new(options.CreateCollectionOptions)
	klineHistoryOpts.SetCapped(true)
	klineHistoryOpts.SetMaxDocuments(500)
	klineHistoryOpts.SetSizeInBytes(1024 * 1024 * 5)

	for _, coin := range models.COIN_LIST {
		for _, period := range periods {
			klineTableName := coin["pair"] + "_kline_" + period
			config.GlobalMongo.DBHandle.Collection(klineTableName).Drop(context.TODO())
			config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), klineTableName, klineHistoryOpts)
			utils.ServiceInfo("mongo init kline collection:", klineTableName)
		}

		mbpTableName := coin["pair"] + "_mbp"
		mbpHistoryOpts := new(options.CreateCollectionOptions)
		mbpHistoryOpts.SetCapped(true)
		mbpHistoryOpts.SetMaxDocuments(30)
		mbpHistoryOpts.SetSizeInBytes(1024 * 1024 * 5)
		config.GlobalMongo.DBHandle.Collection(mbpTableName + "_buy").Drop(context.TODO())
		config.GlobalMongo.DBHandle.Collection(mbpTableName + "_sell").Drop(context.TODO())
		config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), mbpTableName+"_buy", mbpHistoryOpts)
		config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), mbpTableName+"_sell", mbpHistoryOpts)
		utils.ServiceInfo("mongo init mbp collection:", mbpTableName)

		tradeDetailTableName := coin["pair"] + "_tradedetail"
		config.GlobalMongo.DBHandle.Collection(tradeDetailTableName + "_buy").Drop(context.TODO())
		config.GlobalMongo.DBHandle.Collection(tradeDetailTableName + "_sell").Drop(context.TODO())
		config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), tradeDetailTableName+"_buy", mbpHistoryOpts)
		config.GlobalMongo.DBHandle.CreateCollection(context.TODO(), tradeDetailTableName+"_sell", mbpHistoryOpts)
		utils.ServiceInfo("mongo init trade detail collection:", tradeDetailTableName)
	}

	utils.ServiceInfo("mongo init over")
}
