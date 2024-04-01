package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	// Dump the request including headers
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Request:")
	fmt.Println(string(requestDump))
}

func main() {
	http.HandleFunc("/", getRoot)
	http.ListenAndServe(":8080", nil)
}
