package requiredjson

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Example struct {
		A int    `json:"a,required"`
		B string `json:"b,required"`
		C int    `json:"c,omitempty"`
		D string `json:"d,omitempty"`
	}
	for _, tt := range []struct {
		raw      string
		expected Example
		valid    bool
	}{
		{`{"a": 1, "b":"asdf", "c": 2, "d":"fdsa"}`, Example{A: 1, B: "asdf", C: 2, D: "fdsa"}, true},
		{`{"a": 1, "b":"asdf", "d":""}`, Example{A: 1, B: "asdf", D: ""}, true},
		{`{"a": 1, "b":"", "c": 2}`, Example{A: 1, B: "", C: 2}, true},
		{`{"a": 1}`, Example{A: 1}, false},
		{`{b":""}`, Example{B: ""}, false},
	} {
		var e Example
		err := Unmarshal([]byte(tt.raw), &e)
		if err != nil && tt.valid {
			t.Errorf("Unwanted error, unmarshal %s: %v", tt.raw, err)
		}
		if err == nil && !tt.valid {
			t.Errorf("Wanted an error, raw %s", tt.raw)
		}
	}
}
