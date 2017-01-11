package models

import (
	"reflect"
)

type DataGrid struct {
	Data      interface{} `json:"rows"`    //数据
	PageIndex int         `json:"page"`    //当前页码
	PageTotal int         `json:"total"`   //总页数
	DataTotal int64       `json:"records"` //总记录数
}

type dTable struct {
	Key   string
	Value interface{}
}

/**
*    将内容封装到可以解析的类型DataGrid中
**/
func GetDataGrid(data interface{}, pageindex, pagetotal int, datatotal int64) *DataGrid {
	val := reflect.ValueOf(data)

	if val.Type().String() == "map[string]interface {}" || val.Type().String() == "map[string]string" {
		dt := make([]*dTable, 0)
		for k, v := range data.(map[string]interface{}) {
			dt = append(dt, &dTable{k, v})
		}
		data = dt
	}

	return &DataGrid{Data: data, PageIndex: pageindex, PageTotal: pagetotal, DataTotal: datatotal}
}
