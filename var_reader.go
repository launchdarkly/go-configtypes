package configtypes

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// VarReader reads string values from named variables, such as environment variables, and
// translates them into values of any supported type. It accumulates errors as it goes.
//
// The supported types are any type that implements TextUnmarshaler (which includes all of the Opt
// and Req types defined in this package), and also the primitive types bool, int, float64, and
// string.
//
// You may specify the variable name for each target value programmatically, or use struct field
// tags as described in ReadStruct(), or both.
type VarReader struct {
	values map[string]string
	result *ValidationResult
	prefix string
	suffix string
}

// NewVarReaderFromEnvironment creates a VarReader that reads from environment variables.
func NewVarReaderFromEnvironment() *VarReader {
	r := &VarReader{result: new(ValidationResult)}
	vars := os.Environ()
	r.values = make(map[string]string, len(vars))
	for _, s := range vars {
		k, v := parseVar(s)
		r.values[k] = v
	}
	return r
}

// NewVarReaderFromValues creates a VarReader that reads from the specified name-value map.
func NewVarReaderFromValues(values map[string]string) *VarReader {
	r := &VarReader{result: new(ValidationResult)}
	r.values = make(map[string]string, len(values))
	for k, v := range values {
		r.values[k] = v
	}
	return r
}

// Result returns a ValidationResult containing all of the errors encountered so far.
func (r VarReader) Result() ValidationResult {
	return *r.result
}

// Read attempts to read an environment variable into a target value.
//
// The varName may be modified by any previous calls to WithVarNamePrefix or WithVarNameSuffix.
//
// If the variable exists, Read attempts to set the target value as follows: through the
// TextMarshaler interface, if that is implemented; otherwise, if it is a supported built-in type,
// it uses the corresponding Opt type (such as OptBool) to parse the value. If the type is not
// supported, it records an error.
//
// If the variable does not exist, the method does nothing; it does not modify the target value.
// Use ReadRequired or ReadStruct, or call Validate afterward, if you want missing values to be
// treated as errors.
//
// The method returns true if the variable was found (regardless of whether unmarshaling succeeded)
// or false if it was not found.
func (r *VarReader) Read(varName string, target interface{}) bool {
	return r.readInternal(varName, target, false)
}

// ReadRequired is the same as Read, except that if the variable was not found, it records an
// error for that variable name.
func (r *VarReader) ReadRequired(varName string, target interface{}) bool {
	return r.readInternal(varName, target, true)
}

func (r *VarReader) readInternal(varName string, target interface{}, required bool) bool {
	setter := setterForTarget(target)
	if setter == nil {
		r.AddError(ValidationPath{varName}, varReaderBadTargetTypeError(target))
		return false
	}
	s, ok := r.get(varName)
	if !ok {
		if required {
			r.AddError(ValidationPath{varName}, errRequired())
		}
		return false
	}
	err := setter([]byte(s))
	if err != nil {
		r.AddError(ValidationPath{varName}, err)
	}
	return true
}

// ReadStruct uses reflection to populate any exported fields of the target struct that have a tag
// of `conf:"VAR_NAME"`. The behavior for each of these fields is the same as for Read(), unless you
// specify `conf:"VAR_NAME,required"` in which case it behaves like ReadRequired(). If the recursive
// parameter is true, then ReadStruct will be called recursively on any embedded structs.
//
//	   type myStruct struct {
//         MyOptBool               OptBool `conf:"VAR1"`
//         MyPrimitiveBool         bool    `conf:"VAR2"`
//         MyRequiredBool          bool    `conf:"VAR3,required"`
//     }
//     s := myStruct{}
//     r := NewVarReaderFromEnvironment()
//     r.ReadStruct(&myStruct)
//
// In the above example, each field behaves slightly differently. MyOptBool can have three states:
// undefined (VAR1 was not set), true, or false. MyPrimitiveBool is a simple bool so there is no way
// to distinguish between its default value and "not set". MyRequiredBool is a simple bool but will
// cause VarReader to log an error if the variable is not set.
func (r *VarReader) ReadStruct(target interface{}, recursive bool) {
	ok := r.readStructFields(target, recursive)
	if !ok {
		r.AddError(nil, errors.New("ReadStruct was called on something other than a struct pointer"))
	}
}

func (r *VarReader) readStructFields(target interface{}, recursive bool) bool {
	refStruct, ok := getReflectValueForStructPtr(target)
	if !ok {
		return false
	}

	structType := refStruct.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldInType := structType.Field(i)
		if !isFieldExported(fieldInType) {
			continue
		}
		tagInfo, err := getFieldTagInfo(fieldInType)
		if err != nil {
			r.AddError(ValidationPath{fieldInType.Name}, err)
			continue
		}
		fieldInInstancePtr := refStruct.FieldByName(fieldInType.Name)
		if fieldInInstancePtr.Kind() != reflect.Ptr {
			fieldInInstancePtr = fieldInInstancePtr.Addr()
		}
		fieldValuePtr := fieldInInstancePtr.Interface()
		switch {
		case tagInfo.varName == "":
			if recursive {
				r.readStructFields(fieldValuePtr, true) // harmless if this isn't a struct
			}
		case tagInfo.required:
			r.ReadRequired(tagInfo.varName, fieldValuePtr)
		default:
			r.Read(tagInfo.varName, fieldValuePtr)
		}
	}

	return true
}

// WithVarNamePrefix returns a new VarReader based on the current one, which accumulates errors
// in the same ValidationResult, but with the given prefix added to all variable names.
//
//     r0 := NewVarReaderFromValues(map[string]string{"b_x": "2", "b_y": "3"})
//     r1 := r.WithVarNamePrefix("b_")
//     r1.Read(&x, "x")  // x is set to "2"
func (r *VarReader) WithVarNamePrefix(prefix string) *VarReader {
	return &VarReader{
		values: r.values,
		result: r.result,
		prefix: prefix + r.prefix,
		suffix: r.suffix,
	}
}

// WithVarNameSuffix returns a new VarReader based on the current one, which accumulates errors
// in the same ValidationResult, but with the given suffix added to all variable names.
//
//     r0 := NewVarReaderFromValues(map[string]string{"a_x": "2", "b_x": "3"})
//     r1 := r.WithVarNameSuffix("_x")
//     r1.Read(&b, "b")  // b is set to "3"
func (r *VarReader) WithVarNameSuffix(suffix string) *VarReader {
	return &VarReader{
		values: r.values,
		result: r.result,
		prefix: r.prefix,
		suffix: r.suffix + suffix,
	}
}

// FindPrefixedValues finds all named values in the VarReader that have the specified name prefix,
// and returns a name-value map of only those, with the prefixes removed.
//
//     r := NewVarReaderFromValues(map[string]string{"a": "1", "b_x": "2", "b_y": "3"})
//     values := r.FindPrefixedValues("b_")
//     // values == { "x": "2", "y": "3" }
func (r *VarReader) FindPrefixedValues(prefix string) map[string]string {
	ret := make(map[string]string)
	for n, v := range r.values {
		if strings.HasPrefix(n, prefix) {
			ret[strings.TrimPrefix(n, prefix)] = v
		}
	}
	return ret
}

// AddError records an error in the VarReader's result.
func (r *VarReader) AddError(path ValidationPath, e error) {
	r.result.AddError(r.transformPath(path), e)
}

func (r VarReader) get(varName string) (string, bool) {
	value, found := r.values[r.prefix+varName+r.suffix]
	return value, found
}

func (r VarReader) transformPath(path ValidationPath) ValidationPath {
	if r.prefix == "" && r.suffix == "" {
		return path
	}
	ret := make(ValidationPath, len(path))
	copy(ret, path)
	ret[len(ret)-1] = r.prefix + ret[len(ret)-1] + r.suffix
	return ret
}

func parseVar(s string) (string, string) {
	p := strings.Index(s, "=")
	return s[:p], s[p+1:]
}

func setterForTarget(target interface{}) func(data []byte) error {
	if tu, ok := target.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText
	}
	switch p := target.(type) {
	case *bool:
		return func(data []byte) error {
			var v OptBool
			if err := v.UnmarshalText(data); err != nil {
				return err
			}
			if v.IsDefined() {
				*p = v.GetOrElse(false)
			}
			return nil
		}
	case *int:
		return func(data []byte) error {
			var v OptInt
			if err := v.UnmarshalText(data); err != nil {
				return err
			}
			if v.IsDefined() {
				*p = v.GetOrElse(0)
			}
			return nil
		}
	case *float64:
		return func(data []byte) error {
			var v OptFloat64
			if err := v.UnmarshalText(data); err != nil {
				return err
			}
			if v.IsDefined() {
				*p = v.GetOrElse(0)
			}
			return nil
		}
	case *string:
		return func(data []byte) error {
			*p = string(data)
			return nil
		}
	}
	return nil
}

func varReaderBadTargetTypeError(target interface{}) error {
	return fmt.Errorf("could not read into value of type %T", target)
}
