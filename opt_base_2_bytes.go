package configtypes

import (
	"encoding/json"

	"github.com/alecthomas/units"

	"github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

// OptBase2Bytes represents an optional parameter which, if present, must be a
// valid units.Base2Bytes.
//
// There is no additional validation. The value is stored as a
// *units.Base2Bytes pointer, whose value is always copied when accessed; a nil
// pointer indicates the lack of a value (the same as OptBase2Bytes{}).
//
// Converting to or from a string uses the standard behavior for
// units.Base2Bytes.String() and units.ParseBase2Bytes().
//
// When converting to or from JSON, the value must be either a JSON null or a
// JSON string.
//
// See the package documentation for the general contract for methods that have
// no specific documentation here.
type OptBase2Bytes struct {
	hasValue bool
	size     units.Base2Bytes
}

func NewOptBase2Bytes(size units.Base2Bytes) OptBase2Bytes {
	return OptBase2Bytes{hasValue: true, size: size}
}

func NewOptBase2BytesFromString(sizeAsString string) (OptBase2Bytes, error) {
	if sizeAsString == "" {
		return OptBase2Bytes{}, nil
	}
	size, err := units.ParseBase2Bytes(sizeAsString)
	if err == nil {
		return NewOptBase2Bytes(size), nil
	}
	return OptBase2Bytes{}, errBase2BytesFormat()
}

func (o OptBase2Bytes) IsDefined() bool {
	return o.hasValue
}

func (o OptBase2Bytes) GetOrElse(orElseValue units.Base2Bytes) units.Base2Bytes {
	if o.IsDefined() {
		return o.size
	}

	return orElseValue
}

// Get returns the value if it is defined.
//
// The result of this method is only valid if IsDefined() returns true.
func (o OptBase2Bytes) Get() units.Base2Bytes {
	return o.size
}

func (o OptBase2Bytes) String() string {
	if o.IsDefined() {
		return o.size.String()
	}
	return ""
}

func (o OptBase2Bytes) MarshalText() ([]byte, error) {
	if o.IsDefined() {
		return []byte(o.String()), nil
	}
	return nil, nil
}

func (o *OptBase2Bytes) UnmarshalText(data []byte) error {
	parsed, err := NewOptBase2BytesFromString(string(data))
	if err == nil {
		*o = parsed
	}
	return err
}

func (o OptBase2Bytes) MarshalJSON() ([]byte, error) {
	if o.IsDefined() {
		return json.Marshal(o.size.String())
	}
	return json.Marshal(nil)
}

func (o *OptBase2Bytes) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	var err error
	if err = v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptBase2Bytes{}
		return nil
	case v.IsString():
		*o, err = NewOptBase2BytesFromString(v.StringValue())
		return err
	default:
		return errBase2BytesFormat()
	}
}
