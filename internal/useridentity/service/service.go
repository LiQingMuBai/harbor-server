package service

import (
	useridentityrepo "cointrade/internal/useridentity/repo"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
	"cointrade/utils"
)

const (
	stateSuccess     = 0
	stateFailed      = 1
	stateSystemError = 9999999
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	Update(uid int, data db.DB_PARAMS)
}

type Notifier interface {
	IncrementNotify(typ int, num int)
}

type Service struct {
	repo     useridentityrepo.Repository
	user     UserGateway
	notifier Notifier
}

func NewService(repo useridentityrepo.Repository, user UserGateway, notifier Notifier) *Service {
	return &Service{
		repo:     repo,
		user:     user,
		notifier: notifier,
	}
}

func (s *Service) AuthLv1(uid int, authInfo *userdomain.AuthLv1Request) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if uinfo.AuthLv >= 1 {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	if authInfo.CardBack == "" || authInfo.CardFront == "" || authInfo.Phone == "" || authInfo.CardType == 0 {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgInvalidParams
		return rs
	}
	insertData := db.DB_PARAMS{
		"uid":           uid,
		"realname":      authInfo.Name,
		"inid":          authInfo.IdCard,
		"card_front":    authInfo.CardFront,
		"card_back":     authInfo.CardBack,
		"card_hand":     authInfo.HandCard,
		"process_state": 0,
		"createtime":    utils.GetNow(),
		"phone":         authInfo.Phone,
		"card_type":     authInfo.CardType,
	}
	one, _ := s.repo.FetchLv1ByUID(uid)
	if one != nil {
		if one["process_state"].ToInt() == 2 {
			_ = s.repo.UpdateLv1ByUID(uid, insertData)
		} else {
			rs.State = stateFailed
			rs.Msg = shareddomain.MsgFailed
			return rs
		}
	} else {
		_, _ = s.repo.InsertLv1(insertData)
	}
	s.notifier.IncrementNotify(3, 1)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) AuthLv2(uid int, rq *userdomain.AuthLv2Request) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if uinfo.AuthLv >= 2 || uinfo.AuthLv == 0 {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	if rq.FarmilyName == "" || rq.WalletAddress == "" || rq.Address == "" {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgInvalidParams
		return rs
	}
	data := db.DB_PARAMS{
		"uid":               uid,
		"farmily_name":      rq.FarmilyName,
		"relation":          rq.Relation,
		"address":           rq.Address,
		"contact":           rq.Contact,
		"wallet_address":    rq.WalletAddress,
		"chaintype":         rq.ChainType,
		"second_card_front": rq.Second_card_front,
		"second_card_back":  rq.Second_card_Hand,
		"second_card_hand":  rq.Second_card_Hand,
		"createtime":        utils.GetNow(),
		"state":             0,
	}
	one, _ := s.repo.FetchLv2ByUID(uid)
	if one != nil {
		if one["state"].ToInt() == 2 {
			_ = s.repo.UpdateLv2ByID(one["id"].Value, data)
		} else {
			rs.State = stateFailed
			rs.Msg = shareddomain.MsgFailed
			return rs
		}
	} else {
		_, _ = s.repo.InsertLv2(data)
	}
	s.notifier.IncrementNotify(4, 1)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) GetAuthInfo(uid int) map[int]interface{} {
	lv1Info, _ := s.repo.FetchLv1RowByUID(uid)
	lv2Info, _ := s.repo.FetchLv2RowByUID(uid)
	return map[int]interface{}{1: lv1Info, 2: lv2Info}
}

func (s *Service) AdminSaveLv1(uid int, realName string, inid string, cardFront string, cardBack string, cardHand string) error {
	if uid == 0 {
		return ErrInvalidUID
	}
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		return ErrUserNotFound
	}
	if exists := s.repo.CountLv1(db.DB_PARAMS{"uid": uid}); exists > 0 {
		return ErrAuthAlreadyExists
	}
	if len(realName) < 4 || len(realName) > 40 {
		return ErrRealNameTooLong
	}
	if cardFront == "" || cardBack == "" || cardHand == "" {
		return ErrCardPhotosRequired
	}
	insertData := db.DB_PARAMS{
		"uid":           uid,
		"realname":      realName,
		"inid":          inid,
		"card_front":    cardFront,
		"card_back":     cardBack,
		"card_hand":     cardHand,
		"process_state": 1,
		"createtime":    utils.GetNow(),
		"passtime":      utils.GetNow(),
	}
	_, err := s.repo.InsertLv1(insertData)
	return err
}

func (s *Service) AdminDeleteAuth(id int, tp int) error {
	if id == 0 {
		return ErrInvalidID
	}
	switch tp {
	case 1:
		one, _ := s.repo.FetchLv1ByID(id)
		if one == nil {
			return ErrAuthNotFound
		}
		if err := s.repo.DeleteLv1ByID(id); err != nil {
			return err
		}
		s.user.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp - 1})
		return nil
	case 2:
		one, _ := s.repo.FetchLv2ByID(id)
		if one == nil {
			return ErrAuthNotFound
		}
		if err := s.repo.DeleteLv2ByID(id); err != nil {
			return err
		}
		s.user.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp - 1})
		return nil
	default:
		return ErrInvalidType
	}
}

func (s *Service) AdminReviewAuth(id int, tp int, processState int, reason string) error {
	if id == 0 {
		return ErrInvalidID
	}
	switch tp {
	case 1:
		one, _ := s.repo.FetchLv1ByID(id)
		if one == nil {
			return ErrAuthNotFound
		}
		up := db.DB_PARAMS{"process_state": processState, "reason": reason, "passtime": utils.GetNow()}
		if err := s.repo.UpdateLv1ByID(id, up); err != nil {
			return err
		}
		if processState == 1 {
			s.user.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp, "credit_coin": 80})
		}
		return nil
	case 2:
		one, _ := s.repo.FetchLv2ByID(id)
		if one == nil {
			return ErrAuthNotFound
		}
		up := db.DB_PARAMS{"state": processState, "reason": reason, "passtime": utils.GetNow()}
		if err := s.repo.UpdateLv2ByID(id, up); err != nil {
			return err
		}
		if processState == 1 {
			s.user.Update(one.Get("uid").ToInt(), db.DB_PARAMS{"auth_lv": tp, "credit_coin": 100})
		}
		return nil
	default:
		return ErrInvalidType
	}
}
