package configtypes

import (
	"encoding/json"
	"strings"

	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

// OptStringList represents an optional parameter that is a slice of string values. A nil slice is
// equivalent to an empty OptStringList{}, and is not the same as a zero-length slice.
//
// In a configuration file, this can be represented either as a single comma-delimited string or as
// multiple parameters with the same name. In environment variables, it is always represented as a
// single comma-delimited string since there can be only one value for each variable name.
//
// When converting from JSON, the value can be null, a single string, or an array of strings. When
// converting to JSON, the value will be null or an array of strings.
//
// String slices are always copied when getting or setting this type, to ensure that configuration
// structs do not expose mutable data.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptStringList struct {
	hasValue bool
	values   []string
}

func NewOptStringList(values []string) OptStringList {
	if values == nil {
		return OptStringList{}
	}
	v := make([]string, len(values))
	copy(v, values)
	return OptStringList{hasValue: true, values: v}
}

func NewOptStringListFromString(s string) OptStringList {
	if s == "" {
		return OptStringList{}
	}
	if strings.Contains(s, ",") {
		return NewOptStringList(strings.Split(s, ","))
	}
	return NewOptStringList([]string{s})
}

func (o OptStringList) IsDefined() bool {
	return o.hasValue
}

func (o OptStringList) Values() []string {
	if !o.hasValue {
		return nil
	}
	v := make([]string, len(o.values))
	copy(v, o.values)
	return v
}

func (o OptStringList) String() string {
	if len(o.values) == 0 {
		return ""
	}
	return strings.Join(o.values, ",")
}

func (o OptStringList) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *OptStringList) UnmarshalText(data []byte) error {
	parsed := NewOptStringListFromString(string(data))
	if o.hasValue {
		o.values = append(o.values, parsed.values...)
	} else {
		*o = parsed
	}
	return nil
}

func (o OptStringList) MarshalJSON() ([]byte, error) {
	if o.hasValue {
		return json.Marshal(o.values)
	}
	return json.Marshal(nil)
}

func (o *OptStringList) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	if err := v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch v.Type() {
	case ldvalue.NullType:
		*o = OptStringList{}
	case ldvalue.StringType:
		*o = NewOptStringList([]string{v.StringValue()})
	case ldvalue.ArrayType:
		values := make([]string, 0, v.Count())
		for i := 0; i < v.Count(); i++ {
			elem := v.GetByIndex(i)
			if !elem.IsString() {
				return errStringListJSONFormat()
			}
			values = append(values, elem.StringValue())
		}
		*o = NewOptStringList(values)
	default:
		return errStringListJSONFormat()
	}
	return nil
}
