package windowsprocessutils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
	"strings"
)

type WmicResultProcess struct {
	Caption         string
	CommandLine     string
	Description     string
	ExecutablePath  string
	Name            string
	ParentProcessId int
	ProcessId       int
}

func (w *WmicResultProcess) ExeEquals(exePath string) bool {
	trimCharsForExe := "'\""
	return strings.EqualFold(strings.Trim(w.ExecutablePath, trimCharsForExe), strings.Trim(exePath, trimCharsForExe))
}

func (w *WmicResultProcess) Kill() {
	proc, err := os.FindProcess(w.ProcessId)
	CheckError(err)

	err = proc.Kill()
	CheckError(err)
}
