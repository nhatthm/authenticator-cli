package app

import (
	"errors"

	"github.com/spf13/cobra"
)

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
