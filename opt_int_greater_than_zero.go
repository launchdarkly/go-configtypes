package configtypes

// OptIntGreaterThanZero represents an optional int parameter which, if defined, must be greater than zero.
//
// This is the same as OptInt, but with additional validation for the constructor and unmarshalers. It
// is impossible (except with reflection) for code outside this package to construct an instance of this
// type with a defined value that is zero or negative.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptIntGreaterThanZero struct {
	opt OptInt
}

func NewOptIntGreaterThanZero(value int) (OptIntGreaterThanZero, error) {
	return optIntGreaterThanZeroFromOptInt(NewOptInt(value))
}

func NewOptIntGreaterThanZeroFromString(s string) (OptIntGreaterThanZero, error) {
	o, err := NewOptIntFromString(s)
	if err != nil {
		return OptIntGreaterThanZero{}, err
	}
	return optIntGreaterThanZeroFromOptInt(o)
}

func optIntGreaterThanZeroFromOptInt(o OptInt) (OptIntGreaterThanZero, error) {
	if !o.IsDefined() || o.GetOrElse(0) > 0 {
		return OptIntGreaterThanZero{o}, nil
	}
	return OptIntGreaterThanZero{}, errMustBeGreaterThanZero()
}

func (o OptIntGreaterThanZero) IsDefined() bool {
	return o.opt.IsDefined()
}

func (o OptIntGreaterThanZero) GetOrElse(orElseValue int) int {
	return o.opt.GetOrElse(orElseValue)
}

func (o *OptIntGreaterThanZero) UnmarshalText(data []byte) error {
	var opt OptInt
	if err := opt.UnmarshalText(data); err != nil {
		return err
	}
	value, err := optIntGreaterThanZeroFromOptInt(opt)
	if err == nil {
		*o = value
	}
	return err
}

func (o OptIntGreaterThanZero) String() string {
	return o.opt.String()
}

func (o OptIntGreaterThanZero) MarshalText() ([]byte, error) {
	return o.opt.MarshalText()
}

func (o *OptIntGreaterThanZero) UnmarshalJSON(data []byte) error {
	var opt OptInt
	if err := opt.UnmarshalJSON(data); err != nil {
		return err
	}
	value, err := optIntGreaterThanZeroFromOptInt(opt)
	if err == nil {
		*o = value
	}
	return err
}

func (o OptIntGreaterThanZero) MarshalJSON() ([]byte, error) {
	return o.opt.MarshalJSON()
}
