# pgxmap

Custom map to use with [pgx](https://github.com/jackc/pgx) help with `map[string]interface{}` decoding and encoding correctly map items types (e.g. `int64`) when saved and retrieved from Postgres DB. 

## Why 

`pgx` stores JSON-marshalled objects in DB as JSONB type. For Go `structs` with priory know types, `pgx` driver encodes Go types to JSON and back automatically. However, when unmarshalling into a Go interface (e.g. `map[string]interface{}`) Go decoding has to do some guessing to map certain types when scanning DB selection results (e.g. `int64` vs `float64`). There isn't really much `pgx` can actually do here since JavaScript itself does not distinguish between integers and floats (see [this](https://tools.ietf.org/html/rfc7159#section-6) for details). There are similar issues with time precision as well.

> The full demo of the error resulting from an idiomatic implementation is available [here](examples/idiom/main.go).

There are 3 diff ways you can deal with this issue:

1. Avoid using map with interface value altogether
2. Use map with all values as strings and parse on return
3. Use the `driver.Valuer` interface to do your own encoding and decoding

The `pgxmap` uses the 3rd approach and ensures that the `int`, `float` derivative types (`int64`, `int32`, `uint16`, `float64` etc.) and `time` in your map are automatically correctly encoded and decoded to ensure exact same type/precision while still providing that basic Go map-like functionality. 

## Usage 

The basic usage of the `pgxmap` library falls into 3 patterns: 

* Config map creation 
* Saving it to DB
* Retrieving it from DB

### Creation 

> For illustration purposes, this example creates a config map with only 3 items with the different types. See [map_test.go](./map_test.go) for broader example.

```go
	bigInt := time.Now().UnixNano()

	m := pgxmap.ConfigMap{
		"int64":   bigInt,
		"uint32":  uint32(bigInt),
		"float64": float64(bigInt) / float64(3),
	}
```

### Saving 

Assuming the DB table schema looks something like this (key is the `JSONB` column type): 

```sql
CREATE TABLE IF NOT EXISTS demo (
    id varchar NOT NULL PRIMARY KEY,
    cf JSONB
);
```

Saving the map then is as simple as using it as yet another parameter in your `Exec` method:

> Assuming the `sql` variable holding your `SELECT` statement is already defined.

```go
func save(ctx context.Context, p *pgxpool.Pool, id string, m pgxmap.ConfigMap) error {
	// input validation skipped for brevity 
	if _, err := pool.Exec(ctx, sql, id, m); err != nil {
		return errors.Wrap(err, "error inserting row")
	}
	return nil
}
```

### Retrieval 

Similarly, during retrieval you simply provide a pointer to the `ConfigMap` and use the `row.Scan` method:

> Assuming the `sql` variable holding your `INSERT` statement is already defined.

```go
func get(ctx context.Context, p *pgxpool.Pool, id string) (pgxmap.ConfigMap, error) {
    // input validation skipped for brevity 
	var m pgxmap.ConfigMap
	row := pool.QueryRow(ctx, sql, id)
	if err := row.Scan(&cm); err != nil && err != pgx.ErrNoRows {
		return nil, errors.Wrap(err, "error scanning row")
	}

	return cm, nil
}
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
