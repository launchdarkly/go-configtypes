package configtypes

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustOptURLAbsolute(u *url.URL) OptURLAbsolute {
	o, err := NewOptURLAbsolute(u)
	if err != nil {
		panic(err)
	}
	return o
}

func TestOptURLAbsolute(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptURLAbsolute{}
		assertIsDefined(t, false, unsetValue)
		assert.Nil(t, unsetValue.Get())

		nilValue, err := NewOptURLAbsolute(nil)
		assert.NoError(t, err)
		assert.Equal(t, unsetValue, nilValue)
	})

	t.Run("defined value", func(t *testing.T) {
		absValue, err := NewOptURLAbsolute(absoluteURL)
		assert.NoError(t, err)
		assertIsDefined(t, true, absValue)
		assert.Equal(t, absoluteURL, absValue.Get())
	})

	t.Run("invalid value", func(t *testing.T) {
		relValue, err := NewOptURLAbsolute(relativeURL)
		assert.Error(t, err)
		assert.Equal(t, OptURLAbsolute{}, relValue)
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptURLAbsoluteFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptURLAbsolute{}, absoluteURLString: mustOptURLAbsolute(absoluteURL),
	})

	assertConvertFromText(t, &OptURLAbsolute{}, stringCtor, map[string]interface{}{
		"": OptURLAbsolute{}, absoluteURLString: mustOptURLAbsolute(absoluteURL),
	})

	assertConvertFromTextFails(t, &OptURLAbsolute{}, stringCtor, errURLNotAbsolute(),
		relativeURLString,
	)

	assertConvertFromTextFails(t, &OptURLAbsolute{}, stringCtor, errURLFormat(),
		malformedURLString,
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptURLAbsolute{}, quoteJSONString(absoluteURLString): mustOptURLAbsolute(absoluteURL),
	})

	assertConvertFromJSON(t, &OptURLAbsolute{}, map[string]interface{}{
		`null`: OptURLAbsolute{}, quoteJSONString(absoluteURLString): mustOptURLAbsolute(absoluteURL),
	})

	assertConvertFromJSONFails(t, &OptURLAbsolute{},
		`true`, `0.5`, quoteJSONString(relativeURLString), quoteJSONString(malformedURLString), `[]`, `{}`)
}
