package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"testing"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "test_web"
	password = "pass"
	dbname   = "test_snippetbox"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Set the maximum number of concurrently open connections. Setting this to
	// less than or equal to 0 will mean there is no maximum limit. If the maximum
	// number of open connections is reached and a new connection is needed, Go will
	// wait until one of the connections is freed and becomes idle. From a
	// user perspective, this means their HTTP request will hang until a connection
	// is freed.
	db.SetMaxOpenConns(95)
	// Set the maximum number of idle connections in the pool. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(5)
	return db, nil
}


func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := openDB(psqlInfo)
	if err != nil {
		t.Fatal(err)
	}
	// Read the setup SQL script from file and execute the statements.
	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
