package cmd

import (
	"fmt"
	"os"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del [flags] <label> [label...]",
	Aliases: []string{"delete"},
	Args:  cobra.MinimumNArgs(1),
	Short: "Delete secret(s) from the Linux keyring.",
	Long:  `Delete one or more secrets from the Linux keyring by label.`,
	Run: func(cmd *cobra.Command, args []string) {
		collection, err := secrets.Collection(collection)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to get the default keyring: %v\n", err)
			os.Exit(1)
		}
		if err := collection.Unlock(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error unlocking the keyring: %v\n", err)
			os.Exit(1)
		} else {
			for _, label := range args {
				if err := collection.Delete(application, label); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Unable to delete secret '%s': %v\n", args[0], err)
					os.Exit(1)
				}
			}
		}
	},
}
