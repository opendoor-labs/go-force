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

// ParseSFJSON allows a slice of the expected json type to be passed
// in that will be populated with a object(s) parsed out of msg.  The
// motivation for this is that SFDC returns a single object (e.g. {})
// when only one object matches the query, but an array of objects
// (e.g. [{}, {}, ...]) when multiple objects match the query.  Instead
// of requiring client code to know that SFDC returns different types
// based on number of results, allow clients to always pass in a slice.
//
// Summary of parsing results given out object types:
//   single json object, out *slice -> *[T1]
//   json array of objects, out *slice -> *[T1, T2, ...]
//   single json object, out struct -> T
//   json array of objects, out slice -> error
func ParseSFJSON(msg []byte, out interface{}) error {
	err := forcejson.Unmarshal(msg, out)

	if err == nil {
		return nil
	}

	if !IsSlicePtr(out) {
		return errors.Wrap(err, "'out' is not a pointer to a slice")
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
