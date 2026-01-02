/*
Built by Adithya Kaushik
*/
package cmd
import (
	"os"
	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
	Use: `GoExcel`,
	Short: `Excel Engine in Go`,
	Long: `Excel Engine in Go without the UI`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


