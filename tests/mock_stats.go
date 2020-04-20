// MIT License

// Copyright (c) 2020 Mauricio Antunes

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", StatsServer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// StatsServer is a mock server to expose dummy NGINX-RTMP stats
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
	if _, err := w.Write(data); err != nil {
		log.Fatal(err)
	}
}
