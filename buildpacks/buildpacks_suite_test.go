package buildpacks_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildpacks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildpacks Suite")
}
