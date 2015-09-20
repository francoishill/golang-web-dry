package windowsprocessutils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"strconv"
	"strings"
)

func parseWmicOutput(lines []string) WmicResultProcessSlice {
	results := []*WmicResultProcess{}

	var currentResult *WmicResultProcess
	firstKeyName := "Caption"
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		indexEqualSign := strings.Index(trimmedLine, "=")
		if indexEqualSign == -1 {
			continue
		}

		if strings.HasPrefix(trimmedLine, firstKeyName+"=") {
			currentResult = &WmicResultProcess{}
			results = append(results, currentResult)
		}

		keyName := trimmedLine[:indexEqualSign]
		val := trimmedLine[indexEqualSign+1:]

		switch {
		case strings.EqualFold(keyName, "Caption"):
			currentResult.Caption = val
			break
		case strings.EqualFold(keyName, "CommandLine"):
			currentResult.CommandLine = val
			break
		case strings.EqualFold(keyName, "Description"):
			currentResult.Description = val
			break
		case strings.EqualFold(keyName, "ExecutablePath"):
			currentResult.ExecutablePath = val
			break
		case strings.EqualFold(keyName, "Name"):
			currentResult.Name = val
			break
		case strings.EqualFold(keyName, "ParentProcessId"):
			tmpInt, err := strconv.ParseInt(val, 10, 32)
			CheckError(err)
			currentResult.ParentProcessId = int(tmpInt)
			break
		case strings.EqualFold(keyName, "ProcessId"):
			tmpInt, err := strconv.ParseInt(val, 10, 32)
			CheckError(err)
			currentResult.ProcessId = int(tmpInt)
			break
		}
	}

	return WmicResultProcessSlice(results)
}
