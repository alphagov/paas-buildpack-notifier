package buildpacks

import (
	"fmt"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type BuildpackRelease struct {
	Defaults     map[string]string
	Dependencies map[string][]string
}
type BuildpackReleases map[string]BuildpackRelease

func (b BuildpackReleases) key(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}
func (b BuildpackReleases) Add(name, version string, release BuildpackRelease) {
	b[b.key(name, version)] = release
}
func (b BuildpackReleases) Get(name, version string) (BuildpackRelease, error) {
	if !b.Has(name, version) {
		return BuildpackRelease{}, fmt.Errorf("unable to find release for %s", b.key(name, version))
	}

	return b[b.key(name, version)], nil
}
func (b BuildpackReleases) Has(name, version string) bool {
	_, ok := b[b.key(name, version)]
	return ok
}

type BuildpackDependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
type BuildpackManifest struct {
	Name            string                `yaml:"language"`
	DefaultVersions []BuildpackDependency `yaml:"default_versions"`
	Dependencies    []BuildpackDependency `yaml:"dependencies"`
}

func FetchBuildpackRelease(repoURL string) (BuildpackRelease, error) {
	release := BuildpackRelease{}

	resp, err := http.Get(repoURL)
	if err != nil {
		return release, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return release, fmt.Errorf("error getting %s: %s", resp.Request.URL, resp.Status)
	}

	manifest := BuildpackManifest{}
	err = yaml.NewDecoder(resp.Body).Decode(&manifest)
	if err != nil {
		return release, err
	}

	release.Defaults = map[string]string{}
	for _, dep := range manifest.DefaultVersions {
		release.Defaults[dep.Name] = dep.Version
	}
	release.Dependencies = map[string][]string{}
	for _, dep := range manifest.Dependencies {
		release.Dependencies[dep.Name] = append(release.Dependencies[dep.Name], dep.Version)
	}

	return release, err
}

func FetchBuildpackReleases(current, new BuildpackVersions) (BuildpackReleases, error) {
	releases := BuildpackReleases{}

	for _, buildpacks := range []BuildpackVersions{current, new} {
		for name, version := range buildpacks {
			if releases.Has(name, version) {
				continue
			}

			repoName := strings.Replace(name, "_buildpack", "-buildpack", 1)
			repoURL := fmt.Sprintf(
				"https://raw.githubusercontent.com/cloudfoundry/%s/%s/manifest.yml",
				repoName, version,
			)
			release, err := FetchBuildpackRelease(repoURL)
			if err != nil {
				fmt.Println("WARNING:", err)
				// return releases, err
			}

			releases.Add(name, version, release)
		}
	}

	return releases, nil
}
