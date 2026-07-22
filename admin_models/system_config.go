package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"html"
	"math"
	"strings"
	"time"
)

func (s *SystemModel) NotifyList() *AdminResponse {
	rs := notify.NOTIFY.NotifyList()
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"notify": rs,
			"unread": s.Unread(),
			"color":  &map[int]string{1: "primary", 2: "success", 3: "info", 4: "danger"},
		},
	}
}

func (s *SystemModel) ClearNotify(tp string) *AdminResponse {
	notify.NOTIFY.ClearNotify(tp)
	return &AdminResponse{State: SUCCESS, Data: ""}
}

func (s *SystemModel) Unread() *P {
	recharge := config.GlobalDB.GetCount(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"state": 0, "direction": 1})
	withdraw := config.GlobalDB.GetCount(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"state": 0, "direction": 2})
	auth := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH, db.DB_PARAMS{"process_state": 0})
	advAuth := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"state": 0})
	return &P{
		"78": withdraw,
		"75": recharge,
		"74": withdraw + recharge,
		"91": advAuth,
		"79": auth,
		"90": advAuth + auth,
	}
}

func (s *SystemModel) NoticePos() *P {
	return &P{
		"index":         "首页最新公告",
		"list":          "公告列表",
		"banner":        "首页banner",
		"mining_banner": "矿机banner",
	}
}

func (s *SystemModel) ModePair() *P {
	return &P{"1": "真实交易", "2": "虚拟交易"}
}

func (s *SystemModel) UserTypePair() *P {
	return &P{"1": "真实用户", "2": "内部用户"}
}

func (s *SystemModel) OnlinePair() *P {
	return &P{"1": "在线", "0": "离线"}
}

func (s *SystemModel) TradePair() *P {
	return &P{"1": "永续合约", "2": "交割合约", "3": "币币交易"}
}

func (s *SystemModel) DirectPair() *P {
	return &P{"1": "买涨", "2": "买跌"}
}

func (s *SystemModel) TradeStatePair() *P {
	return &P{"1": "持仓中", "2": "已结算"}
}

func (s *SystemModel) DelegateType() *P {
	return &P{"1": "买入", "2": "卖出"}
}

func (s *SystemModel) ApproveState() *P {
	return &P{"0": "等待到账", "1": "已到账", "2": "失败"}
}

func (s *SystemModel) DelegateStatePair() *P {
	return &P{"1": "委托中", "2": "已成交"}
}

func (s *SystemModel) UserStatus() *P {
	return &P{"1": "正常", "0": "禁用"}
}

func (s *SystemModel) COINLOG_TYPELIST() *map[int]string {
	return &map[int]string{
		models.COIN_LOG_USER_RECHARGE:         "用户充值",
		models.COIN_LOG_USER_WITHDRAW:         "用户提现",
		models.COIN_LOG_USER_PROFIT:           "用户矿机收益",
		models.COIN_LOG_USER_CLOSE:            "用户平仓",
		models.COIN_LOG_USER_DELEGATE:         "用户委托",
		models.COIN_LOG_USER_DELEGATE_SUCCESS: "用户委托成功",
		models.COIN_LOG_USER_CANCLE:           "用户撤单",
		models.COIN_LOG_USER_BUY_MINING:       "用户购买矿机",
		models.COIN_LOG_USER_MINING_PROFIT:    "用户挖矿获利",
		models.COIN_LOG_USER_CLEAR_INCOME:     "用户提取下级返利收入",
		models.COIN_LOG_USER_WITHDRAW_FAILD:   "用户提现失败",
		models.COIN_LOG_USER_EXCHANGE:         "用户兑换",
		models.COIN_LOG_BB_TRADE:              "币币交易",
		models.COIN_LOG_EXPLODE_TRADE:         "交割交易",
		models.COIN_LOG_KEEP_TRADE:            "永续交易",
		models.COIN_LOG_BACKEND:               "资产增加",
		models.COIN_LOG_LORA_IN:               "贷款通过",
	}
}

func (s *SystemModel) RuleType() *P {
	return &P{
		"recharge": "充值规则",
		"withdraw": "提现规则",
		"referer":  "推广规则",
		"mining":   "矿机规则",
		"about":    "关于我们",
		"help":     "隐私政策",
		"legal":    "合规性",
		"company":  "公司信息",
	}
}

func (s *SystemModel) LoanState() *P {
	return &P{"0": "申请中", "1": "已放币", "2": "已拒绝"}
}

func (s *SystemModel) WithdrawStatus() *P {
	return &P{"1": "正常", "0": "禁止"}
}

func (s *SystemModel) UserAuthLevel() map[int]string {
	return map[int]string{0: "未认证", 1: "初级认证", 2: "高级认证"}
}

func (s *SystemModel) UserStatePair() P {
	return P{"0": "未处理", "1": "已通过", "2": "已驳回"}
}

func (s *SystemModel) AuthProccess() *P {
	return &P{"0": "待审 ", "1": "成功", "2": "失败"}
}

func (s *SystemModel) WithdrawPair() *P {
	return &P{"0": "未处理", "1": "通过", "2": "驳回"}
}

func (s *SystemModel) RelationAuth() *P {
	return &P{"1": "同学", "2": "亲人", "3": "朋友"}
}

func (s *SystemModel) MinerState() *P {
	return &P{"0": "收益中", "1": "已结束", "2": "预约中", "3": "已过期"}
}

func (s *SystemModel) Setting(rq ...P) *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	keyP := make(P, 0)
	for _, row := range list {
		keyP[row.Get("key").ToString()] = row.Get("value").ToString()
	}
	if len(rq[0]) == 0 {
		return &AdminResponse{State: SUCCESS, Data: keyP}
	}
	s.Config = keyP
	config.GlobalDB.Delete(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{"_": "1"})
	insert := make([]string, 0)
	for k, v := range rq[0] {
		insert = append(insert, fmt.Sprintf("('%s', '%s')", k, v))
	}
	_, err := config.GlobalDB.Execute(fmt.Sprintf(" INSERT INTO %s(`key`, `value`) VALUES%s", models.DB_TABLE_SYSTEMCONFIG, strings.Join(insert, ", ")))
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "更新配置信息失败！"}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_SITE_CONFIG})
	return &AdminResponse{State: SUCCESS, Data: "更新成功!"}
}

func (s *SystemModel) SettingGet(key string) *db.DBValue {
	if len(s.Config) == 0 {
		s.LoadSiteConfig()
	}
	return s.Config.Ts().Get(key)
}

func (s *SystemModel) LoadSiteConfig() {
	s.Config = make(P, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, item := range list {
		s.Config[item["key"].ToString()] = item["value"].ToString()
	}
}

func (s *SystemModel) SiteCount(rq P) P {
	t := rq.Ts()
	where := make([]string, 0)
	group := ""
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("daytime between  %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("sum").ToInt(); v == 0 {
		group = " GROUP BY daytime "
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SITECOUNT, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{
		"SUM(withdraw) as withdraw",
		"SUM(withdraw_num) as withdraw_num",
		"SUM(recharge) AS recharge",
		"SUM(register_num) AS register_num",
		"SUM(pro_num) AS pro_num",
		"SUM(trade) AS trade",
		"SUM(trade_profit) AS trade_profit",
		"SUM(minning_count) AS minning_count",
		"SUM(first_recharge) AS first_recharge",
		"SUM(minning_profit) as minning_count",
		"SUM(first_recharge_num) AS first_recharge_num",
		"SUM(close_num) AS close_num",
		"SUM(open_num) as open_num",
		"ANY_VALUE(daytime) as daytime",
	}, group, utils.Order(t.Get("sort", "daytime desc").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_SITECOUNT, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	rs := make(P, 0)
	if group != "" {
		rq["sum"] = 1
		total := s.SiteCount(rq)
		rs["total"] = total["list"]
	}
	l := make([]*SiteCount, 0)
	for _, item := range list {
		si := new(SiteCount)
		item.SetObj(si)
		l = append(l, si)
	}
	rs["list"] = l
	rs["count"] = count
	return rs
}

func (s *SystemModel) NoticeList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("pos").ToString(); v != "" {
		where = append(where, fmt.Sprintf("pos = '%s'", v))
	}
	if v := t.Get("lang").ToString(); v != "" {
		where = append(where, fmt.Sprintf("lang = '%s'", v))
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_NOTICE, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_NOTICE, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	nlist := make([]*Notice, 0)
	for _, v := range list {
		n := new(Notice)
		v.SetObj(n)
		nlist = append(nlist, n)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":      nlist,
			"pos_list":  s.NoticePos(),
			"lang_list": s.LangeList(),
			"count":     count,
		},
	}
}

func (s *SystemModel) DelNotice(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{State: ERROR, Data: "请确认一个要删除的公告!"}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_NOTICE, db.DB_PARAMS{"id": id}); err == nil {
		return &AdminResponse{State: SUCCESS, Data: "删除公告信息成功!"}
	}
	config.GlobalRedis.Delete("notice")
	return &AdminResponse{State: ERROR, Data: "删除公告信息失败!"}
}

func (s *SystemModel) OpNotice(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("title").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "公告标题必填!"
		return rs
	}
	if v := t.Get("pos").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "位置信息必填"
		return rs
	}
	if v := t.Get("content").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "公告内容不能为空！"
		return rs
	}
	in := P{
		"title":   t.Get("title").ToString(),
		"pos":     t.Get("pos").ToString(),
		"content": t.Get("content").ToString(),
		"pic":     t.Get("pic").ToString(),
	}
	if t.Get("pubtime").ToString() != "" {
		in["pubtime"] = utils.TimeToint64(t.Get("pubtime").ToString())
	} else {
		in["pubtime"] = utils.GetNow()
	}
	if v := t.Get("lang").ToString(); v != "" {
		in["lang"] = v
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_NOTICE, in, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_NOTICE, in)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作公告信息失败！"
		return rs
	}
	cid := models.MODEL_ASSETS.MakeCacheId("notice", t.Get("id").ToInt())
	config.GlobalRedis.Del(models.HASH_NOTICE, cid)
	rs.State = SUCCESS
	rs.Data = "操作公告信息成功!"
	return rs
}

func (s *SystemModel) RoleList() map[int]*Role {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_ROLE, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make(map[int]*Role, 0)
	for _, role := range list {
		r := new(Role)
		role.SetObj(r)
		rs[role.Get("id").ToInt()] = r
	}
	return rs
}

func (s *SystemModel) StatictisCount() *AdminResponse {
	baseDay := time.Unix(int64(utils.GetNow()-(30*86400)), 0).Local()
	where := fmt.Sprintf(" daytime >= %d ", time.Date(baseDay.Year(), baseDay.Month(), baseDay.Day(), 0, 0, 0, 0, time.Local).Local().Unix())
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SITECOUNT, db.DB_PARAMS{"_": where}, db.DB_FIELDS{
		"SUM(withdraw) as withdraw",
		"SUM(recharge) AS recharge",
		"SUM(register_num) AS register_num",
		"SUM(pro_num) AS pro_num",
		"SUM(trade) AS trade",
		"SUM(trade_profit) AS trade_profit",
		"SUM(minning_count) AS minning_count",
		"SUM(minning_profit) as minning_profit",
		"ANY_VALUE(daytime) as daytime",
	}, " GROUP BY daytime ", " ORDER BY daytime asc")
	rs := make(map[string]interface{}, 0)
	feild := map[string][]string{
		"withdraw_recharge": {"withdraw", "recharge"},
		"register_pro":      {"register_num", "pro_num"},
		"trade_profit":      {"trade", "trade_profit"},
		"mining_profit":     {"minning_count", "minning_profit"},
	}
	transfer := map[string]string{
		"withdraw":       "提现",
		"recharge":       "充值",
		"register_num":   "注册",
		"pro_num":        "有效",
		"trade":          "交易额",
		"trade_profit":   "交易利润",
		"minning_count":  "矿机投资",
		"minning_profit": "矿机返利",
	}
	value := map[string]interface{}{
		"withdraw_recharge": "",
		"register_pro":      "",
		"trade_profit":      "",
		"mining_profit":     "",
	}
	totalData := make(map[string]float64, 0)
	for project, needField := range feild {
		tmp := make(map[string]interface{}, 0)
		legend := make([]string, 0)
		data := make([][]int, 0)
		for _, nk := range needField {
			nd := make([]int, 0)
			day := make([]string, 0)
			for _, item := range list {
				day = append(day, fmt.Sprintf("%d号", time.Unix(int64(item.Get("daytime").ToInt()), 0).Local().Day()))
				for ik, iv := range item {
					if nk == ik {
						nd = append(nd, iv.ToInt())
						totalData[ik] = math.Ceil(totalData[ik] + iv.ToFloat())
					}
				}
			}
			tmp["xAxis"] = day
			data = append(data, nd)
			tmp["Data"] = data
			legend = append(legend, transfer[nk])
		}
		tmp["legend"] = legend
		value[project] = tmp
	}
	rs["total_data"] = totalData
	rs["project"] = value
	return &AdminResponse{State: SUCCESS, Data: rs}
}

func (s *SystemModel) DelRule(rq P) *AdminResponse {
	t := rq.Ts()
	id := t.Get("id").ToInt()
	rs := new(AdminResponse)
	if id == 0 {
		rs.State = 0
		rs.Data = "请确认一个要删除的信息"
		return rs
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_RULE_TEXT, db.DB_PARAMS{"id": id})
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_RULE_TEXT, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one != nil {
		cacheID := models.MODEL_ASSETS.MakeCacheId(one.Get("rule_type").ToString(), one.Get("lang").ToString())
		config.GlobalRedis.Delete(cacheID)
	}
	if err != nil {
		rs.State = 0
		rs.Data = "删除信息失败！"
		return rs
	}
	rs.State = 1
	rs.Data = "删除信息成功!"
	return rs
}

func (s *SystemModel) Rulelist(rq P) *AdminResponse {
	t := rq.Ts()
	condition := make(db.DB_PARAMS)
	if v := t.Get("rule_type").ToString(); v != "" {
		condition["rule_type"] = v
	}
	if v := t.Get("lang").ToString(); v != "" {
		condition["lang"] = v
	}
	count := config.GlobalDB.GetCount(models.DB_TABLE_RULE_TEXT, condition)
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_RULE_TEXT, condition, db.DB_FIELDS{}, utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	for l, item := range list {
		item["content"] = html.UnescapeString(item["content"])
		list[l] = item
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"list":          list,
			"total":         count,
			"langlist":      s.LangeList(),
			"position_list": s.RuleType(),
		},
	}
}

func (s *SystemModel) RuleHandler(rq P) *AdminResponse {
	rs := new(AdminResponse)
	t := rq.Ts()
	if t.Get("lang").ToString() == "" {
		rs.State = ERROR
		rs.Data = "语言不能为空"
		return rs
	}
	if t.Get("rule_type").ToString() == "" {
		rs.State = ERROR
		rs.Data = "文案类型不能为空"
		return rs
	}
	insertData := P{
		"rule_type": t.Get("rule_type").ToString(),
		"lang":      t.Get("lang").ToString(),
		"content":   t.Get("content").ToString(),
		"note":      t.Get("note").ToString(),
	}
	var err error
	if t.Get("id").ToInt() == 0 {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_RULE_TEXT, insertData)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RULE_TEXT, insertData, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "文案信息失败可能已经存在该语种信息!"
		return rs
	}
	cacheID := models.MODEL_ASSETS.MakeCacheId(t.Get("rule_type").ToString(), t.Get("lang").ToString())
	config.GlobalRedis.Del(models.HASH_RULE_TEXT, cacheID)
	rs.State = SUCCESS
	rs.Data = "操作文案信息成功!"
	return rs
}

func (s *SystemModel) Delmsg(snID string) *AdminResponse {
	if snID == "" {
		return &AdminResponse{State: ERROR, Data: "消息不存在!"}
	}
	one, err := config.GlobalDB.FetchOne(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"_": fmt.Sprintf("(sn_id = '%s' OR id = '%s')", snID, snID)}, db.DB_FIELDS{})
	if err != nil {
		return &AdminResponse{State: ERROR, Data: "当前信息不存在!"}
	}
	config.GlobalRedis.PushQueue(models.HASH_USER_MESSAGE, db.DB_PARAMS{"cmd": models.MESSAGE_TYPE_MSG_CANCEL, "uid": one.Get("uid").ToInt(), "content": "cancel"})
	config.GlobalDB.Delete(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"sn_id": snID})
	return &AdminResponse{State: SUCCESS, Data: "cancel成功!"}
}
