package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	secrets "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_collection"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [flags] <label> <secret string>",
	Args:  cobra.ExactArgs(2),
	Short: "Set a secret in the Linux keyring",
	Long: `Set secret string as a secret in the Linux keyring with the corresponding label.
If the secret string is "-", lkru reads it from standard input.
Use -b or --base64 to encode the secret string as base64 before storing it.
`,
	Run: func(cmd *cobra.Command, args []string) {
		collection, err := secrets.Collection(collection)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Unable to get the default keyring: %v\n", err)
			os.Exit(1)
		}
		if err := collection.Unlock(); err == nil {
			if len(args) == 2 && args[1] == "-" {
				scanner := bufio.NewScanner(os.Stdin)

				var lines []string
				for {
					scanner.Scan()
					line := scanner.Text()
					if len(line) == 0 {
						break
					}
					lines = append(lines, line)
				}
				if err := scanner.Err(); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Unable to read the secret '%s' from standard input: %v\n", args[0], err)
					os.Exit(1)
				}
				args[1] = strings.Join(lines, "\n")
			}
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
