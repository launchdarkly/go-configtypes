package configtypes

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVarReader(t *testing.T) {
	t.Run("can read from environment", func(t *testing.T) {
		withCleanEnvVars(func() {
			os.Setenv("NAME1", "value1")
			os.Setenv("NAME2", "value2=xyz")
			r := NewVarReaderFromEnvironment()

			var v1, v2, v3 mockTextUnmarshaler
			assert.True(t, r.Read("NAME1", &v1))
			assert.True(t, r.Read("NAME2", &v2))
			assert.False(t, r.Read("NAME3", &v3))
			assert.Equal(t, "value1", v1.value)
			assert.Equal(t, "value2=xyz", v2.value)
		})
	})

	t.Run("reads into TextUnmarshaler", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME": "value"})
		var v mockTextUnmarshaler
		assert.True(t, r.Read("NAME", &v))

		assert.True(t, v.inited)
		assert.Equal(t, "value", v.value)
		assert.Nil(t, r.Result().GetError())
	})

	t.Run("reads into TextUnmarshaler and gets error", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME": "value"})
		var v mockTextUnmarshaler
		v.err = errors.New("sorry")
		assert.True(t, r.Read("NAME", &v))

		assert.Equal(t, ValidationError{Path: ValidationPath{"NAME"}, Err: v.err}, r.Result().GetError())
	})

	t.Run("reads into simple types", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{
			"BOOL":      "true",
			"INT":       "1",
			"FLOAT":     "1.5",
			"STR":       "x",
			"BAD_BOOL":  "x",
			"BAD_INT":   "1.5",
			"BAD_FLOAT": "x",
		})
		var b, bb bool
		var i, bi int
		var f, bf float64
		var s string
		assert.True(t, r.Read("BOOL", &b))
		assert.True(t, r.Read("INT", &i))
		assert.True(t, r.Read("FLOAT", &f))
		assert.True(t, r.Read("STR", &s))
		assert.True(t, r.Read("BAD_BOOL", &bb))
		assert.True(t, r.Read("BAD_INT", &bi))
		assert.True(t, r.Read("BAD_FLOAT", &bf))
		assert.Equal(t, true, b)
		assert.Equal(t, 1, i)
		assert.Equal(t, float64(1.5), f)
		assert.Equal(t, "x", s)
		assert.Equal(t,
			[]ValidationError{
				{Path: ValidationPath{"BAD_BOOL"}, Err: errBoolFormat()},
				{Path: ValidationPath{"BAD_INT"}, Err: errIntFormat()},
				{Path: ValidationPath{"BAD_FLOAT"}, Err: errFloatFormat()},
			},
			r.Result().Errors(),
		)
	})

	t.Run("does not read into unknown types", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME": "value"})
		var f float32
		var s struct{}
		assert.False(t, r.Read("NAME", &f))
		assert.False(t, r.Read("NAME", &s))
		assert.Equal(t, []ValidationError{
			{ValidationPath{"NAME"}, varReaderBadTargetTypeError(&f)},
			{ValidationPath{"NAME"}, varReaderBadTargetTypeError(&s)},
		}, r.Result().Errors())
	})

	t.Run("Read does nothing for a missing optional value", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME": "value"})
		var v OptString
		assert.False(t, r.Read("UNKNOWN", &v))

		assert.Nil(t, r.Result().GetError())
	})

	t.Run("ReadRequired records an error for a missing value", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME": "value"})
		var v1, v2 OptString
		assert.True(t, r.ReadRequired("NAME", &v1))
		assert.False(t, r.ReadRequired("UNKNOWN", &v2))

		assert.Equal(t, []ValidationError{
			{ValidationPath{"UNKNOWN"}, errRequired()},
		}, r.Result().Errors())
	})

	t.Run("errors are accumulated", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME1": "value1", "NAME2": "value2"})
		var v1, v2 mockTextUnmarshaler
		v1.err = errors.New("err1")
		v2.err = errors.New("err2")
		assert.True(t, r.Read("NAME1", &v1))
		assert.True(t, r.Read("NAME2", &v2))

		assert.Equal(t, []ValidationError{{ValidationPath{"NAME1"}, v1.err}, {ValidationPath{"NAME2"}, v2.err}},
			r.Result().Errors())
	})

	t.Run("ReadStruct", func(t *testing.T) {
		t.Run("sets struct fields with conf tag", func(t *testing.T) {
			s := testStructWithTags1{FieldWithNoTag: "shouldn't change", F4: "shouldn't change either"}
			r := NewVarReaderFromValues(map[string]string{"STRING_VAR": "s", "DURATION_VAR": "1s", "BAD_INT_VAR": "x"})
			r.ReadStruct(&s, false)

			assert.Equal(t,
				testStructWithTags1{FieldWithNoTag: s.FieldWithNoTag, F1: "s", F2: NewOptDuration(time.Second), F4: s.F4},
				s)

			result := r.Result()
			assert.Equal(t, []ValidationError{{ValidationPath{"BAD_INT_VAR"}, errIntFormat()}}, result.Errors())
		})

		t.Run("enforces requiredness for fields with conf tag", func(t *testing.T) {
			s := testStructWithTags2{F1: "before1", F2: "before2", F3: "before3"}
			r := NewVarReaderFromValues(map[string]string{"STRING_VAR": "s"})
			r.ReadStruct(&s, false)

			assert.Equal(t,
				testStructWithTags2{F1: "s", F2: s.F2, F3: s.F3},
				s)

			result := r.Result()
			assert.Equal(t, []ValidationError{{ValidationPath{"NOT_SET_VAR1"}, errRequired()}}, result.Errors())
		})

		t.Run("logs error for invalid conf tag", func(t *testing.T) {
			s := testStructWithBadTag{}
			r := NewVarReaderFromValues(map[string]string{"STRING_VAR": "s"})
			r.ReadStruct(&s, false)
			assert.Error(t, r.Result().GetError())
		})

		t.Run("ignores nested structs when recursive is false", func(t *testing.T) {
			s := testStructWithNestedVars{
				F0: "oldF0",
				Nested: testStructWithTags1{
					F1: "oldF1",
				},
			}
			r := NewVarReaderFromValues(map[string]string{"TOP_LEVEL_VAR": "newF0", "STRING_VAR": "newF1"})
			r.ReadStruct(&s, false)
			assert.NoError(t, r.Result().GetError())
			assert.Equal(t, "newF0", s.F0)
			assert.Equal(t, "oldF1", s.Nested.F1)
		})

		t.Run("reads into nested structs when recursive is true", func(t *testing.T) {
			s := testStructWithNestedVars{
				F0: "oldF0",
				Nested: testStructWithTags1{
					F1: "oldF1",
				},
			}
			r := NewVarReaderFromValues(map[string]string{"TOP_LEVEL_VAR": "newF0", "STRING_VAR": "newF1"})
			r.ReadStruct(&s, true)
			assert.NoError(t, r.Result().GetError())
			assert.Equal(t, "newF0", s.F0)
			assert.Equal(t, "newF1", s.Nested.F1)
		})

		t.Run("rejects parameter that is not a struct pointer", func(t *testing.T) {
			var n int
			r1 := NewVarReaderFromValues(nil)
			r1.ReadStruct(&n, false)
			assert.Error(t, r1.Result().GetError())

			s := testStructWithTags1{}
			r2 := NewVarReaderFromValues(nil)
			r2.ReadStruct(s, false) // invalid because it's not a pointer
			assert.Error(t, r2.Result().GetError())

			var o OptBool // the Opt types are technically structs, but ReadStruct wouldn't make sense for them
			r3 := NewVarReaderFromValues(nil)
			r3.ReadStruct(&o, false)
			assert.Error(t, r3.Result().GetError())
		})
	})

	t.Run("WithVarNamePrefix", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"PRE_NAME": "value"})
		r1 := r.WithVarNamePrefix("PRE_")
		var s string
		assert.True(t, r1.Read("NAME", &s))
		assert.Equal(t, s, "value")
		assert.NoError(t, r.Result().GetError())

		var n int
		r1.Read("NAME", &n)
		assert.Equal(t, []ValidationError{
			{ValidationPath{"PRE_NAME"}, errIntFormat()},
		}, r.Result().Errors())
	})

	t.Run("WithVarNameSuffix", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"NAME_SUF": "value"})
		r1 := r.WithVarNameSuffix("_SUF")
		var s string
		assert.True(t, r1.Read("NAME", &s))
		assert.Equal(t, s, "value")

		var n int
		r1.Read("NAME", &n)
		assert.Equal(t, []ValidationError{
			{ValidationPath{"NAME_SUF"}, errIntFormat()},
		}, r.Result().Errors())
	})

	t.Run("FindPrefixedValues", func(t *testing.T) {
		r := NewVarReaderFromValues(map[string]string{"a": "1", "b_x": "2", "b_y": "3"})
		values := r.FindPrefixedValues("b_")
		assert.Equal(t, map[string]string{"x": "2", "y": "3"}, values)
	})
}

type testStructWithTags1 struct {
	FieldWithNoTag                 string
	unexportedFieldShouldBeIgnored string      `conf:"STRING_VAR"`
	F1                             string      `conf:"STRING_VAR"`
	F2                             OptDuration `conf:"DURATION_VAR"`
	F3                             int         `conf:"BAD_INT_VAR"`
	F4                             string      `conf:"NOT_SET_VAR"`
}

type testStructWithTags2 struct {
	F1 string `conf:"STRING_VAR"`
	F2 string `conf:"NOT_SET_VAR1,required"`
	F3 string `conf:"NOT_SET_VAR2"`
}

type testStructWithBadTag struct {
	F1 string `conf:"STRING_VAR,whatever"`
}

type testStructWithNestedVars struct {
	F0     string `conf:"TOP_LEVEL_VAR"`
	Nested testStructWithTags1
}

func withCleanEnvVars(action func()) {
	oldVars := os.Environ()
	os.Clearenv()
	defer func() {
		os.Clearenv()
		for _, v := range oldVars {
			name, value := parseVar(v)
			os.Setenv(name, value)
		}
	}()
	action()
}

type mockTextUnmarshaler struct {
	value  string
	inited bool
	err    error
}

func (m *mockTextUnmarshaler) UnmarshalText(data []byte) error {
	if m.err != nil {
		return m.err
	}
	m.value = string(data)
	m.inited = true
	return nil
}

type mockRequiredTextUnmarshaler struct {
	value  string
	inited bool
}

func (m *mockRequiredTextUnmarshaler) IsDefined() bool {
	return m.inited
}

func (m *mockRequiredTextUnmarshaler) IsRequired() bool {
	return true
}

func (m *mockRequiredTextUnmarshaler) UnmarshalText(data []byte) error {
	m.value = string(data)
	m.inited = true
	return nil
}
