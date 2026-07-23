package service

import (
	creditlogdomain "cointrade/internal/domain/creditlog"
	productdomain "cointrade/internal/domain/product"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	productrepo "cointrade/internal/product/repo"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	stateSuccess     = 0
	stateFailed      = 1
	stateSystemError = 9999999
	circleTime       = 24 * 60 * 60
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
}

type ProductCatalog interface {
	GetProductList() []*productdomain.ProductInfo
}

type Service struct {
	repo    productrepo.Repository
	user    UserGateway
	catalog ProductCatalog
}

func New(repo productrepo.Repository, user UserGateway, catalog ProductCatalog) *Service {
	return &Service{
		repo:    repo,
		user:    user,
		catalog: catalog,
	}
}

func (s *Service) GetAcceptPids(uid int) map[int]int {
	rs := make(map[int]int)
	list, err := s.repo.FetchAcceptStates(uid)
	if err != nil {
		return rs
	}
	for _, item := range list {
		rs[item["product_id"].ToInt()] = item["state"].ToInt()
	}
	return rs
}

func (s *Service) CheckProductAccept(uid int, pid int) bool {
	pinfo := s.GetProductInfo(pid)
	if pinfo == nil {
		return false
	}
	one, _ := s.repo.FetchAcceptRecord(uid, pid)
	if one == nil {
		return pinfo.IsPublic != 0
	}
	return one["state"].ToInt() != 2
}

func (s *Service) MakeSn(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", "M", timestr, uidstr, 10+rand.Intn(89))
}

func (s *Service) GetProductList() []*productdomain.ProductInfo {
	rs := make([]*productdomain.ProductInfo, 0)
	list, err := s.repo.FetchOpenProducts()
	if err != nil {
		return rs
	}
	for _, item := range list {
		tmp := new(productdomain.ProductInfo)
		tmp.Id = item["id"].ToInt()
		tmp.Name = item["name"].ToString()
		tmp.Type = item["type"].ToInt()
		tmp.Rate = item["rate"].ToFloat()
		tmp.Profit = item["profit"].ToFloat()
		tmp.Circle = item["circle"].ToInt()
		tmp.Logo = item["logo"].ToString()
		tmp.Desc = item["desc"].ToString()
		tmp.Profile = new(productdomain.ProductProfile)
		if err := item["profile"].GetObject(tmp.Profile); err != nil {
			tmp.Profile = nil
		}
		tmp.Min = item["min"].ToFloat()
		tmp.Max = item["max"].ToFloat()
		tmp.PerLimit = item["per_limit"].ToInt()
		tmp.IsPublic = item["is_public"].ToInt()
		tmp.RateMin = item["rate_min"].ToFloat()
		rs = append(rs, tmp)
	}
	return rs
}

func (s *Service) Buy(uid int, rq *productdomain.BuyRequest) *shareddomain.BaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
	pinfo := s.GetProductInfo(rq.Pid)

	rs := new(shareddomain.BaseResponse)
	if pinfo == nil {
		rs.State = stateSystemError
		rs.Msg = "no this product"
		return rs
	}
	count := s.repo.CountUserProductOrders(uid, rq.Pid)
	if count >= pinfo.PerLimit && pinfo.PerLimit > 0 {
		rs.State = productdomain.BUY_STATE_PER_LIMIT
		rs.Msg = "you buy limit"
		return rs
	}
	if pinfo.Type == productdomain.MTYPE_R {
		if pinfo.Min > 0 && rq.Amount < pinfo.Min {
			rs.State = productdomain.BUY_STATE_MIN
			rs.Msg = "min"
			return rs
		}
		if pinfo.IsPublic == 1 {
			if uinfo.Credit < rq.Amount {
				rs.State = productdomain.BUY_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
			if pinfo.Max > 0 && rq.Amount > pinfo.Max {
				rs.State = productdomain.BUY_STATE_MAX
				rs.Msg = "max"
				return rs
			}
		} else {
			ids := s.GetROrder(uid)
			if len(ids) > 0 {
				rs.State = productdomain.BUY_V_GETED
				rs.Msg = "geted"
				return rs
			}
		}
	} else {
		rq.Amount = 0
		if s.repo.HasUserProductOrder(uid, pinfo.Id) {
			rs.State = productdomain.BUY_V_GETED
			rs.Msg = "geted"
			return rs
		}
	}

	ntime := utils.GetNow()
	insertData := db.DB_PARAMS{
		"uid":             uid,
		"pid":             rq.Pid,
		"sn":              s.MakeSn(uid),
		"amount":          rq.Amount,
		"circle":          pinfo.Circle,
		"createtime":      ntime,
		"profittime":      ntime + circleTime,
		"endtime":         ntime + circleTime*pinfo.Circle,
		"rate_min":        pinfo.RateMin,
		"dispatch_amount": rq.Amount,
		"type":            pinfo.Type,
		"unlocktime":      0,
	}
	if pinfo.IsPublic == 1 {
		insertData["state"] = 0
	} else {
		insertData["state"] = 2
	}
	if pinfo.Type == productdomain.MTYPE_V {
		insertData["profit"] = pinfo.Profit
	} else {
		insertData["profit"] = (float64(pinfo.Rate) / float64(100)) * rq.Amount
	}

	if err := s.repo.InsertOrder(insertData); err == nil {
		if pinfo.IsPublic == 0 {
			rs.State = stateSuccess
			rs.Msg = "success"
			return rs
		}
		if s.user.AddCredit(uid, &userdomain.CreditValue{
			Credit:          -1 * rq.Amount,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: creditlogdomain.COIN_LOG_USER_BUY_MINING,
			UserCoinLogInfo: creditlogdomain.QueueCreditLog{
				Credit:     -1 * rq.Amount,
				LockCredit: 0,
				Sn:         utils.GetJsonValue(insertData["sn"]),
				CreateTime: ntime,
			},
			TeamCoinLogType: creditlogdomain.TEAM_LOG_MINING,
			TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
				MiningCount: rq.Amount,
				CreateTime:  ntime,
			},
		}) {
			rs.State = stateSuccess
			rs.Msg = "success"
			return rs
		}
	}

	rs.State = stateSystemError
	rs.Msg = "system error"
	return rs
}

func (s *Service) GetROrder(uid int) []int {
	rs := make([]int, 0)
	list, _ := s.repo.FetchReservedOrderProductIDs(uid)
	for _, item := range list {
		rs = append(rs, item["pid"].ToInt())
	}
	return rs
}

func (s *Service) GetProductInfo(pid int) *productdomain.ProductInfo {
	for _, item := range s.catalog.GetProductList() {
		if item.Id == pid {
			return item
		}
	}
	return nil
}

func (s *Service) GetOrderList(uid int, rq productdomain.OrderListRequest) *shareddomain.PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	if rq.State >= 0 {
		condition["state"] = rq.State
	}
	count := s.repo.CountOrders(condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	offset := (rq.Page - 1) * rq.Limit
	list, _ := s.repo.FetchOrders(condition, offset, rq.Limit)
	ls := make([]map[string]interface{}, 0)
	for _, item := range list {
		entry := make(map[string]interface{}, 0)
		item.SetInterface(&entry)
		if item.Get("state").ToInt() == 0 {
			entry["expiredtime"] = item.Get("createtime").ToInt() + item.Get("circle").ToInt()*86400
		}
		ls = append(ls, entry)
	}
	rs := new(shareddomain.PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.Total = count
	rs.PageTotal = pagesize
	rs.List = ls
	return rs
}

func (s *Service) Unlock(uid int, sn string) *shareddomain.BaseResponse {
	ntime := utils.GetNow()
	rs := new(shareddomain.BaseResponse)
	one := s.GetOrderInfo(uid, sn)
	if one == nil {
		rs.State = stateFailed
		rs.Msg = "no this order"
		return rs
	}
	if ntime < one["endtime"].ToInt() {
		rs.State = stateFailed
		rs.Msg = "no this order"
		return rs
	}

	data := db.DB_PARAMS{"state": 1, "unlocktime": utils.GetNow()}
	_ = s.repo.UpdateOrderByID(one["id"].ToInt(), data)
	if s.user.AddCredit(uid, &userdomain.CreditValue{
		Credit:          one["amount"].ToFloat(),
		LockCredit:      0,
		VCrdit:          0,
		LockVCredit:     0,
		UserCoinLogType: creditlogdomain.COIN_LOG_USER_MINING_BACK,
		UserCoinLogInfo: creditlogdomain.QueueCreditLog{
			Credit:     one["amount"].ToFloat(),
			LockCredit: 0,
			Sn:         one["sn"].ToString(),
			CreateTime: utils.GetInt(utils.GetJsonValue(data["unlocktime"])),
		},
		TeamCoinLogType: 0,
		TeamCoinLogInfo: nil,
	}) {
		rs.State = stateSuccess
		rs.Msg = "success"
		return rs
	}

	rs.State = stateSystemError
	rs.Msg = "system error"
	return rs
}

func (s *Service) GetOrderInfo(uid int, sn string) db.DBValues {
	one, _ := s.repo.FetchOrder(uid, sn)
	return one
}

func (s *Service) GetOrderCount(uid int) map[string]interface{} {
	rs := make(map[string]interface{})
	rs["count"] = s.repo.CountOrders(db.DB_PARAMS{"uid": uid, "state": 0})

	one, _ := s.repo.FetchOpenOrderSummary(uid)
	allinfo, _ := s.repo.FetchOrderHistorySummary(uid)
	if one != nil {
		rs["day_profit"] = one["day_profit"].ToFloat()
		rs["all_amount"] = one["all_amount"].ToFloat()
	}
	if allinfo != nil {
		rs["history_profit"] = allinfo["history_profit"].ToFloat()
	}
	return rs
}
