package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/parser/fakes"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/test_helpers"

	. "github.com/cloudfoundry-incubator/cli-plugin-repo/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("handles", func() {

	var (
		fakeParser *fakes.FakeYamlParser
		testLogger *test_helpers.TestLogger
		resp       *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		fakeParser = &fakes.FakeYamlParser{}
		testLogger = test_helpers.NewTestLogger()
		resp = httptest.NewRecorder()
	})

	Describe("ListPlugins()", func() {
		It("logs any error returns by the parser, and returns empty body", func() {
			fakeParser.ParseReturns(models.Plugins{}, errors.New("bad yaml file"))
			h := NewServerHandles(fakeParser, testLogger)
			h.ListPlugins(resp, &http.Request{})

			Ω(testLogger.ContainsSubstring([]string{"Error parsing repo-index", "bad yaml file"})).To(Equal(true))
			Ω(resp.Body.Len()).To(Equal(0))
		})

		It("marshals PluginModels into json object", func() {
			model := models.Plugins{
				Plugins: []models.Plugin{
					models.Plugin{
						Name:        "plugin1",
						Description: "none",
						Binaries: []models.Binary{
							models.Binary{
								Platform: "osx",
								Url:      "asdf123",
							},
						},
					},
				},
			}

			fakeParser.ParseReturns(model, nil)
			h := NewServerHandles(fakeParser, testLogger)
			h.ListPlugins(resp, &http.Request{})

			var respondedModel models.Plugins
			err := json.Unmarshal(resp.Body.Bytes(), &respondedModel)
			Ω(err).ToNot(HaveOccurred())
			Ω(respondedModel.Plugins[0].Name).To(Equal("plugin1"))
			Ω(respondedModel.Plugins[0].Description).To(Equal("none"))
			Ω(respondedModel.Plugins[0].Binaries[0].Platform).To(Equal("osx"))
			Ω(respondedModel.Plugins[0].Binaries[0].Url).To(Equal("asdf123"))
		})
	})

})
