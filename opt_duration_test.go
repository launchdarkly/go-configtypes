package configtypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptDuration(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptDuration{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, time.Hour, unsetValue.GetOrElse(time.Hour))
	})

	t.Run("defined value", func(t *testing.T) {
		minuteValue := NewOptDuration(time.Minute)
		assertIsDefined(t, true, minuteValue)
		assert.Equal(t, time.Minute, minuteValue.GetOrElse(time.Hour))

		zeroValue := NewOptDuration(0)
		assertIsDefined(t, true, zeroValue)
		assert.Equal(t, time.Duration(0), zeroValue.GetOrElse(time.Hour))

		negativeValue := NewOptDuration(-1 * time.Millisecond)
		assertIsDefined(t, true, negativeValue)
		assert.Equal(t, -1*time.Millisecond, negativeValue.GetOrElse(time.Hour))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptDurationFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"":         OptDuration{},
		"3ms":      NewOptDuration(3 * time.Millisecond),
		"3s":       NewOptDuration(3 * time.Second),
		"3m0s":     NewOptDuration(3 * time.Minute),
		"3h0m0s":   NewOptDuration(3 * time.Hour),
		"1m30s":    NewOptDuration(time.Minute + 30*time.Second),
		"1h10m30s": NewOptDuration(time.Hour + 10*time.Minute + 30*time.Second),
		"-1s":      NewOptDuration(-1 * time.Second),
	})

	assertConvertFromText(t, &OptDuration{}, stringCtor, map[string]interface{}{
		"":         OptDuration{},
		"3ms":      NewOptDuration(3 * time.Millisecond),
		"3s":       NewOptDuration(3 * time.Second),
		"3m":       NewOptDuration(3 * time.Minute),
		"3h":       NewOptDuration(3 * time.Hour),
		"1m30s":    NewOptDuration(time.Minute + 30*time.Second),
		"1h10m30s": NewOptDuration(time.Hour + 10*time.Minute + 30*time.Second),
		"-1s":      NewOptDuration(-1 * time.Second),
	})

	assertConvertFromTextFails(t, &OptDuration{}, stringCtor, errDurationFormat(),
		"1", "x", "1x", ":30",
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptDuration{},
		`"3s"`: NewOptDuration(3 * time.Second),
	})

	assertConvertFromJSON(t, &OptDuration{}, map[string]interface{}{
		`null`: OptDuration{},
		`"3s"`: NewOptDuration(3 * time.Second),
	})

	assertConvertFromJSONFails(t, &OptDuration{},
		`true`, `1`, `"x"`, `[]`, `{}`)
}
