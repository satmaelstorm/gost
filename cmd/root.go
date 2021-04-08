package cmd

import (
	"github.com/fatih/color"
	"github.com/satmaelstorm/envviper"
	"github.com/spf13/cobra"
)

const (
	//no color flag name
	NoColor = "no-color"
	//verbose flag full name
	Verbose = "verbose"
	//soft launch flag full name
	SoftLaunch = "soft-launch"
	//aliases parameter name
	Aliases = "aliases"
	//additional aliases parameter name
	AliasesAdd = "aliases-add"
)

var cfg *envviper.EnvViper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gost",
	Short: "Go Start project - helper for start go projects",
	Long: "Go Start (gost) - utility for help you to start go project without type some `go get ...` commands." +
		"\nCommand has some aliases for popular and useful go modules, and has some aliases for bundles of modules",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cfg.GetBool("no-color") {
			color.NoColor = true
		}
		if cfg.GetBool("s") {
			_, _ = color.New(color.FgHiMagenta, color.Bold).Println("Use soft Launch")
			_, _ = color.New(color.Reset).Println("")
		}
	},
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(modCmd, startCmd)
	rootCmd.PersistentFlags().BoolP(Verbose, "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP(SoftLaunch, "s", false, "soft launch - only print commands")
	rootCmd.PersistentFlags().Bool(NoColor, false, "disable color output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = envviper.NewEnvViper()
	cfg.SetEnvParamsSimple("GOST")
	_ = cfg.BindPFlag("v", rootCmd.PersistentFlags().Lookup(Verbose))
	_ = cfg.BindPFlag("s", rootCmd.PersistentFlags().Lookup(SoftLaunch))
	_ = cfg.BindPFlag(NoColor, rootCmd.PersistentFlags().Lookup(NoColor))
}
