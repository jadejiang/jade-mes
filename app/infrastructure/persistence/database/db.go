package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// gorm dialect - mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbInst *gorm.DB

// GetDB returns database instance
func GetDB() *gorm.DB {
	if dbInst == nil {
		log.Fatalln(errNotInited)
	}

	return dbInst
}

func loadDBConfig(config *viper.Viper) *DBConfig {
	var mysqlCnf DBConfig
	config.UnmarshalKey("mysql", &mysqlCnf)

	return &mysqlCnf
}

func initDBConn(config *DBConfig) string {
	connStr := fmt.Sprintf(
		"%s:%s@(%s:%v)/%s?charset=utf8mb4,utf8&parseTime=True",
		config.User, config.Password, config.Host, config.Port, config.DBname)

	return connStr
}

func initDB(connSettings interface{}, config *viper.Viper) {
	db, err := gorm.Open("mysql", connSettings.(string))
	if err != nil {
		log.Fatalf("%s: <%s>\n", errConnFailed, err)
	}

	maxIdleConns := config.GetInt("mysql.max_idle_connections")
	maxOpenConns := config.GetInt("mysql.max_open_connections")

	mysqlConn := db.DB()
	mysqlConn.SetMaxIdleConns(maxIdleConns)
	mysqlConn.SetMaxOpenConns(maxOpenConns)

	isReleaseMode := config.GetBool("release_mode")
	if !isReleaseMode {
		db.LogMode(true)
	}

	dbInst = db
}
