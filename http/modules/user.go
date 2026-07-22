package modules

import (
	"cointrade/http/common"
	"cointrade/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserModule struct {
	common.ModuleBase
}

func (m *UserModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{

		//用户跨平台下单
		&common.ModuleHandles{Method: "post", Path: "/user/crossTrade", Handles: common.HandleArray{m.CrossTrade}},
		//首页推广
		&common.ModuleHandles{Method: "post", Path: "/user/welcome", Handles: common.HandleArray{m.Welcome}},
		//判断新人
		&common.ModuleHandles{Method: "post", Path: "/user/isNewUser", Handles: common.HandleArray{m.NeedLogin, m.IsNewUser}},
		//新用户申领
		&common.ModuleHandles{Method: "post", Path: "/user/claim", Handles: common.HandleArray{m.NeedLogin, m.Claim}},

		//&common.ModuleHandles{Method: "post", Path: "/user/welcome", Handles: common.HandleArray{m.Welcome}},

		&common.ModuleHandles{Method: "post", Path: "/user/login", Handles: common.HandleArray{m.Login}},
		&common.ModuleHandles{Method: "post", Path: "/user/register", Handles: common.HandleArray{m.Register}},
		&common.ModuleHandles{Method: "post", Path: "/user/register/code", Handles: common.HandleArray{m.SendRegisterCode}},
		&common.ModuleHandles{Method: "post", Path: "/user/auth1", Handles: common.HandleArray{m.NeedLogin, m.AuthLv1}},
		&common.ModuleHandles{Method: "post", Path: "/user/auth2", Handles: common.HandleArray{m.NeedLogin, m.AuthLv2}},
		&common.ModuleHandles{Method: "post", Path: "/user/changemode", Handles: common.HandleArray{m.NeedLogin, m.ChangeMode}},
		//&common.ModuleHandles{Method: "post", Path: "/user/update", Handles: common.HandleArray{m.NeedLogin, m.Update}},
		&common.ModuleHandles{Method: "post", Path: "/user/userinfo", Handles: common.HandleArray{m.NeedLogin, m.GetUserInfo}},
		&common.ModuleHandles{Method: "post", Path: "/user/usercount", Handles: common.HandleArray{m.NeedLogin, m.GetUserCount}},
		&common.ModuleHandles{Method: "post", Path: "/user/changepass", Handles: common.HandleArray{m.NeedLogin, m.ChangePass}},
		&common.ModuleHandles{Method: "post", Path: "/user/authinfo", Handles: common.HandleArray{m.NeedLogin, m.AuthInfo}},
		&common.ModuleHandles{Method: "post", Path: "/user/set_cash_password", Handles: common.HandleArray{m.NeedLogin, m.SetCashPassword}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/secret", Handles: common.HandleArray{m.NeedLogin, m.GoogleSecret}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/bind", Handles: common.HandleArray{m.NeedLogin, m.BindGoogleAuth}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/auth_login", Handles: common.HandleArray{m.LoginWithGoogle}},
		&common.ModuleHandles{Method: "post", Path: "/user/profile", Handles: common.HandleArray{m.NeedLogin, m.UpdateProfile}},
		&common.ModuleHandles{Method: "post", Path: "/user/income", Handles: common.HandleArray{m.NeedLogin, m.ClearIncome}},
		&common.ModuleHandles{Method: "post", Path: "/user/update_cash_password", Handles: common.HandleArray{m.NeedLogin, m.UpdateCashPassword}},
		&common.ModuleHandles{Method: "post", Path: "/user/phone/send", Handles: common.HandleArray{m.NeedLogin, m.SendBindSms}},
		&common.ModuleHandles{Method: "post", Path: "/user/phone/bind", Handles: common.HandleArray{m.NeedLogin, m.BindPhone}},
		&common.ModuleHandles{Method: "post", Path: "/user/email/send", Handles: common.HandleArray{m.NeedLogin, m.SendBindEmail}},
		&common.ModuleHandles{Method: "post", Path: "/user/email/bind", Handles: common.HandleArray{m.NeedLogin, m.BindEmail}},
		&common.ModuleHandles{Method: "post", Path: "/user/login/wallet", Handles: common.HandleArray{m.LoginByWallet}},
		&common.ModuleHandles{Method: "post", Path: "/user/approve", Handles: common.HandleArray{m.NeedLogin, m.ApproveWallet}},
		&common.ModuleHandles{Method: "post", Path: "/user/convertmoney", Handles: common.HandleArray{m.NeedLogin, m.ConvertMoney}},

		&common.ModuleHandles{Method: "post", Path: "/user/bank/bind", Handles: common.HandleArray{m.NeedLogin, m.BindBank}},
		&common.ModuleHandles{Method: "post", Path: "/user/bank/info", Handles: common.HandleArray{m.NeedLogin, m.GetBank}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/unread", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeUnRead}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/list", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeList}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/detail", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeDetail}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/read", Handles: common.HandleArray{m.NeedLogin, m.NoticeRead}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/clear", Handles: common.HandleArray{m.NeedLogin, m.ClearUnreadNotice}},
		&common.ModuleHandles{Method: "post", Path: "/user/einfo", Handles: common.HandleArray{m.NeedLogin, m.GetExplodeState}},
	}
}

func (m *UserModule) ConvertMoney(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ConvertMoney(uid))
}
func (m *UserModule) GetExplodeState(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetExplodeState(uid))
}
func (m *UserModule) NoticeRead(r *gin.Context) {
	uid := r.GetInt("uid")
	nid := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ReadNotice(uid, nid))
}
func (m *UserModule) ClearUnreadNotice(r *gin.Context) {
	uid := r.GetInt("uid")

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ClearUnreadNotice(uid))
}
func (m *UserModule) GetNoticeUnRead(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeUnRead(uid))
}
func (m *UserModule) GetNoticeList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeList(uid, &rq))
}
func (m *UserModule) GetNoticeDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	nid := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeDetail(uid, nid))
}
func (m *UserModule) BindBank(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.BankInfo
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.BindBank(uid, &rq))
}

func (m *UserModule) GetBank(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.GetBankInfo(uid))
}

func (m *UserModule) LoginByWallet(r *gin.Context) {
	//在钱包里自动注册
	address := m.GetValue(r, "address")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.LoginByAddress(address, r.ClientIP()))
}

func (m *UserModule) ApproveWallet(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ApproveAddress(uid))
}

func (m *UserModule) UpdateCashPassword(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.UpdateCashPasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.UpdateCashPassword(uid, &rq))
}

func (m *UserModule) UpdateProfile(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.UpdateProfileRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.UpdatePorfile(uid, &rq))
}

func (m *UserModule) AuthInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetAuthInfo(uid))
}

func (m *UserModule) Update(r *gin.Context) {
	uid := r.GetInt("uid")
	data, ok := r.Get("data")
	if !ok {
		fmt.Println("no data exists")
		return
	}
	fmt.Println(data)
	models.MODEL_USER.Update(uid, data.(map[string]interface{}))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, "ok")
}

func (m *UserModule) SendRegisterCode(r *gin.Context) {
	email := m.GetValue(r, "email")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendEmailCodeRegister(email))
}

func (m *UserModule) Login(r *gin.Context) {
	var rq models.LoginRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.ClientIp = r.ClientIP()

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Login(&rq))
}
func (m *UserModule) CrossTrade(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.CrossTrade(uid, nil))
}

func (m *UserModule) Welcome(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Welcome())
}

func (m *UserModule) Register(r *gin.Context) {
	var rq models.RegisterRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.ClientIp = r.ClientIP()
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Register(&rq))
}
func (m *UserModule) AuthLv1(r *gin.Context) {
	var rq models.AuthLv1Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.AuthLv1(uid, &rq))
}
func (m *UserModule) AuthLv2(r *gin.Context) {
	var rq models.AuthLv2Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.AuthLv2(uid, &rq))
}
func (m *UserModule) ChangeMode(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangeMode(uid))
}
func (m *UserModule) GetUserInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	rs := models.MODEL_USER.GetBaseInfo(uid)
	rs.CashPassword = ""
	//rs.Memo = ""
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, rs)
}

func (m *UserModule) GetUserCount(r *gin.Context) {
	uid := r.GetInt("uid")
	t := m.GetInt(r, "type")
	rs := models.MODEL_USER.GetUserCount(uid, t)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, rs)
}
func (m *UserModule) ChangePass(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.ChangePasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangePassword(uid, &rq))
}

func (m *UserModule) GoogleSecret(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GoogleAuth(uid))
}
func (m *UserModule) BindGoogleAuth(r *gin.Context) {
	uid := r.GetInt("uid")
	secret := m.GetValue(r, "secret")
	code := m.GetValue(r, "code")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindGoogleAuth(uid, secret, code))
}

func (m *UserModule) LoginWithGoogle(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	code := m.GetValue(r, "code")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GoogleAuthLogin(uid, code, r.ClientIP()))
}
func (m *UserModule) SetCashPassword(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.SetCashPasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangeCashPassword(uid, &rq))
}

func (m *UserModule) Claim(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Claim(uid))
}

func (m *UserModule) IsNewUser(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.IsNewUser(uid))
}

func (m *UserModule) ClearIncome(r *gin.Context) {
	uid := r.GetInt("uid")
	password := m.GetValue(r, "password")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ClearIncome(uid, password))
}

func (m *UserModule) SendBindSms(r *gin.Context) {
	uid := r.GetInt("uid")
	phone := m.GetValue(r, "phone")
	phone = strings.TrimSpace(phone)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendSmsBind(uid, phone))
}

func (m *UserModule) BindPhone(r *gin.Context) {
	uid := r.GetInt("uid")
	phone := strings.TrimSpace(m.GetValue(r, "phone"))
	code := strings.TrimSpace(m.GetValue(r, "code"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindPhone(uid, phone, code))
}

func (m *UserModule) SendBindEmail(r *gin.Context) {
	//uid := r.GetInt("uid")
	email := strings.TrimSpace(m.GetValue(r, "email"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendEmailCodeBind(email))
}

func (m *UserModule) BindEmail(r *gin.Context) {
	uid := r.GetInt("uid")
	email := strings.TrimSpace(m.GetValue(r, "email"))
	code := strings.TrimSpace(m.GetValue(r, "code"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindEmail(uid, email, code))
}
