package app

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

// Response is custom response
type Response struct {
	Code    int         `json:"code"`
	Reason  string      `json:"reason"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageList struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}

// Json response
func Json(c *gin.Context, data interface{}, err error) {
	var resp Response
	if err != nil {
		_ = c.Error(err)
		se := errors.FromError(err)
		message := strings.Split(se.Message, ":")[0]
		message = strings.Split(message, "\n")[0]
		resp = Response{
			Code:    int(se.Code),
			Message: message,
			Reason:  se.Reason,
			Data:    data,
		}
	} else {
		resp = Response{
			Code:    200,
			Message: "success",
			Reason:  "success",
			Data:    data,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// DataFromReader writes the specified reader into the body stream and updates the HTTP code.
func DataFromReader(c *gin.Context, filename string, contentLength int64, reader io.Reader) {
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, filename),
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

// Data writes some data into the body stream and updates the HTTP code.
func Data(c *gin.Context, filename string, data []byte) {
	contentType := "application/octet-stream"
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, contentType, data)
}

// File writes the specified file into the body stream in an efficient way.
func File(c *gin.Context, filepath string) {
	c.File(filepath)
}

// String writes some string into the body stream and updates the HTTP code.
func String(c *gin.Context, data string, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.String(http.StatusOK, data)
}

// Bytes writes some data into the body stream and updates the HTTP code.
func Bytes(c *gin.Context, contentType string, data []byte) {
	c.Data(http.StatusOK, contentType, data)
}
