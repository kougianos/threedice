// Package domain holds the dice rules: rolling, valid products, and odds.
package domain

import (
	"crypto/rand"
	"math/big"

	"github.com/shopspring/decimal"
)

const diceSides = 6

// Roller produces a single die value. Production uses SecureRoller; tests
// substitute a deterministic implementation, mirroring the Java suite's
// @MockitoBean on DiceEngine.
type Roller interface {
	Roll() int32
}

// SecureRoller draws from a cryptographically secure source, matching the
// Java service's use of java.security.SecureRandom.
type SecureRoller struct{}

// Roll returns a value between 1 and 6 inclusive.
func (SecureRoller) Roll() int32 {
	n, err := rand.Int(rand.Reader, big.NewInt(diceSides))
	if err != nil {
		// The system CSPRNG is unavailable; there is no safe way to settle a bet.
		panic("dice: secure random source unavailable: " + err.Error())
	}
	return int32(n.Int64()) + 1
}

// validProducts holds the 40 values reachable as a product of three six-sided
// dice, out of the 216 combinations.
var validProducts = func() map[int32]struct{} {
	m := make(map[int32]struct{})
	for a := int32(1); a <= diceSides; a++ {
		for b := int32(1); b <= diceSides; b++ {
			for c := int32(1); c <= diceSides; c++ {
				m[a*b*c] = struct{}{}
			}
		}
	}
	return m
}()

// IsValidProduct reports whether value is achievable as a product of three dice.
func IsValidProduct(value int32) bool {
	_, ok := validProducts[value]
	return ok
}

// ValidProductCount is exposed for tests asserting the 40-of-216 invariant.
func ValidProductCount() int { return len(validProducts) }

var (
	oddsLow  = decimal.NewFromInt(2)
	oddsMid  = decimal.NewFromInt(5)
	oddsHigh = decimal.NewFromInt(2)
)

// Odds returns the payout multiplier for a product value:
//
//	Product < 9        -> x2
//	9 <= Product < 120 -> x5
//	Product >= 120     -> x2
func Odds(productValue int32) decimal.Decimal {
	switch {
	case productValue < 9:
		return oddsLow
	case productValue < 120:
		return oddsMid
	default:
		return oddsHigh
	}
}
