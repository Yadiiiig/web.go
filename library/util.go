package library

import "fmt"

func isOrdered(arr []Variable) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i].Start < arr[i-1].Start {
			return false
		}
	}
	return true
}

func PrintStructure(s Structure) {
	fmt.Printf("Start: %d\n", s.Start)
	fmt.Printf("End: %d\n", s.End)
	fmt.Printf("Content: %s\n", string(s.Content))

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
