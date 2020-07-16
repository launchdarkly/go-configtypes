package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptBool(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptBool{}
		assertIsDefined(t, false, unsetValue)
		assert.False(t, unsetValue.GetOrElse(false))
		assert.True(t, unsetValue.GetOrElse(true))
	})

	t.Run("defined value", func(t *testing.T) {
		trueValue := NewOptBool(true)
		assertIsDefined(t, true, trueValue)
		assert.True(t, trueValue.GetOrElse(true))
		assert.True(t, trueValue.GetOrElse(false))

		falseValue := NewOptBool(false)
		assertIsDefined(t, true, falseValue)
		assert.False(t, falseValue.GetOrElse(true))
		assert.False(t, falseValue.GetOrElse(false))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptBoolFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptBool{}, "true": NewOptBool(true), "false": NewOptBool(false),
	})

	assertConvertFromText(t, &OptBool{}, stringCtor, map[string]interface{}{
		"":     OptBool{},
		"true": NewOptBool(true), "false": NewOptBool(false),
		"1": NewOptBool(true), "0": NewOptBool(false),
		"yes": NewOptBool(true), "no": NewOptBool(false),
	})

	assertConvertFromTextFails(t, &OptBool{}, stringCtor, errBoolFormat(),
		"maybe", "2",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptBool{}, `true`: NewOptBool(true), `false`: NewOptBool(false),
	})

	assertConvertFromJSON(t, &OptBool{}, map[string]interface{}{
		`null`: OptBool{}, `true`: NewOptBool(true), `false`: NewOptBool(false),
	})

	assertConvertFromJSONFails(t, &OptBool{},
		`1`, `"x"`, `[]`, `{}`)
}
