package models

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	LOAN_STATE_NOAPPROVE = 7001 //没有授权
	LOAN_STATE_NOOVER    = 7002 //上次的借贷还没有完成
	LOAN_STATE_NOPROOF   = 7003 //提供的资料不完善
)

type LoanModel struct {
	ModelBase
}
type LoanOrderListRequest struct {
	State int `json:"state"`
	PageBaseRequest
}
type LoanOrderRequest struct {
	Amount      float64 `json:"amount"`       //金额
	Circle      int     `json:"circle"`       //周期
	HouseProof  string  `json:"house_proof"`  //房产证明
	IncomeProof string  `json:"income_proof"` //收入证明
	BankProof   string  `json:"bank_proof"`   //银行证明
	IdCard      string  `json:"id_card"`      //身份证件
}

func (m *LoanModel) MakeSn(uid int) string { //创建订单号
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s%s%s%d", "L", timestr, uidstr, 10+rand.Intn(89))
}
func (m *LoanModel) GetRateByCircle(circle int) float64 {
	if f, ok := LOAN_PRODUCT_LIST[circle]; ok {
		return f
	}
	return 0
}
func (m *LoanModel) Loan(uid int, rq *LoanOrderRequest) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil || uinfo.WalletAddress == "" {
		rs.State = LOAN_STATE_NOAPPROVE
		rs.Msg = "no wallet info"
		return rs
	}
	if rq.BankProof == "" || rq.HouseProof == "" || rq.IdCard == "" || rq.IncomeProof == "" {
		rs.State = LOAN_STATE_NOPROOF
		rs.Msg = "no proof"
		return rs
	}
	rate := m.GetRateByCircle(rq.Circle)
	if rate == 0 {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"uid": uid, "state": 0}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = LOAN_STATE_NOOVER
		rs.Msg = "the pre order is not finish"
		return rs
	}
	erc := new(lib.EthLib)
	erc.Type = "usdt"
	erc.CreateClient()
	defer erc.Close()
	if b, e := erc.CheckApprove(uinfo.WalletAddress, config.GlobalConfig.GetValue("approve_wallet").ToString()); e == nil && b {
		//已经授权了 开始产生订单入库
		sn := m.MakeSn(uid)
		insertData := db.DB_PARAMS{"uid": uid}
		insertData["sn"] = sn
		insertData["wallet_address"] = uinfo.WalletAddress
		//insertData["loan_product_id"] = pid
		insertData["amount"] = rq.Amount
		insertData["circle"] = rq.Circle
		insertData["rate"] = rate
		insertData["day_interest"] = float64(rq.Amount) * rate / float64(100)
		//insertData["all_interest"] = utils.GetFloat(pinfo["amount"]) * utils.GetFloat(pinfo["rate"]) / float64(100) * utils.GetFloat(pinfo["circle"])
		insertData["process"] = 0
		insertData["createtime"] = utils.GetNow()
		insertData["finishtime"] = 0
		insertData["state"] = 0
		insertData["house_proof"] = rq.HouseProof
		insertData["income_proof"] = rq.IncomeProof
		insertData["bank_proof"] = rq.BankProof
		insertData["id_card"] = rq.IdCard
		config.GlobalDB.InsertData(DB_TABLE_LOAN_ORDER, insertData)
		rs.State = STATE_SUCCESS
		rs.Msg = "success"
		return rs

	} else {
		rs.State = LOAN_STATE_NOAPPROVE
		rs.Msg = "not approve"
		return rs
	}

}
func (m *LoanModel) GetOrderList(uid int, rq *LoanOrderListRequest) *PageBaseResponse { //获得贷款订单列表
	condition := db.DB_PARAMS{"uid": uid}
	if rq.State > -1 {
		condition["state"] = rq.State
	}
	limit := 15
	count := config.GlobalDB.GetCount(DB_TABLE_LOAN_ORDER, condition)
	pagesize := int(math.Ceil(float64(count) / float64(limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page < 1 {
		rq.Page = 1
	}

	list, _ := config.GlobalDB.FetchRows(DB_TABLE_LOAN_ORDER, condition, db.DB_FIELDS{}, "order by id desc", fmt.Sprintf("limit %d,%d", (rq.Page-1)*limit, limit))
	return &PageBaseResponse{
		BaseResponse: BaseResponse{
			State: STATE_SUCCESS,
			Msg:   "OK",
		},
		Total:     count,
		PageTotal: pagesize,
		Limit:     limit,
		List:      list,
		Page:      rq.Page,
	}
}

func (m *LoanModel) GetLoanInfo(uid int, sn string) db.DB_ROW_RESULT { //取得单笔贷款信息
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}

func (m *LoanModel) GetAllLoanAmount(uid int) map[string]float64 { //获取借贷信息
	all_amount_info, _ := config.GlobalDB.FetchOne(DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"uid": uid, "_": "(state=1 or state=4)"}, db.DB_FIELDS{"SUM(amount) as all_amount"})
	in_amount_info, _ := config.GlobalDB.FetchOne(DB_TABLE_LOAN_ORDER, db.DB_PARAMS{"uid": uid, "state": 1}, db.DB_FIELDS{"SUM(amount) as all_amount"})
	return map[string]float64{"loan_history": all_amount_info["all_amount"].ToFloat(), "loan_prcessing": in_amount_info["all_amount"].ToFloat()}
}
