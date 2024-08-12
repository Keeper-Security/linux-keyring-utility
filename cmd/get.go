package cmd

import (
	"fmt"
	"os"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.ExactArgs(1),
	Short: "Gets a secret from the Linux keyring.",
	Long:  `Get a secret from the Linux keyring by it's label and return the value as a string.`,
	Run: func(cmd *cobra.Command, args []string) {
		collection, err := secrets.DefaultCollection()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to get the default keyring: %v\n", err)
			os.Exit(1)
		}
		if err := collection.Unlock(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error unlocking the keyring: %v\n", err)
			os.Exit(1)
		} else {
			if secret, err := collection.Get(rootCmd.Name(), args[0]); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Unable to get secret '%s': %v\n", args[0], err)
				os.Exit(1)
			} else {
				fmt.Println(string(secret))
			}
		}
	},
}
