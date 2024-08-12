package cmd

import (
	"fmt"
	"os"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Args:  cobra.ExactArgs(1),
	Short: "Delete a secret from the Linux keyring.",
	Long:  `Delete a secret from the Linux keyring by it's label.`,
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
			if err := collection.Delete(rootCmd.Name(), args[0]); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Unable to delete secret '%s': %v\n", args[0], err)
				os.Exit(1)
			}
		}
	},
}