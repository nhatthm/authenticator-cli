package app

import (
	"errors"

	"github.com/bool64/ctxd"
	"github.com/spf13/cobra"
)

var (
	errAccountIsRequired       = errors.New("account is required")
	errNoAccessToAccount       = errors.New("no access to the account")
	errTOTPSecretIsRequired    = errors.New("totp secret is required")
	errTOTPSecretConfirmNeeded = errors.New("need to confirm totp secret")
	errTOTPSecretMismatch      = errors.New("totp secret does not match")
)

func accountCommand(logger *ctxd.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
		Long:  "Manage accounts",
	}

	cmd.AddCommand(
		accountAddCommand(logger),
		accountDeleteCommand(),
	)

	return cmd
}
