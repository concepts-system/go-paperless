package api

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/concepts-system/go-paperless/errors"
)

type errorResponse struct {
	Message string      `json:"message"`
	Cause   string      `json:"cause,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorHandler defines an error handling function for API errors.
func ErrorHandler(err error, c echo.Context) {
	var response errorResponse
	var statusCode int

	if httpError, ok := err.(*echo.HTTPError); ok {
		message := httpError.Error()
		if messageStr, ok := httpError.Message.(string); ok {
			message = messageStr
		}

		response = errorResponse{Message: message}
		statusCode = httpError.Code
	} else {
		statusCode = getStatusCode(err)
		response = errorResponse{Message: err.Error()}

		if context := errors.GetContext(err); context != nil {
			response.Details = context
		}

		if cause := errors.Cause(err).Error(); cause != response.Message {
			response.Cause = errors.Cause(err).Error()
		}
	}

	c.JSON(statusCode, response)
}

func getStatusCode(err error) int {
	switch errors.GetType(err) {
	case errors.BadRequest:
		return http.StatusBadRequest
	case errors.Unauthorized:
		return http.StatusUnauthorized
	case errors.Forbidden:
		return http.StatusForbidden
	case errors.NotFound:
		return http.StatusNotFound
	case errors.Conflict:
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}
