@echo off
cls

echo Can build again without docker once fixed in Golang 1.10
REM https://github.com/golang/go/issues/24232
REM go build -o fileserver-hidden.exe -ldflags "-H windowsgui"

docker run -v "%GOPATH%\src":/go/src -w /go/src/github.com/francoishill/golang-web-dry/zip/examples/fileserver -e GOOS=windows -e GOARC=amd64 golang:1.9 go build -o fileserver-hidden.exe -ldflags "-H windowsgui"
docker run -v "%GOPATH%\src":/go/src -w /go/src/github.com/francoishill/golang-web-dry/zip/examples/fileserver -e GOOS=windows -e GOARC=amd64 golang:1.9 go build -o fileserver-visible.exe

REM pause