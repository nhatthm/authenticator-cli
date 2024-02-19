package app

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/bool64/ctxd"
	"github.com/spf13/cobra"
	"go.nhat.io/authenticator"
	"go.nhat.io/exec"
)

func otpCommand(logger *ctxd.Logger) *cobra.Command {
	copyToClipboard := false

	cmd := &cobra.Command{
		Use:   "otp",
		Short: "Generate OTP",
		Long:  "Generate OTP",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateOTP(cmd.Context(), cmd.OutOrStdout(), args[0], args[1], copyToClipboard, *logger)
		},
	}

	cmd.Flags().BoolVar(&copyToClipboard, "copy", false, "copy the generated otp code to the clipboard")

	return cmd
}

func generateOTP(ctx context.Context, stdout io.Writer, namespace, account string, copyToClipboard bool, logger ctxd.Logger) error {
	otp, err := authenticator.GenerateTOTP(ctx, namespace, account)
	if err != nil {
		return err
	}

	if !copyToClipboard {
		_, _ = fmt.Fprintln(stdout, otp)

		return nil
	}

	return copyOTP(ctx, otp.String(), logger)
}

func copyOTP(ctx context.Context, otp string, logger ctxd.Logger) error {
	var (
		command string
		args    []string
	)

	switch runtime.GOOS {
	case "darwin":
		command = "pbcopy"

	case "linux":
		command = "xclip"
		args = []string{"-selection", "c"}

	case "windows":
		command = "clip"

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS) //nolint: goerr113
	}

	_, err := exec.RunWithContext(ctx, command,
		exec.WithArgs(args...),
		exec.WithStdout(io.Discard),
		exec.WithStderr(io.Discard),
		exec.WithStdin(strings.NewReader(otp)),
		exec.WithLogger(logger),
	)

	return err //nolint: wrapcheck
}
