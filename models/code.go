package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CodeModel struct{}

const (
	CODE_TYPE_EMAIL      = 1
	CODE_TYPE_SMS        = 2
	EMAIL_EXPIRE_TIME    = 5 * 60 //邮件验证码有效期
	SMS_EXPIRE_TIME      = 3 * 60 //手机验证码有效期
	EMAIL_CACHE_REGISTER = 1      //注册邮件
	EMAIL_CACHE_VERDIFY  = 2      //验证邮件
	SMS_CACHE_BIND       = 1      //绑定手机
	SMS_CACHE_VERDIFY    = 2      //验证码
)

func (m *CodeModel) SendEmailCodeRegister(email string) bool {
	//发送注册时的邮件验证码
	if !utils.CheckEmail(email) {
		return false
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"email": email}, db.DB_FIELDS{"id"})
	if one != nil {
		return false
	}
	code := m.CreateCode()
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_REGISTER, 0, email)
	rs := config.GlobalRedis.Get(cacheid)
	if rs != "" {
		return false
	}
	config.GlobalRedis.Set(cacheid, code, EMAIL_EXPIRE_TIME*time.Second)
	go func() {
		for i := 0; i < 3; i++ {
			if m.SendEmail(email, config.SYSTEM_CONFIG.SiteName+" Register Code", fmt.Sprintf(config.SYSTEM_CONFIG.SiteName+" Register code is [%s]", code)) {
				break
			}
		}
	}()
	return true
}
func (m *CodeModel) SendEmailCodeBind(email string) bool {
	//发送注册时的邮件验证码
	if !utils.CheckEmail(email) {
		return false
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"email": email}, db.DB_FIELDS{"id"})
	if one != nil {
		return false
	}
	code := m.CreateCode()
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_VERDIFY, 0, email)
	rs := config.GlobalRedis.Get(cacheid)
	if rs != "" {
		return false
	}
	config.GlobalRedis.Set(cacheid, code, EMAIL_EXPIRE_TIME*time.Second)
	go func() {
		for i := 0; i < 3; i++ {
			if m.SendEmail(email, config.SYSTEM_CONFIG.SiteName+" Verdify Code", fmt.Sprintf(config.SYSTEM_CONFIG.SiteName+" Verdify code is [%s]", code)) {
				break
			}
		}
	}()
	return true
}
func (m *CodeModel) GetEmailCodeBind(email string) string {
	//获得注册时的验证码
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_VERDIFY, 0, email)

	rs := config.GlobalRedis.Get(cacheid)
	//config.GlobalRedis.Delete(cacheid)

	return rs
}
func (m *CodeModel) SendUserEmailCode(uid int, t int) bool {
	//发送用户邮件验证码
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil || uinfo.Email == "" {
		return false
	}
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_VERDIFY, uid, uinfo.Email)
	rs := config.GlobalRedis.Get(cacheid)
	if rs != "" {
		return false
	}
	code := m.CreateCode()
	config.GlobalRedis.Set(cacheid, code, EMAIL_EXPIRE_TIME*time.Second)
	go func() {
		for i := 0; i < 3; i++ {
			if m.SendEmail(uinfo.Email, config.SYSTEM_CONFIG.SiteName+" Verdify Code", fmt.Sprintf(config.SYSTEM_CONFIG.SiteName+" Verdify code is [%s]", code)) {
				break
			}
		}
	}()
	return true
}

func (m *CodeModel) SendSmsBind(uid int, phone string) bool {
	//发送绑定手机时的验证码
	/*uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return false
	}*/
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"phone": phone}, db.DB_FIELDS{"id"})
	if one != nil { //已经被绑定过了
		return false
	}
	cacheid := m.GetCacheName(CODE_TYPE_SMS, SMS_CACHE_BIND, uid, phone)
	rs := config.GlobalRedis.Get(cacheid)
	if rs != "" {
		return false
	}
	code := m.CreateCode()
	go func() {
		for i := 0; i < 3; i++ {
			if m.SendSms(phone, fmt.Sprintf(config.SYSTEM_CONFIG.SiteName+" Bind code is [%s]", code)) {
				break
			}
		}
	}()
	config.GlobalRedis.Set(cacheid, code, SMS_EXPIRE_TIME*time.Second)

	return true
}

func (m *CodeModel) SendUserSmsCode(uid int, t int) bool {
	//发送用户手机验证码
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil || uinfo.Phone == "" {
		return false
	}
	code := m.CreateCode()

	cacheid := m.GetCacheName(CODE_TYPE_SMS, SMS_CACHE_VERDIFY, uid, uinfo.Phone)
	rs := config.GlobalRedis.Get(cacheid)
	if rs != "" {
		return false
	}
	go func() {
		for i := 0; i < 3; i++ {
			if m.SendSms(uinfo.Phone, fmt.Sprintf(config.SYSTEM_CONFIG.SiteName+" Verdify code is [%s]", code)) {
				break
			}
		}
	}()
	config.GlobalRedis.Set(cacheid, code, SMS_EXPIRE_TIME*time.Second)
	return true
}

func (m *CodeModel) GetBindSmsCode(uid int, phone string) string {
	//获得绑定手机时的验证码
	cacheid := m.GetCacheName(CODE_TYPE_SMS, SMS_CACHE_BIND, uid, phone)

	rs := config.GlobalRedis.Get(cacheid)
	//config.GlobalRedis.Delete(cacheid)

	return rs
}

func (m *CodeModel) GetUserSmsCode(uid int, t int) string {
	//获得用户手机验证码
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil || uinfo.Phone == "" {
		return ""
	}
	cacheid := m.GetCacheName(CODE_TYPE_SMS, SMS_CACHE_VERDIFY, uid, uinfo.Phone)

	rs := config.GlobalRedis.Get(cacheid)
	//config.GlobalRedis.Delete(cacheid)

	return rs
}

func (m *CodeModel) GetEmailCodeRegister(email string) string {
	//获得注册时的验证码
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_REGISTER, 0, email)

	rs := config.GlobalRedis.Get(cacheid)
	//config.GlobalRedis.Delete(cacheid)

	return rs
}

func (m *CodeModel) GetEmailCode(uid int, t int) string {
	//获得用户邮件验证码
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil || uinfo.Email == "" {
		return ""
	}
	cacheid := m.GetCacheName(CODE_TYPE_EMAIL, EMAIL_CACHE_VERDIFY, uid, uinfo.Email)
	rs := config.GlobalRedis.Get(cacheid)
	//config.GlobalRedis.Delete(cacheid)

	return rs
}
func (m *CodeModel) CreateCode() string {
	n := 100000 + rand.Intn(900000)
	return strconv.Itoa(n)
}

func (m *CodeModel) SendEmail(email, title, content string) bool {
	err := utils.SendToMail(
		config.GlobalConfig.GetValue("smtp_user").ToString(),
		config.GlobalConfig.GetValue("smtp_pass").ToString(),
		config.GlobalConfig.GetValue("smtp_host").ToString(),
		email,
		title,
		content,
		"html",
	)
	return err == nil
}

func (m *CodeModel) SendSms(phone, content string) bool {
	return config.GlobalSMS.SendSmsA(phone, content)
}

func (m *CodeModel) GetCacheName(code_t int, type_t int, uid int, to string) string {
	prefix := "code"
	switch code_t {
	case CODE_TYPE_EMAIL:
		prefix = prefix + "_email"
	case CODE_TYPE_SMS:
		prefix = prefix + "_sms"
	}
	return fmt.Sprintf(prefix+"_%d_%d_%s", type_t, uid, to)
}
