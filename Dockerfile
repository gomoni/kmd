FROM registry.opensuse.org/opensuse/bci/golang:latest as builder

# Install needed build-time for a builder image
RUN zypper --non-interactive install --no-recommends \
  tesseract-ocr-devel \
  tesseract-ocr-traineddata-ces \
  gcc-c++ \
  leptonica-devel

WORKDIR /src
#USER nobody

COPY . /src

RUN go mod download
RUN go build github.com/gomoni/kmd/cmd/kmd -o /src/kmd

FROM registry.opensuse.org/opensuse/bci/bci-minimal:latest

RUN zypper --non-interactive install --no-recommends \
    tesseract-ocr-traineddata-ces \
	libtesseract5 \
	libleptonica6

WORKDIR /app
COPY --from=builder /src/kmd /app

#USER nobody
ENTRYPOINT ["/app/kmd"]
