package parser

import (
	"reflect"

	"github.com/opendoor-labs/go-force/forcejson"
	"github.com/pkg/errors"
)

func IsSlicePtr(i interface{}) (isptr bool) {
	defer func() {
		recover()
	}()

	slice := reflect.ValueOf(i).Elem()

	// panic if slice doesn't behave like a slice
	slice.Slice(0, 0)

	return true
}

func ParseSFJSON(msg []byte, out interface{}) error {
	err := forcejson.Unmarshal(msg, out)

	if err == nil {
		return nil
	}

	if !IsSlicePtr(out) {
		return err
	}

	// pointer to the slice
	slice := reflect.ValueOf(out).Elem()

	slice.Set(reflect.MakeSlice(slice.Type(), 0, 1))

	t := slice.Type().Elem()
	i := reflect.New(t).Interface()

	err = forcejson.Unmarshal(msg, &i)

	v := reflect.ValueOf(i)

	if err != nil {
		return errors.Wrap(err, "Unmarshal single object failed")
	}

	slice.Set(reflect.Append(slice, v.Elem()))
	return nil
}
