package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Nikit-S/micro/balance-api/db"
)

type Balance struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	Balance int `json:"balance"`
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
			http.Error(responsew, err.Error(), http.StatusInternalServerError)
			return nil
		}
		*bs = append(*bs, b)
	}
	return bs
}

func GetBalance(id int, responsew http.ResponseWriter) *Balance {
	qu := "SELECT * FROM avito.balance WHERE id = " + strconv.Itoa(id)
	db.DB.L.Println("Querry:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer res.Close()
	b := &Balance{}
	res.Next()
	err = res.Scan(&b.ID, &b.UserID, &b.Balance)
	if err != nil {
		db.DB.L.Println("Scan: ", err.Error())
		http.Error(responsew, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return b

}

func GetBalanceByUserId(id int, responsew http.ResponseWriter) *Balance {
	qu := "SELECT * FROM avito.balance WHERE user_id = " + strconv.Itoa(id)
	db.DB.L.Println("Querry:", qu)
	res, err := db.DB.Database.Query(qu)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		//http.Error(responsew, err.Error(), http.StatusBadRequest)
		return nil
	}

	defer res.Close()
	b := &Balance{}
	res.Next()
	err = res.Scan(&b.ID, &b.UserID, &b.Balance)
	if err != nil {
		db.DB.L.Println("Scan: ", err.Error())
		//http.Error(responsew, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return b
}

func AddBalance(b int, userid int, responsew http.ResponseWriter) {
	query := fmt.Sprintf("INSERT INTO avito.balance (user_id, balance) VALUES(%d, %d);", userid, b)
	db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Query(query)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
	}

	defer res.Close()

}

func (b *Balance) Update(responsew http.ResponseWriter) {
	query := fmt.Sprintf("UPDATE avito.balance SET balance=%d WHERE id=%d AND user_id=%d;", b.Balance, b.ID, b.UserID)
	db.DB.L.Println("Querry:", query)
	res, err := db.DB.Database.Query(query)
	if err != nil {
		db.DB.L.Println("Querry:", err.Error())
		http.Error(responsew, err.Error(), http.StatusBadRequest)
	}
	defer res.Close()
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

var BalanceList = Balances{
	&Balance{
		ID:      1,
		Balance: 10000,
		UserID:  234,
	},
	&Balance{
		ID:      2,
		Balance: 20000,
		UserID:  24,
	},
}
