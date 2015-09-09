package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/zip/ziputils"
	"io/ioutil"
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

type timer struct {
	startTime time.Time
}

func (t *timer) printDuration() {
	duration := time.Now().Sub(t.startTime)
	log.Println("Duration:", duration)
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

func (a *appContext) downloadFile(serverUrl, localPath, remotePath string) {
	resp, err := http.Get(serverUrl + "?file=" + url.QueryEscape(remotePath))
	CheckError(err)
	defer resp.Body.Close()

	checkServerResponse(resp)

	out, err := os.Create(localPath)
	CheckError(err)
	defer out.Close()

	log.Println("Now starting to download file of size", humanize.IBytes(uint64(resp.ContentLength)), "to path:", localPath)
	_, err = io.Copy(out, resp.Body)
	CheckError(err)
}

func (a *appContext) downloadDirectory(serverUrl, localPath, remotePath string) {
	defer (&timer{time.Now()}).printDuration()

	resp, err := http.Get(serverUrl + "?dir=" + url.QueryEscape(remotePath))
	CheckError(err)
	defer resp.Body.Close()

	checkServerResponse(resp)

	ziputils.SaveZipDirectoryReaderToFolder(resp.Body, localPath)
}

func (a *appContext) uploadFile(serverUrl, localPath, remotePath string) {
	defer (&timer{time.Now()}).printDuration()

	file, err := os.OpenFile(localPath, 0, 0600)
	CheckError(err)

	log.Println("Now starting to upload file of size", humanize.IBytes(uint64(a.getFileSize(file))), "from path:", localPath)
	resp, err := http.Post(serverUrl+"?file="+url.QueryEscape(remotePath), "application/octet-stream", file)
	CheckError(err)

	checkServerResponse(resp)
}

func (a *appContext) deleteFile(serverUrl, remotePath string) {
	defer (&timer{time.Now()}).printDuration()

	req, err := http.NewRequest("DELETE", serverUrl+"?file="+url.QueryEscape(remotePath), nil)
	CheckError(err)

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	checkServerResponse(resp)
}

func (a *appContext) uploadDirectory(serverUrl, localPath, remotePath string) {
	defer (&timer{time.Now()}).printDuration()

	log.Println("Now starting to upload directory ", "from path:", localPath)
	checkResponseFunc := checkServerResponse
	ziputils.UploadDirectoryToUrl(serverUrl+"?dir="+url.QueryEscape(remotePath), "application/octet-stream", localPath, checkResponseFunc)
}

func (a *appContext) deleteDirectory(serverUrl, remotePath string) {
	defer (&timer{time.Now()}).printDuration()

	req, err := http.NewRequest("DELETE", serverUrl+"?dir="+url.QueryEscape(remotePath), nil)
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
	case "DOWNLOADFILE":
		localPath := c2.RequireGlobalString("localpath")
		a.downloadFile(serverUrl, localPath, remotePath)
		break
	case "DOWNLOADFOLDER":
		localPath := c2.RequireGlobalString("localpath")
		a.downloadDirectory(serverUrl, localPath, remotePath)
		break
	case "UPLOADFILE":
		localPath := c2.RequireGlobalString("localpath")
		a.uploadFile(serverUrl, localPath, remotePath)
		break
	case "UPLOADFOLDER":
		localPath := c2.RequireGlobalString("localpath")
		a.uploadDirectory(serverUrl, localPath, remotePath)
		break
	case "DELETEFILE":
		a.deleteFile(serverUrl, remotePath)
		break
	case "DELETEFOLDER":
		a.deleteDirectory(serverUrl, remotePath)
		break
	default:
		panic("Unknown mode '" + mode + "'")
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(fmt.Sprintf("ERROR: %+v", r))
		}
	}()

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
