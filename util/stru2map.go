package util

import (
	"context"
	"reflect"

	"github.com/lanwenhong/lgobase/logger"
)

func getStructType(struc interface{}) reflect.Type {
	sType := reflect.TypeOf(struc)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}
	return sType
}

func Stru2Map(ctx context.Context, struc interface{}) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	sType := getStructType(struc)
	structVal := reflect.ValueOf(struc)
	for i := 0; i < sType.NumField(); i++ {
		if !structVal.Field(i).IsZero() {
			structFieldName := sType.Field(i).Name
			tag_name := sType.Field(i).Tag.Get("form")
			logger.Debugf(ctx, "tag_name: %s", tag_name)
			ret[tag_name] = structVal.FieldByName(structFieldName).Interface()
		}
	}
	return ret, nil
}
