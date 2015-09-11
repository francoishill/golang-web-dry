package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/errors/stacktraces/prettystacktrace"
	"github.com/francoishill/golang-web-dry/zip/ziputils"
	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/go-log/logger"
	"net/http"
	"os"
	"strings"
)

//
//PULL REQUEST START -- https://github.com/ian-kent/go-log
//
type multipleAppender struct {
	currentLayout   layout.Layout
	listOfAppenders []appenders.Appender
}

func Multiple(layout layout.Layout, appenders ...appenders.Appender) appenders.Appender {
	return &multipleAppender{
		listOfAppenders: appenders,
		currentLayout:   layout,
	}
}

func (this *multipleAppender) Layout() layout.Layout {
	return this.currentLayout
}

func (this *multipleAppender) SetLayout(l layout.Layout) {
	this.currentLayout = l
}

func (this *multipleAppender) Write(level levels.LogLevel, message string, args ...interface{}) {
	for _, appender := range this.listOfAppenders {
		appender.Write(level, message, args...)
	}
}

//
//PULL REQUEST END -- https://github.com/ian-kent/go-log
//

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct {
	l logger.Logger
}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	l.l.Debug(fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	log.Info(fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	log.Error(fmt.Sprintf(msg, args...))
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
		a.logger.Error("Stack: %s", prettystacktrace.GetPrettyStackTrace())
		http.Error(w, fmt.Sprintf("Internal server error: %+v", r), http.StatusInternalServerError)
		req.Body.Close()
	}
}

func (a *appContext) getPathFromRequest(r *http.Request) string {
	err := r.ParseForm()
	CheckError(err)

	path := r.FormValue("path")
	if path == "" {
		panic("Cannot find 'path' query parameter...")
	}
	return strings.TrimRight(path, ` /\`)
}

func (a *appContext) getRequiredQueryValue(r *http.Request, keyName string) string {
	err := r.ParseForm()
	CheckError(err)

	val := r.FormValue(keyName)
	if val == "" {
		panic("Cannot find '" + keyName + "' query parameter...")
	}
	return val
}

func (a *appContext) getFileOrFolderFromRequest(r *http.Request) (path string, isDir bool) {
	err := r.ParseForm()
	CheckError(err)

	saveFilePath := r.FormValue("path")
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

	return r.FormValue("filefilter")
}

func (a *appContext) isDir(path string) bool {
	p, err := os.Open(path)
	CheckError(err)
	defer p.Close()
	info, err := p.Stat()
	CheckError(err)
	return info.IsDir()
}

func (a *appContext) handler(w http.ResponseWriter, r *http.Request) {
	defer a.recoveryFunc(w, r, "ERROR in handler: %+v")

	if r.Method == "POST" {
		path, isDir := a.getFileOrFolderFromRequest(r)

		if isDir {
			a.logger.Info("Receiving directory (zipped) %s", path)
			ziputils.SaveTarReaderToPath(a.logger, r.Body, path)
		} else {
			a.logger.Info("Receiving file to %s", path)
			ziputils.SaveReaderToFile(a.logger, r.Body, path)
		}
	} else if r.Method == "GET" {
		path := a.getPathFromRequest(r)

		if a.isDir(path) {
			a.logger.Info("Sending directory %s", path)
			walkContext := ziputils.NewDirWalkContext(a.getDirFileFilterPatternFromRequest(r))
			ziputils.UploadDirectoryToHttpResponseWriter(a.logger, w, path, walkContext)
		} else {
			a.logger.Info("Sending file %s", path)
			ziputils.UploadFileToHttpResponseWriter(a.logger, w, path)
		}
	} else if r.Method == "DELETE" {
		path := a.getPathFromRequest(r)

		if a.isDir(path) {
			a.logger.Info("Deleting directory %s", path)
			walkContext := ziputils.NewDirWalkContext(a.getDirFileFilterPatternFromRequest(r))
			walkContext.DeleteDirectory(path)
		} else {
			a.logger.Info("Deleting file %s", path)
			err := os.Remove(path)
			CheckError(err)
		}
	} else if r.Method == "PUT" {
		action := a.getRequiredQueryValue(r, "action")
		switch strings.ToLower(action) {
		case "move":
			oldPath := a.getPathFromRequest(r)
			newPath := a.getRequiredQueryValue(r, "newpath")

			err := os.Rename(oldPath, newPath)
			CheckError(err)
			break
		default:
			panic("Unsupported action '" + action + "'")
		}
	} else if r.Method == "HEAD" {
		path := a.getPathFromRequest(r)

		a.logger.Info("Sending stats for path %s", path)

		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			w.Header().Set("EXISTS", "0")
			return
		}
		CheckError(err)

		w.Header().Set("EXISTS", "1")

		if info.IsDir() {
			w.Header().Set("IS_DIR", "1")
		} else {
			w.Header().Set("IS_DIR", "0")
		}
	} else {
		panic("Unsupported method " + r.Method)
	}
}

func getLogger() logger.Logger {
	logger := log.Logger()

	layoutToUse := layout.Pattern("%d [%p] %m") //date, level/priority, message

	rollingFileAppender := appenders.RollingFile("rolling-log.log", true)
	rollingFileAppender.MaxBackupIndex = 5
	rollingFileAppender.MaxFileSize = 20 * 1024 * 1024 // 20 MB
	rollingFileAppender.SetLayout(layoutToUse)

	consoleAppender := appenders.Console()
	consoleAppender.SetLayout(layoutToUse)
	logger.SetAppender(
		Multiple( //appenders.Multiple( ONCE PULL REQUEST OF ABOVE IS IN
			layoutToUse,
			rollingFileAppender,
			consoleAppender,
		))

	return logger
}

func MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	port := c2.RequireGlobalString("port")

	l := getLogger()
	defaultLogger := &defaultLogger{
		l,
	}
	h := &appContext{defaultLogger}

	http.HandleFunc("/", h.handler)

	l.Info("Now serving FileServer on port %s (process id is %d)", port, os.Getpid())
	l.Fatal(fmt.Sprintf("%s", http.ListenAndServe(":"+port, nil)))
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
