package library

import (
	"fmt"
	"log"
	"net/http"
)

func Start(args []string, fns []Function, acts []Action) {
	var files []File
	var err error

	if len(args) < 2 {
		fmt.Println("No parameters specified")
	} else if args[1] == "gen" {
		files, err := Scan(args[2], false)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			PrintStructure(file.Internal)
		}

	} else if args[1] == "run" {
		files, err = read(args[2])
		if err != nil {
			log.Fatal(err)
		}
	} else if args[1] == "dev" {
		files, err = Scan(args[2], true)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(":sunglas:")
	}

	fnsm, actsm := Map(fns, acts)
	fmt.Println(fnsm, actsm)
	for k := range files {
		files[k].Add(fnsm, actsm)
		GenerateEndpoints(&files[k])
		PrintStructure(files[k].Internal)
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
