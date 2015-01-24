package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {

	It("Default server port to 8080", func() {
		result := Repo()
		Eventually(result.Out).Should(Say(":8080"))
		result.Kill()
	})

	It("Default server address to 0.0.0.0", func() {
		result := Repo()
		Eventually(result.Out).Should(Say("0.0.0.0:"))
		result.Kill()
	})

	It("-p flag sets server port", func() {
		result := Repo("-p", "8888")
		Eventually(result.Out).Should(Say(":8888"))
		result.Kill()
	})

	It("-h flag sets server port", func() {
		result := Repo("-h", "127.0.0.1")
		Eventually(result.Out).Should(Say("127.0.0.1:"))
		result.Kill()
	})
})

func Repo(args ...string) *Session {
	path, err := Build("github.com/cloudfoundry-incubator/cli-plugin-repo/main")
	Expect(err).NotTo(HaveOccurred())

	session, err := Start(exec.Command(path, args...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
