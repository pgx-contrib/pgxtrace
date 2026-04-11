package pgxtrace_test

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxtrace"
)

func ExampleCompositeQueryTracer() {
	config, err := pgxpool.ParseConfig(os.Getenv("PGX_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	config.ConnConfig.Tracer = pgxtrace.CompositeQueryTracer{}

	conn, err := pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.TODO(), "SELECT * from customer")
	if err != nil {
		panic(err)
	}
	// close the rows
	defer rows.Close()

	// Customer struct must be defined
	type Customer struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}

	for rows.Next() {
		customer, err := pgx.RowToStructByName[Customer](rows)
		if err != nil {
			panic(err)
		}

		fmt.Println(customer.FirstName)
	}
}
