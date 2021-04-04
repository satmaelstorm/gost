package cmd

import (
	"github.com/satmaelstorm/envviper"
	"github.com/satmaelstorm/gost/app"
	"github.com/spf13/cobra"
	"os"
)

const PackageName = "package-name"

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Create new project and load modules",
	Long:  "Create new project and load modules by aliases\n",
	Run: func(cmd *cobra.Command, args []string) {
		startCommand(cfg, args)
	},
}

func init() {
	startCmd.Flags().String(PackageName, "", "package name for new project")
}

func startCommand(cfg *envviper.EnvViper, args []string) {
	verboseLevel := 0
	if cfg.GetBool("v") {
		verboseLevel = 1
	}
	wd, ok := app.StartCommand(os.Stdout, os.Stderr, cfg.GetString(PackageName), verboseLevel, cfg.GetBool("s"))
	if ok {
		modCommand(cfg, args)
	}
	if wd != "" && verboseLevel > 0 {
		_ = os.Chdir(wd)
	}
	if !ok {
		os.Exit(1)
	}
}
