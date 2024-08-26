package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var application = "lkru"
var collection = "login"

var rootCmd = &cobra.Command{
	Use:   "lkru [flags] <get|set|del>",
	Short: "Linux Keyring Utility (lkru)",
	Long: `lkru is a Linux Keyring Utility.
It manages secrets in a Linux Keyring using the collection interface of the D-Bus Secrets API.
It has a trivial set, get, and delete interface where set always creates and overwrites.
There is no list or search functionality.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&application, "application", "a", application, "The application name to use.")
	rootCmd.PersistentFlags().StringVarP(&collection, "collection", "c", collection, "The collection name to use.")
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(delCmd)
}
