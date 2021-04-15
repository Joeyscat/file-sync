package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joeyscat/file-sync/internal"
)

var dir = flag.String("dir", "", "Dir to watch")
var u = flag.String("url", "", "Server addr for upload")

var d = flag.Bool("daemon", false, "Runs as a daemon ")

func main() {
	watch := internal.NewNotifyFile()

	flag.Parse()
	if dir == nil || *dir == "" {
		log.Fatal("arg dir must not be empty")
	}
	if u == nil || *u == "" {
		log.Fatal("arg url must not be empty")
	}
	_, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	watch.WatchDir(*dir, "xx")

	go func(file *internal.NotifyFile) {
		for {
			select {
			case p := <-watch.Path:
				{
					fmt.Printf("Upload: %+v\n", p.Path)
					result, err := internal.Upload(p.Path, *u)
					if err != nil {
						log.Fatal(err)
					}
					log.Println(result)
				}
			}
		}
	}(watch)

	select {}
	return
}
