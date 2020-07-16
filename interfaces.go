package configtypes

import (
	"encoding"
	"encoding/json"
	"fmt"
)

// SingleValue is the common interface for all types in this package that represent a single optional
// or required value.
//
// All such types can be in a defined or an empty state, and can be converted to text or JSON.
type SingleValue interface {
	fmt.Stringer
	encoding.TextMarshaler
	json.Marshaler

	// IsDefined returns true if the value is defined, or false if it is empty.
	IsDefined() bool
}

// Validation is an optional interface for any type that you wish to have custom validation behavior
// when calling Validate(interface{}) or ValidateFields(interface{}).
type Validation interface {
	// Validate is called by the validation logic to check the validity of this value.
	Validate() ValidationResult
}
