FROM frolvlad/alpine-glibc

# Install Dependencies
RUN apk --update upgrade
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    tesseract-ocr \
    tesseract-ocr-data-deu \
    imagemagick
RUN rm -rf /var/cache/apk/*

# Prepare Executable
ADD go-paperless /
ENTRYPOINT /go-paperless
