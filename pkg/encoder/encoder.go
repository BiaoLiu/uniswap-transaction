package encoder

import (
	"encoding/json"
	nethttp "net/http"
	"strings"

	"github.com/BiaoLiu/go-i18n"
	kjson "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"

	"uniswap-transaction/pkg/net/httputil"
)

func init() {
	kjson.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   false,
	}
}

// Response is custom response
type Response struct {
	Code    int         `json:"code"`
	Reason  string      `json:"reason"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseEncoder is encode response func.
func ResponseEncoder(w nethttp.ResponseWriter, r *nethttp.Request, v interface{}) error {
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	reply := Response{
		Code:    200,
		Message: "success",
		Reason:  "success",
	}
	err = json.Unmarshal(data, &reply.Data)
	if err != nil {
		return err
	}
	data, err = codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	w.WriteHeader(nethttp.StatusOK)
	_, _ = w.Write(data)
	return nil
}

// ErrorEncoder is encode error func.
func ErrorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	var message string
	se := errors.FromError(err)
	reply := Response{
		Code:    int(se.Code),
		Data:    nil,
		Message: se.Message,
		Reason:  se.Reason,
	}
	if errors.IsInternalServer(err) {
		message = "System Error"
	} else {
		accept := r.Header.Get("accept-language")
		lang, err := language.Parse(accept)
		if err != nil {
			lang = language.English
		}
		tag, _, _ := i18n.Supported.Match(lang)
		switch tag {
		case language.Chinese, language.SimplifiedChinese, language.TraditionalChinese:
			message = strings.Split(se.Message, " ")[0]
		default:
			message = strings.Split(se.Message, ".")[0]
		}
		message = strings.Split(se.Message, ":")[0]
	}
	reply.Message = message
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(nethttp.StatusOK)
		return
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	w.WriteHeader(nethttp.StatusOK)
	_, _ = w.Write(data)
}
