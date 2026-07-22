package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/http/common"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) financeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/recharge_list", Handles: common.HandleArray{m.CheckLogin, m.RechargeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/recharge_op", Handles: common.HandleArray{m.CheckLogin, m.ReviewRecharge}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_list", Handles: common.HandleArray{m.CheckLogin, m.AddrList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_op", Handles: common.HandleArray{m.CheckLogin, m.OpAddr}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_del", Handles: common.HandleArray{m.CheckLogin, m.DelAddr}},
		&common.ModuleHandles{Method: "post", Path: "/admin/addr_openclose", Handles: common.HandleArray{m.OpenAddr}},
		&common.ModuleHandles{Method: "post", Path: "/admin/user_wallet", Handles: common.HandleArray{m.CheckLogin, m.ListUserWallets}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_userwallet", Handles: common.HandleArray{m.CheckLogin, m.DeleteUserWallet}},
		&common.ModuleHandles{Method: "post", Path: "/admin/withdraw_list", Handles: common.HandleArray{m.CheckLogin, m.WithdrawList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/withdraw_op", Handles: common.HandleArray{m.CheckLogin, m.ReviewWithdraw}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_list", Handles: common.HandleArray{m.CheckLogin, m.UserAuthList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_del", Handles: common.HandleArray{m.CheckLogin, m.DeletePrimaryUserAuth}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth_op", Handles: common.HandleArray{m.CheckLogin, m.ReviewPrimaryUserAuth}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_list", Handles: common.HandleArray{m.ListAdvancedUserAuth}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_op", Handles: common.HandleArray{m.CheckLogin, m.ReviewAdvancedUserAuth}},
		&common.ModuleHandles{Method: "post", Path: "/admin/save_withdraw", Handles: common.HandleArray{m.CheckLogin, m.SaveWithdraw}},
		&common.ModuleHandles{Method: "post", Path: "/admin/uauth2_del", Handles: common.HandleArray{m.CheckLogin, m.DeleteAdvancedUserAuth}},
		&common.ModuleHandles{Method: "post", Path: "admin/approve_recharge", Handles: common.HandleArray{m.CheckLogin, m.ListRechargeApprovals}},
		&common.ModuleHandles{Method: "post", Path: "admin/loanorder_list", Handles: common.HandleArray{m.CheckLogin, m.ListLoanOrders}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_loan", Handles: common.HandleArray{m.CheckLogin, m.DeleteLoanOrder}},
		&common.ModuleHandles{Method: "post", Path: "admin/audit_loan", Handles: common.HandleArray{m.CheckLogin, m.ReviewLoanOrder}},
		&common.ModuleHandles{Method: "post", Path: "admin/collect_wallet", Handles: common.HandleArray{m.CheckLogin, m.CollectAddress}},
		&common.ModuleHandles{Method: "post", Path: "admin/collect_list", Handles: common.HandleArray{m.CheckLogin, m.CollectList}},
		&common.ModuleHandles{Method: "post", Path: "admin/save_uauth", Handles: common.HandleArray{m.CheckLogin, m.SaveAuth}},
		&common.ModuleHandles{Method: "post", Path: "admin/applycoin_list", Handles: common.HandleArray{m.CheckLogin, m.ListCoinApplications}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_apply", Handles: common.HandleArray{m.CheckLogin, m.ReviewCoinApplication}},
		&common.ModuleHandles{Method: "post", Path: "admin/delapplycoin", Handles: common.HandleArray{m.CheckLogin, m.DeleteCoinApplication}},
		&common.ModuleHandles{Method: "post", Path: "admin/transfer_op", Handles: common.HandleArray{m.CheckLogin, m.ReviewTransfer}},
		&common.ModuleHandles{Method: "post", Path: "admin/transfer_list", Handles: common.HandleArray{m.CheckLogin, m.TransferList}},
		&common.ModuleHandles{Method: "post", Path: "admin/opuserAssetWallet", Handles: common.HandleArray{m.CheckLogin, m.UpdateUserAssetWallet}},
	}
}

func (m *AdminUserModule) UpdateUserAssetWallet(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UpdateUserAssetWallet(rq))
}

func (m *AdminUserModule) DeleteCoinApplication(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteCoinApplication(id))
}

func (m *AdminUserModule) TransferList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.TransferList(rq))
}

func (m *AdminUserModule) ReviewTransfer(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewTransfer(rq))
}

func (m *AdminUserModule) ReviewCoinApplication(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewCoinApplication(rq))
}

func (m *AdminUserModule) ListCoinApplications(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListCoinApplications(rq))
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

func (m *AdminUserModule) ListRechargeApprovals(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListRechargeApprovals(rq))
}

func (m *AdminUserModule) ReviewLoanOrder(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewLoanOrder(rq))
}

func (m *AdminUserModule) DeleteLoanOrder(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteLoanOrder(rq))
}

func (m *AdminUserModule) ListLoanOrders(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListLoanOrders(rq))
}

func (m *AdminUserModule) ListAdvancedUserAuth(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListAdvancedUserAuth(rq))
}

func (m *AdminUserModule) ReviewAdvancedUserAuth(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewUserAuth(rq, 2))
}

func (m *AdminUserModule) DeleteAdvancedUserAuth(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: 2000, Data: adminmodels.MODEL_USER.DeleteUserAuth(id, 2)})
}

func (m *AdminUserModule) ReviewPrimaryUserAuth(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewUserAuth(rq, 1))
}

func (m *AdminUserModule) DeletePrimaryUserAuth(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteUserAuth(id, 1))
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

func (m *AdminUserModule) ReviewWithdraw(r *gin.Context) {
	id := m.GetInt(r, "id")
	state := m.GetInt(r, "state")
	info := m.GetValue(r, "info")
	password := m.GetValue(r, "password")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewWithdraw(id, state, info, password))
}

func (m *AdminUserModule) SaveWithdraw(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveWithdraw(rq))
}

func (m *AdminUserModule) ListUserWallets(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListUserWallets(rq))
}

func (m *AdminUserModule) DeleteUserWallet(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteUserWallet(id))
}

func (m *AdminUserModule) OpenAddr(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpenAddr(rq))
}

func (m *AdminUserModule) DelAddr(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DeleteWalletAddress(id))
}

func (m *AdminUserModule) AddrList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.WalletAddressList(rq))
}

func (m *AdminUserModule) OpAddr(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.SaveWalletAddress(rq))
}

func (m *AdminUserModule) ReviewRecharge(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ReviewRecharge(rq))
}

func (m *AdminUserModule) RechargeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.RechargeList(rq))
}
