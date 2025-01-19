package govalidity

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hoitek/Govalidity/govalidityt"
	"net/http"
	"reflect"
	"strings"

	"github.com/hoitek/Govalidity/govaliditybody"
	"github.com/hoitek/Govalidity/govalidityconv"
	"github.com/hoitek/Govalidity/govaliditym"
)

type (
	Schema                 = map[string]interface{}
	ValidationErrors       = map[string][]error
	ValidityResponseErrors = map[string][]string
	Body                   = map[string]interface{}
	Queries                = map[string]interface{}
	Params                 = map[string]string
)

func isJson(s string) bool {
	var j map[string]interface{}
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		return false
	}
	return true
}

func convertToMap(s string) map[string]interface{} {
	var j map[string]interface{}
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		return nil
	}
	return j
}

func convertToJson(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func sanitizeDataMapToJson(dataMap map[string]interface{}) map[string]interface{} {
	for k, v := range dataMap {
		switch v.(type) {
		case string:
			if isJson(v.(string)) {
				dataMap[k] = convertToMap(v.(string))
				sanitizeDataMapToJson(dataMap[k].(map[string]interface{}))
			}
		}
	}
	return dataMap
}

var validationErrors = ValidationErrors{}
var isValid = true

func validateByJson(baseDataMap map[string]interface{}, dataMap map[string]interface{}, validations Schema, structData interface{}) ValidationErrors {
	for k, v := range validations {
		value, ok := dataMap[k]
		if !ok {
			switch v.(type) {
			case *Validator:
				value = ""
				if v.(*Validator).DefaultValue != "" {
					value = v.(*Validator).DefaultValue
					dataMap[k] = value
				}
			case Schema:
				value = Schema{}
			}
		}
		switch v.(type) {
		case Schema:
			valueStr := fmt.Sprintf("%s", value)
			temp := map[string]interface{}{}
			if isJson(valueStr) {
				jsonData := convertToMap(valueStr)
				jsn, ok := jsonData[k]
				if ok {
					validateByJson(baseDataMap, jsn.(Schema), v.(Schema), &temp)
				}
			} else {
				validateByJson(baseDataMap, value.(Schema), v.(Schema), &temp)
			}
		case *Validator:
			validator := v.(*Validator)
			validator.Value = value
			if (!ok && !validator.IsOptional) || ok {
				errs := validator.run()
				if len(errs) > 0 {
					isValid = false
					validationErrors[validator.Field] = errs
				}
			}
		}
	}
	baseDataMap = sanitizeDataMapToJson(baseDataMap)
	otherTypeDataMap := map[string]interface{}{}
	for k, v := range baseDataMap {
		switch value := v.(type) {
		case *govalidityt.File, []*govalidityt.File:
			otherTypeDataMap[k] = value
			delete(baseDataMap, k)
		}
	}
	if isValid {
		bytes, err := json.Marshal(baseDataMap)
		if err != nil {
			return ValidationErrors{
				"error": {
					errors.New("Input Data is Invalid"),
				},
			}
		}
		err = json.Unmarshal(bytes, &structData)
		if err != nil {
			errMsg := err.Error()
			errMsg = strings.ReplaceAll(errMsg, "json: cannot unmarshal number into Go struct field", "")
			errMsg = strings.ReplaceAll(errMsg, "json: cannot unmarshal string into Go struct field", "")
			errMsg = strings.ReplaceAll(errMsg, "of type string", "")
			errMsg = strings.TrimSpace(errMsg)
			errMsg = "Check data type of these fields: " + errMsg
			return ValidationErrors{
				"error": {
					errors.New(errMsg),
				},
			}
		}
		for k, v := range otherTypeDataMap {
			switch convertedStructData := structData.(type) {
			case map[string]interface{}:
				convertedStructData[k] = v
			default:
				structValue := reflect.ValueOf(structData)

				// Check if structData is a pointer to a struct with reflect
				if structValue.Kind() == reflect.Ptr {
					structValue = structValue.Elem()
					// Check if structData is a pointer to a struct with reflect
					if structValue.Kind() == reflect.Ptr {
						structValue = structValue.Elem()
					}
					// Check if structData is a struct with reflect
					if structValue.Kind() == reflect.Struct {
						// Get Field By Tag Name and Set Value
						structType := structValue.Type()
						for i := 0; i < structType.NumField(); i++ {
							field := structType.Field(i)
							if field.Tag.Get("json") == k {
								structFieldValue := structValue.FieldByName(field.Name)
								if structFieldValue.IsValid() {
									structFieldValue.Set(reflect.ValueOf(v))
								}
							}
						}
					}
				}

				// Check if structData is a struct with reflect
				if structValue.Kind() == reflect.Struct {
					// Get Field By Tag Name and Set Value
					structType := structValue.Type()
					for i := 0; i < structType.NumField(); i++ {
						field := structType.Field(i)
						if field.Tag.Get("json") == k {
							structFieldValue := structValue.FieldByName(field.Name)
							if structFieldValue.IsValid() {
								structFieldValue.Set(reflect.ValueOf(v))
							}
						}
					}
				}
			}
		}
	}
	return validationErrors
}

func ValidateQueries(r *http.Request, validations Schema, structData interface{}) ValidationErrors {
	validationErrors = ValidationErrors{}
	isValid = true
	baseDataMap := Queries{}
	dataMap := Queries{}
	queries := r.URL.Query()

	for k, v := range queries {
		if len(v) > 0 {
			baseDataMap[k] = v[0]
		} else {
			baseDataMap[k] = ""
		}
	}

	for k, v := range baseDataMap {
		value := v.(string)
		if isJson(value) {
			dataMap[k] = `{"` + k + `":` + value + `}`
		} else {
			dataMap[k] = value
		}
	}

	return validateByJson(baseDataMap, dataMap, validations, structData)
}

func isSliceOfString(slice []interface{}) bool {
	isString := false
	for _, sliceItem := range slice {
		_, err := govalidityconv.ToNumber(sliceItem)
		if err != nil {
			isString = true
		}

		if isString {
			break
		}
	}
	return isString
}

func convertToSliceOfString(slice []interface{}) []string {
	result := []string{}
	for _, sliceItem := range slice {
		result = append(result, fmt.Sprintf("%s", sliceItem))
	}
	return result
}

func convertToSliceOfNumber(slice []interface{}) ([]float64, error) {
	result := []float64{}
	for _, sliceItem := range slice {
		num, err := govalidityconv.ToNumber(sliceItem)
		if err != nil || num == nil {
			return []float64{}, errors.New("All of slice items should be number")
		}
		result = append(result, *num)
	}
	return result, nil
}

func getStructFieldType(structData interface{}, fieldName string) reflect.Type {
	structValue := reflect.ValueOf(structData)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	if structValue.Kind() == reflect.Struct {
		structType := structValue.Type()
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			if field.Tag.Get("json") == fieldName {
				return field.Type
			}
		}
	}
	return nil
}

func ValidateBody(r *http.Request, validations Schema, structData interface{}) ValidationErrors {
	validationErrors = ValidationErrors{}
	isValid = true
	dataMap := Body{}
	var baseDataMap Body
	err := govaliditybody.Bind(r, &baseDataMap)
	if err != nil {
		return ValidationErrors{
			"UnknownError": []error{
				err,
			},
		}
	}

	// Remove keys that in struct has tags format binary
	structValue := reflect.ValueOf(structData)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}
	if structValue.Kind() == reflect.Struct {
		structType := structValue.Type()
		for k, _ := range baseDataMap {
			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				if field.Tag.Get("json") == k {
					if strings.Contains(field.Tag.Get("openapi"), "format:binary") {
						delete(baseDataMap, k)
					}
				}
			}
		}
	}

	if r.MultipartForm != nil {
		fileKeys := make([]string, 0, len(r.MultipartForm.File))
		for key := range r.MultipartForm.File {
			fileKeys = append(fileKeys, key)
		}
		for _, v := range fileKeys {
			file, header, err := r.FormFile(v)
			if err == nil {
				if len(r.MultipartForm.File[v]) > 1 {
					// Get multipart.File and *multipart.FileHeader for each files
					files := []*govalidityt.File{}
					for _, fileHeader := range r.MultipartForm.File[v] {
						file, err := fileHeader.Open()
						if err != nil {
							errs := ValidationErrors{}
							errs[v] = []error{err}
							return errs
						}
						fileData := &govalidityt.File{
							File:   &file,
							Header: fileHeader,
						}
						files = append(files, fileData)
					}
					baseDataMap[v] = files
				} else {
					fieldType := getStructFieldType(structData, v)
					if fieldType != nil && fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Ptr && fieldType.Elem().Elem().Name() == "File" {
						baseDataMap[v] = []*govalidityt.File{
							&govalidityt.File{
								File:   &file,
								Header: header,
							},
						}
					} else {
						baseDataMap[v] = &govalidityt.File{
							File:   &file,
							Header: header,
						}
					}
				}
			}
		}
	}

	for k, v := range baseDataMap {
		switch v.(type) {
		case string:
			if isJson(v.(string)) {
				dataMap[k] = `{"` + k + `":` + v.(string) + `}`
			} else {
				dataMap[k] = v.(string)
			}
		case float64:
			dataMap[k] = v.(float64)
		case []interface{}:
			slice := v.([]interface{})
			isString := isSliceOfString(slice)
			sanitizedSlice := []interface{}{}
			if isString {
				cSlice := convertToSliceOfString(slice)
				for _, cSliceItem := range cSlice {
					sanitizedSlice = append(sanitizedSlice, cSliceItem)
				}
			} else {
				cSlice, err := convertToSliceOfNumber(slice)
				if err != nil {
					sanitizedSlice = []interface{}{}
				}
				for _, cSliceItem := range cSlice {
					sanitizedSlice = append(sanitizedSlice, cSliceItem)
				}
			}
			jsonData := convertToJson(sanitizedSlice)
			dataMap[k] = jsonData
		case *govalidityt.File:
			dataMap[k] = v.(*govalidityt.File)
		case []*govalidityt.File:
			dataMap[k] = v.([]*govalidityt.File)
		}
	}

	return validateByJson(baseDataMap, dataMap, validations, structData)
}

func ValidateParams(params Params, validations Schema, structData interface{}) ValidationErrors {
	validationErrors = ValidationErrors{}
	isValid = true
	baseDataMap := Body{}
	dataMap := Body{}

	for k, v := range params {
		baseDataMap[k] = v
	}

	for k, v := range baseDataMap {
		value := v.(string)
		if isJson(value) {
			dataMap[k] = `{"` + k + `":` + value + `}`
		} else {
			dataMap[k] = value
		}
	}

	return validateByJson(baseDataMap, dataMap, validations, structData)
}

func SetDefaultErrorMessages(v *govaliditym.Validations) {
	govaliditym.SetMessages(v)
}

func SetFieldLabels(l *govaliditym.Labels) {
	govaliditym.SetFieldLables(l)
}

func (v *Validator) run() []error {
	var errs []error
	for _, validation := range v.Validations {
		switch v.Value.(type) {
		case *govalidityt.File:
			isValid, err := validation.Fn(v.Field, v.Value)
			if !isValid {
				errs = append(errs, err)
			}
		case []*govalidityt.File:
			isValid, err := validation.Fn(v.Field, v.Value)
			if !isValid {
				errs = append(errs, err)
			}
		default:
			number, errConv := govalidityconv.ToNumber(v.Value)
			str := ""
			if errConv == nil && number != nil {
				str = fmt.Sprintf("%v", *number)
			} else {
				str = v.Value.(string)
			}
			isValid, err := validation.Fn(v.Field, str)
			if !isValid {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func DumpErrors(errs ValidationErrors) ValidityResponseErrors {
	errMap := map[string][]string{}
	for k, vErrs := range errs {
		for _, e := range vErrs {
			errMap[k] = append(errMap[k], e.Error())
		}
	}
	return errMap
}
