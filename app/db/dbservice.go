package db

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/dsn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstances = make(map[string]*gorm.DB)

func GetDBConnection(wrapdsn string) (*gorm.DB, error) {
	dsn, ok := dsn.GetConfig().DnsServers[wrapdsn]
	if !ok {
		return nil, fmt.Errorf("dsn not found")
	}

	if db, exists := dbInstances[dsn]; exists {
		return db, nil
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", wrapdsn)
	}

	dbInstances[dsn] = db
	return db, nil
}
