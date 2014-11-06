package formam

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const TAG_NAME = "formam"

// A decoder holds the values from form, the 'reflect' value of main struct
// and the 'reflect' value of current path
type pathMap struct {
	m reflect.Value

	key   string
	value reflect.Value
}

type pathMaps []*pathMap

func (ma pathMaps) find(id reflect.Value, key string) *pathMap {
	for _, v := range ma {
		if v.m == id && v.key == key {
			return v
		}
	}
	return nil
}

type decoder struct {
	main reflect.Value
	curr reflect.Value

	maps pathMaps

	field string
	value string
	index int
}

// NewDecoder generates a decoder struct with url.Values and struct provided by argument
func Decode(vv url.Values, dst interface{}) error {
	main := reflect.ValueOf(dst)
	if main.Kind() != reflect.Ptr || main.Elem().Kind() != reflect.Struct {
		return errors.New("formam: is not a pointer to struct")
	}
	d := &decoder{main: main.Elem()}
	for k, v := range vv {
		d.field = k
		d.value = v[0]
		if err := d.begin(); err != nil {
			return err
		}
	}
	for _, v := range d.maps {
		k := reflect.New(v.m.Type().Key()).Elem()
		d.curr = k
		d.value = v.key
		if err := d.decode(); err != nil {
			return err
		}
		v.m.SetMapIndex(d.curr, v.value)
	}
	d.maps = []*pathMap{}
	return nil
}

// decode prepare the path of the current key of map to walk through it
func (d *decoder) begin() (err error) {
	d.curr = d.main
	fields := strings.Split(d.field, ".")
	for i, field := range fields {
		b := strings.IndexAny(field, "[")
		if b != -1 {
			// is a array
			e := strings.IndexAny(field, "]")
			if e == -1 {
				return errors.New("formam: bad syntax array")
			}
			d.field = field[:b]
			if d.index, err = strconv.Atoi(field[b+1 : e]); err != nil {
				return errors.New("formam: the index of array not is a number")
			}
			if len(fields) == i+1 {
				return d.end()
			}
			if d.curr, err = d.walk(); err != nil {
				return
			}
		} else {
			// not is a array
			d.field = field
			d.index = -1
			if len(fields) == i+1 {
				return d.end()
			}
			if d.curr, err = d.walk(); err != nil {
				return
			}
		}
	}
	return
}

// walk traverse the path to the final field for set the value
func (d *decoder) walk() (reflect.Value, error) {
	switch d.curr.Kind() {
	case reflect.Struct:
		if err := d.findField(); err != nil {
			return d.curr, err
		}
	case reflect.Map:
		d.currentMap()
	case reflect.Ptr:
		d.curr.Set(reflect.New(d.curr.Type().Elem()))
		d.curr = d.curr.Elem()
		return d.walk()
	}
	if d.index != -1 {
		// should be a array...
		switch d.curr.Kind() {
		case reflect.Slice, reflect.Array:
			if d.curr.Len() <= d.index {
				d.expandSlice()
			}
			d.curr = d.curr.Index(d.index)
		default:
			return d.curr, fmt.Errorf("formam: the field \"%v\" not should be a array", d.field)
		}
	}
	return d.curr, nil
}

// end find the last field for decode/set its value correspondent
func (d *decoder) end() error {
	if d.curr.Kind() == reflect.Struct {
		if err := d.findField(); err != nil {
			return err
		}
	}
	return d.decode()
}

// decode set the value in its field
func (d *decoder) decode() error {
	switch d.curr.Kind() {
	case reflect.Map:
		d.currentMap()
		return d.decode()
	case reflect.Slice, reflect.Array:
		if d.curr.Len() <= d.index {
			d.expandSlice()
		}
		d.curr = d.curr.Index(d.index)
		return d.decode()
	case reflect.String:
		d.curr.SetString(d.value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, err := strconv.ParseInt(d.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid signed integer number", d.field)
		} else {
			d.curr.SetInt(num)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if num, err := strconv.ParseUint(d.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid unsigned integer number", d.field)
		} else {
			d.curr.SetUint(num)
		}
	case reflect.Float32, reflect.Float64:
		if num, err := strconv.ParseFloat(d.value, d.curr.Type().Bits()); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid float number", d.field)
		} else {
			d.curr.SetFloat(num)
		}
	case reflect.Bool:
		switch d.value {
		case "true", "1":
			d.curr.SetBool(true)
		case "false", "0":
			d.curr.SetBool(false)
		default:
			return fmt.Errorf("formam: the value \"%v\" is not a valid boolean", d.field)
		}
	case reflect.Interface:
		d.curr.Set(reflect.ValueOf(d.value))
	case reflect.Ptr:
		d.curr.Set(reflect.New(d.curr.Type().Elem()))
		d.curr = d.curr.Elem()
		return d.decode()
	default:
		return fmt.Errorf("formam: not supported type for field \"%v\"", d.field)
	}
	return nil
}

// findField find a field by its name, if it is not found,
// then retry the search examining the tag "formam" of every field of struct
func (d *decoder) findField() error {
	num := d.curr.NumField()
	for i := 0; i < num; i++ {
		field := d.curr.Type().Field(i)
		if field.Name == d.field {
			// check if the field's name is equal
			d.curr = d.curr.Field(i)
			return nil
		} else if field.Anonymous {
			// if the field is a anonymous struct, then iterate over its fields
			d.curr = d.curr.FieldByIndex(field.Index)
			return d.findField()
		} else if d.field == field.Tag.Get(TAG_NAME) {
			// is not found yet, then retry by its tag name "formam"
			d.curr = d.curr.Field(i)
			return nil
		}
	}
	return fmt.Errorf("formam: not found the field \"%v\"", d.field)
}

// expandSlice expand the length and capacity of the current slice
func (d *decoder) expandSlice() {
	sli := reflect.MakeSlice(d.curr.Type(), d.index+1, d.index+1)
	reflect.Copy(sli, d.curr)
	d.curr.Set(sli)
}

// currentMap get in d.curr the map concrete for decode the current value
func (d *decoder) currentMap() {
	typ := d.curr.Type()
	if d.curr.IsNil() {
		d.curr.Set(reflect.MakeMap(typ))
		v := reflect.New(typ.Elem()).Elem()
		d.maps = append(d.maps, &pathMap{d.curr, d.field, v})
		d.curr = v
	} else if a := d.maps.find(d.curr, d.field); a == nil {
		v := reflect.New(typ.Elem()).Elem()
		d.maps = append(d.maps, &pathMap{d.curr, d.field, v})
		d.curr = v
	} else {
		d.curr = a.value
	}
}
