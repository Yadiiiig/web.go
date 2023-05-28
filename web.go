package library

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

func GenerateEndpoints(file *File) error {
	var err error

	err = file.FunctionHandlers()
	if err != nil {
		return err
	}

	err = file.ActionHandlers()
	if err != nil {
		return err
	}

	return nil
}

func (f *File) FunctionHandlers() error {
	fmt.Println("Generated (action) endpoint: ", fmt.Sprintf("/%s", f.Name))

	http.HandleFunc(fmt.Sprintf("/%s", f.Name), func(w http.ResponseWriter, r *http.Request) {
		total := len(f.Internal.Functions)
		results := make(chan map[string]interface{}, total)

		var wg sync.WaitGroup

		w, end := AddHeaders(w, r.Method)
		if end {
			return
		}

		for i := 0; i < len(f.Internal.Functions); i++ {
			wg.Add(1)
			go func(it int, w http.ResponseWriter, r *http.Request) {
				name, res, err := f.Internal.Functions[it].Run(w, r)
				if err != nil {
					log.Println(err)
				}

				result := make(map[string]interface{})
				result[name] = res

				results <- result
				wg.Done()
			}(i, w, r)
		}

		wg.Wait()
		close(results)

		err := HandleResult(w, results)
		if err != nil {
			log.Println(err)
		}
	})
	return nil
}

func (f *File) ActionHandlers() error {
	for i := 0; i < len(f.Internal.Requests); i++ {
		fmt.Println("Generated (action) endpoint: ", fmt.Sprintf("/%s/%s", f.Name, f.Internal.Requests[i].Name))

		http.HandleFunc(fmt.Sprintf("/%s/%s", f.Name, f.Internal.Requests[i].Name), func(w http.ResponseWriter, r *http.Request) {
			w, end := AddHeaders(w, r.Method)
			if end {
				return
			}

			var data map[string]interface{}

			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Failed to parse JSON data: %v", err)
				return
			}

			name, res, err := f.Internal.Requests[i-1].Run(data, w, r)
			if err != nil {
				log.Println(err)
			}

			err = HandleAction(w, name, res)
			if err != nil {
				log.Println(err)
			}

		})
	}

	return nil
}

func HandleResult(w http.ResponseWriter, results chan map[string]interface{}) error {
	tmp := make(map[string]interface{})

	if results == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	for r := range results {
		for k, v := range r {
			mp := FieldMapping(k, v)
			AppendMap(tmp, mp)
		}
	}

	json, err := json.Marshal(tmp)
	if err != nil {
		return err

	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func HandleAction(w http.ResponseWriter, name string, res interface{}) error {
	tmp := FieldMapping(name, res)

	json, err := json.Marshal(tmp)
	if err != nil {
		return err

	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func FieldMapping(name string, data interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	t := reflect.TypeOf(data)

	if t.Kind() != reflect.Struct {
		fields[name] = data
		return fields
	}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("web")

		fieldPath := fieldType.Name
		if tag != "" {
			fieldPath = strings.ToLower(tag)
		}
		if name != "" {
			fieldPath = name + "." + fieldPath
		}

		if fieldType.Type.Kind() == reflect.Struct {
			subFields := FieldMapping(fieldPath, field.Interface())
			for subKey, value := range subFields {
				key := subKey
				if fieldPath != "" {
					key = fieldPath + "." + subKey
				}
				fields[key] = value
			}
		} else {
			fields[fieldPath] = field.Interface()
		}
	}

	return fields
}

func AppendMap(dest, src map[string]interface{}) {
	for key, value := range src {
		dest[key] = value
	}
}

func AddHeaders(w http.ResponseWriter, method string) (http.ResponseWriter, bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

	if method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)

		return w, true
	}

	return w, false
}
