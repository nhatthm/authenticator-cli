package updater_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/authenticator-cli/internal/updater"
)

func TestIgnoreReleaseVersion(t *testing.T) {
	t.Parallel()

	assert.True(t, updater.IgnoreReleaseVersion("dev-foo"))
	assert.True(t, updater.IgnoreReleaseVersion("edge-foo"))
	assert.True(t, updater.IgnoreReleaseVersion("PR-foo"))
	assert.False(t, updater.IgnoreReleaseVersion("1.2.3"))
}
