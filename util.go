package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func write(name string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s.json", name), file, 0644)
}

func writeHTML(name, output string) error {
	return ioutil.WriteFile(fmt.Sprintf("%s.html", name), []byte(output), 0644)

}

func read(name string) ([]File, error) {
	file, err := os.Open(fmt.Sprintf("%s.json", name))
	if err != nil {
		return []File{}, err
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []File{}, err
	}

	files := []File{}

	json.Unmarshal(b, &files)

	return files, nil
}

func isOrdered(arr []Variable) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i].Start < arr[i-1].Start {
			return false
		}
	}
	return true
}

func replaceIndices(org *string, input string, start, end int) {
	*org = fmt.Sprintf("%s%s%s", (*org)[:start], input, (*org)[end+1:])
}

func PrintStructure(s Structure) {
	fmt.Printf("Start: %d\n", s.Start)
	fmt.Printf("End: %d\n", s.End)
	fmt.Printf("Content: %s\n", s.Formatted)

	fmt.Println("Collections:")
	for _, v := range s.Collections {
		fmt.Printf("\tName: %s\n", v)
	}

	fmt.Println("Vars:")
	for _, v := range s.Vars {
		fmt.Printf("\tValue: %s\n", v.Value)
		fmt.Printf("\tStart: %d\n", v.Start)
		fmt.Printf("\tEnd: %d\n", v.End)
	}

	fmt.Println("Functions:")
	for _, v := range s.Functions {
		fmt.Printf("\tPath: %s\n", v.Path)
		fmt.Printf("\tVarsOrder: %v\n", v.VarsOrder)
	}

	fmt.Println("Requests:")
	for _, v := range s.Requests {
		fmt.Printf("\tName: %s\n", v.Name)
		fmt.Printf("\tParams: %v\n", v.Params)
	}

}

func PrintContentCheck(s Structure, content string) {
	//fmt.Println(string(content))
	fmt.Println(content)
	content = content[s.End+1:]
	fmt.Println(s.End)
	for _, v := range s.Vars {
		fmt.Printf("%s\n", string(content[v.Start:v.End]))
	}
}
