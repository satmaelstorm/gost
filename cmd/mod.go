package cmd

import (
	"github.com/satmaelstorm/envviper"
	"github.com/satmaelstorm/gost/app"
	"github.com/spf13/cobra"
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "Add go modules with aliases",
	Long:  "Add go modules with aliases. Aliases list:\n" + app.GetDefaultAliasesHelp(),
	Run: func(cmd *cobra.Command, args []string) {
		modCommand(cfg, args)
	},
}

func modCommand(cfg *envviper.EnvViper, args []string) {
	executor := new(app.GoGetModules)
	if cfg.GetBool("v") {
		executor.VerboseLevel(1)
	}
	if cfg.GetBool("s") {
		executor.AsSoftLaunch()
	}
	executor.Run(args, app.GetDefaultAliases())
}
