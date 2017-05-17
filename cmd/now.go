package cmd

import (
	"fmt"
	veevaLib "github.com/chrisbenson/viiv-veeva-integration/pkg/veeva"
	"github.com/spf13/cobra"
)

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Manual command-line interface to start Veeva integration.",
	Long: `
	Manual command-line interface to start Veeva integration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 1 {
			err = veevaLib.Start(args[0])
		} else {
			err = veevaLib.Start("")
		}
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(nowCmd)
}
