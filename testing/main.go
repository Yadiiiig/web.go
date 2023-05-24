package main

import (
	"fmt"
	"log"
	"os"
	web "webmaker/library"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No parameters specified")
	} else if os.Args[1] == "gen" {
		err := web.Scan(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

	} else if os.Args[1] == "run" {

	} else {

	}
	//	library.GenFile("http://localhost:8080/foo", library.GenRequest("remove", "POST", "http://localhost:8080/foo/remove", "follow", []string{"user.id", "user.last_visit"}))
}
