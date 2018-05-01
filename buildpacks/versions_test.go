package buildpacks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/alphagov/paas-buildpack-notifier/buildpacks"
)

var _ = Describe("Versions", func() {
	Describe("BuildpacksFromFile", func() {
		It("should parse fixture file", func() {
			versions, err := BuildpacksFromFile(`fixtures/buildpacks.json`)
			Expect(err).ToNot(HaveOccurred())
			Expect(versions).To(Equal(BuildpackVersions{
				"binary_buildpack":     "v1.0.16",
				"go_buildpack":         "v1.8.18",
				"java_buildpack":       "v4.8",
				"nodejs_buildpack":     "v1.6.16",
				"php_buildpack":        "v4.3.48",
				"python_buildpack":     "v1.6.8",
				"ruby_buildpack":       "v1.7.11",
				"staticfile_buildpack": "v1.4.21",
			}))
		})
	})
})
