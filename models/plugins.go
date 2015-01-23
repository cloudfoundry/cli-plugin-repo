package models

import "io"

type PluginModel interface {
	PopulateModel(interface{})
	PluginsModel() []Plugin
}

type Plugins struct {
	Plugins []Plugin
	logger  io.Writer
}

type Plugin struct {
	Name        string
	Description string
	Binaries    []Binary
}

type Binary struct {
	Platform string
	Url      string
	Checksum string
}

func NewPlugins(logger io.Writer) PluginModel {
	return &Plugins{
		logger: logger,
	}
}

func (p *Plugins) PluginsModel() []Plugin {
	return p.Plugins
}

func (p *Plugins) PopulateModel(input interface{}) {
	if contents, ok := input.(map[interface{}]interface{})["plugins"].([]interface{}); ok {
		for _, plugin := range contents {
			p.Plugins = append(p.Plugins, p.extractPlugin(plugin))
		}
	} else {
		p.logger.Write([]byte("unexpected yaml structure, 'plugins' field not found."))
	}
}

func (p *Plugins) extractPlugin(rawData interface{}) Plugin {
	plugin := Plugin{}
	for k, v := range rawData.(map[interface{}]interface{}) {
		switch k.(string) {
		case "name":
			plugin.Name = v.(string)
		case "description":
			plugin.Description = v.(string)
		case "binaries":
			for _, binary := range v.([]interface{}) {
				plugin.Binaries = append(plugin.Binaries, p.extractBinaries(binary))
			}
		default:
			p.logger.Write([]byte("unexpected field in plugins: " + k.(string)))
		}
	}
	return plugin
}

func (p *Plugins) extractBinaries(input interface{}) Binary {
	binary := Binary{}
	for k, v := range input.(map[interface{}]interface{}) {
		switch k.(string) {
		case "platform":
			binary.Platform = v.(string)
		case "url":
			binary.Url = v.(string)
		case "checksum":
			binary.Checksum = v.(string)
		default:
			p.logger.Write([]byte("unexpected field in binaries: " + k.(string)))
		}
	}
	return binary
}
