package app_test

import (
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	var (
		path string
		err  error
	)

	BeforeSuite(func() {
		path, err = gexec.Build("github.com/cloudfoundry-incubator/cli-plugin-repo/")
		Expect(err).NotTo(HaveOccurred())
	})

	Repo := func(args ...string) *gexec.Session {
		session, err := gexec.Start(exec.Command(path, args...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		return session
	}

	RepoWithPORT := func(port string) *gexec.Session {
		cmd := exec.Command(path)
		cmd.Env = []string{"PORT=" + port}
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		return session
	}

	It("uses env var 'PORT' as listening port if found", func() {
		result := RepoWithPORT("12212")
		Eventually(result.Out).Should(gbytes.Say(":12212"))
		result.Kill()
	})

	It("Default server port to 8080 if env var 'PORT' is not found", func() {
		result := Repo()
		Eventually(result.Out).Should(gbytes.Say(":8080"))
		result.Kill()
	})

	It("Default server address to 0.0.0.0", func() {
		result := Repo()
		Eventually(result.Out).Should(gbytes.Say("0.0.0.0:"))
		result.Kill()
	})

	It("-p flag sets server port", func() {
		result := Repo("-p", "8888")
		Eventually(result.Out).Should(gbytes.Say(":8888"))
		result.Kill()
	})

	It("-n flag sets server port", func() {
		result := Repo("-n", "127.0.0.1")
		Eventually(result.Out).Should(gbytes.Say("127.0.0.1:"))
		result.Kill()
	})
})
