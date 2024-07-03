package main_test

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"code.cloudfoundry.org/cli-plugin-repo/sort/yamlsorter"
	"code.cloudfoundry.org/cli-plugin-repo/web"
	"github.com/blang/semver"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

// NamesToSkip provides a list of plugins that were created prior to the naming
// rules being established. No new plugins should be added to this list.
var NamesToSkip = []string{
	"apigee-broker-plugin",
	"app-autoscaler-plugin",
	"Buildpack Management",
	"Buildpack Usage",
	"cf-aklogin",
	"cf-icd-plugin",
	"cf-predix-analytics-plugin",
	"Cloud Deployment Plugin",
	"Copy Autoscaler",
	"Copy Env",
	"doctor",
	"Download Droplet",
	"fastpush",
	"Firehose Plugin",
	"mysql-plugin",
	"Scaleover",
	"Service Instance Logging",
	"spring-cloud-dataflow-for-pcf",
	"Targets",
	"Usage Report",
	"whoami-plugin",
	"wildcard_plugin",
}

func ShouldValidatePluginName(pluginName string) bool {
	for _, pluginToSkip := range NamesToSkip {
		if pluginName == pluginToSkip {
			return false
		}
	}
	return true
}

func GetNameFromBinary(pluginrunner string, b []byte) string {
	// Write bytes to disk and make it runnable
	tmpfile, err := ioutil.TempFile("", "plugin-meta")
	Expect(err).ToNot(HaveOccurred())
	Expect(tmpfile.Close()).ToNot(HaveOccurred())
	pathToPlugin := tmpfile.Name()
	// defer os.Remove(pathToPlugin)
	err = ioutil.WriteFile(pathToPlugin, b, 0777)
	Expect(err).ToNot(HaveOccurred())

	// Required to chmod the tempfile
	err = os.Chmod(pathToPlugin, 0777)
	Expect(err).ToNot(HaveOccurred())

	cmd := exec.Command(pluginrunner, pathToPlugin)
	outBuff := new(bytes.Buffer)
	session, err := Start(cmd, outBuff, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(Exit(0))

	nameBuffer, err := ioutil.ReadAll(outBuff)
	Expect(err).ToNot(HaveOccurred())

	return strings.TrimSpace(string(nameBuffer))
}

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
		var pluginBytes []byte

		BeforeEach(func() {
			var err error
			pluginBytes, err = ioutil.ReadFile("repo-index.yml")
			Expect(err).NotTo(HaveOccurred())

			err = yaml.Unmarshal(pluginBytes, &plugins)
			Expect(err).NotTo(HaveOccurred())
		})

		It("the yaml file is sorted", func() {
			var yamlSorter yamlsorter.YAMLSorter

			sortedBytes, err := yamlSorter.Sort(pluginBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(sortedBytes).To(Equal(pluginBytes), "file is not sorted; please run 'go run sort/main.go repo-index.yml'.\n")
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
				version, err := semver.Make(plugin.Version)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Plugin '%s' has a non-semver version", plugin.Name))
				Expect(version.Validate()).ToNot(HaveOccurred(), fmt.Sprintf("Plugin '%s' has a non-semver version", plugin.Name))
			}
		})

		It("validates the platforms for every binary", func() {
			for _, plugin := range plugins.Plugins {
				for _, binary := range plugin.Binaries {
					Expect(web.ValidPlatforms).To(
						ContainElement(binary.Platform),
						fmt.Sprintf(
							"Plugin '%s' contains a platform '%s' that is invalid. Please use one of the following: '%s'",
							plugin.Name,
							binary.Platform,
							strings.Join(web.ValidPlatforms, ", "),
						))
				}
			}
		})

		It("requires HTTPS for all downloads", func() {
			for _, plugin := range plugins.Plugins {
				for _, binary := range plugin.Binaries {
					Expect(binary.Url).To(
						MatchRegexp("^https|ftps"),
						fmt.Sprintf(
							"Plugin '%s' links to a Binary's URL '%s' that cannot be downloaded over SSL (begins with https/ftps). Please provide a secure download link to your binaries. If you are unsure how to provide one, try out GitHub Releases: https://help.github.com/articles/creating-releases",
							plugin.Name,
							binary.Url,
						))
				}
			}
		})

		It("every binary download had a matching sha1", func() {
			if os.Getenv("BINARY_VALIDATION") != "true" {
				Skip("Skipping SHA1 binary checking. To enable, set the BINARY_VALIDATION env variable to 'true'")
			}

			// Binary will be cleaned up in AfterSuite
			runnerPath, err := Build("code.cloudfoundry.org/cli-plugin-repo/pluginrunner")
			Expect(err).ToNot(HaveOccurred())

			fmt.Println("\nRunning Binary Validations, this could take 10+ minutes")

			for _, plugin := range plugins.Plugins {
				for _, binary := range plugin.Binaries {
					fmt.Printf("%s checking %s %s\n", time.Now().Format(time.RFC3339), plugin.Name, binary.Platform)
					Expect(binary.Url).To(ContainSubstring(plugin.Version), fmt.Sprintf("%s's URL must be versioned - %s is missing version number", plugin.Name, binary.Url))

					var err error
					resp, err := http.Get(binary.Url)
					Expect(err).NotTo(HaveOccurred())

					// If there's a network error, retry exactly once for this plugin binary.
					switch resp.StatusCode {
					case http.StatusInternalServerError,
						http.StatusBadGateway,
						http.StatusServiceUnavailable,
						http.StatusGatewayTimeout:
						Expect(resp.Body.Close()).To(Succeed())
						resp, err = http.Get(binary.Url)
						Expect(err).NotTo(HaveOccurred())
					}

					defer resp.Body.Close()
					b, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(resp.StatusCode).To(And(BeNumerically(">=", 200), BeNumerically("<", 400)),
						"Failed to retrieve '%s', can't compute SHA from URL %s", plugin.Name, binary.Url)

					s := sha1.Sum(b)
					Expect(hex.EncodeToString(s[:])).To(Equal(binary.Checksum),
						fmt.Sprintf("Plugin '%s' has an invalid checksum for platform '%s'", plugin.Name, binary.Platform))

					if binary.Platform == "linux64" && ShouldValidatePluginName(plugin.Name) {
						binaryName := GetNameFromBinary(runnerPath, b)
						Expect(binaryName).To(Equal(plugin.Name),
							fmt.Sprintf("The plugin name provided by in the 'repo-index.yml' must match the plugin name returned from the plugin binary when 'SendMetadata' is called.\nBinary Name: %s\nYAML Name: %s\n", binaryName, plugin.Name))
					}
				}
			}
		})
	})
})
