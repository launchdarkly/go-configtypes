package configtypes

import "time"

// OptDurationNonNegative represents an optional time.Duration parameter which, if defined, must be
// greater than or equal to zero.
//
// This is the same as OptDuration, but with additional validation for the constructor and unmarshalers.
// It is impossible (except with reflection) for code outside this package to construct an instance of this
// type with a defined value that is negative.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptDurationNonNegative struct {
	opt OptDuration
}

func NewOptDurationNonNegative(value time.Duration) (OptDurationNonNegative, error) {
	return optDurationNonNegativeFromOptDuration(NewOptDuration(value))
}

func NewOptDurationNonNegativeFromString(s string) (OptDurationNonNegative, error) {
	o, err := NewOptDurationFromString(s)
	if err != nil {
		return OptDurationNonNegative{}, err
	}
	return optDurationNonNegativeFromOptDuration(o)
}

func optDurationNonNegativeFromOptDuration(o OptDuration) (OptDurationNonNegative, error) {
	if !o.IsDefined() || o.GetOrElse(0) >= 0 {
		return OptDurationNonNegative{o}, nil
	}
	return OptDurationNonNegative{}, errMustBeNonNegative()
}

func (o OptDurationNonNegative) IsDefined() bool {
	return o.opt.IsDefined()
}

func (o OptDurationNonNegative) GetOrElse(orElseValue time.Duration) time.Duration {
	return o.opt.GetOrElse(orElseValue)
}

func (o *OptDurationNonNegative) UnmarshalText(data []byte) error {
	var opt OptDuration
	if err := opt.UnmarshalText(data); err != nil {
		return err
	}
	value, err := optDurationNonNegativeFromOptDuration(opt)
	if err == nil {
		*o = value
	}
	return err
}

func (o OptDurationNonNegative) String() string {
	return o.opt.String()
}

func (o OptDurationNonNegative) MarshalText() ([]byte, error) {
	return o.opt.MarshalText()
}

func (o *OptDurationNonNegative) UnmarshalJSON(data []byte) error {
	var opt OptDuration
	if err := opt.UnmarshalJSON(data); err != nil {
		return err
	}
	value, err := optDurationNonNegativeFromOptDuration(opt)
	if err == nil {
		*o = value
	}
	return err
}

func (o OptDurationNonNegative) MarshalJSON() ([]byte, error) {
	return o.opt.MarshalJSON()
}
