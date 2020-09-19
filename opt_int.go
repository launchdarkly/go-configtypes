package configtypes

import (
	"encoding/json"
	"strconv"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptInt represents an optional int parameter.
//
// In Go, an unset int normally defaults to zero. This type allows application code to distinguish
// between zero and the absence of a value, in case some other default behavior is desired.
//
// When converting to or from JSON, the value must be either a JSON null or a JSON number that is an integer.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptInt struct {
	v ldvalue.OptionalInt
}

func NewOptInt(value int) OptInt {
	return OptInt{v: ldvalue.NewOptionalInt(value)}
}

func NewOptIntFromString(s string) (OptInt, error) {
	if s == "" {
		return OptInt{}, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return OptInt{}, errIntFormat()
	}
	return NewOptInt(n), nil
}

func (o OptInt) IsDefined() bool {
	return o.v.IsDefined()
}

func (o OptInt) GetOrElse(orElseValue int) int {
	return o.v.OrElse(orElseValue)
}

func (o OptInt) String() string {
	if o.IsDefined() {
		return strconv.Itoa(o.v.IntValue())
	}
	return ""
}

func (o OptInt) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *OptInt) UnmarshalText(data []byte) error {
	value, err := NewOptIntFromString(string(data))
	if err == nil {
		*o = value
	}
	return err
}

func (o OptInt) MarshalJSON() ([]byte, error) {
	if o.IsDefined() {
		return json.Marshal(o.v.IntValue())
	}
	return json.Marshal(nil)
}

func (o *OptInt) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	if err := v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptInt{}
	case v.IsInt():
		*o = NewOptInt(v.IntValue())
	default:
		return errIntFormat()
	}
	return nil
}
