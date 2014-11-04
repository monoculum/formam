package formam

import (
	"testing"
	"net/http"
	"strings"
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
	Map map[string]map[string]struct{
		Id string
}
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
	Map map[string]string
	Int int `form:"int"`
	Bool bool
	Anony
}

var urlStr = "http://www.monoculum.com/search?Nest.Children[0].Id=lol&Nest.Children[0].Lol=lol&mierda=cojonudo&Map.es_es.sii.Id=javier&Slice[0]=1&Slice[1]=2&int=20&Bool=true"

func TestDecode(t *testing.T) {
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader("z=post&both=y"))
	req.Header.Set("Content-Type", "application/x-www-form-encoded; param=value");
	req.ParseForm()

	test := &Test{}
	err := Decode(req.Form, test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("FINALIZADO DECODE: ", test)
}



var (
	values = url.Values{
		"Nest.Children.0.Id": []string{"lol"},
		"Nest.Children.0.Lol": []string{"lol"},
		"Map.es_Es.Id": []string{"javier"},
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
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader("z=post&both=y"))
	req.Header.Set("Content-Type", "application/x-www-form-encoded; param=value");
	req.ParseForm()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test := new(Test)
		if err := Decode(req.Form, test); err != nil {
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
		"Map": {"es_Es": {"Id": "emilio"}},
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
