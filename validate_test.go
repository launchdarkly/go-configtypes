package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	t.Run("rejects parameter that is not a struct", func(t *testing.T) {
		anInt := 3
		assert.Error(t, ValidateStruct(anInt, false).GetError())
		assert.Error(t, ValidateStruct(&anInt, false).GetError())
		assert.NoError(t, ValidateStruct(structToValidateNoRequirements{}, false).GetError())
	})

	t.Run("allows optional fields to be set or unset", func(t *testing.T) {
		s1 := structToValidateNoRequirements{}
		assert.NoError(t, ValidateStruct(&s1, false).GetError())

		s2 := structToValidateNoRequirements{Int: NewOptInt(3), Str: "x"}
		assert.NoError(t, ValidateStruct(&s2, false).GetError())
	})

	t.Run("requires required fields to be set", func(t *testing.T) {
		s1 := structToValidateWithRequirements{}
		assert.Equal(t, []ValidationError{
			{Path: ValidationPath{"Int"}, Err: errRequired()},
			{Path: ValidationPath{"Str"}, Err: errRequired()},
		}, ValidateStruct(&s1, false).Errors())

		s2 := structToValidateWithRequirements{Str: "x"}
		assert.Equal(t, []ValidationError{
			{Path: ValidationPath{"Int"}, Err: errRequired()},
		}, ValidateStruct(&s2, false).Errors())

		s3 := structToValidateWithRequirements{Int: NewOptInt(3), Str: "x"}
		assert.NoError(t, ValidateStruct(&s3, false).GetError())
	})

	t.Run("ignores nested struct when recursive is false", func(t *testing.T) {
		s := structWithNestedStructWithRequirements{TopLevelInt: NewOptInt(3)}
		assert.NoError(t, ValidateStruct(&s, false).GetError())
	})

	t.Run("validates nested struct when recursive is true", func(t *testing.T) {
		s := structWithNestedStructWithRequirements{TopLevelInt: NewOptInt(3)}
		assert.Equal(t, []ValidationError{
			{Path: ValidationPath{"Nested", "Int"}, Err: errRequired()},
			{Path: ValidationPath{"Nested", "Str"}, Err: errRequired()},
		}, ValidateStruct(&s, true).Errors())
	})

	t.Run("both top-level and nested fields are validated when recursive is true", func(t *testing.T) {
		s := structWithNestedStructWithRequirements{}
		assert.Equal(t, []ValidationError{
			{Path: ValidationPath{"TopLevelInt"}, Err: errRequired()},
			{Path: ValidationPath{"Nested", "Int"}, Err: errRequired()},
			{Path: ValidationPath{"Nested", "Str"}, Err: errRequired()},
		}, ValidateStruct(&s, true).Errors())
	})

	t.Run("logs error for invalid conf tag", func(t *testing.T) {
		s := testStructWithBadTag{}
		assert.Error(t, ValidateStruct(&s, false).GetError())
	})
}

type mockValidation struct {
	result ValidationResult
}

func (m *mockValidation) Validate() ValidationResult {
	return m.result
}

type structToValidateNoRequirements struct {
	Int OptInt
	Str string
}

type structToValidateWithRequirements struct {
	Int                        OptInt `conf:",required"`
	Str                        string `conf:",required"`
	ignoreThisNonExportedField OptInt `conf:",required"` //nolint:unused,structcheck
}

type structWithNestedStructWithRequirements struct {
	TopLevelInt OptInt `conf:",required"`
	Nested      structToValidateWithRequirements
}
