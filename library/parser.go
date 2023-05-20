package library

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

type File struct {
	Name     string
	Internal Structure
	Content  []byte
}

type Structure struct {
	Start int
	End   int

	Content []byte

	Vars        []Variable
	Collections []string

	Functions []Fn
}

type Variable struct {
	Value string
	Start int
	End   int
}

type Action func() (string, error)

type Fn struct {
	Run  Action
	Path string

	VarsOrder []int
}

// map functions so they can be assigned to files
func MapFunctions(actions []Action) map[string]Fn {
	m := make(map[string]Fn)

	for _, v := range actions {
		path := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
		name := strings.ToLower(path[strings.LastIndex(path, ".")+1:])

		m[name] = Fn{
			Run:  v,
			Path: path,
		}
	}

	return m
}

/*
1 assign start & end
2 cut content
3 retrieve structures
4 find variables
5 add formatters in content
6 check for correct order
7 create function link
*/
func (f *File) Parse(fns map[string]Fn) {
	s := Structure{}

	for k, v := range f.Content {
		if v == '<' {
			s.Start, s.End = 0, k-1
			break
		}
	}

	s.Content = f.Content[s.End+1:]

	value := ""
	cols := []string{}

	for _, v := range f.Content[s.Start:s.End] {
		if v == '\n' {
			cols = append(cols, value)
			value = ""
		} else if v == '-' {
			value = ""
		} else {
			value += string(v)
		}
	}

	vars := []Variable{}
	start, value, found := 0, "", false

	for k, v := range s.Content {
		if v == '{' {
			start = k
			found = true
		} else if v == '}' {
			vars = append(vars, Variable{
				value,
				start,
				k + 1,
			})

			if !strings.Contains(value, ".") {
				cols = append(cols, value)

			}

			start = 0
			found = false
			value = ""
		} else if found {
			value += string(v)
		}
	}

	s.Vars = vars
	s.Collections = cols

	for _, v := range s.Vars {
		format(&s.Content, v.Start, v.End, "%v")
	}

	if !isOrdered(s.Vars) {
		sort.Slice(s.Vars, func(i, j int) bool {
			return s.Vars[i].Start < s.Vars[j].Start
		})
	}

	functions := []Fn{}
	for _, c := range s.Collections {
		fn := Fn{}

		for k, v := range s.Vars {
			if strings.HasPrefix(v.Value, fmt.Sprintf("%s.", c)) || strings.HasPrefix(v.Value, fmt.Sprintf("%s", c)) {
				fn = fns[c]
				fn.VarsOrder = append(fn.VarsOrder, k)
			}
		}

		functions = append(functions, fn)
	}

	s.Functions = functions
	f.Internal = s
}

func format(slice *[]byte, start, end int, newBytes string) {
	if start < 0 || start >= len(*slice) || end < start || end > len(*slice) {
		fmt.Println("Invalid start or end index")
		return
	}

	replaceLen := end - start
	newSlice := []byte(fmt.Sprintf(newBytes))

	if replaceLen == len(newSlice) {
		copy((*slice)[start:end], newSlice)
	} else if len(newSlice) < replaceLen {
		copy((*slice)[start:start+len(newSlice)], newSlice)
		copy((*slice)[start+len(newSlice):end], (*slice)[end:])

		*slice = (*slice)[:len(*slice)-(replaceLen-len(newSlice))]
	} else {
		*slice = append((*slice)[:start], append(newSlice, (*slice)[end:]...)...)
	}
}
