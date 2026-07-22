package adminmodel

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *SystemModel) KickUser(uid int) *AdminResponse {
	var sid string
	config.GlobalRedis.GetObject(models.HASH_USER_SESSION_ID, fmt.Sprintf("%d", uid), &sid)
	if sid != "" {
		config.GlobalRedis.Del(models.HASH_USER_SESSION, sid)
		config.GlobalRedis.Del(models.HASH_USER_SESSION_ID, fmt.Sprintf("%d", uid))
		models.MODEL_USER.Update(uid, db.DB_PARAMS{"online": 0})
		return &AdminResponse{State: SUCCESS, Data: "踢人成功"}
	}
	return &AdminResponse{State: ERROR, Data: "踢人失败"}
}

func (s *SystemModel) ExplodeController(rq P) *AdminResponse {
	t := rq.Ts()
	sn := t.Get("sn").ToString()
	rs := &AdminResponse{State: ERROR}
	if sn == "" {
		rs.Data = "控制项错误!"
		return rs
	}
	order, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
	if order.Get("state").ToInt() == 1 {
		rs.Data = "该笔订单已经完成，无法控制"
		return rs
	}
	if order.Get("trade_type").ToInt() != 2 {
		rs.Data = "非交割不可控!"
		return rs
	}
	controItem := db.DB_PARAMS{
		"startime":          0,
		"endtime":           0,
		"sn":                order.Get("sn").ToString(),
		"pair":              "",
		"type":              "",
		"controller_type":   "explode_trade",
		"value":             t.Get("result").ToInt(),
		"controller_status": 1,
	}
	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controItem, db.DB_PARAMS{"sn": order.Get("sn").ToString()})
	config.GlobalMongo.FindAndReplace("explode_control", bson.M{"result_time": 0, "sn": sn, "result": t.Get("result").ToInt()}, bson.M{"sn": sn})
	return &AdminResponse{State: SUCCESS, Data: "提交控制成功!"}
}

func (s *SystemModel) KlineController(rq P) *AdminResponse {
	t := rq.Ts()
	pair := t.Get("pair").ToString()
	rs := &AdminResponse{State: ERROR}
	coin, _ := config.GlobalDB.FetchOne(models.DB_TABLE_COINS, db.DB_PARAMS{"pair": pair}, db.DB_FIELDS{})
	if coin == nil {
		rs.Data = "币不存在"
		return rs
	}
	if t.Get("value").ToString() == "" {
		rs.Data = "控制值不能为空"
		return rs
	}
	timeData := t.Get("time").ToArray()
	if len(timeData) < 2 {
		rs.Data = "控制时间不对"
		return rs
	}
	starttime := utils.TimeToint64(timeData[0].ToString())
	endtime := utils.TimeToint64(timeData[1].ToString())
	if endtime < int64(utils.GetNow()) {
		rs.Data = "当前控制时间已经到期!"
		return rs
	}
	controlist := config.GlobalMongo.GetList(models.COIN_CONTROLLER, bson.M{"pair": coin.Get("pair").ToString()}, bson.M{}, 100)
	if len(controlist) > 0 {
		for _, item := range controlist {
			checkStart := utils.GetInt(fmt.Sprintf("%v", item["startime"]))
			checkEnd := utils.GetInt(fmt.Sprintf("%v", item["endtime"]))
			if starttime >= int64(checkStart) && starttime <= int64(checkEnd) && t.Get("sn").ToString() != fmt.Sprintf("%v", item["sn"]) {
				rs.Data = "控制区间不能重复，请重新调整时间!"
				return rs
			}
		}
	}
	sn := fmt.Sprintf("%s_%d_%d", coin.Get("pair").ToString(), starttime, endtime)
	controItem := P{
		"startime":        starttime,
		"endtime":         endtime,
		"sn":              sn,
		"pair":            coin.Get("pair").ToString(),
		"controller_type": "coin_trade",
		"type":            t.Get("type").ToInt(),
		"value":           t.Get("value").ToFloat(),
	}
	close := t.Get("open_price").ToFloat()
	if close == 0 {
		nowprice := models.MODEL_SYSTEM.GetLastCoinInfo(coin.Get("pair").ToString())
		if nowprice == nil {
			rs.Data = "币种当前价格获取失败，无法添加控制"
			return rs
		}
		close = nowprice["close"].(float64)
	}
	if t.Get("type").ToInt() == 2 {
		controItem["dist_price"] = fmt.Sprintf("%.4f", close+close*t.Get("value").ToFloat()/100)
	} else {
		controItem["dist_price"] = t.Get("value").ToFloat()
	}
	controItem["now_price"] = close
	controItem["open_price"] = close
	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controItem, db.DB_PARAMS{"sn": sn})
	return &AdminResponse{State: SUCCESS, Data: "已提交控制!"}
}

func (s *SystemModel) GenerateData(p P) float64 {
	t := p.Ts()
	nowPrice := t.Get("now_price").ToFloat()
	openPrice := t.Get("open_price").ToFloat()
	if nowPrice == 0 {
		nowPrice = openPrice
	}
	controllerOne := config.GlobalMongo.GetOne(models.COIN_CONTROLLER, bson.M{"sn": t.Get("sn").ToString()}, bson.M{})
	if controllerOne == nil {
		return 0
	}
	if e, ok := controllerOne["change_price"]; ok {
		nowPrice = utils.GetFloat(fmt.Sprintf("%v", e))
	}
	if nowPrice == 0 {
		return 0
	}
	start := t.Get("startime").ToInt()
	now := utils.GetNow()
	distPrice := t.Get("dist_price").ToFloat()
	endtime := t.Get("endtime").ToInt()
	diffPrice := distPrice - openPrice
	divDiffTime := (endtime - start) / 10
	divDiffPrice := utils.GetFloat(fmt.Sprintf("%.4f", diffPrice/10))
	if now > endtime {
		return 0
	}
	nowBlock := math.Ceil(float64((now - start)) / float64(divDiffTime))
	blockDistPrice := openPrice + nowBlock*divDiffPrice
	blockTimeDiff := nowBlock*float64(divDiffTime) - float64((now - start))
	blockDiffPercent := utils.GetFloat(fmt.Sprintf("%.2f", blockTimeDiff/float64(divDiffTime))) * 100
	if blockDiffPercent > 20 || blockTimeDiff == 0 {
		rand.Seed(time.Now().UnixNano())
		for {
			fornum := 1
			flag := 1
			num := rand.Intn(700)
			randflag := rand.Intn(100)
			if randflag > 50 {
				flag = -1
			}
			floatper := float64(num) / 10000 * blockDistPrice * float64(flag)
			floatper = utils.GetFloat(fmt.Sprintf("%.4f", floatper))
			nowPrice = utils.GetFloat(fmt.Sprintf("%.4f", blockDistPrice+floatper))
			if nowPrice < distPrice || fornum > 50 {
				break
			}
			fornum++
		}
	} else {
		per := (blockDistPrice - nowPrice) / float64(blockTimeDiff)
		nowPrice = utils.GetFloat(fmt.Sprintf("%.4f", nowPrice+per))
	}
	controllerOne["change_price"] = nowPrice
	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controllerOne, bson.M{"sn": t.Get("sn").ToString()})
	if nowPrice == distPrice && nowBlock == 10 {
		startt := now
		for j := 0; j < 15; j++ {
			startt += 1
			config.GlobalMongo.FindAndReplace("kline_control", bson.M{"pair": t.Get("pair").ToString(), "timemap": startt, "price": nowPrice}, bson.M{"pair": t.Get("pair").ToString(), "timemap": now})
		}
		return nowPrice
	}
	config.GlobalMongo.FindAndReplace("kline_control", bson.M{"pair": t.Get("pair").ToString(), "timemap": now + 1, "price": nowPrice}, bson.M{"pair": t.Get("pair").ToString(), "timemap": now})
	return nowPrice
}

func (s *SystemModel) ControllerTradeList() *AdminResponse {
	list := config.GlobalMongo.GetList(models.COIN_CONTROLLER, bson.M{"controller_type": bson.M{"$ne": ""}}, nil, 100)
	seq := make([]primitive.M, 0)
	for _, item := range list {
		if item["controller_type"].(string) == "coin_trade" && item["endtime"].(int64) > int64(utils.GetNow()) {
			item["startime"] = utils.IntTimeToString(int64(utils.GetInt(fmt.Sprintf("%v", item["startime"]))))
			item["endtime"] = utils.IntTimeToString(int64(utils.GetInt(fmt.Sprintf("%v", item["endtime"]))))
			seq = append(seq, item)
		}
		if item["controller_type"].(string) == "explode_trade" {
			seq = append(seq, item)
		}
	}
	return &AdminResponse{State: SUCCESS, Data: seq}
}

func (s *SystemModel) DeleteController(rq P) *AdminResponse {
	ts := rq.Ts()
	if ts.Get("id").ToString() == "" {
		return &AdminResponse{State: ERROR, Data: "参数错误， 无法删除"}
	}
	id, _ := primitive.ObjectIDFromHex(ts.Get("id").ToString())
	one := config.GlobalMongo.GetOne(models.COIN_CONTROLLER, bson.M{"_id": id}, nil)
	if one != nil {
		config.GlobalMongo.DBHandle.Collection("kline_control").DeleteOne(context.TODO(), bson.M{"pair": one["pair"]})
	}
	config.GlobalMongo.DBHandle.Collection(models.COIN_CONTROLLER).DeleteOne(context.TODO(), bson.M{"_id": id})
	return &AdminResponse{State: SUCCESS, Data: "删除成功!"}
}

func (s *SystemModel) UserLevelListCount(rq P) db.DBValues {
	t := rq.Ts()
	uid := t.Get("uid").ToInt()
	if uid == 0 {
		return nil
	}
	where := []string{fmt.Sprintf("uid = %d", uid)}
	if v := rq.Ts().Get("daytime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("daytime BETWEEN %d AND %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("level").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" level = %d", v))
	}
	list, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER_LEVEL_COUNT, db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"SUM(recharge) as recharge_count",
		"SUM(withdraw) AS withdraw_count",
		"SUM(trade) AS trade_count",
		"SUM(mining_count) AS mining_count",
		"SUM(register_num) AS register_num",
		"SUM(pro_num) AS pro_num",
		"SUM(trade_profit) AS trade_profit",
		"SUM(mining_profit) AS mining_profit",
	})
	return list
}

func (s *SystemModel) ControllerState() P {
	return P{
		"1": "全部赢", "2": "全部输", "3": "买涨赢", "4": "买跌赢", "5": "买涨赢买跌输", "6": "买跌赢买涨输",
	}
}
