package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/zip/ziputils"
	"log"
	"net/http"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	log.Println("[D] " + fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	log.Println("[I] " + fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	log.Println("[E] " + fmt.Sprintf(msg, args...))
}

type cliExtendedContext struct {
	*cli.Context
}

func (c *cliExtendedContext) RequireGlobalString(flagName string) string {
	val := c.GlobalString(flagName)
	if strings.TrimSpace(val) == "" {
		panic("Flag '" + flagName + "' is empty")
	}
	return val
}

type appContext struct {
	logger Logger
}

func (a *appContext) recoveryFunc(w http.ResponseWriter, req *http.Request, errorMessageSinglePlaceholder string) {
	if r := recover(); r != nil {
		a.logger.Error(errorMessageSinglePlaceholder, r)
		http.Error(w, fmt.Sprintf("Internal server error: %+v", r), 500)
		req.Body.Close()
	}
}

func (a *appContext) getFileOrFolderFromRequest(r *http.Request) (path string, isDir bool) {
	err := r.ParseForm()
	CheckError(err)

	saveFilePath := r.FormValue("file")
	saveDirPath := r.FormValue("dir")
	if saveFilePath == "" && saveDirPath == "" {
		panic("Cannot find 'file' or 'dir' query parameters...")
	} else if saveFilePath != "" && saveDirPath != "" {
		panic("Cannot specify both 'file' or 'dir' query parameters...")
	}

	if saveFilePath != "" {
		path = saveFilePath
		isDir = false
		return
	} else {
		path = saveDirPath
		isDir = true
		return
	}
}

func (a *appContext) getDirFileFilterPatternFromRequest(r *http.Request) string {
	err := r.ParseForm()
	CheckError(err)

	dirWalkFileFilter := r.FormValue("filefilter")
	return dirWalkFileFilter
}

func (a *appContext) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer a.recoveryFunc(w, r, "ERROR in handler: %+v")

		err := r.ParseForm()
		CheckError(err)

		path, isDir := a.getFileOrFolderFromRequest(r)

		if isDir {
			a.logger.Info("Receiving directory (zipped) %s", path)
			ziputils.SaveTarDirectoryReaderToFolder(a.logger, r.Body, path)
		} else {
			a.logger.Info("Receiving file to %s", path)
			ziputils.SaveReaderToFile(a.logger, r.Body, path)
		}
	} else if r.Method == "GET" {
		defer a.recoveryFunc(w, r, "ERROR in handler: %+v")

		err := r.ParseForm()
		CheckError(err)

		path, isDir := a.getFileOrFolderFromRequest(r)

		if isDir {
			a.logger.Info("Sending directory %s", path)
			walkContext := ziputils.NewDirWalkContext(a.getDirFileFilterPatternFromRequest(r))
			ziputils.UploadDirectoryToHttpResponseWriter(a.logger, w, path, walkContext)
		} else {
			a.logger.Info("Sending file %s", path)
			ziputils.UploadFileToHttpResponseWriter(a.logger, w, path)
		}
	} else if r.Method == "DELETE" {
		defer a.recoveryFunc(w, r, "ERROR in handler: %+v")

		path, isDir := a.getFileOrFolderFromRequest(r)

		if isDir {
			a.logger.Info("Deleting directory %s", path)
			walkContext := ziputils.NewDirWalkContext(a.getDirFileFilterPatternFromRequest(r))
			walkContext.DeleteDirectory(path)
		} else {
			a.logger.Info("Deleting file %s", path)
			err := os.Remove(path)
			CheckError(err)
		}
	} else {
		defer a.recoveryFunc(w, r, "ERROR in handler: %+v")

		panic("Unsupported method " + r.Method)
	}
}

func MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	port := c2.RequireGlobalString("port")

	logger := &defaultLogger{}
	h := &appContext{logger}

	http.HandleFunc("/", h.handler)

	log.Println(fmt.Sprintf("Now serving FileServer on port %s (process id is %d)", port, os.Getpid()))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	app := cli.NewApp()
	app.Name = "copyserver"
	app.Usage = "A http server to allow clients to upload and download files"
	app.Action = MainAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port,p",
			Value: "60878",
			Usage: "The port of the server",
		},
	}
	app.Run(os.Args)
}
