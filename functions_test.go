package configtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrElse(t *testing.T) {
	assert.Equal(t, "true", StringOrElse(NewOptBool(true), "none"))
	assert.Equal(t, "none", StringOrElse(OptBool{}, "none"))
}
