package utils

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}

func GenerateMongoDBID() string {
	return primitive.NewObjectID().Hex()
}
