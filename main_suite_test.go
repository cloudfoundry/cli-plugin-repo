package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var buildPath string
var _ = SynchronizedBeforeSuite(func() []byte {
	path, buildErr := gexec.Build("code.cloudfoundry.org/cli-plugin-repo")
	Expect(buildErr).NotTo(HaveOccurred())
	return []byte(path)
}, func(data []byte) {
	buildPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})
