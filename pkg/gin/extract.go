package gin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"

	"uniswap-transaction/pkg/gin/app"
	"uniswap-transaction/pkg/utils"
)

// ExtractArgs returns the string of the req
func ExtractArgs(ctx context.Context) (query, args string) {
	if ginCtx, ok := FromGinContext(ctx); ok {
		query = ginCtx.Request.URL.RawQuery
		if len(query) > 128 {
			query = query[:128]
		}
		obj := app.FromBindContext(ginCtx)
		bytes, _ := json.Marshal(obj)
		if len(bytes) > 1024 {
			bytes = bytes[:1024]
		}
		args = utils.BytesToString(bytes)
	}
	return query, args
}

// ExtractError returns the string of the error
func ExtractError(err error) (log.Level, string) {
	if err != nil {
		if ginErr, ok := err.(*gin.Error); ok {
			err = ginErr.Err
		}
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
