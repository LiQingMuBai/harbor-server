package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"strconv"
)

const defaultSMSAPIURL = "http://api.wftqm.com/api/sms/mtsend"

type SmsAPI struct {
	AppId  string
	AppKey string
}

func (m *SmsAPI) CreateSign() (string, int) {
	n := GetNow()
	signstr := fmt.Sprintf("%s%s%d", m.AppId, m.AppKey, n)
	return Md5(signstr), n
}

func (m *SmsAPI) SendSms(phone string, content string) bool {
	apiURL := strings.TrimSpace(os.Getenv("SMS_API_URL"))
	if apiURL == "" {
		apiURL = defaultSMSAPIURL
	}

	sign, timestamp := m.CreateSign()
	headers := map[string][]string{
		"Content-Type": {"application/json;charset=UTF-8"},
		"Sign":         {sign},
		"Timestamp":    {strconv.Itoa(timestamp)},
		"Api-Key":      {m.AppId},
	}
	data := map[string]interface{}{
		"appId":    m.AppId,
		"numbers":  phone,
		"content":  content,
		"senderId": "112233",
	}
	s := HttpJsonPost(apiURL, headers, data)
	//utils.Log(s)
	if s == "" {
		return false
	}
	var obj interface{}
	err := json.Unmarshal([]byte(s), &obj)
	if err != nil {
		return false
	}
	obj_map, ok := obj.(map[string]interface{})
	if !ok {
		return false
	}
	_, ok = obj_map["status"]
	if !ok {
		return false
	}
	if GetJsonValue(obj_map["status"]) != "0" {
		return false
	}
	return true
}

func (m *SmsAPI) SendSmsA(phone string, content string) bool {
	apiURL := strings.TrimSpace(os.Getenv("SMS_API_URL"))
	if apiURL == "" {
		apiURL = defaultSMSAPIURL
	}

	appID := strings.TrimSpace(os.Getenv("SMS_API_ID"))
	if appID == "" {
		appID = strings.TrimSpace(os.Getenv("SMSID"))
	}
	if appID == "" {
		appID = m.AppId
	}

	appKey := strings.TrimSpace(os.Getenv("SMS_API_KEY"))
	if appKey == "" {
		appKey = strings.TrimSpace(os.Getenv("SMSKEY"))
	}
	if appKey == "" {
		appKey = m.AppKey
	}

	if appID == "" || appKey == "" {
		return false
	}

	rs := HttpFormPost(apiURL, map[string]interface{}{"appkey": appID, "secretkey": appKey, "phone": phone, "content": content})
	ServiceInfo("sms api response:", rs)
	var s interface{}
	err := json.Unmarshal([]byte(rs), &s)
	if err != nil {
		return false
	}
	if mp, ok := s.(map[string]interface{}); !ok {
		return false
	} else {
		if n, b := mp["code"]; !b {
			return false
		} else {
			if n == "0" {
				return true
			}
		}
	}
	return false
}

func HttpFormPost(url string, params map[string]interface{}) string {
	/*if params == nil {
		return ""
	}*/
	postdata := BuildHttpQuery(params)
	headers := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	rq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(postdata))
	if err != nil {
		return ""
	}

	defer rq.Body.Close()

	rq.Header = headers
	//rq.Header.Add("Connection", "keep-alive")
	//rq.Header.Add("Content-Type","application/json;charset=UTF-8")
	//utils.Log(ioutil.ReadAll(rq.Body))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)

	c.Transport = tr
	rp, err := c.Do(rq)
	if err != nil {
		Log(err.Error())
		return ""
	}
	defer rp.Body.Close()
	b, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		return ""
	}
	return string(b)
	//rp, err := c.Post(url, "application/json;charset=UTF-8", strings.NewReader(postdata))
}
func BuildHttpQuery(params map[string]interface{}) string {
	arr := make([]string, 0)
	for k, v := range params {
		arr = append(arr, fmt.Sprintf("%s=%s", k, url.QueryEscape(GetJsonValue(v))))
	}
	return strings.Join(arr, "&")
}
