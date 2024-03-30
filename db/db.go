package db

import (
	"context"
	"github.com/induzo/gocom/database/pginit/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func New(dbURI string) *pgxpool.Pool {
	pgi, err := pginit.New(
		dbURI, pginit.WithDecimalType(),
		pginit.WithGoogleUUIDType(),
	)

	if err != nil {
		log.Fatal("Can't connect to db", err.Error())
	}

	pool, err := pgi.ConnPool(context.TODO())

	if err != nil {
		log.Fatal("Can't connect to db", err.Error())
	}

	return pool
}
