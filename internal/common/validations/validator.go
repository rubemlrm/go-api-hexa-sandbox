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

type ValidationFunc func(string, ...Option) (*Validator, error)

type Option func(v *Validator) error

func New(locale string, options ...Option) (*Validator, error) {
	lt := en.New()
	uni := ut.New(lt, lt)

	validate := vl.New()
	trans, _ := uni.GetTranslator(locale)

	validator := &Validator{
		validate,
		trans,
	}

	for _, o := range options {
		err := o(validator)
		if err != nil {
			return nil, err
		}
	}

	return validator, nil
}

func WithCustomFieldLabel(label string) Option {
	return func(v *Validator) error {
		if label == "" {
			return fmt.Errorf("custom field label is required")
		}
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get(label), ",", 1)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		return nil
	}
}
func WithCustomTranslation(tag string, message string) Option {
	return func(v *Validator) error {
		if tag == "" {
			return fmt.Errorf("tag name can't be empty")
		}

		if message == "" {
			return fmt.Errorf("message can't be empty")
		}
		err := v.validate.RegisterTranslation(tag, v.translations, func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		}, func(ut ut.Translator, fe vl.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
		return err
	}
}

func WithCustomValidationRule(tag string, fn vl.Func) Option {
	return func(v *Validator) error {
		if tag == "" {
			return fmt.Errorf("tag name can't be empty")
		}

		if fn == nil {
			return fmt.Errorf("function can't be nil")
		}
		return v.validate.RegisterValidation(tag, fn)
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
