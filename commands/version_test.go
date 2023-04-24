package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runVersion(t *testing.T) {
	err := runVersion()
	assert.Equal(t, err, nil)
}
