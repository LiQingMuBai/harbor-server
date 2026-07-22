package taskshell

import (
	"bytes"
	adminmodel "cointrade/adminmodel"
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"compress/gzip"
	"context"
	"math"
	"math/rand"

	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

// 获取KLINE数据
var perList = "1min, 5min, 15min, 30min, 60min, 4hour, 1day, 1mon, 1week, 1year"
var ch chan int

type KlineData struct {
	Data   *config.MConfig
	Pair   string
	Period string
}

var ControlPriceStruct map[string]map[int]float64 //Kline控制结构
var CoinInfoMap map[string]db.DB_ROW_RESULT
var MonthDataMap map[string]map[int]map[int64]float64 //月K线预定存储数组 币种SYMBOL 日期时间戳 收盘价
var DataChan chan *KlineData
var MAXFLOAT = 3.0
var MINPRICE = 0.225
var MAXPRICE = 0.235

/*
	func ControlPriceMapUpdate() {
		//价格控制更新线程

		for {
			for _, v := range models.COIN_LIST {
				ControlPriceStruct[v["pair"]] = models.MODEL_SYSTEM.GetControlKline(v["pair"])
			}
			//fmt.Println("ControlPriceStruct", ControlPriceStruct)
			time.Sleep(500 * time.Millisecond)
		}
	}
*/
func connect() *websocket.Conn {

	dailer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	so, _, err := dailer.Dial(HUOBI_WS, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	go func() {
		id := utils.Md5(utils.RandName())
		coins := models.COIN_LIST
		perListArr := strings.Split(perList, ",")
		for _, v := range coins {
			pair := v["pair"]
			if v["isnative"] != "1" {
				//pair = v["vpair"]
				if pair == "usdtusdt" || pair == "usdcusdt" {
					continue
				}
				go CreateCoinKline(pair)
				continue
			}

			for _, p := range perListArr {
				_ = so.WriteJSON(RequestStruct{Sub: "market." + pair + ".kline." + strings.TrimSpace(p), Id: id}) //订阅k线数据
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					fmt.Println("send sub error:", err.Error())
				}
			}

		}
	}()
	return so
}
func GetKine() {
	MonthDataMap = make(map[string]map[int]map[int64]float64) //初始化月K线数据存储
	CoinInfoMap = make(map[string]db.DB_ROW_RESULT)
	for _, v := range models.COIN_LIST {
		CoinInfoMap[v["pair"]] = v
	}
	//go ControlPriceMapUpdate() //控制结果线程
	go DataFunc()
	go ControllerKlineQueue() //K线控制结果生成

	so := connect()
	for {
		if so == nil {
			so = connect()
		} else {
			break
		}
	}
	go func(so *websocket.Conn) { //单独抛出线程以免PING超时
		//time.Sleep(1 * time.Second) //停顿一秒
		for {
			_, buf, e := so.ReadMessage()
			if e != nil {
				fmt.Println(e.Error())
				so = connect()
				continue
			}
			greader, err := gzip.NewReader(bytes.NewReader(buf))
			if err != nil {
				continue
			}

			greader.Close()
			buf, e = ioutil.ReadAll(greader)
			if e != nil {
				fmt.Println(e.Error())
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
					so.WriteJSON(map[string]interface{}{"pong": mp.GetValue("ping").ToInt()})
				}
				if ok := mp.GetValue("ch"); ok != nil {
					ch := ok.ToString()
					tmp := strings.Split(ch, ".")
					if len(tmp) > 2 && tmp[2] == "kline" {
						//GetKlineData(mp, strings.TrimSpace(tmp[1]), tmp[3])
						pair := strings.TrimSpace(tmp[1])
						if v, ok := PAIR_MAP[pair]; ok {
							pair = v
						}
						DataChan <- &KlineData{
							Data:   mp,
							Pair:   pair,
							Period: tmp[3],
						}
					}

				}

				//fmt.Println(string(buf))
			}
			time.Sleep(10 * time.Millisecond) //休眠10毫秒
		}
	}(so)

	<-ch
}
func DataFunc(a ...int) {
	for {
		obj := <-DataChan

		GetKlineData(obj.Data, obj.Pair, obj.Period, a...)
	}
}

func GetKlineData(mp *config.MConfig, pair string, t string, a ...int) {

	ts := mp.GetValue("ts").ToInt()
	//ntime := utils.GetNow()
	insertData := mp.GetValue("tick").ToConfig().ConfigMap
	insertData["ts"] = ts
	insertData["period"] = t
	insertData["pair"] = pair
	//coininfo, ok := CoinInfoMap[pair]
	/*if ok && coininfo["is_f"] == "1" {
		f_price := utils.FormatFloatA(utils.GetFloat(CoinInfoMap[pair]["f_price"]), utils.GetInt(CoinInfoMap[pair]["dnum"]))
		high := f_price + float64(float64(1+rand.Intn(9))/float64(100))*f_price
		low := f_price - float64(float64(1+rand.Intn(9))/float64(100))*f_price
		close := (high + low) / 2

		insertData["high"] = utils.FormatFloatA(high, utils.GetInt(CoinInfoMap[pair]["dnum"]))
		if t != "1min" {
			openinfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": t}, nil)
			if openinfo != nil {
				if v, o := openinfo["close"]; o {
					f_price = v.(float64)
				}
			}
		}
		insertData["open"] = f_price
		insertData["low"] = utils.FormatFloatA(low, utils.GetInt(CoinInfoMap[pair]["dnum"]))
		insertData["close"] = utils.FormatFloatA(close, utils.GetInt(CoinInfoMap[pair]["dnum"]))
	}
	controllerPrice := config.GlobalMongo.GetOne("kline_control", bson.M{"pair": pair, "timemap": ntime}, bson.M{})
	if controllerPrice != nil && coininfo["is_f"] == "1" {
		fmt.Println("获取控制行情", controllerPrice["pair"], "控制价格", controllerPrice["price"])
		//if price, o := pricemap[ntime]; o && price > 0 {
		insertData["close"] = utils.FormatFloatA(utils.GetFloat(fmt.Sprintf("%v", controllerPrice["price"])), utils.GetInt(CoinInfoMap[pair]["dnum"]))
		f_price := insertData["close"].(float64)
		high := f_price + float64(float64(10+rand.Intn(20))/float64(100))*f_price
		low := f_price - float64(float64(10+rand.Intn(20))/float64(100))*f_price
		insertData["high"] = utils.FormatFloatA(high, utils.GetInt(CoinInfoMap[pair]["dnum"]))
		insertData["low"] = utils.FormatFloatA(low, utils.GetInt(CoinInfoMap[pair]["dnum"]))

		//}
	}
	if pair == "mkdusdt" && utils.GetFloat(CoinInfoMap[pair]["f_price"]) > 0 {
		f_price := utils.FormatFloatA(utils.GetFloat(CoinInfoMap[pair]["f_price"]), utils.GetInt(CoinInfoMap[pair]["dnum"]))

		openinfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": t}, nil)
		if openinfo != nil {
			if v, o := openinfo["close"]; o {
				f_price = v.(float64)
			}
		}

		insertData["open"] = f_price
	}*、
	/*if pricemap, ok := ControlPriceStruct[pair]; ok && pricemap != nil {

		if price, o := pricemap[ntime]; o && price > 0 {
			insertData["close"] = utils.FormatFloatA(price, utils.GetInt(CoinInfoMap[pair]["dnum"]))
			f_price := insertData["close"].(float64)
			high := f_price + float64(float64(1+rand.Intn(9))/float64(100))*f_price
			low := f_price - float64(float64(1+rand.Intn(9))/float64(100))*f_price
			insertData["high"] = utils.FormatFloatA(high, utils.GetInt(CoinInfoMap[pair]["dnum"]))
			insertData["low"] = utils.FormatFloatA(low, utils.GetInt(CoinInfoMap[pair]["dnum"]))
		}
	} else {
		//fmt.Println(" pair ..............", pair)
	}*/
	id := insertData["id"]
	if len(a) == 0 {
		config.GlobalMongo.FindAndReplace("lastkline", insertData, bson.M{"pair": pair, "period": t})
	}

	/*one := config.GlobalMongo.GetOne("kline", bson.D{{"id", id}, {"pair", pair}, {"period", t}}, bson.M{"_id": -1})
	if one != nil {
		fmt.Println("kline exists11111111111:", one)
		e := config.GlobalMongo.Update("kline", insertData, bson.M{"_id": one["_id"].(primitive.ObjectID)})
		if e != nil {
			fmt.Println(e.Error())
		}
		return
	}

	config.GlobalMongo.InsertData("kline", insertData)*/

	config.GlobalMongo.FindAndReplace(pair+"_kline_"+t, insertData, bson.D{{"id", id}, {"pair", pair}, {"period", t}})
	//utils.Log(insertData)
}
func RecoverKline2() {
	starttime := utils.GetNow() - 23*60*60
	endtime := utils.GetNow()
	//var PERIOD_LIST []string = []string{"1min", "15min", "30min", "60min", "4hour"}
	for starttime <= endtime {
		for _, v := range models.PERIOD_LIST {
			id := float64(CreateId(starttime, v))
			lastklineinfo := config.GlobalMongo.GetOne("mkdusdt_kline_"+v, bson.D{{"pair", "mkdusdt"}, {"period", v}, {"id", id}}, nil)
			if lastklineinfo != nil {
				lastId := float64(CreateId(starttime, v, 1))
				preklineinfo := config.GlobalMongo.GetOne("mkdusdt_kline_"+v, bson.D{{"pair", "mkdusdt"}, {"period", v}, {"id", lastId}}, nil)
				if lastklineinfo["close"].(float64) < 0.30 || lastklineinfo["open"].(float64) < 0.2 || lastklineinfo["vol"].(float64) < 0 {
					lastklineinfo["close"] = 0.30 + float64(rand.Int31n(10))/float64(1000)
					lastklineinfo["high"] = 0.30 + float64(rand.Int31n(10))/float64(1000)
					lastklineinfo["low"] = 0.30 - float64(rand.Int31n(10))/float64(1000)
					if preklineinfo != nil {
						lastklineinfo["open"] = preklineinfo["close"]
					} else {
						lastklineinfo["open"] = 0.30
					}
					lastklineinfo["vol"] = float64(1000+rand.Intn(3000)) * lastklineinfo["close"].(float64)
					fmt.Println(lastklineinfo)
					config.GlobalMongo.FindAndReplace("mkdusdt_kline_"+v, lastklineinfo, bson.D{{"id", id}, {"pair", "mkdusdt"}, {"period", v}})

				}
			}

		}

		//dayid := float64(CreateId(starttime, "1day"))

		starttime = starttime + 60
	}
	for _, v := range models.PERIOD_LIST {
		lastklineinfo := config.GlobalMongo.GetOne("mkdusdt_kline_"+v, bson.D{{"pair", "mkdusdt"}, {"period", v}}, bson.M{"id": -1})
		config.GlobalMongo.FindAndReplace("lastkline", map[string]interface{}{"id": lastklineinfo["id"], "close": lastklineinfo["close"], "open": lastklineinfo["open"], "high": lastklineinfo["high"], "low": lastklineinfo["low"]}, bson.M{"pair": "mkdusdt", "period": v})
	}
	fmt.Println("over")

}
func RecoverCoinKline(pair string, starttime int, endtime int) {
	var PERIOD_LIST []string = []string{"1min", "15min", "30min", "60min", "4hour"}
	start_price := 0.05
	min_price := 0.06
	max_price := 0.3100
	high_rate := 0.5
	low_rate := 0.3
	//diff_price := (max_price - min_price) / 8
	//min_price = max_price + diff_price
	up_rate := 55
	heart := 0.25
	coin, _ := CoinInfoMap[pair]
	//open_flag := false
	base_amount := 300
	//day_n := time.Unix(int64(starttime), 0).Day()
	endtime = utils.GetNow() - 60
	//lastKline := make(map[string]map[string]interface{})
	for starttime <= endtime {
		//nd := time.Unix(int64(starttime), 0).Day()

		insertData := make(map[string]interface{})
		flag := 1
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rn := 1 + r.Intn(100) //双随机比较 主要破掉大数定理
		//rb = 70
		float_rate := float64(10+rand.Intn(int(heart*100))) / float64(10000) //震荡幅度控制在%1以内 最低0.1%
		fmt.Println("up_rate:", up_rate, "rn:", rn)
		if rn > up_rate {
			flag = -1
		}
		insertData["open"] = start_price
		start_price = start_price + start_price*float64(flag)*float_rate //生成当前close价格
		if start_price < min_price || start_price > max_price {
			//fmt.Println("rate:", math.Abs(start_price-open_price)/open_price)
			if start_price > max_price {
				start_price = start_price - start_price*(float_rate) //开始反向
			} else {
				if start_price < min_price {
					start_price = start_price + start_price*(float_rate) //开始反向
				}

			}

		}

		insertData["ts"] = time.Now().UnixMilli()

		insertData["pair"] = pair
		for _, v := range PERIOD_LIST {
			id := float64(CreateId(starttime, v))
			insertData["id"] = id
			insertData["period"] = v
			lastId := float64(CreateId(starttime, v, 1))
			lastklineinfo := config.GlobalMongo.GetOne(pair+"_kline_"+v, bson.D{{"pair", pair}, {"period", v}, {"id", lastId}}, nil)

			if lastklineinfo == nil {
				//lastklineinfo = config.GlobalMongo.GetOne(pair+"_kline_"+v, bson.D{{"pair", pair}, {"period", v}}, bson.M{"id": -1})

				insertData["amount"] = 100 + rand.Intn(3000)
				insertData["count"] = 50 + rand.Intn(3000)
				if lastklineinfo == nil {
					insertData["open"] = start_price + float64(flag)*start_price*(float_rate)
				} else {
					insertData["open"] = utils.FormatFloatA(lastklineinfo["close"].(float64), utils.GetInt(coin["dnum"]))
					//insertData["open"] = lastklineinfo["high"].(float64)

				}

				insertData["close"] = start_price
				high := start_price + float64(float64(1+rand.Intn(int(1*100)))/float64(10000))*start_price
				low := start_price - float64(float64(1+rand.Intn(int(0.5*100)))/float64(10000))*start_price
				insertData["high"] = utils.FormatFloatA(high, utils.GetInt(coin["dnum"]))
				insertData["low"] = utils.FormatFloatA(low, utils.GetInt(coin["dnum"]))
				insertData["vol"] = float64(insertData["count"].(int)) * start_price
				fmt.Println(insertData)
				//config.GlobalMongo.FindAndReplace("lastkline", insertData, bson.M{"pair": pair, "period": v})
				config.GlobalMongo.FindAndReplace(pair+"_kline_"+v, insertData, bson.D{{"id", id}, {"pair", pair}, {"period", v}})
			} else {

				old_id, b := lastklineinfo["id"].(float64)
				if !b {
					old_id = float64(lastklineinfo["id"].(int32))
				}

				if old_id == id {
					//open = lastInfo["open"].(float64)
					if lastklineinfo != nil {
						insertData["open"] = lastklineinfo["open"].(float64)
						insertData["amount"] = lastklineinfo["amount"].(int32) + int32(rand.Intn(base_amount))
						insertData["count"] = lastklineinfo["count"].(int32) + int32(rand.Intn(base_amount))
					} else {
						insertData["amount"] = int32(rand.Intn(base_amount))
						insertData["count"] = int32(rand.Intn(base_amount))
					}
				} else {
					insertData["amount"] = 100 + rand.Intn(base_amount)
					insertData["count"] = 50 + rand.Intn(base_amount)
					insertData["open"] = lastklineinfo["close"].(float64)
				}
				insertData["close"] = utils.FormatFloatA(start_price, utils.GetInt(coin["dnum"]))
				high := start_price + float64(float64(1+rand.Intn(int(high_rate*100)))/float64(10000))*start_price
				low := start_price - float64(float64(1+rand.Intn(int(low_rate*100)))/float64(10000))*start_price
				insertData["high"] = utils.FormatFloatA(high, utils.GetInt(coin["dnum"]))
				insertData["low"] = utils.FormatFloatA(low, utils.GetInt(coin["dnum"]))
				count, b := insertData["count"].(int)
				if !b {
					count = int(insertData["count"].(int32))
				}
				insertData["vol"] = float64(count) * start_price * 10
				//config.GlobalMongo.FindAndReplace("lastkline", insertData, bson.M{"pair": pair, "period": v})
				config.GlobalMongo.FindAndReplace(pair+"_kline_"+v, insertData, bson.D{{"id", id}, {"pair", pair}, {"period", v}})
			}

		}
		time.Sleep(1 * time.Millisecond)
		randtime := 60
		starttime = starttime + int(randtime)
	}

}
func CreateCoinKline(pair string) {
	coin, ok := CoinInfoMap[pair]
	var KlineConfig models.CoinKlineConfig
	fmt.Println(coin)
	if !ok {
		return
	}
	fmt.Println("=============>", coin)

	start_price := utils.GetFloat(coin["f_price"])
	heart := 1.0
	max_price := 0.235
	min_price := 0.210
	high_rate := 5.0
	low_rate := 5.0
	up_rate := 50
	base_amount := 300
	//rb := 1 + rand.Intn(100)
	for {
		//start_price = utils.GetFloat(coin["f_price"])
		ntime := utils.GetNow()
		//LoadMonthData(pair, int64(ntime))
		//daytime := CreateId(ntime, "1day")

		coin = CoinInfoMap[pair]
		err := json.Unmarshal([]byte(coin["kline_config"]), &KlineConfig)
		fmt.Println(KlineConfig)

		if err == nil {
			if KlineConfig.MaxPrice <= 0 || KlineConfig.MinPrice <= 0 {
				continue
			}
			heart = KlineConfig.Heart
			max_price = KlineConfig.MaxPrice
			min_price = KlineConfig.MinPrice
			high_rate = KlineConfig.HighRate
			low_rate = KlineConfig.LowRate
			up_rate = KlineConfig.UpRate
			base_amount = KlineConfig.BaseAmount
		}

		/*if (ntime % (1 * 60)) == 0 { //每5分钟更换一次对比数 形成不规则振幅 破除大数定理导致的长期振幅阵型一致的问题
			//rand.Seed(time.Now().UnixMicro())
			rb = 1 + rand.Intn(100)
		}*/
		keeppriceInfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": "1day"}, nil)
		if start_price <= 0 {
			fmt.Println("system error")
			continue
		}

		if keeppriceInfo != nil {
			last_price := keeppriceInfo["close"].(float64)
			open_price := keeppriceInfo["open"].(float64)
			float_rate := float64(10+rand.Intn(int(heart*100))) / float64(10000) //震荡幅度控制在%1以内 最低0.1%

			flag := 1
			//rand.Seed(time.Now().UnixMicro()) //
			rn := 1 + rand.Intn(100) //双随机比较 主要破掉大数定理
			//rb = 70
			fmt.Println("up_rate:", up_rate)
			if rn > up_rate {
				flag = -1
			}
			start_price = last_price + open_price*float64(flag)*float_rate //生成当前close价格
			if start_price <= 0 {
				fmt.Println("system error")
				continue
			}
			//要控制涨跌幅不能大于开盘价的20%
			//dayId := CreateId(ntime, "1day")
			/*if dayId != 1672790400 {
				MAXFLOAT = 0.05
				//MINPRICE = open_price - open_price*MAXFLOAT
			}*/
			if start_price < min_price || start_price > max_price {
				//fmt.Println("rate:", math.Abs(start_price-open_price)/open_price)
				if start_price > max_price {
					start_price = max_price - open_price*float_rate //开始反向
				} else {
					if start_price < min_price {
						start_price = min_price + start_price*(float_rate) //开始反向
					}
				}

			}
		}
		for _, v := range models.PERIOD_LIST {
			insertData := make(map[string]interface{})
			openinfo := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": v}, nil)
			if openinfo == nil {
				//开始初始化第一条进去
				continue
				insertData["id"] = float64(CreateId(ntime, v))
				insertData["ts"] = time.Now().UnixMilli()
				insertData["period"] = v
				insertData["pair"] = pair
				insertData["amount"] = 100 + rand.Intn(base_amount)
				insertData["count"] = 50 + rand.Intn(base_amount)
				insertData["open"] = start_price
				insertData["close"] = start_price
				high := start_price + float64(float64(1+rand.Intn(int(high_rate*100)))/float64(10000))*start_price
				low := start_price - float64(float64(1+rand.Intn(int(low_rate*100)))/float64(10000))*start_price
				insertData["high"] = utils.FormatFloatA(high, utils.GetInt(coin["dnum"]))
				insertData["low"] = utils.FormatFloatA(low, utils.GetInt(coin["dnum"]))
				insertData["vol"] = float64(insertData["count"].(int)) * start_price
				if insertData["open"].(float64) == 0 || insertData["close"].(float64) == 0 {
					break
				}
				fmt.Println(insertData)
				config.GlobalMongo.FindAndReplace("lastkline", insertData, bson.M{"pair": pair, "period": v})
				config.GlobalMongo.FindAndReplace(pair+"_kline_"+v, insertData, bson.D{{"id", float64(CreateId(ntime, v))}, {"pair", pair}, {"period", v}})
			} else {
				id := CreateId(ntime, v)
				lastInfo := openinfo
				open := start_price
				old_id, b := lastInfo["id"].(float64)
				if !b {
					old_id = float64(lastInfo["id"].(int32))
				}
				lastklineinfo := config.GlobalMongo.GetOne(pair+"_kline_"+v, bson.D{{"id", float64(id)}, {"pair", pair}, {"period", v}}, nil)
				if old_id == float64(id) {
					open = lastInfo["open"].(float64)
					if lastklineinfo != nil {
						insertData["amount"] = lastklineinfo["amount"].(int32) + int32(rand.Intn(base_amount))
						insertData["count"] = lastklineinfo["count"].(int32) + int32(rand.Intn(base_amount))
					} else {
						insertData["amount"] = int32(rand.Intn(base_amount))
						insertData["count"] = int32(rand.Intn(base_amount))
					}
				} else {
					insertData["amount"] = 100 + rand.Intn(base_amount)
					insertData["count"] = 50 + rand.Intn(base_amount)
					open = lastInfo["close"].(float64)
				}
				insertData["ts"] = time.Now().UnixMilli()
				insertData["period"] = v
				insertData["pair"] = pair

				insertData["open"] = open
				insertData["close"] = utils.FormatFloatA(start_price, utils.GetInt(coin["dnum"]))
				if open == 0 || insertData["close"].(float64) == 0 {
					break
				}
				high := start_price + float64(float64(1+rand.Intn(int(high_rate*100)))/float64(10000))*start_price
				low := start_price - float64(float64(1+rand.Intn(int(low_rate*100)))/float64(10000))*start_price
				insertData["high"] = utils.FormatFloatA(high, utils.GetInt(coin["dnum"]))
				insertData["low"] = utils.FormatFloatA(low, utils.GetInt(coin["dnum"]))

				count, b := insertData["count"].(int32)
				if !b {
					count = int32(insertData["count"].(int))
				}
				insertData["vol"] = float64(count) * start_price
				controllerPrice := config.GlobalMongo.GetOne("kline_control", bson.M{"pair": pair, "timemap": ntime}, bson.M{})
				if controllerPrice != nil {
					fmt.Println("获取控制行情", controllerPrice["pair"], "控制价格", controllerPrice["price"])
					//if price, o := pricemap[ntime]; o && price > 0 {
					insertData["close"] = utils.FormatFloatA(utils.GetFloat(fmt.Sprintf("%v", controllerPrice["price"])), utils.GetInt(CoinInfoMap[pair]["dnum"]))
					f_price := insertData["close"].(float64)
					high := f_price + float64(float64(1+rand.Intn(int(high_rate*100)))/float64(10000))*f_price
					low := f_price - float64(float64(1+rand.Intn(int(low_rate*100)))/float64(10000))*f_price
					insertData["high"] = utils.FormatFloatA(high, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					insertData["low"] = utils.FormatFloatA(low, utils.GetInt(CoinInfoMap[pair]["dnum"]))
					insertData["open"] = open
					//}
				}
				insertData["id"] = float64(id)
				//fmt.Println(insertData)
				start_price = insertData["close"].(float64) //内存存储改变，以防MONGO读不到。。然后有问题。。。会造成回到初始值
				if start_price == 0 || open == 0 {
					break
				}
				if math.Abs((start_price-open)/open) > 0.5 { //超出理性范围
					break
				}
				config.GlobalMongo.FindAndReplace("lastkline", insertData, bson.M{"pair": pair, "period": v})
				config.GlobalMongo.FindAndReplace(pair+"_kline_"+v, insertData, bson.D{{"id", float64(id)}, {"pair", pair}, {"period", v}})
			}
		}
		sleeptime := 2 + rand.Intn(15)
		time.Sleep(time.Duration(sleeptime) * time.Second)
	}
}
func LoadMonthData(pair string, ntime int64) {
	n := time.Unix(ntime, 0)
	month := n.Month()
	year := n.Year()
	MothdDataPath := fmt.Sprintf("../kline_month/%s/%d.%d.txt", pair, year, month)
	_, b := MonthDataMap[pair]
	if !b {
		MonthDataMap[pair] = make(map[int]map[int64]float64)
	}
	_, b = MonthDataMap[pair][int(month)]
	if !b {
		MonthDataMap[pair][int(month)] = make(map[int64]float64)
	} else {
		return
	}
	if utils.FileExists(MothdDataPath) {

		content := utils.ReadAllFile(MothdDataPath)
		content_arr := strings.Split(content, "\n")
		for _, v := range content_arr {
			v_arr := strings.Split(strings.TrimSpace(v), ",")
			if len(v_arr) != 2 {
				continue
			}
			daytime := GetTimeIntFromDayString(strings.TrimSpace(v_arr[0]))
			if daytime == 0 {
				continue
			}
			close_price := utils.GetFloat(v_arr[1])
			if close_price <= 0 {
				continue
			}
			MonthDataMap[pair][int(month)][daytime] = close_price
		}
	}

}
func GetTimeIntFromDayString(day string) int64 {
	day_arr := strings.Split(day, "/")
	if len(day_arr) != 3 {
		return 0
	}
	y := utils.GetInt(day_arr[0])
	m := utils.GetInt(day_arr[1])
	d := utils.GetInt(day_arr[2])
	ntime := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local).Unix()
	return ntime
}
func CreateId(ntime int, t string, last ...int) int {
	t = strings.ToLower(t)
	b := false
	if len(last) > 0 {
		b = true
	}
	switch t {
	case "1day":
		if b {
			ntime = ntime - 24*60*60
		}
		return ntime - ntime%(24*60*60)
	case "1min":
		if b {
			ntime = ntime - 60
		}
		return ntime - ntime%60
	case "5min":
		if b {
			ntime = ntime - 300
		}
		return ntime - ntime%300
	case "15min":
		if b {
			ntime = ntime - 900
		}
		return ntime - ntime%(15*60)
	case "30min":
		if b {
			ntime = ntime - 1800
		}
		return ntime - ntime%(30*60)
	case "60min":
		if b {
			ntime = ntime - 3600
		}
		return ntime - ntime%(60*60)
	case "4hour":
		if b {
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
		new_utime := time.Unix(int64(ntime), 0)
		id := int(time.Date(new_utime.Year(), new_utime.Month(), new_utime.Day(), 0, 0, 0, 0, time.Local).Unix())
		return id
	}
	return 0
}
func ControllerKlineQueue() {

	for {
		contro_list := config.GlobalMongo.GetList(models.COIN_CONTROLLER, bson.M{"controller_type": bson.M{`$ne`: ""}}, nil, 100)
		for _, item := range contro_list {
			if b, err := json.Marshal(item); err != nil {
				fmt.Println("解析控制失败......")
				continue
			} else {
				data := make(adminmodel.P, 0)
				if err = json.Unmarshal(b, &data); err != nil {
					fmt.Println("任务解析失败， 跳过!")
					continue
				}

				re := data.Ts()

				if re.Get("endtime").ToInt() < utils.GetNow() { //结束时间大于当前时间退出
					config.GlobalMongo.DBHandle.Collection(models.COIN_CONTROLLER).DeleteOne(context.TODO(), bson.M{"pair": re.Get("pair").ToString()})
					config.GlobalMongo.DBHandle.Collection("kline_control").DeleteMany(context.TODO(), bson.M{"sn": re.Get("sn").ToString()})
					continue
				}
				fmt.Println("获取控制", re.Get("pair").ToString())

				if re.Get("startime").ToInt() > utils.GetNow() { //如果未到时间休息2s
					//fmt.Println("未开始.....")
					//time.Sleep(time.Second * 1)
					continue
				}

				//if re.Get("now_price").ToFloat() == 0 {
				coininfo := models.MODEL_SYSTEM.GetLastCoinInfo(re.Get("pair").ToString())
				data["now_price"] = coininfo["close"].(float64)
				if re.Get("open_price").ToFloat() == 0 {
					item["open_price"] = coininfo["close"].(float64)
					data["open_price"] = coininfo["close"].(float64)
					config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, item, bson.M{"sn": re.Get("sn").ToString()})
				}
				//}
				adminmodel.SYSTEM_MODEL.GenerateData(data)

			}
			time.Sleep(time.Second * 1) //休息2s再生成

		}
	}

}
