package db

/*
数据库操作扩展类 更加人性化的操作数据库
*/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type DBValues map[string]*DBValue
type DBValue struct {
	Value interface{}
}

func (dv *DBValues) GetMap() map[string]interface{} {
	rs := make(map[string]interface{})
	for k, v := range *dv {
		rs[k] = v.Value
	}
	return rs
}
func (dv *DBValues) SetObj(obj interface{}) {
	original := make(map[string]interface{}, 0)
	dataType := reflect.TypeOf(obj)
	if dataType.Kind() != reflect.Ptr {
		return
	}
	num := dataType.Elem().NumField()
	for i := 0; i < num; i++ {

		field := dataType.Elem().Field(i)
		name := field.Tag.Get("json")
		if vs, ok := (*dv)[name]; ok {
			switch field.Type.String() {
			case "int":
				original[name] = vs.ToInt()
			case "string":
				original[name] = vs.ToString()
			case "float32":
			case "float64":
				original[name] = vs.ToFloat()
			default:
				original[name] = vs.Value
			}
		}
	}
	if byt, err := json.Marshal(original); err == nil {
		json.Unmarshal(byt, obj)
	}

}

func (dv *DBValues) SetInterface(obj interface{}) {
	a := obj.(*map[string]interface{})
	for k, v := range *dv {
		(*a)[k] = v.ToString()
	}
}

func (dv *DBValues) Get(key string) *DBValue {
	rs := new(DBValue)
	if v, ok := (*dv)[key]; ok {
		return v
	}
	return rs
}

type MysqlWorker struct {
	Mysql
}

func (db *MysqlWorker) FetchOne(tablename string, condition DB_PARAMS, fields DB_FIELDS, others ...string) (DBValues, error) {
	rs, err := db.FetchAll(tablename, condition, fields, others...)
	if err == nil && len(rs) > 0 {
		return rs[0], nil
	}
	return nil, err
}
func (db *MysqlWorker) FetchAll(tablename string, condition DB_PARAMS, fields DB_FIELDS, others ...string) ([]DBValues, error) {
	field_str := "*"
	if len(fields) > 0 {
		field_str = strings.Join(fields, ",")
	}
	sql := fmt.Sprintf("select %s from %s where 1", field_str, tablename)

	if len(condition) > 0 {
		and := db.paraseConditionParams(condition)
		if and != "" {
			sql = sql + " and "
		}
		sql = sql + and
	}
	//println(fmt.Sprintf("fetchrows ====: %s", sql))
	if len(others) > 0 {
		other_str := strings.Join(others, " ")
		sql = sql + " " + other_str
	}
	rows, err := db.GetResult(sql)
	if err == nil {
		return GetDBvalues(rows), nil
	}
	return nil, err
}

func (db *MysqlWorker) JoinTable(tablename string, jointable string, on string, condition DB_PARAMS, fields DB_FIELDS, other ...string) ([]DBValues, error) {
	field_str := "*"
	if len(fields) > 0 {
		field_str = strings.Join(fields, ",")
	}
	j_str := []byte(jointable)
	check_join := strings.ToUpper(string(j_str[:12]))
	join_str := "LEFT JOIN "
	where := " 1=1"
	if strings.Contains(check_join, "JOIN") {
		join_str = ``
	}

	if len(condition) > 0 {
		where = db.paraseConditionParams(condition)
		if where != "" {
			where = "WHERE " + where
		}
	}
	sql := fmt.Sprintf("SELECT %s FROM %s %s  %s ON %s  %s", field_str, tablename, join_str, jointable, on, where)
	if len(other) > 0 {
		sql += strings.Join(other, "")
	}
	//fmt.Println(sql)
	rows, err := db.GetResult(sql)
	if err == nil {
		return GetDBvalues(rows), nil
	}
	return nil, err
}

func (db *MysqlWorker) JoinCount(tablename string, jointable string, on string, condition DB_PARAMS, other ...string) int {
	where := " 1=1"
	join_str := `LEFT JOIN`
	if strings.Contains(jointable, "JOIN") {
		join_str = ``
	}

	if len(condition) > 0 {
		where = db.paraseConditionParams(condition)
		if where != "" {
			where = "WHERE " + where
		}
	}
	ord := ""
	if len(other) > 0 {
		ord = strings.Join(other, " ")
	}
	sql := fmt.Sprintf("SELECT %s FROM %s %s  %s ON %s  %s %s", "COUNT(*) as _total", tablename, join_str, jointable, on, where, ord)
	rows, err := db.GetResult(sql)
	if err == nil {
		n, _ := strconv.Atoi(GetRecords(rows)[0]["_total"])
		return n
	}
	return 0
}

func GetDBvalues(rows *sql.Rows) []DBValues {

	rs := make([]DBValues, 0)
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
		list := make(map[string]*DBValue)
		err1 := rows.Scan(values_pointer...)
		if err1 != nil {
			continue
		}
		i = 0
		for i < n {
			list[fields[i]] = &DBValue{Value: values[i]}
			i = i + 1
		}
		rs = append(rs, list)

	}

	return rs
}

func (v *DBValue) ToString() string {
	return unescapeStr(GetValue(v.Value))
}

func (v *DBValue) ToBool() bool {
	b := v.Value.(string)

	return b == "true"
}

func (v *DBValue) ToInt() int {

	s := strings.TrimSpace(v.ToString())
	if s == "" {
		return 0
	}
	i := strings.Index(s, ".")
	if i > 0 {
		//如果被转换成了浮点则要去除浮点
		s = s[0:i]
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func (v *DBValue) ToMap() DBValues {
	rs := make(DBValues)
	s := v.ToString()
	if s == "" {
		return rs
	}
	var tmp interface{}
	err := json.Unmarshal([]byte(s), &tmp)
	if err != nil {
		return nil
	}
	t := reflect.TypeOf(tmp)
	if t.Kind() == reflect.Map {
		for key, value := range tmp.(map[string]interface{}) {
			rs[key] = &DBValue{Value: value}
		}
		return rs
	}
	return nil
}

func (v *DBValue) Get(key string, def ...interface{}) *DBValue {
	mp := v.ToMap()
	if mp == nil {
		return &DBValue{Value: def[0]}
	}
	if v, ok := mp[key]; ok {
		str := fmt.Sprintf("%v", v.Value) //强转为string
		if str != "" {
			return v
		}
	}
	if len(def) > 0 {
		return &DBValue{Value: def[0]}
	}
	return &DBValue{Value: ""}
}

func (v *DBValue) ToArray() []*DBValue {
	s := v.ToString()
	var tmp interface{}
	err := json.Unmarshal([]byte(s), &tmp)
	if err != nil {
		return nil
	}
	t := reflect.TypeOf(tmp)
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		rs := make([]*DBValue, 0)
		for _, value := range tmp.([]interface{}) {
			rs = append(rs, &DBValue{Value: value})
		}
		return rs
	}
	return nil
}

func (v *DBValue) ToFloat() float64 {
	s := v.ToString()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func (v *DBValue) GetObject(o interface{}) error {
	s := v.ToString()
	return json.Unmarshal([]byte(s), o)
}
