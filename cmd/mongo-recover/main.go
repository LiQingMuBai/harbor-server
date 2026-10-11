package main

import (
	"cointrade/config"
	"cointrade/models"
	"cointrade/task/taskshell"
	"cointrade/utils"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type klineRecordData struct {
	ID     float64 `json:"id" bson:"id"`
	Open   float64 `json:"open" bson:"open"`
	Close  float64 `json:"close" bson:"close"`
	Low    float64 `json:"low" bson:"low"`
	High   float64 `json:"high" bson:"high"`
	Amount float64 `json:"amount" bson:"amount"`
	Vol    float64 `json:"vol" bson:"vol"`
	Count  int     `json:"count" bson:"count"`
	Pair   string  `json:"pair" bson:"pair"`
	Period string  `json:"period" bson:"period"`
	Ts     int64   `json:"ts" bson:"ts"`
}

type huobiKlineResp struct {
	Status string             `json:"status"`
	Data   []huobiKlineDetail `json:"data"`
}

type huobiKlineDetail struct {
	ID     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int     `json:"count"`
}

func main() {
	if err := utils.SetupServiceLogger("mongo-recover"); err != nil {
		log.Fatal(err)
	}
	utils.ServiceStartupBanner("harbor-server mongo-recover", "service", "task", "mode", "mongo-recover")
	if err := models.InitData(); err != nil {
		log.Fatal(err)
	}
	mongoURI := config.GlobalMongo.URI
	if mongoURI == "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%d", config.GlobalMongo.Host, config.GlobalMongo.Port)
	}
	utils.ServiceInfo("mongodb connected:", "uri", maskMongoURI(mongoURI), "db", config.GlobalMongo.DBName)
	rebuildKlineHistory()
}

func maskMongoURI(uri string) string {
	at := strings.Index(uri, "@")
	if at < 0 {
		return uri
	}
	schemeEnd := strings.Index(uri, "://")
	if schemeEnd < 0 {
		return uri
	}
	userInfo := uri[schemeEnd+3 : at]
	colon := strings.Index(userInfo, ":")
	if colon < 0 {
		return uri[:schemeEnd+3] + userInfo + "***@" + uri[at+1:]
	}
	return uri[:schemeEnd+3] + userInfo[:colon] + ":***@" + uri[at+1:]
}

func rebuildKlineHistory() {
	for _, v := range models.COIN_LIST {
		if v["isnative"] == "0" {
			continue
		}
		pair := v["pair"]
		for _, period := range models.PERIOD_LIST {
			utils.ServiceInfo("rebuild kline history start:", pair, period)
			table := pair + "_kline_" + period
			if err := config.GlobalMongo.DBHandle.Collection(table).Drop(context.TODO()); err != nil {
				utils.ServiceError("drop kline collection failed:", table, err)
				continue
			}
			records := fetchHuobiKline(pair, period)
			for _, record := range records {
				id := float64(taskshell.CreateId(int(record.ID), period))
				addRecord := klineRecordData{
					ID:     id,
					Open:   record.Open,
					Close:  record.Close,
					Low:    record.Low,
					High:   record.High,
					Amount: record.Amount,
					Vol:    record.Vol,
					Count:  record.Count,
					Pair:   pair,
					Period: period,
					Ts:     record.ID * 1000,
				}
				config.GlobalMongo.InsertData(table, addRecord)
			}
			utils.ServiceInfo("rebuild kline history done:", pair, period, "count", len(records))
		}
	}
}

func fetchHuobiKline(pair, period string) []huobiKlineDetail {
	url := fmt.Sprintf(taskshell.HUOBI_API_HTTP_URL+"/market/history/kline?period=%s&size=500&symbol=%s", period, pair)
	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.ServiceError("build kline request failed:", pair, period, err)
		return nil
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	resp, err := c.Do(rq)
	if err != nil {
		utils.ServiceError("kline request failed:", pair, period, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ServiceError("read kline body failed:", pair, period, err)
		return nil
	}
	var result huobiKlineResp
	if err := json.Unmarshal(body, &result); err != nil {
		utils.ServiceError("unmarshal kline response failed:", pair, period, err)
		return nil
	}
	return result.Data
}
