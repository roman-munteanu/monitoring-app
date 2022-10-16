package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	// httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type App struct {
	ctx context.Context
	db  *sql.DB
}

type Hero struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

var app App

func init() {
	// db, err := sql.Open("mysql", "user:password@tcp(localhost:13306)/monitoringdb?charset=utf8")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}

	app = App{
		ctx: context.Background(),
		db:  db,
	}
}

func main() {
	defer app.shutdown()
	// providing the tracer manually
	ddAgentAddr := os.Getenv("DATADOG_AGENT_ADDR")
	tracer.Start(
		tracer.WithEnv("dev"),
		tracer.WithService("monitoring-app"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithAgentAddr(ddAgentAddr))
	defer tracer.Stop()

	// mux := httptrace.NewServeMux()
	mux := muxtrace.NewRouter()

	mux.HandleFunc("/", handleIndex)
	mux.Handle("/heroes", heroesHandler())
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(os.Getenv("APP_SERVER_ADDR"), nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) shutdown() {
	err := a.db.Close()
	checkErr(err)
}

func handleIndex(rw http.ResponseWriter, req *http.Request) {
	// span := tracer.StartSpan("index.request", tracer.ResourceName("/"))
	// defer span.Finish()
	// span.SetTag("http.url", req.URL.Path)
	// span.SetTag("index_tag_key", "index_tag_value")

	_, err := io.WriteString(rw, "<h1>Monitoring App</h1>")
	checkErr(err)
}

func heroesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// span := tracer.StartSpan("heroes.request", tracer.ResourceName("/heroes"))
		// defer span.Finish()
		// span.SetTag("http.url", req.URL.Path)
		// span.SetTag("heroes_tag_key", "heroes_tag_value")

		heroes, err := app.findAll()
		if err != nil {
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(rw).Encode(heroes)
		if err != nil {
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
		}
	})
}

func (a *App) findAll() ([]Hero, error) {
	rows, err := a.db.Query(`SELECT id, name, category FROM heroes`)
	checkErr(err)
	defer rows.Close()

	var records []Hero
	for rows.Next() {
		hero := Hero{}
		err := rows.Scan(&hero.ID, &hero.Name, &hero.Category)
		if err != nil {
			return []Hero{}, err
		}
		records = append(records, hero)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
