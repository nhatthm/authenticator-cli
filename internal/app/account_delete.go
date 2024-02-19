package app

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.nhat.io/authenticator"

	"github.com/nhatthm/authenticatorcli/internal/sudo"
)

func accountDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <namespace> <account>",
		Short: "Delete an account",
		Long:  "Delete an account",
		Args: func(_ *cobra.Command, args []string) error {
			if l := len(args); l == 0 {
				return errNamespaceAndAccountAreRequired
			} else if l == 1 {
				return errAccountIsRequired
			} else if l > 2 {
				return errTooManyArgs
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return deleteAccount(args[0], args[1])
		},
	}

	return cmd
}

func deleteAccount(namespace, account string) error {
	hasAccess := sudo.Check()
	if !hasAccess {
		return errNoAccessToAccount
	}

	confirm := false

	input := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete %q account from %q?", account, namespace)).
		Description("This action cannot be undone.").
		Value(&confirm)

	err := huh.NewForm(huh.NewGroup(input)).Run()
	if err != nil {
		if !errors.Is(err, huh.ErrUserAborted) {
			err = fmt.Errorf("failed to confirm: %w", err)
		}

		return err
	}

	if !confirm {
		_, _ = fmt.Println(color.YellowString("operation canceled"))

		return nil
	}

	err = authenticator.DeleteAccountInNamespace(namespace, account)
	if err != nil {
		return err
	}

	fmt.Println(color.GreenString("✓"), "done")

	return nil
}
