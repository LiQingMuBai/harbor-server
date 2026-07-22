package taskshell

import (
	"cointrade/config"
	"cointrade/models"
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
			fmt.Println(v["pair"])
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
	if err != nil {
		fmt.Printf(" %+v \n", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("出错了===>", pair, "|", period)
	}
	obj := config.GlobalConfig.GetConfigFromJson(string(body))
	if obj.GetValue("data") == nil {
		return
	}
	mp := obj.GetValue("data").Value.([]interface{})
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
