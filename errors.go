package configtypes

import (
	"errors"
)

// Error is a type tag for all errors returned by this package.
type Error error

func errBoolFormat() Error {
	return errors.New("not a valid boolean value (must be true/false, yes/no, or 0/1)")
}

func errDurationFormat() Error {
	return errors.New(
		`not a valid duration (must use format "1ms", "1s", "1m", etc.)`,
	)
}

func errIntFormat() Error {
	return errors.New("not a valid integer")
}

func errFloatFormat() Error {
	return errors.New("not a valid number")
}

func errMustBeGreaterThanZero() Error {
	return errors.New("value must be greater than zero")
}

func errMustBeNonEmptyString() Error {
	return errors.New("value must not be an empty string")
}

func errMustBeNonNegative() Error {
	return errors.New("value must not be negative")
}

func errRequired() Error {
	return errors.New("value is required")
}

func errStringListJSONFormat() Error {
	return errors.New("string list value must be a string, an array of strings, or null")
}

func errURLFormat() Error {
	return errors.New("not a valid URL/URI")
}

func errURLNotAbsolute() Error {
	return errors.New("must be an absolute URL/URI")
}

func errValidateNonStruct() Error {
	return errors.New("Validate was called with a parameter that was not a struct pointer") //nolint:stylecheck
}
