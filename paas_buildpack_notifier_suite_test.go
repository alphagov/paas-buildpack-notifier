package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPaasBuildpackNotifier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PaasBuildpackNotifier Suite")
}
