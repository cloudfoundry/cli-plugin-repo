package main

type CliRpcCmd struct {
	PluginMetadata *PluginMetadata
}

type PluginMetadata struct {
	Name string
}

func NewRPCCommand() *CliRpcCmd {
	return &CliRpcCmd{
		PluginMetadata: &PluginMetadata{},
	}
}

func (rpcCmd *CliRpcCmd) SetPluginMetadata(pluginMetadata PluginMetadata, retVal *bool) error {
	rpcCmd.PluginMetadata = &pluginMetadata
	*retVal = true
	return nil
}
