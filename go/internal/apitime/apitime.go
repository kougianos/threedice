// Package apitime formats timestamps the way java.time.Instant does.
package apitime

import "time"

// Time marshals as ISO-8601 with fractional seconds in groups of 3 digits
// (0, 3, 6 or 9), matching java.time.Instant.toString().
//
// Go's RFC3339Nano trims every trailing zero, so an instant at .100Z would go
// out as .1Z and one at .075360Z as .07536Z, where Jackson emits .100Z and
// .075360Z. Both are valid ISO-8601, but they are not the same bytes.
type Time struct {
	time.Time
}

func New(t time.Time) Time { return Time{t} }

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.format() + `"`), nil
}

func (t Time) format() string {
	u := t.UTC()
	switch nanos := u.Nanosecond(); {
	case nanos == 0:
		return u.Format("2006-01-02T15:04:05Z")
	case nanos%1_000_000 == 0:
		return u.Format("2006-01-02T15:04:05.000Z")
	case nanos%1_000 == 0:
		return u.Format("2006-01-02T15:04:05.000000Z")
	default:
		return u.Format("2006-01-02T15:04:05.000000000Z")
	}
}
