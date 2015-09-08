package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/zip/ziputils"
	"strings"

	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
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

func (a *appContext) getFileSize(file *os.File) int64 {
	fi, err := file.Stat()
	CheckError(err)
	return fi.Size()
}

func (a *appContext) downloadFile(serverUrl, localPath, remotePath string) {
	resp, err := http.Get(serverUrl + "?file=" + url.QueryEscape(remotePath))
	CheckError(err)
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	CheckError(err)
	defer out.Close()

	log.Println("Now starting to download file of size", humanize.IBytes(uint64(resp.ContentLength)), "to path:", localPath)
	startTime := time.Now()
	_, err = io.Copy(out, resp.Body)
	CheckError(err)

	duration := time.Now().Sub(startTime)
	log.Println("Duration:", duration)
}

func (a *appContext) downloadDirectory(serverUrl, localPath, remotePath string) {
	resp, err := http.Get(serverUrl + "?zipdir=" + url.QueryEscape(remotePath))
	CheckError(err)
	defer resp.Body.Close()

	ziputils.SaveZipDirectoryReaderToFolder(resp.Body, localPath)
}

func (a *appContext) uploadFile(serverUrl, localPath, remotePath string) {
	file, err := os.OpenFile(localPath, 0, 0600)
	CheckError(err)

	log.Println("Now starting to upload file of size", humanize.IBytes(uint64(a.getFileSize(file))), "from path:", localPath)
	startTime := time.Now()

	_, err = http.Post(serverUrl+"?file="+url.QueryEscape(remotePath), "application/octet-stream", file)
	CheckError(err)

	duration := time.Now().Sub(startTime)
	log.Println("Duration:", duration)
}
func (a *appContext) uploadDirectory(serverUrl, localPath, remotePath string) {
	log.Println("Now starting to upload directory ", "from path:", localPath)
	startTime := time.Now()

	ziputils.UploadDirectoryToUrl(serverUrl+"?zipdir="+url.QueryEscape(remotePath), "application/octet-stream", localPath)

	duration := time.Now().Sub(startTime)
	log.Println("Duration:", duration)
}

func (a *appContext) MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	mode := c2.RequireGlobalString("mode")

	serverUrl := c2.RequireGlobalString("serverurl")
	localPath := c2.RequireGlobalString("localpath")
	remotePath := c2.RequireGlobalString("remotepath")

	switch mode {
	case "DOWNLOADFILE":
		a.downloadFile(serverUrl, localPath, remotePath)
		break
	case "DOWNLOADFOLDER":
		a.downloadDirectory(serverUrl, localPath, remotePath)
		break
	case "UPLOADFILE":
		a.uploadFile(serverUrl, localPath, remotePath)
		break
	case "UPLOADFOLDER":
		a.uploadDirectory(serverUrl, localPath, remotePath)
		break
	default:
		panic("Unknown mode '" + mode + "'")
	}
}

func main() {
	context := &appContext{&defaultLogger{}}

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
	}
	app.Run(os.Args)
}
