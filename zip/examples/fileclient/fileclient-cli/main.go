package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"

	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/zip/examples/fileclient"
)

const (
	AppVersion = "0.0.2"
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

func (a *appContext) MainAction(c *cli.Context) {
	c2 := &cliExtendedContext{c}

	mode := c2.RequireGlobalString("mode")

	serverUrl := c2.RequireGlobalString("serverurl")
	remotePath := c2.RequireGlobalString("remotepath")

	client := fileclient.New(a.logger)

	defer (&timer{a.logger, time.Now()}).printDuration()
	switch mode {
	case "DOWNLOAD":
		localPath := c2.RequireGlobalString("localpath")
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		err := client.DownloadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern)
		CheckError(err)
		break
	case "UPLOAD":
		localPath := c2.RequireGlobalString("localpath")
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		err := client.UploadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern)
		CheckError(err)
		break
	case "DELETE":
		dirFileFilterPattern := c.GlobalString("filefilter") //Not required
		err := client.DeleteDirFiltered(serverUrl, remotePath, dirFileFilterPattern)
		CheckError(err)
		break
	case "STATS":
		stats, err := client.Stats(serverUrl, remotePath)
		CheckError(err)

		if stats.Exists {
			a.logger.Info("STATS_EXISTS=1")
		} else {
			a.logger.Info("STATS_EXISTS=0")
		}

		if stats.IsDir {
			a.logger.Info("STATS_IS_DIR=1")
		} else {
			a.logger.Info("STATS_IS_DIR=0")
		}

		break
	case "MOVE":
		newRemotePath := c.GlobalString("newpath") //Not required
		err := client.Move(serverUrl, remotePath, newRemotePath)
		CheckError(err)
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

	logger.i.Println("VERSION " + AppVersion)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("%+v", r)
			os.Exit(2)
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
		cli.StringFlag{
			Name:  "newpath,np",
			Value: "",
			Usage: "The new path, this is only currently applicable to the 'MOVE' method.",
		},
	}
	app.Version = AppVersion
	app.Run(os.Args)
}
