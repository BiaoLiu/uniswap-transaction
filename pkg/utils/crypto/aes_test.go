package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEncrypt(t *testing.T) {
	key := "asuperstrong16bi"
	d := []byte("{\"levelState\":\"test\", \"achievement\":\"test\", \"totalScore\":\"1000\", \"lastOpenedLevel\":1, \"isMute\":0}\n")

	data, err := AesEncrypt(key, d)
	assert.NoError(t, err)
	text := base64.StdEncoding.EncodeToString(data)
	fmt.Println(text)

	mwb, _ := base64.StdEncoding.DecodeString(text)
	val, err := AesDecrypt(key, mwb)
	assert.NoError(t, err)
	fmt.Println(string(val))
}
