package models

type JqGrid struct {
	Rows    interface{} `json:"rows"`    //数据
	Page    int         `json:"page"`    //当前页码
	Total   int         `json:"total"`   //总页数
	Records int         `json:"records"` //总记录数
}
