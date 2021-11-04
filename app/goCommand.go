package app

import (
	"context"
	"github.com/fatih/color"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func goCommand(args ...string) *exec.Cmd {
	return exec.Command("go", args...)
}

func goMod(subCommand string) *exec.Cmd {
	return goCommand("mod", subCommand)
}

func goGetWithUpdate(args ...string) *exec.Cmd {
	newArgs := make([]string, 1)
	newArgs[0] = "-u"
	newArgs = append(newArgs, args...)
	return goGet(newArgs...)
}

func goGet(args ...string) *exec.Cmd {
	newArgs := make([]string, 1)
	newArgs[0] = "get"
	newArgs = append(newArgs, args...)
	return goCommand(newArgs...)
}

func execCommand(command *exec.Cmd, verboseLevel int, softLaunch bool) (string, error) {
	if softLaunch {
		return command.String(), nil
	}
	buf := new(strings.Builder)
	buf.WriteString(color.HiGreenString("Executing: ") + command.String() + "\n")
	out, err := command.CombinedOutput()
	if verboseLevel > 0 {
		buf.WriteString(string(out) + "\n")
	}
	return buf.String(), err
}

func execCommands(
	outIo io.Writer,
	errIo io.Writer,
	commands []*exec.Cmd,
	verboseLevel int,
	softLaunch bool,
) {
	for _, cmd := range commands {
		r, err := execCommand(cmd, verboseLevel, softLaunch)
		_, _ = outIo.Write([]byte(r + "\n"))
		if err != nil {
			writeError(errIo, err.Error())
		}
	}
}

func execCommandsParallel(
	outIo io.Writer,
	errIo io.Writer,
	commands []*exec.Cmd,
	verboseLevel int,
	softLaunch bool,
	threads int,
) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	wg.Add(threads)
	outChan := make(chan *exec.Cmd, threads)
	writeGreen(outIo, "Start workers")
	for i := 0; i < threads; i++ {
		go execWorker(ctx, i, outChan, wg, outIo, errIo, verboseLevel, softLaunch)
	}
	writeGreen(outIo, "Workers started")
	lastCommand := make([]*exec.Cmd, 1)
	for _, cmd := range commands {
		if strings.HasSuffix(cmd.String(), "go mod download") {
			lastCommand[0] = cmd
		} else {
			outChan <- cmd
		}
	}
	cancel()
	wg.Wait()
	if lastCommand[0] != nil {
		execCommands(outIo, errIo, lastCommand, verboseLevel, softLaunch)
	}
	writeGreen(outIo, "Execution completed")
}

func execWorker(
	ctx context.Context,
	workerNum int,
	inChan <-chan *exec.Cmd,
	wg *sync.WaitGroup,
	outIo io.Writer,
	errIo io.Writer,
	verboseLevel int,
	softLaunch bool,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-inChan:
			writeGreen(outIo, "Executing by worker #" + strconv.Itoa(workerNum)+" "+cmd.String())
			r, err := execCommand(cmd, verboseLevel, softLaunch)
			_, _ = outIo.Write([]byte(r + "\n"))
			if err != nil {
				writeError(errIo, err.Error())
			}
		}
	}
}
