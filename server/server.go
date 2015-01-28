package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RepoServer interface {
	Serve()
	Stop()
	Port() string
}

type server struct {
	port     int
	listener net.Listener
	handles  ServerHandles
}

func NewRepoServer(port int, addr string, handles ServerHandles) RepoServer {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}

	return server{
		port:     port,
		listener: listener,
		handles:  handles,
	}

}

func (s server) registerHandlers() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/list").HandlerFunc(s.handles.ListPlugins)
	r.Methods("GET").Path("/").HandlerFunc(redirectBase)

	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
	http.Handle("/", r)
}

func redirectBase(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusFound)
}

func (s server) Serve() {
	s.registerHandlers()
	http.Serve(s.listener, nil)
}

func (s server) Stop() {
	s.listener.Close()
}

func (s server) Port() string {
	return strconv.Itoa(s.listener.Addr().(*net.TCPAddr).Port)
}
