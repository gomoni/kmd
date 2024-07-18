# Kmd

> Karol Michal Daemon

HTTP/CLI interface for `tesseract-ocr` written in Go.

## Architecture

Client/server.

 * A client `kmc` is a simple static binary which talks to server. Curl would be
   enough too.
 * Server `kmd` is an HTTP server listening on unix socket `/run/user/1000/kmd.sock`. It uses
   https://github.com/klippa-app/go-pdfium?tab=readme-ov-file#webassembly to
   render pdf as png and gotesseract for the final OCR. It has a native dependencies
   on a tesseract-ocr, so is expected to run inside Docker.

## TODO

 1. finish the refactoring of stuff to internal - mainly params
 3. more tests

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

## Server

```sh
# build
$ mage build

# run server via systemd socket activation
# will be handier once Docker files will be ready
$ systemd-socket-activate -l ${XDG_RUNTIME_DIR:-/run/user/`id`}/kmd.sock ./kmd

$ curl --unix-socket /run/user/1000/kmd.sock http://localhost/
version: 5.4.1
languages:
 * ces
 * eng
```

## Client

TODO

# Why OCR?

> I've said this before and I'll say it again: PDF sucks
> Rajesh Koothrappali, Ph.D.

Library from Russ Cox: https://pkg.go.dev/rsc.io/pdf parses one of the testing PDF as

```txt
{Font:Arial CE FontSize:9.168 X:146.04 Y:663.5 W:0 S:2} {Font:Arial CE
FontSize:9.168 X:146.076672 Y:663.5 W:0 S:^@} {Font:Arial CE FontSize:9.168
X:146.076672 Y:663.5 W:0 S:^C} {Font:Arial CE FontSiz e:9.168 X:146.076672
Y:663.5 W:0 S:^@} {Font:Arial CE FontSize:9.168 X:146.076672 Y:663.5 W:0 S:^W} 
```

And an amazing github.com/pdfcpu/pdfcpu

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

Commercial offerings exists, haven't tried
https://unidoc.io/post/pdf-text-extraction-in-golang-with-unipdf/
https://docs.apryse.com/documentation/go/guides/features/extraction/text-extract/
