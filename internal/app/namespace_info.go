package app

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.nhat.io/authenticator"

	"github.com/nhatthm/authenticatorcli/internal/sudo"
)

func namespaceInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <namespace>",
		Short: "Get a namespace info",
		Long:  "Get a namespace info",
		Args: func(_ *cobra.Command, args []string) error {
			if l := len(args); l == 0 {
				return errNamespaceIsRequired
			} else if l > 1 {
				return errTooManyArgs
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return showNamespaceInfo(args[0])
		},
	}

	return cmd
}

func showNamespaceInfo(namespace string) error {
	hasAccess := sudo.Check()
	if !hasAccess {
		return errNoAccessToNamespace
	}

	n, err := authenticator.GetNamespace(namespace)
	if err != nil {
		return err
	}

	fmt.Println(color.YellowString("Namespace:"), namespace)

	if l := len(n.Accounts); l == 0 {
		fmt.Println(color.YellowString("No account found"))

		return nil
	} else if l == 1 {
		fmt.Println(color.YellowString("Account:"))
	} else {
		fmt.Println(color.YellowString("Accounts:"))
	}

	for _, account := range n.Accounts {
		fmt.Println("  ", account)
	}

	return nil
}
