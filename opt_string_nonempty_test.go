package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptStringNonEmpty(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptStringNonEmpty{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, "z", unsetValue.GetOrElse("z"))

		emptyString := NewOptStringNonEmpty("")
		assert.Equal(t, unsetValue, emptyString)
	})

	t.Run("defined value", func(t *testing.T) {
		fullString := NewOptStringNonEmpty("a")
		assertIsDefined(t, true, fullString)
		assert.Equal(t, "a", fullString.GetOrElse("z"))
	})

	stringCtor := func(input string) (interface{}, error) {
		o := NewOptStringNonEmpty(input)
		return o, nil // can't fail
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptStringNonEmpty{}, "a": NewOptStringNonEmpty("a"),
	})

	assertConvertFromText(t, &OptStringNonEmpty{}, stringCtor, map[string]interface{}{
		"": OptStringNonEmpty{}, "a": NewOptStringNonEmpty("a"),
	})

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptStringNonEmpty{}, `"a"`: NewOptStringNonEmpty("a"),
	})

	assertConvertFromJSON(t, &OptStringNonEmpty{}, map[string]interface{}{
		`null`: OptStringNonEmpty{}, `"a"`: NewOptStringNonEmpty("a"),
	})

	assertConvertFromJSONFails(t, &OptStringNonEmpty{},
		`true`, `1`, `""`, `[]`, `{}`)
}
