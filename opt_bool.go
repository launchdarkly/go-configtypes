package configtypes

import (
	"encoding/json"
	"strings"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptBool represents an optional boolean parameter.
//
// In Go, an unset boolean normally defaults to false. This type allows application code to distinguish
// between false and the absence of a value, in case some other default behavior is desired.
//
// When setting this value from a string representation, the following case-insensitive string values
// are allowed for true/false: "true"/"false", "0"/"1", "yes"/"no". An empty string value is converted
// to an empty OptBool{}.
//
// Converting to a string always produces "true", "false", or "".
//
// When converting to or from JSON, the value must be either a JSON null or a JSON boolean.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptBool struct {
	v ldvalue.OptionalBool
}

func NewOptBool(value bool) OptBool {
	return OptBool{v: ldvalue.NewOptionalBool(value)}
}

func NewOptBoolFromString(s string) (OptBool, error) {
	if s == "" {
		return OptBool{}, nil
	}
	if s == "1" || strings.EqualFold(s, "true") || strings.EqualFold(s, "yes") {
		return NewOptBool(true), nil
	}
	if s == "0" || s == "" || strings.EqualFold(s, "false") || strings.EqualFold(s, "no") {
		return NewOptBool(false), nil
	}
	return OptBool{}, errBoolFormat()
}

func (o OptBool) IsDefined() bool {
	return o.v.IsDefined()
}

func (o OptBool) GetOrElse(orElseValue bool) bool {
	return o.v.OrElse(orElseValue)
}

func (o OptBool) String() string {
	b, _ := o.MarshalText()
	return string(b)
}

func (o OptBool) MarshalText() ([]byte, error) {
	if o.IsDefined() {
		if o.v.BoolValue() {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	}
	return nil, nil
}

func (o *OptBool) UnmarshalText(data []byte) error {
	parsed, err := NewOptBoolFromString(string(data))
	if err == nil {
		*o = parsed
	}
	return err
}

func (o OptBool) MarshalJSON() ([]byte, error) {
	if o.IsDefined() {
		return json.Marshal(o.v.BoolValue())
	}
	return json.Marshal(nil)
}

func (o *OptBool) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	if err := v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptBool{}
	case v.IsBool():
		*o = NewOptBool(v.BoolValue())
	default:
		if v.IsString() {
			return errBoolFormat()
		}
		return errBoolFormat()
	}
	return nil
}
