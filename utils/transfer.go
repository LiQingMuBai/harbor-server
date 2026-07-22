package utils

import (
	"cointrade/lib/db"
	"encoding/json"
	"reflect"
)

var Transfer TransferValue

type TransferValue struct {
	Value    map[string]*db.DBValue
	OldValue interface{}
}

func (v *TransferValue) TransferByDbValue(data interface{}) map[string]*db.DBValue {
	v.OldValue = data
	rs := make(map[string]*db.DBValue, 0)
	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Map {
		str, err := json.Marshal(data)
		if err != nil {
			return nil
		}
		d := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(str), &d)

		for k, item := range d {
			isbool := false
			vi := reflect.TypeOf(item)

			switch vi.Kind() {
			case reflect.String:

				if item.(string) != "" {
					isbool = true
				}
			case reflect.Float64:
				//utils.Log("transfer, int ", k+" : ", reflect.ValueOf(item).IsNil())
				//if item.(float64) != 0 {
				isbool = true
				//}
			case reflect.Int:

				if item.(int64) != 0 {
					isbool = true
				}
			case reflect.Slice:

				if len(item.([]interface{})) != 0 {
					isbool = true
				}
			default:
				Log("transfer, item ", item)
			}
			if isbool {
				rs[k] = &db.DBValue{Value: item}
			}

		}

	} else if t.Kind() == reflect.String {
		rs["_"] = &db.DBValue{Value: data}
	}
	v.Value = rs
	return rs
}

func (v *TransferValue) TransferToOrignalByDbValues(data []db.DBValues) []map[string]string {
	rs := make([]map[string]string, 0)
	for i := 0; i < len(data); i++ {
		dbvalue := make(map[string]string)
		for k, v := range data[i] {
			dbvalue[k] = v.ToString()
		}
		rs = append(rs, dbvalue)
	}
	return rs
}
func (v *TransferValue) GetMapByDbValue(key string) *db.DBValue {

	if _, ok := v.Value[key]; ok {
		return v.Value[key]
	}
	return &db.DBValue{Value: ""}
}
