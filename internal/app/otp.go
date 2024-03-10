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
	cfg := otpConfig{
		Namespace: defaultNamespace,
	}

	cmd := &cobra.Command{
		Use:   "otp",
		Short: "Generate OTP",
		Long:  "Generate OTP",
		Args:  exactAccountArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.Output = cmd.OutOrStdout()

			return generateOTP(cmd.Context(), cfg, args[0], *logger)
		},
	}

	cmd.Flags().StringVarP(&cfg.Namespace, "namespace", "n", cfg.Namespace, "namespace")
	cmd.Flags().BoolVar(&cfg.CopyToClipboard, "copy", false, "copy the generated otp code to the clipboard")

	return cmd
}

type otpConfig struct {
	Namespace       string
	CopyToClipboard bool
	Output          io.Writer
}

func generateOTP(ctx context.Context, cfg otpConfig, account string, logger ctxd.Logger) error {
	otp, err := authenticator.GenerateTOTP(ctx, cfg.Namespace, account,
		authenticator.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	if !cfg.CopyToClipboard {
		_, _ = fmt.Fprintln(cfg.Output, otp)

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
