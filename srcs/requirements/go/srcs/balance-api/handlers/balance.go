package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/balance"
)

type Balance struct {
	l *log.Logger
}

func NewBalance(l *log.Logger) *Balance {
	return &Balance{l}
}

func (bh *Balance) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		bh.l.Printf("Got a balance MethodGet request\n")
		b := &balance.Balance{}
		e := json.NewDecoder(r.Body)
		e.Decode(b)
		bh.l.Printf("bal user_id: %d\n", b.UserID)
		err := balance.GetBalanceByUserId(b.UserID, rw).ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	default:
		bh.l.Printf("Got a balance default request\n")
		http.Error(rw, "Invalid method", http.StatusMethodNotAllowed)
	}
}
