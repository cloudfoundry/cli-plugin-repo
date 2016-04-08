package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/cloudfoundry-incubator/cli-plugin-repo/models"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/parser"
	"github.com/cloudfoundry-incubator/cli-plugin-repo/server"
)

var flagPort int
var flagAddr string
var logger io.Writer

func Start() {
	logger = os.Stdout //turn this into a logger soon

	flag.StringVar(&flagAddr, "n", "0.0.0.0", "Address the server to listen on")
	flag.IntVar(&flagPort, "p", 8080, "Port the server to listen on")
	flag.Parse()

	if port := os.Getenv("PORT"); port != "" {

		v, err := strconv.Atoi(port)
		if err != nil {
			logger.Write([]byte("Error getting port from VCAP: " + err.Error()))
		} else {
			flagPort = v
			flagAddr = ""
		}
	}

	model := models.NewPlugins(logger)

	yamlParser := parser.NewYamlParser("repo-index.yml", logger, model)
	handles := server.NewServerHandles(yamlParser, logger)

	repoServer := server.NewRepoServer(flagPort, flagAddr, handles)

	fmt.Printf("\nServer is listening on %s:%d\n", flagAddr, flagPort)

	repoServer.Serve()
}
