# Go Paperless

As the name suggests, _Go Paperless_ is an application for enabling to manage classic paper documents digitally. It provides functionality to index and analyze scanned documents and provide searchable PDFs.

The main goal is to help people with organizing their paperwork.

## Functionalities

1. User management
2. Document management
   - Upload of scans
   - Indexing of scanned documents
   - Text recognition of scanned documents
   - Creation of searchable PDFs based on scans

## Dependencies

Next to the implicit _Go_ dependecies itself, _Go Paperless_ relies on some third-party software:

1. _Go Paperless_ will require some common database system for storing its data. Right now, [SQLite](https://www.sqlite.org/index.html) and [Postgres](https://www.postgresql.org/) are supported.
2. Recognizing, analyzing and the creation of searchable PDFs are done by [Tesseract OCR](https://github.com/tesseract-ocr/).
3. As above tasks need some time for processing, they are done asynchronously. On various user actions, _Go Paperless_ will send async jobs to the [Faktory](https://github.com/contribsys/faktory) job processor. In a second step it will fetch jobs from there and do the expensive work in background.
4. For full-text search, documents are indexed using [Bleve](https://github.com/blevesearch/bleve).

## Configuration

_Go Paperless_ is fully configurable through environment variables. See [config.go](common/config.go) for all configuration options.

## Running locally

Executing _Go Paperless_ locally is easy if you have docker installed. Run the following series of commands in order to start the application:

```sh
$ make build-docker && docker-compose up
```

Use the provided [Postman](https://www.getpostman.com/) collection and start try out the software. Visit [localhost:7420](http://localhost:7420) for the _Faktory_ job dashboard.
