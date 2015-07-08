package builder

import (
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/v4run/bob/blogger"
)

/**
 * Builder struct.
 * appName   : name of the build file. Name of the root directory.
 * dir       : directory for building.
 * lastBuild : last time build was invoked.
 */
type Builder struct {
	appName   string
	dir       string
	lastBuild time.Time
}

/**
 * Returns a new builder with appropriate app name.
 *
 */
func NewBuilder(appName, dir string) Builder {
	if runtime.GOOS == "windows" && !strings.HasPrefix(appName, ".exe") {
		appName += ".exe"
	}
	return Builder{appName: appName, dir: dir}
}

func (b *Builder) LastBuild() time.Time {
	return b.lastBuild
}

func (b *Builder) SetLastBuild(lb time.Time) {
	b.lastBuild = lb
}

func (b *Builder) AppName() string {
	return b.appName
}

func (b *Builder) Build() bool {
	blogger.Info().Command("build").Message("Started.").Log()
	command := exec.Command("go", "build", "-o", b.appName)
	command.Dir = b.dir
	out, err := command.CombinedOutput()
	b.SetLastBuild(time.Now())
	if err != nil {
		blogger.Error().Command("build").Message("Failed.").Log()
		blogger.Error().Message(string(out), err.Error()).Log()
		return false
	}
	blogger.Info().Command("build").Message("Successful.").Log()
	return true
}
