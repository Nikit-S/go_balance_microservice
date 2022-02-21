package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nikit-S/micro/balance-api/db"
	"github.com/Nikit-S/micro/balance-api/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	db.ConnectDB(l)

	defer db.DB.Database.Close()
	servemux := http.NewServeMux()

	balancehandler := handlers.NewBalance(l)
	balanceloghandler := handlers.NewBalanceLog(l)
	transactionhandler := handlers.NewTransaction(l)

	servemux.Handle("/balance/", balancehandler)
	servemux.Handle("/balancelog/", balanceloghandler)
	servemux.Handle("/transaction/", transactionhandler)

	server := &http.Server{
		Addr:         ":9090",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      servemux,
	}
	server.ListenAndServe()
}
