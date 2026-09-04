package version

import "runtime/debug"

const (
	defaultVersion = "dev"
	develVersion   = "(devel)"
)

func Resolve(buildVersion string) string {
	if buildVersion != defaultVersion {
		return buildVersion
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return defaultVersion
	}

	if info.Main.Version != "" && info.Main.Version != develVersion {
		return info.Main.Version
	}

	return defaultVersion
}
