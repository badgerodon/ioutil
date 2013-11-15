package ioutil

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestMultiReadSeeker(t *testing.T) {
	a := "abc"
	b := "def"
	c := "ghi"

	rdr := NewMultiReadSeeker(strings.NewReader(a), strings.NewReader(b), strings.NewReader(c))

	bs, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Errorf("Expected no error got: %v", err)
	}
	if string(bs) != a+b+c {
		t.Errorf("Expected readers to be concatenated got: ", string(bs))
	}
}
