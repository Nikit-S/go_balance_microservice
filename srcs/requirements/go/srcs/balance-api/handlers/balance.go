package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Nikit-S/micro/balance-api/data/balance"
)

type Balance struct {
	l *log.Logger
}

func NewBalance(l *log.Logger) *Balance {
	return &Balance{l}
}

func (b *Balance) ServeHTTP(responsew http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		b.l.Printf("Got a balance MethodGet request\n")
		id := getId(responsew, request)
		//b.l.Printf("ID: %s\n", id)
		b.getBalance(id, responsew, request)
	case http.MethodPost:
		b.l.Printf("Got a balance MethodPost request\n")
	case http.MethodPut:
		b.l.Printf("Got a balance MethodPut request\n")
	default:
		b.l.Printf("Got a balance default request\n")

	}

	//fmt.Fprintf(responsew, "You have sent me a balance request with body: %s\n", body)
}

func getId(responsew http.ResponseWriter, request *http.Request) int {
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(request.URL.Path, -1)
	if len(g) != 1 || len(g[0]) != 2 {
		http.Error(responsew, "Invalid URI1", http.StatusBadRequest)
		return -1
	}
	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(responsew, "Invalid URI2", http.StatusBadRequest)
		return -1
	}

	return id
}

func (b *Balance) getBalanceList(responsew http.ResponseWriter, request *http.Request) {
	balancelist := balance.GetBalanceList(responsew)

	err := balancelist.ToJSON(responsew)
	if err != nil {
		http.Error(responsew, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (b *Balance) getBalance(id int, responsew http.ResponseWriter, request *http.Request) {
	balancelist := balance.GetBalance(id, responsew)
	if balancelist == nil {
		return
	}
	err := balancelist.ToJSON(responsew)
	if err != nil {
		http.Error(responsew, "Unable to marshal json", http.StatusInternalServerError)
	}
}
