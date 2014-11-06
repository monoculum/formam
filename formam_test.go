package formam

import (
	"net/url"
	"testing"
	"fmt"
)

type Anonymous struct {
	Int int `formam:"int"`
}

type Test struct {
	Nest struct {
		Children []struct {
			Id string
			Name string
		}
	}
	String string
	Slice  []int
	Map map[string][]string
	Bool bool
	Ptr  *string
	Tag  string `formam:"tag"`
	Anonymous
}

var valuesFormam = url.Values{
	"Nest.Children[0].Id": []string{"monoculum_id"},
	"Nest.Children[0].Name": []string{"Monoculum"},
	"Map.es_Es[0]": []string{"javier"},
	"Map.es_Es[1]": []string{"javier"},
	"Map.es_Es[2]": []string{"javier"},
	"Map.es_Es[3]": []string{"javier"},
	"Map.es_Es[4]": []string{"javier"},
	"Map.es_Es[5]": []string{"javier"},
	"String": []string{"cojonudo"},
	"Slice[0]": []string{"1"},
	"Slice[1]": []string{"2"},
	"int": []string{"1"},
	"Bool": []string{"true"},
	"tag": []string{"tagged"},
}


func TestDecode(t *testing.T) {
	test := &Test{}
	err := Decode(valuesFormam, test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("END: ", test)
}
