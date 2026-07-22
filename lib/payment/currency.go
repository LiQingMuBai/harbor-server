package payment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetCurrencyRate() map[string]float64 {
	//通过接口获取到汇率
	apiKey := strings.TrimSpace(os.Getenv("EXCHANGE_RATE_API_KEY"))
	if apiKey == "" {
		return nil
	}
	apiurl := "https://v6.exchangerate-api.com/v6/" + apiKey + "/latest/USD"
	c := new(http.Client)
	rp, err := c.Get(apiurl)
	if err != nil {
		return nil
	}
	content, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		return nil
	}

	var js interface{}
	err = json.Unmarshal(content, &js)
	rs := make(map[string]float64)
	rates, ok := js.(map[string]interface{})
	if !ok {
		return nil
	}
	rate_map, ok := rates["conversion_rates"]
	if ok {
		for k, v := range rate_map.(map[string]interface{}) {
			rs[k] = v.(float64)
		}
	} else {
		return nil
	}

	return rs
}
