package main

import (
	"log"
	"net/http"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

type defaultHandler struct{}

func (dh *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my inventory app"))
}

func main() {

	router := http.NewServeMux()

	router.Handle("/", &defaultHandler{})

	server := http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	log.Println("Server started and listening on PORT :9090")
	if err := server.ListenAndServe(); err != nil {
		log.Panic("Internal Server Error")
	}
}
