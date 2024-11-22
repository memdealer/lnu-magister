package db

import (
	"github.com/hashicorp/go-memdb"
)

// Schema for the database
var Schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"runner": {
			Name: "runner",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Name"},
				},
				"hostName": {
					Name:    "hostName",
					Unique:  false, // Not unique since multiple runners can be on the same host
					Indexer: &memdb.StringFieldIndex{Field: "HostName"},
				},
			},
		},
		"host": {
			Name: "host",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "HostName"}, // Using HostName as the primary key (id)
				},
			},
		},
	},
}
