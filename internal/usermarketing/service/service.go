package service

import (
	creditlogdomain "cointrade/internal/domain/creditlog"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	usermarketingrepo "cointrade/internal/usermarketing/repo"
	"cointrade/lib/db"
	"cointrade/utils"
	"time"
)

const (
	stateSuccess = 0
	stateFailed  = 1
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
	ClearCache(uid int)
	EncodePassword(password string) string
}

type Service struct {
	repo usermarketingrepo.Repository
	user UserGateway
}

func NewService(repo usermarketingrepo.Repository, user UserGateway) *Service {
	return &Service{
		repo: repo,
		user: user,
	}
}

func (s *Service) IsNewUser(uid int) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	if exists := s.repo.Count("users", db.DB_PARAMS{"id": uid, "approve_state": 1}); exists == 0 {
		rs.State = stateSuccess
	} else {
		rs.State = stateFailed
	}
	return rs
}

func (s *Service) Claim(uid int) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	if exists := s.repo.Count("users_mechanism", db.DB_PARAMS{"user_id": uid}); exists == 0 {
		utils.ServiceInfo("new user claim bonus:", uid, "add 10 usdt")
		s.user.AddCredit(uid, &userdomain.CreditValue{
			Credit:          100,
			UserCoinLogType: 4100001,
			UserCoinLogInfo: creditlogdomain.QueueCreditLog{
				Credit:     100,
				CoinType:   "usdt",
				CreateTime: utils.GetNow(),
			},
		})
		if err := s.repo.Insert("users_mechanism", db.DB_PARAMS{"user_id": uid}); err != nil {
			utils.ServiceError("insert user mechanism record failed:", err)
		}
		utils.ServiceInfo("cache new user mechanism record:", uid)
	}
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) ClearIncome(uid int, cashPassword string) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	nowDay := time.Now().Day()
	oldDay := time.Unix(int64(uinfo.ClearIncomeTime), 0).Day()
	if uinfo.CashPassword == "" || uinfo.CashPassword != s.user.EncodePassword(cashPassword) || oldDay == nowDay {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	if uinfo.MiningIncome+uinfo.RechargeIncome > 0 {
		if s.user.AddCredit(uid, &userdomain.CreditValue{
			Credit:          uinfo.MiningIncome + uinfo.RechargeIncome,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: 4100004,
			UserCoinLogInfo: creditlogdomain.QueueCreditLog{
				Credit:     uinfo.MiningIncome + uinfo.RechargeIncome,
				LockCredit: 0,
				Sn:         "",
				CreateTime: utils.GetNow(),
			},
		}) {
			_ = s.repo.AddValue("users", map[string]float64{
				"recharge_income": -1 * uinfo.RechargeIncome,
				"mining_income":   -1 * uinfo.MiningIncome,
			}, db.DB_PARAMS{"id": uid})
			s.user.ClearCache(uid)
		}
	}
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}
