# Kmd

> Karol Michal Daemon

HTTP/CLI interface for `tesseract-ocr` written in Go.

## Architecture

Client/server - server has a Tesseracr-OCR and everything else installed, while
client is a simple static binary which talks to server. Curl would be enough
too.

 * `kmd` - waits on systemd socket activation fd or opens a unix socket at
   `/run/user/1000/kmd.sock` and listens there
 * `kmc` - CLI part of a `kmd`

## TODO

 1. finish the refactoring of stuff to internal - mainly params
 2. move cmd/render to kmd/kmc
 3. more tests

 * test https://github.com/klippa-app/go-pdfium
    especially https://github.com/klippa-app/go-pdfium?tab=readme-ov-file#webassembly
 * convert pdf to png
  + poppler-utils: pdftoppm input.pdf outputname -png
  + image-magic: convert -density 150 input.pdf[666] -quality 90 output.png
  + ghostscript?
 * tests
 * openSUSE package for mage
 * Dockerfile + public docker imaaazzz + GHA updating the shit
 * install a systemctl file(s)
 * do not hardcode unix path + make it configurable
 * HTTP/Accept for server - implement text/plain and application/json at least
 * errors reporting - maybe terrasect can't report errors other way than printing it?

```
Error in pixReadStream: Pdf reading is not supported
Leptonica Error in pixRead: pix not read: /tmp/orcserver1026615170
```

# Usage

```sh
# build
mage build
# or manually
go build github.com/gomoni/kmd/cmd/kmd
go build github.com/gomoni/kmd/cmd/kmc

# run server via systemd socket activation
# will be handier once Docker files will be ready
systemd-socket-activate -l /run/user/1000/kmd.sock ./kmd

# or directly
./kmdd

# run client
./kmc
version: 5.4.0
languages:
 * ces
 * eng
# run ocr
./kmc ocr internal/ocr/testdata/hello.png
```

# Why OCR?

> I've said this before and I'll say it again: PDF sucks!

Not there're tons of Go libraries for PDF anyway. And I am not crazy enough to
build it using anything else

Library from Russ Cox: https://pkg.go.dev/rsc.io/pdf parses the testing PDF as

```txt
{Font:Arial CE FontSize:9.168 X:146.04 Y:663.5 W:0 S:2} {Font:Arial CE
FontSize:9.168 X:146.076672 Y:663.5 W:0 S:^@} {Font:Arial CE FontSize:9.168
X:146.076672 Y:663.5 W:0 S:^C} {Font:Arial CE FontSiz e:9.168 X:146.076672
Y:663.5 W:0 S:^@} {Font:Arial CE FontSize:9.168 X:146.076672 Y:663.5 W:0 S:^W} 
```

Great github.com/pdfcpu/pdfcpu

```
BT
/F1 9.24 Tf
1 0 0 1 410.86 709.66 Tm
/GS11 gs
0 g
/GS12 gs
0 G
[<00290044004E>-3<0057>-4<005800550044000300FE0011>] TJ
```

Commercial offering, haven't tried
https://unidoc.io/post/pdf-text-extraction-in-golang-with-unipdf/
https://docs.apryse.com/documentation/go/guides/features/extraction/text-extract/
