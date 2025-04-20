package postgres

import (
	"time"

	"gorm.io/gorm"
)

type DbCluster struct {
	master  *Connection
	slaves  []*Connection
	counter uint64
}

type Connection struct {
	config DBConfig
	db     *gorm.DB
}

type DBConfig struct {
	Host                   string
	Port                   string
	Username               string
	Password               string
	Dbname                 string
	MaxOpenConnections     int
	MaxIdleConnections     int
	ConnMaxLifetime        time.Duration
	DebugMode              bool
	PrepareStmt            bool
	SkipDefaultTransaction bool
}
