package etl

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool

func init() {
	fmt.Printf("PGHOST: '%s'\nPGUSER: '%s'\nPGPASSWORD: '%s'\nPGDATABASE: '%s'", getEnv("PGHOST", "127.0.0.1"), getEnv("PGUSER", "instruments"), getEnv("PGPASSWORD", "raspberry"), getEnv("PGDATABASE", "picobank"))

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     getEnv("PGHOST", "127.0.0.1"),
			User:     getEnv("PGUSER", "instruments"),
			Password: getEnv("PGPASSWORD", "raspberry"),
			Database: getEnv("PGDATABASE", "picobank"),
		},
		MaxConnections: 10,
		// AfterConnect:   afterConnectCallback,
	}
	var err error
	pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		fmt.Println("Unable to create connection pool", err)
		os.Exit(1)
	}

	fmt.Println("\t Connexion a la base établie")
}

// Connection retrieve connection from the pool
func Connection() *pgx.Conn {
	conn, err := pool.Acquire()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		panic(err)
	}
	//defer Release(conn)

	return conn
}

// Release release connection to the pool
func Release(conn *pgx.Conn) {
	pool.Release(conn)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
