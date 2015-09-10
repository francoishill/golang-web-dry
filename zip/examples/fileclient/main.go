package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/zip/ziputils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct {
	d *log.Logger
	i *log.Logger
	e *log.Logger
}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	l.d.Println(fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	l.i.Println(fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	l.e.Println(fmt.Sprintf(msg, args...))
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

type timer struct {
	logger    Logger
	startTime time.Time
}

func (t *timer) printDuration() {
	duration := time.Now().Sub(t.startTime)
	t.logger.Debug("Duration %s", duration.String())
}

type appContext struct {
	logger Logger
}

func (a *appContext) getFileSize(file *os.File) int64 {
	fi, err := file.Stat()
	CheckError(err)
	return fi.Size()
}

func checkServerResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			panic(fmt.Sprintf("The server returned status code %d but could not read response body. Error: %s", e.Error()))
		}
		panic(fmt.Sprintf("Server status code %d with response %s", resp.StatusCode, string(b)))
	}
}

func (a *appContext) download(serverUrl, localPath, remotePath, dirFileFilterPattern string) {
	defer (&timer{a.logger, time.Now()}).printDuration()

	var fileFilterQueryPart = ""
	if dirFileFilterPattern != "" {
		fileFilterQueryPart = "&filefilter=" + url.QueryEscape(dirFileFilterPattern)
	}

	resp, err := http.Get(serverUrl + "?path=" + url.QueryEscape(remotePath) + fileFilterQueryPart)
	CheckError(err)
	defer resp.Body.Close()

	checkServerResponse(resp)

	ziputils.SaveTarReaderToPath(a.logger, resp.Body, localPath)
}

func (a *appContext) uploadFile(serverUrl, localPath, remotePath string) {
	defer (&timer{a.logger, time.Now()}).printDuration()

	file, err := os.OpenFile(localPath, 0, 0600)
	CheckError(err)

	a.logger.Debug("Now starting to upload local file '%s' of size %s to remote path '%s'", localPath, humanize.IBytes(uint64(a.getFileSize(file))), remotePath)
	resp, err := http.Post(serverUrl+"?file="+url.QueryEscape(remotePath), "application/octet-stream", file)
	CheckError(err)

	checkServerResponse(resp)
}

func (a *appContext) uploadDirectory(serverUrl, localPath, remotePath, dirFileFilterPattern string) {
	defer (&timer{a.logger, time.Now()}).printDuration()

	a.logger.Debug("Now starting to upload local directory '%s' to remote '%s", localPath, remotePath)
	checkResponseFunc := checkServerResponse
	walkContext := ziputils.NewDirWalkContext(dirFileFilterPattern)
	ziputils.UploadDirectoryToUrl(a.logger, serverUrl+"?dir="+url.QueryEscape(remotePath), "application/octet-stream", localPath, walkContext, checkResponseFunc)
}

func (a *appContext) isDir(path string) bool {
	p, err := os.Open(path)
	CheckError(err)
	defer p.Close()
	info, err := p.Stat()
	CheckError(err)
	return info.IsDir()
}

func (a *appContext) upload(serverUrl, localPath, remotePath, dirFileFilterPattern string) {
	if a.isDir(localPath) {
		a.uploadDirectory(serverUrl, localPath, remotePath, dirFileFilterPattern)
	} else {
		a.uploadFile(serverUrl, localPath, remotePath)
	}
}

func (a *appContext) delete(serverUrl, remotePath, dirFileFilterPattern string) {
	defer (&timer{a.logger, time.Now()}).printDuration()

	var fileFilterQueryPart = ""
	if dirFileFilterPattern != "" {
		fileFilterQueryPart = "&filefilter=" + url.QueryEscape(dirFileFilterPattern)
	}

	req, err := http.NewRequest("DELETE", serverUrl+"?path="+url.QueryEscape(remotePath)+fileFilterQueryPart, nil)
	CheckError(err)

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	checkServerResponse(resp)
}

func (a *appContext) MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	mode := c2.RequireGlobalString("mode")

	serverUrl := c2.RequireGlobalString("serverurl")
	remotePath := c2.RequireGlobalString("remotepath")

	switch mode {
	case "DOWNLOAD":
		localPath := c2.RequireGlobalString("localpath")
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		a.download(serverUrl, localPath, remotePath, dirFileFilterPattern)
		break
	case "UPLOAD":
		localPath := c2.RequireGlobalString("localpath")
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		a.upload(serverUrl, localPath, remotePath, dirFileFilterPattern)
		break
	case "DELETE":
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		a.delete(serverUrl, remotePath, dirFileFilterPattern)
		break
	default:
		panic("Unknown mode '" + mode + "'")
	}
}

func main() {
	logger := &defaultLogger{
		d: log.New(os.Stdout, "[D] ", log.Ldate|log.Ltime|log.Lshortfile),
		i: log.New(os.Stdout, "[I] ", log.Ldate|log.Ltime|log.Lshortfile),
		e: log.New(os.Stderr, "[E] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	context := &appContext{logger}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("%+v", r)
		}
	}()

	app := cli.NewApp()
	app.Name = "copyclient"
	app.Usage = "A http client to copy files to a server"
	app.Action = context.MainAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mode,m",
			Value: "",
			Usage: "The mode of the action (UPLOADFILE, DOWNLOADFILE, UPLOADFOLDER, DOWNLOADFOLDER)",
		},
		cli.StringFlag{
			Name:  "serverurl,s",
			Value: "",
			Usage: "Url to the file server",
		},
		cli.StringFlag{
			Name:  "localpath,l",
			Value: "",
			Usage: "The full absolute LOCAL path (file/folder)",
		},
		cli.StringFlag{
			Name:  "remotepath,r",
			Value: "",
			Usage: "The full absolute REMOTE path (file/folder)",
		},
		cli.StringFlag{
			Name:  "filefilter,ff",
			Value: "",
			Usage: "The golang filepath filter pattern (for file base name), see http://golang.org/pkg/path/filepath/#Match",
		},
	}
	app.Run(os.Args)
}
