package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EndpointNotFound struct {
}

type ResourceNotFound struct {
}

type MethodNotAllowed struct {
}

type MissingFields struct {
}

type InvalidRequest struct {
}

type InternalError struct {
}

func (e *EndpointNotFound) Error() string {
	return "Endpoint not found"
}

func (e *ResourceNotFound) Error() string {
	return "Resource not found"
}

func (e *MethodNotAllowed) Error() string {
	return "Method not allowed"
}

func (e *MissingFields) Error() string {
	return "Missing required fields"
}

func (e *InvalidRequest) Error() string {
	return "Invalid request"
}

func (e *InternalError) Error() string {
	return "Internal error"
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()

		// Check if there are any errors in the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Handle different types of errors
			switch err.Err.(type) {
			case *EndpointNotFound:
				HandleEndpointNotFound(c)
			case *ResourceNotFound:
				HandleResourceNotFound(c)
			case *MethodNotAllowed:
				HandleMethodNotAllowed(c)
			case *MissingFields:
				HandleMissingFields(c)
			case *InvalidRequest:
				HandleInvalidRequest(c)
			case *InternalError:
				HandleInternalError(c)
			}

			c.Abort()
			return
		}
	}
}

func HandleEndpointNotFound(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: (&EndpointNotFound{}).Error(),
	})
}

func HandleResourceNotFound(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: (&ResourceNotFound{}).Error(),
	})
}

func HandleMethodNotAllowed(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: (&MethodNotAllowed{}).Error(),
	})
}

func HandleMissingFields(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": (&MissingFields{}).Error(),
	})
}

func HandleInvalidRequest(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: (&InvalidRequest{}).Error(),
	})
}

func HandleInternalError(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: (&InternalError{}).Error(),
	})
}
