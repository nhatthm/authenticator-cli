package app

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bool64/ctxd"
	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.nhat.io/authenticator"
	"go.nhat.io/otp"
)

func accountAddCommand(logger *ctxd.Logger) *cobra.Command {
	var cfg addAccountConfig

	cmd := &cobra.Command{
		Use:   "add [-n <namespace>] [--qr </path/to/qr-code-image>]",
		Short: "Add a new account",
		Long:  "Add a new account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addAccountToNamespace(cmd.Context(), cfg, *logger)
		},
	}

	cmd.Flags().StringVarP(&cfg.Namespace, "namespace", "n", cfg.Namespace, "namespace")
	cmd.Flags().StringVar(&cfg.QRCodeFile, "qr", "", "qr code")

	return cmd
}

type addAccountConfig struct {
	Namespace  string
	QRCodeFile string
}

func addAccountToNamespace(ctx context.Context, cfg addAccountConfig, logger ctxd.Logger) error { //nolint: cyclop,funlen
	var (
		namespace, account, totpSecret string
		err                            error
	)

	if cfg.QRCodeFile != "" {
		totpSecret, account, _, err = authenticator.ParseTOTPQRCode(cfg.QRCodeFile)
		if err != nil {
			logger.Error(ctx, "failed to parse QR code", "err", err)

			return err //nolint: wrapcheck
		}
	}

	allNamespaces, err := authenticator.GetAllNamespaceIDs()
	if err != nil {
		logger.Error(ctx, "failed to get all namespaces", "err", err)
	} else {
		logger.Debug(ctx, "available namespaces", "namespaces", allNamespaces)
	}

	namespace, account, totpSecret, err = getUserInput(ctx, cfg.Namespace, account, totpSecret, allNamespaces, logger)

	switch {
	case err != nil:
		return err

	case namespace == "":
		return errNamespaceIsRequired

	case account == "":
		return errAccountIsRequired

	case totpSecret == "":
		return errTOTPSecretIsRequired
	}

	fmt.Println(color.YellowString("Namespace:"), namespace)
	fmt.Println(color.YellowString("Account:"), account)
	fmt.Println()

	if !slices.Contains(allNamespaces, namespace) {
		err = authenticator.CreateNamespace(namespace, namespace)
		if err != nil {
			logger.Error(ctx, "failed to create namespace", "err", err)

			return err
		}

		logger.Debug(ctx, "namespace created", "namespace", namespace)
	}

	if err := authenticator.AddAccountToNamespace(namespace, account); err != nil {
		return err //nolint: wrapcheck
	}

	if err := authenticator.SetTOTPSecret(namespace, account, otp.TOTPSecret(totpSecret)); err != nil {
		return err //nolint: wrapcheck
	}

	fmt.Println(color.GreenString("✓"), "done")

	return nil
}

func getUserInput( //nolint: funlen
	ctx context.Context,
	defaultNamespace string,
	defaultAccount string,
	defaultTOTPSecret string,
	allNamespaces []string,
	logger ctxd.Logger,
) (string, string, string, error) {
	var (
		namespace, account, totpSecret, confirmTOTPSecret string
		err                                               error
	)

	fields := make([]huh.Field, 0, 4)
	namespace = defaultNamespace
	account = defaultAccount
	totpSecret = defaultTOTPSecret

	if namespace == "" {
		fields = append(fields, huh.NewInput().
			Title("Namespace").
			Prompt("? ").
			Suggestions(allNamespaces).
			Validate(func(s string) error {
				if s == "" {
					return errNamespaceIsRequired
				}

				return nil
			}).
			Value(&namespace),
		)
	}

	fields = append(fields, huh.NewInput().
		Title("Account").
		Prompt("? ").
		Validate(func(s string) error {
			if s == "" {
				return errAccountIsRequired
			}

			return nil
		}).
		Value(&account),
	)

	if totpSecret == "" {
		fields = append(fields,
			huh.NewInput().
				Title("TOTP Secret").
				Prompt("? ").
				Password(true).
				Validate(func(s string) error {
					if s == "" {
						return errTOTPSecretIsRequired
					}

					return nil
				}).
				Value(&totpSecret),
			huh.NewInput().
				Title("Confirm TOTP Secret").
				Prompt("? ").
				Password(true).
				Validate(func(s string) error {
					if s == "" {
						return errTOTPSecretConfirmNeeded
					} else if s != totpSecret {
						return errTOTPSecretMismatch
					}

					return nil
				}).
				Value(&confirmTOTPSecret),
		)
	}

	err = huh.NewForm(huh.NewGroup(fields...)).Run()
	if err != nil {
		if !errors.Is(err, huh.ErrUserAborted) {
			logger.Error(ctx, "failed to get user input", "err", err)
		}

		return "", "", "", err
	}

	return namespace, account, totpSecret, nil
}
