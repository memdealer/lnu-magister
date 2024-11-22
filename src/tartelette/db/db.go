package db

import (
	"github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
)

func InitDatabase() *memdb.MemDB {
	// Create the DB schema
	db, err := memdb.NewMemDB(Schema)

	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
