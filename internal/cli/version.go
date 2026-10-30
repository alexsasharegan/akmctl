package cli

import "runtime/debug"

func init() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		Version = resolveVersion(Version, bi)
	}
}

// resolveVersion names the build. A version set at link time (release archives) wins;
// a module version (go install of a tag) comes next; a development build names the
// commit it was built from (semver build metadata).
func resolveVersion(v string, bi *debug.BuildInfo) string {
	if v != "0.0.0-dev" {
		return v
	}
	if mv := bi.Main.Version; len(mv) > 1 && mv[0] == 'v' && mv != "(devel)" {
		return mv[1:]
	}
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" && len(s.Value) >= 7 {
			return v + "+" + s.Value[:7]
		}
	}
	return v
}
