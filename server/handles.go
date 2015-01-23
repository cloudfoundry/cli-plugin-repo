package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/parser"
)

type ServerHandles interface {
	ListPlugins(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	yamlParser parser.YamlParser
	logger     io.Writer
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
		panic("Error parsing repo-index.yml: " + err.Error())
	}

	b, err := json.Marshal(pluginModel)
	if err != nil {
		h.logger.Write([]byte("Error marshalling plugin models: " + err.Error()))
	}

	w.Write(b)
}
