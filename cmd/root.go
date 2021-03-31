package cmd

import (
	"github.com/satmaelstorm/envviper"
	"github.com/spf13/cobra"
)


var cfg *envviper.EnvViper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gost",
	Short: "Go Start project - helper for start go projects",
	Long: "Go Start (gost) - utility for help you to start go project without type some `go get ...` commands." +
		"\nCommand has some aliases for popular and useful go modules, and has some aliases for bundles of modules",
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(modCmd)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("softlaunch", "s", false, "soft launch - only print commands")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = envviper.NewEnvViper()
	cfg.SetEnvParamsSimple("GOST")
	cfg.BindPFlag("v", rootCmd.PersistentFlags().Lookup("verbose"))
	cfg.BindPFlag("s", rootCmd.PersistentFlags().Lookup("softlaunch"))
}
