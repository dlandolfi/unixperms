package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// db needs to be a global var for the handlers to access them
var db *sql.DB
var err error

const port = "8080"

func main() {
	// note: not using short declaration b/c db needs to be global
	db, err = sql.Open("mysql", "test:passwordz@tcp(localhost:3306)/permissions")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ping db to verify connection
	// func (db *DB) Ping() error
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	// print statement to show db connection was successful
	fmt.Println("DB connect! Such leet!")

	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/perms", perms)
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	fmt.Println("Listening on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "unix permission guide ")
}

func perms(w http.ResponseWriter, req *http.Request) {
	// func (r *Request) FormValue(key string) string
	pm := req.FormValue("q") // i.e. localhost:8080/?q=whatever
	rows, err := db.Query("select perm_string from perms where perm_number=?;", pm)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var name string
	// Next prepares the next result row for reading with the Scan method.
	// It returns true on success, or false if there is no next result row or an error happened while preparing it.
	// Err should be consulted to distinguish between the two cases.
	// Every call to Scan, even the first one, must be preceded by a call to Next.
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, name)
	}
}
