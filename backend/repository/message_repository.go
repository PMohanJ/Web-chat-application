package repository

import (
	"context"
	"log"

	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageRepository struct {
	database   mongo.Database
	collection string
}

func NewMessageRepository(db mongo.Database, collection string) domain.MessageRepository {
	return &messageRepository{
		database:   db,
		collection: collection,
	}
}

func (mr *messageRepository) Create(ctx context.Context, message domain.Message) (primitive.ObjectID, error) {
	collection := mr.database.Collection(mr.collection)

	insId, err := collection.InsertOne(ctx, message)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedId, ok := insId.(primitive.ObjectID)
	if !ok {
		log.Panic("Type assertion failed")
	}
	return insertedId, nil
}

func (mr *messageRepository) FetchById(ctx context.Context, field string, id primitive.ObjectID) ([]bson.M, error) {
	collection := mr.database.Collection(mr.collection)

	matchStage := MatchStageBySingleField(field, id)

	lookupStage := LookUpStage("user", "sender", "_id", "sender")

	projectStage := ProjectStage("sender.password", "created_at",
		"updated_at", "sender.created_at", "sender.updated_at")

	cursor, err := collection.Aggregate(ctx, mongo.MongoPipeline{matchStage, lookupStage, projectStage})
	if err != nil {
		return nil, err
	}

	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func (mr *messageRepository) UpdateById(ctx context.Context, filter primitive.D, update primitive.D) error {
	collection := mr.database.Collection(mr.collection)

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (mr *messageRepository) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	collection := mr.database.Collection(mr.collection)

	filter := MatchStageBySingleField("_id", id)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
