package formam

import (
	"fmt"
	"net/url"
	"testing"
)

type Field string

type Laboratory struct {
	MapWithStruct3Key map[Field]string
	Jodete            Field
}

var valss = url.Values{
	"Jodete": []string{"2006-01-02"},
	//"MapWithStruct3Key[ID.ID]": []string{"struct key in map"},
}

func TestLaboratory(t *testing.T) {
	var m Laboratory
	dec := NewDecoder(nil)
	dec.RegisterCustomType(func(vals []string) (interface{}, error) {
		return Field("value changed by custom type"), nil
	}, Field(""))
	err := dec.Decode(valss, &m)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("RESULT1: ", m)
}
