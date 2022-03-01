package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/balance"
	"github.com/Nikit-S/micro/balance-api/data/transaction"
	"github.com/Nikit-S/micro/balance-api/db"
	"github.com/go-playground/validator"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	l *log.Logger
}

func NewTransaction(l *log.Logger) *Transaction {
	return &Transaction{l}
}

func (th *Transaction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		th.l.Printf("Got a transaction MethodGet request\n")
		t := &transaction.Transaction{}
		err := json.NewDecoder(r.Body).Decode(t)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError) //http.StatusBadRequest
		}
		err = th.validateRequest(t)
		if err != nil {
			http.Error(rw, "Wrong json format", http.StatusBadRequest)
			return
		}
		th.getTransaction(t.ID, rw, r)
	case http.MethodPost:
		th.l.Printf("Got a transaction MethodPost request\n")
		t := th.addTransaction(rw, r)
		if t == nil {
			return
		}
		tx, err := db.DB.Database.Begin()
		if err != nil {
			db.DB.L.Println("Begin:", err.Error())
			http.Error(rw, "Connection to mysql:", http.StatusInternalServerError)
		}
		err = th.execTransaction(t, rw, tx)
		if err != nil {
			db.DB.L.Println("Transfer:", err.Error())
		}
		t.ToJSON(rw)
	default:
		th.l.Printf("Got a transaction default request\n")
		http.Error(rw, "Invalid method", http.StatusMethodNotAllowed) //http.StatusBadRequest
	}
}

func (th *Transaction) validateRequest(t *transaction.Transaction) error {
	validate := validator.New()
	err := validate.Struct(t)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		th.l.Println("------ List of tag fields with error ---------")

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.StructField())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------")
		}
		return err
	}
	return nil
}

func (th *Transaction) addTransaction(rw http.ResponseWriter, r *http.Request) *transaction.Transaction {
	th.l.Println("POST (add Transaction)")
	t := &transaction.Transaction{}
	err := t.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to UNmarshal json", http.StatusBadRequest)
		return nil
	}
	err = th.validateRequest(t)
	if err != nil {
		http.Error(rw, "Wrong json format", http.StatusBadRequest)
		return nil
	}
	th.l.Println("body is valid")
	transaction.AddTransaction(t, rw)
	return t

}

func (t *Transaction) getTransaction(id int, responsew http.ResponseWriter, request *http.Request) {
	transactionlist := transaction.GetTransaction(id, responsew)
	if transactionlist == nil {
		http.Error(responsew, "Id is out of ramge", http.StatusInternalServerError)
	}
	err := transactionlist.ToJSON(responsew)
	if err != nil {
		http.Error(responsew, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (th *Transaction) execTransaction(t *transaction.Transaction, rw http.ResponseWriter, tx *sql.Tx) error {
	bt := balance.GetBalanceByUserIdTx(t.UserID, rw, tx)
	th.l.Println("Exec transaction")
	if t.From == "balance" {
		if t.UserID == t.FromID {
			th.l.Println("Self balance")
			tx.Commit()
			return nil
		}
		err := th.fromBalance(t, rw, tx)
		if err != nil {
			return err
		}
	}
	if bt == nil && t.Amount.Cmp(decimal.Zero) == +1 {
		err := th.createBalance(t, rw, tx)
		if err != nil {
			return err
		}
		return nil
	} else if bt != nil && !decimal.Sum(t.Amount, bt.Balance).LessThan(decimal.Zero) {
		err := th.makeTransaction(t, rw, tx, bt)
		if err != nil {
			return err
		}
		return nil
	}
	tx.Rollback()
	return fmt.Errorf("something went wrong")
}

func (th *Transaction) makeTransaction(t *transaction.Transaction, rw http.ResponseWriter, tx *sql.Tx, bt *balance.Balance) error {
	bt.Balance = bt.Balance.Add(t.Amount)
	err := bt.Update(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	t.Status = 1
	err = t.Update()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (th *Transaction) createBalance(t *transaction.Transaction, rw http.ResponseWriter, tx *sql.Tx) error {
	th.l.Println("Gonna make a balance")
	err := balance.AddBalance(t.Amount, t.UserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	t.Status = 1
	err = t.Update()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (th *Transaction) fromBalance(t *transaction.Transaction, rw http.ResponseWriter, tx *sql.Tx) error {
	bf := balance.GetBalanceByUserIdTx(t.FromID, rw, tx)
	if bf == nil {
		http.Error(rw, "Get balance:", http.StatusNotFound)
		tx.Commit()
		return fmt.Errorf("No balance %d", t.UserID)
	}
	t.Amount = t.Amount.Abs()
	if bf.Balance.Cmp(t.Amount) == -1 {
		http.Error(rw, "Sender Balance < Transaction Amount:", http.StatusForbidden)
		tx.Commit()
		return fmt.Errorf("Sender Balance < Transaction Amount")
	}
	bf.Balance = bf.Balance.Sub(t.Amount)
	err := bf.Update(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
