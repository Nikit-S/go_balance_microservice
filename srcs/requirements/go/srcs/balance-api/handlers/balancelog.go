package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/balance"
	"github.com/Nikit-S/micro/balance-api/data/transaction"
)

type BalanceLog struct {
	l *log.Logger
}

func NewBalanceLog(l *log.Logger) *BalanceLog {
	return &BalanceLog{l}
}

func (blh *BalanceLog) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		blh.l.Printf("Got a balancelog MethodGet request\n")
		bal := &balance.Balance{}
		e := json.NewDecoder(r.Body)
		e.Decode(bal)
		err := transaction.GetTransactionsByUserId(bal.UserID, rw).ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	default:
		blh.l.Printf("Got a balance default request\n")
		http.Error(rw, "Invalid method", http.StatusMethodNotAllowed)
	}
}
