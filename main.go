package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	db, err := sql.Open("mysql", "user:password@tcp(localhost:13306)/monitoringdb?charset=utf8")
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

	http.HandleFunc("/", handleIndex)
	http.Handle("/heroes", heroesHandler())
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":3001", nil))
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
	_, err := io.WriteString(rw, "<h1>Monitoring App</h1>")
	checkErr(err)
}

func heroesHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
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
