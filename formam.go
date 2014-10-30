package formam

import (
	"net/http"
	"strings"
	"errors"
	"reflect"
	"strconv"
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
				return d.end(name, value)
			} else {
				d.curr, err = d.walk(name, index)
				if err != nil {
					return err
				}
			}
		} else {
			// not is a array
			if len(fields) == i+1 {
				return d.end(field, value)
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
	v := d.curr.FieldByName(field)
	if v.Kind() == reflect.Invalid {
		for i := 0; i < d.curr.NumField(); i++ {
			f := d.curr.Type().Field(i).Tag.Get(TAG_NAME)
			if field == f {
				d.curr = d.curr.Field(i)
				break
			}
		}
	} else {
		if index != -1 {
			// should be a array...
			switch v.Kind() {
			case reflect.Slice, reflect.Array:
				len := v.Len()
				if len <= index {
					v.Set(reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), len+1, len+1)))
				}
				d.curr = v.Index(index)
			default:
			}
		} else {
			d.curr = v
		}
	}
	return d.curr, nil
}

func (d *decoder) end(field, value string) error {
	var v reflect.Value
	switch d.curr.Kind() {
	case reflect.Struct:
		v = d.curr.FieldByName(field)
		if v.Kind() == reflect.Invalid {
			for i := 0; i < d.curr.NumField(); i++ {
				f := d.curr.Type().Field(i).Tag.Get(TAG_NAME)
				if field == f {
					v = d.curr.Field(i)
					break
				}
			}
		}
	case reflect.Map:
		if d.curr.IsNil() {
			d.curr.Set(reflect.MakeMap(d.curr.Type()))
		}
		d.curr.SetMapIndex(reflect.ValueOf(field), reflect.ValueOf(value))
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	}

	return nil
}
