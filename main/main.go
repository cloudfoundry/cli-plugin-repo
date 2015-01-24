package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/parser"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/server"
)

var flagPort int
var flagAddr string

func init() {
	flag.StringVar(&flagAddr, "h", "0.0.0.0", "Address the server to listen on")
	flag.IntVar(&flagPort, "p", 8080, "Port the server to listen on")
	flag.Parse()
}

func main() {

	model := models.NewPlugins(os.Stdout)

	yamlParser := parser.NewYamlParser("repo-index.yml", os.Stdout, model)
	handles := server.NewServerHandles(yamlParser, os.Stdout)

	repoServer := server.NewRepoServer(flagPort, flagAddr, handles)
	fmt.Printf("\nServer is listening on %s:%d\n", flagAddr, flagPort)

	repoServer.Serve()
}
