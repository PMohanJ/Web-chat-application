package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	UpdateOne(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	DeleteMany(context.Context, interface{}) (int64, error)
	Aggregate(context.Context, interface{}) (Cursor, error)
	Drop(context.Context) error
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

type MongoPipeline []bson.D

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
	cursor, err := mcl.coll.Find(ctx, filter)
	return &mongoCursor{cur: cursor}, err
}

func (mcl *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleresult := mcl.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleresult}
}

func (mcl *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	insertId, err := mcl.coll.InsertOne(ctx, document)
	return insertId.InsertedID, err
}

func (mcl *mongoCollection) InsertMany(ctx context.Context, documents []interface{}) ([]interface{}, error) {
	insertId, err := mcl.coll.InsertMany(ctx, documents)
	return insertId.InsertedIDs, err
}

func (mcl *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	updateResult, err := mcl.coll.UpdateOne(ctx, filter, update)
	return updateResult, err
}

func (mcl *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	updateResult, err := mcl.coll.UpdateMany(ctx, filter, update)
	return updateResult, err
}

func (mcl *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	deletedCount, err := mcl.coll.DeleteOne(ctx, filter)
	return deletedCount.DeletedCount, err
}

func (mcl *mongoCollection) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	deletedCount, err := mcl.coll.DeleteMany(ctx, filter)
	return deletedCount.DeletedCount, err
}

func (mcl *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error) {
	cursor, err := mcl.coll.Aggregate(ctx, pipeline)
	return &mongoCursor{cur: cursor}, err
}

func (mcl *mongoCollection) Drop(ctx context.Context) error {
	return mcl.Drop(ctx)
}

func (mc *mongoCursor) All(ctx context.Context, result interface{}) error {
	return mc.cur.All(ctx, result)
}

func (mc *mongoCursor) Decode(v interface{}) error {
	return mc.cur.Decode(v)
}

func (mc *mongoCursor) Next(ctx context.Context) bool {
	return mc.cur.Next(ctx)
}

func (mc *mongoCursor) Close(ctx context.Context) error {
	return mc.cur.Close(ctx)
}

func (msr *mongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
}
