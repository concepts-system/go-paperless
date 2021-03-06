package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/labstack/echo/v4"
)

const (
	pagesFormKey = "pages[]"
)

type documentRouter struct {
	documentService application.DocumentService
}

// NewDocumentRouter creates a new router for document management using the given
// document service.
func NewDocumentRouter(documentService application.DocumentService) Router {
	return &documentRouter{
		documentService: documentService,
	}
}

// DefineRoutes defines the routes for document management.
func (r *documentRouter) DefineRoutes(group *echo.Group, auth *AuthMiddleware) {
	apiGroup := group.Group("/api", auth.RequireScope(application.TokenScopeAPI))

	documentGroup := apiGroup.Group("/documents", auth.RequireAuthentication())
	documentGroup.GET("", r.getDocuments)
	// documentGroup.GET("/search", searchDocuments)
	documentGroup.POST("", r.createDocument)
	documentGroup.GET("/:id", r.getDocument)
	// documentGroup.PUT("/:id", updateDocument)
	// // documentGroup.DELETE("/:id", deleteDocument)
	// documentGroup.GET("/:id/content", getDocumentContent)

	pageGroup := documentGroup.Group("/:id/pages")
	pageGroup.GET("", r.getDocumentPages)
	pageGroup.POST("/content", r.addPagesToDocument)
	pageGroup.GET("/:pageNumber", r.getDocumentPage)
	// // pageGroup.PUT("/:pageNumber", updateDocumentPage)
	// // pageGroup.DELETE("/:pageNumber", deleteDocumentPage)
	pageGroup.GET("/:pageNumber/content", r.getDocumentPageContent)
	// // pageGroup.PUT("/:pageNumber/content", updatePageContent)
}

/* Handlers */

func (r *documentRouter) getDocuments(ec echo.Context) error {
	c, _ := ec.(*context)
	pr := c.BindPaging()

	documents, totalCount, err := r.documentService.GetUserDocuments(
		*c.Username,
		pr.ToDomainPageRequest(),
	)

	if err != nil {
		return err
	}

	serializer := documentListSerializer{c, documents}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func (r *documentRouter) createDocument(ec echo.Context) error {
	c, _ := ec.(*context)
	validator := newDocumentValidator()

	if err := validator.Bind(c); err != nil {
		return err
	}

	document, err := r.documentService.CreateNewDocument(*c.Username, &validator.document)
	if err != nil {
		return err
	}

	serializer := documentSerializer{c, document}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func (r *documentRouter) getDocument(ec echo.Context) error {
	c, _ := ec.(*context)
	documentNumber, err := r.bindDocumentNumber(ec)
	if err != nil {
		return err
	}

	document, err := r.documentService.GetUserDocumentByDocumentNumber(*c.Username, documentNumber)
	if err != nil {
		return err
	}

	serializer := documentSerializer{c, document}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *documentRouter) getDocumentPages(ec echo.Context) error {
	c, _ := ec.(*context)
	documentNumber, err := r.bindDocumentNumber(c)
	if err != nil {
		return err
	}

	pr := c.BindPaging()
	pages, totalCount, err := r.documentService.GetUserDocumentPagesByDocumentNumber(*c.Username, documentNumber, pr.ToDomainPageRequest())
	if err != nil {
		return err
	}

	serializer := documentPageListSerializer{c, pages}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func (r *documentRouter) getDocumentPage(ec echo.Context) error {
	c, _ := ec.(*context)
	documentNumber, err := r.bindDocumentNumber(c)
	if err != nil {
		return err
	}

	pageNumber, err := r.bindPageNumber(c)
	if err != nil {
		return err
	}

	page, err := r.documentService.GetUserDocumentPageByDocumentNumberAndPageNumber(*c.Username, documentNumber, pageNumber)
	if err != nil {
		return err
	}

	serializer := documentPageSerializer{c, true, page}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *documentRouter) addPagesToDocument(ec echo.Context) error {
	c, _ := ec.(*context)
	documentNumber, err := r.bindDocumentNumber(c)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return application.BadRequestError.Newf("Expecting valid multipart form with '%s' containing at least one file", pagesFormKey)
	}

	files := form.File[pagesFormKey]
	pages, err := r.documentService.AddPagesToUserDocument(*c.Username, documentNumber, files)
	if err != nil {
		return err
	}

	serializer := documentPageListSerializer{c, pages}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func (r *documentRouter) getDocumentPageContent(ec echo.Context) error {
	c, _ := ec.(*context)
	documentNumber, err := r.bindDocumentNumber(c)
	if err != nil {
		return err
	}

	pageNumber, err := r.bindPageNumber(c)
	if err != nil {
		return err
	}

	page, err := r.documentService.GetUserDocumentPageByDocumentNumberAndPageNumber(*c.Username, documentNumber, pageNumber)
	if err != nil {
		return err
	}

	content, err := r.documentService.GetUserDocumentPageContent(*c.Username, documentNumber, pageNumber)
	if err != nil {
		return err
	}

	extension, mimeType := r.getPageContentFileInfos(page)
	title := string(page.Document.Title)
	if strings.TrimSpace(title) == "" {
		title = fmt.Sprint(documentNumber)
	}

	return c.BinaryAttachment(
		mimeType,
		fmt.Sprintf("%s - %d.%s", title, pageNumber, extension),
		-1,
		content,
	)
}

/* Helper Methods */

func (r *documentRouter) bindDocumentNumber(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id <= 0 {
		return 0, application.BadRequestError.New("Document ID has to be a positive integer")
	}

	return uint(id), nil
}

func (r *documentRouter) bindPageNumber(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("pageNumber"), 10, 32)

	if err != nil || id <= 0 {
		return 0, application.BadRequestError.New("Page number has to be a non-negative integer")
	}

	return uint(id), nil
}

func (r *documentRouter) getPageContentFileInfos(page *domain.DocumentPage) (string, string) {
	switch page.Type {
	case domain.PageTypeTIFF:
		return ".tiff", "image/tiff"
	default:
		return "bin", "application/octet-stream"
	}
}
