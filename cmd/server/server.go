package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var up = flag.String("up", "", "path for upload file")
var port = flag.Int("p", 8001, "server port")

// TODO daemon

func main() {
	flag.Parse()
	log.Printf("path for upload file: %s", *up)

	startServer()
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	ff, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}
	defer ff.Close()
	log.Printf("Upload file: %s", header.Filename)

	newFile, err := os.Create(*up + "/" + header.Filename)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, ff)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintf(w, "OK")
}

func startServer() {
	http.HandleFunc("/upload", uploadFile)

	addr := fmt.Sprintf("localhost:%d", *port)
	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}
