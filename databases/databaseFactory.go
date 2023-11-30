package databases

import (
	"log"
	"webScraping/databases/mongoDatabase"
)

const (
	MongoDB = iota
	//PostGress
)

type Database interface {
	ConnectToDBServer() error
	AddSingleRow(collectionName string, object any) error
	CloseConnection()
	ExecuteQuery(collectionName string, query []any) (any, error)
}

type dbFctory struct {
	clusterPath  string
	databaseName string
	logger       *log.Logger
}

func NewDBFctory(clusterPath string, databaseName string, logger *log.Logger) *dbFctory {
	return &dbFctory{clusterPath: clusterPath, databaseName: databaseName, logger: logger}
}
func (dbF *dbFctory) CreateDB(choosenDataBase int) Database {
	switch choosenDataBase {
	case 0:
		return mongoDatabase.NewDB(dbF.clusterPath, dbF.databaseName, dbF.logger)
	// case 1:
	// 	return nil //NewPostGress(dbF.clusterPath, dbF.databaseName, dbF.logger)
	default:
		return nil
	}
}
