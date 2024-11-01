// hellogo.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func getFrontpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go!")
}

func main() {
	http.HandleFunc("/", getFrontpage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
