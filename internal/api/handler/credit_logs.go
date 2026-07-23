package handler

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *CreditModule) logRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/credit/logs", Handles: common.HandleArray{m.NeedLogin, m.LogList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/usercount", Handles: common.HandleArray{m.NeedLogin, m.UserCount}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/levelcount", Handles: common.HandleArray{m.NeedLogin, m.UserLevelCount}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/income", Handles: common.HandleArray{m.NeedLogin, m.IncomeLog}},
	}
}

func (m *CreditModule) IncomeLog(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.IncomeLog(uid, &rq))
}

func (m *CreditModule) LogList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.CoinLogRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetList(uid, &rq))
}

func (m *CreditModule) UserCount(r *gin.Context) {
	uid := r.GetInt("uid")
	timeType := r.GetInt("type")
	switch timeType {
	case models.LOG_TIMETYPE_ALL:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountSum(uid))
	case models.LOG_TIMETYPE_DAY:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountDay(uid))
	case models.LOG_TIMETYPE_MONTH:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountMonth(uid))
	}
}

func (m *CreditModule) UserLevelCount(r *gin.Context) {
	uid := r.GetInt("uid")
	timeType := r.GetInt("type")
	switch timeType {
	case models.LOG_TIMETYPE_ALL:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountSum(uid))
	case models.LOG_TIMETYPE_DAY:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountDay(uid))
	case models.LOG_TIMETYPE_MONTH:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountMonth(uid))
	}
}
