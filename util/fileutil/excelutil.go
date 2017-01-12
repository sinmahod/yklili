package fileutil

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"reflect"
	"strconv"
	"time"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}

type ExcelData struct {
	Name     string `tag:"姓名"`
	DayCount string `tag:"全勤工作天数"`
}

func ReadExcel2(filename string) {
	var eds []ExcelData
	ReadExcel3(filename, &eds)
}

func getFullName(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}

func ReadExcel3(filename string, obj interface{}) error {

	val := reflect.ValueOf(obj)
	ind := reflect.Indirect(val)
	fn := "" //结构体类型FullNane

	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("传递的对象必须为结构体数组")
	}

	if ind.Kind() == reflect.Slice {
		if ind.Type().Elem().Kind() != reflect.Struct {
			return fmt.Errorf("传递的对象必须为结构体数组")
		}
		fn = getFullName(ind.Type().Elem())
	} else {
		return fmt.Errorf("传递的对象必须为结构体数组")
	}

	slice := ind //复制一份结构体

	fmt.Println(fn)

	////////////////

	// v := reflect.ValueOf(obj).Elem() //结构体属性值

	// var t reflect.Type

	// fmt.Println(v.String())

	// if vType := v.Kind().String(); vType == "struct" {
	// 	t = v.Type()
	// } else {
	// 	return fmt.Errorf("传递参数必须是一个结构体[%S]", vType)
	// }

	// //对应关系设置
	// nameMap := make(map[string]string)

	// for i := 0; i < t.NumField(); i++ {
	// 	structColumn := t.Field(i)
	// 	if columnTag := structColumn.Tag.Get("tag"); columnTag != "" {
	// 		nameMap[columnTag] = structColumn.Name
	// 	} else {
	// 		nameMap[structColumn.Name] = structColumn.Name
	// 	}
	// }

	// xlFile, err := xlsx.OpenFile(filename)
	// if err != nil {
	// 	return err
	// }
	// if len(xlFile.Sheets) > 0 {
	// 	sheet := xlFile.Sheets[0]
	// 	for line, row := range sheet.Rows {
	// 		if line == 0 {
	// 			for i, cell := range row.Cells {
	// 				text, _ := cell.String()
	// 				if name, ok := nameMap[text]; ok {
	// 					delete(nameMap, text)
	// 					nameMap[strconv.Itoa(i)] = name
	// 				}
	// 			}
	// 		} else {
	// 			for i, cell := range row.Cells {
	// 				text, _ := cell.String()
	// 				if name, ok := nameMap[strconv.Itoa(i)]; ok {
	// 					//写入数据
	// 					setField(v, name, text)
	// 				}
	// 			}

	// 		}
	// 	}
	// 	fmt.Println(obj)
	// }
	return nil
}

func ReadExcel(filename string, obj interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj).Elem() //结构体属性值

	var t reflect.Type

	fmt.Println(v.String())

	if vType := v.Kind().String(); vType == "struct" {
		t = v.Type()
	} else {
		return nil, fmt.Errorf("传递参数必须是一个结构体[%S]", vType)
	}

	//对应关系设置
	nameMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		structColumn := t.Field(i)
		if columnTag := structColumn.Tag.Get("tag"); columnTag != "" {
			nameMap[columnTag] = structColumn.Name
		} else {
			nameMap[structColumn.Name] = structColumn.Name
		}
	}

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	if len(xlFile.Sheets) > 0 {
		sheet := xlFile.Sheets[0]
		for line, row := range sheet.Rows {
			if line == 0 {
				for i, cell := range row.Cells {
					text, _ := cell.String()
					if name, ok := nameMap[text]; ok {
						delete(nameMap, text)
						nameMap[strconv.Itoa(i)] = name
					}
				}
			} else {
				for i, cell := range row.Cells {
					text, _ := cell.String()
					if name, ok := nameMap[strconv.Itoa(i)]; ok {
						//写入数据
						setField(v, name, text)
					}
				}

			}
		}
		fmt.Println(obj)
	}
	return nil, nil
}

// 用map的值替换结构的值
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

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype + "\t 值：" + value)
}
