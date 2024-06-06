package configtypes

import (
	"encoding/json"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

// OptDuration represents an optional time.Duration parameter. Any time.Duration value is allowed; if you
// want to allow only positive values (which is usually desirable), use OptDurationNonNegative.
//
// When setting this value from a string representation, it uses time.ParseDuration, so the allowed formats
// include "9ms" (milliseconds), "9s" (seconds), "9m" (minutes), or combinations such as "1m30s". Converting
// to a string uses similar rules, as implemented by Duration.String().
//
// When converting to or from JSON, an empty value is null, and all other values are strings.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptDuration struct {
	hasValue bool
	value    time.Duration
}

func NewOptDuration(value time.Duration) OptDuration {
	return OptDuration{hasValue: true, value: value}
}

func NewOptDurationFromString(s string) (OptDuration, error) {
	if s == "" {
		return OptDuration{}, nil
	}
	value, err := time.ParseDuration(s)
	if err == nil {
		return NewOptDuration(value), nil
	}
	return OptDuration{}, errDurationFormat()
}

func (o OptDuration) IsDefined() bool {
	return o.hasValue
}

func (o OptDuration) GetOrElse(orElseValue time.Duration) time.Duration {
	if !o.hasValue {
		return orElseValue
	}
	return o.value
}

func (o OptDuration) String() string {
	if !o.hasValue {
		return ""
	}
	return o.value.String()
}

func (o OptDuration) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *OptDuration) UnmarshalText(data []byte) error {
	opt, err := NewOptDurationFromString(string(data))
	if err == nil {
		*o = opt
	}
	return err
}

func (o OptDuration) MarshalJSON() ([]byte, error) {
	if o.hasValue {
		return json.Marshal(o.String())
	}
	return json.Marshal(nil)
}

func (o *OptDuration) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	var err error
	if err = v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptDuration{}
		return nil
	case v.IsString():
		*o, err = NewOptDurationFromString(v.StringValue())
		return err
	default:
		return errDurationFormat()
	}
}
