package configtypes

import (
	"fmt"
	"strings"
)

// ValidationResult accumulates errors from field validation.
type ValidationResult struct {
	errors []ValidationError
}

// ValidationPath represents a field name or nested series of field names for a ValidationError.
//
// For instance, if you call Validate on a struct whose field A is a struct that has an invalid
// field B, the path for the error would be ValidationPath{"A", "B"}.
type ValidationPath []string

// String converts the path to a dot-delimited string.
func (p ValidationPath) String() string {
	if len(p) == 0 {
		return ""
	}
	return strings.Join(p, ".")
}

// ValidationError represents an invalid value condition for a parsed value or a struct field.
type ValidationError struct {
	Path ValidationPath
	Err  error
}

// Error returns the error description, including the path if specified.
func (v ValidationError) Error() string {
	if len(v.Path) == 0 {
		return v.Err.Error()
	}
	return fmt.Sprintf("%s: %s", v.Path, v.Err.Error())
}

// String is equivalent to Error.
func (v ValidationError) String() string {
	return v.Error()
}

// ValidationAggregateError is the type returned by ValidationResult.GetError() if there were
// multiple errors.
type ValidationAggregateError []ValidationError

// Error for ValidationAggregateError returns a comma-delimited list of error descriptions.
func (v ValidationAggregateError) Error() string {
	ss := make([]string, 0, len(v))
	for _, err := range v {
		ss = append(ss, err.String())
	}
	return strings.Join(ss, ", ")
}

// String is equivalent to Error.
func (v ValidationAggregateError) String() string {
	return v.Error()
}

// OK returns true if there are no errors.
func (r ValidationResult) OK() bool {
	return len(r.errors) == 0
}

// Errors returns a copied slice of all errors in this result.
func (r ValidationResult) Errors() []ValidationError {
	ret := make([]ValidationError, len(r.errors))
	copy(ret, r.errors)
	return ret
}

// GetError returns a single error representing all errors in the result, or nil if there were none.
//
// If not nil, the return value will be either a ValidationError or a ValidationAggregateError.
func (r ValidationResult) GetError() error {
	if len(r.errors) == 0 {
		return nil
	}
	if len(r.errors) == 1 {
		return r.errors[0]
	}
	return ValidationAggregateError(r.Errors())
}

// AddError adds a ValidationError to the result.
func (r *ValidationResult) AddError(path ValidationPath, e error) {
	if e != nil {
		r.errors = append(r.errors, ValidationError{Path: path, Err: e})
	}
}

// AddAll adds all errors from another result, optionally adding a prefix to each path.
func (r *ValidationResult) AddAll(prefixPath ValidationPath, other ValidationResult) {
	for _, e := range other.errors {
		r.errors = append(r.errors, ValidationError{Path: append(prefixPath, e.Path...), Err: e.Err})
	}
}
