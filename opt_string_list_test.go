package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptStringList(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptStringList{}
		assertIsDefined(t, false, unsetValue)
		assert.Nil(t, unsetValue.Values())

		nilValue := NewOptStringList(nil)
		assert.Equal(t, unsetValue, nilValue)
	})

	t.Run("defined value", func(t *testing.T) {
		emptyList := NewOptStringList([]string{})
		assertIsDefined(t, true, emptyList)
		assert.Len(t, emptyList.Values(), 0)

		fullList := NewOptStringList([]string{"a"})
		assertIsDefined(t, true, fullList)
		assert.Equal(t, []string{"a"}, fullList.Values())
	})

	stringCtor := func(input string) (interface{}, error) {
		o := NewOptStringListFromString(input)
		return o, nil // can't fail
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"":    OptStringList{},
		"a":   NewOptStringList([]string{"a"}),
		"a,b": NewOptStringList([]string{"a", "b"}),
	})

	// Need to call assertConvertFromText separately for each of these test cases because of the
	// additive behavior of UnmarshalText for this type.
	assertConvertFromText(t, &OptStringList{}, stringCtor, map[string]interface{}{
		"": OptStringList{},
	})
	assertConvertFromText(t, &OptStringList{}, stringCtor, map[string]interface{}{
		"a": NewOptStringList([]string{"a"}),
	})
	assertConvertFromText(t, &OptStringList{}, stringCtor, map[string]interface{}{
		"a,b": NewOptStringList([]string{"a", "b"}),
	})

	t.Run("multiple values with UnmarshalText", func(t *testing.T) {
		var o OptStringList
		assert.NoError(t, o.UnmarshalText([]byte("a")))
		assert.NoError(t, o.UnmarshalText([]byte("b")))
		assert.Equal(t, []string{"a", "b"}, o.Values())
	})

	assertConvertToJSON(t, map[string]SingleValue{
		`null`:      OptStringList{},
		`["a","b"]`: NewOptStringList([]string{"a", "b"}),
	})

	assertConvertFromJSON(t, &OptStringList{}, map[string]interface{}{
		`null`:      OptStringList{},
		`"a"`:       NewOptStringList([]string{"a"}),
		`["a","b"]`: NewOptStringList([]string{"a", "b"}),
	})

	assertConvertFromJSONFails(t, &OptStringList{},
		`true`, `1`, `["a",1]`, `{}`)
}
