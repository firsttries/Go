package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

type Tv struct { //Tv Name of Tables and columns
	id          int16
	brand       string
	manufacture string
	model       string
	year        int16
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlserver", "sqlserver://sa:helpline@127.0.0.1:1433?database=TV_WB") //connect to base
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/tv", tvsIndex)         // show data of base
	http.HandleFunc("/tv/show", tvsShow)     // show on request
	http.HandleFunc("/tv/create", tvsCreate) //create data
	http.ListenAndServe(":7777", nil)        //server port
}

func tvsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM tv") // select from table
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	tv := make([]*Tv, 0)
	for rows.Next() {
		tvM := new(Tv)
		err := rows.Scan(&tvM.id, &tvM.brand, &tvM.manufacture, &tvM.model, &tvM.year)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		tv = append(tv, tvM)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, tvM := range tv {
		fmt.Fprintf(w, "%d, %s, %s, %s, %d\n", tvM.id, tvM.brand, tvM.manufacture, tvM.model, tvM.year) //output
	}
}

func tvsShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	brand := r.FormValue("brand")
	if brand == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	row := db.QueryRow("SELECT * FROM tv WHERE brand = @p1", brand) // @p1 for MSSQL, $1 for postgre

	tvM := new(Tv)
	err := row.Scan(&tvM.id, &tvM.brand, &tvM.manufacture, &tvM.model, &tvM.year)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%d, %s, %s, %s, %d\n", tvM.id, tvM.brand, tvM.manufacture, tvM.model, tvM.year)
}

func tvsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	brand := r.FormValue("brand")
	manufacture := r.FormValue("manufacture")
	model := r.FormValue("model")

	if brand == "" || manufacture == "" || model == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	year, err := strconv.ParseInt(r.FormValue("year"), 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	result, err := db.Exec("INSERT INTO tv VALUES(@p2, @p3, @p4, @p5)", id, brand, manufacture, model, year)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Tv %s %s created successfully (%d row affected)\n", brand, model, rowsAffected)
}
