package validation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var v = validator.New()

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func init() {
	v.RegisterValidation("decimal_gt_zero", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(decimal.Decimal)
		if !ok {
			return false
		}

		return value.GreaterThan(decimal.Zero)
	})
}

func Validate[T any](w http.ResponseWriter, input T) bool {
	if err := v.Struct(input); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			errs := make([]ValidationError, len(valErrs))

			for i, fieldErr := range valErrs {
				errs[i] = ValidationError{
					Field: fieldErr.Field(),
					Tag:   fieldErr.Tag(),
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"errors": errs,
			})

			return false
		}

		http.Error(
			w,
			`{"error":"Validation failed"}`,
			http.StatusInternalServerError,
		)

		return false
	}

	return true
}