package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

func addr() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":8888"
	} else {
		port = fmt.Sprintf(":%s", port)
	}
	return port
}

func renderIndex(w http.ResponseWriter) error {
	tplt, err := template.ParseFiles("./templates/index.template")
	if err != nil {
		return err
	}
	errExec := tplt.Execute(w, nil)
	return errExec
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	if err := renderIndex(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.ListenAndServe(addr(), r)
}
