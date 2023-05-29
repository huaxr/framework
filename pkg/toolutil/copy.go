// Author: huaxr
// Time:   2021/7/19 上午11:44
// Git:    huaxr

package toolutil

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

// shallow copy
func CopyShallowStruct(src, dst interface{}) {
	sVal := reflect.ValueOf(src).Elem()
	dVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < sVal.NumField(); i++ {
		value := sVal.Field(i)
		name := sVal.Type().Field(i).Name

		dValue := dVal.FieldByName(name)
		if dValue.IsValid() == false {
			continue
		}
		dValue.Set(value)
	}
}

// deep copy
func CopyDeepStruct(src, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &dst)
}

func Clone(a, b interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(a); err != nil {
		return err
	}
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}
