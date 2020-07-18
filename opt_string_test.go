package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptString(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptString{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, "z", unsetValue.GetOrElse("z"))
	})

	t.Run("defined value", func(t *testing.T) {
		emptyString := NewOptString("")
		assertIsDefined(t, true, emptyString)
		assert.Equal(t, "", emptyString.GetOrElse("z"))

		fullString := NewOptString("a")
		assertIsDefined(t, true, fullString)
		assert.Equal(t, "a", fullString.GetOrElse("z"))
	})

	stringCtor := func(input string) (interface{}, error) {
		o := NewOptString(input)
		return o, nil // can't fail
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": NewOptString(""), "a": NewOptString("a"),
	})

	assertConvertFromText(t, &OptString{}, stringCtor, map[string]interface{}{
		"": NewOptString(""), "a": NewOptString("a"),
	})

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptString{}, `""`: NewOptString(""), `"a"`: NewOptString("a"),
	})

	assertConvertFromJSON(t, &OptString{}, map[string]interface{}{
		`null`: OptString{}, `""`: NewOptString(""), `"a"`: NewOptString("a"),
	})

	assertConvertFromJSONFails(t, &OptString{},
		`true`, `1`, `[]`, `{}`)
}
