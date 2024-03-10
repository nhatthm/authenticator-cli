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
	cfg := deleteAccountConfig{
		Namespace: defaultNamespace,
	}

	cmd := &cobra.Command{
		Use:   "delete <namespace> <account>",
		Short: "Delete an account",
		Long:  "Delete an account",
		Args:  exactAccountArgs(),
		RunE: func(_ *cobra.Command, args []string) error {
			return deleteAccount(cfg.Namespace, args[0])
		},
	}

	cmd.Flags().StringVarP(&cfg.Namespace, "namespace", "n", cfg.Namespace, "namespace")

	return cmd
}

type deleteAccountConfig struct {
	Namespace string
}

func deleteAccount(namespace, account string) error {
	if namespace == "" {
		return errNamespaceIsRequired
	} else if account == "" {
		return errAccountIsRequired
	}

	hasAccess := sudo.Check()
	if !hasAccess {
		return errNoAccessToAccount
	}

	confirm := false

	input := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete %q account in %q namespace?", account, namespace)).
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

	err = authenticator.DeleteAccount(namespace, account)
	if err != nil {
		return err
	}

	fmt.Println(color.GreenString("✓"), "done")

	return nil
}
