package configtypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mustOptDurationNonNegative(value time.Duration) OptDurationNonNegative {
	o, err := NewOptDurationNonNegative(value)
	if err != nil {
		panic(err)
	}
	return o
}

func TestOptDurationNonNegative(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptDurationNonNegative{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, time.Hour, unsetValue.GetOrElse(time.Hour))
	})

	t.Run("defined value", func(t *testing.T) {
		minuteValue, err := NewOptDurationNonNegative(time.Minute)
		assert.NoError(t, err)
		assertIsDefined(t, true, minuteValue)
		assert.Equal(t, time.Minute, minuteValue.GetOrElse(time.Hour))

		zeroValue, err := NewOptDurationNonNegative(0)
		assert.NoError(t, err)
		assertIsDefined(t, true, zeroValue)
		assert.Equal(t, time.Duration(0), zeroValue.GetOrElse(time.Hour))
	})

	t.Run("invalid value", func(t *testing.T) {
		negativeValue, err := NewOptDurationNonNegative(-1 * time.Millisecond)
		assert.Equal(t, errMustBeNonNegative(), err)
		assert.Equal(t, OptDurationNonNegative{}, negativeValue)
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptDurationNonNegativeFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"":         OptDurationNonNegative{},
		"3ms":      mustOptDurationNonNegative(3 * time.Millisecond),
		"3s":       mustOptDurationNonNegative(3 * time.Second),
		"3m0s":     mustOptDurationNonNegative(3 * time.Minute),
		"3h0m0s":   mustOptDurationNonNegative(3 * time.Hour),
		"1m30s":    mustOptDurationNonNegative(time.Minute + 30*time.Second),
		"1h10m30s": mustOptDurationNonNegative(time.Hour + 10*time.Minute + 30*time.Second),
	})

	assertConvertFromText(t, &OptDurationNonNegative{}, stringCtor, map[string]interface{}{
		"":         OptDurationNonNegative{},
		"3ms":      mustOptDurationNonNegative(3 * time.Millisecond),
		"3s":       mustOptDurationNonNegative(3 * time.Second),
		"3m":       mustOptDurationNonNegative(3 * time.Minute),
		"3h":       mustOptDurationNonNegative(3 * time.Hour),
		"1m30s":    mustOptDurationNonNegative(time.Minute + 30*time.Second),
		"1h10m30s": mustOptDurationNonNegative(time.Hour + 10*time.Minute + 30*time.Second),
	})

	assertConvertFromTextFails(t, &OptDurationNonNegative{}, stringCtor, errDurationFormat(),
		"1", "x", "1x", ":30",
	)

	assertConvertFromTextFails(t, &OptDurationNonNegative{}, stringCtor, errMustBeNonNegative(),
		"-1s",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptDurationNonNegative{},
		`"3s"`: mustOptDurationNonNegative(3 * time.Second),
	})

	assertConvertFromJSON(t, &OptDurationNonNegative{}, map[string]interface{}{
		`null`: OptDurationNonNegative{},
		`"3s"`: mustOptDurationNonNegative(3 * time.Second),
	})

	assertConvertFromJSONFails(t, &OptDurationNonNegative{},
		`true`, `1`, `"x"`, `"-1s"`, `[]`, `{}`)
}
