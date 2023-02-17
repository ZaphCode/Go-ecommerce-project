package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func ItemInSlice[T comparable](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PTR[T any](v T) *T {
	return &v
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	)

	s := make([]rune, n)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func GetStructFields(s interface{}) ([]string, error) {
	if s == nil {
		return nil, fmt.Errorf("only accepts not nil *structs")
	}

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only accepts structs")
	}

	fields := make([]string, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, v.Type().Field(i).Name)
	}
	return fields, nil
}

func GetStructField(strc interface{}, fldName string) (interface{}, error) {
	v := reflect.ValueOf(strc)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	field := v.FieldByName(fldName)

	if !field.IsValid() {
		return nil, fmt.Errorf("field %q not found in struct", fldName)
	}

	if !field.CanInterface() {
		return nil, fmt.Errorf("field %q is unexported and cannot be accessed", fldName)
	}

	return field.Interface(), nil
}

func SetStructField(strc interface{}, fldName string, val interface{}) error {
	v := reflect.ValueOf(strc)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	fld := v.FieldByName(fldName)

	if !fld.IsValid() {
		return fmt.Errorf("field %q not found in struct", fldName)
	}

	if !fld.CanSet() {
		return fmt.Errorf("field %q is unexported or cannot be set", fldName)
	}

	if !fld.Type().AssignableTo(reflect.TypeOf(val)) {
		return fmt.Errorf("value of type %T cannot be assigned to field %q of type %v", val, fldName, fld.Type())
	}

	fld.Set(reflect.ValueOf(val))

	return nil
}

func IsSameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func UpdateStructFields(strc interface{}, uf map[string]interface{}) error {
	v := reflect.ValueOf(strc)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	for fldName, val := range uf {
		fv, err := GetStructField(strc, fldName)

		if err != nil {
			return err
		}

		if !IsSameType(fv, val) {
			return fmt.Errorf("value of type %T cannot be assigned to field %q of type %T", val, fldName, fv)
		}
	}

	for fldName, val := range uf {
		if err := SetStructField(strc, fldName, val); err != nil {
			return err
		}
	}

	return nil
}

func IsZeroValue(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func StructToMap(strc interface{}) map[string]any {
	v := reflect.ValueOf(strc)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	fldMap := make(map[string]any)

	for i := 0; i < v.NumField(); i++ {
		fldName := v.Type().Field(i).Name
		fldVal := v.Field(i)

		if !fldVal.IsZero() {
			if fldVal.Kind() == reflect.Ptr {
				fldVal = fldVal.Elem()
			}
			fldMap[fldName] = fldVal.Interface()
		}
	}

	return fldMap
}
