package models

import (
	userportalrepo "cointrade/internal/userportal/repo"
	userportalservice "cointrade/internal/userportal/service"
)

var userPortalSvc = userportalservice.NewService(
	userportalrepo.NewDBRepository(),
)
