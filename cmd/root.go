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
	Long: ``,
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = envviper.NewEnvViper()
	cfg.SetEnvParamsSimple("GOST")
	cfg.BindPFlag("v", rootCmd.PersistentFlags().Lookup("verbose"))
}
