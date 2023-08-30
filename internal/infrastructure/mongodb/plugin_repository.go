// internal/infrastructure/db/mongodb.go
package mongodb

import (
	"context"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PluginRepositoryMongo struct {
	collection *mongo.Collection
}

func NewPluginRepositoryMongo(db *mongo.Database) (domain.IPluginRepository, error) {
	collection := db.Collection(PluginsCollectionName)
	return &PluginRepositoryMongo{
		collection: collection,
	}, nil
}

func (r *PluginRepositoryMongo) Create(plugin *domain.Plugin) error {
	_, err := r.collection.InsertOne(context.Background(), plugin)
	if err != nil {
		return err
	}

	return nil
}

func (r *PluginRepositoryMongo) Update(plugin *domain.Plugin) error {
	filter := bson.M{"_id": plugin.Id}
	update := bson.M{"$set": bson.M{
		"name": plugin.Name,
	}}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *PluginRepositoryMongo) FindById(id string) (*domain.Plugin, error) {
	var plugin domain.Plugin
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(context.Background(), filter).Decode(&plugin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Plugin not found
		}
		return nil, err
	}

	return &plugin, nil
}

func (r *PluginRepositoryMongo) Find(filter interface{}) ([]*domain.Plugin, error) {
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var plugins []*domain.Plugin
	for cursor.Next(context.Background()) {
		var plugin domain.Plugin
		if err := cursor.Decode(&plugin); err != nil {
			return nil, err
		}
		plugins = append(plugins, &plugin)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return plugins, nil
}
