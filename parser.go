package library

import (
	"fmt"
	"net/http"
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
	Start, End int

	Content []byte

	Formatted  string
	Formatters []Formatter

	Vars        []Variable
	Collections []string

	Functions []Fn `json:"-"`
	Requests  []Request
}

type Request struct {
	Name string
	Run  Action `json:"-"`

	Start, End int

	Params []string
}

type Variable struct {
	Value string
	Start int
	End   int
}

type Function func(w http.ResponseWriter, r *http.Request) (string, interface{}, error)
type Action func(params map[string]interface{}, w http.ResponseWriter, r *http.Request) (string, interface{}, error)

type Fn struct {
	Run  Function
	Path string

	VarsOrder []int
}

type Act struct {
	Run  Action
	Path string
}

type Formatter struct {
	Name string
	Var  bool

	Start, End int
}

func Map(fns []Function, acts []Action) (map[string]Fn, map[string]Act) {
	fn := MapFunctions(fns)
	ac := MapActions(acts)

	return fn, ac
}

func (f *File) Add(fns map[string]Fn, acts map[string]Act) {
	functions := []Fn{}
	for _, c := range f.Internal.Collections {
		fn := Fn{}

		for k, v := range f.Internal.Vars {
			if strings.HasPrefix(v.Value, fmt.Sprintf("%s.", c)) || strings.HasPrefix(v.Value, fmt.Sprintf("%s", c)) {
				fn = fns[c]
				fn.VarsOrder = append(fn.VarsOrder, k)
			}
		}

		functions = append(functions, fn)
	}

	for k := range f.Internal.Requests {
		f.Internal.Requests[k].Run = acts[f.Internal.Requests[k].Name].Run
	}

	f.Internal.Functions = functions
}

func MapFunctions(fns []Function) map[string]Fn {
	m := make(map[string]Fn)

	for _, v := range fns {
		path := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
		name := strings.ToLower(path[strings.LastIndex(path, ".")+1 : len(path)-2])

		m[name] = Fn{
			Run:  v,
			Path: path,
		}
	}

	return m
}

func MapActions(acts []Action) map[string]Act {
	m := make(map[string]Act)

	for _, v := range acts {
		path := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
		name := strings.ToLower(path[strings.LastIndex(path, ".")+1 : len(path)-3])

		m[name] = Act{
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
func (f *File) Parse() error {
	s := Structure{}

	if strings.Contains(string(f.Content), "<html>") {
		return fmt.Errorf("Template cannot contain html tags")
	}

	for k, v := range f.Content {
		if v == '<' {
			s.Start, s.End = 0, k-1
			break
		}
	}

	value := ""
	cols := []string{}

	if s.End >= 0 {
		s.Content = f.Content[s.End+1:]
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
	}

	vars := []Variable{}
	requests := []Request{}
	start, value, found := 0, "", false

	for k, v := range s.Content {
		if v == '{' {
			start = k
			found = true
		} else if v == '}' {
			if !strings.Contains(value, "(") && !strings.Contains(value, ":") {
				vars = append(vars, Variable{
					value,
					start,
					k,
				})
			}

			if strings.Contains(value, "/") {
				cols = append(cols, value[:strings.Index(value, "/")])
			} else if strings.Contains(value, "(") {
				ind := strings.Index(value, "(")
				requests = append(requests, Request{
					Name:   value[:strings.Index(value, "(")],
					Params: strings.Split(value[ind+1:], ","),
					Start:  start,
					End:    k + 1,
				})
			} else if strings.Contains(value, ":") {
				ind := strings.Index(value, ":")
				vars = append(vars, Variable{
					value[:ind],
					start,
					k,
				})
			} else if !strings.Contains(value, ".") {
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
	s.Requests = requests
	s.Formatted = string(s.Content)

	formatters, indices := createFormatters(s.Vars, s.Requests)

	s.Formatted = format(s.Formatted, indices)
	s.Formatters = formatters

	f.Internal = s

	return nil
}

// TODO add functionality
func ParseIndex(content []byte) (string, error) {
	str := string(content)
	start := strings.Index(str, "<body>")
	end := start + 6

	if start >= 0 {
		str = fmt.Sprintf("%s%s%s", str[:start], body, str[end:])
	} else {
		return str, fmt.Errorf("Body malformatted, needs a body")
	}

	return str, nil
}
func format(input string, indices [][]int) string {
	for _, index := range indices {
		start := index[0]
		end := index[1]

		if start < 0 || end >= len(input) || start > end {
			continue
		}

		replacement := "%v"
		input = input[:start] + replacement + input[end+1:]
	}

	return input
}

func createFormatters(vars []Variable, reqs []Request) ([]Formatter, [][]int) {
	sorted := make([]Formatter, 0, len(vars)+len(reqs))
	indices := [][]int{}

	for _, variable := range vars {
		indices = append(indices, []int{variable.Start, variable.End})
		sorted = append(sorted, Formatter{
			Name:  variable.Value,
			Var:   true,
			Start: variable.Start,
			End:   variable.End,
		})
	}

	for _, request := range reqs {
		indices = append(indices, []int{request.Start, request.End})
		sorted = append(sorted, Formatter{
			Name: request.Name,
			Var:  false,

			Start: request.Start,
			End:   request.End,
		})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Start != sorted[j].Start {
			return sorted[i].Start < sorted[j].Start
		}
		return sorted[i].End < sorted[j].End
	})

	sortIndexes(indices)

	return sorted, indices
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
