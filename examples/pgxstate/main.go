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

type testStruct struct {
	String string    `json:"s"`
	Number int64     `json:"n"`
	Float  float64   `json:"f"`
	Time   time.Time `json:"t"`
	Bool   bool      `json:"b"`
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", connStr)
	}
	defer pool.Close()

	n := time.Now().UnixNano()
	s := &pgxmap.State[testStruct]{
		Data: testStruct{
			String: "hello",
			Number: n,
			Float:  float64(n) / float64(time.Second),
			Time:   time.Now(),
			Bool:   true,
		},
	}

	id := fmt.Sprintf("id-%d", n)

	if err = save(ctx, pool, id, s); err != nil {
		log.Fatalf("error saving state: %s", err)
	}

	s2, err := get[testStruct](ctx, pool, id)
	if err != nil {
		log.Fatalf("error getting state: %s", err)
	}

	if s2 == nil {
		log.Fatalf("unable to find state with id: %s", id)
	}

	fmt.Printf("expected: %+v\n", s.Data)

	fmt.Printf("got:      %+v\n", s2.Data)
}

func get[T any](ctx context.Context, p *pgxpool.Pool, id string) (*pgxmap.State[T], error) {
	if p == nil {
		return nil, errors.New("nil pool")
	}

	row := p.QueryRow(ctx, selectSQL, id)

	var s pgxmap.State[T]
	if err := row.Scan(&s); err != nil && err != pgx.ErrNoRows {
		return nil, errors.Wrap(err, "error scanning row")
	}

	return &s, nil
}

func save[T any](ctx context.Context, p *pgxpool.Pool, id string, s *pgxmap.State[T]) error {
	if p == nil || s == nil {
		return errors.New("pool and cm required")
	}

	if _, err := p.Exec(ctx, insertSQL, id, s); err != nil {
		return errors.Wrap(err, "error inserting row")
	}

	return nil
}
