package taskshell

import (
	"cointrade/config"
	"cointrade/models"
	"cointrade/utils"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const HUOBI_API_HTTP_URL = "https://api-aws.huobi.pro"

func GetHistoryData() {
	for _, v := range models.COIN_LIST {
		if v["isnative"] == "0" {
			continue
		}
		for _, p := range models.PERIOD_LIST {
			utils.ServiceInfo("start history sync:", v["pair"], p)
			go GetHistorySymbol(v["pair"], p)
		}

		//return
	}
}
func GetHistorySymbol(pair string, period string) {
	url := fmt.Sprintf(HUOBI_API_HTTP_URL+"/market/history/kline?period=%s&size=500&symbol=%s", period, pair)
	rq, _ := http.NewRequest("GET", url, nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)
	c.Transport = tr
	resp, err := c.Do(rq)
	//resp, err := http.Get(url)
	if err != nil || resp == nil {
		if err != nil {
			utils.ServiceError("history request failed:", pair, period, err)
		}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		utils.ServiceError("read history body failed:", pair, period, err)
	}
	obj := config.GlobalConfig.GetConfigFromJson(string(body))
	if obj == nil {
		return
	}
	dataValue := obj.GetValue("data")
	if dataValue == nil {
		return
	}
	mp, ok := dataValue.Value.([]interface{})
	if !ok {
		return
	}
	n := len(mp)
	n = n - 1
	for n > 0 {
		item_obj := config.ConfigValue{
			Value: map[string]interface{}{
				"ts":   time.Now().UnixMilli(),
				"tick": mp[n],
			},
		}
		//fmt.Println(item_obj.ToConfig())
		DataChan <- &KlineData{
			Data:   item_obj.ToConfig(),
			Pair:   pair,
			Period: period,
		}
		n--
	}
	//fmt.Println(body)
}
