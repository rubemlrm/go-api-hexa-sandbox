package validations

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Validater interface {
	RegisterCustomTranslation(tag string, message string) error
	RegisterCustomValidationRule(tag string, fn validator.Func) error

	Check(val interface{}) bool
	CheckWithTranslations(val interface{}) error
}

type Validator struct {
	validate     *validator.Validate
	translations ut.Translator
}

func NewValidator(locale string) *Validator {
	lt := en.New()
	uni := ut.New(lt, lt)

	validate := validator.New(validator.WithRequiredStructEnabled())
	trans, _ := uni.GetTranslator(locale)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validate,
		trans,
	}
}

func (v *Validator) RegisterCustomTranslation(tag string, message string) error {
	err := v.validate.RegisterTranslation(tag, v.translations, func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})

	if err != nil {
		return err
	}
	return nil
}

func (v *Validator) RegisterCustomValidationRule(tag string, fn validator.Func) error {
	if tag == "" {
		return fmt.Errorf("tag name can't be empty")
	}

	if fn == nil {
		return fmt.Errorf("function can't be nil")
	}

	err := v.validate.RegisterValidation(tag, fn)
	if err != nil {
		return err
	}
	return nil
}

func (v *Validator) Check(val interface{}) error {
	err := v.validate.Struct(val)
	if err != nil {
		return nil
	}
	return nil
}

func (v *Validator) CheckWithTranslations(val interface{}) error {
	err := v.validate.Struct(val)
	if err != nil {
		trs := err.(validator.ValidationErrors)
		return errors.New(fmt.Sprintf("failed to validate: %s", trs.Translate(v.translations)))
	}
	return nil
}
