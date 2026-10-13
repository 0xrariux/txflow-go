package decimal

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"
)

func TestPrecisionAndCanonicalization(t *testing.T) {
	for _, tc := range []struct{ in, want string }{{"-000.000", "0"}, {"00012.3400", "12.34"}, {"0.000000000000000001", "0.000000000000000001"}, {"9007199254740993", "9007199254740993"}} {
		d, err := Parse(tc.in)
		if err != nil || d.String() != tc.want {
			t.Fatal(tc, err, d)
		}
		b, _ := json.Marshal(d)
		var copy Decimal
		if err := json.Unmarshal(b, &copy); err != nil || copy.Cmp(d) != 0 {
			t.Fatal(string(b), err)
		}
	}
	a := MustParse("9007199254740993.000000000000000001")
	b := MustParse("0.000000000000000001")
	if a.Add(b).Sub(a).Cmp(b) != 0 {
		t.Fatal("lost precision")
	}
	if got := b.Half().String(); got != "0.0000000000000000005" {
		t.Fatal(got)
	}
	if a.String() != "9007199254740993.000000000000000001" {
		t.Fatal("mutation")
	}
	if (Decimal{}).String() != "0" {
		t.Fatal("zero value")
	}
}
func TestRejectMalformedNumbers(t *testing.T) {
	for _, s := range []string{"", "1e2", "NaN", "Infinity", " 1", "1 ", "+1", ".1", "1.", "--1", "1.2.3", strings.Repeat("9", 4097)} {
		if _, err := Parse(s); err == nil {
			t.Fatal(s)
		}
	}
	for _, s := range []string{"null", "1", "true", "`bad`"} {
		var d Decimal
		if err := json.Unmarshal([]byte(s), &d); err == nil {
			t.Fatal(s)
		}
	}
	d := MustParse("7")
	if err := json.Unmarshal([]byte(`"invalid"`), &d); err == nil || d.String() != "7" {
		t.Fatal("failed decode modified value")
	}
}
func FuzzDecimal(f *testing.F) {
	for _, s := range []string{"0", "-12.500", "0.000000000000000001", "9007199254740993", "1e9"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		d, err := Parse(s)
		if err != nil {
			return
		}
		r, ok := new(big.Rat).SetString(s)
		if !ok {
			t.Fatal("accepted invalid rational")
		}
		if d.rat().Cmp(r) != 0 {
			t.Fatal("changed number")
		}
		if d.Half().Add(d.Half()).Cmp(d) != 0 {
			t.Fatal("inexact half")
		}
		raw, _ := json.Marshal(d)
		var out Decimal
		if err := json.Unmarshal(raw, &out); err != nil || out.Cmp(d) != 0 {
			t.Fatal("round trip")
		}
	})
}
