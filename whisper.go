package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func readDBPassword() string {
	b, err := os.ReadFile("dbpasswd")
	if err != nil {
		fmt.Println("Failed to read dbpasswd")
	}
	return string(b)
}

const (
	DBHost = "localhost"
	DBPort = "5423"
	DBUser = "postgres"
	DB     = "whisper"
)

var (
	DBPasswd = readDBPassword()
)

func NewUser(w http.ResponseWriter, r *http.Request) {
	c := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=enable", "localhost", 5432, "postgres", DBPasswd, "whisper")
	db, err := sql.Open("postgres", c)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()
	w.WriteHeader(http.StatusCreated)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/newuser", NewUser)
	http.ListenAndServe(":3000", nil)
}
