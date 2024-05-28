package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bool64/ctxd"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/nhatthm/authenticator-cli/internal/updater"
	"github.com/nhatthm/authenticator-cli/internal/version"
)

var (
	errMissingGithubToken      = errors.New("missing github token, please use --token or set the GITHUB_TOKEN environment variable")
	errFailedToRetrieveRelease = errors.New("failed to retrieve release to update")

	defaultUpdateTimeout = 3 * time.Minute
)

// selfUpdateCommand creates a new self-update command.
func selfUpdateCommand(logger *ctxd.Logger) *cobra.Command {
	cfg := selfUpdateConfig{
		version: "latest",
	}

	cmd := &cobra.Command{
		Use:   "self-update [<version>]",
		Short: "Check for new version of " + version.Info().AppName,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				cfg.version = args[0]
			}

			cfg.output = cmd.ErrOrStderr()
			locator := updater.NewReleaseLocator(cfg.token, defaultUpdateTimeout)

			return selfUpdate(cmd.Context(), cfg, locator, *logger)
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.token, "token", "", "github token that will be used to check for new versions. If not set GITHUB_TOKEN environment variable value is used")
	cmd.PersistentFlags().BoolVar(&cfg.force, "force", false, "update even if the same version is already installed")
	cmd.PersistentFlags().BoolVar(&cfg.dryRun, "dry-run", false, "check the version available but do not run actual update")

	return cmd
}

type selfUpdateConfig struct {
	token   string
	version string
	force   bool
	dryRun  bool
	output  io.Writer
}

func selfUpdate(ctx context.Context, cfg selfUpdateConfig, l updater.ReleaseLocator, log ctxd.Logger) error {
	// Find the release.
	updateTo, err := findRelease(ctx, cfg.output, l, version.Info().AppName, cfg.version)
	if err != nil {
		if errors.Is(err, updater.ErrNoGithubToken) {
			return errMissingGithubToken
		}

		log.Error(ctx, "failed to retrieve release", "error", err)

		return errFailedToRetrieveRelease
	}

	if cfg.force || updateTo.Name != version.Info().Version {
		// Fetch the release and update.
		if !cfg.dryRun {
			configureDownloader(cfg.output, cfg.token, defaultUpdateTimeout) //nolint: contextcheck

			if err := updater.SelfUpdate(ctx, updateTo); err != nil {
				log.Error(ctx, "failed to update",
					"error", err,
					"current_version", version.Info().Version,
					"new_version", updateTo.Name,
				)

				return fmt.Errorf("failed to update to %s: %w", updateTo.Name, err)
			}
		}

		_, _ = fmt.Fprintln(cfg.output, color.HiGreenString("⠿ Updated to version %s", updateTo.Name))

		return nil
	}

	_, _ = fmt.Fprintln(cfg.output, color.HiYellowString("⠿ Already up to date"))

	return nil
}

func findRelease(ctx context.Context, out io.Writer, l updater.ReleaseLocator, appName, version string) (r updater.Release, err error) {
	msg := "Checking for new versions of " + appName
	pb := newProgressBar(out, -1, msg)

	defer func() {
		_ = pb.Finish() //nolint: errcheck

		if err == nil {
			_, _ = fmt.Fprintf(out, "⠿ Found version: %s\n", r.Name)
		}
	}()

	return l.FindRelease(ctx, version) //nolint: wrapcheck
}

func configureDownloader(out io.Writer, token string, timeout time.Duration) {
	updater.ConfigureDownloader(token, timeout, func(c *http.Client) {
		c.Transport = downloadProgress(c.Transport, out, "Downloading Update")
	})
}

func newProgressBar(out io.Writer, max int64, desc string, options ...progressbar.Option) *progressbar.ProgressBar {
	options = append([]progressbar.Option{
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWidth(80), //nolint: gomnd
		progressbar.OptionFullWidth(),
		progressbar.OptionSetWriter(out),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionThrottle(65 * time.Millisecond), //nolint: gomnd
		progressbar.OptionSpinnerType(14),                 //nolint: gomnd
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	}, options...)

	return progressbar.NewOptions64(max, options...)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func downloadProgress(next http.RoundTripper, out io.Writer, desc string) roundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)
		if err != nil {
			return resp, err //nolint: wrapcheck
		}

		contentType := resp.Header.Get("Content-Type")
		contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64) //nolint: errcheck

		if !strings.HasPrefix(contentType, "application/") || contentLength < 1 {
			return resp, err //nolint: wrapcheck
		}

		pb := newProgressBar(out, -1, desc, progressbar.OptionShowBytes(true))
		body := progressbar.NewReader(resp.Body, pb)

		resp.Body = &reader{
			ReadCloser: &body,
			pb:         pb,
			onComplete: func(pb *progressbar.ProgressBar) {
				_ = pb.Finish() //nolint: errcheck
				_, _ = fmt.Fprintf(out, "⠿ Download Complete (%s)\n", byteCountSI(contentLength))
			},
		}

		return resp, err //nolint: wrapcheck
	}
}

type reader struct {
	io.ReadCloser
	pb         *progressbar.ProgressBar
	onComplete func(pb *progressbar.ProgressBar)
}

func (r *reader) Close() error {
	defer r.onComplete(r.pb)

	return r.ReadCloser.Close() //nolint: wrapcheck
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0

	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
