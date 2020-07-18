package configtypes

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	relativeURLString  = "relative/url"
	absoluteURLString  = "http://absolute/url"
	malformedURLString = "::"
)

var (
	relativeURL, _ = url.Parse(relativeURLString)
	absoluteURL, _ = url.Parse(absoluteURLString)
)

func TestOptURL(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptURL{}
		assertIsDefined(t, false, unsetValue)
		assert.Nil(t, unsetValue.Get())

		nilValue := NewOptURL(nil)
		assert.Equal(t, unsetValue, nilValue)
	})

	t.Run("defined value", func(t *testing.T) {
		relValue := NewOptURL(relativeURL)
		assertIsDefined(t, true, relValue)
		assert.Equal(t, relativeURL, relValue.Get())

		absValue := NewOptURL(absoluteURL)
		assertIsDefined(t, true, absValue)
		assert.Equal(t, absoluteURL, absValue.Get())
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptURLFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptURL{}, relativeURLString: NewOptURL(relativeURL), absoluteURLString: NewOptURL(absoluteURL),
	})

	assertConvertFromText(t, &OptURL{}, stringCtor, map[string]interface{}{
		"": OptURL{}, relativeURLString: NewOptURL(relativeURL), absoluteURLString: NewOptURL(absoluteURL),
	})

	assertConvertFromTextFails(t, &OptURL{}, stringCtor, errURLFormat(),
		malformedURLString,
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptURL{}, quoteJSONString(relativeURLString): NewOptURL(relativeURL),
		quoteJSONString(absoluteURLString): NewOptURL(absoluteURL),
	})

	assertConvertFromJSON(t, &OptURL{}, map[string]interface{}{
		`null`: OptURL{}, quoteJSONString(relativeURLString): NewOptURL(relativeURL),
		quoteJSONString(absoluteURLString): NewOptURL(absoluteURL),
	})

	assertConvertFromJSONFails(t, &OptURL{},
		`true`, `0.5`, quoteJSONString(malformedURLString), `[]`, `{}`)
}
