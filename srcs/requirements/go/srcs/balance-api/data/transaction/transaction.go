package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Nikit-S/micro/balance-api/db"
)

type Transaction struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Amount  int    `json:"amount"`
	Status  int    `json:"status"`
	From    string `json:"from"`
	FromID  int    `json:"from_id"`
	Comment string `json:"comment"`
}

type Transactions []*Transaction

var TransactionList = Transactions{}

func (t *Transactions) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Transaction) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (t *Transaction) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Transactions) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func GetTransactionList() *Transactions {
	return &TransactionList
}

func GetTransaction(id int, responsew http.ResponseWriter) *Transaction {
	qu := "SELECT * FROM avito.transaction WHERE id = " + strconv.Itoa(id)
	db.DB.L.Println("Querry:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer res.Close()
	t := &Transaction{}
	res.Next()
	err = res.Scan(&t.UserID, &t.ID, &t.Amount, &t.Status, &t.From, &t.FromID, &t.Comment)
	if err != nil {
		db.DB.L.Println("Scan: ", err.Error())
		http.Error(responsew, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return t
}

func GetTransactionsByUserId(id int, responsew http.ResponseWriter) *Transactions {
	qu := "SELECT * FROM avito.transaction WHERE user_id = " + strconv.Itoa(id)
	db.DB.L.Println("Querry:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer res.Close()
	t := &Transaction{}
	ts := &Transactions{}
	for res.Next() {
		err = res.Scan(&t.ID, &t.UserID, &t.Amount, &t.Status, &t.From, &t.FromID, &t.Comment)
		if err != nil {
			db.DB.L.Println("Scan: ", err.Error())
			http.Error(responsew, err.Error(), http.StatusInternalServerError)
			return nil
		}
		*ts = append(*ts, t)
	}

	return ts
}

func AddTransaction(trans *Transaction, responsew http.ResponseWriter) error {

	query := fmt.Sprintf("INSERT INTO avito.transaction (user_id, amount, status, `from`, from_id, comment) VALUES(%d, %d, %d, '%s', %d, '%s');", trans.UserID, trans.Amount, trans.Status, trans.From, trans.FromID, trans.Comment)
	db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Exec(query)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		db.DB.L.Println("Last id:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return err
	}
	db.DB.L.Printf("Result: %s\n", id)
	trans.ID = int(id)
	return nil
}

func (trans *Transaction) getNextID() int {
	if len(TransactionList) == 0 {
		return 1
	}
	lp := TransactionList[len(TransactionList)-1]
	return lp.ID + 1
}

func (t *Transaction) Update(responsew http.ResponseWriter) {
	query := fmt.Sprintf("UPDATE avito.transaction SET status=%d WHERE id=%d AND user_id=%d;", t.Status, t.ID, t.UserID)
	db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Query(query)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
	}
	defer res.Close()
}
