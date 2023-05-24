package library

import (
	"fmt"
	"log"
	"os"
)

func Scan(location string) error {
	tmp, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	// fnsm := MapFunctions(fns)
	// actsm := MapActions(acts)

	files := []File{}
	for _, e := range tmp {
		name := e.Name()[:len(e.Name())-5]
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", location, e.Name()))
		if err != nil {
			return err
		}

		file := File{
			Name:    name,
			Content: content,
		}

		file.Parse()
		files = append(files, file)
		// PrintStructure(file.Internal)

		//GenerateEndpoints(&file)
	}

	return write("build", files)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
