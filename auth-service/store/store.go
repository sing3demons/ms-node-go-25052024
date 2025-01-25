package store

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStore(ctx context.Context) *mongo.Client {
	url := os.Getenv("MONGO_URI")

	loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetLoggerOptions(loggerOptions))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	return client
}

type Store interface {
	Close(ctx context.Context)
	Database(name string) Database
	Disconnect(ctx context.Context) error
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
}

type store struct {
	Client *mongo.Client
}

func New(client *mongo.Client) Store {
	return &store{Client: client}
}

func (s *store) Close(ctx context.Context) {
	s.Client.Disconnect(ctx)
}

func (s *store) Database(name string) Database {
	return &mongoDatabase{database: s.Client.Database(name)}
}

func (s *store) Disconnect(ctx context.Context) error {
	return s.Client.Disconnect(ctx)
}

func (s *store) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return s.Client.StartSession(opts...)
}

type Database interface {
	Name() string
	Collection(name string) Collection
}

type mongoDatabase struct {
	database *mongo.Database
}

func (d *mongoDatabase) Name() string {
	return d.database.Name()
}

func (d *mongoDatabase) Collection(name string) Collection {
	return NewCollection(d.database.Collection(name))
}

type Collection interface {
	Name() string
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

func NewCollection(collection Collection) Collection {
	return &mongoCollection{collection: collection}
}

type mongoCollection struct {
	collection Collection
}

func (c *mongoCollection) Name() string {
	return c.collection.Name()
}

type SingleResult interface {
	Decode(v interface{}) error
}

func (c *mongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.collection.FindOne(ctx, filter, opts...)
}

func (s *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	return s.collection.Find(ctx, filter, opts...)
}

func (s *mongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return s.collection.InsertOne(ctx, document, opts...)
}

func (s *mongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return s.collection.InsertMany(ctx, documents, opts...)
}

func (s *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return s.collection.UpdateOne(ctx, filter, update, opts...)
}

func (s *mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return s.collection.UpdateMany(ctx, filter, update, opts...)
}

func (s *mongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return s.collection.DeleteOne(ctx, filter, opts...)
}

func (s *mongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return s.collection.DeleteMany(ctx, filter, opts...)
}
