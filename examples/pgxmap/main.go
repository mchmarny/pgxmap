package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mchmarny/pgxmap"
	"github.com/pkg/errors"
)

const (
	selectSQL = `SELECT cm FROM example WHERE id = $1`
	insertSQL = `INSERT INTO example (id, cm) VALUES ($1, $2)`
)

var (
	// example: "postgresql://demo:demo@localhost:5432/demo"
	connStr = os.Getenv("CONN_STR")
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", connStr)
	}
	defer pool.Close()

	utc := time.Now().UTC()
	nano := utc.UnixNano()

	m1 := pgxmap.ConfigMap{
		"int":     int(1),
		"int8":    int8(2),
		"int16":   int16(3),
		"int32":   int32(4),
		"int64":   nano,
		"uint":    uint(5),
		"uint8":   uint8(6),
		"uint16":  uint16(7),
		"uint32":  uint32(8),
		"uint64":  uint64(9),
		"float32": float32(123456.0987),
		"float64": float64(nano) / float64(0.3),
		"bool":    true,
		"string":  "string",
		"utc":     utc,
	}

	id := fmt.Sprintf("id-%d", nano)

	if err = save(ctx, pool, id, m1); err != nil {
		log.Fatalf("error saving config: %s", err)
	}

	m2, err := get(ctx, pool, id)
	if err != nil {
		log.Fatalf("error getting config: %s", err)
	}

	if len(m1) != len(m2) {
		log.Fatalf("configs have different length (want:%d, got:%d)", len(m1), len(m2))
	}

	for k, v := range m1 {
		fmt.Printf("%s - saved: %v as %T, got: %v as %T\n", k, v, v, m2[k], m2[k])
	}
}

func get(ctx context.Context, p *pgxpool.Pool, id string) (pgxmap.ConfigMap, error) {
	if p == nil {
		return nil, errors.New("nil pool")
	}

	row := p.QueryRow(ctx, selectSQL, id)

	var cm pgxmap.ConfigMap
	if err := row.Scan(&cm); err != nil && err != pgx.ErrNoRows {
		return nil, errors.Wrap(err, "error scanning row")
	}

	return cm, nil
}

func save(ctx context.Context, p *pgxpool.Pool, id string, m pgxmap.ConfigMap) error {
	if p == nil || m == nil {
		return errors.New("pool and cm required")
	}

	if _, err := p.Exec(ctx, insertSQL, id, m); err != nil {
		return errors.Wrap(err, "error inserting row")
	}

	return nil
}
