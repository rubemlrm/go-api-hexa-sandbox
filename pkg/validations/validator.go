package validations

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vl "github.com/go-playground/validator/v10"
)

type Validater interface {
	Check(val interface{}) bool
	CheckWithTranslations(val interface{}) error
}

type Validator struct {
	validate     *vl.Validate
	translations ut.Translator
}

type validatorError struct {
	error
}

func New(locale string, options ...func(*Validator, *validatorError)) (*Validator, error) {
	lt := en.New()
	uni := ut.New(lt, lt)

	validate := vl.New()
	trans, _ := uni.GetTranslator(locale)

	validator := &Validator{
		validate,
		trans,
	}

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	for _, o := range options {
		ev := &validatorError{}
		o(validator, ev)
		if ev.error != nil {
			return nil, ev.error
		}
	}

	return validator, nil
}

func WithCustomTranslation(tag string, message string) func(*Validator, *validatorError) {
	return func(v *Validator, ev *validatorError) {
		if tag == "" {
			ev.error = fmt.Errorf("tag name can't be empty")
			return
		}

		if message == "" {
			ev.error = fmt.Errorf("message can't be empty")
			return
		}
		ev.error = v.validate.RegisterTranslation(tag, v.translations, func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		}, func(ut ut.Translator, fe vl.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
	}
}

func WithCustomValidationRule(tag string, fn vl.Func) func(*Validator, *validatorError) {
	return func(v *Validator, ev *validatorError) {
		if tag == "" {
			ev.error = fmt.Errorf("tag name can't be empty")
			return
		}

		if fn == nil {
			ev.error = fmt.Errorf("function can't be nil")
			return
		}
		ev.error = v.validate.RegisterValidation(tag, fn)
	}
}

func (v *Validator) Check(val interface{}) error {
	err := v.validate.Struct(val)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to validate: %s", err))
	}
	return nil
}

func (v *Validator) CheckWithTranslations(val interface{}) error {
	err := v.validate.Struct(val)
	if err != nil {
		var trs vl.ValidationErrors
		errors.As(err, &trs)
		return errors.New(fmt.Sprintf("failed to validate: %s", trs.Translate(v.translations)))
	}
	return nil
}
