// Package formam implements functions to decode values of a html form.
package formam

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const tagName = "formam"

// A pathMap holds the values of a map with its key and values correspondent
type pathMap struct {
	m reflect.Value

	key   string
	value reflect.Value

	path string
}

// a pathMaps holds the values for each key
type pathMaps []*pathMap

// find find and get the value by the given key
func (ma pathMaps) find(id reflect.Value, key string) *pathMap {
	for _, v := range ma {
		if v.m == id && v.key == key {
			return v
		}
	}
	return nil
}

// A decoder holds the values from form, the 'reflect' value of main struct
// and the 'reflect' value of current path
type decoder struct {
	main reflect.Value

	curr   reflect.Value
	value  string
	values []string

	path  string
	field string
	index int

	maps pathMaps
}

// Decode decodes the url.Values into a element that must be a pointer to a type provided by argument
func Decode(vs url.Values, dst interface{}) error {
	main := reflect.ValueOf(dst)
	if main.Kind() != reflect.Ptr {
		return fmt.Errorf("formam: the value passed for decode is not a pointer but a %v", main.Kind())
	}
	dec := &decoder{main: main.Elem()}

	// iterate over the form's values and decode it
	for k, v := range vs {
		dec.path = k
		dec.field = k
		dec.values = v
		dec.value = v[0]
		if dec.value != "" {
			if err := dec.begin(); err != nil {
				return err
			}
		}
	}

	// set values of each maps
	for _, v := range dec.maps {
		typ := v.m.Type()
		key := typ.Key()

		// check if the key implements the UnmarshalText interface
		var val reflect.Value
		if key.Kind() == reflect.Ptr {
			val = reflect.New(key.Elem())
		} else {
			val = reflect.New(key).Elem()
		}
		dec.value = v.key
		ok, err := dec.unmarshalText(val)
		if ok {
			v.m.SetMapIndex(val, v.value)
			continue
		}
		if err != nil {
			return err
		}
		// if the key not implements the UnmarahalText interface then...
		switch key.Kind() {
		case reflect.String:
			v.m.SetMapIndex(reflect.ValueOf(v.key), v.value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(v.key, 10, 64)
			if err != nil {
				return fmt.Errorf("formam: the key's value \"%v\" in path \"%v\" is not a valid signed integer number", v.key, v.path)
			}
			switch key.Kind() {
			case reflect.Int:
				v.m.SetMapIndex(reflect.ValueOf(int(num)), v.value)
			case reflect.Int8:
				v.m.SetMapIndex(reflect.ValueOf(int8(num)), v.value)
			case reflect.Int16:
				v.m.SetMapIndex(reflect.ValueOf(int16(num)), v.value)
			case reflect.Int32:
				v.m.SetMapIndex(reflect.ValueOf(int32(num)), v.value)
			case reflect.Int64:
				v.m.SetMapIndex(reflect.ValueOf(int64(num)), v.value)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			num, err := strconv.ParseUint(v.key, 10, 64)
			if err != nil {
				return fmt.Errorf("formam: the key's value \"%v\" in path \"%v\" is not a valid unsigned integer number", v.key, v.path)
			}
			switch key.Kind() {
			case reflect.Uint:
				v.m.SetMapIndex(reflect.ValueOf(uint(num)), v.value)
			case reflect.Uint8:
				v.m.SetMapIndex(reflect.ValueOf(uint8(num)), v.value)
			case reflect.Uint16:
				v.m.SetMapIndex(reflect.ValueOf(uint16(num)), v.value)
			case reflect.Uint32:
				v.m.SetMapIndex(reflect.ValueOf(uint32(num)), v.value)
			case reflect.Uint64:
				v.m.SetMapIndex(reflect.ValueOf(uint64(num)), v.value)
			}
			/*
				case reflect.Float32, reflect.Float64:
					dec.value = strings.Replace(dec.value, ",", ".", -1)
					num, err := strconv.ParseFloat(dec.value, key.Bits())
					if err != nil {
						return fmt.Errorf("formam: the key's value \"%v\" in path \"%v\" is not a valid float number", v.key, v.path)
					}
					switch key.Kind() {
					case reflect.Float32:
						v.m.SetMapIndex(reflect.ValueOf(float32(num)), v.value)
					case reflect.Float64:
						v.m.SetMapIndex(reflect.ValueOf(float64(num)), v.value)
					}
			*/
		case reflect.Bool:
			switch dec.value {
			case "true", "on", "1":
				v.m.SetMapIndex(reflect.ValueOf(true), v.value)
			case "false", "off", "0":
				v.m.SetMapIndex(reflect.ValueOf(false), v.value)
			default:
				return fmt.Errorf("formam: the key's value \"%v\" in path \"%v\" is not a valid boolean", v.value, v.path)
			}
		case reflect.Array:
			return fmt.Errorf("formam: the type Array is not implemented yet for key in maps...")
		default:
			return fmt.Errorf("formam: the key's type \"%v\" in path \"%v\" is not a valid type", key.Kind().String(), v.path)
		}
	}

	dec.maps = make(pathMaps, 0)
	return nil
}

// begin prepare the current path to walk through it
func (dec *decoder) begin() (err error) {
	dec.curr = dec.main
	fields := strings.Split(dec.field, ".")
	for i, field := range fields {
		b := strings.IndexAny(field, "[")
		if b != -1 {
			// is a array
			e := strings.IndexAny(field, "]")
			if e == -1 {
				return fmt.Errorf("formam: bad syntax array in the field \"%v\" of path \"%v\"", dec.field, dec.path)
			}
			dec.field = field[:b]
			if dec.index, err = strconv.Atoi(field[b+1 : e]); err != nil {
				return fmt.Errorf("formam: the index of array is not a number in the field \"%v\" of path \"%v\"", dec.field, dec.path)
			}
			if len(fields) == i+1 {
				return dec.end()
			}
			if err = dec.walk(); err != nil {
				return
			}
		} else {
			// not is a array
			dec.field = field
			dec.index = -1
			if len(fields) == i+1 {
				return dec.end()
			}
			if err = dec.walk(); err != nil {
				return
			}
		}
	}
	return
}

// walk traverses the current path until to the last field
func (dec *decoder) walk() error {
	// check if is a struct or map
	switch dec.curr.Kind() {
	case reflect.Struct:
		if err := dec.findStructField(); err != nil {
			return err
		}
	case reflect.Map:
		dec.currentMap()
	}
	// check if the struct or map is a interface
	if dec.curr.Kind() == reflect.Interface {
		dec.curr = dec.curr.Elem()
	}
	// check if the struct or map is a pointer
	if dec.curr.Kind() == reflect.Ptr {
		if dec.curr.IsNil() {
			dec.curr.Set(reflect.New(dec.curr.Type().Elem()))
		}
		dec.curr = dec.curr.Elem()
	}
	// finally, check if there are access to slice/array or not...
	if dec.index != -1 {
		switch dec.curr.Kind() {
		case reflect.Slice, reflect.Array:
			if dec.curr.Len() <= dec.index {
				dec.expandSlice(dec.index + 1)
			}
			dec.curr = dec.curr.Index(dec.index)
		default:
			return fmt.Errorf("formam: the field \"%v\" in path \"%v\" has a index for array but it is not", dec.field, dec.path)
		}
	}
	return nil
}

// end finds the last field for decode its value correspondent
func (dec *decoder) end() error {
	if dec.curr.Kind() == reflect.Struct {
		if err := dec.findStructField(); err != nil {
			return err
		}
	}
	if dec.value == "" {
		return nil
	}
	return dec.decode()
}

// decode sets the value in the last field found by end function
func (dec *decoder) decode() error {
	if ok, err := dec.unmarshalText(dec.curr); ok || err != nil {
		return err
	}

	switch dec.curr.Kind() {
	case reflect.Map:
		dec.currentMap()
		return dec.decode()
	case reflect.Slice, reflect.Array:
		if dec.index == -1 {
			// not has index, so to decode all values in the slice/array
			dec.expandSlice(len(dec.values))
			tmp := dec.curr
			for i, v := range dec.values {
				dec.curr = tmp.Index(i)
				dec.value = v
				if err := dec.decode(); err != nil {
					return err
				}
			}
		} else {
			// has index, so to decode value by index indicated
			if dec.curr.Len() <= dec.index {
				dec.expandSlice(dec.index + 1)
			}
			dec.curr = dec.curr.Index(dec.index)
			return dec.decode()
		}
	case reflect.String:
		dec.curr.SetString(dec.value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, err := strconv.ParseInt(dec.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" should be a valid signed integer number", dec.field, dec.path)
		} else {
			dec.curr.SetInt(num)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if num, err := strconv.ParseUint(dec.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" should be a valid unsigned integer number", dec.field, dec.path)
		} else {
			dec.curr.SetUint(num)
		}
	case reflect.Float32, reflect.Float64:
		if num, err := strconv.ParseFloat(dec.value, dec.curr.Type().Bits()); err != nil {
			return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" should be a valid float number", dec.field, dec.path)
		} else {
			dec.curr.SetFloat(num)
		}
	case reflect.Bool:
		switch dec.value {
		case "true", "on", "1":
			dec.curr.SetBool(true)
		case "false", "off", "0":
			dec.curr.SetBool(false)
		default:
			return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" is not a valid boolean", dec.field, dec.path)
		}
	case reflect.Interface:
		dec.curr.Set(reflect.ValueOf(dec.value))
	case reflect.Ptr:
		dec.curr.Set(reflect.New(dec.curr.Type().Elem()))
		dec.curr = dec.curr.Elem()
		return dec.decode()
	case reflect.Struct:
		switch dec.curr.Interface().(type) {
		case time.Time:
			t, err := time.Parse("2006-01-02", dec.value)
			if err != nil {
				return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" is not a valid datetime", dec.field, dec.path)
			}
			dec.curr.Set(reflect.ValueOf(t))
		case url.URL:
			u, err := url.Parse(dec.value)
			if err != nil {
				return fmt.Errorf("formam: the value of field \"%v\" in path \"%v\" is not a valid url", dec.field, dec.path)
			}
			dec.curr.Set(reflect.ValueOf(*u))
		default:
			return fmt.Errorf("formam: not supported type for field \"%v\" in path \"%v\"", dec.field, dec.path)
		}
	default:
		return fmt.Errorf("formam: not supported type for field \"%v\" in path \"%v\"", dec.field, dec.path)
	}

	return nil
}

// findField finds a field by its name, if it is not found,
// then retry the search examining the tag "formam" of every field of struct
func (dec *decoder) findStructField() error {
	var anon reflect.Value

	num := dec.curr.NumField()
	for i := 0; i < num; i++ {
		field := dec.curr.Type().Field(i)
		if field.Name == dec.field {
			// check if the field's name is equal
			dec.curr = dec.curr.Field(i)
			return nil
		} else if field.Anonymous {
			// if the field is a anonymous struct, then iterate over its fields
			tmp := dec.curr
			dec.curr = dec.curr.FieldByIndex(field.Index)
			if err := dec.findStructField(); err != nil {
				dec.curr = tmp
				continue
			}
			// field in anonymous struct is found,
			// but first it should found the field in the rest of struct
			// (a field with same name in the current struct should have preference over anonymous struct)
			anon = dec.curr
			dec.curr = tmp
		} else if dec.field == field.Tag.Get(tagName) {
			// is not found yet, then retry by its tag name "formam"
			dec.curr = dec.curr.Field(i)
			return nil
		}
	}

	if anon.IsValid() {
		dec.curr = anon
		return nil
	}

	return fmt.Errorf("formam: not found the field \"%v\" in the path \"%v\"", dec.field, dec.path)
}

// expandSlice expands the length and capacity of the current slice
func (dec *decoder) expandSlice(length int) {
	n := reflect.MakeSlice(dec.curr.Type(), length, length)
	reflect.Copy(n, dec.curr)
	dec.curr.Set(n)
}

// currentMap gets in d.curr the map concrete for decode the current value
func (dec *decoder) currentMap() {
	n := dec.curr.Type()
	if dec.curr.IsNil() {
		dec.curr.Set(reflect.MakeMap(n))
		m := reflect.New(n.Elem()).Elem()
		dec.maps = append(dec.maps, &pathMap{dec.curr, dec.field, m, dec.path})
		dec.curr = m
	} else if a := dec.maps.find(dec.curr, dec.field); a == nil {
		m := reflect.New(n.Elem()).Elem()
		dec.maps = append(dec.maps, &pathMap{dec.curr, dec.field, m, dec.path})
		dec.curr = m
	} else {
		dec.curr = a.value
	}
}

var (
	timeType  = reflect.TypeOf(time.Time{})
	timePType = reflect.TypeOf(&time.Time{})
)

// unmarshalText returns a boolean and error. The boolean is true if the
// value implements TextUnmarshaler, and false if not.
func (dec *decoder) unmarshalText(v reflect.Value) (bool, error) {
	// check if implements the interface
	m, ok := v.Interface().(encoding.TextUnmarshaler)
	addr := v.CanAddr()
	if !ok && !addr {
		return false, nil
	} else if addr {
		return dec.unmarshalText(v.Addr())
	}
	// skip if the type is time.Time
	n := v.Type()
	if n.ConvertibleTo(timeType) || n.ConvertibleTo(timePType) {
		return false, nil
	}
	// return result
	return true, m.UnmarshalText([]byte(dec.value))
}
