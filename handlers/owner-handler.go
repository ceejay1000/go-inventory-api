package handlers

import (
	"net/http"
)

type Owner struct{}

func (own *Owner) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Owner route reached"))

	// httpMethod := r.Method

}
