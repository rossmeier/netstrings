package netstrings

import (
	"strings"
	"testing"
)

func TestNetstrings(t *testing.T) {

	var encoded = [][]byte{
		[]byte("Hello world!"),
		[]byte(""),
		[]byte("Fo,o"),
		[]byte("☎"),
	}

	out, err := Encode(encoded...)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if string(out) != "12:Hello world!,0:,4:Fo,o,3:☎," {
		t.Error(string(out))
		t.Fail()
	}

	decoded, err := Decode(out)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !equal(encoded, decoded) {
		t.Error(encoded)
		t.Error(decoded)
		t.Fail()
	}

	decoded, err = Decode([]byte(strings.Replace(string(out), ",", "", 2)))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !equal(encoded, decoded) {
		t.Error(encoded)
		t.Error(decoded)
		t.Fail()
	}

}

func equal(x, y [][]byte) bool {

	if x == nil && y == nil {
		return true
	}

	if x == nil || y == nil {
		return false
	}

	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if string(x[i]) != string(y[i]) {
			return false
		}
	}

	return true
}
