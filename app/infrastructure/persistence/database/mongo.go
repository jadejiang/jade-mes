package database

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongodbInst *mongo.Database

// GetMongoDB ...
func GetMongoDB() *mongo.Database {
	if mongodbInst == nil {
		log.Fatalln(errNotInited)
	}

	return mongodbInst
}

func loadMongoDBConfig(config *viper.Viper) *DBConfig {
	var mongodbCnf DBConfig
	config.UnmarshalKey("mongodb", &mongodbCnf)

	return &mongodbCnf
}

func initMongoDBConn(config *DBConfig) string {
	//var authStr string
	//if config.User != "" {
	//	authStr = fmt.Sprintf("%v:%v@", config.User, config.Password)
	//}
	//var dbStr string
	//if config.DBname != "" {
	//	dbStr = fmt.Sprintf("/%v", config.DBname)
	//}
	//
	//return fmt.Sprintf("mongodb://%v%v:%v%v", authStr, config.Host, config.Port, dbStr)
	return config.Url
}

func initMongoDB(connSettings interface{}, config *DBConfig) {
	connStr, ok := connSettings.(string)
	if !ok {
		log.Fatalf("invalid mongodb connection string type, expected: [string], got: [%v]", reflect.TypeOf(connSettings).String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatalf("%s: <%s>\n", errConnFailed, err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("%s: <%s>\n", errConnFailed, err)
	}

	mongodbInst = client.Database(config.DBname)
}
