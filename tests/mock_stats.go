package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", StatsServer)
	http.ListenAndServe(":8080", nil)
}

func StatsServer(w http.ResponseWriter, r *http.Request) {
	xmlFile, err := os.Open("tests/stats.xml")
	if err != nil {
		fmt.Printf("Error %s", err)
		os.Exit(1)
	}

	defer xmlFile.Close()

	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(data)
}
