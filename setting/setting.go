package setting

import (
	"reflect"
	"strings"

	"github.com/bionicosmos/aegle/config"
	"github.com/bionicosmos/aegle/model"
)

func Get[T any](fieldPath string) (t T) {
	setting, err := model.FindSetting(fieldPath)
	if err != nil {
		return
	}
	return find[T](fieldPath, []any{setting, config.C})
}

func find[T any](fieldPath string, structs []any) (t T) {
	for _, s := range structs {
		field := findField(fieldPath, reflect.ValueOf(s))
		if field != nil {
			return field.(T)
		}
	}
	return
}

func findField(fieldPath string, value reflect.Value) any {
	i := strings.IndexByte(fieldPath, '.')
	if i == -1 {
		return getFieldValue(
			value.FieldByName(fieldPath),
			func(value reflect.Value) any {
				return value.Interface()
			},
		)
	}
	return getFieldValue(
		value.FieldByName(fieldPath[:i]),
		func(value reflect.Value) any {
			return findField(fieldPath[i+1:], value)
		},
	)
}

func getFieldValue(value reflect.Value, f func(value reflect.Value) any) any {
	if value.Kind() == reflect.Invalid {
		return nil
	}
	return f(value)
}
