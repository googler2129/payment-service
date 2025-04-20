package cluster

import (
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
)

type Db struct {
	*postgres.DbCluster
}

var dbInstance *Db

func GetCluster() *Db {
	return dbInstance
}

func SetCluster(cluster *postgres.DbCluster) {
	dbInstance = &Db{cluster}
}
