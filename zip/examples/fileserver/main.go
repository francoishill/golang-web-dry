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
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct{}

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

func (h *appContext) recoveryFunc(req *http.Request, errorMessageSinglePlaceholder string) {
	if r := recover(); r != nil {
		h.logger.Error(errorMessageSinglePlaceholder, r)
		req.Body.Close()
	}
}

func (h *appContext) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer h.recoveryFunc(r, "ERROR in handler: %+v")

		err := r.ParseForm()
		CheckError(err)

		saveFilePath := r.FormValue("file")
		saveZipDirPath := r.FormValue("zipdir")
		if saveFilePath == "" && saveZipDirPath == "" {
			panic("Cannot find 'file' or 'zipdir' query parameters...")
		} else if saveFilePath != "" && saveZipDirPath != "" {
			panic("Cannot specify both 'file' or 'zipdir' query parameters...")
		}

		if saveFilePath != "" {
			//Receiving a file
			h.logger.Info("Receiving file to %s", saveFilePath)
			ziputils.SaveReaderToFile(r.Body, saveFilePath)
		} else if saveZipDirPath != "" {
			h.logger.Info("Receiving directory (zipped) %s", saveZipDirPath)
			ziputils.SaveZipDirectoryReaderToFolder(r.Body, saveZipDirPath)
		}
	} else {
		defer h.recoveryFunc(r, "ERROR in handler: %+v")

		err := r.ParseForm()
		CheckError(err)

		downloadFilePath := r.FormValue("file")
		downloadZipDirPath := r.FormValue("zipdir")
		if downloadFilePath == "" && downloadZipDirPath == "" {
			panic("Cannot find 'file' or 'zipdir' query parameters...")
		} else if downloadFilePath != "" && downloadZipDirPath != "" {
			panic("Cannot specify both 'file' or 'zipdir' query parameters...")
		}

		if downloadFilePath != "" {
			h.logger.Info("Sending file %s", downloadFilePath)
			ziputils.UploadFileToHttpResponseWriter(w, downloadFilePath)
		} else if downloadZipDirPath != "" {
			h.logger.Info("Sending directory %s", downloadZipDirPath)
			ziputils.UploadDirectoryToHttpResponseWriter(w, downloadZipDirPath)
		}
	}
}
func MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	port := c2.RequireGlobalString("port")

	h := &appContext{&defaultLogger{}}

	http.HandleFunc("/", h.handler)

	log.Println("Now serving on", port)
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
