package service

import (
	creditrepo "cointrade/internal/credit/repo"
	assetsdomain "cointrade/internal/domain/assets"
	creditdomain "cointrade/internal/domain/credit"
	creditlogdomain "cointrade/internal/domain/creditlog"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const transferListLimit = 15

type TransferUserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
	RechargeByApprove(uid int, amount float64) *shareddomain.BaseResponse
}

type TransferAssetGateway interface {
	GetOneAsset(uid int, coin string) *assetsdomain.AssetInfo
	AddAssets(uid int, asset *assetsdomain.Assets) bool
}

type TransferConfigGateway interface {
	GetMaxWithdrawNum() int
}

type TransferService struct {
	repo   creditrepo.TransferRepository
	user   TransferUserGateway
	asset  TransferAssetGateway
	config TransferConfigGateway
}

func NewTransferService(
	repo creditrepo.TransferRepository,
	user TransferUserGateway,
	asset TransferAssetGateway,
	config TransferConfigGateway,
) *TransferService {
	return &TransferService{
		repo:   repo,
		user:   user,
		asset:  asset,
		config: config,
	}
}

func (s *TransferService) AddWallet(uid int, rq *creditdomain.WalletAddressRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	if rq.Address == "" || rq.CoinType == "" || rq.Contract == "" {
		rs.State = 1
		rs.Msg = "faild"
		return rs
	}
	if s.repo.CountWallet(rq.CoinType, rq.Contract) >= 10 {
		rs.State = 1
		rs.Msg = "too more"
		return rs
	}
	insertData := db.DB_PARAMS{
		"uid":        uid,
		"createtime": utils.GetNow(),
		"cointype":   rq.CoinType,
		"contract":   rq.Contract,
		"address":    rq.Address,
		"title":      rq.Title,
	}
	_ = s.repo.InsertWallet(insertData)
	rs.State = 0
	rs.Msg = "success"
	return rs
}

func (s *TransferService) DeleteWallet(uid int, id int) *shareddomain.BaseResponse {
	_ = s.repo.DeleteWallet(uid, id)
	return &shareddomain.BaseResponse{State: 0, Msg: "success"}
}

func (s *TransferService) GetWalletList(uid int) db.DB_LIST_RESULT {
	list, _ := s.repo.FetchWalletList(uid)
	return list
}

func (s *TransferService) Transfer(uid int, trans *creditdomain.TransferRequest) *shareddomain.BaseResponse {
	ntime := utils.GetNow()
	uinfo := s.user.GetBaseInfo(uid)
	trans.Coin = strings.ToLower(trans.Coin)
	toAddress := ""

	if uinfo == nil {
		return &shareddomain.BaseResponse{State: stateSystemError, Msg: "SYSTEM ERROR"}
	}
	if trans.Amount <= 0 {
		return &shareddomain.BaseResponse{State: creditdomain.RECHARGE_STATE_MIN, Msg: "too min"}
	}
	if strings.Index(trans.Coin, "usdt") >= 0 {
		trans.Coin = "usdt"
	}

	sn := s.makeTransferOrderSN(uid)
	insertData := db.DB_PARAMS{"uid": uid, "sn": sn, "createtime": ntime}
	assetInfo := s.asset.GetOneAsset(uid, trans.Coin)

	if trans.Coin != "usdt" && trans.Coin != "usdc" && assetInfo == nil {
		return &shareddomain.BaseResponse{State: 1, Msg: "error assets"}
	}
	if trans.Coin != "usdt" && trans.Coin != "usdc" && assetInfo != nil {
		if assetInfo.IsTrans != 1 {
			return &shareddomain.BaseResponse{State: 1, Msg: "error assets  12"}
		}
		toAddress = assetInfo.Address
	} else if trans.ToAddress != "" {
		toAddress = trans.ToAddress
	} else {
		toAddress = uinfo.WalletAddress
	}

	if trans.Direction == creditdomain.TRANSFER_DIRECTION_OUT {
		toAddress = trans.ToAddress
		if toAddress == "" {
			return &shareddomain.BaseResponse{State: creditdomain.RECHARGE_STATE_ERROR_ADDRESS, Msg: "error address"}
		}
		day := time.Now()
		today := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local).Unix()
		count := s.repo.CountTodayOutTransfer(uinfo.Id, today)
		if count > s.config.GetMaxWithdrawNum() {
			return &shareddomain.BaseResponse{State: creditdomain.RECHARGE_STATE_ERROR_MAX_WITHDRAW, Msg: "max withdraw"}
		}

		if trans.Coin != "usdt" {
			if assetInfo.Count < trans.Amount {
				return &shareddomain.BaseResponse{State: creditdomain.RECHARGE_STATE_ERROR_MONEY, Msg: "not enough assets"}
			}
			s.asset.AddAssets(uid, &assetsdomain.Assets{
				Coin:    trans.Coin,
				Pair:    trans.Coin + "usdt",
				Num:     -1 * trans.Amount,
				LockNum: trans.Amount,
				Mode:    1,
			})
			s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: creditlogdomain.COIN_LOG_USER_WITHDRAW,
				UserCoinLogInfo: creditlogdomain.QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   trans.Coin,
				},
			})
		} else {
			if trans.Amount > uinfo.Credit {
				return &shareddomain.BaseResponse{State: creditdomain.RECHARGE_STATE_ERROR_MONEY, Msg: "not enough assets"}
			}
			if !s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: creditlogdomain.COIN_LOG_EXCHANGE_ACCOUNT_OUT,
				UserCoinLogInfo: creditlogdomain.QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   "usdt",
				},
			}) {
				return &shareddomain.BaseResponse{State: stateSystemError, Msg: "ERROR addcredit"}
			}
		}
	}

	insertData["to_address"] = toAddress
	insertData["direction"] = trans.Direction
	insertData["amount"] = trans.Amount
	insertData["coin_symbol"] = trans.Coin
	if err := s.repo.InsertTransfer(insertData); err == nil {
		return &shareddomain.BaseResponse{State: stateSuccess, Msg: "ok"}
	}
	return &shareddomain.BaseResponse{State: stateSystemError, Msg: "error"}
}

func (s *TransferService) ExchangeAccount(uid int, rq *creditdomain.ExchangeAccountRequest) *shareddomain.BaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		return &shareddomain.BaseResponse{State: stateSystemError, Msg: "SYSTEM ERROR"}
	}
	if rq.Drection == creditdomain.EXCHANGE_DIRECTION_CONTRACT {
		return s.user.RechargeByApprove(uid, rq.Amount)
	}
	if rq.Drection == creditdomain.EXCHANGE_DIRECTION_ACCOUNT {
		if uinfo.IsWithDraw == 0 {
			return &shareddomain.BaseResponse{State: creditdomain.WIDTHDRAW_STATE_ERROR_LOCKED, Msg: uinfo.WithDrawMsg}
		}
		return s.Transfer(uid, &creditdomain.TransferRequest{
			Coin:      "usdt",
			Amount:    rq.Amount,
			Direction: creditdomain.TRANSFER_DIRECTION_OUT,
			ToAddress: uinfo.WalletAddress,
		})
	}
	return nil
}

func (s *TransferService) ExchangeAccount2(uid int, rq *creditdomain.ExchangeAccountRequest2) *shareddomain.BaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		return &shareddomain.BaseResponse{State: stateSystemError, Msg: "SYSTEM ERROR"}
	}
	if uinfo.IsWithDraw == 0 {
		return &shareddomain.BaseResponse{State: creditdomain.WIDTHDRAW_STATE_ERROR_LOCKED, Msg: uinfo.WithDrawMsg}
	}

	amount, err := strconv.ParseFloat(rq.Amount, 64)
	if err != nil {
		return &shareddomain.BaseResponse{State: stateSystemError, Msg: "SYSTEM ERROR"}
	}

	if rq.Symbol == "" || strings.ToLower(rq.Symbol) != "usdc" {
		rq.Symbol = "usdt"
	}
	return s.Transfer(uid, &creditdomain.TransferRequest{
		Coin:      strings.ToLower(rq.Symbol),
		Amount:    amount,
		Direction: creditdomain.TRANSFER_DIRECTION_OUT,
		ToAddress: rq.Address,
	})
}

func (s *TransferService) TransferLogs(uid int, rq *creditdomain.TransFerLogsRequest) *shareddomain.PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	if rq.Direction > -1 {
		condition["direction"] = rq.Direction
	}
	count := s.repo.CountTransfer(condition)
	pagesize := int(math.Ceil(float64(count) / float64(transferListLimit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	list, _ := s.repo.FetchTransferRows(condition, fmt.Sprintf("limit %d,%d", (rq.Page-1)*transferListLimit, transferListLimit))
	rs := new(shareddomain.PageBaseResponse)
	rs.Limit = transferListLimit
	rs.State = stateSuccess
	rs.Msg = "ok"
	rs.Total = count
	rs.PageTotal = pagesize
	rs.Page = rq.Page
	rs.List = list
	return rs
}

func (s *TransferService) TransferDetail(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := s.repo.FetchTransferDetail(uid, sn)
	return one
}

func (s *TransferService) makeTransferOrderSN(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", creditdomain.TRANSFER_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}
