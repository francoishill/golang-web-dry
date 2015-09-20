package windowsprocessutils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindProcessesByName(exePath string) WmicResultProcessSlice {
	exeNameOnly := filepath.Base(exePath)

	cmd := exec.Command(
		"wmic",
		"process",
		"where",
		fmt.Sprintf("name='%s'", exeNameOnly),
		"get",
		"Caption,CommandLine,Description,ExecutablePath,Name,ParentProcessId,ProcessId",
		"/format:list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Unable to run wmic for name '%s'. Error was: %s. CombinedOutput was: %s", exeNameOnly, err.Error(), string(output)))
	}

	return parseWmicOutput(strings.Split(string(output), "\n"))
}
