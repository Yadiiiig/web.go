package library

import (
	"fmt"
	"log"
	"net/http"
)

func Start(args []string, fns []Function, acts []Action) {
	var files []File
	var err error

	settings, err := BuildSettings()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < 2 {
		fmt.Println("No parameters specified")
	} else if args[1] == "gen" {
		files, err := Scan(args[2], false)
		if err != nil {
			log.Fatal(err)
		}

		err = settings.GenLibrary(files)
		if err != nil {
			log.Fatal(err)
		}

		return

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

	for k := range files {
		files[k].Add(fnsm, actsm)
		GenerateEndpoints(&files[k])
	}

	log.Fatal(http.ListenAndServe(settings.Endpoint, nil))
}
