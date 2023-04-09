package repository

import (
	"context"
	"log"

	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type chatRepository struct {
	database   mongo.Database
	collection string
}

func NewChatRepository(db mongo.Database, collection string) domain.ChatRepository {
	return &chatRepository{
		database:   db,
		collection: collection,
	}
}

func (cr *chatRepository) Create(ctx context.Context, chat domain.Chat) (primitive.ObjectID, error) {
	collection := cr.database.Collection(cr.collection)

	insId, err := collection.InsertOne(ctx, chat)

	insertedId, ok := insId.(primitive.ObjectID)
	if !ok {
		log.Fatalln("Type assetion failed")
	}

	return insertedId, err
}

func (cr *chatRepository) FetchById(ctx context.Context, insertedId primitive.ObjectID) ([]bson.M, error) {
	collection := cr.database.Collection(cr.collection)

	matchStage := MatchStageBySingleField("_id", insertedId)

	lookupStage := LookUpStage("user", "users", "_id", "users")

	projectStage := ProjectStage("users.password", "created_at",
		"updated_at", "users.created_at", "users.updated_at")

	cursor, err := collection.Aggregate(ctx, mongo.MongoPipeline{matchStage, lookupStage, projectStage})
	if err != nil {
		return nil, err
	}

	var chat []bson.M
	if err := cursor.All(ctx, &chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (cr *chatRepository) FindByFilter(ctx context.Context, filter interface{}) error {
	collection := cr.database.Collection(cr.collection)

	var chat domain.Chat
	err := collection.FindOne(ctx, filter).Decode(&chat)

	return err
}

func (cr *chatRepository) FetchByFilter(ctx context.Context, filter primitive.D) ([]bson.M, error) {
	collection := cr.database.Collection(cr.collection)

	lookupStage := LookUpStage("user", "users", "_id", "users")

	projectStage := ProjectStage("users.password", "created_at",
		"updated_at", "users.created_at", "users.updated_at")

	var chat []bson.M
	cur, err := collection.Aggregate(ctx, mongo.MongoPipeline{filter, lookupStage, projectStage})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (cr *chatRepository) FetchWithLatestMessage(ctx context.Context, filter primitive.D) ([]bson.M, error) {
	collection := cr.database.Collection(cr.collection)

	lookupStage := LookUpStage("user", "users", "_id", "users")
	lookupStageLatestMessage := LookUpStage("message", "latestMessage", "_id", "latestMessage")
	projectStage := ProjectStage("users.password", "created_at",
		"updated_at", "users.created_at", "users.updated_at")

	cursor, err := collection.Aggregate(ctx, mongo.MongoPipeline{filter, lookupStage, lookupStageLatestMessage, projectStage})
	if err != nil {
		return nil, err
	}

	var chats []bson.M
	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}
