package configtypes

import (
	"encoding/json"
	"strconv"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptFloat64 represents an optional float64 parameter.
//
// In Go, an unset float64 normally defaults to zero. This type allows application code to distinguish
// between zero and the absence of a value, in case some other default behavior is desired.
//
// When converting to or from JSON, the value must be either a JSON null or a JSON number.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptFloat64 struct {
	hasValue bool
	value    float64
}

func NewOptFloat64(value float64) OptFloat64 {
	return OptFloat64{hasValue: true, value: value}
}

func NewOptFloat64FromString(s string) (OptFloat64, error) {
	if s == "" {
		return OptFloat64{}, nil
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return OptFloat64{}, errFloatFormat()
	}
	return NewOptFloat64(n), nil
}

func (o OptFloat64) IsDefined() bool {
	return o.hasValue
}

func (o OptFloat64) GetOrElse(orElseValue float64) float64 {
	if o.hasValue {
		return o.value
	}
	return orElseValue
}

func (o OptFloat64) String() string {
	if o.hasValue {
		return strconv.FormatFloat(o.value, 'f', -1, 64)
	}
	return ""
}

func (o OptFloat64) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *OptFloat64) UnmarshalText(data []byte) error {
	value, err := NewOptFloat64FromString(string(data))
	if err == nil {
		*o = value
	}
	return err
}

func (o OptFloat64) MarshalJSON() ([]byte, error) {
	if o.hasValue {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *OptFloat64) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	if err := v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptFloat64{}
	case v.IsNumber():
		*o = NewOptFloat64(v.Float64Value())
	default:
		return errFloatFormat()
	}
	return nil
}
