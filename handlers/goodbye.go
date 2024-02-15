package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type GoodBye struct {
	l *log.Logger
}

func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

func (gb *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	gb.l.Println("GoodBye")
	fmt.Fprintln(rw, "GoodBye from", r.URL.Path)
}
