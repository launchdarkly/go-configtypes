package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustOptIntGreaterThanZero(n int) OptIntGreaterThanZero {
	o, err := NewOptIntGreaterThanZero(n)
	if err != nil {
		panic(err)
	}
	return o
}

func TestOptIntGreaterThanZero(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptIntGreaterThanZero{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, 999, unsetValue.GetOrElse(999))
	})

	t.Run("defined value", func(t *testing.T) {
		oneValue, err := NewOptIntGreaterThanZero(1)
		assert.NoError(t, err)
		assertIsDefined(t, true, oneValue)
		assert.Equal(t, 1, oneValue.GetOrElse(0))
	})

	t.Run("invalid value", func(t *testing.T) {
		zeroValue, err := NewOptIntGreaterThanZero(0)
		assert.Equal(t, errMustBeGreaterThanZero(), err)
		assert.Equal(t, OptIntGreaterThanZero{}, zeroValue)

		negativeValue, err := NewOptIntGreaterThanZero(-1)
		assert.Equal(t, errMustBeGreaterThanZero(), err)
		assert.Equal(t, OptIntGreaterThanZero{}, negativeValue)
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptIntGreaterThanZeroFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptIntGreaterThanZero{}, "100": mustOptIntGreaterThanZero(100),
	})

	assertConvertFromText(t, &OptIntGreaterThanZero{}, stringCtor, map[string]interface{}{
		"": OptIntGreaterThanZero{}, "100": mustOptIntGreaterThanZero(100),
	})

	assertConvertFromTextFails(t, &OptIntGreaterThanZero{}, stringCtor, errIntFormat(),
		"-", "0.5", "x",
	)

	assertConvertFromTextFails(t, &OptIntGreaterThanZero{}, stringCtor, errMustBeGreaterThanZero(),
		"0", "-1",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptIntGreaterThanZero{}, `100`: mustOptIntGreaterThanZero(100),
	})

	assertConvertFromJSON(t, &OptIntGreaterThanZero{}, map[string]interface{}{
		`null`: OptIntGreaterThanZero{}, `100`: mustOptIntGreaterThanZero(100),
	})

	assertConvertFromJSONFails(t, &OptIntGreaterThanZero{},
		`true`, `0`, `-1`, `0.5`, `"x"`, `[]`, `{}`)
}
