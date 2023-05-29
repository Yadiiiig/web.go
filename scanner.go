package library

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Scan(location string, dev bool) ([]File, error) {
	var tmp []os.DirEntry
	var err error

	if location == "." {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		tmp, err = os.ReadDir(path)
	} else {
		tmp, err = os.ReadDir(location)
	}

	if err != nil {
		log.Fatal(err)
	}

	index := ""
	files := []File{}

	// add index.html check, otherwise add default template
	// make sure to add a formatter inbetween body tags
	for _, e := range tmp {
		if e.IsDir() || !strings.Contains(e.Name(), ".html") || e.Name() == "index.html" {
			continue
		}

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

		err = file.Parse()
		if err != nil {
			return files, err
		}

		files = append(files, file)
	}

	if index == "" {
		index = defaultIndex
	}

	for k, file := range files {
		output := []interface{}{}
		for _, v := range file.Internal.Formatters {
			if v.Var {
				output = append(output, fmt.Sprintf(token, v.Name))
			} else {
				output = append(output, fmt.Sprintf(onClick, v.Name))
			}
		}

		files[k].Internal.Formatted = fmt.Sprintf(file.Internal.Formatted, output...)

		_ = index
		// surround with index html
		err = writeHTML(file.Name, fmt.Sprintf(index, file.Name, files[k].Internal.Formatted))
		if err != nil {
			return files, err
		}
	}

	if dev {
		return files, nil
	} else {
		return files, write("runtime", files)
	}
}
