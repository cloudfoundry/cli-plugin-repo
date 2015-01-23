package parser_test

import (
	"github.com/cloudfoundry-incubator/cli-plugin-repo/models/fakes"
	. "github.com/cloudfoundry-incubator/cli-plugin-repo/parser"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {

	var (
		yparser   YamlParser
		logger    *test_helpers.TestLogger
		fakeModel *fakes.FakePluginModel
	)

	BeforeEach(func() {
		logger = test_helpers.NewTestLogger()
		fakeModel = &fakes.FakePluginModel{}
	})

	Describe("Parse()", func() {

		It("logs error if file does not exist", func() {
			yparser = NewYamlParser("../path/to/nowhere/bad.yml", logger, fakeModel)
			_, err := yparser.Parse()
			Ω(err).To(HaveOccurred())
			Ω(logger.ContainsSubstring([]string{"File does not exist"})).To(Equal(true))
		})

		It("logs error if file is not a valid yaml", func() {
			yparser = NewYamlParser("../fixtures/parser/bad.yml", logger, fakeModel)
			_, err := yparser.Parse()
			Ω(err).To(HaveOccurred())
			Ω(logger.ContainsSubstring([]string{"Failed to decode document"})).To(Equal(true))
		})

		It("validate a yaml file", func() {
			yparser = NewYamlParser("../fixtures/parser/test.yml", logger, fakeModel)
			_, err := yparser.Parse()
			Ω(err).ToNot(HaveOccurred())
		})

		It("calls models.PopulateModel", func() {
			yparser = NewYamlParser("../fixtures/parser/test.yml", logger, fakeModel)
			_, err := yparser.Parse()
			Ω(err).ToNot(HaveOccurred())

			Ω(fakeModel.PopulateModelCallCount()).To(Equal(1))
		})

		It("passes parsed yml raw data to models.PluginModel", func() {
			yparser = NewYamlParser("../fixtures/parser/test.yml", logger, fakeModel)
			_, err := yparser.Parse()
			Ω(err).ToNot(HaveOccurred())

			rawData := fakeModel.PopulateModelArgsForCall(0)
			Ω(rawData.(map[interface{}]interface{})["plugins"]).ShouldNot(Equal(nil))
			Ω(rawData.(map[interface{}]interface{})["plugins"].([]interface{})[0].(map[interface{}]interface{})["name"]).Should(Equal("plugin1"))
			Ω(rawData.(map[interface{}]interface{})["plugins"].([]interface{})[1].(map[interface{}]interface{})["name"]).Should(Equal("plugin2"))
		})
	})
})
