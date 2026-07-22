package models

import (
	"cointrade/config"
	"encoding/json"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *SystemModel) GetLastCoinInfo(pair string) primitive.M {
	return config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": "1min"}, nil)
}

func (m *SystemModel) GetLastCoinData() map[string]interface{} {
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = config.GlobalMongo.GetOne("lastkline", bson.M{"pair": v["pair"], "period": "1day"}, nil)
	}
	return rs
}

func (m *SystemModel) GetCoinHistoryData() map[string]interface{} {
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = config.GlobalMongo.GetList("tradedata", bson.M{"pair": v["pair"]}, bson.M{"createtime": -1}, 500)
	}
	return rs
}

func (m *SystemModel) GetCoinMbp() map[string]interface{} {
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = map[string]interface{}{
			"sell": config.GlobalMongo.GetList(v["pair"]+"_mbp_sell", bson.M{}, nil, 0),
			"buy":  config.GlobalMongo.GetList(v["pair"]+"_mbp_buy", bson.M{}, nil, 0),
		}
	}
	return rs
}

func (m *SystemModel) GetCoinTradeDetail() map[string]interface{} {
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = map[string]interface{}{
			"sell": config.GlobalMongo.GetList(v["pair"]+"_tradedetail_sell", bson.M{}, nil, 0),
			"buy":  config.GlobalMongo.GetList(v["pair"]+"_tradedetail_buy", bson.M{}, nil, 0),
		}
	}
	return rs
}

func (m *SystemModel) GetCoinLastKline() map[string]map[string]interface{} {
	rs := make(map[string]map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = make(map[string]interface{})
		for _, period := range PERIOD_LIST {
			rs[v["pair"]][period] = config.GlobalMongo.GetOne("lastkline", bson.M{"pair": v["pair"], "period": period}, nil)
		}
	}
	return rs
}

func (m *SystemModel) GetCoinHitoryKline() map[string]map[string]interface{} {
	rs := make(map[string]map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = make(map[string]interface{})
		for _, period := range PERIOD_LIST {
			rs[v["pair"]][period] = config.GlobalMongo.GetList(v["pair"]+"_kline_"+period, bson.M{}, nil, 0)
		}
	}
	return rs
}

func (m *SystemModel) GetControlExplode(sn string) int32 {
	one := config.GlobalMongo.GetOne("explode_control", bson.M{"sn": sn}, nil)
	if one != nil && one["result_time"].(int32) == 0 {
		return one["result"].(int32)
	}
	return 0
}

func (m *SystemModel) GetControlKline(pair string) map[int]float64 {
	one := config.GlobalMongo.GetOne("kline_control", bson.M{"pair": pair}, nil)
	if one == nil {
		return nil
	}

	rows := make([]map[string]interface{}, 0)
	data, _ := json.Marshal(one["timemap"])
	err := json.Unmarshal(data, &rows)
	result := make(map[int]float64, 0)
	if err == nil && len(rows) > 0 {
		for _, item := range rows {
			now, _ := strconv.Atoi(item["Key"].(string))
			price, _ := item["Value"].(float64)
			result[now] = price
		}
		return result
	}
	return map[int]float64{}
}
