package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-std/internal/config"
	"go-std/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	env, _ := config.Config()
	conn, err := pgx.Connect(context.Background(), env.GetString("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("tls config: ", tlsConfig)
	defer conn.Close(context.Background())

	q := sqlc.New(conn)

	// Create a slice to hold 10000 authors
	authors := make([]sqlc.CreateAuthorsParams, 10000)
	for i := 0; i < 10000; i++ {
		authors[i] = sqlc.CreateAuthorsParams{
			Name: fmt.Sprintf("Author %d", i+1),
			Bio:  pgtype.Text{String: fmt.Sprintf("Bio for author %d", i+1), Valid: true},
		}
	}

	// Time the batch insert operation
	start := time.Now()
	q.CreateAuthors(context.Background(), authors)
	elapsed := time.Since(start)
	fmt.Printf("Batch insert of 10000 authors took: %s\n", elapsed)

	author, err := q.GetAuthor(context.Background(), 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAuthor failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(author.Name)
}
