package em

import (
	"context"
	"emsrv/pkg/db"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

/*type FieldError struct {
	Field      string                `json:"field"`
	Error      string                `json:"error"`
	Constraint *FieldErrorConstraint `json:"constraint,omitempty"`
}

type FieldErrorConstraint struct {
	Max int `json:"max,omitempty"` // Max value for field.
	Min int `json:"min,omitempty"` // Min value for field.
}

type Validator struct {
	fields []FieldError
	err    error
}

func (v *Validator) SetInternalError(err error) {
	v.err = err
}

func (v *Validator) CheckBasic(ctx context.Context, item interface{}) {
	v.SetInternalError(nil)
	err := validate.StructCtx(ctx, item)
	if err == nil {
		return
	}

	var playgroundValidationErrors validator.ValidationErrors

	if errors.As(err, &playgroundValidationErrors) {
		for _, fieldError := range playgroundValidationErrors {
			v.fields = append(v.fields, NewFieldError(fieldError))
		}
	} else {
		v.SetInternalError(err)
	}
}*/

func (em EmService) isValid(ctx context.Context, person *db.Person) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, person); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		var validationError *validator.ValidationErrors
		if errors.As(err, &invalidValidationError) {
			return fmt.Errorf("[isValid] invalid validation: %w", err)
		}
		if errors.As(err, &validationError) {
			//fmt.Printf("[isValid] validation failed: ", err)
			return fmt.Errorf("[isValid] validation failed: %w", err)
		}
		return err
	}
	return nil
}

func (em EmService) isValidUpdate(ctx context.Context, person *db.UpdatePerson) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, person); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		var validationError *validator.ValidationErrors
		if errors.As(err, &invalidValidationError) {
			return fmt.Errorf("[isValid] invalid validation: %w", err)
		}
		if errors.As(err, &validationError) {
			//fmt.Printf("[isValid] validation failed: ", err)
			return fmt.Errorf("[isValid] validation failed: %w", err)
		}
		return err
	}
	return nil
}
