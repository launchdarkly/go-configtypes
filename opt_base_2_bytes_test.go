package configtypes

import (
	"testing"

	"github.com/alecthomas/units"
	"github.com/stretchr/testify/assert"
)

const (
	gigString           = "10GiB"
	megString           = "3MiB"
	malformedSizeString = "7gb"
)

var (
	gigBytes, _ = units.ParseBase2Bytes(gigString)
	megBytes, _ = units.ParseBase2Bytes(megString)
)

func TestOptBase2Bytes(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		unsetValue := OptBase2Bytes{}
		assertIsDefined(t, false, unsetValue)
		assert.Equal(t, 3*units.KiB, unsetValue.GetOrElse(3*units.KiB))
		assert.Equal(t, gigBytes, unsetValue.GetOrElse(gigBytes))
	})

	t.Run("defined value", func(t *testing.T) {
		gigValue := NewOptBase2Bytes(gigBytes)
		assertIsDefined(t, true, gigValue)
		assert.Equal(t, gigBytes, gigValue.Get())
		assert.Equal(t, gigBytes, gigValue.GetOrElse(megBytes))

		megValue := NewOptBase2Bytes(megBytes)
		assertIsDefined(t, true, megValue)
		assert.Equal(t, megBytes, megValue.Get())
		assert.Equal(t, megBytes, megValue.GetOrElse(gigBytes))
	})

	stringCtor := func(input string) (interface{}, error) {
		o, err := NewOptBase2BytesFromString(input)
		return o, err
	}

	assertConvertToText(t, map[string]textMarshalerAndStringer{
		"": OptBase2Bytes{}, gigString: NewOptBase2Bytes(gigBytes), megString: NewOptBase2Bytes(megBytes),
	})

	assertConvertFromText(t, &OptBase2Bytes{}, stringCtor, map[string]interface{}{
		"": OptBase2Bytes{}, gigString: NewOptBase2Bytes(gigBytes), megString: NewOptBase2Bytes(megBytes),
	})

	assertConvertFromTextFails(t, &OptBase2Bytes{}, stringCtor, errBase2BytesFormat(),
		malformedSizeString,
	)

	assertConvertToJSON(t, map[string]SingleValue{
		`null`: OptBase2Bytes{}, quoteJSONString(gigString): NewOptBase2Bytes(gigBytes),
		quoteJSONString(megString): NewOptBase2Bytes(megBytes),
	})

	assertConvertFromJSON(t, &OptBase2Bytes{}, map[string]interface{}{
		`null`: OptBase2Bytes{}, quoteJSONString(gigString): NewOptBase2Bytes(gigBytes),
		quoteJSONString(megString): NewOptBase2Bytes(megBytes),
	})

	assertConvertFromJSONFails(t, &OptBase2Bytes{},
		`true`, `0.5`, quoteJSONString(malformedSizeString), `[]`, `{}`)
}
