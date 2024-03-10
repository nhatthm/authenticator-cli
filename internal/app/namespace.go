package app

import (
	"errors"

	"github.com/spf13/cobra"
)

const defaultNamespace = "default"

var (
	errNamespaceIsRequired = errors.New("namespace is required")
	errNoAccessToNamespace = errors.New("no access to the namespace")
)

func namespaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage namespaces",
		Long:  "Manage namespaces",
	}

	cmd.AddCommand(
		namespaceInfoCommand(),
		namespaceDeleteCommand(),
	)

	return cmd
}

func exactNamespaceArgs() cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if l := len(args); l == 0 {
			return errNamespaceIsRequired
		} else if l > 1 {
			return errTooManyArgs
		}

		return nil
	}
}
