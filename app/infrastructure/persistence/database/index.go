package database

import (
	"log"

	"jade-mes/config"
)

const (
	mysql = iota
	postgresql
	sqlte3
	sqlserver
	mongodb
	redis
	errInvalidDB  = "invalid db type"
	errNotInited  = "db instance is not inited"
	errConnFailed = "connect to db failed"
)

// DBConfig defines the database' config schema
type DBConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
	DBname   string
	Url      string
}

func init() {
	println("initing database...")

	settings := config.GetConfig()

	mysqlConf := loadDBConfig(settings)
	// mongodbConf := loadMongoDBConfig(settings)

	dbConnSettings := connStrFactory(mysql, mysqlConf)
	initDB(dbConnSettings, settings)

	initRedis(settings)

	// mongodbConnSettings := connStrFactory(mongodb, mongodbConf)
	// initMongoDB(mongodbConnSettings, mongodbConf)
}

func connStrFactory(dbType int, conf *DBConfig) interface{} {
	var connSettings interface{}

	switch dbType {
	case mysql:
		connSettings = initDBConn(conf)
	case mongodb:
		connSettings = initMongoDBConn(conf)
	default:
		log.Fatalf("%s: <%v>\n", errInvalidDB, dbType)
	}

	return connSettings
}
