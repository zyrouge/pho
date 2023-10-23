package utils

import "reflect"

func SetStructField(target any, field string, to any) bool {
	rTarget := reflect.ValueOf(target).Elem()
	rTargetType := rTarget.Type()
	rField, ok := rTargetType.FieldByName(field)
	if !ok {
		return false
	}
	rTo := reflect.ValueOf(to)
	rToType := rTo.Type()
	if rField.Type == rToType {
		rTarget.FieldByIndex(rField.Index).Set(rTo.Convert(rField.Type))
	}
	return true
}
