package configtypes

import (
	"fmt"
	"reflect"
	"strings"
)

// This file contains internal helpers for reflection-based functionality.

type fieldTagInfo struct {
	varName  string
	required bool
}

func getFieldTagInfo(field reflect.StructField) (fieldTagInfo, error) {
	ret := fieldTagInfo{}
	tagStr := strings.TrimSpace(field.Tag.Get("conf"))
	if tagStr == "" {
		return ret, nil
	}
	parts := strings.Split(tagStr, ",")
	ret.varName = strings.TrimSpace(parts[0])
	for i := 1; i < len(parts); i++ {
		p := strings.TrimSpace(parts[i])
		switch p {
		case "required":
			ret.required = true
		default:
			return ret, fmt.Errorf("unrecognized field tag option %q", p)
		}
	}
	return ret, nil
}

func isFieldExported(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func getReflectValueForStruct(value interface{}) (reflect.Value, bool) {
	refValue := reflect.ValueOf(value)
	if refValue.Kind() == reflect.Struct {
		if _, ok := refValue.Interface().(SingleValue); ok {
			// Our own Opt types are technically structs, but we don't want to treat them as structs
			return reflect.Value{}, false
		}
		return refValue, true
	}
	if refValue.Kind() == reflect.Ptr {
		return getReflectValueForStruct(refValue.Elem().Interface())
	}
	return reflect.Value{}, false
}

func getReflectValueForStructPtr(value interface{}) (reflect.Value, bool) {
	refValue := reflect.ValueOf(value)
	if refValue.Kind() != reflect.Ptr || refValue.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	if _, ok := refValue.Elem().Interface().(SingleValue); ok {
		// Our own Opt types are technically structs, but we don't want to treat them as structs
		return reflect.Value{}, false
	}
	return refValue.Elem(), true
}
