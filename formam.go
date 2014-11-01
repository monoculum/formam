package formam

import (
	"strings"
	"errors"
	"reflect"
	"strconv"
	"fmt"
	"net/url"
)

const TAG_NAME = "formam"

// A decoder holds the values from form, the 'reflect' value of main struct
// and the 'reflect' value of current path
type decoder struct {
	main reflect.Value
	curr reflect.Value

	key   string
	value string
	index int
}

// NewDecoder generates a decoder struct with url.Values and struct provided by argument
func Decode(v url.Values, dst interface {}) error {
	main := reflect.ValueOf(dst)
	if main.Kind() != reflect.Ptr || main.Elem().Kind() != reflect.Struct {
		return errors.New("formam: is not a pointer to struct")
	}
	d := &decoder{main: main.Elem()}
	for k, v := range v {
		d.key = k
		d.value = v[0]
		if err := d.begin(); err != nil {
			return err
		}
	}
	return nil
}

// decode prepare the path of the current key of map to walk through it
func (d *decoder) begin() (err error) {
	d.curr = d.main
	fields := strings.Split(d.key, ".")
	for i, field := range fields {
		b := strings.IndexAny(field, "[")
		if b != -1 {
			// is a array
			e := strings.IndexAny(field, "]")
			if e == -1 {
				return errors.New("formam: bad syntax array")
			}
			d.key = field[:b]
			d.index, err = strconv.Atoi(field[b+1:e])
			if err != nil {
				return errors.New("formam: the index of array not is a number")
			}
			if len(fields) == i+1 {
				return d.end()
			} else {
				d.curr, err = d.walk()
				if err != nil {
					return err
				}
			}
		} else {
			// not is a array
			d.key = field
			d.index = -1
			if len(fields) == i+1 {
				return d.end()
			} else {
				d.curr, err = d.walk()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// walk traverse the path to the final field for set the value
func (d *decoder) walk() (reflect.Value, error) {
	if err := d.findField(); err != nil {
		return d.curr, err
	}
	if d.index != -1 {
		// should be a array...
		switch d.curr.Kind() {
		case reflect.Slice, reflect.Array:
			len := d.curr.Len()
			if len <= d.index {
				len = len-d.index+1
				d.curr.Set(reflect.AppendSlice(d.curr, reflect.MakeSlice(d.curr.Type(), len, len)))
			}
			d.curr = d.curr.Index(d.index)
		default:
			return d.curr, fmt.Errorf("formam: the field \"%v\" not should be a array", d.key)
		}
	}
	return d.curr, nil
}

// end the last field for set its value correspondent
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
		if d.curr.IsNil() {
			d.curr.Set(reflect.MakeMap(d.curr.Type()))
		}
		//d.curr.MapIndex()
		d.curr.SetMapIndex(reflect.ValueOf(d.key), reflect.ValueOf(d.value))
	case reflect.Slice, reflect.Array:
		if d.curr.Len() <= d.index {
			sl := reflect.MakeSlice(d.curr.Type(), d.index+1, d.index+1)
			reflect.Copy(sl, d.curr)
			d.curr.Set(sl)
		}
		d.curr = d.curr.Index(d.index)
		return d.decode()
	case reflect.String:
		d.curr.SetString(d.value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, err := strconv.ParseInt(d.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid signed integer number", d.key)
		} else {
			d.curr.SetInt(num)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if num, err := strconv.ParseUint(d.value, 10, 64); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid unsigned integer number", d.key)
		} else {
			d.curr.SetUint(num)
		}
	case reflect.Float32, reflect.Float64:
		if num, err := strconv.ParseFloat(d.value, d.curr.Type().Bits()); err != nil {
			return fmt.Errorf("formam: the value \"%v\" should be a valid float number", d.key)
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
			return fmt.Errorf("formam: the value \"%v\" is not a valid boolean", d.key)
		}
	case reflect.Interface:
	default:
		return fmt.Errorf("formam: not supported type for field \"%v\"", d.key)
	}
	return nil
}

// findField find a field by its name, if it is not found,
// then retry the search examining the tag "formam" of every field of struct
func (d *decoder) findField() error {
	num := d.curr.NumField()
	for i := 0; i < num; i++ {
		field := d.curr.Type().Field(i)
		if field.Name == d.key {
			// check if the field's name is equal
			d.curr = d.curr.Field(i)
			return nil
		} else if field.Anonymous {
			// if the field is anonymous, then iterate over its sub fields
			d.curr = d.curr.FieldByIndex(field.Index)
			return d.findField()
		} else if d.key == field.Tag.Get(TAG_NAME) {
			// is not found yet, then retry by its tag name "formam"
			d.curr = d.curr.Field(i)
			return nil
		}
	}
	return fmt.Errorf("formam: not found the field \"%v\"", d.key)
}
