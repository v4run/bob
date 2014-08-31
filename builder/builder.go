package builder

import (
	"runtime"
	"strings"
	"os/exec"
	"time"
)

/**
 * Builder struct.
 * appName: name of the build file. Name of the root directory.
 * dir: directory for building.
 * lastBuild: last time build was invoked.
 */
type Builder struct {
	appName string
	dir string
	lastBuild time.Time
}

/**
 * Returns a new builder with appropriate app name.
 *
 */
func NewBuilder(appName, dir string) Builder {
	if (runtime.GOOS == "windows" && !strings.HasPrefix(appName, ".exe")) {
		appName += ".exe"
	}
	return Builder{appName: appName, dir: dir}
}

func (b *Builder)LastBuild() time.Time {
	return b.lastBuild
}

func (b *Builder)SetLastBuild(lb time.Time) {
	b.lastBuild = lb
}

func (b *Builder) AppName() string {
	return b.appName
}

func (b *Builder) Build() ([]byte, error) {
	command := exec.Command("go", "build", "-o", b.appName)
	command.Dir = b.dir
	b.SetLastBuild(time.Now())
	return command.CombinedOutput()
}
