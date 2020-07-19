package web

import (
	"net/http"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type errorResponse struct {
	StatusCode int         `json:"status"`
	Title      string      `json:"title"`
	Cause      string      `json:"cause,omitempty"`
	Details    interface{} `json:"details,omitempty"`
}

func newErrorResponse(stausCode int, err error) errorResponse {
	resp := errorResponse{
		StatusCode: stausCode,
		Title:      err.Error(),
	}

	application.RemoveErrorType(err)

	if context := errors.GetContext(err); context != nil {
		resp.Details = context
	}

	if cause := errors.RootCause(err); cause != err && cause.Error() != err.Error() {
		resp.Cause = cause.Error()
	}

	return resp
}

func errorHandler(err error, c echo.Context) {
	var response errorResponse

	switch err.(type) {
	case *domain.Error:
		response = handleDomainError(err.(*domain.Error))
		break
	case *echo.HTTPError:
		response = handleHTTPError(err.(*echo.HTTPError))
		log.Error(err)
		break
	default:
		response = handleErrorGeneric(err)
		log.Error(err)
	}

	c.JSON(response.StatusCode, response)
}

func handleDomainError(err *domain.Error) errorResponse {
	return newErrorResponse(http.StatusBadRequest, err)
}

func handleHTTPError(err *echo.HTTPError) errorResponse {
	message := err.Error()
	if messageStr, ok := err.Message.(string); ok {
		message = messageStr
	}

	return errorResponse{StatusCode: err.Code, Title: message}
}

func handleErrorGeneric(err error) errorResponse {
	return newErrorResponse(getStatusCode(err), err)
}

func getStatusCode(err error) int {
	switch application.GetErrorType(err) {
	case application.BadRequestError:
		return http.StatusBadRequest
	case application.UnauthorizedError:
		return http.StatusUnauthorized
	case application.ForbiddenError:
		return http.StatusForbidden
	case application.NotFoundError:
		return http.StatusNotFound
	case application.ConflictError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
