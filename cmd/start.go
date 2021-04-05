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
		initConfigToViperAliasesFiles(cmd, cfg)
		_ = cfg.BindPFlag(PackageName, cmd.Flags().Lookup(PackageName))
		startCommand(cfg, args)
	},
}

func init() {
	startCmd.Flags().String(PackageName, "", "package name for new project")
	startCmd.Flags().String(Aliases, "", "aliases and bundles file to rewrite default aliases")
	startCmd.Flags().String(AliasesAdd, "", "aliases and bundles file to add to defaults or loaded with --aliases")
}

func initConfigToViperAliasesFiles(cmd *cobra.Command, cfg *envviper.EnvViper) {
	_ = cfg.BindPFlag(Aliases, cmd.Flags().Lookup(Aliases))
	_ = cfg.BindPFlag(AliasesAdd, cmd.Flags().Lookup(AliasesAdd))
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
