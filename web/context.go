package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	// HeaderContentType is the name for the conten type HTTP header.
	HeaderContentType = "Content-Type"
	// HeaderContentDisposition is the name for the conten disposition HTTP header.
	HeaderContentDisposition = "Content-Disposition"
	// HeaderContentLength is the name for the conten disposition HTTP header.
	HeaderContentLength = "Content-Length"
)

// Context extends the standard echo context by API relevant fields.
type (
	context struct {
		echo.Context
		Username *string
		Roles    []string
		Scopes   []string
	}

	pageResponse struct {
		Size       int         `json:"size"`
		Offset     int         `json:"offset"`
		TotalCount int64       `json:"totalCount"`
		Data       interface{} `json:"data"`
	}
)

// IsAuthenticated returns a boolean value indicating whether
// the given context is authenticated.
func (c *context) IsAuthenticated() bool {
	return c.Username != nil
}

// HasRole returns a boolean value indicating whether the given context
// claims the given role.
func (c *context) HasRole(role string) bool {
	if !c.IsAuthenticated() {
		return false
	}

	for _, claimedRole := range c.Roles {
		if role == claimedRole {
			return true
		}
	}

	return false
}

// BindPaging tries to derive pagination info from the current request. The method falls back
// to default values (Offset: 0, Size: 10) if some arguments are missing or wrong.
func (c *context) BindPaging() pageRequest {
	pageRequest := pageRequest{}
	_ = c.Bind(&pageRequest)

	if pageRequest.Offset < 0 {
		pageRequest.Offset = 0
	}

	if pageRequest.Size <= 0 {
		pageRequest.Size = defaultPageSize
	} else if pageRequest.Size > maxPageSize {
		pageRequest.Size = maxPageSize
	}

	pageRequest.Sort = strings.TrimSpace(pageRequest.Sort)

	return pageRequest
}

// Page sends a page response.
func (c *context) Page(
	status int,
	page pageRequest,
	totalCount int64,
	data []interface{},
) error {
	response := pageResponse{
		Size:       len(data),
		Offset:     page.Offset,
		TotalCount: totalCount,
		Data:       data,
	}

	return c.JSON(status, response)
}

// BindAndValidate binds and validates the given object from the current context.
func (c *context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return application.BadRequestError.Newf("Invalid request: %s", err.Error())
	}

	if err := c.Validate(i); err != nil {
		customError := application.BadRequestError.New("Validation failed")

		for _, validationError := range err.(validator.ValidationErrors) {
			customError = errors.AddContext(
				customError,
				makeFirstLowerCase(validationError.Field()),
				validationError.Tag(),
			)
		}

		return customError
	}

	return nil
}

// BinaryAttachment sends a binary stream as an attachment with the given content type and disposition.
func (c *context) BinaryAttachment(
	contentType,
	fileName string,
	contentLength int64,
	content io.ReadCloser,
) error {
	response := c.Response()
	headers := response.Header()
	defer content.Close()

	response.Status = http.StatusOK
	headers.Set(HeaderContentType, contentType)
	headers.Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if contentLength >= 0 {
		headers.Set(HeaderContentLength, strconv.FormatInt(contentLength, 10))
	}

	_, err := io.Copy(response.Writer, content)
	return err
}

// extendedContext defines an echo middleware for using the extended context.
func extendedContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		c := context{Context: ec}
		return h(&c)
	}
}

func makeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}

	binary := []byte(s)

	start := bytes.ToLower([]byte{binary[0]})
	rest := binary[1:]

	return string(bytes.Join([][]byte{start, rest}, nil))
}
