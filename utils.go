package formam

import (
	"encoding"
	"reflect"
	"time"
)

var (
	timeType  = reflect.TypeOf(time.Time{})
	timePType = reflect.TypeOf(&time.Time{})
)

// unmarshalText returns a boolean and error. The boolean is true if the
// value implements TextUnmarshaler, and false if not.
func unmarshalText(v reflect.Value, val string) (bool, error) {
	// check if implements the interface
	m, ok := v.Interface().(encoding.TextUnmarshaler)
	addr := v.CanAddr()
	if !ok && !addr {
		return false, nil
	} else if addr {
		return unmarshalText(v.Addr(), val)
	}
	// skip if the type is time.Time
	n := v.Type()
	if n.ConvertibleTo(timeType) || n.ConvertibleTo(timePType) {
		return false, nil
	}
	// return result
	return true, m.UnmarshalText([]byte(val))
}
