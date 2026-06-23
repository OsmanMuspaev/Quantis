package user_service

import (
	"context"
	"time"

	"auth/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RolesResponse struct {
	Roles []string `bson:"roles" json:"roles"`
}

type GetRolesRequest struct {
	UserId string `json:"user_id"`
}

func GetUserRoles(userIDStr string) (RolesResponse, error) {
	var roles RolesResponse

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return roles, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{
		"roles": 1,
		"_id":   0,
	})

	err = storage.GetUserCollection().FindOne(
		ctx,
		bson.M{"_id": userID},
		findOptions,
	).Decode(&roles)

	if err != nil {
		return roles, err
	}
	return roles, nil
}
