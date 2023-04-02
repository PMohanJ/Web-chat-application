package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	Find(context.Context, interface{}) (Cursor, error)
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	InsertMany(context.Context, []interface{}) ([]interface{}, error)
	UpdateOne(context.Context, interface{}, interface{}) (interface{}, error)
	UpdateMany(context.Context, interface{}, interface{}) ([]interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	DeleteMany(context.Context, interface{}) (int64, error)
	Aggregrate(context.Context, interface{}) (Cursor, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type Cursor interface {
	All(context.Context, interface{}) error
	Decode(interface{}) error
	Next(context.Context) bool
	Close(context.Context) error
}

type Client interface {
	Connect(context.Context) error
	Database(string) Database
	Disconnect(context.Context) error
	Ping(context.Context) error
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCursor struct {
	cur *mongo.Cursor
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func NewClient(connection string) (Client, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(connection).SetServerAPIOptions(serverAPIOptions)

	c, err := mongo.NewClient(clientOptions)
	return &mongoClient{cl: c}, err
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (mc *mongoClient) Database(dbName string) Database {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, readpref.Primary())
}

func (md *mongoDatabase) Collection(collectionName string) Collection {
	collection := md.db.Collection(collectionName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mcl *mongoCollection) Find(ctx context.Context, filter interface{}) (Cursor, error) {

}

func (mcl *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {

}

func (mcl *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {

}

func (mcl *mongoCollection) InsertMany(ctx context.Context, document []interface{}) ([]interface{}, error) {

}

func (mcl *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error) {

}

func (mcl *mongoCollection) UpdareMany(ctx context.Context, filter interface{}, update interface{}) ([]interface{}, error) {

}

func (mcl *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {

}

func (mcl *mongoCollection) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {

}

func (mcl *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error) {

}
