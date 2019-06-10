package database

import (
	"data_base/presentation/logger"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/xlab/closer"
	"time"
)

const (
	maxConnections = 100
	acquireTimeout = 5 * time.Second
)

func loadConfiguration() (pgxConfig pgx.ConnConfig) {
	pgxConfig.Host = "localhost"
	pgxConfig.User = "mac"
	pgxConfig.Password = "1209qawsed"
	pgxConfig.Database = "postgres"
	pgxConfig.Port = 5432
	return
}

type databaseManager struct {
	dataBase *pgx.ConnPool
}

var database *databaseManager

func init() {
	pgxConfig := loadConfiguration()
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig, MaxConnections: maxConnections, AcquireTimeout: acquireTimeout}

	dataBase, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		logger.Fatal.Println(err.Error())
		return
	}

	fmt.Println("DB connection opened")

	database = &databaseManager{
		dataBase: dataBase,
	}

	closer.Bind(closeConnection)
}

func closeConnection() {
	database.dataBase.Close()
	fmt.Println("DB connection closed")
}

func GetInstance() *databaseManager {
	return database
}
