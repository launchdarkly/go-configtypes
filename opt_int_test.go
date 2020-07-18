package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptInt(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptInt{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, 999, unsetValue.GetOrElse(999))
	})

	t.Run("defined value", func(t *testing.T) {
		zeroValue := NewOptInt(0)
		assertIsDefined(t, true, zeroValue)
		assert.Equal(t, 0, zeroValue.GetOrElse(999))

		negativeValue := NewOptInt(-1)
		assertIsDefined(t, true, negativeValue)
		assert.Equal(t, -1, negativeValue.GetOrElse(999))

		oneValue := NewOptInt(1)
		assertIsDefined(t, true, oneValue)
		assert.Equal(t, 1, oneValue.GetOrElse(0))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptIntFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptInt{}, "0": NewOptInt(0), "100": NewOptInt(100), "-100": NewOptInt(-100),
	})

	assertConvertFromText(t, &OptInt{}, stringCtor, map[string]interface{}{
		"": OptInt{}, "0": NewOptInt(0), "100": NewOptInt(100), "-100": NewOptInt(-100),
	})

	assertConvertFromTextFails(t, &OptInt{}, stringCtor, errIntFormat(),
		"-", "0.5", "x",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptInt{}, `0`: NewOptInt(0), `100`: NewOptInt(100), `-100`: NewOptInt(-100),
	})

	assertConvertFromJSON(t, &OptInt{}, map[string]interface{}{
		`null`: OptInt{}, `0`: NewOptInt(0), `100`: NewOptInt(100), `-100`: NewOptInt(-100),
	})

	assertConvertFromJSONFails(t, &OptInt{},
		`true`, `0.5`, `"x"`, `[]`, `{}`)
}
