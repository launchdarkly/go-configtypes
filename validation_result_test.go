package configtypes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationResult(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		assert.True(t, ValidationResult{}.OK())
		assert.Equal(t, []ValidationError{}, ValidationResult{}.Errors())
	})

	t.Run("AddError", func(t *testing.T) {
		err1, err2 := errors.New("err1"), errors.New("err2")

		var r ValidationResult
		r.AddError(nil, err1)
		r.AddError(ValidationPath{"x"}, err2)

		assert.False(t, r.OK())
		assert.Equal(t, []ValidationError{{nil, err1}, {ValidationPath{"x"}, err2}}, r.Errors())
	})

	t.Run("AddAll", func(t *testing.T) {
		err1, err2 := errors.New("err1"), errors.New("err2")

		var r ValidationResult
		r.AddError(nil, err1)

		var sub ValidationResult
		sub.AddError(ValidationPath{"b"}, err2)

		r.AddAll(ValidationPath{"a"}, sub)

		assert.Equal(t, []ValidationError{{nil, err1}, {ValidationPath{"a", "b"}, err2}}, r.Errors())
	})

	t.Run("aggregate Error", func(t *testing.T) {
		err1, err2 := errors.New("err1"), errors.New("err2")

		var r1 ValidationResult
		assert.Nil(t, r1.GetError())

		var r2 ValidationResult
		r2.AddError(nil, err1)
		assert.Equal(t, ValidationError{Err: err1}, r2.GetError())

		var r3 ValidationResult
		r3.AddError(ValidationPath{"a"}, err1)
		assert.Equal(t, ValidationError{Path: ValidationPath{"a"}, Err: err1}, r3.GetError())

		var r4 ValidationResult
		r4.AddError(nil, err1)
		r4.AddError(ValidationPath{"a"}, err2)
		assert.Equal(t,
			ValidationAggregateError{{Err: err1}, {Path: ValidationPath{"a"}, Err: err2}},
			r4.GetError(),
		)
	})
}

func TestValidationError(t *testing.T) {
	e1 := ValidationError{Path: nil, Err: errors.New("message")}
	assert.Equal(t, "message", e1.String())

	e2 := ValidationError{Path: ValidationPath{"a"}, Err: errors.New("message")}
	assert.Equal(t, "a: message", e2.String())

	e3 := ValidationError{Path: ValidationPath{"a", "b"}, Err: errors.New("message")}
	assert.Equal(t, "a.b: message", e3.String())
}

func TestValidationAggregateError(t *testing.T) {
	e1 := ValidationAggregateError{}
	assert.Equal(t, "", e1.String())

	e2 := ValidationAggregateError{{Err: errors.New("message")}}
	assert.Equal(t, "message", e2.String())

	e3 := ValidationAggregateError{{Err: errors.New("message1")}, {Path: ValidationPath{"a"}, Err: errors.New("message2")}}
	assert.Equal(t, "message1, a: message2", e3.String())
}

func TestValidationPath(t *testing.T) {
	assert.Equal(t, "", ValidationPath(nil).String())
	assert.Equal(t, "", ValidationPath{}.String())
	assert.Equal(t, "a", ValidationPath{"a"}.String())
	assert.Equal(t, "a.b.c", ValidationPath{"a", "b", "c"}.String())
}
