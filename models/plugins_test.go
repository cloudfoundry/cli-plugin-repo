package models_test

import (
	"os"

	. "github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	var (
		parsedYaml  interface{}
		pluginModel PluginModel
	)

	Context("When raw data is valid", func() {
		BeforeEach(func() {
			parsedYaml = map[interface{}]interface{}{
				"plugins": []interface{}{
					map[interface{}]interface{}{
						"name":        "test1",
						"description": "n/a",
						"binaries": []interface{}{
							map[interface{}]interface{}{
								"platform": "osx",
								"url":      "example.com/plugin",
								"checksum": "abcdefg",
							},
						},
					},
					map[interface{}]interface{}{
						"name":        "test2",
						"description": "n/a",
						"binaries": []interface{}{
							map[interface{}]interface{}{
								"platform": "windows",
								"url":      "example.com/plugin",
								"checksum": "abcdefg",
							},
							map[interface{}]interface{}{
								"platform": "linux32",
								"url":      "example.com/plugin",
								"checksum": "abcdefg",
							},
						},
					},
				},
			}

			pluginModel = NewPlugins(os.Stdout)
			pluginModel.PopulateModel(parsedYaml)
		})

		It("populates the plugin model with raw data", func() {
			data := pluginModel.PluginsModel()
			Ω(len(data.Plugins)).To(Equal(2))
			Ω(data.Plugins[0].Name).To(Equal("test1"))
			Ω(data.Plugins[0].Binaries[0].Platform).To(Equal("osx"))
			Ω(data.Plugins[1].Name).To(Equal("test2"))
			Ω(data.Plugins[1].Binaries[1].Platform).To(Equal("linux32"))
		})
	})

	Context("When raw data contains unknown field", func() {
		var (
			logger *test_helpers.TestLogger
		)

		BeforeEach(func() {
			parsedYaml = map[interface{}]interface{}{
				"plugins": []interface{}{
					map[interface{}]interface{}{
						"name":          "test1",
						"description":   "n/a",
						"unknown_field": "123",
					},
				},
			}

			logger = test_helpers.NewTestLogger()
			pluginModel = NewPlugins(logger)
			pluginModel.PopulateModel(parsedYaml)
		})

		It("logs error to terminal", func() {
			data := pluginModel.PluginsModel()
			Ω(len(data.Plugins)).To(Equal(1))
			Ω(logger.ContainsSubstring([]string{"unexpected field", "unknown_field"})).To(Equal(true))
		})
	})

})
