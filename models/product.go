package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"time"
)

//质押挖矿 挖矿没有虚拟模式
type ProductModel struct {
	ModelBase
}

const (
	BUY_STATE_CREDIT    = 400001 //余额不足
	BUY_STATE_MIN       = 400002 //低于最小限额
	BUY_STATE_MAX       = 400003 //高于最大限额
	BUY_STATE_PER_LIMIT = 400003 //超出每人的购买限制
	BUY_V_GETED         = 400004 //已经体验过了
	BUY_STATE_LOCKED    = 40005  //未允许购买此款矿机

	MTYPE_V = 2 //体验矿机
	MTYPE_R = 1 //真实矿机
)

type ProductProfile struct { //机器属性结构
	Algorithm   string `json:"algorithm"`   //算法
	MathPower   string `json:"mathpower"`   //算力
	GPowW       string `json:"gpoww"`       //官方功率
	Factory     string `json:"factory"`     //厂家
	WallW       string `json:"wallw"`       //墙上功率
	Size        string `json:"size"`        //尺寸
	Weight      string `json:"weight"`      //重量
	Temperature string `json:"temperature"` //温度
	Humidity    string `json:"humidity"`    //湿度
	Chan        string `json:"chan"`        //链
}

type ProductInfo struct {
	Id       int             `json:"id"`       //	id
	Name     string          `json:"name"`     //矿机名称
	Type     int             `json:"type"`     //矿机类型 是否体验类型
	Rate     float64         `json:"rate"`     //收益比例
	RateMin  float64         `json:"rate_min"` //收益比例下线
	Profit   float64         `json:"profit"`   //每日产出
	Circle   int             `json:"circle"`   //周期
	Price    float64         `json:"price"`
	Logo     string          `json:"logo"`      //矿机图片
	Desc     string          `json:"desc"`      //描述
	Profile  *ProductProfile `json:"profile"`   //矿机属性
	PerLimit int             `json:"per_limit"` //每人限购
	IsOpen   int             `json:"isopen"`
	Min      float64         `json:"min"`       //最低投入额度 0为无限制
	Max      float64         `json:"max"`       //最大投入额度 0为无限制
	UserMin  float64         `json:"user_min"`  //用户最低余额
	IsPublic int             `json:"is_public"` //是否公开的
}
type BuyRequest struct {
	Pid    int     `json:"pid"`    //产品ID
	Amount float64 `json:"amount"` //购买额度
}

type OrderListRequest struct {
	PageBaseRequest
	State int `json:"state"` //是否已经解锁
}

func (m *ProductModel) GetAcceptPids(uid int) map[int]int { //获取用户的授权
	rs := make(map[int]int)
	list, err := config.GlobalDB.FetchAll(DB_TABLE_MINING_ACCEPT, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if err != nil {
		return rs
	}
	for _, v := range list {
		rs[v["product_id"].ToInt()] = v["state"].ToInt()
	}
	return rs
}
func (m *ProductModel) CheckProcutAccept(uid int, pid int) bool {
	pinfo := m.GetProductInfo(pid)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ACCEPT, db.DB_PARAMS{"uid": uid, "product_id": pid}, db.DB_FIELDS{})

	if one == nil {
		return !(pinfo.IsPublic == 0)
	}
	if one["state"].ToInt() == 2 {
		return false
	}

	return true
}
func (m *ProductModel) MakeSn(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", "M", timestr, uidstr, 10+rand.Intn(89))
}
func (m *ProductModel) GetProductList() []*ProductInfo {
	//获得所有矿机
	rs := make([]*ProductInfo, 0)
	list, err := config.GlobalDB.FetchAll(DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"isopen": 1}, db.DB_FIELDS{})
	if err != nil {
		return rs
	}
	for _, v := range list {
		tmp := new(ProductInfo)
		tmp.Id = v["id"].ToInt()
		tmp.Name = v["name"].ToString()
		tmp.Type = v["type"].ToInt()
		tmp.Rate = v["rate"].ToFloat()
		tmp.Profit = v["profit"].ToFloat()
		tmp.Circle = v["circle"].ToInt()
		tmp.Logo = v["logo"].ToString()
		tmp.Desc = v["desc"].ToString()
		tmp.Profile = new(ProductProfile)
		e := v["profile"].GetObject(tmp.Profile)
		if e != nil {
			tmp.Profile = nil
		}
		tmp.Min = v["min"].ToFloat()
		tmp.Max = v["max"].ToFloat()
		tmp.PerLimit = v["per_limit"].ToInt()
		tmp.IsPublic = v["is_public"].ToInt()
		tmp.RateMin = v["rate_min"].ToFloat()
		rs = append(rs, tmp)
	}
	return rs
}
func (m *ProductModel) Buy(uid int, rq *BuyRequest) *BaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	pinfo := m.GetProductInfo(rq.Pid)

	rs := new(BaseResponse)
	if pinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "no this product"
		return rs
	}
	count := config.GlobalDB.GetCount(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid, "pid": rq.Pid})
	if count >= pinfo.PerLimit && pinfo.PerLimit > 0 {
		rs.State = BUY_STATE_PER_LIMIT
		rs.Msg = "you buy limit"
		return rs
	}
	/*if !m.CheckProcutAccept(uid, rq.Pid) {
		rs.State = BUY_STATE_LOCKED
		rs.Msg = "not accept"
		return rs
	}*/
	if pinfo.Type == MTYPE_R {
		if pinfo.Min > 0 && rq.Amount < pinfo.Min {
			rs.State = BUY_STATE_MIN
			rs.Msg = "min"
			return rs
		}
		if pinfo.IsPublic == 1 {
			if uinfo.Credit < rq.Amount {
				rs.State = BUY_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}

			if pinfo.Max > 0 && rq.Amount > pinfo.Max {
				rs.State = BUY_STATE_MAX
				rs.Msg = "max"
				return rs
			}
		} else {
			ids := m.GetROrder(uid)
			if len(ids) > 0 {
				rs.State = BUY_V_GETED
				rs.Msg = "geted"
				return rs
			}
		}

	} else { //判断已经得到过体验矿机了
		rq.Amount = 0
		one, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"pid": pinfo.Id, "uid": uid}, db.DB_FIELDS{"id"})
		if one != nil {
			rs.State = BUY_V_GETED
			rs.Msg = "geted"
			return rs
		}
	}
	ntime := utils.GetNow()
	insertData := db.DB_PARAMS{}
	insertData["uid"] = uid
	insertData["pid"] = rq.Pid
	insertData["sn"] = m.MakeSn(uid)
	insertData["amount"] = rq.Amount
	insertData["circle"] = pinfo.Circle
	insertData["createtime"] = ntime
	insertData["profittime"] = ntime + CIRCLE_TIME
	insertData["endtime"] = ntime + CIRCLE_TIME*pinfo.Circle
	insertData["rate_min"] = pinfo.RateMin
	insertData["dispatch_amount"] = rq.Amount
	if pinfo.IsPublic == 1 {
		insertData["state"] = 0
	} else {
		insertData["state"] = 2
	}

	if pinfo.Type == MTYPE_V {
		insertData["profit"] = pinfo.Profit
	} else {
		insertData["profit"] = (float64(pinfo.Rate) / float64(100)) * rq.Amount
	}
	insertData["type"] = pinfo.Type
	insertData["unlocktime"] = 0
	_, err := config.GlobalDB.InsertData(DB_TABLE_MINING_ORDER, insertData)
	if err == nil {
		if pinfo.IsPublic == 0 {
			rs.State = STATE_SUCCESS
			rs.Msg = "success"
			return rs
		}
		if MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          -1 * rq.Amount,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: COIN_LOG_USER_BUY_MINING,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     -1 * rq.Amount,
				LockCredit: 0,
				Sn:         utils.GetJsonValue(insertData["sn"]),
				CreateTime: ntime,
			},
			TeamCoinLogType: TEAM_LOG_MINING,
			TeamCoinLogInfo: QueueTeamLog{
				MiningCount: rq.Amount,
				CreateTime:  ntime,
			},
		}) {
			rs.State = STATE_SUCCESS
			rs.Msg = "success"
			return rs
		}
	}
	rs.State = STATE_SYSTEM_ERROR
	rs.Msg = "system error"
	return rs
}
func (m *ProductModel) GetROrder(uid int) []int { //获取用户正在预约中的产品ID
	rs := make([]int, 0)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid, "_": "(state=2 or state=4)"}, db.DB_FIELDS{"pid"})
	for _, v := range list {
		rs = append(rs, v["pid"].ToInt())
	}
	return rs

}

func (m *ProductModel) GetProductInfo(pid int) *ProductInfo {
	//获得单个矿机信息
	for _, v := range MINPRODUCT_LIST {
		if v.Id == pid {
			return v
		}
	}
	return nil
}
func (m *ProductModel) GetOrderList(uid int, rq OrderListRequest) *PageBaseResponse { //取得矿机列表
	condition := db.DB_PARAMS{"uid": uid}
	if rq.State >= 0 {
		condition["state"] = rq.State
	}
	count := config.GlobalDB.GetCount(DB_TABLE_MINING_ORDER, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	offset := (rq.Page - 1) * rq.Limit
	limitstr := fmt.Sprintf("limit %d,%d", offset, rq.Limit)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_MINING_ORDER, condition, db.DB_FIELDS{}, "order by createtime desc", limitstr)
	ls := make([]map[string]interface{}, 0)
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		if item.Get("state").ToInt() == 0 {
			i["expiredtime"] = item.Get("createtime").ToInt() + item.Get("circle").ToInt()*86400
		}
		ls = append(ls, i)
	}
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.Total = count
	rs.PageTotal = pagesize
	rs.List = ls
	return rs
}
func (m *ProductModel) Unlock(uid int, sn string) *BaseResponse {
	//矿机解锁
	ntime := utils.GetNow()
	rs := new(BaseResponse)
	one := m.GetOrderInfo(uid, sn)
	if one == nil {
		rs.State = STATE_FAILD
		rs.Msg = "no this order"
		return rs
	}
	if ntime < utils.GetInt(one["endtime"]) { //还没到解锁时间
		rs.State = STATE_FAILD
		rs.Msg = "no this order"
		return rs
	}
	data := db.DB_PARAMS{"state": 1, "unlocktime": utils.GetNow()}
	config.GlobalDB.UpdateData(DB_TABLE_MINING_ORDER, data, db.DB_PARAMS{"id": one["id"]})
	if MODEL_USER.AddCredit(uid, &CreditValue{
		Credit:          utils.GetFloat(one["amount"]),
		LockCredit:      0,
		VCrdit:          0,
		LockVCredit:     0,
		UserCoinLogType: COIN_LOG_USER_MINING_BACK,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     utils.GetFloat(one["amount"]),
			LockCredit: 0,
			Sn:         one["sn"],
			CreateTime: utils.GetInt(utils.GetJsonValue(data["unlocktime"])),
		},
		TeamCoinLogType: 0,
		TeamCoinLogInfo: nil,
	}) {
		//解锁利润并返回用户本金和利润
		rs.State = STATE_SUCCESS
		rs.Msg = "success"
		return rs
	}
	rs.State = STATE_SYSTEM_ERROR
	rs.Msg = "system error"
	return rs
}

func (m *ProductModel) GetOrderInfo(uid int, sn string) db.DB_ROW_RESULT { //获得单个矿机订单信息
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}
func (m *ProductModel) GetOrderCount(uid int) map[string]interface{} {
	rs := make(map[string]interface{})
	count := config.GlobalDB.GetCount(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid, "state": 0})
	rs["count"] = count
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid, "state": 0}, db.DB_FIELDS{"SUM(profit) as day_profit,sum(amount) as all_amount"})
	allinfo, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"SUM(allprofit) as history_profit"})
	if one != nil {

		rs["day_profit"] = one["day_profit"].ToFloat()
		//rs["history_profit"] = one["history_profit"].ToFloat()
		rs["all_amount"] = one["all_amount"].ToFloat()
	}
	if allinfo != nil {
		rs["history_profit"] = allinfo["history_profit"].ToFloat()
	}
	return rs
}
