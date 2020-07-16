package configtypes

import (
	"reflect"
)

// ValidateStruct checks whether all of a struct's exported fields are valid according to the
// tag-based required field rule. If the recursive parameter is true, then ValidateStruct will be
// called recursively on any embedded structs in exported fields.
//
// The required field rule is that if any field has a "conf:" field tag that includes ",required",
// it must have a value that is not the zero value for that type. Therefore, any required field
// that uses an Opt type must be in the "defined" state (since its zero value is the "empty"
// state); a required int field must be non-zero; a required string field must not be ""; etc.
//
// The returned ValidationResult can contain any number of errors.
//
// Calling ValidateStruct with a parameter that is not a struct or struct pointer returns an error
// result.
func ValidateStruct(value interface{}, recursive bool) ValidationResult {
	refStruct, ok := getReflectValueForStruct(value)
	if ok {
		return validateFields(refStruct, recursive)
	} else {
		var result ValidationResult
		result.AddError(nil, errValidateNonStruct())
		return result
	}
}

func validateFields(refStruct reflect.Value, recursive bool) ValidationResult {
	var result ValidationResult
	structType := refStruct.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldInType := structType.Field(i)
		if !isFieldExported(fieldInType) {
			continue
		}
		fieldInInstance := refStruct.FieldByName(fieldInType.Name)
		refFieldStruct, fieldIsStruct := getReflectValueForStruct(fieldInInstance.Interface())
		if fieldIsStruct {
			if recursive {
				result.AddAll(ValidationPath{fieldInType.Name}, validateFields(refFieldStruct, true))
			}
		} else {
			tagInfo, err := getFieldTagInfo(fieldInType)
			if err == nil {
				if tagInfo.required && fieldInInstance.IsZero() {
					result.AddError(ValidationPath{fieldInType.Name}, errRequired())
				}
			} else { // invalid field tag, log an error for it
				result.AddError(ValidationPath{fieldInType.Name}, err)
			}
		}
	}

	return result
}
