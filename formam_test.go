package formam

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/ajg/form"
	"net/url"
)

type Test struct {
	Nest struct {
		Children []struct {
			Id string
			Lol string
		}
	}
	Mierda string `formam:"mierda" form:"mierda"`
	Map map[string]string
}

func TestDecode(t *testing.T) {
	//req, _ := http.NewRequest("POST", "http://www.monoculum.com/search?main=foo&childs[0]=bar&childs[1]=buz&nest.childs[0].id=lol", strings.NewReader("z=post&both=y"))
	req, _ := http.NewRequest("POST", "http://www.monoculum.com/search?Nest.Children[0].Id=lol&Nest.Children[0].Lol=lol&mierda=cojonudo&Map.es_es=titanic", strings.NewReader("z=post&both=y"))
	req.Header.Set("Content-Type", "application/x-www-form-encoded; param=value");

	test := &Test{}
	decoder, err := NewDecoder(req, test)
	if err != nil {
		t.Error(err)
	}
	if err := decoder.Decode(); err != nil {
		t.Error(err)
	}
	fmt.Println("FINALIZADO DECODE: ", test)
}



var (
	values = url.Values{
		"Nest.Children.0.Id": []string{"lol"},
		"Nest.Children.0.Lol": []string{"lol"},
		"Map.es_Es": []string{"coooño"},
		"mierda": []string{"cojonudo"},
	}
)

func BenchmarkAJGForm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ne := new(Test)
		if err := form.DecodeValues(ne, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ne := new(Test)
		dec := schema.NewDecoder()
		if err := dec.Decode(ne, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFormam(b *testing.B) {
	req, _ := http.NewRequest("POST", "http://www.monoculum.com/search?Nest.Children[0].Id=lol&Nest.Children[0].Lol=lol&mierda=cojonudo&Map.es_es=titanic", strings.NewReader("z=post&both=y"))
	req.Header.Set("Content-Type", "application/x-www-form-encoded; param=value");

	for i := 0; i < b.N; i++ {
		test := new(Test)
		decoder, err := NewDecoder(req, test)
		if err != nil {
			b.Error(err)
		}
		if err := decoder.Decode(); err != nil {
			b.Error(err)
		}
	}
}
