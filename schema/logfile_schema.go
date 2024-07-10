package schema

import "github.com/hashicorp/go-memdb"

const TABLE_NAME = "log"

func LogfileSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			TABLE_NAME: {
				Name: TABLE_NAME,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Index"},
					},
					"operation": {
						Name:    "operation",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Operation"},
					},
					"term": {
						Name:    "term",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Term"},
					},
				},
			},
		},
	}
}
