package workwx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkWx(t *testing.T) {
	c := &Config{
		Key: "",
	}
	alert := NewWorkWx(c)
	err := alert.SendMessage("test")
	assert.NoError(t, err)
}
