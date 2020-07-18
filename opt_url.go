package configtypes

import (
	"encoding/json"
	"net/url"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

// OptURL represents an optional parameter which, if present, must be a valid URL.
//
// There is no additional validation, so the URL could be absolute or relative; see OptURL. The value is
// stored as a *URL pointer, whose value is always copied when accessed; a nil pointer indicates the lack
// of a value (the same as OptURL{}).
//
// Converting to or from a string uses the standard behavior for URL.String() and url.Parse().
//
// When converting to or from JSON, the value must be either a JSON null or a JSON string.
//
// See the package documentation for the general contract for methods that have no specific documentation
// here.
type OptURL struct {
	url *url.URL
}

func NewOptURL(url *url.URL) OptURL {
	if url == nil {
		return OptURL{}
	}
	u := *url
	return OptURL{url: &u}
}

func NewOptURLFromString(urlString string) (OptURL, error) {
	if urlString == "" {
		return OptURL{}, nil
	}
	u, err := url.Parse(urlString)
	if err == nil {
		return NewOptURL(u), nil
	}
	return OptURL{}, errURLFormat()
}

func (o OptURL) IsDefined() bool {
	return o.url != nil
}

func (o OptURL) Get() *url.URL {
	if o.url == nil {
		return nil
	}
	u := *o.url
	return &u
}

func (o OptURL) String() string {
	if o.url == nil {
		return ""
	}
	return o.url.String()
}

func (o OptURL) MarshalText() ([]byte, error) {
	if o.url == nil {
		return nil, nil
	}
	return []byte(o.String()), nil
}

func (o *OptURL) UnmarshalText(data []byte) error {
	parsed, err := NewOptURLFromString(string(data))
	if err == nil {
		*o = parsed
	}
	return err
}

func (o OptURL) MarshalJSON() ([]byte, error) {
	if o.url != nil {
		return json.Marshal(o.url.String())
	}
	return json.Marshal(nil)
}

func (o *OptURL) UnmarshalJSON(data []byte) error {
	var v ldvalue.Value
	var err error
	if err = v.UnmarshalJSON(data); err != nil {
		return err
	}
	switch {
	case v.IsNull():
		*o = OptURL{}
		return nil
	case v.IsString():
		*o, err = NewOptURLFromString(v.StringValue())
		return err
	default:
		return errURLFormat()
	}
}
