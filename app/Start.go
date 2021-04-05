package app

import (
	"github.com/fatih/color"
	"io"
	"os"
	"regexp"
)

var rePackageName = regexp.MustCompile(`^[A-Za-z][-._\w]*$`)

func StartCommand(
	outIo io.Writer,
	errIo io.Writer,
	pn string,
	verboseLevel int,
	isSoft bool,
) (workingDir string, ok bool) {
	if "" == pn {
		writeError(os.Stderr, "package-name is required")
		return "", false
	}
	if !rePackageName.MatchString(pn) {
		writeError(os.Stderr, "package-name must be "+rePackageName.String())
	}
	wd, err := os.Getwd()
	if err != nil {
		writeError(os.Stderr, err.Error())
		return "", false
	}

	if !isSoft {
		_, _ = outIo.Write([]byte(color.HiGreenString("Execute: ") + "mkdir " + pn + "\n"))
		err = os.Mkdir(pn, 0755)
		if err != nil {
			writeError(os.Stderr, err.Error())
			return "", false
		}
	} else {
		_, _ = outIo.Write([]byte("mkdir " + pn + "\n"))
	}

	if !isSoft {
		_, _ = outIo.Write([]byte(color.HiGreenString("Execute: ") + "chdir " + pn + "\n"))
		err = os.Chdir(pn)
		if err != nil {
			writeError(os.Stderr, err.Error())
			return "", false
		}
	} else {
		_, _ = outIo.Write([]byte("chdir " + pn + "\n"))
	}

	r, err := execCommand(goMod("init"), verboseLevel, isSoft)
	_, _ = outIo.Write([]byte(r + "\n"))
	if err != nil {
		writeError(errIo, err.Error())
		return wd, false
	}
	return wd, true
}
