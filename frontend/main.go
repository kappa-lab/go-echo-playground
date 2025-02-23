package main

import (
	_ "embed"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexPage)

	// 3000ポートしかCORSで許可していない
	// 3000ポート以外で起動するとCORSエラーになる
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//go:embed index.html
var indexHTML string

func indexPage(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(indexHTML))
	if err != nil {
		log.Fatal(err)
	}
}
