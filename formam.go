package formam

import (
	"net/http"
	"strings"
	"errors"
	"reflect"
	"strconv"
	"fmt"
)

const TAG_NAME = "formam"

type decoder struct {
	r   *http.Request
	dst interface {}

	main reflect.Value
	curr reflect.Value
}

func NewDecoder(r *http.Request, dst interface {}) (*decoder, error) {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Ptr {
		return nil, errors.New("formam: is not a pointer to struct")
	}
	if value.Elem().Kind() != reflect.Struct {
		return nil, errors.New("formam: is not to struct")
	}
	return &decoder{r: r, dst: dst, main: value.Elem()}, nil
}

func (d *decoder) Decode() error {
	d.r.ParseForm()
	for k, v := range d.r.Form {
		d.decode(k, v[0])
	}
	return nil
}

func (d *decoder) decode(key, value string) error {
	fields := strings.Split(key, ".")
	d.curr = d.main
	for i, field := range fields {
		b := strings.IndexAny(field, "[")
		if b != -1 {
			// is a array
			e := strings.IndexAny(field, "]")
			if e == -1 {
				return errors.New("formam: bad syntax array")
			}
			name := field[:b]
			index, err := strconv.Atoi(field[b+1:e])
			if err != nil {
				return errors.New("formam: the index of array not is a number")
			}
			if len(fields) == i+1 {
				return d.end(name, value, index)
			} else {
				d.curr, err = d.walk(name, index)
				if err != nil {
					return err
				}
			}
		} else {
			// not is a array
			if len(fields) == i+1 {
				return d.end(field, value, -1)
			} else {
				var err error
				d.curr, err = d.walk(field, -1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *decoder) walk(field string, index int) (reflect.Value, error) {
	d.findField(field)
	if index != -1 {
		// should be a array...
		switch d.curr.Kind() {
		case reflect.Slice, reflect.Array:
			len := d.curr.Len()
			if len <= index {
				len++
				d.curr.Set(reflect.AppendSlice(d.curr, reflect.MakeSlice(d.curr.Type(), len, len)))
			}
			d.curr = d.curr.Index(index)
		default:
			return d.curr, fmt.Errorf("formam: the field \"%v\" not should be a array", field)
		}
	}
	return d.curr, nil
}

func (d *decoder) end(field, value string, index int) error {
	if d.curr.Kind() == reflect.Struct {
		d.findField(field)
	}
	switch d.curr.Kind() {
	case reflect.Map:
		if d.curr.IsNil() {
			d.curr.Set(reflect.MakeMap(d.curr.Type()))
		}
		d.curr.SetMapIndex(reflect.ValueOf(field), reflect.ValueOf(value))
	case reflect.Slice, reflect.Array:

	case reflect.String:
		d.curr.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.Interface:
	default:
		return fmt.Errorf("formam: not supported type for filed \"%v\"", field)
	}
	return nil
}

func (d *decoder) findField(field string) {
	if v := d.curr.FieldByName(field); v.Kind() == reflect.Invalid {
		num := d.curr.NumField()
		for i := 0; i < num; i++ {
			f := d.curr.Type().Field(i).Tag.Get(TAG_NAME)
			if field == f {
				d.curr = d.curr.Field(i)
				break
			}
		}
	} else {
		d.curr = v
	}
}
