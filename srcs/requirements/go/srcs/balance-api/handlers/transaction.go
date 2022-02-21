package handlers

import (
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/balance"
	"github.com/Nikit-S/micro/balance-api/data/transaction"
)

type Transaction struct {
	l *log.Logger
}

func NewTransaction(l *log.Logger) *Transaction {
	return &Transaction{l}
}

func (t *Transaction) ServeHTTP(responsew http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		t.l.Printf("Got a balance MethodGet request\n")
		id := getId(responsew, request)
		//transaction.l.Printf("ID: %s\n", id)
		t.getTransaction(id, responsew, request)
	case http.MethodPost:
		t.l.Printf("Got a balance MethodPost request\n")
		t.addTransaction(responsew, request)
	case http.MethodPut:
		t.l.Printf("Got a balance MethodPut request\n")
	default:
		t.l.Printf("Got a balance default request\n")
	}
	//fmt.Fprintf(responsew, "You have sent me a transaction request with body: %s\n", body)
}

func (trans *Transaction) addTransaction(rw http.ResponseWriter, r *http.Request) {
	trans.l.Println("POST (add Transaction)")

	product := &transaction.Transaction{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to UNmarshal json", http.StatusBadRequest)
	}

	transaction.AddTransaction(product, rw)
	//trans.l.Println("Transaction added")
	trans.execTransaction(product, rw)
	product.ToJSON(rw)
	//trans.l.Printf("Product: %#v", product.Status)

}

func (t *Transaction) getTransactionList(responsew http.ResponseWriter, request *http.Request) {
	transactionlist := transaction.GetTransactionList()

	err := transactionlist.ToJSON(responsew)
	if err != nil {
		http.Error(responsew, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Transaction) getTransaction(id int, responsew http.ResponseWriter, request *http.Request) {
	transactionlist := transaction.GetTransaction(id, responsew)

	if transactionlist == nil {
		http.Error(responsew, "Id is out of ramge", http.StatusInternalServerError)
		return
	}
	err := transactionlist.ToJSON(responsew)
	if err != nil {
		http.Error(responsew, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (t *Transaction) execTransaction(trans *transaction.Transaction, rw http.ResponseWriter) error {
	//t.l.Println("Exec trans")
	b := balance.GetBalanceByUserId(trans.UserID, rw)
	if b == nil && trans.Amount > 0 {
		t.l.Println("Gonna make a balance")
		balance.AddBalance(trans.Amount, trans.UserID, rw)
		trans.Status = 1
		trans.Update(rw)
	}
	if b != nil && trans.Amount+b.Balance >= 0 {
		//t.l.Println("Gonna make a transaction")
		b.Balance += trans.Amount
		b.Update(rw)
		trans.Status = 1
		trans.Update(rw)
	}
	return nil
}
