package main

import (
	"fmt"
	"os"

	"net/http"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/server"
	"github.com/jessevdk/go-flags"
	"github.com/tedsuo/rata"
	"gopkg.in/yaml.v2"
	"sort"
)

type CLIPR struct {
	Port          int    `short:"p" long:"port" default:"8080" description:"Port that the plugin repo listens on"`
	RepoIndexPath string `short:"f" long:"filepath" default:"repo-index.yml" description:"Path to repo-index file"`
}

func main() {
	cmd := &CLIPR{}

	parser := flags.NewParser(cmd, flags.Default)
	parser.NamespaceDelimiter = "-"

	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	err = cmd.Execute(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (cmd *CLIPR) Execute(args []string) error {
	logger := os.Stdout //turn this into a logger soon

	var plugins models.PluginsJson

	b, err := ioutil.ReadFile(cmd.RepoIndexPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &plugins)
	if err != nil {
		return err
	}

	sort.Sort(plugins)
	handles := server.NewServerHandles(plugins, logger)

	handlers := map[string]http.Handler{
		Index: http.FileServer(http.Dir("ui")),
		List:  http.HandlerFunc(handles.ListPlugins),
	}

	router, err := rata.NewRouter(Routes, handlers)

	if err != nil {
		return err
	}

	err = http.ListenAndServe(cmd.bindAddr(), router)

	return err
}

func (cmd *CLIPR) bindAddr() string {
	return fmt.Sprintf(":%d", cmd.Port)
}

const (
	Index = "Index"
	List  = "List"
)

var Routes = rata.Routes([]rata.Route{
	{Path: "/", Method: "GET", Name: Index},
	{Path: "/js/:file", Method: "GET", Name: Index},
	{Path: "/css/:file", Method: "GET", Name: Index},
	{Path: "/font/:file", Method: "GET", Name: Index},
	{Path: "/images/:file", Method: "GET", Name: Index},
	{Path: "/list", Method: "GET", Name: List},
})
