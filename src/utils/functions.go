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

func GetStructAttr(strc interface{}, fieldName string) (reflect.Value, error) {
	zero := reflect.Value{}

	if strc == nil {
		return zero, fmt.Errorf("only accepts not nil *structs")
	}

	v := reflect.ValueOf(strc)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return zero, fmt.Errorf("only accepts structs")
	}

	fld := v.FieldByName(fieldName)

	if !fld.IsValid() {
		return zero, fmt.Errorf("field does'nt exist")
	}

	return fld, nil
}
