package mappers

import (
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/mongodb/models"
)

type PluginMapper struct{}

func NewPluginMapper() *PluginMapper {
	return &PluginMapper{}
}

// Mapping from Plugin model in persistence to Domain Plugin
func (m *PluginMapper) ToDomain(pluginModel *models.PluginModel) *domain.Plugin {
	if pluginModel == nil {
		return nil
	}

	return &domain.Plugin{
		Id:   pluginModel.Id,
		Name: pluginModel.Name,
	}
}

// Map the Plugin of the domain to the Plugin model of the persistence
func (m *PluginMapper) ToPersistence(plugin *domain.Plugin) *models.PluginModel {
	if plugin == nil {
		return nil
	}

	return &models.PluginModel{
		Id:   plugin.Id,
		Name: plugin.Name,
	}
}

// func (m *PluginMapper) ToResponse(job *domain.Plugin) *domain.PluginResponse {
// 	// Implementar el mapeo de Plugin del dominio a PluginResponse
// }
