// Package decimal provides a narrow immutable finite-decimal value.
// It deliberately has no floating-point conversion, generic division,
// rounding or exponential notation. JSON values must be decimal strings.
package decimal

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type Decimal struct{ text string }

// Parse accepts fixed decimal notation with optional minus sign, up to 4096
// characters. This resource limit is an SDK choice, not a TxFlow rule.
func Parse(s string) (Decimal, error) {
	if len(s) == 0 || len(s) > 4096 {
		return Decimal{}, fmt.Errorf("decimal: invalid length")
	}
	neg := strings.HasPrefix(s, "-")
	digits := s
	if neg {
		digits = s[1:]
	}
	parts := strings.Split(digits, ".")
	if len(parts) > 2 || len(parts[0]) == 0 || (len(parts) == 2 && len(parts[1]) == 0) {
		return Decimal{}, fmt.Errorf("decimal: expected finite decimal notation")
	}
	for _, p := range parts {
		for _, ch := range p {
			if ch < '0' || ch > '9' {
				return Decimal{}, fmt.Errorf("decimal: expected digits")
			}
		}
	}
	whole := strings.TrimLeft(parts[0], "0")
	if whole == "" {
		whole = "0"
	}
	frac := ""
	if len(parts) == 2 {
		frac = strings.TrimRight(parts[1], "0")
	}
	text := whole
	if frac != "" {
		text += "." + frac
	}
	if neg && text != "0" {
		text = "-" + text
	}
	return Decimal{text: text}, nil
}
func MustParse(s string) Decimal {
	d, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return d
}
func (d Decimal) String() string {
	if d.text == "" {
		return "0"
	}
	return d.text
}
func (d Decimal) rat() *big.Rat { r, _ := new(big.Rat).SetString(d.String()); return r }
func (d Decimal) scale() int {
	s := d.String()
	i := strings.IndexByte(s, '.')
	if i < 0 {
		return 0
	}
	return len(s) - i - 1
}
func (d Decimal) Cmp(other Decimal) int { return d.rat().Cmp(other.rat()) }
func (d Decimal) Sign() int             { return d.rat().Sign() }
func normalized(s string) Decimal {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	if s == "-0" {
		s = "0"
	}
	return Decimal{text: s}
}
func (d Decimal) Add(other Decimal) Decimal {
	scale := d.scale()
	if other.scale() > scale {
		scale = other.scale()
	}
	return normalized(new(big.Rat).Add(d.rat(), other.rat()).FloatString(scale))
}
func (d Decimal) Sub(other Decimal) Decimal {
	scale := d.scale()
	if other.scale() > scale {
		scale = other.scale()
	}
	return normalized(new(big.Rat).Sub(d.rat(), other.rat()).FloatString(scale))
}

// Half is exact even for odd smallest units; no financial rounding is applied.
func (d Decimal) Half() Decimal {
	return normalized(new(big.Rat).Quo(d.rat(), big.NewRat(2, 1)).FloatString(d.scale() + 1))
}
func (d Decimal) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }
func (d *Decimal) UnmarshalJSON(b []byte) error {
	var s string
	if len(b) == 0 || b[0] != '"' {
		return fmt.Errorf("decimal: JSON decimal must be a string")
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := Parse(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
