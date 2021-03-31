package cmd

import (
	"github.com/satmaelstorm/gost/app"
	"github.com/spf13/cobra"
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "Add go modules with aliases",
	Long: "Add go modules with aliases",
	Run: modCommand,
}

func modCommand(cmd *cobra.Command, args []string) {
	app.GetDefaultAliases()
}