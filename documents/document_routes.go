package documents

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kpango/glg"
	"github.com/labstack/echo"

	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/auth"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	formKeyPages             = "pages[]"
	mimeHeaderKeyContentType = "Content-Type"
	contentTypePdf           = "application/pdf"

	queryParamSort      = "sort"
	queryParamHighlight = "highlight"
)

var (
	validContentTypes   = regexp.MustCompile("^image/[a-zA-Z\\-]+$")
	validHighlightTypes = regexp.MustCompile("^html$")
)

// RegisterRoutes registers all related routes for managing users.
func RegisterRoutes(r *echo.Group) {
	documentGroup := r.Group("/documents", auth.RequireAuthorization())
	documentGroup.GET("", getDocuments)
	documentGroup.GET("/search", searchDocuments)
	documentGroup.POST("", createDocument)
	documentGroup.GET("/:id", getDocument)
	documentGroup.PUT("/:id", updateDocument)
	// documentGroup.DELETE("/:id", deleteDocument)
	documentGroup.GET("/:id/raw", getDocumentContent)

	pageGroup := documentGroup.Group("/:id/pages")
	pageGroup.GET("", getDocumentPages)
	pageGroup.POST("/raw", addPagesToDocument)
	pageGroup.GET("/:pageNumber", getDocumentPage)
	// pageGroup.PUT("/:pageNumber", updateDocumentPage)
	// pageGroup.DELETE("/:pageNumber", deleteDocumentPage)
	// pageGroup.GET("/:pageNumber/raw", getPageContent)
	// pageGroup.PUT("/:pageNumber/raw", updatePageContent)
}

// Document Handlers

func getDocuments(ec echo.Context) error {
	c, _ := ec.(api.Context)
	pr := c.BindPaging()

	documents, totalCount, err := FindDocumentsByUserID(*c.UserID, pr)
	if err != nil {
		return err
	}

	serializer := DocumentListSerializer{c, documents}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func searchDocuments(ec echo.Context) error {
	c, _ := ec.(api.Context)
	pr := c.BindPaging()
	sort := c.QueryParam(queryParamSort)
	highlight := c.QueryParam(queryParamHighlight)

	if len(highlight) > 0 && !isHighlightTypeSupported(highlight) {
		return errors.BadRequest.Newf(
			"The highlight type '%s' is not supported",
			highlight,
		)
	}

	results, totalCount, err := SearchDocuments(*c.UserID, c.QueryParam("query"), pr, sort, highlight)

	if err != nil {
		return err
	}

	return c.Page(http.StatusOK, pr, int64(totalCount), results)
}

func createDocument(ec echo.Context) error {
	c, _ := ec.(api.Context)
	validator := NewDocumentModelValidator()

	if err := validator.Bind(c); err != nil {
		return err
	}

	validator.documentModel.OwnerID = *c.UserID
	validator.documentModel.State = DocumentStateClean

	if err := validator.documentModel.Create(); err != nil {
		return err
	}

	serializer := DocumentSerializer{c, &validator.documentModel}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func getDocument(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	serializer := DocumentSerializer{c, document}
	return c.JSON(http.StatusOK, serializer.Response())
}

func updateDocument(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	validator := NewDocumentModelValidatorFillWith(*document)

	if err := validator.Bind(c); err != nil {
		return err
	}

	if err := validator.documentModel.Save(); err != nil {
		return err
	}

	serializer := DocumentSerializer{c, &validator.documentModel}
	return c.JSON(http.StatusOK, serializer.Response())
}

func getDocumentContent(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	if document.State != DocumentStateClean {
		return errors.BadRequest.New("Document pipeline not complete; try again when document's state is 'CLEAN'")
	}

	pages, err := GetAllPagesByDocumentID(document.ID)
	if err != nil {
		return err
	}

	if len(pages) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	contentFile, err := OpenContent(document.ID, document.ContentID)
	if err != nil {
		return err
	}

	fileInfo, err := contentFile.Stat()
	if err != nil {
		return err
	}

	defer contentFile.Close()

	title := document.Title
	if title == "" {
		title = "Document"
	}

	return c.BinaryAttachment(
		contentTypePdf,
		fmt.Sprintf("%s.pdf", title),
		fileInfo.Size(),
		contentFile,
	)
}

// Page Handlers

func addPagesToDocument(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return errors.BadRequest.New("Expecting valid multipart form with 'pages[]' containing at least one file")
	}

	files := form.File[formKeyPages]
	if len(files) <= 0 {
		return errors.BadRequest.New("Expecting at least on page")
	}

	pages := make([]PageModel, len(files))
	glg.Infof("Appending pages to document %d...", document.ID)

	for idx, page := range files {
		stream, err := page.Open()
		if err != nil {
			return errors.BadRequest.Newf("Failed to read page %d during upload: %s", idx, err.Error())
		}

		contentType := page.Header.Get(mimeHeaderKeyContentType)
		if !isPageContentTypeSupported(contentType) {
			return errors.BadRequest.Newf(
				"The content type '%s' for page %d is not supported; please supply a valid image content type",
				contentType,
				idx,
			)
		}

		glg.Debugf("Appending page %d with type '%s'", idx, contentType)
		page, err := AppendPageToDocument(document, contentType, stream)

		if err != nil {
			return errors.Wrapf(err, "Failed to index page %d", idx)
		}

		pages[idx] = *page
	}

	serializer := PageListSerializer{c, pages}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func getDocumentPages(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	pr := c.BindPaging()
	pages, totalCount, err := FindPagesByDocumentID(document.ID, pr)
	if err != nil {
		return err
	}

	serializer := PageListSerializer{c, pages}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func getDocumentPage(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindDocumentID(ec)
	if err != nil {
		return err
	}

	pageNumber, err := bindPageNumber(c)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(id)
	if err != nil {
		return err
	}

	page, err := GetPageByDocumentIDAndPageNumber(id, pageNumber)
	if err != nil {
		return err
	}

	if !allowedToAccessDocument(c, document) {
		return errors.Forbidden.New("Access to given document denied")
	}

	serializer := PageSerializer{c, page}
	return c.JSON(http.StatusOK, serializer.Response())
}

// Helper Methods

func bindDocumentID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id <= 0 {
		return 0, errors.BadRequest.New("Document ID has to be a positive integer")
	}

	return uint(id), nil
}

func bindPageNumber(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("pageNumber"), 10, 32)

	if err != nil || id < 0 {
		return 0, errors.BadRequest.New("Page number has to be a non-negative integer")
	}

	return uint(id), nil
}

func allowedToAccessDocument(c api.Context, document *DocumentModel) bool {
	return document.OwnerID == *c.UserID
}

func isPageContentTypeSupported(contentType string) bool {
	return validContentTypes.MatchString(contentType)
}

func isHighlightTypeSupported(highlightType string) bool {
	return validHighlightTypes.MatchString(highlightType)
}
