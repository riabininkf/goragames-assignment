package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var InternalServerError = NewErrorResponse(http.StatusInternalServerError, "internal server error")

type (
	Response interface {
		StatusCode() int
		Body() any
	}

	response struct {
		statusCode int
		body       any
	}
)

func (r response) StatusCode() int {
	return r.statusCode
}

func (r response) Body() any {
	return r.body
}

func NewResponse(statusCode int, body any) Response {
	return response{statusCode: statusCode, body: body}
}

func NewEmptyResponse() Response {
	return NewResponse(http.StatusOK, gin.H{})
}

func NewErrorResponse(statusCode int, msg string) Response {
	return NewResponse(statusCode, gin.H{"error": msg})
}

func NewBadRequest(msg string) Response {
	return NewErrorResponse(http.StatusBadRequest, msg)
}
func NewBadNotFound(msg string) Response {
	return NewErrorResponse(http.StatusNotFound, msg)
}
