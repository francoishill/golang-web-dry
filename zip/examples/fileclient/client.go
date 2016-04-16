package fileclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/francoishill/golang-web-dry/zip/ziputils"
)

type Client interface {
	Download(serverUrl, localPath, remotePath string) error
	DownloadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) error
	Upload(serverUrl, localPath, remotePath string) error
	UploadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) error
	Delete(serverUrl, remotePath string) error
	DeleteDirFiltered(serverUrl, remotePath, dirFileFilterPattern string) error
	Move(serverUrl, oldRemotePath, newRemotePath string) error
	Stats(serverUrl, remotePath string) (*Stats, error)
}

func New(simpleLogger ziputils.SimpleLogger) Client {
	return &client{
		simpleLogger: simpleLogger,
	}
}

type client struct {
	simpleLogger ziputils.SimpleLogger
}

func (c *client) Download(serverUrl, localPath, remotePath string) (returnErr error) {
	return c.download(serverUrl, localPath, remotePath, "")
}
func (c *client) DownloadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) (returnErr error) {
	return c.download(serverUrl, localPath, remotePath, dirFileFilterPattern)
}

func (c *client) Upload(serverUrl, localPath, remotePath string) error {
	return c.upload(serverUrl, localPath, remotePath, "")
}
func (c *client) UploadDirFiltered(serverUrl, localPath, remotePath, dirFileFilterPattern string) error {
	return c.upload(serverUrl, localPath, remotePath, dirFileFilterPattern)
}

func (c *client) Delete(serverUrl, remotePath string) error {
	return c.delete(serverUrl, remotePath, "")
}
func (c *client) DeleteDirFiltered(serverUrl, remotePath, dirFileFilterPattern string) error {
	return c.delete(serverUrl, remotePath, dirFileFilterPattern)
}

func (c *client) Move(serverUrl, oldRemotePath, newRemotePath string) error {
	return c.move(serverUrl, oldRemotePath, newRemotePath)
}

func (c *client) Stats(serverUrl, remotePath string) (*Stats, error) {
	return c.getStats(serverUrl, remotePath)
}

func (c *client) checkServerResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if b, e := ioutil.ReadAll(resp.Body); e != nil {
			return fmt.Errorf("The server returned status code %d but could not read response body. Error: %s", resp.StatusCode, e.Error())
		} else {
			return fmt.Errorf("Server status code %d with response %s", resp.StatusCode, string(b))
		}
	}
	return nil
}

func (c *client) isDir(path string) (bool, error) {
	p, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer p.Close()
	info, err := p.Stat()
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (c *client) getFileSize(file *os.File) (int64, error) {
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (c *client) download(serverUrl, localPath, remotePath, dirFileFilterPattern string) (returnErr error) {
	defer CatchPanicAsError(&returnErr)

	var fileFilterQueryPart = ""
	if dirFileFilterPattern != "" {
		fileFilterQueryPart = "&filefilter=" + url.QueryEscape(dirFileFilterPattern)
	}

	resp, err := http.Get(serverUrl + "?path=" + url.QueryEscape(remotePath) + fileFilterQueryPart)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = c.checkServerResponse(resp); err != nil {
		return err
	}

	ziputils.SaveTarReaderToPath(c.simpleLogger, resp.Body, localPath)
	return nil
}

func (c *client) uploadFile(serverUrl, localPath, remotePath string) (returnErr error) {
	defer CatchPanicAsError(&returnErr)

	file, err := os.OpenFile(localPath, 0, 0600)
	if err != nil {
		return fmt.Errorf("Unable to read local file '%s', error: %s", localPath, err.Error())
	}
	defer file.Close()

	fileSize, err := c.getFileSize(file)
	if err != nil {
		return fmt.Errorf("Unable to get size of file '%s', error: %s", localPath, err.Error())
	}

	c.simpleLogger.Debug("Now starting to upload local file '%s' of size %s to remote path '%s'", localPath, humanize.IBytes(uint64(fileSize)), remotePath)
	url := serverUrl + "?path=" + url.QueryEscape(remotePath)
	ziputils.UploadFileToUrl(c.simpleLogger, url, "application/octet-stream", localPath, c.checkServerResponse)
	return nil
}

func (c *client) uploadDirectory(serverUrl, localPath, remotePath, dirFileFilterPattern string) (returnErr error) {
	defer CatchPanicAsError(&returnErr)

	c.simpleLogger.Debug("Now starting to upload local directory '%s' to remote '%s", localPath, remotePath)
	checkResponseFunc := c.checkServerResponse
	walkContext := ziputils.NewDirWalkContext(dirFileFilterPattern)
	ziputils.UploadDirectoryToUrl(c.simpleLogger, serverUrl+"?dir="+url.QueryEscape(remotePath), "application/octet-stream", localPath, walkContext, checkResponseFunc)
	return nil
}

func (c *client) upload(serverUrl, localPath, remotePath, dirFileFilterPattern string) error {
	if isDir, err := c.isDir(localPath); err != nil {
		return err
	} else if isDir {
		return c.uploadDirectory(serverUrl, localPath, remotePath, dirFileFilterPattern)
	} else {
		return c.uploadFile(serverUrl, localPath, remotePath)
	}
}

func (c *client) delete(serverUrl, remotePath, dirFileFilterPattern string) (returnErr error) {
	defer CatchPanicAsError(&returnErr)

	var fileFilterQueryPart = ""
	if dirFileFilterPattern != "" {
		fileFilterQueryPart = "&filefilter=" + url.QueryEscape(dirFileFilterPattern)
	}

	req, err := http.NewRequest("DELETE", serverUrl+"?path="+url.QueryEscape(remotePath)+fileFilterQueryPart, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return c.checkServerResponse(resp)
}

func (c *client) move(serverUrl, oldRemotePath, newRemotePath string) (returnErr error) {
	defer CatchPanicAsError(&returnErr)

	escapedOldPath := url.QueryEscape(oldRemotePath)
	escapedNewPath := url.QueryEscape(newRemotePath)
	req, err := http.NewRequest("PUT", serverUrl+"?action=move&path="+escapedOldPath+"&newpath="+escapedNewPath, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return c.checkServerResponse(resp)
}

func (c *client) getStats(serverUrl, remotePath string) (stats *Stats, returnErr error) {
	defer CatchPanicAsError(&returnErr)

	resp, err := http.Head(serverUrl + "?path=" + url.QueryEscape(remotePath))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = c.checkServerResponse(resp)
	if err != nil {
		return nil, err
	}

	if exists := resp.Header.Get("EXISTS"); exists == "" {
		return nil, fmt.Errorf("Could not find 'EXISTS' header")
	} else if exists == "0" {
		return &Stats{}, nil
	}

	if isDir := resp.Header.Get("IS_DIR"); isDir == "" {
		return nil, fmt.Errorf("Could not find 'IS_DIR' header")
	} else {
		return &Stats{
			Exists: true,
			IsDir:  isDir == "1",
		}, nil
	}
}
