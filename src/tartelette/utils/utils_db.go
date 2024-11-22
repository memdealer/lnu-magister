package utils

import (
	"TartaLette/api/models"
	"github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
)

func FillDbWithStateValuesFromGithub(dbConnection *memdb.MemDB, hostnames []models.Hostname, runners []models.Runner) error {
	txn := dbConnection.Txn(true)
	defer txn.Abort()

	for _, host := range hostnames {
		if err := txn.Insert("host", host); err != nil {
			logrus.Fatal(err)
			return err
		}
	}

	for _, runner := range runners {
		if err := txn.Insert("runner", runner); err != nil {
			logrus.Fatal(err)
			return err
		}
	}

	txn.Commit()
	return nil
}
