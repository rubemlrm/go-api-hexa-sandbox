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

type Validater[T any] interface {
	Check(vf ValidationFunc) ([]map[string]string, error)
}

type StructValidator interface {
	Struct(val interface{}) error
	RegisterTagNameFunc(fn vl.TagNameFunc)
	RegisterTranslation(tag string, trans ut.Translator, registerFn vl.RegisterTranslationsFunc, translationFn vl.TranslationFunc) error
	RegisterValidation(tag string, fn vl.Func, callValidationEvenIfNull ...bool) error
}

type Validator struct {
	Validate     StructValidator
	Translations ut.Translator
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
		v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
		err := v.Validate.RegisterTranslation(tag, v.Translations, func(ut ut.Translator) error {
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
		return v.Validate.RegisterValidation(tag, fn)
	}
}

func (v Validator) ValidateInput(val interface{}) ([]map[string]string, error) {
	err := v.Validate.Struct(val)
	if err != nil {
		var validationErrors vl.ValidationErrors
		hasErrors := errors.As(err, &validationErrors)
		if hasErrors {
			return v.ConvertToMap(err.(vl.ValidationErrors)), nil
		}
		return nil, fmt.Errorf("failed to validate: %s", err)
	}
	return nil, nil
}

func (v Validator) ConvertToMap(errs vl.ValidationErrors) []map[string]string {
	result := make([]map[string]string, 0, len(errs))
	for _, e := range errs {
		result = append(result, map[string]string{
			"field": e.Field(),
			"error": e.Translate(v.Translations),
		})
	}
	return result
}
