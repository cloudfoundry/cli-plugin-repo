package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"code.cloudfoundry.org/cli-plugin-repo/web"
	"github.com/unrolled/secure"
	yaml "gopkg.in/yaml.v2"
)

type CLIPR struct {
	Port          int    `short:"p" long:"port" default:"8080" description:"Port that the plugin repo listens on"`
	RepoIndexPath string `short:"f" long:"filepath" default:"repo-index.yml" description:"Path to repo-index file"`
	ForceSSL      bool   `long:"force-ssl"  description:"Force SSL on every request"`
}

func (cmd *CLIPR) Execute(args []string) error {
	logger := os.Stdout //turn this into a logger soon

	var plugins web.PluginsJson

	b, err := ioutil.ReadFile(cmd.RepoIndexPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &plugins)
	if err != nil {
		return err
	}

	sort.Sort(plugins)

	tmpl, err := template.ParseFiles(filepath.Join("ui", "index.html"))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("ui"))
	mux.Handle("/images/", staticHandler)
	mux.Handle("/css/", staticHandler)
	mux.Handle("/font/", staticHandler)
	mux.Handle("/favicon.ico", http.RedirectHandler("/images/favicon.png", http.StatusFound))
	mux.Handle("/list", web.NewListPluginsHandler(plugins, logger)) // we no longer use for rendering, but keeping around in case others do
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, plugins)
		if err != nil { // should only error if template has syntax errors
			log.Println(err)
		}
	})

	var router http.Handler
	router = mux

	if cmd.ForceSSL {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect:     true,
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		})
		router = secureMiddleware.Handler(router)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", cmd.Port), router)
}
