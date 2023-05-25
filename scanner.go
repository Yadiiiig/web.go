package library

import (
	"fmt"
	"log"
	"os"
)

func Scan(location string, dev bool) ([]File, error) {
	tmp, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	files := []File{}

	for _, e := range tmp {
		name := e.Name()[:len(e.Name())-5]

		content, err := os.ReadFile(fmt.Sprintf("%s/%s", location, e.Name()))
		if err != nil {
			return files, err
		}

		file := File{
			Name:    name,
			Content: content,
		}

		file.Parse()
		files = append(files, file)
	}

	if dev {
		return files, nil
	} else {
		return files, write("build", files)

	}
}
