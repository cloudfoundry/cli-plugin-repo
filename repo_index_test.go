package main_test

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/web"

	"net/url"

	"crypto/sha1"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Database", func() {
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
			if os.Getenv("BINARY_VALIDATION") != "true" {
				Skip("Skipping SHA1 binary checking. To enable, set the BINARY_VALIDATION env variable to 'true'")
			}

			fmt.Println("\nRunning Binary Validations, this could take 10+ minutes")

			for _, plugin := range plugins.Plugins {
				for _, binary := range plugin.Binaries {
					resp, err := http.Get(binary.Url)
					Expect(err).NotTo(HaveOccurred())

					defer resp.Body.Close()
					b, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())

					s := sha1.Sum(b)
					Expect(hex.EncodeToString(s[:])).To(Equal(binary.Checksum), fmt.Sprintf("Plugin '%s' has an invalid checksum for platform '%s'", plugin.Name, binary.Platform))
				}
			}
		})
	})
})
