package updater

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bool64/ctxd"
	"github.com/google/go-github/v60/github"
	"github.com/hellofresh/updater-go/v3"
	"golang.org/x/oauth2"

	"github.com/nhatthm/authenticatorcli/internal/version"
)

const (
	githubOwner = "nhatthm"

	// "dev-" is used for local dev builds.
	devPrefix = "dev-"
	// "edge-" is used for edge builds.
	edgePrefix = "edge-"
	// "PR-" is used for PR builds.
	prPrefix = "PR-"
)

var (
	// ErrNoGithubToken indicates that no GitHub token was provided.
	ErrNoGithubToken = errors.New("no github token provided")
	// ErrReleaseAssetNotFound indicates that release asset was not found for the current platform.
	ErrReleaseAssetNotFound = errors.New("release asset not found for the current platform")
	// ErrReleaseNoAssets indicates that release has no assets.
	ErrReleaseNoAssets = errors.New("release has no assets")
	// ErrReleaseNoTag indicates that release has no name.
	ErrReleaseNoTag = errors.New("release has no tag")
	// ErrReleaseNotLatest indicates that release is not the latest.
	ErrReleaseNotLatest = errors.New("outdated version")
)

var releaseAssetPattern = fmt.Sprintf("-%s-%s", runtime.GOOS, runtime.GOARCH)

// Release is an alias of updater.Release.
type Release = updater.Release

// ReleaseLocator locates releases.
type ReleaseLocator interface {
	LatestRelease(ctx context.Context) (updater.Release, error)
	FindRelease(ctx context.Context, version string) (updater.Release, error)
}

// MissingGitHubTokenLocator locates releases without GitHub token.
type MissingGitHubTokenLocator struct{}

// LatestRelease returns the latest release.
func (MissingGitHubTokenLocator) LatestRelease(context.Context) (updater.Release, error) {
	return updater.Release{}, ErrNoGithubToken
}

// FindRelease returns the release by version.
func (MissingGitHubTokenLocator) FindRelease(context.Context, string) (updater.Release, error) {
	return updater.Release{}, ErrNoGithubToken
}

// DefaultReleaseLocator locates releases for a single app.
type DefaultReleaseLocator struct {
	service *github.RepositoriesService

	owner, repo string
}

// LatestRelease returns the latest release.
func (l DefaultReleaseLocator) LatestRelease(ctx context.Context) (updater.Release, error) {
	r, _, err := l.service.GetLatestRelease(ctx, l.owner, l.repo)
	if err != nil {
		return updater.Release{}, fmt.Errorf("unable to get latest release: %w", err)
	}

	return findReleaseAsset(r)
}

// FindRelease returns the release by version.
func (l DefaultReleaseLocator) FindRelease(ctx context.Context, version string) (updater.Release, error) {
	if version == "" || version == "latest" || version == "stable" {
		return l.LatestRelease(ctx)
	}

	r, _, err := l.service.GetReleaseByTag(ctx, l.owner, l.repo, version)
	if err != nil {
		return updater.Release{}, fmt.Errorf("unable to get release by tag: %w", err)
	}

	return findReleaseAsset(r)
}

// NewReleaseLocator creates a new ReleaseLocator.
func NewReleaseLocator(githubToken string, timeout time.Duration) ReleaseLocator {
	githubToken = ensureGitHubToken(githubToken)
	if githubToken == "" {
		return MissingGitHubTokenLocator{}
	}

	return DefaultReleaseLocator{
		service: newClient(githubToken, timeout),
		owner:   githubOwner,
		repo:    version.Info().RepositoryName,
	}
}

// SelfUpdate update the current executable to the release.
func SelfUpdate(ctx context.Context, release updater.Release) error {
	return updater.SelfUpdate(ctx, release) //nolint: wrapcheck
}

// IgnoreReleaseVersion checks if the version is dev build so that force update should be skipped.
func IgnoreReleaseVersion(version string) bool {
	for _, prefix := range []string{devPrefix, edgePrefix, prPrefix} {
		if strings.HasPrefix(version, prefix) {
			return true
		}
	}

	return false
}

// CheckLatestRelease checks if the running application has the latest available version.
func CheckLatestRelease(ctx context.Context, l ReleaseLocator, currentVersion string, log ctxd.Logger) error {
	if IgnoreReleaseVersion(currentVersion) {
		log.Debug(ctx, "ignored experimental version", "version", currentVersion)

		return nil
	}

	release, err := l.LatestRelease(ctx)
	if err != nil {
		log.Error(ctx, "failed to check for new release", "error", err)

		return err //nolint: wrapcheck
	}

	log.Debug(ctx, "found latest release", "version", release.Name)

	if release.Name != currentVersion {
		return ErrReleaseNotLatest
	}

	return nil
}

// ConfigureDownloader configures the updater.DefaultDownloader.
func ConfigureDownloader(githubToken string, timeout time.Duration, opts ...func(c *http.Client)) {
	githubToken = ensureGitHubToken(githubToken)
	client := newHTTPClient(githubToken, timeout)

	for _, o := range opts {
		o(client)
	}

	updater.DefaultDownloader = updater.NewHTTPDownloader(client)
}

func newClient(githubToken string, timeout time.Duration) *github.RepositoriesService {
	return github.NewClient(newHTTPClient(githubToken, timeout)).Repositories
}

func newHTTPClient(githubToken string, timeout time.Duration) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	tc.Timeout = timeout

	return tc
}

func findReleaseAsset(r *github.RepositoryRelease) (updater.Release, error) {
	releaseName := r.GetName()
	if len(releaseName) == 0 {
		return updater.Release{}, ErrReleaseNoTag
	}

	if len(r.Assets) == 0 {
		return updater.Release{}, ErrReleaseNoAssets
	}

	for _, a := range r.Assets {
		if strings.Contains(a.GetName(), releaseAssetPattern) {
			return updater.Release{
				Name:  releaseName,
				Asset: a.GetName(),
				URL:   a.GetURL(),
			}, nil
		}
	}

	return updater.Release{}, ErrReleaseAssetNotFound
}

func ensureGitHubToken(token string) string {
	if token == "" {
		return os.Getenv("GITHUB_TOKEN")
	}

	return token
}
