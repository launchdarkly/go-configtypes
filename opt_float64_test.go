package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptFloat64(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptFloat64{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, float64(999), unsetValue.GetOrElse(999))
	})

	t.Run("defined value", func(t *testing.T) {
		zeroValue := NewOptFloat64(0)
		assertIsDefined(t, true, zeroValue)
		assert.Equal(t, float64(0), zeroValue.GetOrElse(999))

		negativeValue := NewOptFloat64(-1)
		assertIsDefined(t, true, negativeValue)
		assert.Equal(t, float64(-1), negativeValue.GetOrElse(999))

		oneValue := NewOptFloat64(1.5)
		assertIsDefined(t, true, oneValue)
		assert.Equal(t, float64(1.5), oneValue.GetOrElse(0))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptFloat64FromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptFloat64{}, "0": NewOptFloat64(0), "1.5": NewOptFloat64(1.5),
		"100": NewOptFloat64(100), "-100": NewOptFloat64(-100),
	})

	assertConvertFromText(t, &OptFloat64{}, stringCtor, map[string]interface{}{
		"": OptFloat64{}, "0": NewOptFloat64(0), "1.5": NewOptFloat64(1.5),
		"100": NewOptFloat64(100), "-100": NewOptFloat64(-100),
	})

	assertConvertFromTextFails(t, &OptFloat64{}, stringCtor, errFloatFormat(),
		"-", "x",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptFloat64{}, `0`: NewOptFloat64(0), `1.5`: NewOptFloat64(1.5),
		`100`: NewOptFloat64(100), `-100`: NewOptFloat64(-100),
	})

	assertConvertFromJSON(t, &OptFloat64{}, map[string]interface{}{
		`null`: OptFloat64{}, `0`: NewOptFloat64(0), `1.5`: NewOptFloat64(1.5),
		`100`: NewOptFloat64(100), `-100`: NewOptFloat64(-100),
	})

	assertConvertFromJSONFails(t, &OptFloat64{},
		`true`, `"x"`, `[]`, `{}`)
}
