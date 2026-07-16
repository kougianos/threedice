// Package money renders monetary values the way the Java service does.
package money

import "github.com/shopspring/decimal"

// Amount wraps decimal.Decimal to always marshal at two decimal places.
//
// shopspring's own MarshalJSON quotes the value and trims trailing zeros, so a
// balance of 1040.00 would go out as "1040" rather than the 1040.00 that
// Jackson produces for a BigDecimal read from NUMERIC(12,2).
type Amount struct {
	decimal.Decimal
}

func New(d decimal.Decimal) Amount { return Amount{d} }

func FromString(s string) (Amount, error) {
	d, err := decimal.NewFromString(s)
	return Amount{d}, err
}

var Zero = Amount{decimal.Zero}

func (a Amount) MarshalJSON() ([]byte, error) {
	return []byte(a.StringFixed(2)), nil
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	return a.Decimal.UnmarshalJSON(b)
}
