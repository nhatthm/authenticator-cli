package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.nhat.io/authenticator"

	"github.com/nhatthm/authenticatorcli/internal/sudo"
)

func namespaceDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <namespace>",
		Short: "Delete a namespace",
		Long:  "Delete a namespace",
		Args:  exactNamespaceArgs(),
		RunE: func(_ *cobra.Command, args []string) error {
			return deleteNamespace(args[0])
		},
	}

	return cmd
}

func deleteNamespace(namespace string) error {
	hasAccess := sudo.Check()
	if !hasAccess {
		return errNoAccessToNamespace
	}

	confirm := false

	input := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete %q namespace?", namespace)).
		Description("All accounts in the namespace will be deleted. This action cannot be undone.").
		Value(&confirm)

	err := huh.NewForm(huh.NewGroup(input)).Run()
	if err != nil {
		if !errors.Is(err, huh.ErrUserAborted) {
			err = fmt.Errorf("failed to confirm: %w", err)
		}

		return err
	}

	if !confirm {
		_, _ = fmt.Fprintln(os.Stderr, color.YellowString("operation canceled"))

		return nil
	}

	err = authenticator.DeleteNamespace(namespace)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(os.Stderr, color.GreenString("✓"), "done")

	return nil
}
