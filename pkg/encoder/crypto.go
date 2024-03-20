package encoder

import (
	"encoding/base64"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"

	"uniswap-transaction/pkg/utils/crypto"
)

func AesDecrypt(r *nethttp.Request, aesKey, encryptedStr string, v interface{}) error {
	encrypted, _ := base64.StdEncoding.DecodeString(encryptedStr)
	decrypted, err := crypto.AesDecrypt(aesKey, encrypted)
	if err != nil {
		return err
	}
	codec, _ := http.CodecForRequest(r, "Accept")
	err = codec.Unmarshal(decrypted, &v)
	if err != nil {
		return err
	}
	return nil
}

func AesEncrypt(r *nethttp.Request, aesKey string, v interface{}) (string, error) {
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return "", err
	}
	encrypted, err := crypto.AesEncrypt(aesKey, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
