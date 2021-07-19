package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"calenwu.com/snippetbox/pkg/models/postgres"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "snippetbox"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.CookieStore
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	//f, e := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	//if e != nil {
	//	log.Fatal(e)
	//}
	//defer f.Close()
	infoLog := log.New(os.Stdin, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := openDB(psqlInfo)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	secret := flag.String("secret", "safjkladfuioejlj+32@afi", "Secret")
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static", "./ui/static", "Path to static assets")
	flag.Parse()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.NewCookieStore([]byte(*secret))

	tlsConfig := & tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: [] tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &postgres.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %", srv.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

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
