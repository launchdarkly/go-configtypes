package configtypes

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptStringNonEmpty represents an optional string parameter which, if defined, must be non-empty.
//
// This is the same as OptString, but with additional validation for the constructor and unmarshalers. It
// is impossible (except with reflection) for code outside this package to construct an instance of this
// type with a defined value that is an empty string.
//
// Since an empty string is already the standard way to indicate an empty value for Opt types in this
// package, NewOptStringNonEmpty("") and UnmarshalText([]byte("")) cannot fail; they simply produce an
// empty value. However, JSON conversion is different: an empty value is a JSON null, not a JSON "".
//
// The semantics of this type are mostly the same as a regular Go string, since there is no need to
// distinguish between an empty string value and the absence of a value. It only exists for feature
// parity, so you can use helper methods like StringOrElse. ReqStringNonEmpty may be more useful.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptStringNonEmpty struct {
	opt OptString
}

func NewOptStringNonEmpty(value string) OptStringNonEmpty {
	if value == "" {
		return OptStringNonEmpty{}
	}
	return OptStringNonEmpty{NewOptString(value)}
}

func (o OptStringNonEmpty) IsDefined() bool {
	return o.opt.IsDefined()
}

func (o OptStringNonEmpty) GetOrElse(orElseValue string) string {
	return o.opt.GetOrElse(orElseValue)
}

func (o OptStringNonEmpty) String() string {
	return o.opt.String()
}

func (o *OptStringNonEmpty) UnmarshalText(data []byte) error {
	*o = NewOptStringNonEmpty(string(data))
	return nil
}

func (o OptStringNonEmpty) MarshalText() ([]byte, error) {
	return []byte(o.opt.GetOrElse("")), nil
}

func (o *OptStringNonEmpty) UnmarshalJSON(data []byte) error {
	var s ldvalue.OptionalString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	if s.IsDefined() {
		if s.StringValue() == "" {
			return errMustBeNonEmptyString()
		}
		*o = NewOptStringNonEmpty(s.StringValue())
	} else {
		*o = OptStringNonEmpty{}
	}
	return nil
}

func (o OptStringNonEmpty) MarshalJSON() ([]byte, error) {
	return o.opt.MarshalJSON()
}
