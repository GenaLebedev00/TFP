package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	if install {
		_, err = db.Exec(`CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title VARCHAR(32) NOT NULL,
			comment VARCHAR(128) NOT NULL,
			repeat VARCHAR(128) NOT NULL
		);`)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("CREATE INDEX dates_task ON scheduler (date)")
		if err != nil {
			panic(err)
		}
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		fmt.Println(err)
	}
}
