package app

import (
	"io"
	"os"
)

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
	goPath := os.Getenv("GOPATH")
	writeGreen(os.Stdout, "GOPATH: "+goPath)
	wd, err := os.Getwd()
	if err != nil {
		writeError(os.Stderr, err.Error())
		return "", false
	}
	newDir := goPath + "/src/" + pn
	writeGreen(outIo, "Mkdir "+newDir)
	if !isSoft {
		err = os.Mkdir(newDir, 0755)
		if err != nil {
			writeError(os.Stderr, err.Error())
			return "", false
		}
	}
	writeGreen(outIo, "Chdir "+newDir)
	if !isSoft {
		err = os.Chdir(newDir)
		if err != nil {
			writeError(os.Stderr, err.Error())
			return "", false
		}
	}
	r, err := execCommand(goMod("init"), verboseLevel, isSoft)
	_, _ = outIo.Write([]byte(r + "\n"))
	if err != nil {
		writeError(errIo, err.Error())
		return wd, false
	}
	return wd, true
}
