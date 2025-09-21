package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
	"github.com/alphaloan/vehicle/handler"
)

func main() {
	//datastore.InitializeDatabase("db/migration", "sqlite3://alphaloan.db")

	db, err := sql.Open("sqlite3", "alphaloan.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign_keys ", err)
	}

	defer db.Close()

	loanCustomerStore := datastore.NewLoanCustomerStore(db)
	loanSubmissionStore := datastore.NewLoanSubmissionStore(db)

	loanSubmitHandler := handler.NewLoanSubmitHandler(*loanCustomerStore, *loanSubmissionStore)
	loanSubmissionHandler := handler.NewLoanSubmissionHandler(*loanSubmissionStore)
	loanCustomerHandler := handler.NewLoanCustomerHandler(*loanCustomerStore)

	http.HandleFunc("/api/loan/submit", loanSubmitHandler.HandleSubmitLoan)
	http.HandleFunc("/api/loan/submissions", loanSubmissionHandler.HandleGetAllLoanSubmission)
	http.HandleFunc("/api/loan/submission/tracks", loanSubmissionHandler.HandleSubmissionLoanById)
	http.HandleFunc("/api/loan/customers", loanCustomerHandler.HandleGetAllLoanSubmission)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
