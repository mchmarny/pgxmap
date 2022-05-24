package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const connStr = "postgresql://dev:dev@localhost:5432/cm"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", connStr)
	}
	defer pool.Close()

	nano := time.Now().UnixNano()
	key := "number"

	c1 := map[string]interface{}{
		key: nano, // large enough number for JSON serializer to interpret it as float64
	}

	id := fmt.Sprintf("id-%d", nano)

	if err = save(ctx, pool, id, c1); err != nil {
		log.Fatalf("error saving config: %s", err)
	}

	c2, err := get(ctx, pool, id)
	if err != nil {
		log.Fatalf("error getting config: %s", err)
	}

	if c2 == nil {
		log.Fatalf("no config found")
	}

	val, ok := c2[key]
	if !ok {
		log.Fatalf("no key found: %s", key)
	}

	if val.(int64) != nano {
		log.Fatalf("wrong value: %v, expected: %d", val, nano)
	}
}

const selectSQL = `SELECT cf FROM demo WHERE id = $1`

func get(ctx context.Context, pool *pgxpool.Pool, id string) (map[string]interface{}, error) {
	if pool == nil {
		return nil, errors.New("nil pool")
	}

	row := pool.QueryRow(ctx, selectSQL, id)

	var cm map[string]interface{}
	if err := row.Scan(&cm); err != nil && err != pgx.ErrNoRows {
		return nil, errors.Wrap(err, "error scanning row")
	}

	return cm, nil
}

const insertSQL = `INSERT INTO demo (id, cf) VALUES ($1, $2)`

func save(ctx context.Context, pool *pgxpool.Pool, id string, cm map[string]interface{}) error {
	if pool == nil {
		return errors.New("nil pool")
	}
	if cm == nil {
		return errors.New("nil cm")
	}

	// name, executed, type, config, state
	if _, err := pool.Exec(ctx, insertSQL, id, cm); err != nil {
		return errors.Wrap(err, "error inserting row")
	}

	return nil
}
