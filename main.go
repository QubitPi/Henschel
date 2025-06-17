package main

import (
	"flag"
	"net/http"
)

func main() {
	port := flag.String("port", ":8080", "Webservice port; default to 8080")
	flag.Parse()
	http.ListenAndServe(*port, nil)
}
