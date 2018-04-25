package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type RPCServer struct {
	listener net.Listener
	stopCh   chan struct{}
	Server   *rpc.Server
	RPCCmd   *CliRpcCmd
}

func NewRPCServer(providedServer *rpc.Server, providedCmd *CliRpcCmd) (*RPCServer, error) {
	rpcServer := &RPCServer{
		Server: providedServer,
		RPCCmd: providedCmd,
	}
	err := rpcServer.Server.Register(rpcServer.RPCCmd)
	if err != nil {
		return nil, err
	}

	return rpcServer, nil
}

func (rpcServer *RPCServer) Stop() {
	close(rpcServer.stopCh)
	rpcServer.listener.Close()
}

func (rpcServer *RPCServer) Port() string {
	return strconv.Itoa(rpcServer.listener.Addr().(*net.TCPAddr).Port)
}

func (rpcServer *RPCServer) Start() error {
	var err error

	rpcServer.stopCh = make(chan struct{})

	rpcServer.listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := rpcServer.listener.Accept()
			if err != nil {
				select {
				case <-rpcServer.stopCh:
					return
				default:
					fmt.Println(err)
				}
			} else {
				go rpcServer.Server.ServeConn(conn)
			}
		}
	}()

	return nil
}
