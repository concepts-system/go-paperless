package documents

import (
	"github.com/concepts-system/go-paperless/errors"
	"github.com/concepts-system/go-paperless/jobs"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/kpango/glg"
)

const (
	jobConvertPage   = "page.convert"
	jobRecognizePage = "page.recognize"

	jobIndexDocument    = "document.index"
	jobGenerateDocument = "document.generate"
)

var (
	imageConverter = ImageConverter{}
	ocrEngine      = OcrEngine{}
)

// RegisterWorkers registers all workers related to asynchronous job processing of documents and pages.
func RegisterWorkers(manager *worker.Manager) {
	manager.Register(jobConvertPage, jobHandler(convertPage))
	manager.Register(jobRecognizePage, jobHandler(recognizePage))
	manager.Register(jobIndexDocument, jobHandler(indexDocument))
	manager.Register(jobGenerateDocument, jobHandler(generateDocument))
}

// Job Submission Functions

func submitPageConversionJob(pageID uint) error {
	job := faktory.NewJob(jobConvertPage, pageID)
	return jobs.Client().Push(job)
}

func submitPageRecognitionJob(pageID uint) error {
	job := faktory.NewJob(jobRecognizePage, pageID)
	return jobs.Client().Push(job)
}

func submitDocumentIndexingJob(documentID uint) error {
	job := faktory.NewJob(jobIndexDocument, documentID)
	return jobs.Client().Push(job)
}

func submitDocumentGenerationJob(documentID uint) error {
	job := faktory.NewJob(jobGenerateDocument, documentID)
	return jobs.Client().Push(job)
}

// Job Handler Functions

func convertPage(ctx worker.Context, args ...interface{}) error {
	glg.Infof("Received %s job with ID %s", ctx.JobType(), ctx.Jid())

	pageID, err := getIDFromArgs(args)
	if err != nil {
		return err
	}

	glg.Infof("Converting page %d", pageID)

	if err := imageConverter.ConvertPage(pageID); err != nil {
		return errors.Wrapf(err, "Conversion of page %d failed", pageID)
	}

	if err := submitPageRecognitionJob(pageID); err != nil {
		return errors.Wrapf(err, "Submission of page recognition job for page %d failed", pageID)
	}

	glg.Infof("Sent recognition job for page %d", pageID)
	glg.Successf("Conversion for page %d complete", pageID)
	return nil
}

func recognizePage(ctx worker.Context, args ...interface{}) error {
	glg.Infof("Received %s job with ID %s", ctx.JobType(), ctx.Jid())

	pageID, err := getIDFromArgs(args)
	if err != nil {
		return err
	}

	page, err := FindPageByID(pageID)
	if err != nil {
		return err
	}

	glg.Infof("Running recognition for page %d", pageID)

	text, err := ocrEngine.RecognizePage(page)
	if err != nil {
		return errors.Wrapf(err, "Recognition process for page %d failed", pageID)
	}

	page.Text = text.String()
	page.State = PageStateClean

	if err = page.Save(); err != nil {
		return errors.Wrapf(err, "Saving of page %d failed", pageID)
	}

	pages, err := GetAllPagesByDocumentID(page.DocumentID)
	if err != nil {
		return errors.Wrapf(err, "Retreival of pages for document %d failed", page.DocumentID)
	}

	if allPagesClean(pages) {
		glg.Infof("All pages for document %d clean; sending indexing job...", page.DocumentID)

		document, err := GetDocumentByID(page.DocumentID)
		if err != nil {
			return errors.Wrapf(err, "Failed to find document with ID %d", page.DocumentID)
		}

		document.State = DocumentStateDirty
		if err := document.Save(); err != nil {
			return err
		}

		if err := submitDocumentIndexingJob(page.DocumentID); err != nil {
			return errors.Wrapf(
				err,
				"Submission of document indexing job for document %d failed",
				page.DocumentID,
			)
		}

		glg.Infof("Submitted indexing job for document %d", page.DocumentID)
	}

	glg.Successf("Recognition for page %d complete", pageID)
	return nil
}

func indexDocument(ctx worker.Context, args ...interface{}) error {
	glg.Infof("Received %s job with ID %s", ctx.JobType(), ctx.Jid())

	documentID, err := getIDFromArgs(args)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(documentID)
	if err != nil {
		return err
	}

	if err = IndexDocument(documentID, GetIndex()); err != nil {
		return errors.Wrapf(err, "Indexing job for document %d failed", documentID)
	}

	document.State = DocumentStateIndexed
	if err = document.Save(); err != nil {
		return err
	}

	if err := submitDocumentGenerationJob(documentID); err != nil {
		return err
	}

	glg.Successf("Indexing for document %d complete", documentID)
	return nil
}

func generateDocument(ctx worker.Context, args ...interface{}) error {
	glg.Infof("Received %s job with ID %s", ctx.JobType(), ctx.Jid())

	documentID, err := getIDFromArgs(args)
	if err != nil {
		return err
	}

	document, err := GetDocumentByID(documentID)
	if err != nil {
		return errors.Wrapf(err, "Failed to find document with ID %d", documentID)
	}

	glg.Infof("Generating searchable PDF for document %d", documentID)

	newContentID, err := ocrEngine.GenerateDocument(document)
	if err != nil {
		return err
	}

	if err := DeleteContent(documentID, document.ContentID); err != nil {
		return errors.Wrapf(
			err,
			"Failed to clean-up old content %s for document %d",
			document.ContentID,
			documentID,
		)
	}

	document.ContentID = newContentID
	document.State = DocumentStateClean

	if err := document.Save(); err != nil {
		return errors.Wrapf(err, "Failed to update document %d", documentID)
	}

	glg.Successf("Generation for document %d complete", documentID)
	return nil
}

// Helper Functions

func jobHandler(handler worker.Perform) worker.Perform {
	return func(ctx worker.Context, args ...interface{}) error {
		err := handler(ctx, args...)
		if err != nil {
			glg.Error(err)
		}

		return err
	}
}

func getIDFromArgs(args []interface{}) (uint, error) {
	if len(args) != 1 {
		return 0, errors.BadRequest.New("Invalid arguments given; expecting the ID as single parameter")
	}

	idRaw, ok := args[0].(float64)
	if !ok || idRaw < 0 {
		return 0, errors.BadRequest.New("Invalid ID; expecing a valid, positive number")
	}

	return uint(idRaw), nil
}

func allPagesClean(pages []PageModel) bool {
	for _, page := range pages {
		if page.State != PageStateClean {
			return false
		}
	}

	return true
}
