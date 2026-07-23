package service

import (
	userdomain "cointrade/internal/domain/user"
	shareddomain "cointrade/internal/domain/shared"
	usersecurityrepo "cointrade/internal/usersecurity/repo"
	"cointrade/lib/db"
)

const (
	stateSuccess            = 0
	stateFailed             = 1
	registerStateBadPass    = 2002
	changePassStateOldError = 1002
	bindPhoneStateBound     = 4001
	bindPhoneStateBadCode   = 4002
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	Update(uid int, data db.DB_PARAMS)
}

type CodeGateway interface {
	GetBindSMSCode(uid int, phone string) string
	GetEmailBindCode(email string) string
}

type Service struct {
	repo  usersecurityrepo.Repository
	user  UserGateway
	codes CodeGateway
}

func NewService(repo usersecurityrepo.Repository, user UserGateway, codes CodeGateway) *Service {
	return &Service{
		repo:  repo,
		user:  user,
		codes: codes,
	}
}

func (s *Service) ChangeCashPassword(uid int, rq *userdomain.SetCashPasswordRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	if rq.Password == "" {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgPasswordRequired
		return rs
	}
	if len(rq.Password) < 6 || len(rq.Password) > 20 {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if rq.Password != rq.RePassword {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	s.user.Update(uid, db.DB_PARAMS{"cash_password": rq.Password})
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) UpdateCashPassword(uid int, rq *userdomain.UpdateCashPasswordRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if rq.O_Password != uinfo.CashPassword {
		rs.State = changePassStateOldError
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	if len(rq.N_Password) < 6 || len(rq.N_Password) > 20 {
		rs.State = registerStateBadPass
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if rq.N_Password != rq.R_Password {
		rs.State = registerStateBadPass
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	s.user.Update(uid, db.DB_PARAMS{"cash_password": rq.N_Password})
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) BindPhone(uid int, phone string, code string) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if uinfo.Phone != "" {
		rs.State = bindPhoneStateBound
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"phone": phone}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = bindPhoneStateBound
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	verifyCode := s.codes.GetBindSMSCode(uid, phone)
	if verifyCode != code {
		rs.State = bindPhoneStateBadCode
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	s.user.Update(uid, db.DB_PARAMS{"phone": phone})
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) BindEmail(uid int, email string, code string) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if uinfo.Email != "" && uinfo.Email != "0" {
		rs.State = bindPhoneStateBound
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"email": email}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = bindPhoneStateBound
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	verifyCode := s.codes.GetEmailBindCode(email)
	if verifyCode != code {
		rs.State = bindPhoneStateBadCode
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	s.user.Update(uid, db.DB_PARAMS{"email": email})
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}
