package app

import (
	"github.com/fatih/color"
	"io"
)

func writeError(out io.Writer, error string) {
	_, _ = out.Write([]byte(color.HiRedString("[ERROR]: ") + error + "\n"))
}

func writeGreen(out io.Writer, str string) {
	_, _ = out.Write([]byte(color.HiGreenString(str) + "\n"))
}
