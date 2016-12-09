package modelutil

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

//用map填充结构
func FillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
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
