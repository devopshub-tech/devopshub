package mongodb

import (
	"context"
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PluginVersionRepositoryMongo struct {
	collection *mongo.Collection
}

func NewPluginVersionRepositoryMongo(db *mongo.Database) (domain.IPluginVersionRepository, error) {
	collection := db.Collection(PluginVersionsCollectionName)
	return &PluginVersionRepositoryMongo{
		collection: collection,
	}, nil
}

func (r *PluginVersionRepositoryMongo) Create(version *domain.PluginVersion) error {
	version.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.Background(), version)
	return err
}

func (r *PluginVersionRepositoryMongo) Update(version *domain.PluginVersion) error {
	filter := bson.M{"_id": version.Id}
	update := bson.M{"$set": bson.M{
		"dockerfile": version.Dockerfile,
		"actionyaml": version.ActionYAML,
	}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *PluginVersionRepositoryMongo) FindById(id string) (*domain.PluginVersion, error) {
	var version domain.PluginVersion
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *PluginVersionRepositoryMongo) Find(filter interface{}) ([]*domain.PluginVersion, error) {
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var versions []*domain.PluginVersion
	for cursor.Next(context.Background()) {
		var version domain.PluginVersion
		if err := cursor.Decode(&version); err != nil {
			return nil, err
		}
		versions = append(versions, &version)
	}

	return versions, nil
}
