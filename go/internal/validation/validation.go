// Package validation reproduces the Java service's Bean Validation behaviour:
// the same rules, the same messages, and the same joined detail string.
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	"github.com/kougianos/threedice/go/internal/domain"
	"github.com/kougianos/threedice/go/internal/dto"
	"github.com/kougianos/threedice/go/internal/money"
)

type Validator struct {
	v *validator.Validate
}

func New() (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Report the JSON field name, as Spring reports the record component name.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := v.RegisterValidation("dmin", decimalBound(decimal.Decimal.GreaterThanOrEqual)); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("dmax", decimalBound(decimal.Decimal.LessThanOrEqual)); err != nil {
		return nil, err
	}

	return &Validator{v: v}, nil
}

// messages maps "field.tag" to the exact text the Jakarta annotations produce.
var messages = map[string]string{
	"playerId.required":       "Player ID is required",
	"stake.required":          "Stake is required",
	"stake.dmin":              "Minimum stake is $1.00",
	"stake.dmax":              "Maximum stake is $10,000.00",
	"predictedValue.required": "Predicted value is required",
	"predictedValue.min":      "Predicted value must be at least 1",
	"predictedValue.max":      "Predicted value must be at most 216",
	"idempotencyKey.required": "Idempotency key is required",
	"idempotencyKey.max":      "Idempotency key must be at most 64 characters",
}

// ValidatePlaceBet returns the joined detail string Spring's handler would
// produce, or "" when the request is valid. Field errors are rendered as
// "field: message", sorted, and joined with "; ".
func (val *Validator) ValidatePlaceBet(req dto.PlaceBetRequest) string {
	var details []string

	if err := val.v.Struct(req); err != nil {
		var verrs validator.ValidationErrors
		if !errors.As(err, &verrs) {
			return err.Error()
		}
		for _, fe := range verrs {
			key := fe.Field() + "." + fe.Tag()
			msg, ok := messages[key]
			if !ok {
				msg = fe.Error()
			}
			details = append(details, fe.Field()+": "+msg)
		}
	}

	// Applied outside the tag chain so it stacks with min/max the way
	// @ValidPredictedValue stacks with @Min/@Max. Nulls are left to `required`,
	// matching Bean Validation's skip-on-null.
	if req.PredictedValue != nil && !domain.IsValidProduct(*req.PredictedValue) {
		details = append(details, fmt.Sprintf(
			"predictedValue: Predicted value %d is not a possible product of three dice",
			*req.PredictedValue))
	}

	if len(details) == 0 {
		return ""
	}
	sort.Strings(details)
	return strings.Join(details, "; ")
}

// decimalBound adapts a decimal comparison into a validator func for money.Amount.
func decimalBound(cmp func(decimal.Decimal, decimal.Decimal) bool) validator.Func {
	return func(fl validator.FieldLevel) bool {
		amount, ok := fl.Field().Interface().(money.Amount)
		if !ok {
			return false
		}
		bound, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return cmp(amount.Decimal, bound)
	}
}
