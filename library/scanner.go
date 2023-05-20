package library

import (
	"fmt"
	"log"
	"os"
)

func Scan(location string, actions ...Action) error {
	tmp, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	fns := MapFunctions(actions)

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

		file.Parse(fns)
		PrintStructure(file.Internal)
	}

	return nil
}
