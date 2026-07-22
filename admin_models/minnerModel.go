package adminmodels

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type MinnerModel struct{}

/**
 *	矿机列表
 */
func (m *MinnerModel) MinnerList(rq P) *AdminResponse {
	t := rq.Ts()

	count := config.GlobalDB.GetCount(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{})

	list, _ := config.GlobalDB.FetchAll(models.DB_TABLE_MINING_PRODUCT, db.DB_PARAMS{}, db.DB_FIELDS{}, utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))

	minnerlist := make([]*models.ProductInfo, 0)
	for _, v := range list {
		minner := new(models.ProductInfo)
		v.SetObj(minner)
		minnerlist = append(minnerlist, minner)
	}

	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":  minnerlist,
			"count": count,
		},
	}
}

/**
 * 操作矿机
 */
func (m *MinnerModel) OpMinner(rq *models.ProductInfo) *AdminResponse {
	rs := new(AdminResponse)
	if rq.Name == "" {
		rs.State = PARAM_ERROR
		rs.Data = "矿机名称不能为空!"
		return rs
	}
	if rq.Circle == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "矿机周期必填"
		return rs
	}
	if rq.Rate == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "用户收益比例必填"
		return rs
	}
	if rq.Profile == nil {
		rs.State = PARAM_ERROR
		rs.Data = "矿机基本信息配置不能为空!"
		return rs
	}
	if rq.Price == 0 {
		rs.State = PARAM_ERROR
		rs.Data = "矿机价格不能为空!"
		return rs
	}
	insert := make(P)
	if bytes, err := json.Marshal(rq); err == nil {
		json.Unmarshal(bytes, &insert)
	} else {
		rs.State = PARAM_ERROR
		rs.Data = "解析矿机失败!"
		return rs
	}
	var err error
	if rq.Id == 0 {
		_, err = config.GlobalDB.InsertData(models.DB_TABLE_MINING_PRODUCT, insert)
	} else {
		_, err = config.GlobalDB.UpdateData(models.DB_TABLE_MINING_PRODUCT, insert, db.DB_PARAMS{"id": rq.Id})
	}
	if err != nil {
		rs.State = ERROR
		rs.Data = "操作矿机信息失败!"
	} else {
		config.GlobalRedis.PushQueue(models.QUEUE_RPC_LIST, &models.RpcRequest{Cmd: models.SYSTEM_RELOAD_MINPRODUCT_LIST})
		rs.State = SUCCESS
		rs.Data = "操作成功!"
	}
	return rs
}

/**
 *	用户订单矿机列表
 */
func (m *MinnerModel) UserMinnerOrderList(rq P) *AdminResponse {
	where := make([]string, 0)
	t := rq.Ts()
	if v := t.Get("sn").ToString(); v != "" {
		where = append(where, fmt.Sprintf(" o.sn like '%%%s%%'", v))
	}
	if v := t.Get("email").ToString(); v != "" {
		where = append(where, fmt.Sprintf("u.username like '%%%s%%'", v))
	}
	if v := t.Get("pid").ToInt(); v > 0 {
		where = append(where, fmt.Sprintf(" o.pid = '%d'", v))
	}
	if v := t.Get("createtime").ToArray(); len(v) > 0 {
		where = append(where, fmt.Sprintf("o.createtime between %d and %d", utils.TimeToint64(v[0].ToString()), utils.TimeToint64(v[1].ToString())))
	}

	list, _ := config.GlobalDB.JoinTable(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")}, db.DB_FIELDS{"o.*", "u.username"}, utils.Order(t.Get("sort").ToString()), utils.Limit(t.Get("page").ToInt(), t.Get("limit").ToInt()))
	count := config.GlobalDB.JoinCount(models.DB_TABLE_MINING_ORDER+" as o", models.DB_TABLE_USER+" as u", "o.uid = u.id", db.DB_PARAMS{"_": strings.Join(where, " AND ")})

	return &AdminResponse{
		State: SUCCESS,
		Data: P{
			"list":   list,
			"count":  count,
			"p_list": SYSTEM_MODEL.MinnerPair(),
		},
	}
}
