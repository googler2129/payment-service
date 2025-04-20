package postgres

import (
	"fmt"
	"time"

	"github.com/mercor/payment-service/pkg/log"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDBInstance(master DBConfig, slaves *[]DBConfig) *DbCluster {
	db := getDbInstance(master, slaves)
	return db
}

func getDbInstance(master DBConfig, slaves *[]DBConfig) (instance *DbCluster) {
	slavesCount := len(*slaves)
	instance = &DbCluster{
		master: new(Connection),
		slaves: make([]*Connection, slavesCount),
	}

	instance.master = initDbConnection(master)

	for i := 0; i < slavesCount; i++ {
		instance.slaves[i] = initDbConnection((*slaves)[i])
	}
	return
}

func initDbConnection(config DBConfig) *Connection {
	gormLogger := logger.Default
	if config.DebugMode {
		gormLogger = gormLogger.LogMode(logger.Info)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Dbname,
	)

	gormDB, err := gorm.Open(postgres.Dialector{
		Config: &postgres.Config{
			DriverName: "nrpostgres", // For newrelic tracking
			DSN:        dsn,
		},
	}, &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		PrepareStmt:            config.PrepareStmt,
	})
	if err != nil {
		log.Errorf("Unable to make gorm connection | Error: %v", err)
		panic("Unable to make gorm connection")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Errorf("Unable to get sqlDB from gormDB | Error: %v", err)
		panic("Unable to get sqlDB from gormDB")
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	retries := 3
	for retries > 0 {
		err = sqlDB.Ping()
		if err != nil {
			log.Errorf("Unable to ping database server: %s, waiting 2 seconds before trying %d more times", err.Error(), retries)
			time.Sleep(time.Second * 2)
			retries--
		} else {
			err = nil
			break
		}
	}
	if err != nil {
		panic("Db not initialised")
	}

	conn := Connection{db: gormDB, config: config}
	return &conn
}
