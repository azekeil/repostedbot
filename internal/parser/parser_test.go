package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	assert.Equal(t, "", getCommand("1"))
	assert.Equal(t, "2", getCommand("1 2 "))
}

func TestGetArgs(t *testing.T) {
	assert.Equal(t, "", getArgs("1"))
	assert.Equal(t, "3 4", getArgs("1 2 3 4"))
}
