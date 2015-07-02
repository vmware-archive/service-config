package test_helpers

import (
	. "github.com/onsi/gomega"

	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate() error
}

func setNestedFieldToEmpty(obj interface{}, nestedFieldNames []string) error {

	s := reflect.ValueOf(obj).Elem()
	if s.Type().Kind() == reflect.Slice {
		if s.Len() == 0 {
			return errors.New("Trying to set nested property on empty slice")
		}
		s = s.Index(0)
	}

	currFieldName := nestedFieldNames[0]
	remainingFieldNames := nestedFieldNames[1:]
	field := s.FieldByName(currFieldName)
	if field.IsValid() == false {
		return errors.New(fmt.Sprintf("Field '%s' is not defined", currFieldName))
	}

	if len(remainingFieldNames) == 0 {
		fieldType := field.Type()
		field.Set(reflect.Zero(fieldType))
		return nil
	}
	return setNestedFieldToEmpty(field.Addr().Interface(), remainingFieldNames)
}

func setFieldToEmpty(obj interface{}, fieldName string) error {
	return setNestedFieldToEmpty(obj, strings.Split(fieldName, "."))
}

func IsRequiredField(obj Validator, fieldName string) func() {
	return func() {
		err := setFieldToEmpty(obj, fieldName)
		Expect(err).NotTo(HaveOccurred())

		err = obj.Validate()

		Expect(err).To(HaveOccurred())

		fieldParts := strings.Split(fieldName, ".")
		for _, fieldPart := range fieldParts {
			Expect(err.Error()).To(ContainSubstring(fieldPart))
		}
	}
}

func IsOptionalField(obj Validator, fieldName string) func() {
	return func() {
		err := setFieldToEmpty(obj, fieldName)
		Expect(err).NotTo(HaveOccurred())

		err = obj.Validate()

		Expect(err).NotTo(HaveOccurred())
	}
}
