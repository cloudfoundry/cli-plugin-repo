package main

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
)

func main() {
	pluginLocation := strings.TrimSpace(os.Args[1])

	rpcServer, err := NewRPCServer(rpc.NewServer(), NewRPCCommand())
	if err != nil {
		panic(err)
	}

	pluginName, err := GetMetadata(rpcServer, pluginLocation)
	if err != nil {
		panic(err)
	}
	fmt.Println(pluginName)
}

func GetMetadata(rpcServer *RPCServer, path string) (string, error) {
	err := Run(rpcServer, path, "SendMetadata")
	if err != nil {
		return "", err
	}

	metadata := rpcServer.RPCCmd.PluginMetadata
	return metadata.Name, nil
}

func Run(rpcServer *RPCServer, path string, command string) error {
	err := rpcServer.Start()
	if err != nil {
		return err
	}
	defer rpcServer.Stop()

	cmd := exec.Command(path, rpcServer.Port(), command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
