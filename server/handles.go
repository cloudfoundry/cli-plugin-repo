package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/parser"
)

type ServerHandles interface {
	ListPlugins(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	yamlParser parser.YamlParser
	logger     io.Writer
}

type JsonPluginList struct {
	Plugins []models.Plugin `json:"plugins"`
}

func NewServerHandles(yamlParser parser.YamlParser, logger io.Writer) ServerHandles {
	return &handlers{
		yamlParser: yamlParser,
		logger:     logger,
	}
}

func (h *handlers) ListPlugins(w http.ResponseWriter, r *http.Request) {
	pluginModel, err := h.yamlParser.Parse()
	if err != nil {
		h.logger.Write([]byte("Error parsing repo-index.yml: " + err.Error()))
		return
	}

	b, err := json.Marshal(JsonPluginList{
		Plugins: pluginModel,
	})
	if err != nil {
		h.logger.Write([]byte("Error marshalling plugin models: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
