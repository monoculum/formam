package formam

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		in         error
		wantString string
		wantJSON   string
	}{
		{newError(0, "", "", "oh noes"),
			"formam: oh noes",
			`"formam: oh noes"`},
		{newError(0, "foo", "foo.bar", "oh noes"),
			"formam: field=foo; path=foo.bar: oh noes",
			`"formam: field=foo; path=foo.bar: oh noes"`},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			outString := tt.in.Error()
			if outString != tt.wantString {
				t.Errorf("\nout:  %#v\nwant: %#v\n", outString, tt.wantString)
			}

			j, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			outJSON := string(j)

			if outJSON != tt.wantJSON {
				t.Errorf("\nout:  %#v\nwant: %#v\n", outJSON, tt.wantJSON)
			}
		})
	}
}
