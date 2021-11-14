package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIllegalIdError_Error(t *testing.T) {
	assert.Equal(t, "Illegal ID", (&IllegalIdError{}).Error())
}

func TestInternalError_Error(t *testing.T) {
	assert.Equal(t, "Internal error", (&InternalError{}).Error())
}