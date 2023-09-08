package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "paste",
	Short: "Yet another Paste Bin",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    // This is now defined in MONGODB_URI env var
	// rootCmd.PersistentFlags().String("db-uri", "mongodb://127.0.0.1:27017/paste", "MongoDB URI")
}
