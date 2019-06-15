package server

import (
	"errors"
	"fastrpc/common"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// MakeArgs 返回Args的reflect.Value类型
func MakeArgs(request *common.Request, edcode common.EdCode, service Service) (reflect.Value, error) {
	switch edcode.(type) {
	case common.JSONEdCode:
		reqArgs := request.Args.(map[string]interface{})
		argv := reflect.New(service.ArgType)
		err := MakeArgType(reqArgs, argv)
		if err != nil {
			log.Println(err.Error())
			return reflect.New(nil), err
		}
		if argv.Kind() == reflect.Ptr {
			argv = argv.Elem()
		}
		return argv, nil
	default:
		return reflect.ValueOf(request.Args), errors.New("Unknown edcode")
	}
}

// MakeArgType 用data填充obj
func MakeArgType(data map[string]interface{}, obj reflect.Value) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetField 用map的值替换结构的值
func SetField(obj reflect.Value, name string, value interface{}) error {
	structValue := obj.Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Kind())
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

// TypeConversion 将string类型的value值转换成reflect.Value类型
func TypeConversion(value string, ntype reflect.Kind) (reflect.Value, error) {
	switch ntype {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int:
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	default:
		return reflect.ValueOf(value), errors.New("unknown type：" + ntype.String())
	}
}

