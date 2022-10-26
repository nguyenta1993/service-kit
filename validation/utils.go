package validation

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gogovan-korea/ggx-kr-service-utils/localization"
	"github.com/gogovan-korea/ggx-kr-service-utils/logger"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return validate.Struct(val)
}

func ValidateVariable(s interface{}, rules string) error {
	var value string

	val := reflect.ValueOf(s)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = strconv.FormatInt(val.Int(), 10)
	case reflect.String:
		value = val.String()
	}

	return validate.Var(value, rules)
}

func GetErrorMessage(err validator.FieldError, logger logger.Logger) string {
	logger.Error("Validation Error", zap.Error(err))

	templateData := make(map[string]string)

	if err.Field() != "" {
		templateData["FieldName"] = localization.Localize(strings.ToLower(err.Field()), nil)
	}

	if err.Param() != "" {
		templateData["FieldData"] = localization.Localize(strings.ToLower(err.Param()), nil)
	}

	return localization.Localize(err.ActualTag(), templateData)
}
