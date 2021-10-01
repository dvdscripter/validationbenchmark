package validationbenchmark

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type Valid interface {
	isValid() error
}

type PaymentValidByTag struct {
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

type PaymentValidTypeSafe struct {
	Amount decimal.Decimal `json:"amount"`
}

func (p PaymentValidTypeSafe) Validate() error {
	if !p.Amount.IsPositive() {
		return errors.New("amout must be positive")
	}

	return nil
}

func (p PaymentValidTypeSafe) isValid() error {
	return p.Validate()
}

func requiredGT0(field reflect.Value) interface{} {
	fieldDecimal, ok := field.Interface().(decimal.Decimal)
	if !ok {
		return nil
	}
	if !fieldDecimal.IsPositive() {
		return nil
	}
	v, _ := fieldDecimal.Float64()
	return v
}

func (p PaymentValidByTag) isValid() error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(requiredGT0, decimal.Decimal{})

	return validate.Struct(p)
}

func Validate(data Valid) error {
	return data.isValid()
}

func createByTag(w http.ResponseWriter, r *http.Request) {
	payment := PaymentValidByTag{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := Validate(payment); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}

func createByTypeSafe(w http.ResponseWriter, r *http.Request) {
	payment := PaymentValidTypeSafe{}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := Validate(payment); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}
