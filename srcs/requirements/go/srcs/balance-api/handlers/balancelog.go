package handlers

import (
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/transaction"
)

type BalanceLog struct {
	l *log.Logger
}

func NewBalanceLog(l *log.Logger) *BalanceLog {
	return &BalanceLog{l}
}

func (bl *BalanceLog) ServeHTTP(responsew http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		bl.l.Printf("Got a balancelog MethodGet request\n")
		id := getId(responsew, request)
		//bl.l.Printf("ID: %s\n", id)
		bl.getBalanceLog(id, responsew, request)
	case http.MethodPost:
		bl.l.Printf("Got a balancelog MethodPost request\n")
	case http.MethodPut:
		bl.l.Printf("Got a balancelog MethodPut request\n")
	default:
		bl.l.Printf("Got a balancelog default request\n")

	}

	//fmt.Fprintf(responsew, "You have sent me a balance request with body: %s\n", body)
}

func (bl *BalanceLog) getBalanceLog(id int, responsew http.ResponseWriter, request *http.Request) error {
	bllist := transaction.GetTransactionsByUserId(id, responsew)
	bllist.ToJSON(responsew)
	return nil
}
