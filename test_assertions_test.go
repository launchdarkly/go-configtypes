package configtypes

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

type textMarshalerAndStringer interface {
	encoding.TextMarshaler
	fmt.Stringer
}

func quoteJSONString(s string) string {
	return ldvalue.String(s).JSONString()
}

func assertIsDefined(t *testing.T, shouldBe bool, value SingleValue) {
	assert.Equal(t, shouldBe, value.IsDefined())
}

func assertConvertToText(
	t *testing.T,
	expectedStringFromValue map[string]textMarshalerAndStringer,
) {
	t.Run("convert to text with TextMarshaler and String", func(t *testing.T) {
		for expectedString, value := range expectedStringFromValue {
			bytes, err := value.MarshalText()
			if assert.NoError(t, err) {
				assert.Equal(t, expectedString, string(bytes))
			}

			assert.Equal(t, expectedString, value.String())
		}
	})
}

func assertConvertFromText(
	t *testing.T,
	zeroValue encoding.TextUnmarshaler,
	constructor func(string) (interface{}, error),
	values map[string]interface{},
) {
	t.Run("convert from text with constructor and UnmarshalText - valid values", func(t *testing.T) {
		for input, expected := range values {
			t.Run(input, func(t *testing.T) {
				actual, err := constructor(input)
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)

				err = zeroValue.UnmarshalText([]byte(input))
				assert.NoError(t, err)
				assert.Equal(t, expected, dereferenceIfPointer(zeroValue))
			})
		}
	})
}

func assertConvertFromTextFails(
	t *testing.T,
	zeroValue encoding.TextUnmarshaler,
	constructor func(string) (interface{}, error),
	expectedError Error,
	values ...string,
) {
	t.Run("convert from text with constructor and UnmarshalText - invalid values", func(t *testing.T) {
		for _, input := range values {
			t.Run(input, func(t *testing.T) {
				_, err := constructor(input)
				assert.Equal(t, expectedError, err)

				err = zeroValue.UnmarshalText([]byte(input))
				assert.Equal(t, expectedError, err)
			})
		}
	})
}

func assertConvertToJSON(
	t *testing.T,
	expectedJSONStringFromValue map[string]SingleValue,
) {
	t.Run("convert to JSON with MarshalJSON - valid values", func(t *testing.T) {
		for expectedJSONString, value := range expectedJSONStringFromValue {
			bytes, err := value.MarshalJSON()
			if assert.NoError(t, err) {
				assert.JSONEq(t, expectedJSONString, string(bytes))
			}

			assert.JSONEq(t, expectedJSONString, string(bytes))
		}
	})
}

func assertConvertFromJSON(
	t *testing.T,
	zeroValue json.Unmarshaler,
	jsonValues map[string]interface{},
) {
	t.Run("convert from JSON with UnmarshalJSON - valid values", func(t *testing.T) {
		for input, expected := range jsonValues {
			t.Run(input, func(t *testing.T) {
				err := zeroValue.UnmarshalJSON([]byte(input))
				assert.NoError(t, err)
				assert.Equal(t, expected, dereferenceIfPointer(zeroValue))
			})
		}
	})
}

func assertConvertFromJSONFails(
	t *testing.T,
	zeroValue json.Unmarshaler,
	values ...string,
) {
	values = append(values, `some invalid JSON`)
	t.Run("convert from JSON with UnmarshalJSON - invalid values", func(t *testing.T) {
		for _, input := range values {
			t.Run(input, func(t *testing.T) {
				assert.Error(t, zeroValue.UnmarshalJSON([]byte(input)))
			})
		}
	})
}

func dereferenceIfPointer(value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		return rv.Elem().Interface()
	}
	return value
}
