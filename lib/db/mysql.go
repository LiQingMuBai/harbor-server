package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DEBUG = 1 //开启DEBUG  模式
)

type DB_PARAMS map[string]interface{}
type DB_FIELDS []string
type DB_LIST_RESULT []map[string]string
type DB_ROW_RESULT map[string]string
type Mysql struct {
	host   string
	port   int
	user   string
	pass   string
	dbname string
	db     *sql.DB
}

func (conn *Mysql) SetLinkInfo(host string, port int, user string, pass string, dbname string) {
	conn.host = host
	conn.port = port
	conn.pass = pass
	conn.user = user
	conn.dbname = dbname
}
func (conn *Mysql) Connect() error {
	sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conn.user, conn.pass, conn.host, conn.port, conn.dbname)
	//println(sqlstr)
	db, err := sql.Open("mysql", sqlstr)
	db.SetMaxIdleConns(50)
	conn.db = db
	return err
}
func (conn *Mysql) SetNames() error {
	_, err := conn.db.Exec("set names 'utf8'")
	return err
}
func (conn *Mysql) Execute(sql string) (sql.Result, error) {
	if DEBUG == 1 {
		log.Printf("mysql exec sql: %s", sql)
	}
	//println(sql)
	return conn.db.Exec(sql)
}
func (conn *Mysql) GetResult(sql string) (*sql.Rows, error) {
	//println(sql)
	//if DEBUG == 1 {
	//	fmt.Println(sql)
	//}
	return conn.db.Query(sql)
}

func (conn *Mysql) Close() {
	conn.db.Close()
}
func (conn *Mysql) paraseInsertData(params map[string]interface{}) string {
	keystr := ""
	valuestr := ""
	for k, v := range params {
		keystr = keystr + "`" + k + "`,"
		valuestr = valuestr + "'" + formatValue(v) + "',"
	}
	keystr = keystr[0 : len(keystr)-1]
	valuestr = valuestr[0 : len(valuestr)-1]
	str := fmt.Sprintf("(%s)values(%s)", keystr, valuestr)
	return str
}
func (conn *Mysql) InsertData(tablename string, params map[string]interface{}) (int64, error) {
	partstr := conn.paraseInsertData(params)
	sql := fmt.Sprintf("insert into %s %s", tablename, partstr)
	//println(sql)
	rs, err := conn.Execute(sql)
	if err == nil {
		return rs.LastInsertId()
	} else {
		return 0, err
	}
}
func (conn *Mysql) paraseUpdateParams(params map[string]interface{}) string {

	str := ""
	for k, v := range params {
		str = str + fmt.Sprintf("`%s`='%s',", k, formatValue(v))
	}
	str = str[0 : len(str)-1]
	return str
}
func (conn *Mysql) UpdateData(tablename string, params map[string]interface{}, condition map[string]interface{}) (sql.Result, error) {
	partstr := conn.paraseUpdateParams(params)
	sql := fmt.Sprintf("update %s set %s where 1", tablename, partstr)
	if len(condition) > 0 {
		sql = fmt.Sprintf("update %s set %s where %s", tablename, partstr, conn.paraseConditionParams(condition))
	}
	log.Printf("mysql update sql: %s", sql)
	return conn.Execute(sql)
}
func (conn *Mysql) Delete(tablename string, condition map[string]interface{}) (sql.Result, error) {
	sql := fmt.Sprintf("delete from %s where %s", tablename, conn.paraseConditionParams(condition))
	return conn.Execute(sql)
}

func (conn *Mysql) FetchRowsByJoin(tablename string, feild []string, join []string, condition map[string]interface{}, others ...string) ([]map[string]string, error) {
	field_str := "*"
	if len(feild) > 0 {
		field_str = strings.Join(feild, ",")
	}
	sql := fmt.Sprintf("select %s from %s %s WHERE ", field_str, tablename, strings.Join(join, ","))

	if len(condition) > 0 {
		sql = sql + "  " + conn.paraseConditionParams(condition)
	}
	if len(others) > 0 {
		other_str := strings.Join(others, " ")
		sql = sql + " " + other_str
	}
	//println(sql)
	rows, err := conn.GetResult(sql)
	if err == nil {
		return GetRecords(rows), nil
	}
	return []map[string]string{}, err
}
func (conn *Mysql) FetchRows(tablename string, condition map[string]interface{}, fields []string, others ...string) ([]map[string]string, error) {
	field_str := "*"
	if len(fields) > 0 {
		field_str = strings.Join(fields, ",")
	}
	sql := fmt.Sprintf("select %s from %s where 1", field_str, tablename)

	if len(condition) > 0 {
		parsestr := conn.paraseConditionParams(condition)
		and := ""
		if parsestr != "" {
			and = " and "
		}
		sql = sql + and + conn.paraseConditionParams(condition)
	}
	if len(others) > 0 {
		other_str := strings.Join(others, " ")
		sql = sql + " " + other_str
	}
	rows, err := conn.GetResult(sql)
	if err == nil {
		return GetRecords(rows), nil
	}
	return []map[string]string{}, err

}
func (conn *Mysql) FetchRow(tablename string, condition map[string]interface{}, fields []string, others ...string) (map[string]string, error) {
	rs, err := conn.PageRows(tablename, condition, fields, 0, 1, others...)

	if err == nil {
		if len(rs) > 0 {
			return rs[0], nil
		} else {
			return nil, nil
		}

	}
	return nil, err
}
func (conn *Mysql) paraseConditionParams(params map[string]interface{}) string {

	strlist := make([]string, 0)
	for k, v := range params {
		if k != "_" {
			//fmt.Println("formatValue(v)", formatValue(v))
			//str = str + fmt.Sprintf("`%s`='%s and '", k, formatValue(v))
			strlist = append(strlist, fmt.Sprintf("`%s`='%s'", k, formatValue(v)))
		} else {
			strlist = append(strlist, unescapeStr(formatValue(v)))

		}
	}
	return strings.Join(strlist, " and ")
}
func (conn *Mysql) GetCount(tablename string, condition map[string]interface{}, other ...string) int {
	sql := "select count(*) as _total from " + tablename + " where 1"
	condition_str := conn.paraseConditionParams(condition)
	if condition_str != "" {
		sql = sql + " and " + condition_str
	}
	if len(other) > 0 {
		sql += strings.Join(other, " ")
	}
	rows, err := conn.GetResult(sql)
	if err == nil {
		record := GetRecords(rows)
		if len(record) == 0 {
			return 0
		}
		log.Printf("mysql count record: %+v", record[0])
		n, _ := strconv.Atoi(record[0]["_total"])
		return n
	}
	return 0
}
func (conn *Mysql) PageRows(tablename string, condition map[string]interface{}, fields []string, offset int, limit int, others ...string) ([]map[string]string, error) {
	//分页数据获取
	field_str := "*"
	if len(fields) > 0 {
		field_str = strings.Join(fields, ",")
	}
	sql := fmt.Sprintf("select %s from %s where 1", field_str, tablename)
	if len(condition) > 0 {
		sql = sql + " and " + conn.paraseConditionParams(condition)
	}
	if len(others) > 0 {
		other_str := strings.Join(others, " ")
		sql = sql + " " + other_str
	}
	sql = sql + fmt.Sprintf(" limit %d,%d", offset, limit)
	rows, err := conn.GetResult(sql)
	if err == nil {
		return GetRecords(rows), nil
	}
	return []map[string]string{}, err
}
func (conn *Mysql) AddValue(tablename string, addvalues map[string]float64, condition DB_PARAMS) error { //按给定的键值增加数值
	temp := make([]string, 0)
	for k, v := range addvalues {
		flag := "+"
		if v < 0 {
			flag = "-"
		}
		temp = append(temp, fmt.Sprintf("`%s`=`%s`%s%f", k, k, flag, math.Abs(v)))
	}
	sql := fmt.Sprintf("update %s set %s where 1", tablename, strings.Join(temp, ","))
	if len(condition) > 0 {
		sql = sql + " and " + conn.paraseConditionParams(condition)
	}
	_, err := conn.Execute(sql)
	return err
}
func GetRecords(rows *sql.Rows) []map[string]string {

	rs := make([]map[string]string, 0)
	if rows == nil {
		return rs
	}
	fields, err := rows.Columns()
	if err != nil {
		return rs
	}
	n := len(fields)
	i := 0
	values := make([]interface{}, n)
	values_pointer := make([]interface{}, n)
	for index, _ := range values {
		values_pointer[index] = &values[index]
	}
	for rows.Next() {
		list := make(map[string]string)
		err1 := rows.Scan(values_pointer...)
		if err1 != nil {
			continue
		}
		i = 0
		for i < n {
			list[fields[i]] = unescapeStr(formatValue(values[i]))
			i = i + 1
		}
		rs = append(rs, list)

	}

	return rs
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case int:
		//utils.Log("dsadsadasdas111111111111")
		return strconv.Itoa(v.(int))
	case int32:
		//utils.Log("dsadsadasdas222222222222222222")
		return fmt.Sprintf("%d", v.(int32))
	case int64:
		//utils.Log("dsadsadas333333333333333")
		return fmt.Sprintf("%d", v.(int64))
	case float32:
		//utils.Log("dsadsadasdas444444444444")
		return fmt.Sprintf("%.8f", v.(float64))
	case float64:
		//utils.Log("dsadsadasda5555555555555")
		return fmt.Sprintf("%.8f", v.(float64))

	case string:
		return escapeStr(v.(string))
	case []byte:
		//utils.Log("dsadsadasda666666666")
		return string(v.([]byte))
	}
	s, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(s)
}
func escapeStr(v string) string {
	return strings.Replace(v, "'", "''", -1)
}
func unescapeStr(v string) string {
	return strings.Replace(v, "''", "'", -1)
}
func GetValue(v interface{}) string {
	return formatValue(v)
}
