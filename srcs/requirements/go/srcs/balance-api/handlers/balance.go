package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nikit-S/micro/balance-api/data/balance"
	"github.com/go-playground/validator"
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
		err := bh.validateRequest(b)
		if err != nil {
			http.Error(rw, "Wrong json format", http.StatusBadRequest)
			return
		}
		bh.l.Printf("bal user_id: %d\n", b.UserID)
		err = balance.GetBalanceByUserId(b.UserID, rw).ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	default:
		bh.l.Printf("Got a balance default request\n")
		http.Error(rw, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func (bh *Balance) validateRequest(b *balance.Balance) error {
	validate := validator.New()
	err := validate.Struct(b)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		bh.l.Println("------ List of tag fields with error ---------")

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
