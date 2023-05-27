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

	index := ""
	files := []File{}

	// add index.html check, otherwise add default template
	// make sure to add a formatter inbetween body tags
	for _, e := range tmp {
		name := e.Name()[:len(e.Name())-5]

		content, err := os.ReadFile(fmt.Sprintf("%s/%s", location, e.Name()))
		if err != nil {
			return files, err
		}

		if name == "index" {
			tmp, err := ParseIndex(content)
			if err != nil {
				return files, err
			}

			index = tmp

			continue
		}

		file := File{
			Name:    name,
			Content: content,
		}

		file.Parse()
		files = append(files, file)
	}

	for k, file := range files {
		output := []interface{}{}
		for _, v := range file.Internal.Formatters {
			if v.Var {
				output = append(output, fmt.Sprintf(token, v.Name))
			} else {
				output = append(output, fmt.Sprintf("%s()", v.Name))
			}
		}

		files[k].Internal.Formatted = fmt.Sprintf(file.Internal.Formatted, output...)

		_ = index
		// surround with index html
		err = writeHTML(file.Name, files[k].Internal.Formatted)
		if err != nil {
			return files, err
		}
	}

	if dev {
		return files, nil
	} else {
		return files, write("build", files)
	}
}
