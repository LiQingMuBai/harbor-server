package trclib

import (
	"cointrade/utils"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

const (
	TRANSFER_FLAG_IN  = 1 //进
	TRANSFER_FLAG_OUT = 2 //出
	TRANSFER_FLAG_ALL = 3 //全部
)

var TRCADDRESSMAP = map[string]string{
	"usdt": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	"onsi": "TDpE9cLSCDj7Y3USwYdF59Eze6A7tbNsSE",
}

type TransResult struct {
	Translist []*TransInfo
	NextUrl   string
}

func GetFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 4)
	if err == nil {
		return f
	}
	return 0
}

type SmartContract struct {
	Owner_address    string `json:"owner_address"`
	Contract_address string `json:"contract_address"`
	Data             string `json:"data"`
}
type TronHelper struct {
	ApiKey             string
	Url                string
	ContractInfo       map[string]interface{}
	ABIMAP             map[string]string
	Client             *client.GrpcClient
	TransferHeader     string
	TransferFromHeader string
}
type JsonValue struct {
	Value interface{}
}
type TransInfo struct {
	//交易信息结构提
	FromAddress     string  //发起方
	ToAddress       string  //接收方
	Amount          float64 //金额
	Hex             string  //确认的区块哈希
	Timestamp       int     //产生时间
	ContractAddress string  //合约地址
}

func (j *JsonValue) Int() int64 {
	n, ok := j.Value.(int64)
	if !ok {
		return 0
	}
	return n
}
func (j *JsonValue) Float() float64 {
	f, ok := j.Value.(float64)
	if !ok {
		return 0
	}
	return f
}
func (j *JsonValue) String() string {
	f, ok := j.Value.(string)
	if !ok {
		return ""
	}
	return f
}
func (j *JsonValue) Map() map[string]*JsonValue {
	f, ok := j.Value.(map[string]interface{})
	if !ok {
		return nil
	}
	rs := make(map[string]*JsonValue)
	for k, v := range f {
		rs[k] = &JsonValue{Value: v}
	}
	return rs
}
func (j *JsonValue) DeepMapValue(keys ...string) *JsonValue {
	if len(keys) == 0 {
		return nil
	}
	if j.Map() == nil {
		return nil
	}
	m := j.Map()
	n := len(keys)
	for i := 0; i < n; i++ {
		t, ok := m[keys[i]]
		if !ok {
			return nil
		}
		if i == n-1 {
			return t
		}
		m = t.Map()
	}
	return nil
}
func (j *JsonValue) Array() []*JsonValue {
	f, ok := j.Value.([]interface{})
	if !ok {
		return nil
	}
	rs := make([]*JsonValue, 0)
	for _, v := range f {
		rs = append(rs, &JsonValue{Value: v})
	}
	return rs
}
func GetJsonValue(s string) *JsonValue {
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		utils.Log(err.Error())
		return nil
	}
	return &JsonValue{Value: v}
}
func CreateHttp(apikey string, url string) *TronHelper {
	rs := new(TronHelper)
	rs.ApiKey = apikey
	rs.Url = url
	rs.TransferFromHeader = rs.GetTransferFromFuncHeader()
	rs.TransferHeader = rs.GetTransferFuncHeader()
	return rs
}

/*func (h *TrcHttpService) CreateHeader() map[string]string {
	//设置请求
	headers := map[string]string{
		"Content-Type":     "application/json",
		"TRON-PRO-API-KEY": h.ApiKey,
	}
	return headers
}*/
func (h *TronHelper) GetTransferFuncHeader() string {
	s := crypto.Keccak256([]byte("transfer(address,uint256)"))
	return fmt.Sprintf("%x", s[0:4])
}
func (h *TronHelper) GetTransferFromFuncHeader() string {
	s := crypto.Keccak256([]byte("transferFrom(address, address,uint256)"))
	return fmt.Sprintf("%x", s[0:4])
}
func (h *TronHelper) InitClient() {
	addr := strings.TrimSpace(os.Getenv("TRON_GRPC_ADDR"))
	if addr == "" {
		addr = "grpc.trongrid.io:50052"
	}
	h.Client = client.NewGrpcClient(addr)
	h.Client.SetAPIKey(h.ApiKey)
	h.Client.Start(grpc.WithInsecure())
}
func (h *TronHelper) GetBlockNumber() int64 {
	b, err := h.Client.GetNowBlock()
	if err != nil {
		return 0
	}

	return b.BlockHeader.RawData.Number
}
func (h *TronHelper) ParasTransFromBlock(ts []*api.TransactionExtention) []*TransInfo {
	translist := make([]*TransInfo, 0)
	for _, vv := range ts {
		cs := vv.Transaction.RawData.GetContract()
		if vv.Transaction.RawData.Timestamp == 0 {
			continue
		}
		for _, con := range cs {

			if con.Type == core.Transaction_Contract_TriggerSmartContract {

				var d core.TriggerSmartContract
				err := proto.Unmarshal(con.GetParameter().Value, &d)
				//c, err := json.Marshal(con.Parameter)
				if err != nil {
					utils.Log("jsonerror:", err.Error())
					continue
				}

				tmp := new(TransInfo)

				tmp.FromAddress = address.HexToAddress(fmt.Sprintf("%x", d.OwnerAddress)).String()
				tmp.ContractAddress = address.HexToAddress(fmt.Sprintf("%x", d.ContractAddress)).String()
				if tmp.ContractAddress != TRCADDRESSMAP["usdt"] {
					continue
				}
				tmp.Hex = fmt.Sprintf("%x", vv.Txid)
				tmp.Timestamp = int(vv.Transaction.RawData.Timestamp)
				//utils.Log(vv.Transaction.RawData.Timestamp, ":", time.Unix(int64(vv.Transaction.RawData.Timestamp/1000), 0).Format("2006-01-02 15:04:05"))
				info := h.DecodeAbi(fmt.Sprintf("%x", d.Data))
				if info == nil {
					continue
				}
				tmp.Amount = float64(info[1].(int64)) / math.Pow10(6)
				tmp.ToAddress = info[0].(string)
				//utils.Log(tmp.Hex)
				if tmp.ToAddress[0] != 'T' {
					utils.Log(tmp.Hex)
				}
				translist = append(translist, tmp)

			}
			if con.Type == core.Transaction_Contract_TransferContract {
				//转TRX时
				var d core.TransferContract
				err := proto.Unmarshal(con.GetParameter().Value, &d)
				//c, err := json.Marshal(con.Parameter)
				if err != nil {
					utils.Log("jsonerror:", err.Error())

					continue
				}
				tmp := new(TransInfo)
				tmp.ContractAddress = ""
				tmp.Amount = float64(d.Amount) / math.Pow10(6)
				tmp.FromAddress = address.HexToAddress(fmt.Sprintf("%x", d.OwnerAddress)).String()
				tmp.Hex = fmt.Sprintf("%x", vv.Txid)
				tmp.ToAddress = address.HexToAddress(fmt.Sprintf("%x", d.ToAddress)).String()
				tmp.Timestamp = int(vv.Transaction.RawData.Timestamp)
				//utils.Log(vv.Transaction.RawData.Timestamp, ":", time.Unix(int64(vv.Transaction.RawData.Timestamp/1000), 0).Format("2006-01-02 15:04:05"))
				translist = append(translist, tmp)
			}
		}

	}
	return translist
}
func (h *TronHelper) GetBlockTransListByNum(n int64) ([]*TransInfo, error) {
	b, err := h.Client.GetBlockByNum(n)
	if err != nil {
		return nil, err
	}
	ts := b.GetTransactions()

	return h.ParasTransFromBlock(ts), nil
}
func (h *TronHelper) GetBlockTransList(n int64) ([]*TransInfo, error) {

	translist := make([]*TransInfo, 0)
	rs, err := h.Client.GetBlockByLatestNum(n)
	if err != nil {
		utils.Log(err.Error())
		return translist, err
	}
	list := rs.GetBlock()
	for _, v := range list {

		ts := v.GetTransactions()
		translist = append(translist, h.ParasTransFromBlock(ts)...)
	}
	return translist, nil
}
func (h *TronHelper) Get(apiurl string) string {
	//发起请求
	utils.Log(h.Url + "v1/" + apiurl)
	rq, err := http.NewRequest("GET", h.Url+"v1/"+apiurl, nil)
	if err != nil {
		utils.Log(err.Error())
		return ""
	}

	//rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("TRON-PRO-API-KEY", h.ApiKey)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)
	c.Transport = tr
	rp, err := c.Do(rq)
	//utils.Log(rp.StatusCode)
	if err != nil {
		utils.Log(err.Error())
		return ""
	}
	if rp.StatusCode != 200 {
		utils.Log("dsadsadsadsadsa")
		return ""
	}
	defer rp.Body.Close()
	s, err := ioutil.ReadAll(rp.Body)
	if err == nil {
		return string(s)
	}
	utils.Log(err.Error())
	return ""
}

func (h *TronHelper) Request(apiurl string, params map[string]interface{}) string {

	//发起请求
	data, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	rq, err := http.NewRequest("POST", h.Url+apiurl, strings.NewReader(string(data)))
	if err != nil {
		utils.Log(err.Error())
		return ""
	}

	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("TRON-PRO-API-KEY", h.ApiKey)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := new(http.Client)
	c.Transport = tr
	rp, err := c.Do(rq)

	if err != nil {
		utils.Log(err.Error())
		return ""
	}
	if rp.StatusCode != 200 {
		utils.Log(rp.StatusCode)
		return ""
	}
	defer rp.Body.Close()
	s, err := ioutil.ReadAll(rp.Body)
	if err == nil {
		return string(s)
	}
	utils.Log(err.Error())
	return ""
}

func (h *TronHelper) CreateAccount() map[string]string { //离线生成钱包地址
	//k, _ := btcec.S256()
	priv, err := btcec.NewPrivateKey()
	if err != nil {
		return nil
	}
	if len(priv.ToECDSA().D.Bytes()) != 32 {
		for {
			priv, err = btcec.NewPrivateKey()
			if err != nil {
				continue
			}
			if len(priv.ToECDSA().D.Bytes()) != 32 {
				continue
			}
			break
		}
	}
	pkey := priv.ToECDSA()
	//utils.Log("private key:", fmt.Sprintf("%x", pkey.D.Bytes()))

	address_str := address.PubkeyToAddress(pkey.PublicKey).String()
	//utils.Log("public key:", address_str)

	return map[string]string{"address": string(address_str), "pkey": fmt.Sprintf("%x", pkey.D.Bytes())}
}
func (h *TronHelper) CreateQueryStr(starttime, endtime, flag int, pripage string) string {
	str := "?only_confirmed=true"
	if starttime > 0 {
		str = str + fmt.Sprintf("&min_timestamp=%d", starttime)
	}
	if endtime > 0 {
		str = str + fmt.Sprintf("&max_timestamp=%d", endtime)
	}
	switch flag {
	case TRANSFER_FLAG_IN:
		str = str + fmt.Sprintf("&only_to=true")
	case TRANSFER_FLAG_OUT:
		str = str + fmt.Sprintf("&only_from=true")
	}
	if len(pripage) > 0 {
		str = str + "&fingerprint=" + pripage
	}
	str = str + "&limit=200"
	return str
}
func (h *TronHelper) GetBanlance(address string) map[string]float64 {
	rs := h.Get("accounts/" + address + "?only_confirmed=true")
	jvalue := GetJsonValue(rs)
	mp := jvalue.Map()
	if mp == nil {
		return nil
	}
	arr := mp["data"].Array()
	if arr == nil {
		return nil
	}
	mp_b := arr[0].Map()
	if mp_b == nil {
		return nil
	}
	b, ok := mp_b["balance"]
	if !ok {
		return nil
	}
	r := make(map[string]float64)
	r["trx"] = b.Float() / float64(math.Pow10(6))
	trc20arr := mp_b["trc20"].Array()
	if trc20arr != nil && len(trc20arr) > 0 {
		//开始获得TRC20的余额
		for _, v := range trc20arr {
			_mp := v.Map()
			if _mp != nil {
				for kk, vv := range _mp {
					for kkk, vvv := range TRCADDRESSMAP {
						if kk == vvv {
							r[kkk] = GetFloat(vv.String()) / float64(math.Pow10(6))
						}
					}
				}
			}
		}
	}
	//r := make(map[string]float64)
	return r
}

func (h *TronHelper) GetTransInfo(faddress string, starttime, endtime int, flag int, pripage string) ([]*TransInfo, string) { //获取交易记录 flag 1 进 2 出 3 全部 返回值 交易记录列表 翻页hash

	rs := make([]*TransInfo, 0)
	body := h.Get("accounts/" + faddress + "/transactions" + h.CreateQueryStr(starttime, endtime, flag, pripage))
	if body == "" {
		return rs, ""
	}
	mp := GetJsonValue(body)
	if mp == nil {
		return rs, ""
	}
	if mp.Map() == nil && mp.Map()["data"].Array() == nil {
		return rs, ""
	}
	j := 0
	for _, v := range mp.Map()["data"].Array() {
		v_mp := v.Map()
		//list := v_mp["raw_data"].Map()["contract"].Array()
		list_obj := v.DeepMapValue("raw_data", "contract")
		if list_obj == nil {
			continue
		}
		list := list_obj.Array()
		//utils.Log(len(list))

		tmp := new(TransInfo)
		tmp.Hex = v_mp["txID"].String()
		//times := v_mp["raw_data"].Map()["raw_data]
		times := v.DeepMapValue("raw_data", "timestamp")
		if times != nil {
			tmp.Timestamp = int(times.Float())
		}
		j++
		for i := 0; i < len(list); i++ {
			vv := list[i]
			//contracts_obj := vv.Map()["parameter"].Map()["value"].Map()["contract_address"]
			contracts_obj := vv.DeepMapValue("parameter", "value", "contract_address")
			if contracts_obj != nil {
				//utils.Log(v.DeepMapValue("parameter", "value").Value)
				//cv, ok := vv.Map()["parameter"].Map()["value"].Map()["call_value"]
				//if ok {
				//utils.Log(cv.Float() / math.Pow10(6))
				//}
				tmp.ContractAddress = address.HexToAddress(vv.DeepMapValue("parameter", "value", "contract_address").String()).String()
				data := vv.DeepMapValue("parameter", "value", "data").String()
				crs := h.DecodeAbi(data)
				tmp.Amount = float64(crs[1].(int64)) / math.Pow10(6)
				tmp.ToAddress = crs[0].(string)
				tmp.FromAddress = address.HexToAddress(vv.Map()["parameter"].Map()["value"].Map()["owner_address"].String()).String()
				continue
			}
			//utils.Log(vv.Map()["parameter"].Map()["value"].Map())

			amount_obj := vv.DeepMapValue("parameter", "value", "amount")
			if amount_obj == nil {
				continue
			}
			amount_n := amount_obj.Float()

			tmp.Amount = float64(amount_n) / math.Pow10(6)
			tmp.FromAddress = address.HexToAddress(vv.Map()["parameter"].Map()["value"].Map()["owner_address"].String()).String()
			tmp.ToAddress = address.HexToAddress(vv.Map()["parameter"].Map()["value"].Map()["to_address"].String()).String()
		}
		rs = append(rs, tmp)

		//utils.Log(tmp.FromAddress)
		//tmp.Amount =
	}
	pre, ok := mp.Map()["meta"].Map()["fingerprint"]

	if ok {
		return rs, pre.String()
	}
	return rs, ""

}

func (h *TronHelper) TransTrx(from string, fromkey string, to string, amount float64) (bool, error) { //TRX转账操作
	trans, err := h.Client.Transfer(from, to, int64(amount*math.Pow10(6)))
	if err != nil {
		return false, err
	}
	pkey, _ := crypto.HexToECDSA(fromkey)
	h.SignTransaction(trans.Transaction, pkey)
	rs, err := h.Client.Client.BroadcastTransaction(context.Background(), trans.Transaction, grpc.EmptyCallOption{})
	if rs.Result {
		return true, nil
	}
	return false, err
}
func (h *TronHelper) TransTrc20(contract_address, from string, fromkey string, to string, amount float64) (bool, error) { //TRC20转账操作
	dataobj := []map[string]interface{}{{"address": to}, {"uint256": strconv.Itoa(int(amount * math.Pow10(6)))}}
	//utils.Log(dataobj)
	jsonstr, _ := json.Marshal(dataobj)
	trans, err := h.Client.TriggerContract(from, contract_address, "transfer(address,uint256)", string(jsonstr), 300000000, 0, "", 0)
	//fmt.Printf(string(trans.Txid))
	pkey, _ := crypto.HexToECDSA(fromkey)
	h.SignTransaction(trans.Transaction, pkey)
	rs, err := h.Client.Client.BroadcastTransaction(context.Background(), trans.Transaction, grpc.EmptyCallOption{})
	//utils.Log(trans.Transaction)
	//utils.Log(fmt.Sprintf("%x", trans.Txid))
	if err != nil || !rs.Result {
		utils.Log(err.Error())
		return false, err
	}
	return true, nil
}

func (h *TronHelper) InitContracts() {
	h.ABIMAP = map[string]string{}
	h.ContractInfo = map[string]interface{}{}
	//初始化所有合约
	for k, a := range TRCADDRESSMAP {
		h.ContractInfo[k] = h.GetContract(a)
		utils.Log((h.ContractInfo[k].(*JsonValue)).Map())
		abi_json, _ := json.Marshal((h.ContractInfo[k].(*JsonValue)).Map()["abi"])
		h.ABIMAP[k] = string(abi_json)
	}
	utils.Log(h.ABIMAP)
}
func (h *TronHelper) GetContract(a string) *JsonValue {

	//hex_address := fmt.Sprintf("%x", []byte(a))

	s := h.Request("wallet/getcontract", map[string]interface{}{"value": a, "visible": true})
	return GetJsonValue(s)
}

func (h *TronHelper) SignTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey) ([]byte, error) {

	transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {

		return nil, err

	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	contractList := transaction.GetRawData().GetContract()
	for range contractList {
		signature, err := crypto.Sign(hash, key)
		if err != nil {
			return nil, err
		}
		transaction.Signature = append(transaction.Signature, signature)
	}

	return hash, nil

}
func (h *TronHelper) DecodeAbi(s string) []interface{} {
	if len(s) < 136 {
		return nil
	}
	bs := []byte(s)
	rs := make([]interface{}, 0)
	method_s := s[0:8]
	if method_s != h.TransferHeader && method_s != h.TransferFromHeader {
		//utils.Log(method_s)
		return nil
	}
	p1_b := bs[8:72]
	p2_b := bs[72:136]

	p1_s := h.ClearPad(p1_b, 1)
	p2_s := h.ClearPad(p2_b, 0)

	a := address.HexToAddress(p1_s).String()
	//utils.Log(a)
	if a[0] != 'T' {
		utils.Log(s)
		utils.Log(p1_s)
		utils.Log(a)
	}

	rs = append(rs, a)
	//utils.Log(p1_s)
	//utils.Log(a)
	//utils.Log(p2_s)
	n, _ := strconv.ParseInt("0x"+p2_s, 0, 64)
	rs = append(rs, n)

	if method_s == h.TransferFromHeader {
		p3_b := bs[136:200]
		p3_s := h.ClearPad(p3_b, 0)
		n, _ := strconv.ParseInt("0x"+p3_s, 0, 64)
		rs[1] = n
		p2_s := h.ClearPad(p2_b, 1)
		rs[0] = address.HexToAddress(p2_s).String()

	}
	//utils.Log("nnn", n)
	//utils.Log(rs[0])
	return rs

}
func (h *TronHelper) ClearPad(s []byte, t int) string {
	tmp := make([]byte, 0)
	f := false
	if t == 1 {
		b := s[22:64]
		b[0] = '4'
		b[1] = '1'
		return string(b)
	}
	for i := 0; i < len(s); i++ { //去除PAD 0
		if f {
			//rs = rs + string(s[i])
			tmp = append(tmp, s[i])
			continue
		}
		if s[i] == '0' {
			continue
		} else {
			tmp = append(tmp, s[i])
			f = true
		}
	}
	return string(tmp)
}
func (h *TronHelper) Close() {
	h.Client.Conn.Close()
}
