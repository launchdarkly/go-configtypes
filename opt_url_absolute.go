package configtypes

import (
	"net/url"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptURLAbsolute represents an optional URL parameter which, if defined, must be an absolute URL.
//
// This is the same as OptURL, but with additional validation for the constructor and unmarshalers. It
// is impossible (except with reflection) for code outside this package to construct an instance of this
// type with a defined value that is not a valid absolute URL.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptURLAbsolute struct {
	opt OptURL
}

func NewOptURLAbsolute(url *url.URL) (OptURLAbsolute, error) {
	if url != nil && !url.IsAbs() {
		return OptURLAbsolute{}, errURLNotAbsolute()
	}
	return OptURLAbsolute{opt: NewOptURL(url)}, nil
}

func NewOptURLAbsoluteFromString(urlString string) (OptURLAbsolute, error) {
	opt, err := NewOptURLFromString(urlString)
	if err != nil {
		return OptURLAbsolute{}, err
	}
	return NewOptURLAbsolute(opt.Get())
}

func (o OptURLAbsolute) IsDefined() bool {
	return o.opt.IsDefined()
}

func (o OptURLAbsolute) Get() *url.URL {
	return o.opt.Get()
}

func (o OptURLAbsolute) String() string {
	return o.opt.String()
}

func (o OptURLAbsolute) MarshalText() ([]byte, error) {
	return o.opt.MarshalText()
}

func (o *OptURLAbsolute) UnmarshalText(data []byte) error {
	parsed, err := NewOptURLAbsoluteFromString(string(data))
	if err == nil {
		*o = parsed
	}
	return err
}

func (o OptURLAbsolute) MarshalJSON() ([]byte, error) {
	return o.opt.MarshalJSON()
}

func (o *OptURLAbsolute) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	var err error
	if err = v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptURLAbsolute{}
		return nil
	case v.IsString():
		*o, err = NewOptURLAbsoluteFromString(v.StringValue())
		return err
	default:
		return errURLFormat()
	}
}
