package buildpacks_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/alphagov/paas-buildpack-notifier/buildpacks"
)

var _ = Describe("Releases", func() {
	Describe("BuildpackReleases", func() {
		const (
			rubyName    = "ruby"
			rubyVersion = "v3.3.3"
		)
		var (
			releases    BuildpackReleases
			rubyRelease BuildpackRelease
		)

		BeforeEach(func() {
			releases = BuildpackReleases{}
			rubyRelease = BuildpackRelease{
				Defaults:     map[string]string{"ruby": "v1.1.1"},
				Dependencies: map[string][]string{"ruby": {"v1.1.1", "v2.2.2"}},
			}
		})

		It("supports checking and adding to set", func() {
			By("checking non-existant release", func() {
				Expect(releases.Has(rubyName, rubyVersion)).To(Equal(false))

				release, err := releases.Get(rubyName, rubyVersion)
				Expect(err).To(MatchError("unable to find release for ruby-v3.3.3"))
				Expect(release).To(Equal(BuildpackRelease{}))
			})

			By("adding release", func() {
				releases.Add(rubyName, rubyVersion, rubyRelease)
			})

			By("checking newly added release", func() {
				Expect(releases.Has(rubyName, rubyVersion)).To(Equal(true))

				release, err := releases.Get(rubyName, rubyVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(rubyRelease))
			})
		})
	})

	Describe("FetchBuildpackRelease", func() {
		var (
			server      *ghttp.Server
			manifestURL string
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			manifestURL = server.URL() + "/manifest.yml"
		})

		Describe("reading fixture", func() {
			BeforeEach(func() {
				response, err := ioutil.ReadFile("fixtures/manifest.ruby.yml")
				Expect(err).ToNot(HaveOccurred())

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/manifest.yml"),
						ghttp.RespondWith(http.StatusOK, response),
					),
				)
			})

			It("should parse buildpack manifest", func() {
				release, err := FetchBuildpackRelease(manifestURL)
				Expect(err).ToNot(HaveOccurred())
				Expect(release).To(Equal(BuildpackRelease{
					Defaults: map[string]string{
						"ruby": "2.4.x",
					},
					Dependencies: map[string][]string{
						"bundler":           []string{"1.16.1"},
						"jruby":             []string{"ruby-1.9.3-jruby-1.7.26", "ruby-2.0.0-jruby-1.7.26", "ruby-2.3.3-jruby-9.1.16.0"},
						"node":              []string{"6.14.1", "4.9.1"},
						"openjdk1.8-latest": []string{"1.8.0"},
						"ruby":              []string{"2.2.8", "2.2.9", "2.3.5", "2.3.6", "2.4.2", "2.4.3", "2.5.0", "2.5.1"},
						"rubygems":          []string{"2.7.6"},
						"yarn":              []string{"1.5.1"},
					},
				}))
			})
		})

		Describe("failure", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/manifest.yml"),
						ghttp.RespondWith(http.StatusServiceUnavailable, ""),
					),
				)
			})

			It("should parse buildpack manifest", func() {
				release, err := FetchBuildpackRelease(manifestURL)
				Expect(err).To(MatchError(
					fmt.Sprintf("error getting %s: 503 Service Unavailable", manifestURL),
				))
				Expect(release).To(Equal(BuildpackRelease{}))
			})
		})
	})
})
