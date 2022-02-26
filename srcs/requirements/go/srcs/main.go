package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Nikit-S/micro/balance-api/db"
	"github.com/Nikit-S/micro/balance-api/handlers"
)

func main() {

	l := log.New(os.Stdout, "balance-api ", log.LstdFlags)

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
		Addr:         os.Getenv("SERVICE_HOST") + ":" + os.Getenv("SERVICE_PORT"),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      servemux,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Gracefull termination", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
