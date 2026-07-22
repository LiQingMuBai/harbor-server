package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AdminUserModule struct {
	common.ModuleBase
}

func (m *AdminUserModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/login", Handles: common.HandleArray{m.Login}},
		&common.ModuleHandles{Method: "post", Path: "/admin/logout", Handles: common.HandleArray{m.Logout}},
		&common.ModuleHandles{Method: "get", Path: "/admin/token_info", Handles: common.HandleArray{m.TokenSid}},
		&common.ModuleHandles{Method: "post", Path: "/admin/userlist", Handles: common.HandleArray{m.CheckLogin, m.UserList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_user", Handles: common.HandleArray{m.CheckLogin, m.OpUser}},
		&common.ModuleHandles{Method: "post", Path: "/admin/usercoinlog", Handles: common.HandleArray{m.CheckLogin, m.UserCoinLog}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coin/change", Handles: common.HandleArray{m.CheckLogin, m.ChangeCoin}},

		&common.ModuleHandles{Method: "post", Path: "/admin/user_asset", Handles: common.HandleArray{m.CheckLogin, m.UserAsset}},
		&common.ModuleHandles{Method: "post", Path: "/admin/userinfo", Handles: common.HandleArray{m.CheckLogin, m.UserInfo}},
		&common.ModuleHandles{Method: "post", Path: "/admin/user_levelcount", Handles: common.HandleArray{m.CheckLogin, m.UserLevelCount}},

		&common.ModuleHandles{Method: "post", Path: "/admin/admin_list", Handles: common.HandleArray{m.CheckLogin, m.AdminList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/add_manage", Handles: common.HandleArray{m.CheckLogin, m.AddAdmin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/mean_list", Handles: common.HandleArray{m.CheckLogin, m.MeanRouter}},

		&common.ModuleHandles{Method: "get", Path: "/admin/mean_router", Handles: common.HandleArray{m.CheckLogin, m.MeanRouter}},
		&common.ModuleHandles{Method: "get", Path: "/admin/auth_router", Handles: common.HandleArray{m.AuthRouter}},
		&common.ModuleHandles{Method: "post", Path: "/admin/role_list", Handles: common.HandleArray{m.CheckLogin, m.RoleList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_role", Handles: common.HandleArray{m.CheckLogin, m.HandlerRole}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_mean", Handles: common.HandleArray{m.CheckLogin, m.HandlerMean}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_mean", Handles: common.HandleArray{m.CheckLogin, m.DelMean}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coin_list", Handles: common.HandleArray{m.CheckLogin, m.CoinList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coindesc_list", Handles: common.HandleArray{m.CheckLogin, m.CoinDescList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_coindesc", Handles: common.HandleArray{m.CheckLogin, m.OpCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coindesc", Handles: common.HandleArray{m.CheckLogin, m.DelCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/save_coin", Handles: common.HandleArray{m.CheckLogin, m.SaveCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coin", Handles: common.HandleArray{m.CheckLogin, m.DelCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/currency_list", Handles: common.HandleArray{m.CheckLogin, m.CurrencyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_currency", Handles: common.HandleArray{m.CheckLogin, m.OpCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_currency", Handles: common.HandleArray{m.CheckLogin, m.DelCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/explode_list", Handles: common.HandleArray{m.ExplodeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_explode", Handles: common.HandleArray{m.CheckLogin, m.OpExplodeTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_explode", Handles: common.HandleArray{m.CheckLogin, m.DelExplode}},

		&common.ModuleHandles{Method: "post", Path: "/admin/minner_open", Handles: common.HandleArray{m.CheckLogin, m.MinnerOpen}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_list", Handles: common.HandleArray{m.CheckLogin, m.MinnerList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_op", Handles: common.HandleArray{m.CheckLogin, m.OpMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_del", Handles: common.HandleArray{m.CheckLogin, m.DelMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_list", Handles: common.HandleArray{m.CheckLogin, m.NoticeList}},

		&common.ModuleHandles{Method: "post", Path: "/admin/notice_op", Handles: common.HandleArray{m.CheckLogin, m.OPNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_del", Handles: common.HandleArray{m.CheckLogin, m.DelNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/recharge_list", Handles: common.HandleArray{m.CheckLogin, m.RechargeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/recharge_op", Handles: common.HandleArray{m.CheckLogin, m.OpRecharge}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_list", Handles: common.HandleArray{m.CheckLogin, m.AddrList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_op", Handles: common.HandleArray{m.CheckLogin, m.OpAddr}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_del", Handles: common.HandleArray{m.CheckLogin, m.DelAddr}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_openclose", Handles: common.HandleArray{m.OpenAddr}},

		&common.ModuleHandles{Method: "post", Path: "/admin/user_wallet", Handles: common.HandleArray{m.CheckLogin, m.UserWallet}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_userwallet", Handles: common.HandleArray{m.CheckLogin, m.DeluserWallet}},
		&common.ModuleHandles{Method: "post", Path: "/admin/withdraw_list", Handles: common.HandleArray{m.CheckLogin, m.WithdrawList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/withdraw_op", Handles: common.HandleArray{m.CheckLogin, m.WithdrawOp}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_list", Handles: common.HandleArray{m.CheckLogin, m.UserAuthList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_del", Handles: common.HandleArray{m.CheckLogin, m.UauthDel}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_op", Handles: common.HandleArray{m.CheckLogin, m.UauthOp}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_list", Handles: common.HandleArray{m.UserAuth2List}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_op", Handles: common.HandleArray{m.CheckLogin, m.Uauth2Op}},
		&common.ModuleHandles{Method: "post", Path: "/admin/save_withdraw", Handles: common.HandleArray{m.CheckLogin, m.SaveWithdraw}},

		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_del", Handles: common.HandleArray{m.CheckLogin, m.Uauth2Del}},
		&common.ModuleHandles{Method: "post", Path: "/admin/morder_list", Handles: common.HandleArray{m.CheckLogin, m.MorderList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/stop_minner", Handles: common.HandleArray{m.CheckLogin, m.StopMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/close_trade_list", Handles: common.HandleArray{m.CheckLogin, m.CloseTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/hold_trade_list", Handles: common.HandleArray{m.CheckLogin, m.HoldTradeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_now", Handles: common.HandleArray{m.CheckLogin, m.DelegateNow}},
		&common.ModuleHandles{Method: "post", Path: "/admin/manual_delegate_trade", Handles: common.HandleArray{m.CheckLogin, m.ManualDelegateTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_history", Handles: common.HandleArray{m.CheckLogin, m.DelegateHistory}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_del", Handles: common.HandleArray{m.CheckLogin, m.DelegateHistoryDel}},
		&common.ModuleHandles{Method: "post", Path: "/admin/setting", Handles: common.HandleArray{m.CheckLogin, m.Setting}},
		&common.ModuleHandles{Method: "post", Path: "/admin/sitecount", Handles: common.HandleArray{m.CheckLogin, m.SiteCount}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notify_list", Handles: common.HandleArray{m.CheckLogin, m.NotifyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/clearNotify", Handles: common.HandleArray{m.ClearNotify}},
		&common.ModuleHandles{Method: "post", Path: "/admin/statictis_count", Handles: common.HandleArray{m.CheckLogin, m.StatictisCount}},
		&common.ModuleHandles{Method: "post", Path: "/admin/chatlist", Handles: common.HandleArray{m.CheckLogin, m.Chatlist}},
		&common.ModuleHandles{Method: "post", Path: "/admin/custom_msg", Handles: common.HandleArray{m.CheckLogin, m.UserByMessage}},
		&common.ModuleHandles{Method: "post", Path: "/admin/send_msg", Handles: common.HandleArray{m.CheckLogin, m.SendMsg}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulelist", Handles: common.HandleArray{m.CheckLogin, m.Rulelist}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulehandler", Handles: common.HandleArray{m.RuleHandler}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_rule", Handles: common.HandleArray{m.CheckLogin, m.DelRule}},
		&common.ModuleHandles{Method: "post", Path: "admin/kick_user", Handles: common.HandleArray{m.CheckLogin, m.KickUser}},

		&common.ModuleHandles{Method: "post", Path: "admin/controller_list", Handles: common.HandleArray{m.CheckLogin, m.ControllerList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_controller", Handles: common.HandleArray{m.CheckLogin, m.DelController}},

		&common.ModuleHandles{Method: "post", Path: "admin/kline_controller", Handles: common.HandleArray{m.CheckLogin, m.KlineController}},
		&common.ModuleHandles{Method: "post", Path: "admin/explode_controller", Handles: common.HandleArray{m.CheckLogin, m.ExplodeController}},
		&common.ModuleHandles{Method: "post", Path: "admin/accept_list", Handles: common.HandleArray{m.CheckLogin, m.AcceptList}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.OpAcceptMinner}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.DelAcceptMinner}},

		&common.ModuleHandles{Method: "post", Path: "admin/agent_count", Handles: common.HandleArray{m.CheckLogin, m.AgentCount}},
		&common.ModuleHandles{Method: "post", Path: "admin/agent_list", Handles: common.HandleArray{m.CheckLogin, m.AgentList}},
		&common.ModuleHandles{Method: "post", Path: "admin/employer_list", Handles: common.HandleArray{m.CheckLogin, m.EmployerList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_msg", Handles: common.HandleArray{m.CheckLogin, m.Delmsg}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_agent", Handles: common.HandleArray{m.CheckLogin, m.DelAgent}},
		&common.ModuleHandles{Method: "post", Path: "admin/userexplode_controller", Handles: common.HandleArray{m.CheckLogin, m.UserControllerExp}},
		&common.ModuleHandles{Method: "post", Path: "admin/asset", Handles: common.HandleArray{m.CheckLogin, m.Assets}},
		&common.ModuleHandles{Method: "post", Path: "admin/save_parentmemo", Handles: common.HandleArray{m.CheckLogin, m.SaveParentMemo}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_admin", Handles: common.HandleArray{m.CheckLogin, m.DelAdmin}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_loan", Handles: common.HandleArray{m.CheckLogin, m.LoanSetting}},

		&common.ModuleHandles{Method: "post", Path: "admin/approve_recharge", Handles: common.HandleArray{m.CheckLogin, m.ApproveRecharge}},
		&common.ModuleHandles{Method: "post", Path: "admin/loan", Handles: common.HandleArray{m.CheckLogin, m.LoanList}},
		&common.ModuleHandles{Method: "post", Path: "admin/loanorder_list", Handles: common.HandleArray{m.CheckLogin, m.LoanOrderList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_loan", Handles: common.HandleArray{m.CheckLogin, m.DelLoan}},
		&common.ModuleHandles{Method: "post", Path: "admin/audit_loan", Handles: common.HandleArray{m.CheckLogin, m.AuditLoan}},
		&common.ModuleHandles{Method: "post", Path: "admin/collect_wallet", Handles: common.HandleArray{m.CheckLogin, m.CollectAddress}},
		&common.ModuleHandles{Method: "post", Path: "admin/collect_list", Handles: common.HandleArray{m.CheckLogin, m.CollectList}},
		&common.ModuleHandles{Method: "post", Path: "admin/save_uauth", Handles: common.HandleArray{m.CheckLogin, m.SaveAuth}},
		&common.ModuleHandles{Method: "post", Path: "admin/applycoin_list", Handles: common.HandleArray{m.CheckLogin, m.ApplyCoin}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_apply", Handles: common.HandleArray{m.CheckLogin, m.OpApply}},
		&common.ModuleHandles{Method: "post", Path: "admin/delapplycoin", Handles: common.HandleArray{m.CheckLogin, m.DelApplyCoin}},
		&common.ModuleHandles{Method: "post", Path: "admin/transfer_op", Handles: common.HandleArray{m.CheckLogin, m.TransferOp}},
		&common.ModuleHandles{Method: "post", Path: "admin/transfer_list", Handles: common.HandleArray{m.CheckLogin, m.TransferList}},
		&common.ModuleHandles{Method: "post", Path: "admin/msg_list", Handles: common.HandleArray{m.CheckLogin, m.MsgList}},
		&common.ModuleHandles{Method: "post", Path: "admin/send_user_notice", Handles: common.HandleArray{m.CheckLogin, m.SendUserNotice}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_usernotice", Handles: common.HandleArray{m.CheckLogin, m.DelUserNotice}},
		&common.ModuleHandles{Method: "post", Path: "admin/opuserAssetWallet", Handles: common.HandleArray{m.CheckLogin, m.OpuserAssetWallet}},
		&common.ModuleHandles{Method: "post", Path: "admin/audit_accept", Handles: common.HandleArray{m.CheckLogin, m.AuditAccept}},
		&common.ModuleHandles{Method: "post", Path: "admin/spot_delegate", Handles: common.HandleArray{m.CheckLogin, m.SpotDelegate}},
		&common.ModuleHandles{Method: "post", Path: "admin/opspot", Handles: common.HandleArray{m.CheckLogin, m.OpSpot}},
		&common.ModuleHandles{Method: "post", Path: "admin/kline_config", Handles: common.HandleArray{m.CheckLogin, m.kline_config}}, //K线控制
	}
}

func (m *AdminUserModule) OpSpot(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.OpSpot(rq))
}

func (m *AdminUserModule) SpotDelegate(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}

func (m *AdminUserModule) UserInfo(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminmodels.AdminResponse{
		State: adminmodels.SUCCESS,
		Data:  models.MODEL_USER.GetBaseInfo(uid),
	})
}

func (m *AdminUserModule) AuditAccept(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.AuditAccept(rq))
}
func (m *AdminUserModule) OpuserAssetWallet(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpuserAssetWallet(rq))
}
func (m *AdminUserModule) DelUserNotice(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelUserNotice(id))
}
func (m *AdminUserModule) SendUserNotice(r *gin.Context) {
	rq := new(adminmodels.UserNoticeMsg)
	m.ConvertObject(r, rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SendUserNotice(rq))
}
func (m *AdminUserModule) MsgList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MsgList(rq))
}
func (m *AdminUserModule) DelApplyCoin(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelApplyCoin(id))
}
func (m *AdminUserModule) TransferList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.TransferList(rq))
}
func (m *AdminUserModule) TransferOp(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.TransferOp(rq))
}

func (m *AdminUserModule) OpApply(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpApplyCoin(rq))
}

func (m *AdminUserModule) ApplyCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ApplyCoin(rq))
}

func (m *AdminUserModule) SaveAuth(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveAuth(rq))
}

func (m *AdminUserModule) CollectList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CollectLogList(rq))
}
func (m *AdminUserModule) CollectAddress(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CollectAddress(rq))
}
func (m *AdminUserModule) ApproveRecharge(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserApproveRecharge(rq))
}

func (m *AdminUserModule) AuditLoan(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpLuan(rq))

}

func (m *AdminUserModule) DelLoan(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelLoan(rq))
}

func (m *AdminUserModule) LoanList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.LoanList())
}
func (m *AdminUserModule) LoanSetting(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.LoanSetting(rq))
}
func (m *AdminUserModule) LoanOrderList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.LoanOrderList(rq))
}

func (m *AdminUserModule) SaveParentMemo(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveParentMemo(rq))
}
func (m *AdminUserModule) Assets(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserAssetConver(uid))
}

func (m *AdminUserModule) UserControllerExp(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserControllerExp(rq))
}

func (m *AdminUserModule) DelAgent(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_AGENT.DelAgent(id))
}

func (m *AdminUserModule) UserCoinLog(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserCoinLog(rq))
}

func (m *AdminUserModule) Delmsg(r *gin.Context) {
	id := m.GetValue(r, "sn_id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Delmsg(id))
}
func (m *AdminUserModule) EmployerList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_AGENT.EmpoyerList(rq))
}
func (m *AdminUserModule) AgentList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_AGENT.AgentList(rq))
}
func (m *AdminUserModule) AgentCount(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_AGENT.AgentCountList(rq))
}

func (m *AdminUserModule) ChangeCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpCredit(rq))
}

func (m *AdminUserModule) DelAcceptMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelMinneAccept(id))
}

func (m *AdminUserModule) OpAcceptMinner(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpMinnerAccept(rq))
}

func (m *AdminUserModule) AcceptList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.AcceptList(rq))
}
func (m *AdminUserModule) DelController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelController(rq))
}
func (m *AdminUserModule) ControllerList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ControllerTradeList())
}

func (m *AdminUserModule) ExplodeController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ExplodeController(rq))
}
func (m *AdminUserModule) KlineController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.KlineController(rq))
}

func (m *AdminUserModule) KickUser(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.KickUser(uid))
}

func (m *AdminUserModule) DelRule(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelRule(rq))
}

func (m *AdminUserModule) Rulelist(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Rulelist(rq))
}

func (m *AdminUserModule) RuleHandler(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.RuleHandler(rq))
}

func (m *AdminUserModule) SendMsg(r *gin.Context) {
	rq := new(adminmodels.CustomMsg)
	m.ConvertObject(r, rq)
	fmt.Println(rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SendMsg(rq))
}

func (m *AdminUserModule) UserByMessage(r *gin.Context) {

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserByMessage(m.GetInt(r, "uid")))
}

func (m *AdminUserModule) Chatlist(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.CustomServiceList(rq))
}

func (m *AdminUserModule) StatictisCount(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.StatictisCount())
}

func (m *AdminUserModule) ClearNotify(r *gin.Context) {
	tp := m.GetValue(r, "type")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ClearNotify(tp))
}

func (m *AdminUserModule) NotifyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.NotifyList())
}

func (m *AdminUserModule) OpUser(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpUser(rq))
}

func (m *AdminUserModule) UserAuth2List(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserAuthDouble(rq))
}

func (m *AdminUserModule) SiteCount(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: 2000, Data: adminmodels.SYSTEM_MODEL.SiteCount(rq)})
}

func (m *AdminUserModule) Setting(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Setting(rq))
}
func (m *AdminUserModule) UserLevelCount(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserTeamLevelCount(rq))
}
func (m *AdminUserModule) DelegateHistoryDel(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.DelegateHistoryDel(id))
}

func (m *AdminUserModule) DelegateHistory(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}
func (m *AdminUserModule) DelegateNow(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 0
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}

func (m *AdminUserModule) ManualDelegateTrade(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.ManualOperationTrade(uid, sn))
}

func (m *AdminUserModule) HoldTradeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.TradeList(rq, true, ""))
}

func (m *AdminUserModule) CloseTrade(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryCloseRradeList(rq, true, ""))
}
func (m *AdminUserModule) StopMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	pass := m.GetValue(r, "pass")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.StopMinner(id, pass))
}

func (m *AdminUserModule) MorderList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MorderList(rq))
}

func (m *AdminUserModule) Uauth2Op(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UauthOp(rq, 2))
}

func (m *AdminUserModule) Uauth2Del(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: 2000, Data: adminmodels.MODEL_USER.UauthDel(id, 2)})
}
func (m *AdminUserModule) UauthOp(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UauthOp(rq, 1))
}
func (m *AdminUserModule) UauthDel(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UauthDel(id, 1))
}
func (m *AdminUserModule) UserAuthList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserAuthList(rq))
}

func (m *AdminUserModule) WithdrawList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.WithdrawList(rq))
}

func (m *AdminUserModule) WithdrawOp(r *gin.Context) {
	id := m.GetInt(r, "id")
	state := m.GetInt(r, "state")
	info := m.GetValue(r, "info")
	password := m.GetValue(r, "password")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpWithdraw(id, state, info, password))
}

func (m *AdminUserModule) SaveWithdraw(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveWithdraw(rq))
}
func (m *AdminUserModule) UserWallet(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserWallet(rq))
}

func (m *AdminUserModule) DeluserWallet(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeluserWallet(id))
}

func (m *AdminUserModule) OpenAddr(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpenAddr(rq))
}
func (m *AdminUserModule) DelAddr(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelWalletAddress(id))
}
func (m *AdminUserModule) AddrList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.WalletAddressList(rq))
}

func (m *AdminUserModule) OpAddr(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpWalletAddress(rq))
}

func (m *AdminUserModule) OpRecharge(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpRecharge(rq))
}

func (m *AdminUserModule) RechargeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.RechargeList(rq))
}
func (m *AdminUserModule) NoticeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.NoticeList(rq))
}

func (m *AdminUserModule) OPNotice(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpNotice(rq))
}

func (m *AdminUserModule) DelNotice(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelNotice(id))
}

func (m *AdminUserModule) DelMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelMinner(id))
}

func (m *AdminUserModule) OpMinner(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpMinner(rq))
}
func (m *AdminUserModule) MinnerOpen(r *gin.Context) {
	id := m.GetInt(r, "id")
	key := m.GetValue(r, "key")
	isopen := m.GetInt(r, "state")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.MinnerSet(id, key, isopen))
}
func (m *AdminUserModule) MinnerList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.MinnerList(rq))
}

func (m *AdminUserModule) DelExplode(r *gin.Context) {

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelExplodeTrade(m.GetInt(r, "id")))
}
func (m *AdminUserModule) OpExplodeTrade(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpExplodeTrade(rq))
}
func (m *AdminUserModule) ExplodeList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ExplodeTradeList())
}

func (m *AdminUserModule) DelCurrency(r *gin.Context) {
	id := m.GetValue(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCurrency(id))
}

func (m *AdminUserModule) OpCurrency(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCurrency(rq))

}

func (m *AdminUserModule) CurrencyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CurrencyList())
}
func (m *AdminUserModule) DelCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCoin(rq))
}

// 保存币
func (m *AdminUserModule) SaveCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCoin(rq))
}

func (m *AdminUserModule) CoinDescList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CoinDescList(rq))
}

func (m *AdminUserModule) OpCoinDesc(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCoinDesc(rq))
}

func (m *AdminUserModule) DelCoinDesc(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCoinDesc(id))
}

func (m *AdminUserModule) CoinList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CoinList(rq))
}
func (m *AdminUserModule) DelMean(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelMean(rq))
}

func (m *AdminUserModule) HandlerMean(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.HandlerMean(rq))
}

/**
 *	修改用户组
 */
func (m *AdminUserModule) HandlerRole(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.HandlerRole(rq))
}
func (m *AdminUserModule) AddAdmin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AddAdmin(rq))
}

func (m *AdminUserModule) MeanRouter(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MeanRouter())
}
func (m *AdminUserModule) RoleList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.RoleList())
}

func (m *AdminUserModule) AuthRouter(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AuthRouter())
}
func (m *AdminUserModule) AdminList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AdminList())
}

func (m *AdminUserModule) DelAdmin(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelAdmin(id))
}
func (m *AdminUserModule) Login(r *gin.Context) {
	var request adminmodels.LoginRequest
	m.ConvertObject(r, &request)
	request.ClientIp = r.ClientIP()
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.Login(request))
}

func (m *AdminUserModule) Logout(r *gin.Context) {
	sid := r.GetString("sid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: adminmodels.SUCCESS, Data: adminmodels.MODEL_USER.Logout(sid)})
}

func (m *AdminUserModule) UserAsset(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserAssetList(rq))
}
func (m *AdminUserModule) TokenSid(r *gin.Context) {
	fmt.Println(r.GetString("sid"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.TokenInfo(r.GetString("sid")))
}
func (m *AdminUserModule) UserList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserList(rq))
}

func (m *AdminUserModule) CheckLogin(r *gin.Context) {
	sid := r.GetString("sid")
	u := adminmodels.MODEL_USER.SidInfo(sid)
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	if u == nil {
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminmodels.AdminResponse{
			State: 50008,
			Data:  "当前用户的sid无法获取到值",
		})
		r.Abort()
	}
	notify, _ := regexp.Match(`admin\/notify_list`, []byte(r.Request.URL.Path))
	if notify {
		return
	}
	config.GlobalDB.InsertData("operation_log", db.DB_PARAMS{
		"username":   u.UserName,
		"path":       r.Request.URL.Path,
		"createtime": utils.GetNow(),
		"data":       rq,
		"ip":         r.ClientIP(),
	})

}
func (m *AdminUserModule) kline_config(r *gin.Context) {
	var klineConfig models.CoinKlineConfig
	id := m.GetInt(r, "id")
	klineConfig.BaseAmount = m.GetInt(r, "base_amount")
	klineConfig.Heart = m.GetFloat(r, "heart")
	klineConfig.HighRate = m.GetFloat(r, "high_rate")
	klineConfig.LowRate = m.GetFloat(r, "low_rate")
	klineConfig.MaxPrice = m.GetFloat(r, "max_price")
	klineConfig.MinPrice = m.GetFloat(r, "min_price")
	klineConfig.UpRate = m.GetInt(r, "up_rate")
	config.GlobalDB.UpdateData(models.DB_TABLE_COINS, db.DB_PARAMS{"kline_config": klineConfig}, db.DB_PARAMS{"id": id})
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminmodels.AdminResponse{
		State: adminmodels.SUCCESS,
		Data:  "已提交控制!",
	})

}
