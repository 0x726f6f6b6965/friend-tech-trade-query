package helper

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

const (
	TAG_NAME = "abi"
)

func GetDataByAbi(data map[string]interface{}, ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if !v.IsValid() {
		return errors.New("the ptr is invalid")
	}

	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	} else {
		return errors.New("the input ptr is not a pointer")
	}

	var (
		crv     reflect.Value
		curLen  = 0
		element reflect.Value
	)

	if v.Kind() == reflect.Slice {
		curLen = v.Len()
		typ := v.Type().Elem()
		if typ.Kind() == reflect.Ptr {
			element = reflect.New(typ.Elem())
		} else if typ.Kind() == reflect.Struct {
			element = reflect.New(typ).Elem()
		} else {
			return errors.New("the slice ptr is not a support type")
		}
		v.Set(reflect.Append(v, element))
		temp := v.Index(curLen)
		if temp.Kind() == reflect.Ptr && !v.IsNil() {
			temp = temp.Elem()
		} else {
			return errors.New("the slice element is not a pointer")
		}
		crv = temp
	} else if v.Kind() == reflect.Struct {
		crv = v
	} else {
		return errors.New("the ptr is not a support type")
	}

	var (
		err error
	)
	for i := 0; i < crv.NumField(); i++ {
		if tName, exist := crv.Type().Field(i).Tag.Lookup(TAG_NAME); exist {
			val := data[tName]
			if crv.Field(i).Kind() != reflect.ValueOf(val).Kind() {
				val, err = matchType(crv.Field(i).Kind(), val)
				if err != nil {
					return err
				}
			}
			crv.Field(i).Set(reflect.ValueOf(val))
		}
	}
	return nil
}

func matchType(require reflect.Kind, record interface{}) (interface{}, error) {
	switch require {
	case reflect.String:
		switch reflect.TypeOf(record) {
		case reflect.TypeOf(&big.Int{}):
			// big.int
			if val, ok := record.(*big.Int); ok {
				return val.String(), nil
			}
			return nil, errors.New("parse error")
		case reflect.TypeOf(common.Address{}):
			// common.address
			if val, ok := record.(common.Address); ok {
				return val.String(), nil
			}
			return nil, errors.New("parse error")
		}
		return nil, errors.New("the require type is not support")
	}
	return nil, errors.New("the require type is not support")
}
