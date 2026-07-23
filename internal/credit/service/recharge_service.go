package service

import (
	creditrepo "cointrade/internal/credit/repo"
	creditdomain "cointrade/internal/domain/credit"
	creditlogdomain "cointrade/internal/domain/creditlog"
	shareddomain "cointrade/internal/domain/shared"
	systemdomain "cointrade/internal/domain/system"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	stateSuccess     = 0
	stateSystemError = 9999999
)

type RechargeUserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
}

type RechargeSystemGateway interface {
	GetRechargeConfig(cointype string, contract string) *systemdomain.RechargeContractConfig
	GetCoinClosePrice(pair string) float64
}

type RechargeNotifier interface {
	IncrementNotify(typ int, num int)
}

type RechargeService struct {
	repo     creditrepo.RechargeRepository
	user     RechargeUserGateway
	system   RechargeSystemGateway
	notifier RechargeNotifier
}

func NewRechargeService(
	repo creditrepo.RechargeRepository,
	user RechargeUserGateway,
	system RechargeSystemGateway,
	notifier RechargeNotifier,
) *RechargeService {
	return &RechargeService{
		repo:     repo,
		user:     user,
		system:   system,
		notifier: notifier,
	}
}

func (s *RechargeService) MakeOrderSN(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", creditdomain.RECHARGE_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}

func (s *RechargeService) GetAllRechargeAddress() db.DB_LIST_RESULT {
	list, _ := s.repo.FetchRechargeAddresses()
	return list
}

func (s *RechargeService) CreateRecharge(uid int, rq *creditdomain.RechargeRequest) *creditdomain.RechargeResponse {
	rs := new(creditdomain.RechargeResponse)
	uinfo := s.user.GetBaseInfo(uid)
	rechargeConfig := s.system.GetRechargeConfig(rq.CoinType, rq.Contract)
	if rechargeConfig == nil {
		rs.State = stateSystemError
		rs.Msg = "system error"
		return rs
	}
	if uinfo == nil {
		rs.State = creditdomain.RECHARGE_STATE_ERROR_USER
		rs.Msg = "the user is not exists"
		return rs
	}
	if rq.Amount < rechargeConfig.Min {
		rs.State = creditdomain.RECHARGE_STATE_MIN
		rs.Msg = "too small"
		return rs
	}

	rate := 1.0
	cointype := strings.ToLower(rq.CoinType)
	if cointype != "usdt" {
		rate = s.system.GetCoinClosePrice(fmt.Sprintf("%susdt", cointype))
		if rate <= 0 {
			rate = 1.0
		}
	}

	sn := s.MakeOrderSN(uid)
	insertData := db.DB_PARAMS{
		"uid":         uid,
		"sn":          sn,
		"cointype":    rq.CoinType,
		"contract":    rq.Contract,
		"type":        0,
		"credit":      rq.Amount,
		"rate":        rate,
		"fact_credit": rq.Amount * rate,
		"createtime":  utils.GetNow(),
		"info":        rechargeConfig.Address,
		"txid":        "",
		"proof":       rq.Proof,
		"address":     rechargeConfig.Address,
	}
	if err := s.repo.InsertRecharge(insertData); err != nil {
		rs.State = stateSystemError
		rs.Msg = err.Error()
		return rs
	}
	s.notifier.IncrementNotify(2, 1)
	rs.State = stateSuccess
	rs.Msg = "success"
	rs.Info = insertData
	rs.Sn = sn
	return rs
}

func (s *RechargeService) SuccessRecharge(sn string) bool {
	one, _ := s.repo.FetchRechargeBySN(sn)
	if one == nil {
		return false
	}
	ntime := utils.GetNow()
	_ = s.repo.UpdateRechargeByID(one["id"].Value, db.DB_PARAMS{"state": 1, "finishtime": ntime})
	cvalue := &userdomain.CreditValue{
		Credit:          one["fact_credit"].ToFloat(),
		VCrdit:          0,
		LockCredit:      0,
		LockVCredit:     0,
		UserCoinLogType: creditlogdomain.COIN_LOG_USER_RECHARGE,
		UserCoinLogInfo: creditlogdomain.QueueCreditLog{
			Credit:     one["fact_credit"].ToFloat(),
			LockCredit: 0,
			Sn:         sn,
			CreateTime: ntime,
		},
		TeamCoinLogType: creditlogdomain.TEAM_LOG_RECHARGE,
		TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
			Recharge:   one["fact_credit"].ToFloat(),
			CreateTime: ntime,
		},
	}
	return s.user.AddCredit(one["uid"].ToInt(), cvalue)
}

func (s *RechargeService) GetRechargeOrderBySN(sn string) db.DB_ROW_RESULT {
	one, _ := s.repo.FetchRechargeRowBySN(sn)
	return one
}

func (s *RechargeService) GetRechargeList(uid int, rq *shareddomain.PageBaseRequest) *shareddomain.PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	count := s.repo.CountRecharge(condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := s.repo.FetchRechargeRows(condition, limitstr)
	rs := new(shareddomain.PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}

func (s *RechargeService) RechargeInfo(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := s.repo.FetchRechargeInfo(uid, sn)
	return one
}
