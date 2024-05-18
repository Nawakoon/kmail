package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertNoAnyError(t *testing.T, errors ...error) {
	for _, err := range errors {
		assert.NoError(t, err)
	}
}
