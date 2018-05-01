package buildpacks

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type BuildpackVersions map[string]string

func ParseBuildpackVersions(buildpacks []cfclient.Buildpack) (BuildpackVersions, error) {
	buildpackVersions := BuildpackVersions{}
	for _, buildpack := range buildpacks {
		filePrefix := strings.Replace(buildpack.Name, "_buildpack", "-buildpack", 1)
		re := regexp.MustCompile(fmt.Sprintf(`^%s-(v[\d.]+)\.zip`, filePrefix))
		if !re.MatchString(buildpack.Filename) {
			return BuildpackVersions{}, fmt.Errorf("failed to parse buildpack filename: %s", buildpack.Filename)
		}

		buildpackVersions[buildpack.Name] = re.FindStringSubmatch(buildpack.Filename)[1]
	}

	return buildpackVersions, nil
}

func BuildpacksFromCF(client *cfclient.Client) (BuildpackVersions, error) {
	buildpacks, err := client.ListBuildpacks()
	if err != nil {
		return BuildpackVersions{}, err
	}

	return ParseBuildpackVersions(buildpacks)
}

func BuildpacksFromFile(filename string) (BuildpackVersions, error) {
	file, err := os.Open(filename)
	if err != nil {
		return BuildpackVersions{}, err
	}

	var buildpackResponse cfclient.BuildpackResponse
	err = json.NewDecoder(file).Decode(&buildpackResponse)
	if err != nil {
		return BuildpackVersions{}, err
	}

	buildpacks := make([]cfclient.Buildpack, len(buildpackResponse.Resources))
	for i := 0; i < len(buildpacks); i++ {
		buildpacks[i] = buildpackResponse.Resources[i].Entity
	}

	return ParseBuildpackVersions(buildpacks)
}
