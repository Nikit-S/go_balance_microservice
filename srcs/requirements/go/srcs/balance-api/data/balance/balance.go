package balance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Nikit-S/micro/balance-api/db"
	"github.com/shopspring/decimal"
)

type Balance struct {
	ID      int             `json:"id" validate:"-"`
	UserID  int             `json:"user_id" validate:"numeric,gte=1"`
	Balance decimal.Decimal `json:"balance" validate:"numeric"`
}

type Balances []*Balance

func (b *Balances) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

func (b *Balances) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(b)
}

func (b *Balance) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

func (b *Balance) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(b)
}

func GetBalanceList(responsew http.ResponseWriter) *Balances {
	qu := "SELECT * FROM avito.balance"
	db.DB.L.Println("Querry:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer res.Close()
	bs := &Balances{}
	b := &Balance{}
	for res.Next() {

		err = res.Scan(&b.ID, &b.UserID, &b.Balance)
		if err != nil {
			db.DB.L.Println("Scan: ", err.Error())
			http.Error(responsew, err.Error(), http.StatusBadRequest)
			return nil
		}
		*bs = append(*bs, b)
	}
	return bs
}

func GetBalance(id int, responsew http.ResponseWriter) *Balance {
	qu := "SELECT * FROM avito.balance WHERE id = " + strconv.Itoa(id)
	db.DB.L.Println("Querry:", qu)
	tx, err := db.DB.Database.Begin()

	if err != nil {
		db.DB.L.Println("Begin:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer tx.Rollback()

	b := &Balance{}
	res, err := tx.Query(qu) //timeout retry
	if err != nil {
		db.DB.L.Println("Err Query:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}
	res.Next()
	err = res.Scan(&b.ID, &b.UserID, &b.Balance)
	res.Close()
	if err != nil {
		db.DB.L.Println("Scan: ", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	err = tx.Commit()
	if err != nil {
		db.DB.L.Println("Commit: ", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}
	return b

}

func GetBalanceByUserId(id int, responsew http.ResponseWriter) *Balance {
	qu := fmt.Sprintf("SELECT * FROM avito.balance WHERE user_id = %d", id)
	db.DB.L.Println("Query:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Err Query:", err.Error())
		return nil
	}
	defer res.Close()
	b := &Balance{}
	res.Next()
	err = res.Scan(&b.ID, &b.UserID, &b.Balance)

	if err != nil {
		db.DB.L.Println("Scan: ", err.Error())
		return nil
	}
	return b
}

func GetBalanceByUserIdTx(id int, responsew http.ResponseWriter, tx *sql.Tx) *Balance {
	qu := fmt.Sprintf("SELECT * FROM avito.balance WHERE user_id = %d FOR UPDATE", id)
	db.DB.L.Println("Querry:", qu)
	res, err := tx.Query(qu)
	if err != nil {
		db.DB.L.Println("Err Querry:", err.Error())
		return nil
	}

	b := &Balance{}
	res.Next()
	err = res.Scan(&b.ID, &b.UserID, &b.Balance)
	res.Close()
	if err != nil {
		db.DB.L.Println("Err Scan: ", err.Error())
		return nil
	}
	return b
}

func AddBalance(b decimal.Decimal, userid int, tx *sql.Tx) error {
	query := fmt.Sprintf("INSERT INTO avito.balance (user_id, balance) VALUES(%d, %s);", userid, b.String())
	db.DB.L.Println("Querry:", query)
	res, err := tx.Query(query)
	defer res.Close()
	if err != nil {
		db.DB.L.Println("Err Query:", err.Error())
		return err
	}
	return nil
}

func (b *Balance) Update(tx *sql.Tx) error {
	query := fmt.Sprintf("UPDATE avito.balance SET balance=%s WHERE id=%d AND user_id=%d;", b.Balance.String(), b.ID, b.UserID)
	res, err := tx.Query(query)
	if err != nil {
		db.DB.L.Println("Err Query:", err.Error())
		return err
	}
	res.Close()
	return nil
}

var ErorTransactionLimit = fmt.Errorf("Transaction amount exceeded balance amount")
var ErorWrongUserId = fmt.Errorf("Transaction amount exceeded balance amount")

func getNextID() int {
	if len(BalanceList) == 0 {
		return 1
	}
	lp := BalanceList[len(BalanceList)-1]
	return lp.ID + 1
}

var BalanceList = Balances{}
