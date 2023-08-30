package models

import (
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
)

// JobModel represents the MongoDB model for the Job entity.
type JobModel struct {
	Id        string            `bson:"_id,omitempty"`
	Plugin    PluginInfoModel   `bson:"plugin,omitempty"`
	Status    domain.JobStatus  `bson:"status,omitempty"`
	Log       string            `bson:"log,omitempty"`
	Inputs    map[string]string `bson:"inputs,omitempty"`
	CreatedAt time.Time         `bson:"createdAt,omitempty"`
	UpdatedAt time.Time         `bson:"updatedAt,omitempty"`
}

// PluginInfoModel represents the MongoDB model for the PluginInfo entity.
type PluginInfoModel struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
}
