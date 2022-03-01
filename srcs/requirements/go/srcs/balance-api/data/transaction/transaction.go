package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/db"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID      int             `json:"id" validate:"-"`
	UserID  int             `json:"user_id" validate:"required,numeric,gte=1"`
	Amount  decimal.Decimal `json:"amount" validate:"required,numeric"`
	Status  int             `json:"status" validate:"-"`
	From    string          `json:"from" validate:"required"`
	FromID  int             `json:"from_id" validate:"-"`
	Comment string          `json:"comment" validate:"-"`
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

func GetTransaction(id int, responsew http.ResponseWriter) *Transaction {
	qu := fmt.Sprintf("SELECT * FROM avito.transaction WHERE id = %d;", id)
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
		return nil
	}
	return t
}

func GetTransactionsByUserId(id int, responsew http.ResponseWriter) *Transactions {
	qu := fmt.Sprintf("SELECT * FROM avito.transaction WHERE user_id = %d;", id)
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
			http.Error(responsew, err.Error(), http.StatusBadRequest)
			return nil
		}
		*ts = append(*ts, t)
	}

	return ts
}

func AddTransaction(trans *Transaction, responsew http.ResponseWriter) error {

	//db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Exec("INSERT INTO avito.transaction (user_id, amount, status, `from`, from_id, comment) VALUES(?, ?, ?, ?, ?, ?);", trans.UserID, trans.Amount.String(), trans.Status, trans.From, trans.FromID, trans.Comment)
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
	db.DB.L.Printf("Result: %d\n", id)
	trans.ID = int(id)
	return nil
}

func (t *Transaction) Update() error {
	query := fmt.Sprintf("UPDATE avito.transaction SET status=%d WHERE id=%d AND user_id=%d;", t.Status, t.ID, t.UserID)
	db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Query(query)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		return err
	}
	defer res.Close()
	return nil
}
