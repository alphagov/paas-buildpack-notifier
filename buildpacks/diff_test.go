package buildpacks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/alphagov/paas-buildpack-notifier/buildpacks"
)

var _ = Describe("Diff", func() {
	var releases BuildpackReleases

	BeforeEach(func() {
		releases = BuildpackReleases{
			"test-v1.0.0": BuildpackRelease{
				Defaults: map[string]string{
					"modified": "v0.0.1",
				},
				Dependencies: map[string][]string{
					"modified":   []string{"v0.0.1", "v0.0.2"},
					"unmodified": []string{"v0.0.1"},
					"removed":    []string{"v0.0.1"},
				},
			},
			"test-v2.0.0": BuildpackRelease{
				Defaults: map[string]string{
					"modified": "v0.0.2",
				},
				Dependencies: map[string][]string{
					"modified":   []string{"v0.0.2", "v0.0.3"},
					"unmodified": []string{"v0.0.1"},
					"added":      []string{"v0.0.1"},
				},
			},
		}
	})

	It("should compare two different versions", func() {
		diff, err := DiffBuildpackVersions("test", "v1.0.0", "v2.0.0", releases)
		Expect(err).ToNot(HaveOccurred())
		Expect(diff.Changes()).To(BeTrue())
		Expect(diff).To(Equal(BuildpackDiff{
			Defaults: map[string]DefaultDiff{
				"modified": DefaultDiff{
					Added:   "v0.0.2",
					Removed: "v0.0.1",
				},
			},
			Added: map[string][]string{
				"modified": []string{"v0.0.3"},
				"added":    []string{"v0.0.1"},
			},
			Removed: map[string][]string{
				"modified": []string{"v0.0.1"},
				"removed":  []string{"v0.0.1"},
			},
			Overlap: map[string][]string{
				"modified":   []string{"v0.0.2"},
				"unmodified": []string{"v0.0.1"},
			},
		}))
	})

	It("should compare two identical versions", func() {
		diff, err := DiffBuildpackVersions("test", "v1.0.0", "v1.0.0", releases)
		Expect(err).ToNot(HaveOccurred())
		Expect(diff.Changes()).To(BeFalse())
		Expect(diff).To(Equal(BuildpackDiff{}))
	})
})
