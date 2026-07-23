package taskshell

import (
	"bytes"
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

// 火币行情数据接收并保存进MONGOdb
const (
	HUOBI_WS = "wss://api.huobi.pro/ws"
)

type RequestStruct struct {
	Sub string `json:"sub"`
	Id  string `json:"id"`
}
type TradeData struct {
	Type string
	Data *config.MConfig
	Pair string
}

var PAIR_MAP map[string]string
var ChanTradeData chan *TradeData
var c chan int
var LOCK sync.RWMutex
var PairLock map[string]*sync.RWMutex

//var mbp_done chan int
//var detail_done chan int

func Connect() *websocket.Conn {
	dailer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	so, _, err := dailer.Dial(HUOBI_WS, nil)
	if err != nil {
		utils.ServiceError("trade websocket dial failed:", err)
		return nil
	}
	id := utils.Md5(utils.RandName())
	//coins := models.MODEL_SYSTEM.GetAllCoins()
	go func() {
		for _, v := range models.COIN_LIST {
			pair := v["pair"]
			if v["isnative"] != "1" {
				pair = v["vpair"]
			}

			//_ = so.WriteJSON(RequestStruct{Sub: "market." + pair + ".ticker", Id: id}) //订阅行情数据
			//time.Sleep(100 * time.Millisecond)
			LOCK.Lock()
			_ = so.WriteJSON(RequestStruct{Sub: "market." + pair + ".trade.detail", Id: id}) //成交明细
			LOCK.Unlock()
			time.Sleep(100 * time.Millisecond)
			LOCK.Lock()
			_ = so.WriteJSON(RequestStruct{Sub: "market." + pair + ".mbp.150", Id: id}) //20档买卖数据
			LOCK.Unlock()
			time.Sleep(100 * time.Millisecond)
			if err != nil {
				utils.ServiceError("trade subscribe failed:", err)
			}

		}
	}()
	return so
}
func GetTradeData() {

	//u, _ := url.Parse(HUOBI_WS)
	//config.InitGlobal(true)
	//mbp_done = make(chan int)
	//detail_done = make(chan int)
	//err = so.WriteJSON(RequestStruct{Sub: "market.btcusdt.ticker", Id: id})
	PairLock = make(map[string]*sync.RWMutex)
	go DataOp()
	so := reconnectTradeSocket()
	/*coins := models.COIN_LIST
	for _, v := range coins {
		pair := v["pair"]
		if v["isnative"] != "1" {
			pair = v["vpair"]

		}
		PairLock[pair] = new(sync.RWMutex)



	}*/
	ReciveFunc(so)
	//<-c
}

func reconnectTradeSocket() *websocket.Conn {
	for {
		so := Connect()
		if so != nil {
			return so
		}
		utils.ServiceWarn("trade websocket reconnecting")
		time.Sleep(3 * time.Second)
	}
}

func ReciveFunc(so *websocket.Conn) {

	for {
		//开始不停的读取数据
		_, buf, e := so.ReadMessage()
		if e != nil {
			so = reconnectTradeSocket()
			continue
		}
		greader, err := gzip.NewReader(bytes.NewReader(buf))
		if err != nil {
			continue
		}
		buf, e = ioutil.ReadAll(greader)
		greader.Close()
		if e != nil {
			utils.ServiceError("trade websocket read gzip failed:", e)
			continue
		} else {
			//message := string(buf)
			//fmt.Println(string(buf))
			var message_obj interface{}
			json.Unmarshal(buf, &message_obj)
			if message_obj == nil {
				continue
			}
			bvalue := config.ConfigValue{Value: message_obj}
			mp := bvalue.ToConfig()
			if ok := mp.GetValue("ping"); ok != nil { //维持心跳
				LOCK.Lock()
				so.WriteJSON(map[string]interface{}{"pong": mp.GetValue("ping").ToInt()})
				LOCK.Unlock()
			}
			if ok := mp.GetValue("ch"); ok != nil {
				ch := ok.ToString()
				tmp := strings.Split(ch, ".")
				if len(tmp) <= 2 {
					continue
				}
				p := strings.TrimSpace(tmp[1])
				_, ok := PAIR_MAP[p]
				if ok {
					p = PAIR_MAP[p]
				}

				ChanTradeData <- &TradeData{
					Type: tmp[2],
					Data: mp,
					Pair: p,
				}

			}

		}

		time.Sleep(100 * time.Millisecond)
	}
}
func DataOp() {

	for {
		obj := <-ChanTradeData

		switch obj.Type {
		//case "ticker":
		//GetTicks(obj.Data, obj.Pair)

		case "trade":
			GetTradeDetail(obj.Data, obj.Pair)
		case "mbp":
			GetMbp(obj.Data, obj.Pair)
		}
	}
}
func GetTicks(mp *config.MConfig, pair string) { //处理行情数据
	//	fmt.Println("my pair:", pair)
	coinname := strings.TrimSpace(pair)
	insertData := mp.GetValue("tick").ToConfig().ConfigMap
	insertData["pair"] = coinname
	insertData["createtime"] = mp.GetValue("ts").ToInt()
	//fmt.Println("tradedata:", insertData)
	config.GlobalMongo.InsertData("tradedata", insertData)
	config.GlobalMongo.FindAndReplace("lastdata", insertData, bson.M{"pair": pair})
	//fmt.Println("insert result:", rs)

	utils.ServiceInfo("trade ticker updated:", insertData)
}

func GetMbp(mp *config.MConfig, pair string) {
	//处理挂单数据 深度数据
	//fmt.Println(mp.ConfigMap)
	coininfo, _ := CoinInfoMap[pair]
	ticks := mp.GetValue("tick").ToConfig()
	//fmt.Println(ticks.ConfigMap)
	baseClosePrice := 0.0
	hasBaseClosePrice := false
	if coininfo["isnative"] == "0" {
		lastpriceInfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": "1day"}, nil)
		if v, ok := lastpriceInfo["close"]; ok {
			baseClosePrice = v.(float64)
			hasBaseClosePrice = true
		}
	}
	bids, ok := ticks.GetValue("bids").Value.([]interface{})
	ntime := utils.GetNow()
	if ok {
		for _, vv := range bids {
			s, o := vv.([]interface{})
			if !o || len(s) < 2 {
				continue
			}
			if vv.([]interface{})[1].(float64) == 0 {
				continue
			}
			insertData := db.DB_PARAMS{}
			if coininfo["isnative"] == "0" {
				if hasBaseClosePrice {
					insertData["amouut"] = baseClosePrice + float64(float64(1+rand.Intn(9))/float64(100))*baseClosePrice
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
				}
			} else {
				insertData["amouut"] = vv.([]interface{})[0].(float64)
			}
			if pricemap, ok := ControlPriceStruct[pair]; ok && pricemap != nil {

				if price, o := pricemap[ntime]; o && price > 0 {
					//insertData["close"] = utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					f_price := utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					insertData["amouut"] = f_price + float64(float64(1+rand.Intn(9))/float64(100))*f_price
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
				}
			}
			insertData["count"] = vv.([]interface{})[1].(float64)
			//fmt.Println(insertData)
			insertData["ts"] = mp.GetValue("ts").ToInt()
			insertData["pair"] = strings.TrimSpace(pair)
			insertData["type"] = "buy"

			config.GlobalMongo.FindAndReplace(pair+"_mbp_buy", insertData, bson.D{{"pair", pair}, {"amouut", vv.([]interface{})[0].(float64)}, {"type", "buy"}})
			//	utils.Log("bid insert data:", insertData)
		}
	}

	asks, ok := ticks.GetValue("asks").Value.([]interface{})
	//fmt.Println(asks)
	if ok {
		for _, vv := range asks {
			s, o := vv.([]interface{})
			if !o || len(s) < 2 {
				continue
			}
			if vv.([]interface{})[1].(float64) == 0 {
				continue
			}
			insertData := db.DB_PARAMS{}
			if coininfo["isnative"] == "0" {
				if hasBaseClosePrice {
					insertData["amouut"] = baseClosePrice + float64(float64(1+rand.Intn(9))/float64(100))*baseClosePrice
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
				}
			} else {
				insertData["amouut"] = vv.([]interface{})[0].(float64)
			}
			//insertData["amouut"] = vv.([]interface{})[0].(float64)
			if pricemap, ok := ControlPriceStruct[pair]; ok && pricemap != nil {

				if price, o := pricemap[ntime]; o && price > 0 {
					//insertData["close"] = utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					f_price := utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					insertData["amouut"] = f_price + float64(float64(1+rand.Intn(9))/float64(100))*f_price
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
				}
			}
			insertData["count"] = vv.([]interface{})[1].(float64)

			insertData["ts"] = mp.GetValue("ts").ToInt()
			insertData["pair"] = strings.TrimSpace(pair)
			insertData["type"] = "sell"

			config.GlobalMongo.FindAndReplace(pair+"_mbp_sell", insertData, bson.D{{"pair", pair}, {"amouut", vv.([]interface{})[0].(float64)}, {"type", "sell"}})
		}
	}

}
func GetTradeDetail(mp *config.MConfig, pair string) {
	//处理成交数据
	//fmt.Println(mp.ConfigMap)
	ntime := utils.GetNow()
	coininfo := CoinInfoMap[pair]
	ticks := mp.GetValue("tick").ToConfig()
	datas := ticks.GetValue("data").Value.([]interface{})
	baseClosePrice := 0.0
	hasBaseClosePrice := false
	if coininfo["isnative"] == "0" {
		lastpriceInfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": "1day"}, nil)
		if v, ok := lastpriceInfo["close"]; ok {
			baseClosePrice = v.(float64)
			hasBaseClosePrice = true
		}
	}
	//fmt.Println(datas)
	for _, vv := range datas {
		insertData := vv.(map[string]interface{})
		//fmt.Println(insertData)
		if _, ok := insertData["price"]; ok {
			if coininfo["isnative"] == "0" {
				if hasBaseClosePrice {
					insertData["amouut"] = baseClosePrice + float64(float64(1+rand.Intn(9))/float64(100))*baseClosePrice
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
					insertData["price"] = insertData["amouut"]
				}
			}
			if pricemap, ok := ControlPriceStruct[pair]; ok && pricemap != nil {
				if price, o := pricemap[ntime]; o && price > 0 {
					//insertData["close"] = utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					f_price := utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					insertData["amouut"] = f_price + float64(float64(1+rand.Intn(9))/float64(100))*f_price
					insertData["amouut"] = utils.FormatFloatA(insertData["amouut"].(float64), utils.GetInt(coininfo["dnum"]))
					insertData["price"] = insertData["amouut"]
				}
			}
		}
		insertData["pair"] = strings.TrimSpace(pair)
		config.GlobalMongo.InsertData(pair+"_tradedetail_"+insertData["direction"].(string), insertData)
		utils.ServiceInfo("trade detail inserted:", insertData)
	}

}
