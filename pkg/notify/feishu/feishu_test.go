package feishu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeiShu(t *testing.T) {
	c := &Config{
		Key: "502eaa32-41f1-45d3-8dc1-2b77202c2a15",
	}
	alert := NewFeiShu(c)
	content := "test\n\ntest2"
	err := alert.SendMessage(content)
	assert.NoError(t, err)
}
