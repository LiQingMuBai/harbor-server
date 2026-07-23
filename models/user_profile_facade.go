package models

import (
	userdomain "cointrade/internal/domain/user"
	userprofilerepo "cointrade/internal/userprofile/repo"
	userprofileservice "cointrade/internal/userprofile/service"
	"cointrade/lib/db"
)

type userProfileGateway struct{}

func (userProfileGateway) Update(uid int, data db.DB_PARAMS) {
	MODEL_USER.Update(uid, data)
}

func (userProfileGateway) ClearCache(uid int) {
	MODEL_USER.ClearCache(uid)
}

var userProfileSvc = userprofileservice.NewService(
	userprofilerepo.NewDBRepository(),
	userProfileGateway{},
)

func userProfileGetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return userProfileSvc.GetBaseInfo(uid, MODEL_USER.MakeCacheId)
}
