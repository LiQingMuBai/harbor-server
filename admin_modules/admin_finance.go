package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/http/common"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) financeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
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
		&common.ModuleHandles{Method: "post", Path: "admin/approve_recharge", Handles: common.HandleArray{m.CheckLogin, m.ApproveRecharge}},
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
		&common.ModuleHandles{Method: "post", Path: "admin/opuserAssetWallet", Handles: common.HandleArray{m.CheckLogin, m.OpuserAssetWallet}},
	}
}

func (m *AdminUserModule) OpuserAssetWallet(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.OpuserAssetWallet(rq))
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

func (m *AdminUserModule) LoanOrderList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.LoanOrderList(rq))
}

func (m *AdminUserModule) UserAuth2List(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserAuthDouble(rq))
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
