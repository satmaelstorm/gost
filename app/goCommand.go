package app

import (
	"github.com/fatih/color"
	"io"
	"os/exec"
	"strings"
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
			_, _ = errIo.Write([]byte(color.HiRedString("[ERROR]: ") + err.Error() + "\n"))
		}
	}
}
