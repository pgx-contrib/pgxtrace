# pgxtrace

[![CI](https://github.com/pgx-contrib/pgxtrace/actions/workflows/ci.yml/badge.svg)](https://github.com/pgx-contrib/pgxtrace/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/pgx-contrib/pgxtrace?include_prereleases)](https://github.com/pgx-contrib/pgxtrace/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/pgx-contrib/pgxtrace.svg)](https://pkg.go.dev/github.com/pgx-contrib/pgxtrace)
[![License](https://img.shields.io/github/license/pgx-contrib/pgxtrace)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![pgx](https://img.shields.io/badge/pgx-v5-blue)](https://github.com/jackc/pgx)

`CompositeQueryTracer` is a decorator for [pgx v5](https://github.com/jackc/pgx)
that chains multiple tracers together. Assign it to `ConnConfig.Tracer` and
every database operation is dispatched to all registered tracers in order —
query, batch, connect, prepare, and copy-from.

## Installation

```bash
go get github.com/pgx-contrib/pgxtrace
```

## Usage

### Connection pool

```go
config, err := pgxpool.ParseConfig(os.Getenv("PGX_DATABASE_URL"))
if err != nil {
    panic(err)
}

config.ConnConfig.Tracer = pgxtrace.CompositeQueryTracer{
    &myTracer{},
    &anotherTracer{},
}

pool, err := pgxpool.NewWithConfig(context.Background(), config)
if err != nil {
    panic(err)
}
defer pool.Close()
```

### Mixing tracer types

Each element only needs to implement `pgx.QueryTracer`. Elements that also
implement `pgx.ConnectTracer`, `pgx.BatchTracer`, `pgx.PrepareTracer`, or
`pgx.CopyFromTracer` are automatically called for those operations:

```go
config.ConnConfig.Tracer = pgxtrace.CompositeQueryTracer{
    &pgxotel.QueryTracer{Name: "my-service"},  // all five interfaces
    &myAuditLogger{},                           // QueryTracer only
}
```

## Contributing

Contributions are welcome! Please open an issue or pull request.

To set up a development environment with [Nix](https://nixos.org):

```bash
nix develop
```

Or using the provided dev container:

```bash
devcontainer up --workspace-folder . --remove-existing-container
```

Then run the tests:

```bash
go tool ginkgo run -r
```

## License

[MIT](LICENSE)
