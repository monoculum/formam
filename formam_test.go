package formam

import (
	"testing"
	"fmt"
	"github.com/ajg/form"
	"net/url"
	"encoding/json"
)

type Anony struct {
	Int int `formam:"int" json:"int"`
}

type Test struct {
	Nest struct {
		Children []struct {
			Id string
			Lol string
		}
	}
	Mierda string `formam:"mierda" form:"mierda"`
	Slice  []int
	Map map[string][]string
	Bool bool
	Anony
}

type Test1 struct {
	Nest struct {
		Children []struct {
			Id string
			Lol string
		}
	}
	Mierda string `formam:"mierda" form:"mierda"`
	Slice  []int
	Map map[string]struct{
	Id string
	Name string
	Type string
	Class string
}
	Int int `form:"int"`
	Bool bool
	Anony
}


func TestDecode(t *testing.T) {
	test := &Test{}
	err := Decode(valuesFormam, test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("FINALIZADO DECODE: ", test)
}


var (
	valuesFormam = url.Values{
		"Nest.Children[0].Id": []string{"lol"},
		"Nest.Children[0].Lol": []string{"lol"},
		"Map.es_Es[0]": []string{"javier"},
		"Map.es_Es[1]": []string{"javier"},
		"Map.es_Es[2]": []string{"javier"},
		"Map.es_Es[3]": []string{"javier"},
		"Map.es_Es[4]": []string{"javier"},
		"Map.es_Es[5]": []string{"javier"},
		"mierda": []string{"cojonudo"},
		"Slice[0]": []string{"1"},
		"Slice[1]": []string{"2"},
		"int": []string{"1"},
		"Bool": []string{"true"},
	}
	values = url.Values{
		"Nest.Children.0.Id": []string{"lol"},
		"Nest.Children.0.Lol": []string{"lol"},
		"Map.es_Es.Id": []string{"javier"},
		"Map.es_Es.Name": []string{"javier"},
		"Map.es_Es.Type": []string{"javier"},
		"Map.es_Es.Class": []string{"javier"},
		"mierda": []string{"cojonudo"},
		"Slice.0": []string{"1"},
		"Slice.1": []string{"2"},
		"int": []string{"1"},
		"Bool": []string{"true"},
	}
)

func BenchmarkAJGForm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ne := new(Test1)
		if err := form.DecodeValues(ne, values); err != nil {
			b.Error(err)
		}
	}
}

/*
func BenchmarkSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ne := new(Test)
		dec := schema.NewDecoder()
		if err := dec.Decode(ne, values); err != nil {
			b.Error(err)
		}
	}
}
*/

func BenchmarkFormam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test := new(Test)
		if err := Decode(valuesFormam, test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJSON(b *testing.B) {
	val := `
	{
		"Nest":
			{
				"Children": [{"Id": "lol", "Lol":"lol"}]
			},
		"Mierda": "cojonudo",
		"Map": {"es_Es": ["emilio", "javier", "ronaldo", "waldo", "javier", "cinco"]},
		"Slice": [1, 2],
		"int": 20,
		"Bool": true
	}
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test := new(Test)
		if err := json.Unmarshal([]byte(val), &test); err != nil {
			b.Error(err)
		}
	}
}
