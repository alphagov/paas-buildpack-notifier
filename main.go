package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/alphagov/paas-buildpack-notifier/buildpacks"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

func main() {
	var (
		hostname       = flag.String("a", "", "API hostname")
		token          = flag.String("t", "", "API token")
		buildpacksFile = flag.String("b", "", "Buildpacks migrate-to file")
	)
	flag.Parse()

	config := &cfclient.Config{
		ApiAddress: *hostname,
		Token:      *token,
	}
	client, err := cfclient.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	current, err := buildpacks.BuildpacksFromCF(client)
	if err != nil {
		log.Fatal(err)
	}

	new, err := buildpacks.BuildpacksFromFile(*buildpacksFile)
	if err != nil {
		log.Fatal(err)
	}

	releases, err := buildpacks.FetchBuildpackReleases(current, new)
	if err != nil {
		log.Fatal(err)
	}

	apps, err := client.ListAppsByQuery(url.Values{})
	if err != nil {
		log.Fatalln(err)
	}

	usedBuildpacks := map[string]map[string]struct{}{}
	for _, app := range apps {
		if _, ok := usedBuildpacks[app.Buildpack]; !ok {
			usedBuildpacks[app.Buildpack] = map[string]struct{}{}
		}
		usedBuildpacks[app.Buildpack][app.SpaceGuid] = struct{}{}
	}

	for name, _ := range current {
		spaceMap, used := usedBuildpacks[name]
		if !used {
			continue
		}

		spaces := make([]string, len(spaceMap))
		var spaceIndex int
		for guid, _ := range spaceMap {
			spaces[spaceIndex] = guid
			spaceIndex++
		}

		diff, err := buildpacks.DiffBuildpackVersions(name, current[name], new[name], releases)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Name:", name)
		if !diff.Changes() {
			fmt.Println("No changes..")
		} else {
			fmt.Println("Defaults:", diff.Defaults)
			fmt.Println("Added:", diff.Added)
			fmt.Println("Removed:", diff.Removed)
			fmt.Println("Overlap:", diff.Overlap)
			fmt.Println("Spaces:", spaces)
		}
		fmt.Println()
	}
}
