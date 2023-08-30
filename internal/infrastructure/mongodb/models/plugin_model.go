package models

// PluginModel represents the MongoDB model for the Plugin entity.
type PluginModel struct {
	Id             string `bson:"_id,omitempty"`
	Name           string `bson:"name,omitempty"`
	CurrentVersion string `bson:"currentVersion,omitempty"`
}
