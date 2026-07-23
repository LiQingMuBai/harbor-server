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

type WithdrawUserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
}

type WithdrawSystemGateway interface {
	GetRechargeConfig(cointype string, contract string) *systemdomain.RechargeContractConfig
	GetCoinClosePrice(pair string) float64
	GetMinWithdraw() float64
	GetWithdrawFee() float64
}

type WithdrawBankGateway interface {
	GetBankInfo(uid int) *creditdomain.BankInfo
}

type WithdrawNotifier interface {
	IncrementNotify(typ int, num int)
}

type WithdrawService struct {
	repo     creditrepo.WithdrawRepository
	user     WithdrawUserGateway
	system   WithdrawSystemGateway
	bank     WithdrawBankGateway
	notifier WithdrawNotifier
}

func NewWithdrawService(
	repo creditrepo.WithdrawRepository,
	user WithdrawUserGateway,
	system WithdrawSystemGateway,
	bank WithdrawBankGateway,
	notifier WithdrawNotifier,
) *WithdrawService {
	return &WithdrawService{
		repo:     repo,
		user:     user,
		system:   system,
		bank:     bank,
		notifier: notifier,
	}
}

func (s *WithdrawService) MakeOrderSN(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", creditdomain.WITHDRAW_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}

func (s *WithdrawService) CreateWithdraw(uid int, rq *creditdomain.WithdrawRequest) *creditdomain.RechargeResponse {
	rs := new(creditdomain.RechargeResponse)
	sn := s.MakeOrderSN(uid)
	ntime := utils.GetNow()
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = creditdomain.WITHDRAW_STATE_ERROR_USER
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if uinfo.IsWithDraw != 1 {
		rs.State = creditdomain.WITHDRAW_STATE_ERROR_LOCKED
		rs.Msg = shareddomain.MsgWithdrawLocked
		return rs
	}

	cointype := strings.ToLower(rq.CoinType)
	if cointype != "bank" {
		rechargeConfig := s.system.GetRechargeConfig(rq.CoinType, rq.Contract)
		if rechargeConfig == nil {
			rs.State = stateSystemError
			rs.Msg = shareddomain.MsgInternalError
			return rs
		}
		if rq.Amount < rechargeConfig.Min {
			rs.State = creditdomain.WITHDRAW_STATE_MIN
			rs.Msg = shareddomain.MsgAmountTooSmall
			return rs
		}
	} else if rq.Amount < s.system.GetMinWithdraw() {
		rs.State = creditdomain.WITHDRAW_STATE_MIN
		rs.Msg = shareddomain.MsgAmountTooSmall
		return rs
	}

	rate := 1.0
	var bankInfo *creditdomain.BankInfo
	if cointype != "usdt" {
		if cointype != "bank" {
			rate = s.system.GetCoinClosePrice(fmt.Sprintf("%susdt", cointype))
			if rate <= 0 {
				rate = 1.0
			}
		} else {
			bankInfo = s.bank.GetBankInfo(uid)
			if bankInfo == nil {
				rs.State = creditdomain.WITHDRAW_STATE_ERROR_NOTBINDBANK
				rs.Msg = shareddomain.MsgBankNotBound
				return rs
			}
		}
	}

	factCredit := rq.Amount * rate
	if uinfo.Credit < factCredit {
		rs.State = creditdomain.WITHDRAW_STATE_NOTENOUGH
		rs.Msg = shareddomain.MsgInsufficient
		return rs
	}

	insertData := db.DB_PARAMS{
		"uid":         uid,
		"credit":      rq.Amount,
		"rate":        rate,
		"fact_credit": factCredit,
		"cointype":    rq.CoinType,
		"contract":    rq.Contract,
		"address":     rq.Address,
		"fee":         factCredit * s.system.GetWithdrawFee() / float64(100),
		"info":        "",
		"createtime":  ntime,
		"sn":          sn,
		"state":       0,
		"finishtime":  0,
		"memo":        "",
	}
	if cointype == "bank" {
		insertData["type"] = 1
		insertData["bankinfo"] = bankInfo
	}
	if err := s.repo.InsertWithdraw(insertData); err != nil {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}
	if !s.user.AddCredit(uid, &userdomain.CreditValue{
		Credit:          -1 * factCredit,
		LockCredit:      factCredit,
		LockVCredit:     0,
		VCrdit:          0,
		UserCoinLogType: creditlogdomain.COIN_LOG_USER_WITHDRAW,
		UserCoinLogInfo: creditlogdomain.QueueCreditLog{
			Credit:     -1 * factCredit,
			LockCredit: factCredit,
			Sn:         sn,
			CreateTime: ntime,
		},
	}) {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}
	s.notifier.IncrementNotify(1, 1)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgOK
	rs.Info = insertData
	rs.Sn = sn
	return rs
}

func (s *WithdrawService) GetWithdrawList(uid int, rq *shareddomain.PageBaseRequest) *shareddomain.PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	count := s.repo.CountWithdraw(condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := s.repo.FetchWithdrawRows(condition, limitstr)
	rs := new(shareddomain.PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}

func (s *WithdrawService) WithdrawInfo(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := s.repo.FetchWithdrawInfo(uid, sn)
	return one
}
