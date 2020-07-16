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
		value := NewOptDuration(time.Minute)
		assertIsDefined(t, true, value)
		assert.Equal(t, time.Minute, value.GetOrElse(time.Minute))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptDurationFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"":        OptDuration{},
		"3ms":     NewOptDuration(3 * time.Millisecond),
		"3s":      NewOptDuration(3 * time.Second),
		"3m":      NewOptDuration(3 * time.Minute),
		"3h":      NewOptDuration(3 * time.Hour),
		"1:30":    NewOptDuration(time.Minute + 30*time.Second),
		"1:10:30": NewOptDuration(time.Hour + 10*time.Minute + 30*time.Second),
	})

	assertConvertFromText(t, &OptDuration{}, stringCtor, map[string]interface{}{
		"":        OptDuration{},
		"3ms":     NewOptDuration(3 * time.Millisecond),
		"3s":      NewOptDuration(3 * time.Second),
		"3m":      NewOptDuration(3 * time.Minute),
		"3h":      NewOptDuration(3 * time.Hour),
		":30":     NewOptDuration(30 * time.Second),
		"1:30":    NewOptDuration(time.Minute + 30*time.Second),
		"1:10:30": NewOptDuration(time.Hour + 10*time.Minute + 30*time.Second),
	})

	assertConvertFromTextFails(t, &OptDuration{}, stringCtor, errDurationFormat(),
		"1", "x", "1x", "00:30:",
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
