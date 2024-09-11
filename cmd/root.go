package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var application = "lkru"
var collection = "default"
var use_base64 = false

var rootCmd = &cobra.Command{
	Use:   "lkru [flags] <get|set|del>",
	Short: "Linux Keyring Utility (lkru)",
	Long: `Linux Keyring Utility manages secrets in a Linux Keyring.
It uses the collection interface of the D-Bus Secrets API.
It has a trivial set, get, and delete interface.
There is no list or search.

It sets three attributes on the secret to facilitate lookups and namespacing.
1. The secret label becomes the 'id' attribute.
2. The application name is settable and 'lkru' by default.
3. The agent attribute always contains 'lkru (Linux Keyring Utility)'.

The default application name is 'lkru' and the default collection is the 'default' alias.
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
	getCmd.Flags().BoolVarP(&use_base64, "base64", "b", false, "Decode the secret from base64 before printing.")
	setCmd.Flags().BoolVarP(&use_base64, "base64", "b", false, "Encode the secret as base64 before storing.")
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(delCmd)
}
