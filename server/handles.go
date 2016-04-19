package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
)

type handlers struct {
	plugins models.PluginsJson
	logger  io.Writer
}

func NewServerHandles(plugins models.PluginsJson, logger io.Writer) *handlers {
	return &handlers{
		plugins: plugins,
		logger:  logger,
	}
}

func (h *handlers) ListPlugins(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(h.plugins)
	if err != nil {
		h.logger.Write([]byte("Error marshalling plugin models: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
