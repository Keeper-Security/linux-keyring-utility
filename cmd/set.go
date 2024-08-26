package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [flags] <label> <secret string>",
	Args:  cobra.ExactArgs(2),
	Short: "Set a secret in the Linux keyring.",
	Long:  `Set the input string as a secret in the Linux keyring with the corresponding label.
Use -b or --base64 to encode the secret as base64 automatically before storing.
`,
	Run: func(cmd *cobra.Command, args []string) {
		collection, err := secrets.Collection(collection)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to get the default keyring: %v\n", err)
			os.Exit(1)
		}
		if err := collection.Unlock(); err == nil {
			if use_base64 {
				args[1] = base64.StdEncoding.EncodeToString([]byte(args[1]))
			}
			if err := collection.Set(application, args[0], []byte(args[1])); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Unable to create the secret '%s': %v\n", args[0], err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error unlocking the keyring: %v\n", err)
			os.Exit(1)
		}

	},
}
