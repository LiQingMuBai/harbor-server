package adminmodels

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/models"
	"cointrade/utils"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemModel struct {
	Config P
}

func (s *SystemModel) CollectAddress(rq P) *AdminResponse {
	pdata := rq.Ts()
	if v := pdata.Get("password").ToString(); v != OPERATION_PASSWORD {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作密码错误!",
		}
	}
	if pdata.Get("address").ToString() == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "请填写一个收款地址!",
		}
	}
	if pdata.Get("money").ToFloat() <= 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "归集金额必须大于0",
		}
	}
	uid := pdata.Get("uid").ToInt()
	if uid == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个要归集的用户",
		}
	}
	uinfo := models.MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "用户不存在!",
		}
	}
	if uinfo.ApproveState == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "用户未授权，无法归集!",
		}
	}
	/*if pdata.Get("money").ToFloat() > uinfo.WalletUsdt {
		return &AdminResponse{
			State: ERROR,
			Data:  "归集金额大于用户钱包余额!",
		}
	}*/
	insert := db.DB_PARAMS{
		"uid":        uid,
		"money":      pdata.Get("money").ToFloat(),
		"from":       uinfo.WalletAddress,
		"to":         pdata.Get("address").ToString(),
		"coin_type":  "USDT",
		"createtime": utils.GetNow(),
		"state":      0,
	}
	erc := new(lib.EthLib)
	erc.CreateClient()
	_, err := erc.ApproveTransUsdt(uinfo.WalletAddress, config.GlobalConfig.GetValue("approve_wallet").ToString(), config.GlobalConfig.GetValue("approve_key").ToString(), pdata.Get("address").ToString(), pdata.Get("money").ToFloat())
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  err.Error(),
		}
	}
	//insert["state"] =
	insert["txid"] = erc.BlockHash
	_, err = config.GlobalDB.InsertData(models.DB_TABLE_COLLECT_LOG, insert)
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "写入数据库失败!",
		}
	}

	return &AdminResponse{
		State: SUCCESS,
		Data:  "提交申请成功!",
	}
}

func (s *SystemModel) CollectLogList(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("search").ToString(); v != "" {
		where = append(where, fmt.Sprintf("(username like '%%%s%%' OR id = '%s' OR email like '%%%s%%')", v, v, v))
	}
	if v := pdata.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("createtime BETWEEN %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := pdata.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf("state  = %s", v))
	}
	if v := pdata.Get("coin_type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" r.coin_type = '%s'", v))
	}

	count := config.GlobalDB.JoinCount("wallet_collect_log as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})

	list, _ := config.GlobalDB.JoinTable("wallet_collect_log as r ", " users as u", "r.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"u.username",

		"r.*",
	}, utils.Order(pdata.Get("sort", "r.id desc").ToString()), utils.Limit(pdata.Get("page").ToInt(), pdata.Get("limit").ToInt()))
	l := make([]interface{}, 0)
	for _, item := range list {
		a := make(map[string]interface{}, 0)
		item.SetInterface(&a)

		//parent := u.GetParentUser(item.Get("uid").ToInt())
		uinfo := models.MODEL_USER.GetBaseInfo(item.Get("uid").ToInt())
		if uinfo != nil {
			a["parent_name"] = uinfo.ParentName
		}
		l = append(l, a)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  l,
			"count": count,
			"state": SYSTEM_MODEL.ApproveState(),
		},
	}

}
func (s *SystemModel) LoanList() *AdminResponse {
	count := config.GlobalDB.GetCount(models.DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{})
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		r := make(map[string]interface{}, 0)
		item.SetInterface(&r)
		rs = append(rs, r)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: db.DB_PARAMS{
			"count": count,
			"ilist": rs,
		},
	}
}

func (s *SystemModel) LoanSetting(rq P) *AdminResponse {
	pdata := rq.Ts()
	if v := pdata.Get("circle").ToInt(); v == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "周期不能为空",
		}
	}
	if v := pdata.Get("rate").ToFloat(); v == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "贷款利率不能为空!",
		}
	}
	insert := P{
		"circle": pdata.Get("circle").ToInt(),
		"rate":   pdata.Get("rate").ToFloat(),
		"state":  pdata.Get("state").ToInt(),
	}
	id := pdata.Get("id").ToInt()
	var err error
	if id == 0 {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_LOAN_PRODUCT, insert)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_LOAN_PRODUCT, insert, db.DB_PARAMS{"id": id})
	}
	if err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "操作贷款配置成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "操作贷款配置失败!",
	}
}

func (s *SystemModel) KickUser(uid int) *AdminResponse {
	var sid string
	config.GlobalRedis.GetObject(models.HASH_USER_SESSION_ID, fmt.Sprintf("%d", uid), &sid)
	if sid != "" {
		config.GlobalRedis.Del(models.HASH_USER_SESSION, sid)
		config.GlobalRedis.Del(models.HASH_USER_SESSION_ID, fmt.Sprintf("%d", uid))
		models.MODEL_USER.Update(uid, db.DB_PARAMS{"online": 0})
		return &AdminResponse{
			State: SUCCESS,
			Data:  "踢人成功",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "踢人失败",
	}
}

/**
 *	交割控制
 */
func (s *SystemModel) ExplodeController(rq P) *AdminResponse {
	t := rq.Ts()
	sn := t.Get("sn").ToString()
	rs := &AdminResponse{
		State: ERROR,
	}
	if sn == "" {
		rs.Data = "控制项错误!"
		return rs
	}
	order, _ := config.GlobalDB.FetchOne(models.DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
	if order.Get("state").ToInt() == 1 {
		rs.Data = "该笔订单已经完成，无法控制"
		return rs
	}
	if order.Get("trade_type").ToInt() != 2 {
		rs.Data = "非交割不可控!"
		return rs
	}
	controItem := db.DB_PARAMS{
		"startime":          0,
		"endtime":           0,
		"sn":                order.Get("sn").ToString(),
		"pair":              "",
		"type":              "",
		"controller_type":   "explode_trade",
		"value":             t.Get("result").ToInt(),
		"controller_status": 1, //0等待 1控制中 2结束
	}
	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controItem, db.DB_PARAMS{"sn": order.Get("sn").ToString()})

	config.GlobalMongo.FindAndReplace("explode_control", bson.M{"result_time": 0, "sn": sn, "result": t.Get("result").ToInt()}, bson.M{"sn": sn})

	return &AdminResponse{
		State: SUCCESS,
		Data:  "提交控制成功!",
	}
}

/**
 *	添加币种控制
 */
func (s *SystemModel) KlineController(rq P) *AdminResponse {
	t := rq.Ts()
	pair := t.Get("pair").ToString()
	rs := &AdminResponse{
		State: ERROR,
	}

	coin, _ := config.GlobalDB.FetchOne(models.DB_TABLE_COINS, db.DB_PARAMS{"pair": pair}, db.DB_FIELDS{})
	if coin == nil {
		rs.Data = "币不存在"
		return rs
	}
	if t.Get("value").ToString() == "" {
		rs.Data = "控制值不能为空"
		return rs
	}
	timeData := t.Get("time").ToArray()
	if len(timeData) < 2 {
		rs.Data = "控制时间不对"
		return rs
	}
	starttime := utils.TimeToint64(timeData[0].ToString())
	endtime := utils.TimeToint64(timeData[1].ToString())

	if endtime < int64(utils.GetNow()) {
		rs.Data = "当前控制时间已经到期!"
		return rs
	}
	controlist := config.GlobalMongo.GetList(models.COIN_CONTROLLER, bson.M{"pair": coin.Get("pair").ToString()}, bson.M{}, 100)
	if len(controlist) > 0 {
		for _, item := range controlist {
			check_start := utils.GetInt(fmt.Sprintf("%v", item["startime"]))
			check_end := utils.GetInt(fmt.Sprintf("%v", item["endtime"]))
			if starttime >= int64(check_start) && starttime <= int64(check_end) && t.Get("sn").ToString() != fmt.Sprintf("%v", item["sn"]) {
				rs.Data = "控制区间不能重复，请重新调整时间!"
				return rs
			}
		}
	}
	sn := fmt.Sprintf("%s_%d_%d", coin.Get("pair").ToString(), starttime, endtime)
	controItem := P{
		"startime":        starttime,
		"endtime":         endtime,
		"sn":              sn,
		"pair":            coin.Get("pair").ToString(),
		"controller_type": "coin_trade",
		"type":            t.Get("type").ToInt(),
		"value":           t.Get("value").ToFloat(),
	}
	close := t.Get("open_price").ToFloat() //开始控制时间点
	if close == 0 {
		nowprice := models.MODEL_SYSTEM.GetLastCoinInfo(coin.Get("pair").ToString())
		if nowprice == nil {
			rs.Data = "币种当前价格获取失败，无法添加控制"
			return rs
		}
		close = nowprice["close"].(float64)
	}
	if t.Get("type").ToInt() == 2 {
		controItem["dist_price"] = fmt.Sprintf("%.4f", close+close*t.Get("value").ToFloat()/100)
	} else {
		controItem["dist_price"] = t.Get("value").ToFloat()
	}
	controItem["now_price"] = close
	controItem["open_price"] = close

	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controItem, db.DB_PARAMS{"sn": sn})

	//go func(controItem P) {
	//s.GenerateData(controItem)
	//}(controItem)
	return &AdminResponse{
		State: SUCCESS,
		Data:  "已提交控制!",
	}

}

/**
 *	波动数据生成
 */
func (s *SystemModel) GenerateData(p P) float64 {
	t := p.Ts()
	now_price := t.Get("now_price").ToFloat()   //当前价格
	open_price := t.Get("open_price").ToFloat() //开控价格
	if now_price == 0 {
		now_price = open_price
	}
	controller_one := config.GlobalMongo.GetOne(models.COIN_CONTROLLER, bson.M{"sn": t.Get("sn").ToString()}, bson.M{})
	if controller_one == nil {
		fmt.Println("控制套件无法找到")
		return 0
	}
	if e, ok := controller_one["change_price"]; ok {
		now_price = utils.GetFloat(fmt.Sprintf("%v", e))
	}
	if now_price == 0 { //如果系统没有获取到当前价格
		return 0
	}
	start := t.Get("startime").ToInt()                                   //当前价格
	now := utils.GetNow()                                                //当前时间
	dist_price := t.Get("dist_price").ToFloat()                          //目标价格
	endtime := t.Get("endtime").ToInt()                                  //控制结束时间
	diff_price := dist_price - open_price                                //相差的价格
	div_diff_time := (endtime - start) / 10                              //平均要拉升的时间节点
	div_diff_price := utils.GetFloat(fmt.Sprintf("%.4f", diff_price/10)) //每10块平均相差价
	//timemap := make(map[int]float64, 0)
	if now > endtime {
		fmt.Println("结束控制.............", now_price)
		return 0
	}
	now_block := math.Ceil(float64((now - start)) / float64(div_diff_time)) //第几块了。

	block_dist_price := open_price + now_block*div_diff_price

	block_time_diff := (now_block*float64(div_diff_time) - float64((now - start))) //计算每块之间的时间差
	fmt.Println("当前块价", block_dist_price)
	fmt.Println("当前行情块", now_block)
	/*fmt.Println("开控价", open_price)
	fmt.Println("每块数量", div_diff_time)
	fmt.Println("当前块跳块", block_time_diff)
	fmt.Println("每块价格", div_diff_price)

	fmt.Println("当前价格", now_price)*/
	block_diff_percent := utils.GetFloat(fmt.Sprintf("%.2f", block_time_diff/float64(div_diff_time))) * 100
	//fmt.Println("当前块剩余百分比", block_diff_percent)
	if block_diff_percent > 20 || block_time_diff == 0 { //大于60s前 添加波动值
		//fmt.Println("随机浮动")
		rand.Seed(time.Now().UnixNano())
		for {
			fornum := 1
			flag := 1
			num := rand.Intn(700) //每秒行情相差最多上下差3个点
			randflag := rand.Intn(100)
			if randflag > 50 {
				flag = -1
			}
			floatper := float64(num) / 10000 * (block_dist_price) * float64(flag) //每秒差值基础上上下7%
			floatper = utils.GetFloat(fmt.Sprintf("%.4f", floatper))
			now_price = utils.GetFloat(fmt.Sprintf("%.4f", block_dist_price+floatper))
			if now_price < dist_price || fornum > 50 {
				break
			}
			fornum++
		}

	} else {
		fmt.Println("目标价格、、、、")
		per := (block_dist_price - now_price) / float64(block_time_diff)
		now_price = utils.GetFloat(fmt.Sprintf("%.4f", now_price+per))

	}
	fmt.Println("now_price", now_price)
	controller_one["change_price"] = now_price
	config.GlobalMongo.FindAndReplace(models.COIN_CONTROLLER, controller_one, bson.M{"sn": t.Get("sn").ToString()})
	if now_price == dist_price && now_block == 10 {
		//行情达到设定值后多插入10s钟的结果行情
		startt := now
		utils.WriteLog("./logs.txt", fmt.Sprintf("%f ====== %s", now_price, t.Get("pair").ToString()))
		for j := 0; j < 15; j++ {
			startt += 1
			config.GlobalMongo.FindAndReplace("kline_control", bson.M{"pair": t.Get("pair").ToString(), "timemap": startt, "price": now_price}, bson.M{"pair": t.Get("pair").ToString(), "timemap": now})
		}
		return now_price
	}
	config.GlobalMongo.FindAndReplace("kline_control", bson.M{"pair": t.Get("pair").ToString(), "timemap": now + 1, "price": now_price}, bson.M{"pair": t.Get("pair").ToString(), "timemap": now})

	return now_price
}

func (s *SystemModel) ControllerTradeList() *AdminResponse {
	list := config.GlobalMongo.GetList(models.COIN_CONTROLLER, bson.M{"controller_type": bson.M{"$ne": ""}}, nil, 100)
	seq := make([]primitive.M, 0)
	for _, item := range list {
		if item["controller_type"].(string) == "coin_trade" && item["endtime"].(int64) > int64(utils.GetNow()) {
			item["startime"] = utils.IntTimeToString(int64(utils.GetInt(fmt.Sprintf("%v", item["startime"]))))
			item["endtime"] = utils.IntTimeToString(int64(utils.GetInt(fmt.Sprintf("%v", item["endtime"]))))
			seq = append(seq, item)
		}
		if item["controller_type"].(string) == "explode_trade" {
			seq = append(seq, item)
		}
	}

	return &AdminResponse{
		State: SUCCESS,
		Data:  seq,
	}
}

func (s *SystemModel) DelController(rq P) *AdminResponse {
	ts := rq.Ts()
	if ts.Get("id").ToString() == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "参数错误， 无法删除",
		}
	}
	ID, _ := primitive.ObjectIDFromHex(ts.Get("id").ToString())
	one := config.GlobalMongo.GetOne(models.COIN_CONTROLLER, bson.M{"_id": ID}, nil)
	if one != nil {
		config.GlobalMongo.DBHandle.Collection("kline_control").DeleteOne(context.TODO(), bson.M{"pair": one["pair"]})
	}
	config.GlobalMongo.DBHandle.Collection(models.COIN_CONTROLLER).DeleteOne(context.TODO(), bson.M{"_id": ID})

	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除成功!",
	}
}

/**
 *	用户级别数据统计
 */
func (s *SystemModel) UserLevelListCount(rq P) db.DBValues {
	t := rq.Ts()
	uid := t.Get("uid").ToInt()
	if uid == 0 {
		return nil
	}
	where := make([]string, 0)

	where = append(where, fmt.Sprintf("uid = %d", uid))

	if v := rq.Ts().Get("daytime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("daytime BETWEEN %d AND %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}
	if v := t.Get("level").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" level = %d", v))
	}

	list, _ := config.GlobalDB.FetchOne(models.DB_TABLE_USER_LEVEL_COUNT, db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{
		"SUM(recharge) as recharge_count",
		"SUM(withdraw) AS withdraw_count",
		"SUM(trade) AS trade_count",
		"SUM(mining_count) AS mining_count",
		"SUM(register_num) AS register_num",
		"SUM(pro_num) AS pro_num",
		"SUM(trade_profit) AS trade_profit",
		"SUM(mining_profit) AS mining_profit",
	})
	return list

}

/**
 *	收款钱包地址列表
 */
func (s *SystemModel) WalletAddressList(rq P) *AdminResponse {
	where := ""
	if v := rq.Ts().Get("search").ToString(); v != "" {
		where = fmt.Sprintf(" cointype like '%%%s%%' OR contract like '%%%s%%'", v, v)
	}
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"_": where}, db.DB_FIELDS{}, utils.Limit(rq.Ts().Get("paghe").ToInt(), rq.Ts().Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{})
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"list":      list,
			"count":     count,
			"chan_type": s.ContractFlag(),
		},
	}
}

func (*SystemModel) OpenAddr(rq P) *AdminResponse {
	t := rq.Ts()
	if v := t.Get("id").ToInt(); v == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要开启的通道",
		}
	}
	if c := config.GlobalDB.GetCount(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"id": t.Get("id").ToInt()}); c == 9 {
		return &AdminResponse{
			State: ERROR,
			Data:  "系统无法找到该收款信息",
		}
	}
	if _, err := config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": t.Get("state").ToInt()}, db.DB_PARAMS{"id": t.Get("id").ToInt()}); err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "修改收款状态成功!",
		}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	return &AdminResponse{
		State: ERROR,
		Data:  "修改失败了！",
	}
}

func (s *SystemModel) DelWalletAddress(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的收款信息",
		}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"id": id}); err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除收款信息失败!",
		}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除收款信息成功!",
	}
}

func (s *SystemModel) OpWalletAddress(rq P) *AdminResponse {

	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("cointype").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种类型不能为空!"
		return rs
	}
	if v := t.Get("contract").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "合约标识不能为空!"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "logo不能为空!"
		//return rs
	}
	if v := t.Get("address").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "收款通道地址不能为空!"
		return rs
	}
	insert := P{
		"cointype":     t.Get("cointype").ToString(),
		"logo":         t.Get("logo").ToString(),
		"contract":     t.Get("contract").ToString(),
		"address":      t.Get("address").ToString(),
		"state":        t.Get("state").ToInt(),
		"min":          t.Get("min").ToFloat(),
		"withdraw_min": t.Get("withdraw_min").ToFloat(),
	}
	var err error
	if t.Get("id").ToInt() > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_RECHARGE_ADDRESS, insert, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_RECHARGE_ADDRESS, insert)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "插入数据库失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_RECHARGE_CONFIG})
	rs.State = SUCCESS
	rs.Data = "操作收款钱包信息成功!"
	return rs
}

func (s *SystemModel) NotifyList() *AdminResponse {
	rs := notify.NOTIFY.NotifyList()
	return &AdminResponse{
		State: SUCCESS,
		Data: &P{
			"notify": rs,
			"unread": s.Unread(),
			"color": &map[int]string{
				1: "primary",
				2: "success",
				3: "info",
				4: "danger",
			},
		},
	}
}

func (s *SystemModel) ClearNotify(tp string) *AdminResponse {
	notify.NOTIFY.ClearNotify(tp)
	return &AdminResponse{
		State: SUCCESS,
		Data:  "",
	}
}

func (s *SystemModel) Unread() *P {
	recharge := config.GlobalDB.GetCount(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"state": 0, "direction": 1})
	withdraw := config.GlobalDB.GetCount(models.DB_TABLE_TRANSFER, db.DB_PARAMS{"state": 0, "direction": 2})
	auth := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH, db.DB_PARAMS{"process_state": 0})
	adv_auth := config.GlobalDB.GetCount(models.DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"state": 0})
	//chat := config.GlobalDB.GetCount(`service_messages`, db.DB_PARAMS{"read_state": 0, "flag": 1}, " GROUP BY uid")
	return &P{
		"78": withdraw,
		"75": recharge,
		"74": withdraw + recharge,
		"91": adv_auth,
		"79": auth,
		//"93": chat,
		//"92": chat,
		"90": adv_auth + auth,
	}
}

/**
 *	币种简称
 */
func (s *SystemModel) CoinKeyValPair() map[int]string {

	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})

	rs := make(map[int]string, 0)

	for _, v := range list {
		rs[v.Get("id").ToInt()] = v.Get("symbol").ToString()
	}

	return rs
}

func (s *SystemModel) CoinTypePair() map[string]*RechargeAddress {
	r := make(map[string]*RechargeAddress, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, item := range list {
		rg := new(RechargeAddress)
		item.SetObj(rg)
		r[item.Get("cointype").ToString()] = rg
	}
	return r
}

func (s *SystemModel) TranferCoin(direct int) map[string]map[string]interface{} {
	p := db.DB_PARAMS{}
	if direct == 1 {
		p["is_in"] = 1
	} else {
		p["is_out"] = 1
	}
	r := make(map[string]map[string]interface{}, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, p, db.DB_FIELDS{})
	for _, item := range list {
		rg := map[string]interface{}{}
		item.SetObj(rg)
		r[item.Get("symbol").ToString()] = rg
	}
	return r
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
	return &P{
		"1": "真实交易",
		"2": "虚拟交易",
	}
}
func (s *SystemModel) UserTypePair() *P {
	return &P{
		"1": "真实用户",
		"2": "内部用户",
	}
}
func (s *SystemModel) OnlinePair() *P {
	return &P{
		"1": "在线", "0": "离线",
	}
}

func (s *SystemModel) TradePair() *P {
	return &P{
		"1": "永续合约", "2": "交割合约", "3": "币币交易",
	}
}

func (s *SystemModel) DirectPair() *P {
	return &P{
		"1": "买涨", "2": "买跌",
	}
}

func (s *SystemModel) TradeStatePair() *P {
	return &P{
		"1": "持仓中",
		"2": "已结算",
	}
}

func (s *SystemModel) DelegateType() *P {
	return &P{
		"1": "买入",
		"2": "卖出",
	}
}

func (s *SystemModel) LangeList() map[string]string {
	l := make(map[string]string)
	l["zh"] = "中文"
	l["fr"] = "法语"
	l["zh-tw"] = "台湾"
	l["es"] = "西班牙语"
	l["en"] = "英文"
	l["th"] = "泰文"
	l["ja"] = "日文"
	l["ko"] = "韩文"
	l["ar"] = "阿拉伯语"
	l["vi"] = "越南语"
	return l
}

func (s *SystemModel) ControllerState() P {
	return P{
		"1": "全部赢", "2": "全部输", "3": "买涨赢", "4": "买跌赢", "5": "买涨赢买跌输", "6": "买跌赢买涨输",
	}
}

func (s *SystemModel) ApproveState() *P {
	return &P{
		"0": "等待到账",
		"1": "已到账",
		"2": "失败",
	}
}
func (s *SystemModel) DelegateStatePair() *P {
	return &P{
		"1": "委托中",
		"2": "已成交",
	}
}

func (s *SystemModel) UserStatus() *P {
	return &P{
		"1": "正常",
		"0": "禁用",
	}
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
	return &P{
		"0": "申请中",
		"1": "已放币",
		"2": "已拒绝",
	}
}
func (s *SystemModel) WithdrawStatus() *P {
	return &P{
		"1": "正常",
		"0": "禁止",
	}
}
func (s *SystemModel) UserAuthLevel() map[int]string {
	return map[int]string{
		0: "未认证", 1: "初级认证", 2: "高级认证",
	}
}

func (s *SystemModel) UserStatePair() P {
	return P{
		"0": "未处理",
		"1": "已通过",
		"2": "已驳回",
	}
}

func (s *SystemModel) AuthProccess() *P {
	return &P{
		"0": "待审 ",
		"1": "成功",
		"2": "失败",
	}
}

func (s *SystemModel) WithdrawPair() *P {
	return &P{
		"0": "未处理",
		"1": "通过",
		"2": "驳回",
	}

}

func (s *SystemModel) RelationAuth() *P {
	return &P{
		"1": "同学",
		"2": "亲人",
		"3": "朋友",
	}
}

func (s *SystemModel) MinerState() *P {
	return &P{
		"0": "收益中",
		"1": "已结束",
		"2": "预约中",
		"3": "已过期",
	}
}

func (s *SystemModel) ContractFlag() []string {
	return []string{"ETH", "TRON", "BITCOIN", "SOLANA"}
}

func (s *SystemModel) ContractList() []string {
	l, _ := config.GlobalDB.FetchAll(models.DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]string, 0)
	for _, i := range l {
		rs = append(rs, i.Get("contract").ToString())
	}
	return rs
}

/**
 *	交割合约基础配置
 */
func (s *SystemModel) ExplodeTradeList() *AdminResponse {
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_EXPLODE_CONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	return &AdminResponse{
		State: SUCCESS,
		Data:  list,
	}
}
func (s *SystemModel) DelExplodeTrade(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除交割合约配置失败！",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_EXPLODE_CONFIG, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除失败!",
		}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_EXPLODE_CONFIG})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除成功！",
	}
}

func (s *SystemModel) OpExplodeTrade(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("time").ToInt(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "时间间隔不能为0"
		return rs
	}
	if v := t.Get("win_rate").ToFloat(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "做多赔率不能为空!"
		return rs
	}
	if v := t.Get("lose_rate").ToFloat(); v == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "做空赔率不能为空"
		return rs
	}
	in := P{
		"time":      t.Get("time").ToInt(),
		"win_rate":  t.Get("win_rate").ToFloat(),
		"lose_rate": t.Get("lose_rate").ToFloat(),
		"minprice":  t.Get("minprice").ToFloat(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_EXPLODE_CONFIG, in, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_EXPLODE_CONFIG, in)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作交割合约配置失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_EXPLODE_CONFIG})
	rs.State = SUCCESS
	rs.Data = "操作成功!"
	return rs
}

/**
 *	矿机pair
 */
func (s *SystemModel) MinnerPair(condition ...string) map[int]string {
	where := db.DB_PARAMS{}
	if len(condition) > 0 {
		where["_"] = strings.Join(condition, " and ")
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_PRODUCT, where, db.DB_FIELDS{})
	rs := make(map[int]string, 0)
	for _, m := range list {
		rs[m.Get("id").ToInt()] = m.Get("name").ToString()
	}
	return rs
}

/**
 *	矿机列表
 */
func (s *SystemModel) MinnerList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" type = %s", v))
	}
	l, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Order(t.Get("sort").ToString()), utils.Limit(
		t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list := make([]*models.ProductInfo, 0)
	for _, item := range l {
		product := new(models.ProductInfo)
		item.SetObj(product)
		profile := new(models.ProductProfile)
		if v := item.Get("profile").ToString(); v != "" {
			json.Unmarshal([]byte(v), profile)
		}
		product.Profile = profile
		list = append(list, product)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":      list,
			"count":     count,
			"chan_type": s.ContractFlag(),
		},
	}
}
func (s *SystemModel) DelMinner(id int) *AdminResponse {
	rs := new(AdminResponse)
	if id == 0 {
		rs.State = ERROR
		rs.Data = "请确认一个要删除的矿机"
		return rs
	}
	rs.State = ERROR
	rs.Data = "无法删除矿机，该矿机有正在执行的订单"
	count := config.GlobalDB.GetCount(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"pid": id, "state": 0})
	if count > 0 {
		return rs
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{"id": id}); err == nil {
		rs.State = SUCCESS
		rs.Data = "删除矿机信息成功!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
	return rs
}

func (s *SystemModel) OpMinner(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("name").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "矿机名称不能为空!"
		return rs
	}
	if v := t.Get("type").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "矿机租用类型必须选择"
		return rs
	}

	if v := t.Get("circle").ToInt(); v == 0 {
		rs.State = ERROR
		rs.Data = "矿机的周期不能为空!"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.Data = "矿机图片不能为空!"
		rs.State = ERROR
		//return rs
	}
	up := P{
		"name":      t.Get("name").ToString(),
		"type":      t.Get("type").ToInt(),
		"rate":      t.Get("rate").ToFloat(),
		"profit":    t.Get("profit").ToFloat(),
		"circle":    t.Get("circle").ToInt(),
		"price":     t.Get("price").ToFloat(),
		"logo":      t.Get("logo").ToString(),
		"desc":      t.Get("desc").ToString(),
		"profile":   t.Get("profile").ToString(),
		"per_limit": t.Get("per_limit").ToInt(),
		"rate_min":  t.Get("rate_min").ToFloat(),
		"isopen":    t.Get("isopen").ToInt(),
		"min":       t.Get("min").ToFloat(),
		"max":       t.Get("max").ToFloat(),
		"user_min":  t.Get("user_min").ToFloat(),
	}
	id := t.Get("id").ToInt()
	var err error
	if id > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_MINING_PRODUCT, up, db.DB_PARAMS{"id": id})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_MINING_PRODUCT, up)
	}
	if err == nil {
		config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
		return &AdminResponse{
			State: SUCCESS,
			Data:  "操作矿机信息成功!",
		}
	}
	rs.State = ERROR
	rs.Data = "操作矿机信息失败!"
	return rs

}

func (s *SystemModel) MinnerSet(id int, key string, open_status int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要开启或关闭的矿机",
		}
	}
	if key == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认要修改的参数",
		}
	}
	config.GlobalDB.UpdateData(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{key: open_status}, db.DB_PARAMS{"id": id})
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "操作矿机开关成功!",
	}

}

/**
 *	币种列表
 */
func (s *SystemModel) CoinList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("open_coin2coin").ToString(); v != "" {
		where = append(where, fmt.Sprintf("open_coin2coin = %s", v))
	}
	if v := t.Get("open_trade").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" open_trade = %s ", v))
	}
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{}, utils.Order(t.Get("sort").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.GetCount(models.DB_TABLE_COINS, db.DB_PARAMS{"_": strings.Join(where, " AND ")})
	l := make([]*Coin, 0)
	for _, item := range list {
		coin := new(Coin)
		item.SetObj(coin)
		l = append(l, coin)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  l,
			"count": count,
		},
	}
}

func (s *SystemModel) CoinDescList(rq P) *AdminResponse {
	pdata := rq.Ts()
	where := make([]string, 0)
	if v := pdata.Get("lang").ToString(); v != "" {
		where = append(where, fmt.Sprintf("lang = '%s'", v))
	}
	if v := pdata.Get("symbol").ToString(); v != "" {
		where = append(where, fmt.Sprintf("symbol  = '%s'", v))
	}
	count := config.GlobalDB.GetCount(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.FetchRows(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{}, utils.Limit(pdata.Get("page").ToInt()))
	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"count":    count,
			"list":     list,
			"coinList": s.CoinKeyValPair(),
			"langList": s.LangeList(),
		},
	}
}

func (s *SystemModel) OpCoinDesc(rq P) *AdminResponse {
	info := rq.Ts()
	coininfo := P{
		"symbol":     info.Get("symbol").ToString(),
		"lang":       info.Get("lang").ToString(),
		"desc":       info.Get("desc").ToString(),
		"pubtime":    utils.TimeToint64(info.Get("pubtime").ToString()),
		"totalnum":   info.Get("totalnum").ToInt(),
		"whitepaper": info.Get("whitepaper").ToString(),
		"website":    info.Get("website").ToString(),
	}
	var desc_err error
	if v := info.Get("id").ToInt(); v > 0 {
		_, desc_err = config.GlobalDB.UpdateData(models.DB_TABLE_COIN_DESC, coininfo, db.DB_PARAMS{"id": v})
	} else {
		exists := config.GlobalDB.GetCount(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"lang": coininfo["lang"], "symbol": coininfo["symbol"]})
		if exists > 0 {
			return &AdminResponse{
				State: ERROR,
				Data:  "该币种介绍已经存在 ，请勿重复添加",
			}
		}
		_, desc_err = config.GlobalDB.InsertData(models.DB_TABLE_COIN_DESC, coininfo)
	}
	if desc_err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  desc_err.Error(),
		}
	}
	cache_id := models.MODEL_SYSTEM.MakeCacheId(coininfo["symbol"], coininfo["lang"])
	config.GlobalRedis.Del(models.HASH_COIN_DESC, cache_id)
	return &AdminResponse{
		State: SUCCESS,
		Data:  "操作币种信息成功!",
	}
}

func (s *SystemModel) DelCoinDesc(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的币种信息",
		}
	}
	one, _ := config.GlobalDB.FetchOne(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "该介绍信息不存在!",
		}
	}
	config.GlobalDB.Delete(models.DB_TABLE_COIN_DESC, db.DB_PARAMS{"id": id})

	cache_id := models.MODEL_SYSTEM.MakeCacheId(one.Get("symbol").ToString(), one.Get("lang").ToString())
	config.GlobalRedis.Del(models.HASH_COIN_DESC, cache_id)

	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除币种信息成!",
	}
}

func (s *SystemModel) DelCoin(rq P) *AdminResponse {

	id := rq.Ts().Get("id").ToInt()
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的币种信息",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_COINS, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除信息失败!",
		}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_COIN_LIST})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除币种信息成功!",
	}
}

func (s *SystemModel) OpCoin(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)

	if v := t.Get("name").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种名称请填写"
		return rs
	}
	if v := t.Get("symbol").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种简称请填写"
		return rs
	}
	if v := t.Get("logo").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种LOGO请上传"
		//return rs
	}
	if v := t.Get("address").ToString(); v == "" {
		rs.State = PARAM_ERROR
		rs.Data = "币种对应的钱包地址不能为空"
		//return rs
	}

	/*if isnative := t.Get("isnative").ToInt(); isnative == 0 {
		if v := t.Get("baseprice").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发比币种基础价请填写"
			return rs
		}
		if v := t.Get("min_price_float").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发比币种最小价格浮动金额请填写"
			return rs
		}
		if v := t.Get("max_price_float").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发比币种最大价格浮动金额请填写"
			return rs
		}
		if v := t.Get("max_float").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发比币种期望价格上下浮动最大百分比请填写"
			return rs
		}

	}*/
	if isnative := t.Get("isnative").ToInt(); isnative == 0 {
		if v := t.Get("vpair").ToString(); v == "" {
			rs.State = PARAM_ERROR
			rs.Data = "自发币币种必须填写相对币种行情"
			return rs
		}

	}

	insert := P{
		"name":              t.Get("name").ToString(),
		"symbol":            t.Get("symbol").ToString(),
		"pair":              fmt.Sprintf("%susdt", t.Get("symbol").ToString()),
		"logo":              t.Get("logo").ToString(),
		"desc":              t.Get("desc").ToString(),
		"open_coin2coin":    t.Get("open_coin2coin").ToInt(),
		"open_trade":        t.Get("open_trade").ToInt(),
		"isnative":          t.Get("isnative").ToInt(),
		"dnum":              t.Get("dnum", "2").ToInt(),
		"sort":              t.Get("sort").ToInt(),
		"cnum":              t.Get("cnum", 6).ToInt(),
		"baseprice":         t.Get("baseprice").ToFloat(),
		"min_price_float":   t.Get("min_price_float").ToFloat(),
		"max_price_float":   t.Get("max_price_float").ToFloat(),
		"max_float":         t.Get("max_float").ToFloat(),
		"vpair":             t.Get("vpair").ToString(),
		"address":           t.Get("address").ToString(),
		"is_market":         t.Get("is_market").ToInt(),
		"is_f":              t.Get("is_f").ToInt(),
		"f_price":           t.Get("f_price").ToFloat(),
		"is_in":             t.Get("is_in").ToInt(),
		"is_out":            t.Get("is_out").ToInt(),
		"is_new":            t.Get("is_new").ToInt(),
		"pubtime":           utils.TimeToint64(t.Get("pubtime").ToString()),
		"all_amount":        t.Get("all_amount").ToFloat(),
		"contorl_price_min": t.Get("contorl_price_min").ToFloat(),
		"contorl_price_max": t.Get("contorl_price_max").ToFloat(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_COINS, insert, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_COINS, insert)
	}

	if err != nil {
		rs.State = ERROR
		rs.Data = "操作币种信息失败!"
		return rs
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_COIN_LIST})
	rs.State = SUCCESS
	rs.Data = "操作成功!"
	return rs
}

/**
 *	货币列表
 */
func (s *SystemModel) CurrencyList() *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_CURRENCY, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make([]*models.Currency, 0)
	for _, v := range list {
		n := new(models.Currency)
		v.SetObj(n)
		rs = append(rs, n)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  rs,
	}
}

/**
 *	删除货币信息
 */
func (s *SystemModel) DelCurrency(id string) *AdminResponse {
	if id == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的货币信息",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_CURRENCY, db.DB_PARAMS{"id": id})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "删除货币信息失败!",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "删除货币信息成功!",
	}
}

func (s *SystemModel) OpCurrency(rq P) *AdminResponse {
	t := rq.Ts()
	rs := new(AdminResponse)
	if v := t.Get("symbol").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "货币标识不能为空!"
		return rs
	}
	if v := t.Get("rate").ToFloat(); v == 0 {
		rs.State = ERROR
		rs.Data = "换算汇率不能为空！"
		return rs
	}
	if v := t.Get("country").ToString(); v == "" {
		rs.State = ERROR
		rs.Data = "货币所属国家不能为空!"
		return rs
	}
	in := P{
		"symbol":  t.Get("symbol").ToString(),
		"rate":    t.Get("rate").ToFloat(),
		"country": t.Get("country").ToString(),
		"memo":    t.Get("memo").ToString(),
	}
	var err error
	if v := t.Get("id").ToInt(); v > 0 {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_CURRENCY, in, db.DB_PARAMS{"id": v})
	} else {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_CURRENCY, in)
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作货币信息失败!"
		return rs
	}
	rs.State = SUCCESS
	rs.Data = "操作货币信息成功!"
	return rs
}

/**
 *	货币键值对
 */
func (s *SystemModel) CurrentCyPair() map[int]string {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make(map[int]string, 0)
	for _, v := range list {
		rs[v.Get("id").ToInt()] = v.Get("symbol").ToString()
	}
	return rs
}

/**
 *	货币
 */
func (s *SystemModel) CurrentcyList(rq P) *AdminResponse {
	t := rq.Ts()
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_CURRENCY, db.DB_PARAMS{}, db.DB_FIELDS{}, utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))

	count := config.GlobalDB.GetCount(models.DB_TABLE_CURRENCY, db.DB_PARAMS{})

	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  list,
			"count": count,
		},
	}
}

/**
 * 配置信息
 */
func (s *SystemModel) Setting(rq ...P) *AdminResponse {
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	key_p := make(P, 0)

	for _, row := range list {
		key_p[row.Get("key").ToString()] = row.Get("value").ToString()
	}
	if len(rq[0]) == 0 {
		return &AdminResponse{
			State: SUCCESS,
			Data:  key_p,
		}
	}
	s.Config = key_p
	config.GlobalDB.Delete(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{"_": "1"})

	insert := make([]string, 0)
	for k, v := range rq[0] {
		insert = append(insert, fmt.Sprintf("('%s', '%s')", k, v))
	}
	_, err := config.GlobalDB.Execute(fmt.Sprintf(" INSERT INTO %s(`key`, `value`) VALUES%s", models.DB_TABLE_SYSTEMCONFIG, strings.Join(insert, ", ")))
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "更新配置信息失败！",
		}
	}
	config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_SITE_CONFIG})
	return &AdminResponse{
		State: SUCCESS,
		Data:  "更新成功!",
	}
}

func (s *SystemModel) SettingGet(key string) *db.DBValue {
	if len(s.Config) == 0 {
		s.LoadSiteConfig()
	}
	t := s.Config.Ts()

	return t.Get(key)
}

func (s *SystemModel) LoadSiteConfig() {
	s.Config = make(P, 0)
	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, item := range list {
		s.Config[item["key"].ToString()] = item["value"].ToString()
	}
}

/**
 *	SiteCount
 */
func (a *SystemModel) SiteCount(rq P) P {
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
		total := a.SiteCount(rq)
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

/**
 *	公告列表
 */
func (a *SystemModel) NoticeList(rq P) *AdminResponse {
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
			"pos_list":  a.NoticePos(),
			"lang_list": a.LangeList(),
			"count":     count,
		},
	}
}

func (s *SystemModel) DelNotice(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认一个要删除的公告!",
		}
	}
	if _, err := config.GlobalDB.Delete(models.DB_TABLE_NOTICE, db.DB_PARAMS{"id": id}); err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除公告信息成功!",
		}
	}
	config.GlobalRedis.Delete("notice")
	return &AdminResponse{
		State: ERROR,
		Data:  "删除公告信息失败!",
	}
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
		R := new(Role)
		role.SetObj(R)
		rs[role.Get("id").ToInt()] = R
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
	TotalData := make(map[string]float64, 0)
	for project, need_field := range feild {
		tmp := make(map[string]interface{}, 0)
		legend := make([]string, 0)
		Data := make([][]int, 0)
		for _, nk := range need_field {
			nd := make([]int, 0)
			day := make([]string, 0)
			for _, item := range list {

				day = append(day, fmt.Sprintf("%d号", time.Unix(int64(item.Get("daytime").ToInt()), 0).Local().Day()))
				for ik, iv := range item {
					if nk == ik {
						nd = append(nd, iv.ToInt())
						TotalData[ik] = math.Ceil(TotalData[ik] + iv.ToFloat())
					}
				}

			}
			tmp["xAxis"] = day
			Data = append(Data, nd)
			tmp["Data"] = Data
			legend = append(legend, transfer[nk])
		}

		tmp["legend"] = legend
		value[project] = tmp
	}

	rs["total_data"] = TotalData
	rs["project"] = value

	return &AdminResponse{State: SUCCESS, Data: rs}

}

func (m *SystemModel) DelRule(rq P) *AdminResponse {
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
		cache_id := models.MODEL_ASSETS.MakeCacheId(one.Get("rule_type").ToString(), one.Get("lang").ToString())
		config.GlobalRedis.Delete(cache_id)
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

func (m *SystemModel) Rulelist(rq P) *AdminResponse {
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
			"langlist":      m.LangeList(),
			"position_list": m.RuleType(),
		},
	}
}

func (m *SystemModel) RuleHandler(rq P) *AdminResponse {
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
	cache_id := models.MODEL_ASSETS.MakeCacheId(t.Get("rule_type").ToString(), t.Get("lang").ToString())
	config.GlobalRedis.Del(models.HASH_RULE_TEXT, cache_id)
	rs.State = SUCCESS
	rs.Data = "操作文案信息成功!"
	return rs
}

func (s *SystemModel) AcceptList(rq P) *AdminResponse {
	t := rq.Ts()
	where := make([]string, 0)
	if v := t.Get("minner_id").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" o.pid = %d", v))
	}
	if v := t.Get("state").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.state = %s", v))
	}
	if v := t.Get("type").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.type = %s", v))
	}
	if v := t.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" (u.username = '%s' OR o.uid = '%s' OR o.sn='%s')", v, v, v))
	}
	where = append(where, "o.state in(2,3,4)")

	count := config.GlobalDB.JoinCount(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")})
	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{
		"u.credit", "o.*",
	}, " ORDER BY o.id desc")
	rs := make([]map[string]interface{}, 0)
	for _, item := range list {
		i := make(map[string]interface{}, 0)
		item.SetInterface(&i)
		rs = append(rs, i)
	}
	return &AdminResponse{
		State: SUCCESS,
		Data: db.DB_PARAMS{
			"list":   rs,
			"count":  count,
			"minner": s.MinnerPair("is_public = 0"),
			"state": db.DB_PARAMS{
				"2": "预约中",
				"4": "已通过",
				"3": "拒绝",
			},
		},
	}
}

func (s *SystemModel) DelMinneAccept(id int) *AdminResponse {
	if id == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请指定一个要删除的预约",
		}
	}
	_, err := config.GlobalDB.Delete(models.DB_TABLE_MINING_ACCEPT, db.DB_PARAMS{"id": id})
	if err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "删除成功",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  "删除失败",
	}
}

func (s *SystemModel) AuditAccept(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("id").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请确认要操作的信息",
		}
	}
	detail, _ := config.GlobalDB.FetchOne(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"id": t.Get("id").ToInt(), "state": 2}, db.DB_FIELDS{})
	if detail == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作的订单已被审核或不存在!",
		}
	}
	if detail.Get("dispatch_amount").ToFloat() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请先分配订单金额!",
		}
	}
	_, err := config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": t.Get("state").ToInt()}, db.DB_PARAMS{"id": t.Get("id").ToInt()})
	if err == nil {
		return &AdminResponse{
			State: SUCCESS,
			Data:  "操作信息成功!",
		}
	}
	return &AdminResponse{
		State: ERROR,
		Data:  err.Error(),
	}
}

func (s *SystemModel) OpMinnerAccept(rq P) *AdminResponse {
	t := rq.Ts()
	if t.Get("uid").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请填写一个用户!",
		}
	}
	if t.Get("amount").ToFloat() <= 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请填写预约金额!",
		}
	}

	if t.Get("dispatch_amount").ToFloat() <= 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请填写分配金额!",
		}
	}
	if t.Get("expiredtime").ToString() == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "请选择最终截止时间!",
		}
	}
	user := models.MODEL_USER.GetBaseInfo(t.Get("uid").ToInt())
	if user == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前用户不存在!",
		}
	}
	if t.Get("pid").ToInt() == 0 {
		return &AdminResponse{
			State: ERROR,
			Data:  "请选择一个矿机!",
		}
	}
	pinfo := models.MODEL_PRODUCT.GetProductInfo(t.Get("pid").ToInt())

	if pinfo == nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "矿机不存在或该矿机不支持预约!",
		}
	}
	in := db.DB_PARAMS{
		"uid":             user.Id,
		"pid":             pinfo.Id,
		"state":           2,
		"amount":          t.Get("amount").ToFloat(),
		"dispatch_amount": t.Get("dispatch_amount").ToFloat(),
		"expiredtime":     utils.TimeToint64(t.Get("expiredtime").ToString()),
		"circle":          pinfo.Circle,
		"createtime":      utils.GetNow(),
		"rate_min":        pinfo.RateMin,
		"rate":            pinfo.Rate,
		"type":            pinfo.Type,
	}
	id := t.Get("id").ToInt()
	var err error
	if id == 0 {
		in["sn"] = models.MODEL_PRODUCT.MakeSn(user.Id)
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_MINING_ORDER, in)
	} else {

		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_MINING_ORDER, in, db.DB_PARAMS{"id": id})
	}
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "操作预约失败!",
		}
	}
	return &AdminResponse{
		State: SUCCESS,
		Data:  "操作矿机成功!",
	}
}

func (s *SystemModel) Delmsg(sn_id string) *AdminResponse {
	if sn_id == "" {
		return &AdminResponse{
			State: ERROR,
			Data:  "消息不存在!",
		}
	}
	one, err := config.GlobalDB.FetchOne(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"_": fmt.Sprintf("(sn_id = '%s' OR id = '%s')", sn_id, sn_id)}, db.DB_FIELDS{})
	if err != nil {
		return &AdminResponse{
			State: ERROR,
			Data:  "当前信息不存在!",
		}
	}
	config.GlobalRedis.PushQueue(models.HASH_USER_MESSAGE, db.DB_PARAMS{"cmd": models.MESSAGE_TYPE_MSG_CANCEL, "uid": one.Get("uid").ToInt(), "content": "cancel"})

	config.GlobalDB.Delete(models.DB_TABLE_SERVICE_MESSAGE, db.DB_PARAMS{"sn_id": sn_id})

	return &AdminResponse{
		State: SUCCESS,
		Data:  "cancel成功!",
	}
}
