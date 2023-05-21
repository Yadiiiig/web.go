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

	Content   []byte
	Formatted string

	Vars        []Variable
	Collections []string

	Functions []Fn
}

type Variable struct {
	Value string
	Start int
	End   int
}

type Action func() (map[int]interface{}, error)

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
	s.Formatted = string(s.Content)

	ind := [][]int{}
	for _, v := range s.Vars {
		ind = append(ind, []int{v.Start, v.End})
	}

	s.Formatted = format(s.Formatted, ind)

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

func format(input string, indexes [][]int) string {
	sortIndexes(indexes)

	for _, index := range indexes {
		startIndex := index[0]
		endIndex := index[1]

		if startIndex < 0 || endIndex >= len(input) || startIndex > endIndex {
			continue
		}

		replacement := "%v"
		input = input[:startIndex] + replacement + input[endIndex+1:]
	}

	return input
}

func sortIndexes(indexes [][]int) {
	for i := 0; i < len(indexes)-1; i++ {
		for j := 0; j < len(indexes)-i-1; j++ {
			if indexes[j][0] < indexes[j+1][0] {
				indexes[j], indexes[j+1] = indexes[j+1], indexes[j]
			}
		}
	}
}
