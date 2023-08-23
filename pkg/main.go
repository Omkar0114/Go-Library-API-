package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	_"github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)


const (
	API_PATH = "/apis/v1/books"
)

type library struct{
	dbHost, dbPass, dbName string
}

type Book struct{
	ID, Name, ISBN string
}

func main(){
	
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == ""{
		dbPass = "omkark"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == ""{
		dbName = "library"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == ""{
		apiPath = API_PATH
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}


	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getbooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postbook).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)

	
}

func (l library) postbook(w http.ResponseWriter, r *http.Request){
	// read the request from an instance of the book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)

	// open connection
	db := l.openConnection()
	// write the data into database
	insertQuery, err := db.Prepare("INSERT INTO books VALUES(?, ?, ?)")
	if err != nil{
		log.Fatalf("Preparing the db query %s\n", err.Error())
	}

	tx, err := db.Begin()
	if err != nil{
		log.Fatalf("while begining the transactions%s\n", err.Error())

	}

	_, err = tx.Stmt(insertQuery).Exec(book.ID, book.Name, book.ISBN)
	if err != nil{
		log.Fatalf("while executing the statement%s\n", err.Error())
	}

	err = tx.Commit()
	if err != nil{
		log.Fatalf("While commiting the transaction%s\n", err.Error())
	}

	//close the connection
	defer l.closeconnection(db)
}



func (l library) getbooks(w http.ResponseWriter, r *http.Request){
	
	// open connection
	db := l.openConnection()

	// read all the books 
	rows, err := db.Query("SELECT * FROM books")
	if err != nil{
		log.Fatalf("Error while reading books %s\n", err.Error())
	}

	books := []Book{}
	// iterate over the rows
	for rows.Next(){
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("While scanning the row%s\n", err.Error())
		}
		aBook := Book{
			ID: id, 
			Name: name,
			ISBN: isbn,
		}
		books = append(books, aBook)

	}
	// write the response using encoding/json
	json.NewEncoder(w).Encode(books)

	// close the connection
	defer l.closeconnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s" , "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil{
		log.Fatalf("Error while opening connection %s/", err.Error())
	}

	return db

}

func (l library) closeconnection(db *sql.DB){
	err := db.Close()
	if err != nil {
		log.Fatalf("Error while closing the connection%s\n", err.Error())
	}
}