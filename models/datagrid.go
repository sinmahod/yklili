package models

type DataGrid struct {
	Data      interface{} `json:"rows"`    //数据
	PageIndex int         `json:"page"`    //当前页码
	PageTotal int         `json:"total"`   //总页数
	DataTotal int64       `json:"records"` //总记录数
}

/**
*    将内容封装到可以解析的类型DataGrid中
**/
func GetDataGrid(data interface{}, pageindex, pagetotal int, datatotal int64) *DataGrid {
	return &DataGrid{Data: data, PageIndex: pageindex, PageTotal: pagetotal, DataTotal: datatotal}
}
