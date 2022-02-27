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

	l := log.New(os.Stdout, "balance-api ", log.LstdFlags) //создание логирования

	db.ConnectDB(l) // подключение к базе данных

	defer db.DB.Database.Close()
	servemux := http.NewServeMux()

	balancehandler := handlers.NewBalance(l)
	balanceloghandler := handlers.NewBalanceLog(l)
	transactionhandler := handlers.NewTransaction(l)

	servemux.Handle("/balance/", balancehandler)         //прослушка для метода запроса баланса
	servemux.Handle("/balancelog/", balanceloghandler)   //прослушка для метода запроса списка транзакций
	servemux.Handle("/transaction/", transactionhandler) //прослушка для метода проведения тразакций

	server := &http.Server{ //настройки сервера
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
	sigChan := make(chan os.Signal) //плавное отключение за 30 секунд
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Gracefull termination", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
