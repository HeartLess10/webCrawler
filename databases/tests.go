package databases

import (
	"log"
	"os"
	"webScraping/databases/dataTypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDbTestComblixQuery(clusterPath string, databaseName string, collectionName string, logger *log.Logger) {
	factory := NewDBFctory(clusterPath, databaseName, logger)
	database := factory.CreateDB(MongoDB)
	MongoDBTestAddRow(database)
	models := []mongo.WriteModel{
		mongo.NewReplaceOneModel().SetFilter(bson.D{{"name", "hhi"}}).
			SetReplacement(dataTypes.UserDataType{Name: "Cafe Zucchini", Age: 55, Job: "hamham"}), //changeing the docment

		mongo.NewUpdateOneModel().SetFilter(bson.D{{"name", "Cafe Zucchini"}}).
			SetUpdate(bson.D{{"$set", bson.D{{"name", "Zucchini Land"}}}}), //renameing the name
	}
	opts := options.BulkWrite().SetOrdered(true)
	results, err := database.ExecuteQuery(collectionName, []any{models, opts})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(results)
}

func MongoDBTestAddRow(database Database) {
	testUser := dataTypes.UserDataType{Name: "hhi", Age: 23, Job: "programmer"}
	database.AddSingleRow(os.Getenv("DATA_BASE_COLLECTION_NAME"), testUser)
}
