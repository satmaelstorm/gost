package cmd

import (
	"github.com/satmaelstorm/envviper"
	"github.com/satmaelstorm/gost/app"
	"github.com/spf13/cobra"
	"os"
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "Add go modules with aliases",
	Long:  "Add go modules with aliases. Aliases list:\n" + app.GetDefaultAliasesHelp(),
	Run: func(cmd *cobra.Command, args []string) {
		initConfigToViperAliasesFiles(cmd, cfg)
		_ = cfg.BindPFlag(Threads, cmd.Flags().Lookup(Threads))
		modCommand(cfg, args)
	},
}

func init() {
	modCmd.Flags().String(Aliases, "", "aliases and bundles file to rewrite default aliases")
	modCmd.Flags().String(AliasesAdd, "", "aliases and bundles file to add to defaults or loaded with --aliases")
    modCmd.Flags().Int(Threads, 1, "executor threads")
}

func modCommand(cfg *envviper.EnvViper, args []string) {
	if !app.LoadAliasesFromFlags(os.Stdout, os.Stderr, cfg.GetString(Aliases), cfg.GetString(AliasesAdd)) {
		os.Exit(1)
	}
	executor := new(app.GoGetModules)
	if cfg.GetBool("v") {
		executor.VerboseLevel(1)
	}
	if cfg.GetBool("s") {
		executor.AsSoftLaunch()
	}
	executor.SetThreads(cfg.GetInt(Threads))
	executor.Run(args, app.GetAliases())
}
