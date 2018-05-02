package buildpacks

type DefaultDiff struct {
	Added, Removed string
}
type BuildpackDiff struct {
	Defaults map[string]DefaultDiff
	Added    map[string][]string
	Removed  map[string][]string
	Overlap  map[string][]string
}

func (b BuildpackDiff) Changes() bool {
	return len(b.Defaults) > 0 || len(b.Added) > 0 || len(b.Removed) > 0
}

func DiffBuildpackVersions(name, oldVersion, newVersion string, releases BuildpackReleases) (BuildpackDiff, error) {
	diff := BuildpackDiff{}
	if newVersion == oldVersion {
		return diff, nil
	}

	old, err := releases.Get(name, oldVersion)
	if err != nil {
		return diff, err
	}
	new, err := releases.Get(name, newVersion)
	if err != nil {
		return diff, err
	}

	diff.Defaults = diffDefaults(old.Defaults, new.Defaults)
	diff.Added = diffDependencies(old.Dependencies, new.Dependencies)
	diff.Removed = diffDependencies(new.Dependencies, old.Dependencies)
	diff.Overlap = overlapDependencies(old.Dependencies, new.Dependencies)

	return diff, nil
}

func diffDefaults(a, b map[string]string) map[string]DefaultDiff {
	diff := map[string]DefaultDiff{}

	for name, aVersion := range a {
		if bVersion, ok := b[name]; !ok || aVersion != bVersion {
			diff[name] = DefaultDiff{
				Removed: aVersion,
				Added:   bVersion,
			}
		}
	}

	return diff
}

func diffDependencies(a, b map[string][]string) map[string][]string {
	diff := map[string][]string{}

	for name, _ := range b {
		aMap := map[string]struct{}{}
		for _, ver := range a[name] {
			aMap[ver] = struct{}{}
		}

		for _, ver := range b[name] {
			if _, ok := aMap[ver]; !ok {
				diff[name] = append(diff[name], ver)
			}
		}
	}

	return diff
}

func overlapDependencies(a, b map[string][]string) map[string][]string {
	overlap := map[string][]string{}

	for name, _ := range b {
		aMap := map[string]struct{}{}
		for _, ver := range a[name] {
			aMap[ver] = struct{}{}
		}

		for _, ver := range b[name] {
			if _, ok := aMap[ver]; ok {
				overlap[name] = append(overlap[name], ver)
			}
		}
	}

	return overlap
}
