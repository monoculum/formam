# formam

A Go package to decode HTTP form and query parameters.
The only requirement is [Go 1.12](http://golang.org/doc/go1.12) or later.

[![Build Status](https://travis-ci.org/monoculum/formam.svg?branch=master)](https://travis-ci.org/monoculum/formam)
[![GoDoc](https://godoc.org/github.com/monoculum/formam/v3?status.svg)](https://pkg.go.dev/github.com/monoculum/formam/v3)

## Install

```
go get github.com/monoculum/formam/v3
```

## Features

- Infinite nesting for `maps`, `structs` and `slices`.
- Support `UnmarshalText()` interface in values and keys of maps.
- Supported `map` keys are `string`, `int` and variants, `uint` and variants, `uintptr`, `float32`, `float64`, `bool`, `struct`, `custom types` to one of the above types registered by function or `UnmarshalText` method, a `pointer` to one of the above types
- A field with `interface{}` that has a `map`, `struct` or `slice` as value is accessible.
- Decode `time.Time` with format `2006-01-02` by its `UnmarshalText()` method.
- Decode `url.URL`.
- Append to `slice` and `array` types without explicitly indicating an index.
- Register a function for a custom type.

## Performance

You can see the performance in [formam-benchmark](https://github.com/monoculum/formam-benchmark) compared with [ajg/form](https://github.com/ajg/form), [gorilla/schema](https://github.com/gorilla/schema), [go-playground/form](https://github.com/go-playground/form) and [built-in/json](http://golang.org/pkg/encoding/json/).

## Basic usage example

### In form HTML

- Use `.` to access a struct field (e.g. `struct.field1`).
- Use `[<index>]` to access tje specific slice/array index (e.g. `struct.array[0]`). It's not necessary to add an index to append data.
- Use `[<key>]` to access map keys (e.g.. `struct.map[es-ES]`).

```html
<form method="POST">
  <input type="text" name="Name" value="Sony" />
  <input type="text" name="Location.Country" value="Japan" />
  <input type="text" name="Location.City" value="Tokyo" />
  <input type="text" name="Products[0].Name" value="Playstation 4" />
  <input type="text" name="Products[0].Type" value="Video games" />
  <input type="text" name="Products[1].Name" value="TV Bravia 32" />
  <input type="text" name="Products[1].Type" value="TVs" />
  <input type="text" name="Founders[0]" value="Masaru Ibuka" />
  <input type="text" name="Founders[0]" value="Akio Morita" />
  <input type="text" name="Employees" value="90000" />
  <input type="text" name="public" value="true" />
  <input type="url" name="website" value="http://www.sony.net" />
  <input type="date" name="foundation" value="1946-05-07" />
  <input type="text" name="Interface.ID" value="12" />
  <input type="text" name="Interface.Name" value="Go Programming Language" />
  <input type="submit" />
</form>
```

### In Go

You can use the `formam` struct tag to ensure the form values are unmarshalled in the currect struct fields.

```go
type InterfaceStruct struct {
    ID   int
    Name string
}

type Company struct {
  Public     bool      `formam:"public"`
  Website    url.URL   `formam:"website"`
  Foundation time.Time `formam:"foundation"`
  Name       string
  Location   struct {
    Country  string
    City     string
  }
  Products   []struct {
    Name string
    Type string
  }
  Founders   []string
  Employees  int64

  Interface interface{}
}

func MyHandler(w http.ResponseWriter, r *http.Request) error {
  r.ParseForm()

  m := Company{
      // it's is possible to access to the fields although it's an interface field!
      Interface: &InterfaceStruct{},
  }
  dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
  return dec.Decode(r.Form, &m)
}
```

## Types

Supported types in the destination struct are:

- `string`
- `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `slice`, `array`
- `struct` and `struct anonymous`
- `map`
- `interface{}`
- `time.Time`
- `url.URL`
- `custom types` to one of the above types
- a `pointer` to one of the above types

## Custom Marshaling

You can umarshal data and map keys by implementing the `encoding.TextUnmarshaler` interface.

If the forms sends multiple values then only the first value is passed to `UnmarshalText()`, but if the name ends with `[]` then it's called for all values.

## Custom Type

You can register a function for a custom type using the `RegisterCustomType()` method. This will work for any number of given fields or all fields with the given type.

Registered type have preference over the UnmarshalText method unless the `PrefUnmarshalText` option is used.

### All fields

```go
decoder.RegisterCustomType(func(vals []string) (interface{}, error) {
        return time.Parse("2006-01-02", vals[0])
}, []interface{}{time.Time{}}, nil)
```

### Specific fields

```go
package main

type Times struct {
    Timestamp   time.Time
    Time        time.Time
    TimeDefault time.Time
}

func main() {
    var t Timestamp

    dec := NewDecoder(nil)

    // for Timestamp field
    dec.RegisterCustomType(func(vals []string) (interface{}, error) {
            return time.Parse("2006-01-02T15:04:05Z07:00", vals[0])
    }, []interface{}{time.Time{}}, []interface{}{&t.Timestamp{}})

    // for Time field
    dec.RegisterCustomType(func(vals []string) (interface{}, error) {
                return time.Parse("Mon, 02 Jan 2006 15:04:05 MST", vals[0])
    }, []interface{}{time.Time{}}, []interface{}{&t.Time{}})

    // for field that not be Time or Timestamp, e.g. in this example, TimeDefault.
    dec.RegisterCustomType(func(vals []string) (interface{}, error) {
                return time.Parse("2006-01-02", vals[0])
    }, []interface{}{time.Time{}}, nil)

    dec.Decode(url.Values{}, &t)
}
```

## Notes

Version 2 is compatible with old syntax to access to maps (`map.key`), but brackets are the preferred way to access a map (`map[key])`.
