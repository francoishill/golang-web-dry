package RunningCommand

import (
	"bufio"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RunningCommand struct {
	appendLock      sync.RWMutex
	commandDoneLock sync.RWMutex

	OnFinishAction   actionOnFinish
	AdditionalObject interface{}
	WorkingDirectory string
	CommandExePath   string
	CommandArguments []string
	CommandObj       *exec.Cmd
	CurrentFeedback  []string
	IsRunning        bool

	ReadStdOutChannel chan string
	ReadStdErrChannel chan string
	QuitChannel       chan bool

	startTime                time.Time
	CompletionTime           time.Time
	completionTimeAlreadySet bool

	WasKilledBeforeFinish bool

	tmpStdOutDone bool
	tmpStdErrDone bool
}

func (r *RunningCommand) AppendLine(lineStr string, isErrorLine bool) {
	r.appendLock.Lock()
	r.CurrentFeedback = append(r.CurrentFeedback, lineStr)
	r.appendLock.Unlock()

	defer recover()
	if !isErrorLine {
		if r.ReadStdOutChannel != nil {
			r.ReadStdOutChannel <- lineStr
		}
	} else {
		if r.ReadStdErrChannel != nil {
			r.ReadStdErrChannel <- lineStr
		}
	}
}

func (r *RunningCommand) recoverAndPrintGoRoutine(messagePrefix string) {
	recoveredObj := recover()
	if recoveredObj == nil {
		return
	}

	finalPrefix := strings.TrimRight(messagePrefix, ":")
	r.AppendLine(fmt.Sprintf("%s: %+v", finalPrefix, recoveredObj), true)
}

func (r *RunningCommand) onCommandDone() {
	r.commandDoneLock.Lock()
	defer r.commandDoneLock.Unlock()

	if r.OnFinishAction != nil {
		r.OnFinishAction(r)
	}
	r.IsRunning = false
	if !r.completionTimeAlreadySet {
		r.CompletionTime = time.Now()
		r.completionTimeAlreadySet = true
	}
}

func (r *RunningCommand) handleOutputReader(bufOutReader *bufio.Reader) {
	defer r.recoverAndPrintGoRoutine("Unable to finish reading std output of command")

	for {
		lineString, err := bufOutReader.ReadString('\n')
		if err != nil {
			break
		}
		if len(lineString) > 0 {
			r.AppendLine(lineString, false)
		}
	}
}

func (r *RunningCommand) handleErrorReader(bufErrReader *bufio.Reader) {
	defer r.recoverAndPrintGoRoutine("Unable to finish reading std error of command")

	for {
		lineString, err := bufErrReader.ReadString('\n')
		if err != nil {
			break
		}
		if len(lineString) > 0 {
			r.AppendLine(lineString, true)
		}
	}
}

func (r *RunningCommand) Start(timeoutDuration time.Duration) *RunningCommand {
	if r.QuitChannel == nil {
		r.QuitChannel = make(chan bool)
	}

	r.IsRunning = false
	r.tmpStdOutDone = false
	r.tmpStdErrDone = false

	cmd := exec.Command(r.CommandExePath, r.CommandArguments...)

	if r.WorkingDirectory == "" {
		workingDir := filepath.Dir(r.CommandExePath)
		cmd.Dir = workingDir
	} else {
		cmd.Dir = r.WorkingDirectory
	}

	cmd.Stdin = os.Stdin

	reader, err := cmd.StdoutPipe()
	CheckError(err)

	errReader, err := cmd.StderrPipe()
	CheckError(err)

	bufOutReader := bufio.NewReader(reader)
	go r.handleOutputReader(bufOutReader)

	bufErrReader := bufio.NewReader(errReader)
	go r.handleErrorReader(bufErrReader)

	err = cmd.Start()
	CheckError(err)

	r.CommandObj = cmd

	r.startTime = time.Now()
	r.IsRunning = true

	go func() {
		defer func() {
			if r.ReadStdOutChannel != nil {
				close(r.ReadStdOutChannel)
			}
			if r.ReadStdErrChannel != nil {
				close(r.ReadStdErrChannel)
			}
		}()
		defer r.onCommandDone()
		defer r.recoverAndPrintGoRoutine("Unable to finish waiting for command")

		r.AppendLine("Command has started.\n", false)

		err := cmd.Wait()
		if err != nil {
			CheckError(err)
		}

		if r.QuitChannel != nil {
			r.QuitChannel <- true
		}
	}()

	if timeoutDuration > 0 {
		time.AfterFunc(
			timeoutDuration,
			func() {
				defer r.onCommandDone()
				defer r.recoverAndPrintGoRoutine("Unable to kill command")

				r.AppendLine(fmt.Sprintf("Forcefully killing because timeout duration (%s) reached", timeoutDuration), true)
				r.ForceKill()
			},
		)
	}

	return r
}

func (r *RunningCommand) StartAndDoWithFeedback(timeoutDuration time.Duration, doWithFeedbackFunc DoWithFeedbackFunc) {
	r.ReadStdOutChannel = make(chan string)
	r.ReadStdErrChannel = make(chan string)
	r.QuitChannel = make(chan bool)

	r.Start(timeoutDuration)

loop:
	for {
		select {
		case lineStr, ok := <-r.ReadStdOutChannel:
			if !ok {
				break loop
			}
			if doWithFeedbackFunc != nil {
				doWithFeedbackFunc(false, lineStr)
			}
		case lineStr, ok := <-r.ReadStdErrChannel:
			if !ok {
				break loop
			}
			if doWithFeedbackFunc != nil {
				doWithFeedbackFunc(true, lineStr)
			}
		case <-r.QuitChannel:
			defer r.onCommandDone()
			defer func() {
				if rec := recover(); rec != nil {
					r.AppendLine(fmt.Sprintf("Cannot force kill, error: %+v", rec), true)
				}
			}()
			r.ForceKill()
			break loop
		}
	}

	for r.IsRunning {
	}
}

func (r *RunningCommand) ForceKill() {
	if !r.WasKilledBeforeFinish && r.IsRunning {
		r.WasKilledBeforeFinish = true
	}
	if !r.IsRunning {
		return
	}

	err := r.CommandObj.Process.Kill()
	CheckError(err)
}

func New(onFinishAction actionOnFinish, additionalObject interface{}, commandExePath string, commandArguments ...string) *RunningCommand {
	return &RunningCommand{
		OnFinishAction:    onFinishAction,
		AdditionalObject:  additionalObject,
		WorkingDirectory:  "",
		CommandExePath:    commandExePath,
		CommandArguments:  commandArguments,
		CommandObj:        nil,
		CurrentFeedback:   nil,
		IsRunning:         false,
		ReadStdOutChannel: nil,
		ReadStdErrChannel: nil,
		QuitChannel:       nil,
	}
}
