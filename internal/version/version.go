package version

import "runtime/debug"

const (
	dev   = "dev"
	devel = "(devel)"
)

func Resolve(buildVersion string) string {
	if buildVersion != dev {
		return buildVersion
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return dev
	}

	if info.Main.Version != "" && info.Main.Version != devel {
		return info.Main.Version
	}

	return dev
}
