package store

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MockMongoClient is a mock for the MongoDB client
type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	args := m.Called(ctx, rp)
	return args.Error(0)
}
func (m *MockMongoClient) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(mongo.ListDatabasesResult), args.Error(1)
}
func (m *MockMongoClient) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockMongoClient) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	args := m.Called(ctx, opts, fn)
	return args.Error(0)
}
func (m *MockMongoClient) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := m.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}
func (m *MockMongoClient) NumberSessionsInProgress() int {
	args := m.Called()
	return args.Int(0)
}
func (m *MockMongoClient) Timeout() *time.Duration {
	args := m.Called()
	return args.Get(0).(*time.Duration)
}

func (m *MockMongoClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoClient) Close(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoClient) Database(name string) Database {
	args := m.Called(name)
	return args.Get(0).(Database)
}

func (m *MockMongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	args := m.Called(opts)
	return args.Get(0).(mongo.Session), args.Error(1)
}

func (m *MockMongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// MockMongoDatabase is a mock for the MongoDB database
type MockMongoDatabase struct {
	mock.Mock
}

func (m *MockMongoDatabase) Collection(name string) Collection {
	args := m.Called(name)
	return args.Get(0).(Collection)
}

func (m *MockMongoDatabase) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMongoDatabase) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	args := m.Called(ctx, runCommand, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoDatabase) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMongoDatabase) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).([]string), args.Error(1)
}

// MockMongoCollection is a mock for the MongoDB collection
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) Name() string {
	args := m.Called()
	return args.String(0) // Return the first argument as the collection name (string)
}

func (m *MockMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	if result, ok := args.Get(0).(*mongo.DeleteResult); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateMany mocks the UpdateMany method of mongo.Collection
func (m *MockMongoCollection) UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	args := m.Called(ctx, documents, opts)
	return args.Get(0).(*mongo.InsertManyResult), args.Error(1)
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	result := &mongo.UpdateResult{}
	if args.Get(0) != nil {
		result = args.Get(0).(*mongo.UpdateResult)
	}
	return result, args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	// Return the result directly
	return args.Get(0).(*mongo.SingleResult)
}

type MockSession struct {
	mock.Mock
}

func (m *MockSession) StartTransaction(opts ...*options.TransactionOptions) error {
	args := m.Called(opts)
	return args.Error(0)
}

func (m *MockSession) AbortTransaction(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSession) CommitTransaction(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSession) WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
	args := m.Called(ctx, fn, opts)
	return args.Get(0), args.Error(1)
}

func (m *MockSession) EndSession(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockSession) ClusterTime() bson.Raw {
	args := m.Called()
	return args.Get(0).(bson.Raw)
}

func (m *MockSession) OperationTime() *primitive.Timestamp {
	args := m.Called()
	return args.Get(0).(*primitive.Timestamp)
}

func (m *MockSession) Client() *mongo.Client {
	args := m.Called()
	return args.Get(0).(*mongo.Client)
}

func (m *MockSession) ID() bson.Raw {
	args := m.Called()
	return args.Get(0).(bson.Raw)
}

func (m *MockSession) AdvanceClusterTime(clusterTime bson.Raw) error {
	args := m.Called(clusterTime)
	return args.Error(0)
}

func (m *MockSession) AdvanceOperationTime(opTime *primitive.Timestamp) error {
	args := m.Called(opTime)
	return args.Error(0)
}

// Implement the 'session()' method to satisfy the interface
// Assuming this method does not need further customization
func (m *MockSession) session() {
	// This method can be left empty, as it's just part of the interface
	m.Called()
}
