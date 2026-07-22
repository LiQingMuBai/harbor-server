package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KlineRecordData struct {
	ID     float64 `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int     `json:"count"`
	Pair   string  `json:"pair"`
	Period string  `json:"period"`
	Ts     int64   `json:"ts"`
}

func main() {
	mongoURI := strings.TrimSpace(os.Getenv("RECOVER_MONGO_URI"))
	if mongoURI == "" {
		mongoURI = strings.TrimSpace(os.Getenv("MONGO_URI"))
	}
	if mongoURI == "" {
		log.Fatal("missing RECOVER_MONGO_URI")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	pairs := []string{"btcusdt", "ethusdt", "xrpusdt", "solusdt", "adausdt", "zecusdt", "usdcusdt", "trxusdt", "dogeusdt", "shibusdt", "linkusdt", "bchusdt", "pepeusdt", "ltcusdt", "suiusdt", "uniusdt", "bnbusdt"}
	periods := []string{"1min", "5min", "15min", "30min", "60min", "4hour", "1day", "1mon", "1week", "1year"}

	databaseName := "trade"
	for _, pair := range pairs {
		for _, period := range periods {
			table := pair + "_kline_" + period

			client.Database(databaseName)
			err := client.Database(databaseName).Collection(table).Drop(context.TODO())
			if err != nil {
				log.Println("删除table失败", table)
				return
			}

			collection := client.Database(databaseName).Collection(table)
			records := GetHistoryData(pair, period)
			for _, record := range records {
				id := float64(CreateId(int(record.ID), period))
				var addRecord KlineRecordData
				addRecord.High = record.High
				addRecord.Vol = record.Vol
				addRecord.Count = record.Count
				addRecord.Period = period
				addRecord.ID = id
				addRecord.Close = record.Close
				addRecord.Low = record.Low
				addRecord.Pair = pair
				addRecord.Open = record.Open
				addRecord.Amount = record.Amount
				addRecord.Ts = record.ID * 1000

				log.Println("id", addRecord.ID)
				log.Println("Amount", addRecord.Amount)
				log.Println("Vol", addRecord.Vol)
				log.Println("Count", addRecord.Count)
				log.Println("Open", addRecord.Open)
				log.Println("Close", addRecord.Close)
				log.Println("High", addRecord.High)
				log.Println("Low", addRecord.Low)
				log.Println("Pair", addRecord.Pair)
				log.Println("Period", addRecord.Period)
				log.Println("Ts", addRecord.Ts)

				insertResult, err := collection.InsertOne(context.TODO(), addRecord)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Inserted a single document: ", insertResult.InsertedID)
			}
			time.Sleep(100 * time.Nanosecond)
		}
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func CreateId(ntime int, t string, last ...int) int {
	t = strings.ToLower(t)
	hasLast := len(last) > 0
	switch t {
	case "1day":
		if hasLast {
			ntime = ntime - 24*60*60
		}
		return ntime - ntime%(24*60*60)
	case "1min":
		if hasLast {
			ntime = ntime - 60
		}
		return ntime - ntime%60
	case "5min":
		if hasLast {
			ntime = ntime - 300
		}
		return ntime - ntime%300
	case "15min":
		if hasLast {
			ntime = ntime - 900
		}
		return ntime - ntime%(15*60)
	case "30min":
		if hasLast {
			ntime = ntime - 1800
		}
		return ntime - ntime%(30*60)
	case "60min":
		if hasLast {
			ntime = ntime - 3600
		}
		return ntime - ntime%(60*60)
	case "4hour":
		if hasLast {
			ntime = ntime - 4*60*60
		}
		return ntime - ntime%(4*60*60)
	case "1mon":
		utime := time.Unix(int64(ntime), 0)
		return int(time.Date(utime.Year(), utime.Month(), 0, 0, 0, 0, 0, time.Local).Unix())
	case "1year":
		utime := time.Unix(int64(ntime), 0)
		return int(time.Date(utime.Year(), utime.Month(), 0, 0, 0, 0, 0, time.Local).Unix())
	case "1week":
		utime := time.Unix(int64(ntime), 0)
		offset := int(utime.Weekday())
		ntime = ntime - offset*24*60
		newUtime := time.Unix(int64(ntime), 0)
		return int(time.Date(newUtime.Year(), newUtime.Month(), newUtime.Day(), 0, 0, 0, 0, time.Local).Unix())
	}
	return 0
}

type KlineData struct {
	Ch     string             `json:"ch"`
	Status string             `json:"status"`
	Ts     string             `json:"ts"`
	Data   []KilineDataDetail `json:"data"`
}

type KilineDataDetail struct {
	ID     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int     `json:"count"`
}

const huobiAPIHTTPURL = "https://api-aws.huobi.pro"

func GetHistoryData(pair, period string) []KilineDataDetail {
	return GetHistorySymbol(pair, period)
}

func GetHistorySymbol(pair string, period string) []KilineDataDetail {
	log.Println(">>>>>>>>>>>>>>>>>GetHistorySymbol<<<<<<<<<<<<<<")
	fmt.Printf("pair: %s, period: %s\n", pair, period)
	url := fmt.Sprintf(huobiAPIHTTPURL+"/market/history/kline?period=%s&size=500&symbol=%s", period, pair)
	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("network dial...")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)
	c.Transport = tr
	resp, err := c.Do(rq)
	if err != nil {
		fmt.Printf(" %+v \n", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("出错了===>", pair, "|", period)
	}

	fmt.Println(string(body))

	var response KlineData
	json.Unmarshal([]byte(string(body)), &response)

	fmt.Println(response.Ch)
	fmt.Println(response.Ts)
	fmt.Println(response.Status)

	log.Println("====================开始====================")

	for _, result := range response.Data {
		log.Println("open", result.Open)
		log.Println("close", result.Close)
		log.Println("id", result.ID)
		log.Println("amount", result.Amount)
		log.Println("high", result.High)
		log.Println("low", result.Low)
		log.Println("count", result.Count)
		log.Println("vol", result.Vol)
	}
	log.Println(len(response.Data))

	log.Println("=====================结束====================")

	return response.Data
}
