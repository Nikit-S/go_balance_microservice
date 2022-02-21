package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Basic struct {
	l *log.Logger
}

func NewBasic(l *log.Logger) *Basic {
	return &Basic{l}
}

func (basic *Basic) ServeHTTP(responsew http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		basic.l.Println("Error reading body", err)

		http.Error(responsew, "Unable to read request body", http.StatusBadRequest)
		return
	}
	basic.l.Printf("Got a basic request with body: %s\n", body)
	fmt.Fprintf(responsew, "You have sent me a request with body: %s\n", body)
}
