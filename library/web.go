package library

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

func GenerateEndpoints(file *File) {
	http.HandleFunc(fmt.Sprintf("/%s", file.Name), func(w http.ResponseWriter, r *http.Request) {
		amount := len(file.Internal.Functions)
		results := make(chan map[string]interface{}, amount)

		var wg sync.WaitGroup
		wg.Add(amount)

		for i := 0; i < amount; i++ {
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

		fm := []interface{}{}
		for _, v := range file.Internal.Vars {
			if output, ok := tmp[v.Value]; ok {
				fm = append(fm, output)
			} else {
				fm = append(fm, "ERROR")
			}
		}

		snippet := fmt.Sprintf(file.Internal.Formatted, fm[:]...)

		fmt.Fprintf(w, snippet)
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
