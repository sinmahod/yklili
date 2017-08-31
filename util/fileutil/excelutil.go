package fileutil

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// type ExcelData struct {
// 	Name     string `tag:"姓名"`
// 	DayCount int    `tag:"全勤工作天数"`
// 	Sse      float32
// 	www      string
// }
//
// func WriteXLSX_Test() {
// 	eds := make([]ExcelData, 0, 0)
// 	eds = append(eds, ExcelData{"郭亮", 22, 32.123, ""})
// 	eds = append(eds, ExcelData{"郭亮2", 22, 32.123, "aa"})
// 	eds = append(eds, ExcelData{"郭亮3", 22, 32.123, "bb"})

// 	nameMap := make(map[string]string)
// 	nameMap["Name"] = "姓名"
// 	nameMap["DayCount"] = "全勤工作天数"

// 	err := WriteXLSXByMap("/workspace/asd.xlsx", eds, nameMap)
// 	fmt.Println(err)
// }
// 将结构体切片写入xlsx文件--自动提取Title = Tag > Name
func WriteXLSX(filename string, obj interface{}) error {
	return WriteXLSXByMap(filename, obj, nil)
}

// 将结构体切片写入xlsx,根据传递的map确定title以及导出的字段
func WriteXLSXByMap(filename string, obj interface{}, nameMap map[string]string) error {
	if !strings.HasSuffix(filename, ".xlsx") {
		return fmt.Errorf("文件必须是xlsx类型")
	}

	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("传递的对象必须为切片")
	}

	if val.Len() == 0 {
		return fmt.Errorf("传递的对象长度为0")
	}

	e := val.Type().Elem() //得到单个结构体type

	if e.Kind() != reflect.Struct {
		return fmt.Errorf("传递的对象必须为结构体数组指针")
	}

	if nameMap == nil {

		nameMap = make(map[string]string)

		//取得结构体tag与excel表格第一列对应，如果不写tag则取结构体变量名
		for i := 0; i < e.NumField(); i++ {
			structColumn := e.Field(i)
			if columnTag := structColumn.Tag.Get("tag"); columnTag != "" {
				//如果有Tag则将key设置为tag
				nameMap[structColumn.Name] = columnTag
			} else {
				nameMap[structColumn.Name] = structColumn.Name
			}
		}
	}

	return writexlsx(filename, nameMap, &val)
}

func writexlsx(filename string, nameMap map[string]string, val *reflect.Value) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")

	if err != nil {
		return err
	}

	row := sheet.AddRow()

	e := val.Type().Elem() //得到单个结构体type

	if e.Kind() != reflect.Struct {
		return fmt.Errorf("传递的对象必须为结构体数组指针")
	}

	nameList := make([]string, 0)

	// 写入标题(第一行)
	for i := 0; i < e.NumField(); i++ {
		columnName := e.Field(i).Name
		if title, ok := nameMap[columnName]; ok {
			cell := row.AddCell()
			cell.Value = title
			nameList = append(nameList, columnName)
		}
	}

	if len(nameList) == 0 {
		return fmt.Errorf("map与结构的列无法对应")
	}

	// 写入数据
	for i := 0; i < val.Len(); i++ {
		row := sheet.AddRow()
		idxV := val.Index(i)
		for _, columnName := range nameList {
			cell := row.AddCell()
			cell.SetValue(idxV.FieldByName(columnName))
		}
	}

	// 写入文件
	err = file.Save(filename)

	if err != nil {
		return err
	}
	return nil
}

// type ExcelData struct {
// 	Name     string `tag:"姓名"`
// 	DayCount int    `tag:"全勤工作天数"`
// 	Sse      float32
// }
//
// func ReadExcel2(filename string) {
// 	var eds []ExcelData
// 	//ReadXLSX(filename, &eds)
// 	nameMap := make(map[string]string)
// 	nameMap["姓名"] = "Name"
// 	nameMap["全勤工作天数"] = "DayCount"
// 	ReadXLSXByMap(filename, &eds, nameMap)
// 	for i, ed := range eds {
// 		fmt.Printf("%d\t%s\t%d\t%g\n", i+1, ed.Name, ed.DayCount, ed.Sse)
// 	}
// }
// 结构体tag与excel表格第一列对应，如果不写tag则取结构体变量名
func ReadXLSX(filename string, obj interface{}) error {
	return ReadXLSXByMap(filename, obj, nil)
}

// map中记录变量名与标题的对应关系 {"姓名":"Name","全勤工作天数":"DayCount"}
func ReadXLSXByMap(filename string, obj interface{}, nameMap map[string]string) error {
	if !strings.HasSuffix(filename, ".xlsx") {
		return fmt.Errorf("文件必须是xlsx类型")
	}

	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("传递的对象必须为指针")
	}

	ele := val.Elem()

	if ele.Kind() == reflect.Slice {
		if ele.Type().Elem().Kind() != reflect.Struct {
			return fmt.Errorf("传递的对象必须为结构体数组指针")
		}

		if nameMap == nil {
			e := ele.Type().Elem() //得到单个结构体type

			//保存结构体Tag或变量名与列名的对应关系
			nameMap = make(map[string]string)

			//取得结构体tag与excel表格第一列对应，如果不写tag则取结构体变量名
			for i := 0; i < e.NumField(); i++ {
				structColumn := e.Field(i)
				if columnTag := structColumn.Tag.Get("tag"); columnTag != "" {
					//如果有Tag则将key设置为tag
					nameMap[columnTag] = structColumn.Name
				} else {
					nameMap[structColumn.Name] = structColumn.Name
				}
			}
		}

		err := readxlsx(filename, nameMap, &ele)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("传递的对象必须为结构体数组指针")
	}

	if !ele.IsNil() {
		val.Elem().Set(ele)
	}
	return nil
}

func ReadXLSXToDT(filename string) ([][]string, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	dt := make([][]string, 0)

	if len(xlFile.Sheets) > 0 {
		sheet := xlFile.Sheets[0]
		for _, row := range sheet.Rows {
			dr := make([]string, 0)
			for _, cell := range row.Cells {
				text := cell.String()
				dr = append(dr, text)
			}
			dt = append(dt, dr)
		}
	}
	return dt, nil
}

func readxlsx(filename string, nameMap map[string]string, ele *reflect.Value) error {
	v := reflect.New(ele.Type().Elem()).Elem() //创建一个新的结构体value

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}

	if len(xlFile.Sheets) > 0 {
		sheet := xlFile.Sheets[0]
		for line, row := range sheet.Rows {
			if line == 0 {
				exist := false
				for i, cell := range row.Cells {
					text := cell.String()
					if name, ok := nameMap[text]; ok {
						exist = true
						delete(nameMap, text)
						nameMap[strconv.Itoa(i)] = name
					}
				}
				if !exist {
					return fmt.Errorf("Excel与结构体不存在对应列")
				}
			} else {
				for i, cell := range row.Cells {
					text := cell.String()
					if name, ok := nameMap[strconv.Itoa(i)]; ok {
						//填充结构体value
						setField(v, name, text)
					}
				}
				*ele = reflect.Append(*ele, v)
			}
		}
	}
	return nil
}

// 设置结构体的值
func setField(structValue reflect.Value, name string, value interface{}) error {
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		log.Printf("[WARN] 没有找到字段: %s \n", name)
		return nil
	}

	if !structFieldValue.CanSet() {
		log.Printf("[WARN] 字段类型不可被修改：%s \n", name)
		return nil
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		valStr := fmt.Sprintf("%v", value)
		ntype := structFieldValue.Type().Name()
		//如果是空字符串则只能Set到string类型，其他类型跳过
		if valStr == "" && ntype != "string" {
			return nil
		}
		val, err = typeConversion(valStr, ntype) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "bool" {
		i, err := strconv.ParseBool(value)
		return reflect.ValueOf(bool(i)), err
	}

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype + "\t 值：" + value)
}
