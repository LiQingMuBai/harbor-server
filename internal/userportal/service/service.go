package service

import (
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	userportalrepo "cointrade/internal/userportal/repo"
	"cointrade/lib/db"
	"fmt"
)

const (
	stateSuccess = 0
	stateFailed  = 1
)

type Service struct {
	repo userportalrepo.Repository
}

func NewService(repo userportalrepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Welcome2() *userdomain.WelcomeResponse {
	rs := new(userdomain.WelcomeResponse)
	one, _ := s.repo.FetchWelcome()
	if one == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	rs.WelcomeInfo = &userdomain.WelcomeInfo{
		PlatformName: one["platform_name"].ToString(),
		WelcomePage:  one["welcome_page"].ToString(),
	}
	return rs
}

func (s *Service) Welcome() *userdomain.WelcomeResponse {
	rs := new(userdomain.WelcomeResponse)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	welcomeInfo := &userdomain.WelcomeInfo{
		DirectWithdraw: "0",
		LinkWallet:     "0",
	}
	list, _ := s.repo.FetchSystemConfigList()
	for _, item := range list {
		switch item.Get("key").ToString() {
		case "sitename":
			welcomeInfo.PlatformName = item.Get("value").ToString()
		case "domain":
			welcomeInfo.WelcomePage = item.Get("value").ToString()
		case "vip_contact":
			welcomeInfo.VIP = item.Get("value").ToString()
		case "direct_withdraw":
			welcomeInfo.DirectWithdraw = item.Get("value").ToString()
		case "link_wallet":
			welcomeInfo.LinkWallet = item.Get("value").ToString()
		}
	}
	rs.WelcomeInfo = welcomeInfo
	return rs
}

func (s *Service) CrossTrade(uid int, data db.DB_PARAMS) *userdomain.CrossPlatformTradeResponse {
	rs := new(userdomain.CrossPlatformTradeResponse)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	userAsset, _ := s.repo.FetchUserByID(uid)
	var userBalance float64
	if userAsset != nil {
		userBalance = userAsset["credit"].ToFloat()
	}
	amount, _ := data["amount"].(float64)
	if userBalance <= amount {
		rs.State = stateFailed
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
		return rs
	}
	rest := userBalance - amount
	if err := s.repo.UpdateUserByID(uid, db.DB_PARAMS{"credit": rest}); err != nil {
		rs.State = stateFailed
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
		return rs
	}
	if err := s.repo.InsertCrossTrade(data); err != nil {
		rs.State = stateFailed
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
	}
	return rs
}
