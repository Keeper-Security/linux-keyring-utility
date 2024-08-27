package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var application = "lkru"
var collection = "login"
var use_base64 = false

var rootCmd = &cobra.Command{
	Use:   "lkru [flags] <get|set|del>",
	Short: "Linux Keyring Utility (lkru)",
	Long: `lkru is a Linux Keyring Utility.
It manages secrets in a Linux Keyring using the collection interface of the D-Bus Secrets API.
It has a trivial set, get, and delete interface where set creates and always overwrites.
There is no list or search functionality.

It sets attributes on the secret to facilitate namespacing.
The application name is an attribute on the secret.
There is an agent attribute containing 'lkru (Linux Keyring Utility)'.
And the label becomes the _id_ attribute on the secret.

The default application name is 'lkru' and the default collection is 'login'.
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
