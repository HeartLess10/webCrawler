package mongoDatabase

import (
	"context"
	"log"
	"sync"
	"webScraping/globalUtils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDbCluster *mongoDBCluster = nil
	once           sync.Once
)

type mongoDBCluster struct {
	clusterPath  string
	logger       *log.Logger
	client       *mongo.Client //cluster connection
	databaseName string
	context      context.Context
}

// Singleton creation of db connection with thread safe
func NewDB(clusterPath string, databaseName string, logger *log.Logger) *mongoDBCluster {
	once.Do(func() {
		mongoDbCluster = new(mongoDBCluster)
		mongoDbCluster.clusterPath = clusterPath
		mongoDbCluster.logger = logger
		mongoDbCluster.context = context.Background()
	})
	//database name can change for logic purposes
	mongoDbCluster.databaseName = databaseName

	return mongoDbCluster
}

func (mdbc *mongoDBCluster) ConnectToDBServer() error {
	client, err := mongo.Connect(mdbc.context, options.Client().ApplyURI(mdbc.clusterPath))
	if err != nil {
		return globalUtils.CreateError("Failed to connect to mongodb server", mdbc.logger)
	}
	mdbc.client = client
	mdbc.logger.Println("Connected to mongodb server")
	return nil
}

func (mdbc *mongoDBCluster) AddSingleRow(collectionName string, object any) error {
	coll := mdbc.client.Database(mdbc.databaseName).Collection(collectionName)
	var bsonDocument *bson.D //map[string]any //map[string]interface{}

	bsonDocument, err := toBsonDoc(object)
	if err != nil {
		return globalUtils.CreateError("Failed to convert object to bson", mdbc.logger)
	}

	_, err = coll.InsertOne(mdbc.context, *bsonDocument)
	if err != nil {
		return globalUtils.CreateError("Could not add the row to the mongo database", mdbc.logger)
	}
	return nil
}

// expects in query[0] to be []mongo.WriteModel and in query[1] to be *options.BulkWriteOptions
func (mdbc *mongoDBCluster) ExecuteQuery(collectionName string, query []any) (any, error) {
	coll := mdbc.client.Database(mdbc.databaseName).Collection(collectionName)
	models, ok := query[0].([]mongo.WriteModel)

	if !ok {
		return nil, globalUtils.CreateError("Could not convert the query[0] to a write model", mdbc.logger)
	}
	opts, ok := query[1].(*options.BulkWriteOptions)
	if !ok {
		return nil, globalUtils.CreateError("Could not convert the query[1] to a bulk write options", mdbc.logger)
	}

	results, err := coll.BulkWrite(mdbc.context, models, opts)
	if err != nil {
		return nil, globalUtils.CreateError("Could not find the rows in the mongo database", mdbc.logger)
	}
	return results, nil
}

func (mdbc *mongoDBCluster) CloseConnection() {
	mdbc.client.Disconnect(mdbc.context)
	mdbc.logger.Println("Disconnected from mongodb server")
}
