package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	defaultRequestTTL = time.Second * 5

	TagHandler = "handler"

	MethodPost = http.MethodPost
)

type (
	Handler interface {
		Method() string
		Path() string
		Handle(ctx context.Context, req *gin.Context) Response
	}
)

func WrapHandler(parentCtx context.Context, f func(context.Context, *gin.Context) Response) gin.HandlerFunc {
	return func(req *gin.Context) {
		ctx, cancelFunc := context.WithTimeout(parentCtx, defaultRequestTTL)
		defer cancelFunc()

		resp := f(ctx, req)
		req.JSON(resp.StatusCode(), resp.Body())
	}
}
