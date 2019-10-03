package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"

	"github.com/concepts-system/go-paperless/common"
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
	Context struct {
		echo.Context
		UserID   *uint
		Username *string
		Roles    []string
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
func (c Context) IsAuthenticated() bool {
	return c.UserID != nil
}

// HasRole returns a boolean value indicating whether the given context
// claims the given role.
func (c Context) HasRole(role string) bool {
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
func (c Context) BindPaging() common.PageRequest {
	pageRequest := common.PageRequest{}

	c.Bind(&pageRequest)

	if pageRequest.Offset < 0 {
		pageRequest.Offset = 0
	}

	if pageRequest.Size <= 0 {
		pageRequest.Size = common.DefaultPageSize
	} else if pageRequest.Size > common.MaxPageSize {
		pageRequest.Size = common.MaxPageSize
	}

	return pageRequest
}

// Page sends a page response.
func (c Context) Page(
	status int,
	page common.PageRequest,
	totalCount int64,
	data interface{},
) error {
	response := pageResponse{
		Size:       page.Size,
		Offset:     page.Offset,
		TotalCount: totalCount,
		Data:       data,
	}

	return c.JSON(status, response)
}

// BindAndValidate binds and validates the given object from the current context.
func (c Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return errors.BadRequest.Newf("Invalid request: %s", err.Error())
	}

	if err := c.Validate(i); err != nil {
		customError := errors.BadRequest.New("Validation failed")

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
func (c Context) BinaryAttachment(contentType, fileName string, contentLength int64, content io.ReadCloser) error {
	response := c.Response()
	headers := response.Header()

	response.Status = http.StatusOK
	headers.Set(HeaderContentType, contentType)
	headers.Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if contentLength >= 0 {
		headers.Set(HeaderContentLength, strconv.FormatInt(contentLength, 10))
	}

	io.Copy(response.Writer, content)
	return content.Close()
}

// CustomContext defines an echo middleware for using the custom context.
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		c := Context{Context: ec}
		return h(c)
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
