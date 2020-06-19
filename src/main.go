package main

import (
	_ "github.com/pm2go"
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/index", func(rw http.ResponseWriter, req *http.Request){
		rw.Write([]byte("hello, pm2go!\n"))
	})
	log.Fatalln(http.ListenAndServe(":9000", mux ))
}
