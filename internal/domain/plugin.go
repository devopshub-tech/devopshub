package domain

import (
	"time"

	"github.com/devopshub-tech/devopshub/pkg/utils"
)

// Plugin represents a plugin to be used in a job.
type Plugin struct {
	Id             string
	Name           string
	CurrentVersion string // Referencia al ID de la versión actual en Versions collection
}

// PluginVersion represents a specific version of a plugin.
type PluginVersion struct {
	Id         string
	PluginID   string // Reference to the plugin ID in the Plugins entity
	Version    string
	Dockerfile string
	ActionYAML string
	CreatedAt  time.Time
}

// PluginInfo provides information about a plugin.
type PluginInfo struct {
	Name    string
	Version string
}

// NewPlugin creates a new Plugin instance.
func NewPlugin(name string) *Plugin {
	return &Plugin{
		Id:   utils.GenerateUUID(),
		Name: name,
	}
}

// PluginRepository defines the interface to interact with the plugin repository.
type IPluginRepository interface {
	Create(plugin *Plugin) error
	Update(plugin *Plugin) error
	Find(filter interface{}) ([]*Plugin, error)
	FindById(id string) (*Plugin, error)
}

// PluginVersionRepository defines the interface to interact with the plugin version repository.
type IPluginVersionRepository interface {
	Create(pluginVersion *PluginVersion) error
	Update(pluginVersion *PluginVersion) error
	Find(filter interface{}) ([]*PluginVersion, error)
	FindById(id string) (*PluginVersion, error)
}
