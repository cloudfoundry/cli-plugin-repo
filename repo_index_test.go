package main_test

import (
	"io/ioutil"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/web"

	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"fmt"
)

var _ = FDescribe("Database", func() {
	It("correctly parses the current repo-index.yml", func() {
		var plugins web.PluginsJson

		b, err := ioutil.ReadFile("repo-index.yml")
		Expect(err).NotTo(HaveOccurred())

		err = yaml.Unmarshal(b, &plugins)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("validations", func() {
		var plugins web.PluginsJson

		BeforeEach(func() {
			b, err := ioutil.ReadFile("repo-index.yml")
			Expect(err).NotTo(HaveOccurred())

			err = yaml.Unmarshal(b, &plugins)
			Expect(err).NotTo(HaveOccurred())
		})

		It("has every binary link over https", func() {
			for _, plugin := range plugins.Plugins {
				for _, binary := range plugin.Binaries {
					url, err := url.Parse(binary.Url)
					Expect(err).NotTo(HaveOccurred())

					Expect(url.Scheme).To(Equal("https"))
				}
			}
		})

		It("has every version parseable by semver", func() {
			for _, plugin := range plugins.Plugins {
				Expect(plugin.Version).To(MatchRegexp(`^\d+\.\d+\.\d+$`), fmt.Sprintf("Plugin '%s' has a non-semver version", plugin.Name))
			}
		})

		It("every binary download had a matching sha1", func() {

		})
	})
})
