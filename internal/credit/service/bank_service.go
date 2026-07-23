package service

import (
	creditrepo "cointrade/internal/credit/repo"
	creditdomain "cointrade/internal/domain/credit"
	shareddomain "cointrade/internal/domain/shared"
	"cointrade/lib/db"
)

type BankCacheGateway interface {
	GetBankInfo(uid int) *creditdomain.BankInfo
	SetBankInfo(uid int, info *creditdomain.BankInfo)
	DeleteBankInfo(uid int)
}

type BankService struct {
	repo  creditrepo.BankRepository
	cache BankCacheGateway
}

func NewBankService(repo creditrepo.BankRepository, cache BankCacheGateway) *BankService {
	return &BankService{
		repo:  repo,
		cache: cache,
	}
}

func (s *BankService) BindBank(uid int, rq *creditdomain.BankInfo) *shareddomain.BaseResponse {
	if rq.Account == "" || rq.BankAddress == "" || rq.BankName == "" || rq.RealName == "" || rq.RoutNumber == "" || rq.SwiftCode == "" {
		return &shareddomain.BaseResponse{
			State: 1,
			Msg:   "the bank info is valid",
		}
	}

	data := db.DB_PARAMS{
		"uid":          uid,
		"bankname":     rq.BankName,
		"realname":     rq.RealName,
		"account":      rq.Account,
		"router_num":   rq.RoutNumber,
		"swift_code":   rq.SwiftCode,
		"bank_address": rq.BankAddress,
	}

	one, _ := s.repo.FetchBankByUID(uid)
	if one != nil {
		_ = s.repo.UpdateBankByID(one["id"].Value, data)
		s.cache.DeleteBankInfo(uid)
	} else {
		_ = s.repo.InsertBank(data)
	}

	return &shareddomain.BaseResponse{State: stateSuccess, Msg: "ok"}
}

func (s *BankService) GetBankInfo(uid int) *creditdomain.BankInfo {
	if bankInfo := s.cache.GetBankInfo(uid); bankInfo != nil && bankInfo.Account != "" {
		return bankInfo
	}

	one, _ := s.repo.FetchBankByUID(uid)
	if one == nil {
		return nil
	}

	bankInfo := &creditdomain.BankInfo{
		Account:     one["account"].ToString(),
		BankAddress: one["bank_address"].ToString(),
		BankName:    one["bankname"].ToString(),
		RealName:    one["realname"].ToString(),
		RoutNumber:  one["router_num"].ToString(),
		SwiftCode:   one["swift_code"].ToString(),
	}
	s.cache.SetBankInfo(uid, bankInfo)
	return bankInfo
}
