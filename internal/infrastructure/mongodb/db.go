package mongodb

import (
	"context"
	"time"

	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/config"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDB interface {
	GetDatabase() *mongo.Database
	CheckMongoDBHealth() bool
}

type mongoDB struct {
	database *mongo.Database
	config   domain.IConfig
	logger   domain.ILogger
}

func NewMongoDB(connectionString, dbName string) (MongoDB, error) {
	config := config.NewConfig()
	clientOptions := buildClientOptions(connectionString, config)
	client, err := connectToMongo(clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &mongoDB{
		database: database,
		config:   config,
		logger:   logging.NewLogger(),
	}, nil
}

func (db *mongoDB) GetDatabase() *mongo.Database {
	return db.database
}

func (db *mongoDB) CheckMongoDBHealth() bool {
	checkTimeout := time.Duration(db.config.GetDbHealthTimeout()) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	err := db.database.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		db.logger.Errorf("Error pinging MongoDB: %v", err)
		return false
	}
	return true
}

func buildClientOptions(connectionString string, config domain.IConfig) *options.ClientOptions {
	clientOptions := options.Client().ApplyURI(connectionString)

	clientOptions.SetConnectTimeout(time.Duration(config.GetDbConnectionTimeout()) * time.Second)
	clientOptions.SetServerSelectionTimeout(time.Duration(config.GetDbServerSelectionTimeout()) * time.Second)
	clientOptions.SetMaxPoolSize(uint64(config.GetDbMaxPoolSize()))
	clientOptions.SetMinPoolSize(uint64(config.GetDbMinPoolSize()))
	clientOptions.SetReadPreference(mapReadPreference(config.GetDbReadPreference()))
	writeConcern := writeconcern.Majority()
	writeConcern.WTimeout = time.Duration(config.GetDbWriteConcernTimeout()) * time.Millisecond
	clientOptions.SetWriteConcern(writeConcern)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetRetryWrites(config.GetDbRetryWrites())
	clientOptions.SetHeartbeatInterval(time.Duration(config.GetDbHeartbeatInterval()) * time.Second)

	return clientOptions
}

func connectToMongo(clientOptions *options.ClientOptions) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func mapReadPreference(preference string) *readpref.ReadPref {
	switch preference {
	case "primary":
		return readpref.Primary()
	case "secondary":
		return readpref.Secondary()
	case "nearest":
		return readpref.Nearest()
	default:
		return readpref.Primary()
	}
}
