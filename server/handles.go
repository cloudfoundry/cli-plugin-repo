package server

import "net/http"

type ServerHandles interface {
	ListPlugins(w http.ResponseWriter, r *http.Request)
}
