package repository

import (
	"context"
	"log"

	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	collection := ur.database.Collection(ur.collection)

	insId, err := collection.InsertOne(ctx, user)

	insertedId, ok := insId.(primitive.ObjectID)
	if !ok {
		log.Fatalln("Type assertion failed")
	}

	return insertedId, err
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	return user, err
}

func (ur *userRepository) FetchUsers(ctx context.Context, query string) ([]bson.M, error) {
	collection := ur.database.Collection(ur.collection)

	matchStage := bson.D{
		{"$match", bson.D{
			{"$or",
				bson.A{
					bson.D{{Key: "name", Value: primitive.Regex{Pattern: query}}},
					bson.D{{Key: "email", Value: primitive.Regex{Pattern: query}}},
				},
			},
		}},
	}
	projectStage := bson.D{
		{
			"$project", bson.D{
				{Key: "password", Value: 0},
				{Key: "created_at", Value: 0},
				{Key: "updated_at", Value: 0},
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, mongo.MongoPipeline{matchStage, projectStage})
	if err != nil {
		return nil, err
	}

	var users []bson.M
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
