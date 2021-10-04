package formam_test

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/monoculum/formam"
)

type Text string

func (s *Text) UnmarshalText(text []byte) error {
	var n Text
	n = "the string has changed by UnmarshalText method"
	*s = n
	return nil
}

type UUID [16]byte

func (u *UUID) UnmarshalText(text []byte) error {
	if len(text) != 32 {
		return fmt.Errorf("text must be exactly 16 bytes long, got %d bytes", len(text))
	}
	_, err := hex.Decode(u[:], text)
	if err != nil {
		return err
	}
	return nil
}

func (u UUID) String() string {
	buf := make([]byte, 32)
	hex.Encode(buf[:], u[:])
	return string(buf)
}

const unmarshalTextString = "If you see this text, then it's a bug"

type AnonymousID struct {
	ID string
}

type Anonymous struct {
	AnonymousField string
	FieldOverride  string
	*AnonymousID
}

type FieldString string

type TestStruct struct {
	Anonymous
	FieldOverride string

	// traverse
	TraverseStruct struct {
		Field1 [][]struct {
			Field string
		}
		Field2 struct {
			Field int
		}
	}
	TraverseMapByBracket map[string]map[int]map[uint]map[bool]*string
	TraverseMapByPoint   map[string]map[int]map[uint]map[bool]string

	// slices/arrays
	SlicesWithIndex      []string
	SlicesWithoutIndex   []float32
	SlicesMultiDimension [][][][]uintptr
	ArrayWithIndex       [2]interface{}
	ArrayWithoutIndex    [2]bool
	ArrayMultiDimension  [2][2]bool

	// int
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	// uint
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr

	// bool
	Bool bool

	// string
	String string

	// byte
	Byte byte

	// pointer
	Pointer *string
	// pointer to struct
	PointerToStruct *struct{ Field float64 }
	// pointer to map
	PointerToMap *map[string]string
	// pointer to anonymous struct
	PointerToSlice []Anonymous

	// map
	Map map[string]string
	// mp with slice
	MapWithSlice map[string][]string
	// map with multi dimension slice
	MapWithMultiDimensionSlice map[string][][]string
	// map with array
	MapWithArray map[string][2]int
	// map with int key
	MapWithIntKey map[int]string
	// map with int8 key
	MapWithInt8Key map[int8]string
	// map with *int64 key
	MapWithInt64PtrKey map[*int64]string
	// map with uint key
	MapWithUintKey map[uint]string
	// map with uint key
	MapWithUint8Key map[uint8]string
	// map with uint key
	MapWithUint32PtrKey map[*uint32]string
	// map with float key
	MapWithFloatKey map[float32]string
	// map with boolean key
	MapWithBooleanKey map[bool]string
	// map with custom key and decode key by unmarshal key
	MapWithCustomKey map[UUID]string
	// map with custom key pointer and decode key by unmarshal key
	MapWithCustomKeyPointer map[*UUID]string
	// map with time.Time Key
	MapWithStruct1Key map[time.Time]string
	// map with url.URL Key
	MapWithStruct2Key map[url.URL]string
	//MapWithStruct3Key map[struct {ID struct {ID string}]string

	// unmarshal text
	UnmarshalTextString Text
	UnmarshalTextUUID   UUID

	// tag
	Tag    string `formam:"tag"`
	TagOpt string `formam:"tagopt,x"`

	// time
	Time time.Time

	// url
	URL url.URL

	// interface
	Interface interface{}
	// interface with struct as data
	InterfaceStruct interface{}

	// custom type
	CustomType FieldString
	// custom type by field
	Time1       time.Time
	Time2       time.Time
	TimeDefault time.Time
}

type InterfaceStruct struct {
	ID   int
	Name string
}

var vals = url.Values{
	// anonymous
	"AnonymousField": []string{"anonymous field"},
	"FieldOverride":  []string{"field not override"},

	// traverse
	"TraverseStruct.Field1[0][2].Field":            []string{"traverse over structs is recursive"},
	"TraverseStruct.Field2.Field":                  []string{"2"},
	"TraverseMapByBracket[by-bracket][1][2][true]": []string{"traverse over map by bracket is recursive too"},
	"TraverseMapByPoint.by-point.1.2.true":         []string{"traverse over map by point is recursive too"},

	// slices/arrays
	"SlicesWithIndex[0]":               []string{"slice index 0"},
	"SlicesWithIndex[2]":               []string{"slice index 2"},
	"SlicesWithIndex[4]":               []string{"slice index 4"},
	"SlicesWithoutIndex":               []string{"1.111", "2.222", "3.333"},
	"SlicesMultiDimension[0][1][2][3]": []string{"8"},
	"ArrayWithIndex[0]":                []string{"array index 0"},
	"ArrayWithIndex[1]":                []string{"array index 1"},
	"ArrayWithoutIndex":                []string{"true", "true"},
	"ArrayMultiDimension[0][0]":        []string{"true"},
	"ArrayMultiDimension[0][1]":        []string{"true"},
	"ArrayMultiDimension[1][0]":        []string{"true"},
	"ArrayMultiDimension[1][1]":        []string{"true"},

	// int
	"Int":   []string{"-1"},
	"Int8":  []string{"-1"},
	"Int16": []string{"-1"},
	"Int32": []string{"-1"},
	"Int64": []string{"-1"},

	// uint
	"Uint":    []string{"1"},
	"Uint8":   []string{"1"},
	"Uint16":  []string{"1"},
	"Uint32":  []string{"1"},
	"Uint64":  []string{"1"},
	"Uintptr": []string{"10"},

	// bool
	"Bool": []string{"true"},

	// string
	"String": []string{"string"},

	// byte
	"Byte": []string{"20"},

	// pointer
	"Pointer":               []string{"20"},
	"PointerToStruct.Field": []string{"20"},
	"PointerToMap[es]":      []string{"20"},
	"PointerToSlice[0].ID":  []string{"20"},

	// map
	"Map[by.bracket.with.point]":                                []string{"by bracket"},
	"Map.by_point":                                              []string{"by point"},
	"MapWithSlice[slice][0]":                                    []string{"map with slice"},
	"MapWithMultiDimensionSlice[slice][0][1]":                   []string{"map with multidimension slice"},
	"MapWithArray[array][0]":                                    []string{"0"},
	"MapWithArray[array][1]":                                    []string{"1"},
	"MapWithIntKey[-1]":                                         []string{"int key in map"},
	"MapWithInt8Key[-1]":                                        []string{"int8 key in map"},
	"MapWithInt64PtrKey[-1]":                                    []string{"int64 ptr key in map"},
	"MapWithUint8Key[1]":                                        []string{"uint8 ptr key in map"},
	"MapWithUint32PtrKey[1]":                                    []string{"uint32 ptr key in map"},
	"MapWithUintKey[1]":                                         []string{"uint key in map"},
	"MapWithFloatKey[3.14]":                                     []string{"float key in map"},
	"MapWithBooleanKey[true]":                                   []string{"bool key in map"},
	"MapWithCustomKey[11e5bf2d3e403a8c86740023dffe5350]":        []string{"UUID key in map"},
	"MapWithCustomKeyPointer[11e5bf2d3e403a8c86740023dffe5350]": []string{"UUID key pointer in map"},
	"MapWithStruct1Key[2006-01-02]":                             []string{"time.Time key in map"},
	"MapWithStruct2Key[http://www.monoculum.com]":               []string{"url.URL key in map"},

	// unmarshal text
	"UnmarshalTextString": []string{"If you see this text, then it's a bug"},
	"UnmarshalTextUUID":   []string{"11e5bf2d3e403a8c86740023dffe5350"},

	// tag
	"tag":    []string{"string placed by tag"},
	"tagopt": []string{"string placed by tagopt"},

	// time
	"Time": []string{"2016-06-12"},

	// url
	"URL": []string{"http://www.monoculum.com"},

	// interface
	"Interface":            []string{"Germany"},
	"InterfaceStruct.ID":   []string{"1"},
	"InterfaceStruct.Name": []string{"Germany"},

	// custom type
	"CustomType": []string{"if you see this text, then it's a bug"},
	// custom type by field
	"Time1":       []string{"2001-01-01"},
	"Time2":       []string{"2001-01-01"},
	"TimeDefault": []string{"2001-01-01"},
}

func TestDecodeInStruct(t *testing.T) {
	var m TestStruct
	m.InterfaceStruct = &InterfaceStruct{}

	dec := formam.NewDecoder(nil).RegisterCustomType(func(vals []string) (interface{}, error) {
		return FieldString("value changed by custom type"), nil
	}, []interface{}{FieldString("")}, nil)

	dec.RegisterCustomType(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", "2016-01-02")
	}, []interface{}{time.Time{}}, []interface{}{&m.Time1})

	dec.RegisterCustomType(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", "2017-01-02")
	}, []interface{}{time.Time{}}, []interface{}{&m.Time2})

	dec.RegisterCustomType(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", "2018-01-02")
	}, []interface{}{time.Time{}}, []interface{}{})

	err := dec.Decode(vals, &m)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// anonymous struct
	if m.Anonymous.AnonymousField == "" {
		t.Error("AnonymousField is empty")
	}
	if m.Anonymous.FieldOverride != "" {
		t.Error("FieldOverride is full")
	}
	if m.FieldOverride == "" {
		t.Error("FieldOverride is empty")
	}

	// traverse
	// traverse > struct
	if len(m.TraverseStruct.Field1) == 0 {
		t.Error("TraverseStruct.Field1 is empty")
	} else {
		if len(m.TraverseStruct.Field1[0]) != 3 {
			t.Errorf("TraverseStruct.Field1[0] must has 3 as length but has %v", len(m.TraverseStruct.Field1[0]))
		} else {
			if m.TraverseStruct.Field1[0][2].Field == "" {
				t.Error("TraverseStruct.Field1[0][2].Field is empty")
			}
		}
	}
	// traverse > maps
	n, ok := m.TraverseMapByBracket["by-bracket"]
	if ok {
		m, ok := n[1]
		if ok {
			j, ok := m[2]
			if ok {
				g, ok := j[true]
				if ok {
					if *g == "" {
						t.Error("the value of TraverseMapByBracket in the last field is empty")
					}
				} else {
					t.Error("the key \"true\" in TraverseMapByBracket not exists")
				}
			} else {
				t.Error("the key \"2\" in TraverseMapByBracket not exists")
			}
		} else {
			t.Error("the key \"1\" in TraverseMapByBracket not exists")
		}
	} else {
		t.Error("the key \"bracket\" in TraverseMapByBracket not exists")
	}
	u, ok := m.TraverseMapByPoint["by-point"]
	if ok {
		m, ok := u[1]
		if ok {
			j, ok := m[2]
			if ok {
				g, ok := j[true]
				if ok {
					if g == "" {
						t.Error("the value of TraverseMapByPoint in the last field is empty")
					}
				} else {
					t.Error("the key \"true\" in TraverseMapByPoint not exists")
				}
			} else {
				t.Error("the key \"2\" in TraverseMapByPoint not exists")
			}
		} else {
			t.Error("the key \"1\" in TraverseMapByPoint not exists")
		}
	} else {
		t.Error("the key \"by-point\" in TraverseMapByPoint not exists")
	}

	// slices
	if len(m.SlicesWithIndex) != 5 {
		t.Error("the length SlicesWithIndex is not 5")
	}
	if len(m.SlicesWithoutIndex) != 3 {
		t.Error("the length SlicesWithoutIndex is not 3")
	}
	if len(m.SlicesMultiDimension) != 1 {
		t.Error("the length SlicesMultiDimension is not 1")
	}
	if len(m.SlicesMultiDimension[0]) != 2 {
		t.Error("the length SlicesMultiDimension[0] is not 2")
	}
	if len(m.SlicesMultiDimension[0][1]) != 3 {
		t.Error("the length SlicesMultiDimension[0] is not 3")
	}
	// array
	if len(m.ArrayWithIndex) != 2 {
		t.Error("the length ArrayWithIndex is not 2")
	}
	if len(m.ArrayWithoutIndex) != 2 {
		t.Error("the length ArrayWithoutIndex is not 2")
	}
	if len(m.ArrayMultiDimension) != 2 {
		t.Error("the length ArrayMultiDimension is not 2")
	}
	if len(m.ArrayMultiDimension[0]) != 2 {
		t.Error("the length ArrayMultiDimension[0] is not 2")
	}
	if len(m.ArrayMultiDimension[1]) != 2 {
		t.Error("the length ArrayMultiDimension[1] is not 2")
	}

	// int
	if m.Int != -1 {
		t.Error("the length Int is not -1")
	}
	if m.Int8 != -1 {
		t.Error("the length Int8 is not -1")
	}
	if m.Int16 != -1 {
		t.Error("the length Int16 is not -1")
	}
	if m.Int32 != -1 {
		t.Error("the length Int32 is not -1")
	}
	if m.Int64 != -1 {
		t.Error("the length Int64 is not -1")
	}

	// uint
	if m.Uint != 1 {
		t.Error("the length Uint is not 1")
	}
	if m.Uint8 != 1 {
		t.Error("the length Uint8 is not 1")
	}
	if m.Uint16 != 1 {
		t.Error("the length Uint16 is not 1")
	}
	if m.Uint32 != 1 {
		t.Error("the length Uint32 is not 1")
	}
	if m.Uint64 != 1 {
		t.Error("the length Uint64 is not 1")
	}
	if m.Uintptr != 10 {
		t.Error("the length Uintptr is not 10")
	}

	// bool
	if !m.Bool {
		t.Error("Bool is false")
	}

	// string
	if m.String == "" {
		t.Error("String is empty")
	}

	// byte
	if string(m.Byte) == "" {
		t.Error("Byte is empty")
	}

	// pointer
	if m.Pointer == nil {
		t.Error("Pointer is nil")
	} else if *m.Pointer == "" {
		t.Error("Pointer is not nil but is empty")
	}
	if m.PointerToMap == nil {
		t.Error("Pointer is nil")
	} else if len(*m.PointerToMap) == 0 {
		t.Error("PointerToMap is not nil but is empty")
	} else {
		for k := range *m.PointerToMap {
			if (*m.PointerToMap)[k] == "" {
				t.Error("PointerToMap[" + k + "] is empty")
			}
		}
	}
	if m.PointerToSlice == nil {
		t.Error("PointerToSlice is nil")
	} else if len(m.PointerToSlice) == 0 {
		t.Error("PointerToSlice is not nil but is empty")
	} else {
		for i := range m.PointerToSlice {
			if m.PointerToSlice[i].AnonymousID == nil {
				t.Error("PointerToSlice[" + strconv.Itoa(i) + "] is nil")
			} else if m.PointerToSlice[i].ID == "" {
				t.Error("PointerToSlice[" + strconv.Itoa(i) + "].ID is empty")
			}
		}
	}

	// map
	f, ok := m.Map["by.bracket.with.point"]
	if ok {
		if f == "" {
			t.Error("The value in key \"by.bracket.with.point\" of Map is empty")
		}
	} else {
		t.Error("The key \"by.bracket.with.point\" in Map not exists")
	}
	f, ok = m.Map["by_point"]
	if ok {
		if f == "" {
			t.Error("The value in key \"by_point\" of Map is empty")
		}
	} else {
		t.Error("The key \"by_point\" in Map not exists")
	}
	s, ok := m.MapWithSlice["slice"]
	if ok {
		if len(s) == 0 {
			t.Error("The length of key \"slice\" of MapWithSlice is 0")
		} else {
			if s[0] == "" {
				t.Error("The value of key \"slice\" in MapWithSlice is empty")
			}
		}
	} else {
		t.Error("The key \"slice\" in MapWithSlice not exists")
	}
	a, ok := m.MapWithMultiDimensionSlice["slice"]
	if ok {
		if len(a) == 0 {
			t.Error("The length of key \"slice\" of MapWithSlice is 0")
		} else {
			if len(a) == 0 {
				t.Error("The length of MapWithMultiDimensionSlice[slice] is 0")
			} else {
				if len(a[0]) != 2 {
					t.Error("The length of MapWithMultiDimensionSlice[slice][0] is not 2")
				} else {
					if a[0][1] == "" {
						t.Error("The value in MapWithSlice[slice][0][1] is empty")
					}
				}
			}
		}
	} else {
		t.Error("The key \"slice\" in MapWithSlice not exists")
	}
	w, ok := m.MapWithArray["array"]
	if ok {
		if len(w) != 2 {
			t.Error("The length of MapWithArray[array] is not 2")
		}
	} else {
		t.Error("The key \"array\" in MapWithArray not exists")
	}
	q, ok := m.MapWithIntKey[-1]
	if ok {
		if q == "" {
			t.Error("The value of MapWithIntKey[-1] is empty")
		}
	} else {
		t.Error("The key \"-1\" in MapWithIntKey not exists")
	}
	ñ, ok := m.MapWithInt8Key[-1]
	if ok {
		if ñ == "" {
			t.Error("The value of MapWithInt8Key[-1] is empty")
		}
	} else {
		t.Error("The key \"-1\" in MapWithInt8Key not exists")
	}
	if len(m.MapWithInt64PtrKey) == 0 {
		t.Error("The MapWithInt64PtrKey is empty")
	} else {
		for _, v := range m.MapWithInt64PtrKey {
			if v == "" {
				t.Error("The value of MapWithInt64PtrKey[-1] is empty")
			}
		}
	}
	y, ok := m.MapWithUintKey[1]
	if ok {
		if y == "" {
			t.Error("The value of MapWithUintKey[1] is empty")
		}
	} else {
		t.Error("The key \"1\" in MapWithUintKey not exists")
	}
	bb, ok := m.MapWithUint8Key[1]
	if ok {
		if bb == "" {
			t.Error("The value of MapWithUint8Key[1] is empty")
		}
	} else {
		t.Error("The key \"1\" in MapWithUint8Key not exists")
	}
	if len(m.MapWithUint32PtrKey) == 0 {
		t.Error("The MapWithUint32PtrKey is empty")
	} else {
		for _, v := range m.MapWithUint32PtrKey {
			if v == "" {
				t.Error("The value of MapWithUint32PtrKey[1] is empty")
			}
		}
	}
	o, ok := m.MapWithFloatKey[3.14]
	if ok {
		if o == "" {
			t.Error("The value of MapWithFloatKey[3.14] is empty")
		}
	} else {
		t.Error("The key \"3.14\" in MapWithFloatKey not exists")
	}
	b, ok := m.MapWithBooleanKey[true]
	if ok {
		if b == "" {
			t.Error("The value of MapWithFloatKey[true] is empty")
		}
	} else {
		t.Error("The key \"true\" in MapWithFloatKey not exists")
	}
	uuid := UUID{17, 229, 191, 45, 62, 64, 58, 140, 134, 116, 0, 35, 223, 254, 83, 80}
	uu, ok := m.MapWithCustomKey[uuid]
	if ok {
		if uu == "" {
			t.Error("The value of MapWithFloatKey[11e5bf2d3e403a8c86740023dffe5350] is empty")
		}
	} else {
		t.Error("The key \"11e5bf2d3e403a8c86740023dffe5350\" in MapWithCustomKey not exists")
	}
	for k, v := range m.MapWithCustomKeyPointer {
		if k.String() != uuid.String() {
			t.Error("The key in MapWithCustomKeyPointer is not 11e5bf2d3e403a8c86740023dffe5350")
		} else if v == "" {
			t.Error("The value of MapWithCustomKeyPointer[11e5bf2d3e403a8c86740023dffe5350] is empty")
		}
	}
	for k, v := range m.MapWithStruct1Key {
		if k.IsZero() {
			t.Error("The key of MapWithStruct1Key is zero")
		}
		if v == "" {
			t.Error("The value of MapWithStruct1Key[time.Time] is empty")
		}
	}
	for k, v := range m.MapWithStruct2Key {
		if k.String() == "" {
			t.Error("The key of MapWithStruct2Key is empty")
		}
		if v == "" {
			t.Error("The value of MapWithStruct2Key[ur.URL] is empty")
		}
	}

	// unmarshalText
	if m.UnmarshalTextString == unmarshalTextString {
		t.Error("The value of UnmarshalTextString is not correct. It should not to contain the text of the const unmarshalTextString")
	}
	if m.UnmarshalTextUUID.String() != uuid.String() {
		t.Errorf("The value of UnmarshalTextUUID is not 11e5bf2d3e403a8c86740023dffe5350 but %s", m.UnmarshalTextUUID.String())
	}

	// tag
	if m.Tag == "" {
		t.Error("The value of UnmarshalTextString is empty")
	}

	// time
	if m.Time.IsZero() {
		t.Error("The value of Time is zero")
	}

	// interface
	if v, ok := m.Interface.(string); !ok {
		t.Error("The Interface is not string")
	} else if v == "" {
		t.Error("The value of Interface is empty")
	}
	if v, ok := m.InterfaceStruct.(*InterfaceStruct); !ok {
		t.Error("The InterfaceStruct is not InterfaceStruct struct")
	} else {
		if v.ID == 0 {
			t.Error("The value of InterfaceStruct.ID is 0")
		}
		if v.Name == "" {
			t.Error("The value of InterfaceStruct.Name is empty")
		}
	}
	// custom type
	if m.CustomType != "value changed by custom type" {
		t.Error("The value of CustomType is not correct")
	}
	if m.Time1.IsZero() {
		t.Error("The value of Time1 is not correct")
	}
	if m.Time2.IsZero() {
		t.Error("The value of Time2 is not correct")
	}
	if m.TimeDefault.IsZero() {
		t.Error("The value of TimeDefault is not correct")
	}
}

type TestSlice []string

var sliceValues = url.Values{
	"[0]": []string{"spanish"},
	"[1]": []string{"english"},
}

func TestDecodeInSlice(t *testing.T) {
	var t2 TestSlice
	err := formam.Decode(sliceValues, &t2)
	if err != nil {
		t.Error(err)
	}
}

func TestIgnoreUnknownKeys(t *testing.T) {
	s := struct {
		Name  string `formam:"Name"`
		Map   map[string]string
		Slice []string
	}{}
	vals := url.Values{
		"Name":      []string{"Homer"},
		"City":      []string{"Springfield"},
		"Children.": []string{"Bart", "Lisa"},
		"Job[":      []string{"Safety inspector"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		IgnoreUnknownKeys: true,
	})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}
	if s.Name != "Homer" {
		t.Errorf("Expected Homer got %s", s.Name)
	}
}

func TestIgnoreUnknownKeysDoubleEmbedded(t *testing.T) {
	type Embedded0 struct {
		Foo0 string
		Bar0 string
	}

	type Embedded1 struct {
		Foo1 string
		Bar1 string
	}

	s := struct {
		Embedded0
		Embedded1
		Name  string `formam:"Name"`
		Map   map[string]string
		Slice []string
	}{}

	vals := url.Values{
		"Name":      []string{"Homer"},
		"Foo0":      []string{"Bar"},
		"Bar0":      []string{"Foo"},
		"Foo1":      []string{"Bar"},
		"Bar1":      []string{"Foo"},
		"City":      []string{"Springfield"},
		"Children.": []string{"Bart", "Lisa"},
		"Job[":      []string{"Safety inspector"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		IgnoreUnknownKeys: true,
	})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}
	if s.Name != "Homer" {
		t.Errorf("Expected Homer got %s", s.Name)
	}
	if s.Embedded0.Foo0 != "Bar" {
		t.Errorf("Expected Bar got %s", s.Embedded0.Foo0)
	}
	if s.Embedded0.Bar0 != "Foo" {
		t.Errorf("Expected Foo got %s", s.Embedded0.Bar0)
	}
	if s.Embedded1.Foo1 != "Bar" {
		t.Errorf("Expected Bar got %s", s.Embedded1.Foo1)
	}
	if s.Embedded1.Bar1 != "Foo" {
		t.Errorf("Expected Foo got %s", s.Embedded1.Bar1)
	}
}

func TestIgnoreBracketedKeysError(t *testing.T) {
	s := struct {
		Name string `formam:"Name"`
	}{}
	vals := url.Values{
		"Name":      []string{"Homer"},
		"[Wife]":    []string{"Marge"},
		"His[Wife]": []string{"Marge"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{})
	err := dec.Decode(vals, &s)
	if err == nil {
		t.Error("error is not nil")
	}
}

func TestIgnoreBracketedKeysIgnoreError(t *testing.T) {
	s := []string{}
	vals := url.Values{
		"Name":      []string{"Homer"},
		"[Wife]":    []string{"Marge"},
		"His[Wife]": []string{"Marge"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		IgnoreUnknownKeys: true,
	})
	err := dec.Decode(vals, &s)
	if err == nil {
		t.Error("error is not nil")
	}
}

func TestIgnoreBracketedKeysIgnoreStruct(t *testing.T) {
	s := struct {
		Name string `formam:"Name"`
	}{}
	vals := url.Values{
		"[Wife]":    []string{"Marge"},
		"Name":      []string{"Homer"},
		"His[Wife]": []string{"Marge"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		IgnoreUnknownKeys: true,
	})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}
	if s.Name != "Homer" {
		t.Errorf("Expected Homer got %s", s.Name)
	}
}

func TestEmptyString(t *testing.T) {
	s := struct {
		Name string
	}{
		Name: "Homer",
	}
	vals := url.Values{
		"Name": []string{""},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}
	if s.Name == "Homer" {
		t.Errorf("Expected empty string got %s", s.Name)
	}
}

func TestIgnoredStructTag(t *testing.T) {
	s := struct {
		Name     string `formam:"-"`
		Phone    string
		LastName string `formam:"lastname"`
	}{
		Name:     "Homer",
		LastName: "Bouvier",
		Phone:    "555-555-111",
	}
	vals := url.Values{
		"Name":     []string{"Marge"},
		"lastname": []string{"Simpson"},
		"Phone":    []string{"555-333-222"},
	}

	dec := formam.NewDecoder(&formam.DecoderOptions{})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}

	if s.Name != "Homer" {
		t.Errorf("Expected Homer got %s", s.Name)
	}

	if s.LastName != "Simpson" {
		t.Errorf("Expected LastName is Simpson but got %s", s.LastName)
	}

	if s.Phone != "555-333-222" {
		t.Errorf("Expected new phone number '555-333-222' but got %s", s.Phone)
	}
}

// Test for #7
func TestPanic(t *testing.T) {
	s := struct {
		Foo []string `formam:"foo"`
	}{}

	vals := url.Values{
		"foo[]": []string{"5", "6"},
	}

	dec := formam.NewDecoder(&formam.DecoderOptions{})
	err := dec.Decode(vals, &s)
	if err != nil {
		t.Error(err)
	}

	if len(s.Foo) != 2 {
		t.Fatalf("len(s.Foo) is %d", len(s.Foo))
	}
	if s.Foo[0] != "5" {
		t.Errorf("s.Foo[0] is %#v", s.Foo[0])
	}
	if s.Foo[1] != "6" {
		t.Errorf("s.Foo[1] is %#v", s.Foo[1])
	}

}

// #25
func TestOverflow(t *testing.T) {
	s := struct {
		N uint8
	}{}

	vals := url.Values{
		"N": []string{"300"},
	}

	dec := formam.NewDecoder(&formam.DecoderOptions{})
	err := dec.Decode(vals, &s)
	if err == nil {
		t.Fatalf("error is nil")
	}
	if s.N != 0 {
		t.Fatalf("s.N is %d", s.N)
	}

	fErr := err.(*formam.Error)
	if fErr.Code() != formam.ErrCodeRange {
		t.Fatalf("error code is %d", fErr.Code())
	}
}

// #28
func TestArrayIgnore(t *testing.T) {
	type Users []struct {
		Email string `form:"Email"`
	}

	vals := url.Values{
		"[0].Email":          {"a@a.a"},
		"[1].Email":          {"b@b.b"},
		"authenticity_token": {"xxx"},
	}

	var u Users
	dec := formam.NewDecoder(&formam.DecoderOptions{
		IgnoreUnknownKeys: true,
	})
	err := dec.Decode(vals, &u)
	if err != nil {
		t.Fatal(err)
	}

	out := fmt.Sprintf("%v", u)
	want := `[{a@a.a} {b@b.b}]`
	if out != want {
		t.Fatalf("\nout:  %s\nwant: %s", out, want)
	}
}

// #31
func TestArrayLength(t *testing.T) {
	tests := []struct {
		arrayLen  int
		maxSize   int
		wantError string
	}{
		{15000, 0, ""},
		{15999, 0, ""},
		{16000, 0, "array size 16001 is longer than MaxSize 16000"},
		{8, 8, "array size 9 is longer than MaxSize 8"},

		{26000, -1, ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			type Malicious []struct {
				X string
			}

			vals := url.Values{
				fmt.Sprintf("[%d].X", tt.arrayLen): {"yeah nah"},
			}

			var m Malicious
			dec := formam.NewDecoder(&formam.DecoderOptions{MaxSize: tt.maxSize})
			err := dec.Decode(vals, &m)

			if !errorContains(err, tt.wantError) {
				t.Fatalf("wrong error: %s", err)
			}

			if err == nil {
				if len(m) != tt.arrayLen+1 {
					t.Errorf("%d != %d", len(m), tt.arrayLen+1)
				}
			} else {
				if len(m) != 0 {
					t.Errorf("len(n) is not 0: %d", len(m))
				}
			}
		})
	}
}

type StringWithUnmarshalText string

func (s *StringWithUnmarshalText) UnmarshalText(data []byte) error {
	*s = "string by UnmarshalText"
	return nil
}

func TestEnableUnmarshalText(t *testing.T) {
	s := struct {
		Text StringWithUnmarshalText
	}{}
	vals := url.Values{
		"Text": []string{"string by form"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{})
	if err := dec.Decode(vals, &s); err != nil {
		t.Fatalf("error when decode %s", err)
	}
	if s.Text != "string by UnmarshalText" {
		t.Errorf("the UnmarshalText hasn't been called")
	}
}

func TestDisableUnmarshalText(t *testing.T) {
	s := struct {
		Text StringWithUnmarshalText
	}{}
	vals := url.Values{
		"Text": []string{"string by form"},
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		DisableUnmarshalText: true,
	})
	if err := dec.Decode(vals, &s); err != nil {
		t.Fatalf("error when decode %s", err)
	}
	if s.Text != "string by form" {
		t.Errorf("the UnmarshalText has been called")
	}
}

type AppendString string

func (s *AppendString) UnmarshalText(data []byte) error {
	*s += AppendString(data) + "-"
	return nil
}

func TestUnmarshalTextList(t *testing.T) {
	var s struct{ Append AppendString }

	// No [], so use first value.
	err := formam.Decode(url.Values{
		"Append": []string{"banana", "chocolate"},
	}, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Append != "banana-" {
		t.Error(s.Append)
	}

	// [], so call for every value.
	s.Append = ""
	err = formam.Decode(url.Values{
		"Append[]": []string{"banana", "chocolate"},
	}, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Append != "banana-chocolate-" {
		t.Error(s.Append)
	}

	// Append to existing
	err = formam.Decode(url.Values{
		"Append[]": []string{"cinnamon"},
	}, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Append != "banana-chocolate-cinnamon-" {
		t.Error(s.Append)
	}
}

// #40
func TestArrayIndexIsOutOfBoundsError(t *testing.T) {
	s := struct {
		Array [2]string
	}{}

	dec := formam.NewDecoder(nil)
	err := dec.Decode(url.Values{
		"Array[2]": []string{"10"},
	}, &s)
	if !errorContains(err, "array index is out of bounds") {
		t.Fatalf("wrong error: %s", err)
	}
	if s.Array[0] != "" || s.Array[1] != "" {
		t.Errorf("wrong value: %#v", s.Array)
	}

	err = dec.Decode(url.Values{
		"Array[1]": []string{"10"},
	}, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Array[0] != "" || s.Array[1] != "10" {
		t.Errorf("wrong value: %#v", s.Array)
	}
}

// errorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// #42
func TestMapToPtrStruct(t *testing.T) {
	s := struct {
		M map[string]*struct {
			ID string
		}
	}{}

	vals := url.Values{
		"M[key].ID": []string{"M[key].ID"},
	}

	dec := formam.NewDecoder(nil)
	if err := dec.Decode(vals, &s); err != nil {
		t.Fatalf("error when decode %s", err)
	}

	v, ok := s.M["key"]
	if !ok {
		t.Fatal("The key \"key\" in M does not exists")
	}
	if v.ID != "M[key].ID" {
		t.Errorf("The value in key \"key\" of M is incorrect: %q", v.ID)
	}
}

func TestBracketsAndNestedPointers(t *testing.T) {
	var s struct {
		MapStringString    map[string]string
		MapStringPtrStruct map[string]struct {
			ID string
		}
		MapStringMapStringString map[string]map[string]string
	}

	vals := url.Values{
		"MapStringString[a[b][c]d]":             []string{"MapStringString[a[b][c]d]"},
		"MapStringString[name.with.dots]":       []string{"MapStringString[name.with.dots]"},
		"MapStringPtrStruct[k2]ID":              []string{"MapStringPtrStruct[k2]ID"},
		"MapStringMapStringString[a[b[c]d]]q]w": []string{"MapStringMapStringString[a[b[c]d]]q]w"},
	}

	dec := formam.NewDecoder(nil)

	if err := dec.Decode(vals, &s); err != nil {
		t.Fatalf("error when decode %s", err)
	}

	if v, ok := s.MapStringString["a[b][c]d"]; !ok {
		t.Error("The key \"a[b][c]d\" in MapStringString does not exists")
	} else if v != "MapStringString[a[b][c]d]" {
		t.Error("The value in key \"a[b][c]d\" of MapStringString is incorrect")
	}

	if v, ok := s.MapStringString["name.with.dots"]; !ok {
		t.Error("The key \"name.with.dots\" in MapStringString does not exists")
	} else if v != "MapStringString[name.with.dots]" {
		t.Error("The value in key \"name.with.dots\" of MapStringString is incorrect")
	}

	if v, ok := s.MapStringPtrStruct["k2"]; !ok {
		t.Error("The key \"k2\" in MapStringPtrStruct does not exists")
	} else if v.ID != "MapStringPtrStruct[k2]ID" {
		t.Error("The value in key \"k2\" of MapStringPtrStruct is incorrect")
	}

	if v, ok := s.MapStringMapStringString["a[b[c]d]"]; !ok {
		t.Error("The key \"a[b[c]d]\" in MapStringMapStringString does not exists")
	} else if vv, ok := v["q]w"]; !ok {
		t.Error("The key \"q]w\" in MapStringMapStringString[a[b[c]d]] does not exists")
	} else if vv != "MapStringMapStringString[a[b[c]d]]q]w" {
		t.Error("The value in key \"q]w\" of MapStringMapStringString[a[b[c]d]] is incorrect")
	}
}
