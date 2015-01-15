package formam

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

type Anonymous struct {
	Int int `formam:"int"`
}

type PtrStruct struct {
	String *string
}

type Test struct {
	Nest struct {
		Children []struct {
			Id   string
			Name string
		}
	}
	String          string
	Slice           []int
	MapSlice        map[string][]string
	MapMap          map[string]map[string]string
	MapMapMapStruct map[string]map[string]map[string]map[string]struct {
		Recursive bool
	}
	Bool bool
	Ptr  *string
	Tag  string `formam:"tag"`
	Anonymous
	Time time.Time
	URL  url.URL
	PtrStruct *PtrStruct
}

var valuesFormam = url.Values{
	"Nest.Children[0].Id":                                []string{"monoculum_id"},
	"Nest.Children[0].Name":                              []string{"Monoculum"},
	"MapSlice.names[0]":                                  []string{"shinji"},
	"MapSlice.names[2]":                                  []string{"sasuka"},
	"MapSlice.names[4]":                                  []string{"carla"},
	"MapSlice.countries[0]":                              []string{"japan"},
	"MapSlice.countries[1]":                              []string{"spain"},
	"MapSlice.countries[2]":                              []string{"germany"},
	"MapSlice.countries[3]":                              []string{"united states"},
	"MapMap.titles.es_es":                                []string{"El viaje de Chihiro"},
	"MapMap.titles.en_us":                                []string{"The spirit away"},
	"MapMapMapStruct.map.struct.are.recursive.Recursive": []string{"true"},
	"Slice[0]": []string{"1"},
	"Slice[1]": []string{"2"},
	"int":      []string{"1"}, // Int is located inside Anonymous struct
	"Bool":     []string{"true"},
	"tag":      []string{"tagged"},
	"Ptr":      []string{"this is a pointer to string"},
	"Time":     []string{"2006-10-08"},
	"URL":      []string{"https://www.golang.org"},
	"PtrStruct.String": []string{"dashaus"},
}

func TestDecode(t *testing.T) {
	test := &Test{}
	err := Decode(valuesFormam, test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("END: ", test)
}
