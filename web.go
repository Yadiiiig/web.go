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

func GenerateEndpoints(file *File) {
	http.HandleFunc(fmt.Sprintf("/%s", file.Name), func(w http.ResponseWriter, r *http.Request) {
		total := len(file.Internal.Functions) + len(file.Internal.Requests)
		results := make(chan map[string]interface{}, total)

		var wg sync.WaitGroup

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		for i := 0; i < len(file.Internal.Functions); i++ {
			wg.Add(1)
			go func(it int, w http.ResponseWriter, r *http.Request) {
				name, res, err := file.Internal.Functions[it].Run(w, r)
				if err != nil {
					log.Println(err)
				}

				result := make(map[string]interface{})
				result[name] = res

				results <- result
				wg.Done()
			}(i, w, r)
		}

		for i := 0; i < len(file.Internal.Requests); i++ {
			wg.Add(1)
			go func(it int, w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					w.WriteHeader(http.StatusMethodNotAllowed)
					fmt.Fprintf(w, "Invalid request method")
					return
				}

				var data map[string]interface{}

				err := json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Failed to parse JSON data: %v", err)
					return
				}

				name, res, err := file.Internal.Requests[it].Run(data, w, r)
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

		tmp := make(map[string]interface{})
		for r := range results {
			for k, v := range r {
				mp := FieldMapping(k, v)
				AppendMap(tmp, mp)
			}
		}

		fmt.Println(tmp)

		json, err := json.Marshal(tmp)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(json)
		if err != nil {
			log.Fatal(err)
		}
	})
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
