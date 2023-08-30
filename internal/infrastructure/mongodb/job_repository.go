// internal/infrastructure/db/job_repo.go
package mongodb

import (
	"context"
	"time"

	"github.com/devopshub-tech/devopshub/internal/application/mappers"
	"github.com/devopshub-tech/devopshub/internal/domain"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/logging"
	"github.com/devopshub-tech/devopshub/internal/infrastructure/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobRepositoryMongo struct {
	collection *mongo.Collection
	mapper     *mappers.JobMapper
	logger     domain.ILogger
}

// NewJobRepository creates a new instance of job repository.
func NewJobRepository(db *mongo.Database, mapper *mappers.JobMapper) (domain.IJobRepository, error) {
	collection := db.Collection(JobCollectionName)

	return &JobRepositoryMongo{
		collection: collection,
		mapper:     mapper,
		logger:     logging.NewLogger(),
	}, nil
}

// Create creates a new document in the collection.
func (r *JobRepositoryMongo) Create(job *domain.Job) (interface{}, error) {
	jobModel := r.mapper.ToPersistence(job)
	result, err := r.collection.InsertOne(context.Background(), jobModel)
	if err != nil {
		r.logger.Errorf("Error creating job in the repository: %v", err)
		return "", err
	}

	r.logger.Debugf("Job created in the repository. ID: %s", jobModel.Id)
	return result, nil
}

// Find searches for documents that match the filter in MongoDB.
func (r *JobRepositoryMongo) Find(filter interface{}) ([]*domain.Job, error) {
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		r.logger.Errorf("Error finding jobs in the repository: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var jobs []*domain.Job
	for cursor.Next(context.Background()) {
		var job domain.Job
		if err := cursor.Decode(&job); err != nil {
			r.logger.Errorf("Error decoding job from cursor: %v", err)
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Errorf("Cursor error while finding jobs: %v", err)
		return nil, err
	}

	r.logger.Debugf("Found %d jobs in the repository", len(jobs))
	return jobs, nil
}

// FindByID looks for a job by its ID in the collection.
func (r *JobRepositoryMongo) FindById(id string) (*domain.Job, error) {
	filter := bson.M{"_id": id}

	var jobModel models.JobModel
	err := r.collection.FindOne(context.Background(), filter).Decode(&jobModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Debugf("Job not found in the repository. ID: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Error finding job by ID in the repository: %v", err)
		return nil, err
	}

	job := r.mapper.ToDomain(&jobModel)
	r.logger.Debugf("Found job by ID in the repository. ID: %s", id)
	return job, nil
}

// Update updates an existing document in the collection.
func (r *JobRepositoryMongo) Update(job *domain.Job) error {
	jobModel := r.mapper.ToPersistence(job)
	jobModel.UpdatedAt = time.Now()

	update := bson.M{
		"$set": jobModel,
	}
	_, err := r.collection.UpdateByID(context.Background(), jobModel.Id, update)
	if err != nil {
		r.logger.Errorf("Error updating job in the repository: %v", err)
		return err
	}

	r.logger.Debugf("Job updated in the repository. ID: %s", jobModel.Id)
	return nil
}
