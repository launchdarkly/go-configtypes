package configtypes

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

// OptString represents an optional string parameter.
//
// In Go, an unset string normally defaults to "". This type allows application code to distinguish
// between "" and the absence of a value, in case some other default behavior is desired.
//
// When converting to a string, an undefined value becomes "". When converting from a string, any string
// becomes a defined value even if it is "". If you want to treat "" as an undefined value, use
// OptStringNonEmpty.
//
// When converting to or from JSON, the value must be either a JSON null or a JSON string.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptString struct {
	s ldvalue.OptionalString
}

func NewOptString(value string) OptString {
	return OptString{ldvalue.NewOptionalString(value)}
}

func (o OptString) IsDefined() bool {
	return o.s.IsDefined()
}

func (o OptString) GetOrElse(orElseValue string) string {
	return o.s.OrElse(orElseValue)
}

func (o OptString) String() string {
	return o.s.StringValue()
}

func (o OptString) MarshalText() ([]byte, error) {
	return []byte(o.s.StringValue()), nil
}

func (o *OptString) UnmarshalText(data []byte) error {
	*o = NewOptString(string(data))
	return nil // cannot fail
}

func (o OptString) MarshalJSON() ([]byte, error) {
	return o.s.MarshalJSON()
}

func (o *OptString) UnmarshalJSON(data []byte) error {
	var s ldvalue.OptionalString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	*o = OptString{s}
	return nil
}
